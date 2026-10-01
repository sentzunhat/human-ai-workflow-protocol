package download

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func serve(t *testing.T, body []byte) (Fetcher, string) {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(body)
	}))
	t.Cleanup(server.Close)
	return NewHTTPFetcher(), server.URL
}

func hashOf(b []byte) string {
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:])
}

func TestVerifiedFileSuccess(t *testing.T) {
	body := []byte("hello ONNX runtime")
	fetcher, url := serve(t, body)
	dest := filepath.Join(t.TempDir(), "asset.bin")

	if err := VerifiedFile(fetcher, url, hashOf(body), dest); err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(dest)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != string(body) {
		t.Errorf("content mismatch: %q", got)
	}
}

func TestVerifiedFileChecksumMismatch(t *testing.T) {
	fetcher, url := serve(t, []byte("actual content"))
	dest := filepath.Join(t.TempDir(), "asset.bin")

	err := VerifiedFile(fetcher, url, hashOf([]byte("expected content")), dest)
	var mismatch *ErrChecksumMismatch
	if err == nil {
		t.Fatal("expected checksum mismatch error")
	}
	if !asMismatch(err, &mismatch) {
		t.Fatalf("wrong error type: %v", err)
	}
	if _, statErr := os.Stat(dest); !os.IsNotExist(statErr) {
		t.Error("dest file must not exist after a checksum mismatch")
	}
	entries, _ := os.ReadDir(filepath.Dir(dest))
	for _, e := range entries {
		t.Errorf("temp file leaked: %s", e.Name())
	}
}

func asMismatch(err error, target **ErrChecksumMismatch) bool {
	if m, ok := err.(*ErrChecksumMismatch); ok {
		*target = m
		return true
	}
	return false
}

func TestVerifiedFileHTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	dest := filepath.Join(t.TempDir(), "asset.bin")
	if err := VerifiedFile(NewHTTPFetcher(), server.URL, "irrelevant", dest); err == nil {
		t.Fatal("expected error for 404 response")
	}
}

func TestFileSHA256(t *testing.T) {
	body := []byte("some file content")
	path := filepath.Join(t.TempDir(), "f.txt")
	if err := os.WriteFile(path, body, 0o644); err != nil {
		t.Fatal(err)
	}
	got, err := FileSHA256(path)
	if err != nil {
		t.Fatal(err)
	}
	if got != hashOf(body) {
		t.Errorf("FileSHA256 = %s, want %s", got, hashOf(body))
	}
}

// Fetch is an external callback between the initial guard and temp creation.
type substitutingFetcher struct {
	t               *testing.T
	parent, outside string
}

func (f substitutingFetcher) Fetch(string) (io.ReadCloser, error) {
	f.t.Helper()
	if err := os.Remove(f.parent); err != nil {
		f.t.Fatal(err)
	}
	if err := os.Symlink(f.outside, f.parent); err != nil {
		f.t.Skipf("symlinks unavailable: %v", err)
	}
	return io.NopCloser(strings.NewReader("payload")), nil
}
func TestVerifiedFileRejectsParentSubstitutedDuringFetch(t *testing.T) {
	root := t.TempDir()
	parent := filepath.Join(root, "download")
	outside := t.TempDir()
	err := VerifiedFile(substitutingFetcher{t, parent, outside}, "unused", hashOf([]byte("payload")), filepath.Join(parent, "asset.bin"))
	if err == nil || !strings.Contains(err.Error(), "unsafe download destination") {
		t.Fatalf("expected unsafe destination, got %v", err)
	}
	entries, err := os.ReadDir(outside)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 0 {
		t.Fatalf("download wrote outside destination: %v", entries)
	}
}
