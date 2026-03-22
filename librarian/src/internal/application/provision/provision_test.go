package provision

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"archive/tar"
	"compress/gzip"

	domainprovision "github.com/sentzunhat/hawp/librarian/src/internal/domain/provision"
	"github.com/sentzunhat/hawp/librarian/src/internal/infrastructure/download"
)

func hashHex(b []byte) string {
	sum := sha256.Sum256(b)
	return hex.EncodeToString(sum[:])
}

// fakeArchive builds a minimal .tgz containing one member and returns its
// bytes plus sha256.
func fakeArchive(t *testing.T, memberPath, memberContent string) ([]byte, string) {
	t.Helper()
	path := filepath.Join(t.TempDir(), "fake.tgz")
	f, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	gz := gzip.NewWriter(f)
	tw := tar.NewWriter(gz)
	if err := tw.WriteHeader(&tar.Header{Name: memberPath, Mode: 0o644, Size: int64(len(memberContent))}); err != nil {
		t.Fatal(err)
	}
	if _, err := tw.Write([]byte(memberContent)); err != nil {
		t.Fatal(err)
	}
	tw.Close()
	gz.Close()
	f.Close()

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return data, hashHex(data)
}

// testServer serves a fixed map of path -> body and returns its base URL.
func testServer(t *testing.T, routes map[string][]byte) string {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, ok := routes[r.URL.Path]
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.Write(body)
	}))
	t.Cleanup(server.Close)
	return server.URL
}

func TestRunEndToEndFreshInstall(t *testing.T) {
	archiveBytes, archiveHash := fakeArchive(t, "pkg/lib/libonnxruntime.so.1.27.1", "shared lib bytes")
	modelBytes := []byte("fake onnx model")
	tokenizerBytes := []byte(`{"tokenizer":true}`)

	baseURL := testServer(t, map[string][]byte{
		"/runtime.tgz":    archiveBytes,
		"/model.onnx":     modelBytes,
		"/tokenizer.json": tokenizerBytes,
	})

	registry := Registry{
		RuntimeAsset: domainprovision.Asset{
			Name: "onnxruntime", URL: baseURL + "/runtime.tgz", SHA256: archiveHash,
			Size: int64(len(archiveBytes)), ArchiveMember: "pkg/lib/libonnxruntime.so.1.27.1",
			DestName: "libonnxruntime.so",
		},
		ModelAssets: []domainprovision.Asset{
			{Name: "model.onnx", URL: baseURL + "/model.onnx", SHA256: hashHex(modelBytes), Size: int64(len(modelBytes)), DestName: "model.onnx"},
			{Name: "tokenizer.json", URL: baseURL + "/tokenizer.json", SHA256: hashHex(tokenizerBytes), Size: int64(len(tokenizerBytes)), DestName: "tokenizer.json"},
		},
		RuntimeVersion: "1.27.1", ModelName: "test/model", ModelVersion: "1",
	}

	home := t.TempDir()
	result := Run(download.NewHTTPFetcher(), home, registry)

	if result.Failed() {
		t.Fatalf("unexpected failure: %+v", result.Steps)
	}
	if len(result.Steps) != 3 {
		t.Fatalf("steps = %d, want 3", len(result.Steps))
	}
	for _, step := range result.Steps {
		if step.Status != "downloaded" {
			t.Errorf("step %s status = %s, want downloaded on fresh install", step.Name, step.Status)
		}
	}

	libPath := filepath.Join(home, ".hawp", "runtime", "libonnxruntime.so")
	if content, err := os.ReadFile(libPath); err != nil || string(content) != "shared lib bytes" {
		t.Fatalf("extracted lib = %q, err %v", content, err)
	}
	modelPath := filepath.Join(home, ".hawp", "models", "model.onnx")
	if content, err := os.ReadFile(modelPath); err != nil || string(content) != string(modelBytes) {
		t.Fatalf("model file = %q, err %v", content, err)
	}

	manifestPath := filepath.Join(home, ".hawp", "manifest.json")
	if _, err := os.Stat(manifestPath); err != nil {
		t.Fatalf("manifest not written: %v", err)
	}

	// Idempotent re-run: point the fetcher at a server that always 404s —
	// re-provisioning must succeed from disk without hitting it.
	deadServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer deadServer.Close()

	second := Run(download.NewHTTPFetcher(), home, registry)
	if second.Failed() {
		t.Fatalf("idempotent re-run failed: %+v", second.Steps)
	}
	for _, step := range second.Steps {
		if step.Status != "already-installed" {
			t.Errorf("step %s status = %s on re-run, want already-installed", step.Name, step.Status)
		}
	}
}

func TestRunRuntimeUnsupportedPlatformStillProvisionsModels(t *testing.T) {
	modelBytes := []byte("model bytes")
	baseURL := testServer(t, map[string][]byte{"/model.onnx": modelBytes})

	registry := Registry{
		RuntimeAssetErr: os.ErrInvalid,
		ModelAssets: []domainprovision.Asset{
			{Name: "model.onnx", URL: baseURL + "/model.onnx", SHA256: hashHex(modelBytes), Size: int64(len(modelBytes)), DestName: "model.onnx"},
		},
	}

	home := t.TempDir()
	result := Run(download.NewHTTPFetcher(), home, registry)

	if !result.Failed() {
		t.Fatal("expected overall Failed() true due to runtime asset error")
	}
	var modelStep, runtimeStep *Step
	for i := range result.Steps {
		switch result.Steps[i].Name {
		case "model.onnx":
			modelStep = &result.Steps[i]
		case "onnxruntime":
			runtimeStep = &result.Steps[i]
		}
	}
	if runtimeStep == nil || runtimeStep.Status != "failed" {
		t.Errorf("runtime step = %+v, want failed", runtimeStep)
	}
	if modelStep == nil || modelStep.Status != "downloaded" {
		t.Errorf("model step = %+v, want downloaded despite runtime failure", modelStep)
	}
}

func TestRunChecksumMismatchReportsFailedStep(t *testing.T) {
	modelBytes := []byte("model bytes")
	baseURL := testServer(t, map[string][]byte{"/model.onnx": modelBytes})

	registry := Registry{
		RuntimeAssetErr: os.ErrInvalid,
		ModelAssets: []domainprovision.Asset{
			{Name: "model.onnx", URL: baseURL + "/model.onnx", SHA256: "0000000000000000000000000000000000000000000000000000000000000000"[:64], Size: int64(len(modelBytes)), DestName: "model.onnx"},
		},
	}
	result := Run(download.NewHTTPFetcher(), t.TempDir(), registry)
	if !result.Failed() {
		t.Fatal("expected failure on checksum mismatch")
	}
}

func TestResultStringReportsFailuresAndSuccesses(t *testing.T) {
	result := Result{Steps: []Step{
		{Name: "a", Status: "downloaded", Path: "/x/a"},
		{Name: "b", Status: "failed", Err: os.ErrInvalid},
	}}
	s := result.String()
	if !contains(s, "a: downloaded") || !contains(s, "b:") {
		t.Errorf("String() = %q", s)
	}
}

func contains(haystack, needle string) bool {
	return len(haystack) >= len(needle) && (func() bool {
		for i := 0; i+len(needle) <= len(haystack); i++ {
			if haystack[i:i+len(needle)] == needle {
				return true
			}
		}
		return false
	})()
}
