package configure

import (
	"os"
	"path/filepath"
	"testing"
)

func TestMCPConfigureRejectsInvalidArguments(t *testing.T) {
	cwd := t.TempDir()
	for _, args := range [][]string{
		{"--provider"}, {"--provider="}, {"--repo-root="},
		{"--repo-root"}, {"--unknown"}, {"extra"},
	} {
		if err := Run(args, cwd); err == nil {
			t.Errorf("accepted %v", args)
		}
	}
}

func TestConfigureRootRejectsSymlinkedAncestors(t *testing.T) {
	base, outside := t.TempDir(), t.TempDir()
	if err := os.Mkdir(filepath.Join(outside, "repo"), 0o755); err != nil {
		t.Fatal(err)
	}
	link := filepath.Join(base, "redirect")
	if err := os.Symlink(outside, link); err != nil {
		t.Skipf("symlinks unavailable: %v", err)
	}
	for _, candidate := range []string{link, filepath.Join(link, "repo"), filepath.Join("redirect", "repo")} {
		if _, err := resolveConfigureRoot(candidate, base, true); err == nil {
			t.Errorf("accepted symlinked root %q", candidate)
		}
	}
}
