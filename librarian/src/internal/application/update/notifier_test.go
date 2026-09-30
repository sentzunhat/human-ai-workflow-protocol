package update

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSetAutoUpdateRejectsSymlinkedConfig(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	outside := t.TempDir()
	manifest := filepath.Join(outside, "update.json")
	if err := os.WriteFile(manifest, []byte("keep"), 0o644); err != nil {
		t.Fatal(err)
	}
	config := filepath.Join(home, ".hawp", "config")
	if err := os.MkdirAll(config, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(manifest, filepath.Join(config, "update.json")); err != nil {
		t.Skipf("symlinks unavailable: %v", err)
	}
	if err := SetAutoUpdate(false); err == nil {
		t.Fatal("SetAutoUpdate accepted symlinked config")
	}
	content, err := os.ReadFile(manifest)
	if err != nil || string(content) != "keep" {
		t.Fatalf("external config changed: %q, %v", content, err)
	}
}

func TestSaveCacheRejectsSymlinkedCache(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	outside := t.TempDir()
	cache := filepath.Join(outside, "update-check.json")
	if err := os.WriteFile(cache, []byte("keep"), 0o644); err != nil {
		t.Fatal(err)
	}
	cacheDir := filepath.Join(home, ".hawp", "cache")
	if err := os.MkdirAll(cacheDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(cache, filepath.Join(cacheDir, "update-check.json")); err != nil {
		t.Skipf("symlinks unavailable: %v", err)
	}
	saveCache(&updateCache{Current: "old"})
	content, err := os.ReadFile(cache)
	if err != nil || string(content) != "keep" {
		t.Fatalf("external cache changed: %q, %v", content, err)
	}
}
