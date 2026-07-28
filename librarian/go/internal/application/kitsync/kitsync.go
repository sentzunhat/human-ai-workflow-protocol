// Package kitsync is the application service behind `hawp update`'s
// kit + provider-overlay sync: download the release's kit bundle,
// refresh .hawp/kit/, and refresh the detected (or explicitly named)
// provider's update:refresh files, leaving update:skip files (like a
// customized CLAUDE.md) untouched.
package kitsync

import (
	"fmt"
	"os"
	"path/filepath"

	domainkitsync "github.com/sentzunhat/hawp/librarian/go/internal/domain/kitsync"
	"github.com/sentzunhat/hawp/librarian/go/internal/infrastructure/archive"
	"github.com/sentzunhat/hawp/librarian/go/internal/infrastructure/download"
	"github.com/sentzunhat/hawp/librarian/go/internal/infrastructure/githubrelease"
)

// BundleAssetName is the release asset that carries the kit + provider
// bundle, alongside the platform binaries.
const BundleAssetName = "hawp-kit-bundle.tar.gz"

// Result summarizes a sync pass.
type Result struct {
	KitFilesWritten  int
	Providers        map[string]int  // provider name -> files written
	ProviderInstalls map[string]bool // true when freshly installed (not updated)
	SkippedRules     []string
	NoBundleAsset    bool // release has no kit bundle yet (older release)
}

// Sync downloads and verifies the release's kit bundle, refreshes
// repoRoot/.hawp/kit/, and refreshes each provider in providerNames (or,
// if empty, every auto-detected provider). Returns Result{NoBundleAsset:
// true} (not an error) when the release predates this feature.
func Sync(fetcher download.Fetcher, client githubrelease.Client, repo, repoRoot string, providerNames []string) (Result, error) {
	release, err := client.Latest(repo)
	if err != nil {
		return Result{}, err
	}
	asset, ok := release.AssetNamed(BundleAssetName)
	if !ok {
		return Result{NoBundleAsset: true}, nil
	}
	sha256 := asset.SHA256()
	if sha256 == "" {
		return Result{}, fmt.Errorf("release %s asset %q has no computed checksum yet; try again shortly", release.TagName, BundleAssetName)
	}

	tempDir, err := os.MkdirTemp("", "hawp-kit-bundle-*")
	if err != nil {
		return Result{}, err
	}
	defer os.RemoveAll(tempDir)

	archivePath := filepath.Join(tempDir, BundleAssetName)
	if err := download.VerifiedFile(fetcher, asset.DownloadURL, sha256, archivePath); err != nil {
		return Result{}, err
	}

	bundleRoot := filepath.Join(tempDir, "extracted")
	if err := archive.ExtractAll(archivePath, bundleRoot); err != nil {
		return Result{}, err
	}

	manifest, err := domainkitsync.ParseManifest(filepath.Join(bundleRoot, "providers", "manifest.yaml"))
	if err != nil {
		return Result{}, fmt.Errorf("parse bundled manifest: %w", err)
	}

	result := Result{Providers: map[string]int{}, ProviderInstalls: map[string]bool{}}
	result.KitFilesWritten, err = domainkitsync.SyncKit(filepath.Join(bundleRoot, "kit"), repoRoot)
	if err != nil {
		return result, fmt.Errorf("sync kit: %w", err)
	}

	// Resolve the target provider list.
	targets := resolveTargets(providerNames, manifest, repoRoot)

	// Determine which are already installed so we can route install vs update.
	installed := map[string]bool{}
	for _, name := range domainkitsync.DetectProviders(repoRoot, manifest) {
		installed[name] = true
	}

	for _, name := range targets {
		if installed[name] {
			written, skipped, err := domainkitsync.ApplyProviderUpdate(bundleRoot, repoRoot, manifest, name)
			if err != nil {
				return result, fmt.Errorf("update provider %s: %w", name, err)
			}
			result.Providers[name] = written
			result.ProviderInstalls[name] = false
			result.SkippedRules = append(result.SkippedRules, skipped...)
		} else {
			written, seeded, err := domainkitsync.ApplyProviderInstall(bundleRoot, repoRoot, manifest, name)
			if err != nil {
				return result, fmt.Errorf("install provider %s: %w", name, err)
			}
			result.Providers[name] = written
			result.ProviderInstalls[name] = true
			result.SkippedRules = append(result.SkippedRules, seeded...)
		}
	}

	return result, nil
}

// resolveTargets expands providerNames:
//   - empty → no providers (kit-only sync)
//   - contains "all" → every provider in manifest
//   - otherwise → use the list as-is
func resolveTargets(providerNames []string, manifest *domainkitsync.Manifest, _ string) []string {
	for _, name := range providerNames {
		if name == "all" {
			return manifest.AllProviderNames()
		}
	}
	return providerNames
}
