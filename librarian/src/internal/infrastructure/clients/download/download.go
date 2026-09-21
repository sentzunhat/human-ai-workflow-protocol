// Package download fetches a URL, verifies its SHA-256 against an expected
// hash while streaming (never trusting unverified bytes to disk), and
// writes the result atomically.
package download

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// Fetcher retrieves a URL's body. Swappable in tests so no real network
// call is made.
type Fetcher interface {
	Fetch(url string) (io.ReadCloser, error)
}

// HTTPFetcher is the default Fetcher, backed by net/http with a timeout.
type HTTPFetcher struct {
	Client *http.Client
}

// NewHTTPFetcher returns an HTTPFetcher with a sane default timeout for
// large binary assets.
func NewHTTPFetcher() HTTPFetcher {
	return HTTPFetcher{Client: &http.Client{Timeout: 10 * time.Minute}}
}

func (f HTTPFetcher) Fetch(url string) (io.ReadCloser, error) {
	client := f.Client
	if client == nil {
		client = http.DefaultClient
	}
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("download %s: unexpected status %s", url, resp.Status)
	}
	return resp.Body, nil
}

// ErrChecksumMismatch means the downloaded content did not match the
// expected SHA-256.
type ErrChecksumMismatch struct {
	URL      string
	Expected string
	Actual   string
}

func (e *ErrChecksumMismatch) Error() string {
	return fmt.Sprintf("checksum mismatch for %s: expected %s, got %s", e.URL, e.Expected, e.Actual)
}

// VerifiedFile downloads url via fetcher, verifying the stream's SHA-256
// against expectedSHA256 (lowercase hex) as it writes to a temp file
// alongside destPath, then atomically renames into place. On mismatch the
// temp file is removed and no partial content ever lands at destPath.
func VerifiedFile(fetcher Fetcher, url, expectedSHA256, destPath string) error {
	if err := os.MkdirAll(filepath.Dir(destPath), 0o755); err != nil {
		return err
	}

	body, err := fetcher.Fetch(url)
	if err != nil {
		return err
	}
	defer body.Close()

	temp, err := os.CreateTemp(filepath.Dir(destPath), ".download-*")
	if err != nil {
		return err
	}
	tempPath := temp.Name()
	defer os.Remove(tempPath) // no-op once renamed

	hasher := sha256.New()
	if _, err := io.Copy(io.MultiWriter(temp, hasher), body); err != nil {
		temp.Close()
		return err
	}
	if err := temp.Close(); err != nil {
		return err
	}

	actual := hex.EncodeToString(hasher.Sum(nil))
	if actual != expectedSHA256 {
		return &ErrChecksumMismatch{URL: url, Expected: expectedSHA256, Actual: actual}
	}

	return os.Rename(tempPath, destPath)
}

// FileSHA256 computes the SHA-256 of an existing file, for idempotency
// checks before re-downloading.
func FileSHA256(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	hasher := sha256.New()
	if _, err := io.Copy(hasher, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(hasher.Sum(nil)), nil
}
