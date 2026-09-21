package validate

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestRunRejectsSymlinkedDefaultKitTree(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("symlink permissions are not reliable on Windows")
	}

	root := t.TempDir()
	backlog := filepath.Join(root, ".hawp", "work", "BACKLOG.md")
	if err := os.MkdirAll(filepath.Dir(backlog), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(backlog, []byte("# Backlog\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	external := t.TempDir()
	if err := os.WriteFile(filepath.Join(external, "outside.md"), []byte("# outside\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	kitPath := filepath.Join(root, ".hawp", "kit")
	if err := os.Symlink(external, kitPath); err != nil {
		t.Skipf("symlinks unavailable: %v", err)
	}

	err := Run(nil, root)
	if err == nil {
		t.Fatal("accepted symlinked default kit root")
	}
	if !strings.Contains(err.Error(), "unsafe kit validation path") || !strings.Contains(err.Error(), "symlink") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunRejectsSymlinkedDefaultKitDescendant(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("symlink permissions are not reliable on Windows")
	}

	root := t.TempDir()
	backlog := filepath.Join(root, ".hawp", "work", "BACKLOG.md")
	kitPath := filepath.Join(root, ".hawp", "kit")
	if err := os.MkdirAll(filepath.Dir(backlog), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(kitPath, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(backlog, []byte("# Backlog\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	external := filepath.Join(t.TempDir(), "outside.md")
	if err := os.WriteFile(external, []byte("# outside\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(external, filepath.Join(kitPath, "outside.md")); err != nil {
		t.Skipf("symlinks unavailable: %v", err)
	}

	err := Run(nil, root)
	if err == nil {
		t.Fatal("accepted symlinked markdown descendant in default kit")
	}
	if !strings.Contains(err.Error(), "unsafe kit validation path") || !strings.Contains(err.Error(), "symlink") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunRejectsSymlinkedExplicitKitRoot(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("symlink permissions are not reliable on Windows")
	}

	external := t.TempDir()
	link := filepath.Join(t.TempDir(), "kit")
	if err := os.Symlink(external, link); err != nil {
		t.Skipf("symlinks unavailable: %v", err)
	}

	err := Run([]string{"--kit-path", link}, t.TempDir())
	if err == nil {
		t.Fatal("accepted symlinked explicit kit root")
	}
	if !strings.Contains(err.Error(), "unsafe explicit kit validation path") || !strings.Contains(err.Error(), "symlink") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunRejectsSymlinkedAncestorForRelativeExplicitKitPath(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("symlink permissions are not reliable on Windows")
	}

	cwd := t.TempDir()
	outside := t.TempDir()
	if err := os.MkdirAll(filepath.Join(outside, "kit"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(outside, filepath.Join(cwd, "linked")); err != nil {
		t.Skipf("symlinks unavailable: %v", err)
	}

	err := Run([]string{"--kit-path", filepath.Join("linked", "kit")}, cwd)
	if err == nil {
		t.Fatal("accepted explicit kit path below symlinked ancestor")
	}
	if !strings.Contains(err.Error(), "unsafe explicit kit validation path") || !strings.Contains(err.Error(), "symlink") {
		t.Fatalf("unexpected error: %v", err)
	}
}
