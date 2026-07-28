package githubrelease

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func serverReturning(t *testing.T, status int, body string) (Client, func()) {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/repos/owner/repo/releases" {
			t.Errorf("unexpected path: %s (Latest must use the list endpoint, not /releases/latest, so prereleases are visible)", r.URL.Path)
		}
		w.WriteHeader(status)
		w.Write([]byte(body))
	}))
	return Client{BaseURL: server.URL, HTTP: server.Client()}, server.Close
}

// samplePayload mirrors the /releases list endpoint: an array, newest
// first, which is how GitHub actually orders it.
const samplePayload = `[{
  "tag_name": "v0.2.0",
  "assets": [
    {"name": "hawp-darwin-arm64", "browser_download_url": "https://example.test/hawp-darwin-arm64", "digest": "sha256:abc123", "size": 1000},
    {"name": "hawp-linux-amd64", "browser_download_url": "https://example.test/hawp-linux-amd64", "digest": "", "size": 900}
  ]
}]`

func TestLatestParsesReleaseAndAssets(t *testing.T) {
	client, closeServer := serverReturning(t, http.StatusOK, samplePayload)
	defer closeServer()

	release, err := client.Latest("owner/repo")
	if err != nil {
		t.Fatal(err)
	}
	if release.TagName != "v0.2.0" {
		t.Errorf("TagName = %q", release.TagName)
	}
	asset, ok := release.AssetNamed("hawp-darwin-arm64")
	if !ok {
		t.Fatal("expected to find hawp-darwin-arm64 asset")
	}
	if asset.SHA256() != "abc123" {
		t.Errorf("SHA256() = %q, want abc123", asset.SHA256())
	}
	if _, ok := release.AssetNamed("does-not-exist"); ok {
		t.Error("expected AssetNamed to report not-found for unknown asset")
	}
}

func TestAssetSHA256EmptyWhenDigestMissing(t *testing.T) {
	client, closeServer := serverReturning(t, http.StatusOK, samplePayload)
	defer closeServer()

	release, _ := client.Latest("owner/repo")
	asset, _ := release.AssetNamed("hawp-linux-amd64")
	if got := asset.SHA256(); got != "" {
		t.Errorf("SHA256() = %q, want empty for asset with no digest yet", got)
	}
}

// TestLatestIncludesPrereleases is a regression test for the real bug
// found during 4c152ee3's live release test: GitHub's /releases/latest
// excludes prereleases by definition, so a repo whose only releases are
// all prerelease-tagged (as 0.0.x test releases intentionally are) would
// otherwise be permanently invisible to `hawp update`.
func TestLatestIncludesPrereleases(t *testing.T) {
	client, closeServer := serverReturning(t, http.StatusOK, `[{
		"tag_name": "v0.0.2",
		"prerelease": true,
		"assets": [{"name": "hawp-linux-amd64", "browser_download_url": "https://example.test/x", "digest": "sha256:def456", "size": 100}]
	}]`)
	defer closeServer()

	release, err := client.Latest("owner/repo")
	if err != nil {
		t.Fatal(err)
	}
	if release.TagName != "v0.0.2" {
		t.Errorf("TagName = %q, want v0.0.2 (prerelease must still be found)", release.TagName)
	}
}

func TestLatestNoReleasesYet(t *testing.T) {
	client, closeServer := serverReturning(t, http.StatusOK, `[]`)
	defer closeServer()

	_, err := client.Latest("owner/repo")
	if err != ErrNoReleases {
		t.Errorf("err = %v, want ErrNoReleases for an empty releases list", err)
	}
}

func TestLatest404TreatedAsNoReleases(t *testing.T) {
	client, closeServer := serverReturning(t, http.StatusNotFound, `{"message":"Not Found"}`)
	defer closeServer()

	_, err := client.Latest("owner/repo")
	if err != ErrNoReleases {
		t.Errorf("err = %v, want ErrNoReleases", err)
	}
}

func TestLatestUnexpectedStatus(t *testing.T) {
	client, closeServer := serverReturning(t, http.StatusInternalServerError, `{}`)
	defer closeServer()

	if _, err := client.Latest("owner/repo"); err == nil {
		t.Fatal("expected error for 500 response")
	}
}
