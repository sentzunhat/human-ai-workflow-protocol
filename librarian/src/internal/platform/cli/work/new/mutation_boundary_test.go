package newcmd

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewRejectsEmptyRootBeforeWriting(t *testing.T) {
	root := t.TempDir()
	if err := Run([]string{"Should not be created", "--hawp-root="}, root); err == nil {
		t.Fatal("expected empty root to be rejected")
	}
}

func TestNewRejectsExplicitRootWithSymlinkedAncestor(t *testing.T) {
	parent := t.TempDir()
	realRoot := filepath.Join(parent, "repo")
	workDir := filepath.Join(realRoot, ".hawp", "work")
	if err := os.MkdirAll(workDir, 0o755); err != nil {
		t.Fatal(err)
	}
	linkParent := filepath.Join(t.TempDir(), "linked-parent")
	if err := os.Symlink(parent, linkParent); err != nil {
		t.Skipf("symlink capability unavailable (enable Windows Developer Mode or SeCreateSymbolicLinkPrivilege): %v", err)
	}
	if err := Run([]string{"should fail", "--hawp-root=" + filepath.Join(linkParent, "repo", ".hawp")}, parent); err == nil {
		t.Fatal("new accepted an explicit HAWP root beneath a symlinked ancestor")
	}
}
