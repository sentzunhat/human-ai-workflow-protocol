// Package provision orchestrates `hawp init`: create the ~/.hawp/ layout,
// download and verify the ONNX Runtime and embedding model, and record what
// was installed in a manifest so re-running is idempotent.
package provision

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	domainprovision "github.com/sentzunhat/hawp/librarian/src/internal/domain/provision"
	"github.com/sentzunhat/hawp/librarian/src/internal/infrastructure/archive"
	"github.com/sentzunhat/hawp/librarian/src/internal/infrastructure/download"
	"github.com/sentzunhat/hawp/librarian/src/internal/infrastructure/filesystem"
)

// Step describes the outcome of provisioning one asset.
type Step struct {
	Name   string
	Status string // "already-installed" | "downloaded" | "failed"
	Path   string
	Err    error
}

// Result is the full init run outcome.
type Result struct {
	Home  filesystem.HawpHome
	Steps []Step
}

// Registry is the set of assets to provision, injected so callers (tests,
// alternate registries) do not have to go through the global default.
type Registry struct {
	RuntimeAsset    domainprovision.Asset
	RuntimeAssetErr error // set when the current platform has no prebuilt asset
	ModelAssets     []domainprovision.Asset
	RuntimeVersion  string
	ModelName       string
	ModelVersion    string
}

// DefaultRegistry resolves the real ONNX Runtime + embedding model assets
// for the current platform (see internal/domain/provision).
func DefaultRegistry() Registry {
	asset, err := domainprovision.RuntimeAsset()
	return Registry{
		RuntimeAsset:    asset,
		RuntimeAssetErr: err,
		ModelAssets:     domainprovision.ModelAssets,
		RuntimeVersion:  domainprovision.RuntimeVersion,
		ModelName:       domainprovision.ModelName,
		ModelVersion:    domainprovision.ModelVersion,
	}
}

// Run provisions the ~/.hawp/ layout for home using the given registry.
// The runtime and model assets are provisioned independently, so a failure
// in one does not block the other (fail soft, report what's missing).
func Run(fetcher download.Fetcher, home string, registry Registry) Result {
	hawpHome := filesystem.ResolveHawpHome(home)
	result := Result{Home: hawpHome}

	for _, dir := range hawpHome.Dirs() {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			result.Steps = append(result.Steps, Step{Name: "layout:" + dir, Status: "failed", Err: err})
			return result
		}
	}

	manifest, err := domainprovision.LoadManifest(hawpHome.Root)
	if err != nil {
		result.Steps = append(result.Steps, Step{Name: "manifest", Status: "failed", Err: err})
		return result
	}

	if registry.RuntimeAssetErr != nil {
		result.Steps = append(result.Steps, Step{Name: "onnxruntime", Status: "failed", Err: registry.RuntimeAssetErr})
	} else {
		result.Steps = append(result.Steps, provisionRuntimeAsset(fetcher, hawpHome, manifest, registry.RuntimeAsset))
	}

	for _, asset := range registry.ModelAssets {
		result.Steps = append(result.Steps, provisionPlainAsset(fetcher, hawpHome.Models, manifest, asset))
	}

	manifest.RuntimeVersion = registry.RuntimeVersion
	manifest.ModelName = registry.ModelName
	manifest.ModelVersion = registry.ModelVersion
	if err := manifest.Save(hawpHome.Root); err != nil {
		result.Steps = append(result.Steps, Step{Name: "manifest", Status: "failed", Err: err})
	}

	return result
}

// provisionRuntimeAsset downloads the ONNX Runtime archive (verified) and
// extracts just the shared library into hawpHome.Runtime.
func provisionRuntimeAsset(fetcher download.Fetcher, hawpHome filesystem.HawpHome, manifest *domainprovision.Manifest, asset domainprovision.Asset) Step {
	destPath := filepath.Join(hawpHome.Runtime, asset.DestName)

	if manifest.Satisfies(asset) {
		if hash, err := download.FileSHA256(destPath); err == nil && verifiesExtractedLib(hash, destPath) {
			return Step{Name: asset.Name, Status: "already-installed", Path: destPath}
		}
	}

	archivePath := filepath.Join(hawpHome.Downloads, filepath.Base(asset.URL))
	if err := download.VerifiedFile(fetcher, asset.URL, asset.SHA256, archivePath); err != nil {
		return Step{Name: asset.Name, Status: "failed", Err: err}
	}
	if err := archive.ExtractMember(archivePath, asset.ArchiveMember, destPath); err != nil {
		return Step{Name: asset.Name, Status: "failed", Err: err}
	}

	manifest.Assets[asset.Name] = domainprovision.AssetRecord{
		SHA256: asset.SHA256, Size: asset.Size, InstalledAt: time.Now().UTC().Format(time.RFC3339),
	}
	return Step{Name: asset.Name, Status: "downloaded", Path: destPath}
}

// verifiesExtractedLib is a light sanity check: the extracted lib exists
// and is non-empty. The manifest hash covers the source archive, not the
// extracted member, so this just guards against a half-finished extract.
func verifiesExtractedLib(_ string, path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.Size() > 0
}

// provisionPlainAsset downloads a non-archive asset directly into destDir,
// verified against its recorded SHA-256, skipping the download when the
// manifest and on-disk hash both already match (no network needed).
func provisionPlainAsset(fetcher download.Fetcher, destDir string, manifest *domainprovision.Manifest, asset domainprovision.Asset) Step {
	destPath := filepath.Join(destDir, asset.DestName)

	if manifest.Satisfies(asset) {
		if hash, err := download.FileSHA256(destPath); err == nil && hash == asset.SHA256 {
			return Step{Name: asset.Name, Status: "already-installed", Path: destPath}
		}
	}

	if err := download.VerifiedFile(fetcher, asset.URL, asset.SHA256, destPath); err != nil {
		return Step{Name: asset.Name, Status: "failed", Err: err}
	}

	manifest.Assets[asset.Name] = domainprovision.AssetRecord{
		SHA256: asset.SHA256, Size: asset.Size, InstalledAt: time.Now().UTC().Format(time.RFC3339),
	}
	return Step{Name: asset.Name, Status: "downloaded", Path: destPath}
}

// String renders a human-readable report of the run.
func (r Result) String() string {
	lines := "hawp init\n=========\n"
	lines += "home: " + r.Home.Root + "\n\n"
	failed := 0
	for _, step := range r.Steps {
		switch step.Status {
		case "already-installed":
			lines += fmt.Sprintf("✓ %s: already installed (%s)\n", step.Name, step.Path)
		case "downloaded":
			lines += fmt.Sprintf("✓ %s: downloaded and verified (%s)\n", step.Name, step.Path)
		case "failed":
			failed++
			lines += fmt.Sprintf("✗ %s: %v\n", step.Name, step.Err)
		}
	}
	lines += "\n"
	if failed > 0 {
		lines += fmt.Sprintf("%d of %d asset(s) failed. Lexical kit/work commands are unaffected; vector search and context building need these assets.\n", failed, len(r.Steps))
	} else {
		lines += fmt.Sprintf("All %d asset(s) ready.\n", len(r.Steps))
	}
	return lines
}

// Failed reports whether any step failed.
func (r Result) Failed() bool {
	for _, step := range r.Steps {
		if step.Status == "failed" {
			return true
		}
	}
	return false
}
