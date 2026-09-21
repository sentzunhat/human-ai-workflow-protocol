package filesystem

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestRejectSymlinkAncestorsRejectsTargetOutsideRoot(t *testing.T) {
	root := t.TempDir()
	target := filepath.Join(root, "..", "outside")

	err := RejectSymlinkAncestors(root, target)
	if err == nil || !strings.Contains(err.Error(), "outside root") {
		t.Fatalf("expected outside-root rejection, got %v", err)
	}
}
