package update

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	domainupdate "github.com/sentzunhat/hawp/librarian/src/internal/domain/update"
	"github.com/sentzunhat/hawp/librarian/src/internal/infrastructure/download"
	"github.com/sentzunhat/hawp/librarian/src/internal/infrastructure/githubrelease"
)

func hashHex(b []byte) string {
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:])
}

// releaseServer serves /repos/owner/repo/releases (the list endpoint —
// Client.Latest deliberately avoids /releases/latest, which excludes
// prereleases) and a binary download route.
func releaseServer(t *testing.T, tag string, binary []byte, digest string) (githubrelease.Client, string) {
	t.Helper()
	assetName, err := domainupdate.AssetName()
	if err != nil {
		t.Skipf("platform not supported: %v", err)
	}

	var server *httptest.Server
	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/repos/owner/repo/releases":
			fmt.Fprintf(w, `[{"tag_name":%q,"assets":[{"name":%q,"browser_download_url":%q,"digest":%q,"size":%d}]}]`,
				tag, assetName, server.URL+"/download/"+assetName, digest, len(binary))
		case "/download/" + assetName:
			w.Write(binary)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	t.Cleanup(server.Close)
	return githubrelease.Client{BaseURL: server.URL, HTTP: server.Client()}, server.URL
}

func TestCheckReportsUpdateAvailable(t *testing.T) {
	client, _ := releaseServer(t, "v1.0.0", []byte("binary"), "sha256:"+hashHex([]byte("binary")))

	status, err := Check(client, "owner/repo", "v0.9.0")
	if err != nil {
		t.Fatal(err)
	}
	if !status.UpdateAvailable || status.Latest != "v1.0.0" {
		t.Errorf("status = %+v, want update available to v1.0.0", status)
	}
}

func TestCheckReportsUpToDate(t *testing.T) {
	client, _ := releaseServer(t, "v1.0.0", []byte("binary"), "sha256:"+hashHex([]byte("binary")))

	status, err := Check(client, "owner/repo", "v1.0.0")
	if err != nil {
		t.Fatal(err)
	}
	if status.UpdateAvailable {
		t.Error("expected no update available when already on latest")
	}
}

func TestCheckNoReleasesYet(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()
	client := githubrelease.Client{BaseURL: server.URL, HTTP: server.Client()}

	status, err := Check(client, "owner/repo", "dev")
	if err != nil {
		t.Fatal(err)
	}
	if !status.NoReleases {
		t.Error("expected NoReleases when the repo has no published releases")
	}
}

func TestApplyDownloadsVerifiesAndReplaces(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("self-replace is not supported on Windows")
	}
	newBinary := []byte("new binary contents")
	client, _ := releaseServer(t, "v1.0.0", newBinary, "sha256:"+hashHex(newBinary))

	dir := t.TempDir()
	execPath := filepath.Join(dir, "hawp")
	if err := os.WriteFile(execPath, []byte("old binary"), 0o755); err != nil {
		t.Fatal(err)
	}

	applied, err := Apply(download.NewHTTPFetcher(), client, "owner/repo", execPath)
	if err != nil {
		t.Fatal(err)
	}
	if applied != "v1.0.0" {
		t.Errorf("applied = %q, want v1.0.0", applied)
	}

	content, err := os.ReadFile(execPath)
	if err != nil {
		t.Fatal(err)
	}
	if string(content) != string(newBinary) {
		t.Errorf("execPath content = %q, want %q", content, newBinary)
	}

	entries, _ := os.ReadDir(dir)
	for _, e := range entries {
		if e.Name() != "hawp" {
			t.Errorf("staging file leaked: %s", e.Name())
		}
	}
}

func TestApplyChecksumMismatchLeavesOriginalBinaryIntact(t *testing.T) {
	wrongHash := "sha256:" + hashHex([]byte("something else entirely"))
	client, _ := releaseServer(t, "v1.0.0", []byte("new binary"), wrongHash)

	dir := t.TempDir()
	execPath := filepath.Join(dir, "hawp")
	original := []byte("old binary — must survive")
	if err := os.WriteFile(execPath, original, 0o755); err != nil {
		t.Fatal(err)
	}

	if _, err := Apply(download.NewHTTPFetcher(), client, "owner/repo", execPath); err == nil {
		t.Fatal("expected checksum mismatch error")
	}

	content, err := os.ReadFile(execPath)
	if err != nil {
		t.Fatal(err)
	}
	if string(content) != string(original) {
		t.Fatalf("original binary was modified after a failed update: %q", content)
	}
}

func TestApplyMissingAssetForRelease(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `[{"tag_name":"v1.0.0","assets":[]}]`)
	}))
	defer server.Close()
	client := githubrelease.Client{BaseURL: server.URL, HTTP: server.Client()}

	dir := t.TempDir()
	execPath := filepath.Join(dir, "hawp")
	os.WriteFile(execPath, []byte("old"), 0o755)

	if _, err := Apply(download.NewHTTPFetcher(), client, "owner/repo", execPath); err == nil {
		t.Fatal("expected error when release has no matching asset")
	}
}
