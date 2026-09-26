package benchmark

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestRunRejectsSymlinkedSearchDatabaseAncestry(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("symlink permissions are not reliable on Windows")
	}

	root := t.TempDir()
	workDir := filepath.Join(root, ".hawp", "work")
	if err := os.MkdirAll(workDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(workDir, "BACKLOG.md"), []byte("# Backlog\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	external := t.TempDir()
	dbParent := filepath.Join(root, ".hawp", "db")
	if err := os.Symlink(external, dbParent); err != nil {
		t.Skipf("symlinks unavailable: %v", err)
	}

	err := Run(nil, root)
	if err == nil {
		t.Fatal("accepted symlinked search database ancestry")
	}
	if !strings.Contains(err.Error(), "search index path") || !strings.Contains(err.Error(), "symlink") {
		t.Fatalf("unexpected error: %v", err)
	}
}
