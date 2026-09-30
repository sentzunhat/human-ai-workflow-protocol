package work

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDuplicateLinksRejectsSymlinkedWorkRoot(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, ".hawp"), 0o755); err != nil {
		t.Fatal(err)
	}
	outside := t.TempDir()
	if err := os.Symlink(outside, filepath.Join(root, ".hawp", "work")); err != nil {
		t.Skipf("symlink capability unavailable (enable Windows Developer Mode or SeCreateSymbolicLinkPrivilege): %v", err)
	}

	if _, err := PreviewDuplicateLinks(root); err == nil {
		t.Fatal("expected preview to reject symlinked work root")
	} else if !strings.Contains(err.Error(), "symlink") {
		t.Fatalf("expected preview symlink rejection, got %v", err)
	}
	if _, err := ApplyDuplicateLinks(root); err == nil {
		t.Fatal("expected apply to reject symlinked work root")
	} else if !strings.Contains(err.Error(), "symlink") {
		t.Fatalf("expected apply symlink rejection, got %v", err)
	}
}

func TestDuplicateLinksRejectsSymlinkedScopeDirectory(t *testing.T) {
	root := buildDuplicateFixture(t, map[string]string{
		".hawp/work/BACKLOG.md": duplicateBacklogHeader + duplicateBacklogFooter,
	})
	outside := t.TempDir()
	outsidePath := filepath.Join(outside, "BUG-901.md")
	if err := os.WriteFile(outsidePath, []byte("# outside\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(outside, filepath.Join(root, ".hawp", "work", "active")); err != nil {
		t.Skipf("symlink capability unavailable (enable Windows Developer Mode or SeCreateSymbolicLinkPrivilege): %v", err)
	}

	_, err := ApplyDuplicateLinks(root)
	if err == nil {
		t.Fatal("expected symlinked active scope to be rejected")
	}
	if !strings.Contains(err.Error(), "symlink") {
		t.Fatalf("expected symlink rejection, got %v", err)
	}
	data, readErr := os.ReadFile(outsidePath)
	if readErr != nil {
		t.Fatal(readErr)
	}
	if string(data) != "# outside\n" {
		t.Fatal("outside file changed through symlinked active scope")
	}
}

func TestDuplicateLinksRejectsSymlinkedPlanBeforeReading(t *testing.T) {
	root := buildDuplicateFixture(t, map[string]string{
		".hawp/work/BACKLOG.md": duplicateBacklogHeader + duplicateBacklogFooter,
	})
	outside := filepath.Join(t.TempDir(), "outside.md")
	if err := os.WriteFile(outside, []byte("# OUTSIDE_DUPLICATE_SENTINEL\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(root, ".hawp", "work", "active"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(outside, filepath.Join(root, ".hawp", "work", "active", "linked-plan.md")); err != nil {
		t.Skipf("symlink capability unavailable (enable Windows Developer Mode or SeCreateSymbolicLinkPrivilege): %v", err)
	}

	if _, err := PreviewDuplicateLinks(root); err == nil {
		t.Fatal("expected preview to reject symlinked plan")
	} else if !strings.Contains(err.Error(), "symlink") {
		t.Fatalf("expected symlink rejection, got %v", err)
	}
}
