package provision

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	domainprovision "github.com/sentzunhat/hawp/librarian/src/internal/domain/provision"
)

func testManifest() *domainprovision.Manifest {
	return &domainprovision.Manifest{
		Assets: map[string]domainprovision.AssetRecord{
			"runtime": {SHA256: "abc123", Size: 3},
		},
	}
}

func TestSaveRejectsSymlinkManifestWithoutWritingThrough(t *testing.T) {
	root := t.TempDir()
	external := filepath.Join(t.TempDir(), "outside.json")
	if err := os.WriteFile(external, []byte("preserve me\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(external, filepath.Join(root, "manifest.json")); err != nil {
		t.Skipf("symlinks unavailable: %v", err)
	}

	if err := Save(testManifest(), root); err == nil {
		t.Fatal("manifest save followed a symlink destination")
	}
	if got, err := os.ReadFile(external); err != nil || string(got) != "preserve me\n" {
		t.Fatalf("symlink target changed: content=%q err=%v", got, err)
	}
}

func TestSaveReplacesHardLinkedManifestWithoutTruncatingTarget(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("hard-link behavior is not reliable on Windows")
	}
	root := t.TempDir()
	external := filepath.Join(t.TempDir(), "outside.json")
	if err := os.WriteFile(external, []byte("preserve me\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Link(external, filepath.Join(root, "manifest.json")); err != nil {
		t.Skipf("hard links unavailable: %v", err)
	}

	if err := Save(testManifest(), root); err != nil {
		t.Fatal(err)
	}
	if got, err := os.ReadFile(external); err != nil || string(got) != "preserve me\n" {
		t.Fatalf("hard-link target changed: content=%q err=%v", got, err)
	}
	if got, err := os.ReadFile(filepath.Join(root, "manifest.json")); err != nil || string(got) == "preserve me\n" {
		t.Fatalf("manifest destination was not replaced: content=%q err=%v", got, err)
	}
}
