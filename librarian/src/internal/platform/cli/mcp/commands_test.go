package mcp

import (
	"os"
	"path/filepath"
	"testing"
)

func TestMCPExplicitRootIndependentOfCWD(t *testing.T) {
	cwd, target := t.TempDir(), t.TempDir()
	got, err := resolveRoot([]string{"--repo-root", target}, cwd)
	if err != nil || got != target {
		t.Fatalf("got %q, %v; want %q", got, err, target)
	}
	child := filepath.Join(cwd, "project with spaces")
	if err := os.Mkdir(child, 0o755); err != nil {
		t.Fatal(err)
	}
	got, err = resolveRoot([]string{"--repo-root=project with spaces"}, cwd)
	if err != nil || got != child {
		t.Fatalf("relative root: %q, %v", got, err)
	}
	got, err = resolveRoot(nil, cwd)
	if err != nil || got != cwd {
		t.Fatalf("legacy fallback changed: %q, %v", got, err)
	}
}

func TestMCPRootRejectsInvalidInput(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "file")
	if err := os.WriteFile(file, nil, 0o644); err != nil {
		t.Fatal(err)
	}
	for _, args := range [][]string{{"--repo-root"}, {"--repo-root="}, {"--repo-root", file}, {"--repo-root", filepath.Join(dir, "missing")}, {"--unknown"}, {"extra"}} {
		if _, err := resolveRoot(args, dir); err == nil {
			t.Errorf("accepted %v", args)
		}
	}
}

func TestMCPRootRejectsSymlinkedAncestors(t *testing.T) {
	base, outside := t.TempDir(), t.TempDir()
	if err := os.Mkdir(filepath.Join(outside, "repo"), 0o755); err != nil {
		t.Fatal(err)
	}
	link := filepath.Join(base, "redirect")
	if err := os.Symlink(outside, link); err != nil {
		t.Skipf("symlinks unavailable: %v", err)
	}
	for _, candidate := range []string{link, filepath.Join(link, "repo"), filepath.Join("redirect", "repo")} {
		if _, err := resolveRoot([]string{"--repo-root", candidate}, base); err == nil {
			t.Errorf("accepted symlinked root %q", candidate)
		}
	}
}
