package validatecmd

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestRunRejectsSymlinkedDefaultWorkTree(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("symlink permissions are not reliable on Windows")
	}

	root := t.TempDir()
	hawp := filepath.Join(root, ".hawp")
	if err := os.MkdirAll(hawp, 0o755); err != nil {
		t.Fatal(err)
	}

	external := t.TempDir()
	if err := os.WriteFile(filepath.Join(external, "BACKLOG.md"), []byte("# Backlog\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(external, filepath.Join(hawp, "work")); err != nil {
		t.Skipf("symlinks unavailable: %v", err)
	}

	err := Run(nil, root)
	if err == nil {
		t.Fatal("accepted symlinked default work root")
	}
	if !strings.Contains(err.Error(), "unsafe work validation path") || !strings.Contains(err.Error(), "symlink") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunRejectsSymlinkedDefaultWorkDescendant(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("symlink permissions are not reliable on Windows")
	}

	root := t.TempDir()
	workDir := filepath.Join(root, ".hawp", "work")
	if err := os.MkdirAll(filepath.Join(workDir, "active"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(workDir, "BACKLOG.md"), []byte("# Backlog\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	external := filepath.Join(t.TempDir(), "outside.md")
	if err := os.WriteFile(external, []byte("# outside\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(external, filepath.Join(workDir, "active", "outside.md")); err != nil {
		t.Skipf("symlinks unavailable: %v", err)
	}

	err := Run(nil, root)
	if err == nil {
		t.Fatal("accepted symlinked descendant in default work tree")
	}
	if !strings.Contains(err.Error(), "unsafe work validation path") || !strings.Contains(err.Error(), "symlink") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunRejectsSymlinkedExplicitWorkRoot(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("symlink permissions are not reliable on Windows")
	}

	external := t.TempDir()
	if err := os.WriteFile(filepath.Join(external, "BACKLOG.md"), []byte("# Backlog\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	link := filepath.Join(t.TempDir(), "work")
	if err := os.Symlink(external, link); err != nil {
		t.Skipf("symlinks unavailable: %v", err)
	}

	err := Run([]string{"--work-root", link}, t.TempDir())
	if err == nil {
		t.Fatal("accepted symlinked explicit work root")
	}
	if !strings.Contains(err.Error(), "unsafe work validation path") || !strings.Contains(err.Error(), "symlink") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunRejectsSymlinkedAncestorForRelativeExplicitWorkRoot(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("symlink permissions are not reliable on Windows")
	}

	cwd := t.TempDir()
	outside := t.TempDir()
	work := filepath.Join(outside, "work")
	if err := os.MkdirAll(work, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(work, "BACKLOG.md"), []byte("# Backlog\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(outside, filepath.Join(cwd, "linked")); err != nil {
		t.Skipf("symlinks unavailable: %v", err)
	}

	err := Run([]string{"--work-root", filepath.Join("linked", "work")}, cwd)
	if err == nil {
		t.Fatal("accepted explicit work root below symlinked ancestor")
	}
	if !strings.Contains(err.Error(), "unsafe work validation path") || !strings.Contains(err.Error(), "symlink") {
		t.Fatalf("unexpected error: %v", err)
	}
}
