package filesystem_test

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	infrafs "github.com/sentzunhat/hawp/librarian/src/internal/infrastructure/filesystem"
)

func TestGenerateREADMEsRejectsSymlinkedProjectRootBeforeWritingHome(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("symlink permissions are not reliable on Windows")
	}

	home := t.TempDir()
	project := t.TempDir()
	outside := t.TempDir()
	if err := os.Symlink(outside, filepath.Join(project, ".hawp")); err != nil {
		t.Skipf("symlinks unavailable: %v", err)
	}

	if err := infrafs.GenerateREADMEs(home, project); err == nil {
		t.Fatal("expected symlinked project root to be rejected")
	}
	if _, err := os.Stat(filepath.Join(home, ".hawp")); !os.IsNotExist(err) {
		t.Fatalf("home README tree was created before rejecting project path: %v", err)
	}
}
