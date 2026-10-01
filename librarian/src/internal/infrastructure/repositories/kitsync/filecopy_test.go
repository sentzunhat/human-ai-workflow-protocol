package kitsync

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTree(t *testing.T, files map[string]string) string {
	t.Helper()
	root := t.TempDir()
	for rel, content := range files {
		full := filepath.Join(root, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return root
}

func TestFileCopierCopyTree(t *testing.T) {
	fc := NewFileCopier()

	files := map[string]string{
		"src/a.txt":     "alpha",
		"src/b.txt":     "beta",
		"src/sub/c.txt": "gamma",
	}
	srcRoot := writeTree(t, files)
	destRoot := t.TempDir()

	count, err := fc.copyTree(filepath.Join(srcRoot, "src"), filepath.Join(destRoot, "dst"), "")
	if err != nil {
		t.Fatal(err)
	}
	if count != 3 {
		t.Fatalf("expected 3 files copied, got %d", count)
	}

	for _, name := range []string{"a.txt", "b.txt"} {
		got, err := os.ReadFile(filepath.Join(destRoot, "dst", name))
		if err != nil {
			t.Fatal(err)
		}
		if string(got) == "" {
			t.Errorf("dst/%s: expected content", name)
		}
	}
	got, err := os.ReadFile(filepath.Join(destRoot, "dst/sub/c.txt"))
	if err != nil {
		t.Fatal(err)
	}
	if string(got) == "" {
		t.Error("dst/sub/c.txt: expected content")
	}
}

func TestFileCopierCopyTreeWithPattern(t *testing.T) {
	fc := NewFileCopier()

	files := map[string]string{
		"src/hawp-a.md": "a",
		"src/hawp-b.md": "b",
		"src/readme.md": "readme",
	}
	srcRoot := writeTree(t, files)
	destRoot := t.TempDir()

	count, err := fc.copyTree(filepath.Join(srcRoot, "src"), filepath.Join(destRoot, "dst"), "hawp-*.md")
	if err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Fatalf("expected 2 files copied, got %d", count)
	}

	for _, name := range []string{"hawp-a.md", "hawp-b.md"} {
		if _, err := os.ReadFile(filepath.Join(destRoot, "dst", name)); err != nil {
			t.Errorf("dst/%s should exist: %v", name, err)
		}
	}
	if _, err := os.ReadFile(filepath.Join(destRoot, "dst", "readme.md")); err == nil {
		t.Error("dst/readme.md should not exist")
	}
}

func TestFileCopierSeedTree(t *testing.T) {
	fc := NewFileCopier()

	files := map[string]string{
		"src/a.txt": "alpha",
		"src/b.txt": "beta",
	}
	srcRoot := writeTree(t, files)
	destRoot := t.TempDir()

	// Pre-create b.txt at destination.
	dstDir := filepath.Join(destRoot, "dst")
	if err := os.MkdirAll(dstDir, 0o755); err != nil {
		t.Fatal(err)
	}
	existingPath := filepath.Join(dstDir, "b.txt")
	if err := os.WriteFile(existingPath, []byte("existing"), 0o644); err != nil {
		t.Fatal(err)
	}

	count, err := fc.seedTree(filepath.Join(srcRoot, "src"), filepath.Join(destRoot, "dst"), "")
	if err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Fatalf("expected 1 file seeded, got %d", count)
	}

	existingContent, _ := os.ReadFile(existingPath)
	if string(existingContent) == "beta" {
		t.Error("seed should not overwrite existing files")
	}
	gotA, err := os.ReadFile(filepath.Join(destRoot, "dst", "a.txt"))
	if err != nil {
		t.Fatal(err)
	}
	if string(gotA) != "alpha" {
		t.Errorf("dst/a.txt: expected alpha, got %s", string(gotA))
	}
}

func TestFileCopierCopyFileAtomicity(t *testing.T) {
	fc := NewFileCopier()

	files := map[string]string{"src/input.txt": "test-data"}
	srcRoot := writeTree(t, files)
	destRoot := t.TempDir()

	dstPath := filepath.Join(destRoot, "output.txt")
	if err := fc.copyFile(filepath.Join(srcRoot, "src", "input.txt"), dstPath); err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(dstPath)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != "test-data" {
		t.Errorf("output: %q, want %q", string(got), "test-data")
	}

	// No temp file should remain.
	entries, _ := os.ReadDir(destRoot)
	for _, e := range entries {
		if e.Name() == ".kitsync-*" {
			t.Error("leftover temp file after copy")
		}
	}
}
