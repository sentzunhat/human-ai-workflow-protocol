package work

import (
	"os"
	"path/filepath"
	"testing"
)

func TestAdapterReadCollectsBacklogAndMarkdown(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "active"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "BACKLOG.md"), []byte("# Backlog\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "active", "item.md"), []byte("# Plan\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	snapshot, err := NewAdapter().Read(root)
	if err != nil {
		t.Fatal(err)
	}
	if snapshot.BacklogContent != "# Backlog\n" || len(snapshot.Files) != 2 {
		t.Fatalf("snapshot = %+v", snapshot)
	}
}
