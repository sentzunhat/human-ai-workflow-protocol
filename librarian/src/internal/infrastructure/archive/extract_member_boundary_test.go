package archive

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestWriteMemberRejectsSymlinkedParentBeforeReading(t *testing.T) {
	root, outside := t.TempDir(), t.TempDir()
	parent := filepath.Join(root, "redirect")
	if err := os.Symlink(outside, parent); err != nil {
		t.Skipf("symlinks unavailable: %v", err)
	}
	if err := writeMember(filepath.Join(parent, "asset.bin"), strings.NewReader("payload")); err == nil {
		t.Fatal("accepted symlinked extraction parent")
	}
	entries, err := os.ReadDir(outside)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 0 {
		t.Fatalf("extraction wrote outside destination: %v", entries)
	}
}
