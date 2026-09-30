package location

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRootForHomeRejectsSymlinkedModelsRoot(t *testing.T) {
	home := t.TempDir()
	outside := t.TempDir()
	if err := os.MkdirAll(filepath.Join(home, ".hawp"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(outside, filepath.Join(home, ".hawp", "models")); err != nil {
		t.Skipf("symlinks unavailable: %v", err)
	}
	if _, err := rootForHome(home); err == nil {
		t.Fatal("rootForHome accepted symlinked models root")
	}
}
