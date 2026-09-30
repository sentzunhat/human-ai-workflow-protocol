// Package githubrelease reads release metadata from the GitHub REST API.
// GitHub computes a SHA-256 "digest" for every uploaded release asset
// automatically (confirmed against microsoft/onnxruntime's public
// releases), so any published release — regardless of whether its author
// hand-computed checksums — can be verified through this one field.
package githubrelease

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// Asset is one file attached to a release.
type Asset struct {
	Name        string `json:"name"`
	DownloadURL string `json:"browser_download_url"`
	Digest      string `json:"digest"` // "sha256:<hex>", empty if not yet computed
	Size        int64  `json:"size"`
}

// Release is the subset of the GitHub release payload this client needs.
type Release struct {
	TagName string  `json:"tag_name"`
	Assets  []Asset `json:"assets"`
}

// SHA256 extracts the hex digest from Digest's "sha256:<hex>" form, or ""
// if the asset has no computed digest yet (GitHub computes these shortly
// after upload; a brand-new release may briefly have none).
func (a Asset) SHA256() string {
	hex, ok := strings.CutPrefix(a.Digest, "sha256:")
	if !ok {
		return ""
	}
	return hex
}

// AssetNamed returns the asset matching name, or ok=false.
func (r Release) AssetNamed(name string) (Asset, bool) {
	for _, asset := range r.Assets {
		if asset.Name == name {
			return asset, true
		}
	}
	return Asset{}, false
}

// Client fetches release metadata. BaseURL defaults to the real GitHub
// API; tests override it to a local httptest server.
type Client struct {
	BaseURL string
	HTTP    *http.Client
}

// NewClient returns a Client pointed at the real GitHub API.
func NewClient() Client {
	return Client{BaseURL: "https://api.github.com", HTTP: &http.Client{Timeout: 30 * time.Second}}
}

// Latest fetches the most recently published release for "owner/repo",
// including prereleases.
//
// This deliberately does NOT use GitHub's /releases/latest endpoint:
// that endpoint excludes prereleases by definition, which would make
// hawp update blind to prerelease-tagged versions (verified against the
// real API during 4c152ee3's v0.0.1/v0.0.2 test releases, both
// intentionally marked prerelease — /releases/latest 404'd for both).
// /releases (the list endpoint) returns every non-draft release newest
// first, so the first entry is simply the newest.
//
// Draft releases never appear here regardless of endpoint — GitHub
// hides drafts from unauthenticated requests (this Client sends none),
// which is the entire mechanism behind the release workflow's optional
// `draft` input (see b95436f2): a draft build is invisible to every
// hawp update check until a maintainer publishes it. Confirmed
// empirically on this repo with a real draft release.
func (c Client) Latest(repo string) (Release, error) {
	client := c.HTTP
	if client == nil {
		client = http.DefaultClient
	}
	url := fmt.Sprintf("%s/repos/%s/releases?per_page=1", c.BaseURL, repo)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return Release{}, err
	}
	req.Header.Set("Accept", "application/vnd.github+json")

	resp, err := client.Do(req)
	if err != nil {
		return Release{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return Release{}, ErrNoReleases
	}
	if resp.StatusCode != http.StatusOK {
		return Release{}, fmt.Errorf("github releases API: unexpected status %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Release{}, err
	}
	var releases []Release
	if err := json.Unmarshal(body, &releases); err != nil {
		return Release{}, err
	}
	if len(releases) == 0 {
		return Release{}, ErrNoReleases
	}
	return releases[0], nil
}

// ErrNoReleases means the repo has not published any releases yet.
var ErrNoReleases = fmt.Errorf("no published releases found")
