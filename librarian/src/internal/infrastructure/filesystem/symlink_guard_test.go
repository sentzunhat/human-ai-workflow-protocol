package filesystem

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func TestRejectSymlinkedWorkRootRejectsHAWPAncestor(t *testing.T) {
	repoRoot := t.TempDir()
	outside := t.TempDir()
	if err := os.MkdirAll(filepath.Join(outside, "work"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(outside, filepath.Join(repoRoot, ".hawp")); err != nil {
		t.Skipf("symlink capability unavailable (enable Windows Developer Mode or SeCreateSymbolicLinkPrivilege): %v", err)
	}
	if err := RejectSymlinkedWorkRoot(filepath.Join(repoRoot, ".hawp", "work")); err == nil {
		t.Fatal("expected symlinked .hawp ancestor to be rejected")
	}
}

func TestRejectSymlinksInTreeRejectsNestedSymlink(t *testing.T) {
	root := t.TempDir()
	tree := filepath.Join(root, "kit")
	outside := filepath.Join(t.TempDir(), "outside.md")
	if err := os.MkdirAll(tree, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(outside, []byte("outside"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(outside, filepath.Join(tree, "linked.md")); err != nil {
		t.Skipf("symlink capability unavailable (enable Windows Developer Mode or SeCreateSymbolicLinkPrivilege): %v", err)
	}

	if err := RejectSymlinksInTree(root, tree); err == nil {
		t.Fatal("expected nested symlink to be rejected")
	}
}

func TestRejectSymlinksInTreeAllowsMissingTarget(t *testing.T) {
	root := t.TempDir()
	if err := RejectSymlinksInTree(root, filepath.Join(root, "missing", "kit")); err != nil {
		t.Fatalf("missing target should be allowed: %v", err)
	}
}

func TestAtomicWriteFileRejectsSymlinkedDestination(t *testing.T) {
	root := t.TempDir()
	outside := filepath.Join(root, "outside.txt")
	if err := os.WriteFile(outside, []byte("outside"), 0o644); err != nil {
		t.Fatal(err)
	}
	target := filepath.Join(root, "target.txt")
	if err := os.Symlink(outside, target); err != nil {
		t.Skipf("symlink capability unavailable (enable Windows Developer Mode or SeCreateSymbolicLinkPrivilege): %v", err)
	}
	if err := AtomicWriteFile(root, target, []byte("replacement"), 0o644); err == nil {
		t.Fatal("expected symlinked destination to be rejected")
	}
	data, err := os.ReadFile(outside)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "outside" {
		t.Fatalf("outside target was modified: %q", data)
	}
}

func TestAtomicWriteFileAtPathRejectsSymlinkedDestination(t *testing.T) {
	root := t.TempDir()
	outside := filepath.Join(root, "outside.txt")
	if err := os.WriteFile(outside, []byte("outside"), 0o644); err != nil {
		t.Fatal(err)
	}
	target := filepath.Join(root, "nested", "target.txt")
	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(outside, target); err != nil {
		t.Skipf("symlink capability unavailable (enable Windows Developer Mode or SeCreateSymbolicLinkPrivilege): %v", err)
	}
	if err := AtomicWriteFileAtPath(target, []byte("replacement"), 0o644); err == nil {
		t.Fatal("expected symlinked destination to be rejected")
	}
	data, err := os.ReadFile(outside)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "outside" {
		t.Fatalf("outside target was modified: %q", data)
	}
}

func TestAtomicWriteFileAtPathReplacesHardLinkWithoutChangingExternalTarget(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("hard-link behavior is not reliable on Windows")
	}
	root := t.TempDir()
	external := filepath.Join(t.TempDir(), "external.txt")
	if err := os.WriteFile(external, []byte("preserve"), 0o644); err != nil {
		t.Fatal(err)
	}
	target := filepath.Join(root, "new", "target.txt")
	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Link(external, target); err != nil {
		t.Skipf("hard links unavailable: %v", err)
	}
	if err := AtomicWriteFileAtPath(target, []byte("replacement"), 0o644); err != nil {
		t.Fatal(err)
	}
	externalData, err := os.ReadFile(external)
	if err != nil {
		t.Fatal(err)
	}
	if string(externalData) != "preserve" {
		t.Fatalf("external hard-link target was modified: %q", externalData)
	}
	targetData, err := os.ReadFile(target)
	if err != nil {
		t.Fatal(err)
	}
	if string(targetData) != "replacement" {
		t.Fatalf("replacement destination = %q, want replacement", targetData)
	}
}

func TestRejectSymlinkAncestorsRejectsTargetOutsideRoot(t *testing.T) {
	root := t.TempDir()
	target := filepath.Join(root, "..", "outside")

	err := RejectSymlinkAncestors(root, target)
	if err == nil || !strings.Contains(err.Error(), "outside root") {
		t.Fatalf("expected outside-root rejection, got %v", err)
	}
}

func TestRejectSymlinkAncestorsFailsClosedWhenRootCannotBeStatted(t *testing.T) {
	root := filepath.Join(t.TempDir(), "missing-root")
	target := filepath.Join(root, "nested", "file.txt")

	err := RejectSymlinkAncestors(root, target)
	if err == nil || !strings.Contains(err.Error(), "stat root") {
		t.Fatalf("expected root stat error, got %v", err)
	}
}

func TestRejectSymlinksInPathRejectsSymlinkedAncestor(t *testing.T) {
	base := t.TempDir()
	outside := t.TempDir()
	target := filepath.Join(outside, "kit")
	if err := os.MkdirAll(target, 0o755); err != nil {
		t.Fatal(err)
	}
	link := filepath.Join(base, "linked-parent")
	if err := os.Symlink(outside, link); err != nil {
		t.Skipf("symlink capability unavailable (enable Windows Developer Mode or SeCreateSymbolicLinkPrivilege): %v", err)
	}
	if err := RejectSymlinksInPath(filepath.Join(link, "kit")); err == nil {
		t.Fatal("expected symlinked ancestor to be rejected")
	}
}

func TestRootGuardRejectsAncestorBeforeReadOrWrite(t *testing.T) {
	base, outside := t.TempDir(), t.TempDir()
	realRoot := filepath.Join(outside, "repo")
	if err := os.Mkdir(realRoot, 0o755); err != nil {
		t.Fatal(err)
	}
	target := filepath.Join(realRoot, "file")
	if err := os.WriteFile(target, []byte("preserve"), 0o600); err != nil {
		t.Fatal(err)
	}
	link := filepath.Join(base, "redirect")
	if err := os.Symlink(outside, link); err != nil {
		t.Skipf("symlinks unavailable: %v", err)
	}
	root := filepath.Join(link, "repo")
	if err := RejectSymlinkAncestors(root, filepath.Join(root, "missing")); err == nil {
		t.Error("accepted ancestor of root")
	}
	if _, err := ResolveSafeSearchIndexPath(root); err == nil {
		t.Error("accepted redirected search root")
	}
	if err := AtomicWriteFile(root, filepath.Join(root, "file"), []byte("changed"), 0o600); err == nil {
		t.Error("accepted redirected write")
	}
	got, err := os.ReadFile(target)
	if err != nil || string(got) != "preserve" {
		t.Fatalf("external file changed: %q, %v", got, err)
	}
}
