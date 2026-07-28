// Package update is the application service for `hawp update` and
// `hawp version`.
package update

import (
	"fmt"
	"os"
	"path/filepath"

	domainupdate "github.com/sentzunhat/hawp/librarian/go/internal/domain/update"
	"github.com/sentzunhat/hawp/librarian/go/internal/infrastructure/download"
	"github.com/sentzunhat/hawp/librarian/go/internal/infrastructure/githubrelease"
	"github.com/sentzunhat/hawp/librarian/go/internal/infrastructure/selfreplace"
)

// Status is the result of checking for an update, without applying one.
type Status struct {
	Current         string
	Latest          string
	UpdateAvailable bool
	AssetName       string
	NoReleases      bool
}

// Check compares the running version against the latest published release
// for the current platform. NoReleases is set (not an error) when the repo
// has not published any releases yet.
func Check(client githubrelease.Client, repo, currentVersion string) (Status, error) {
	assetName, err := domainupdate.AssetName()
	if err != nil {
		return Status{}, err
	}

	release, err := client.Latest(repo)
	if err == githubrelease.ErrNoReleases {
		return Status{Current: currentVersion, NoReleases: true, AssetName: assetName}, nil
	}
	if err != nil {
		return Status{}, err
	}

	return Status{
		Current:         currentVersion,
		Latest:          domainupdate.CleanVersion(release.TagName),
		UpdateAvailable: domainupdate.IsNewer(currentVersion, release.TagName),
		AssetName:       assetName,
	}, nil
}

// Apply downloads, verifies, and installs the latest release's platform
// asset over execPath (typically os.Executable()). Returns the applied
// version, or an error — the running binary is left untouched on any
// failure (verified download, then atomic same-directory rename).
func Apply(fetcher download.Fetcher, client githubrelease.Client, repo, execPath string) (string, error) {
	assetName, err := domainupdate.AssetName()
	if err != nil {
		return "", err
	}

	release, err := client.Latest(repo)
	if err != nil {
		return "", err
	}

	asset, ok := release.AssetNamed(assetName)
	if !ok {
		return "", fmt.Errorf("release %s has no asset named %q", release.TagName, assetName)
	}
	sha256 := asset.SHA256()
	if sha256 == "" {
		return "", fmt.Errorf("release %s asset %q has no computed checksum yet; try again shortly", release.TagName, assetName)
	}

	stagingPath := filepath.Join(filepath.Dir(execPath), ".hawp-update-staged")
	if err := download.VerifiedFile(fetcher, asset.DownloadURL, sha256, stagingPath); err != nil {
		return "", err
	}
	defer os.Remove(stagingPath) // no-op once Replace renames it away

	if err := selfreplace.Replace(stagingPath, execPath); err != nil {
		return "", err
	}
	return domainupdate.CleanVersion(release.TagName), nil
}
