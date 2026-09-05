package cli

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestMCPConfigureExplicitRootWithoutProvisioning(t *testing.T) {
	root, home := t.TempDir(), t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("USERPROFILE", home)
	name := "hawp"
	if runtime.GOOS == "windows" {
		name += ".exe"
	}
	for _, path := range []string{filepath.Join(root, ".hawp/bin", name), filepath.Join(root, ".hawp/work/BACKLOG.md")} {
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte("fixture"), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	if err := Run([]string{"mcp", "configure", "--repo-root", root, "--provider=claude"}); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(root, ".mcp.json")); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(root, ".hawp/kit")); !os.IsNotExist(err) {
		t.Fatal("kit changed")
	}
	entries, err := os.ReadDir(home)
	if err != nil || len(entries) != 0 {
		t.Fatal("provisioned home", entries, err)
	}
}
