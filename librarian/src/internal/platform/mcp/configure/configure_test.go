package mcp

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func configureFixture(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	for _, path := range []string{hawpBinaryPath(root), filepath.Join(root, ".hawp", "work", "BACKLOG.md")} {
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte("fixture"), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	return root
}

func TestConfigureFreshAndRepeat(t *testing.T) {
	root := configureFixture(t)
	for i := 0; i < 2; i++ {
		if err := Configure(root, []string{"claude", "codex", "cursor"}); err != nil {
			t.Fatal(err)
		}
	}
	for _, name := range []string{".mcp.json", ".codex/config.toml", ".cursor/mcp.json"} {
		data, err := os.ReadFile(filepath.Join(root, name))
		if err != nil || !strings.Contains(string(data), "--repo-root") {
			t.Fatalf("%s: %s %v", name, data, err)
		}
	}
	data, err := os.ReadFile(filepath.Join(root, ".hawp/work/BACKLOG.md"))
	if err != nil || string(data) != "fixture" {
		t.Fatal("modified backlog", err)
	}
	if _, err := os.Stat(filepath.Join(root, ".hawp/kit")); !os.IsNotExist(err) {
		t.Fatal("created kit")
	}
}

func TestConfigurePreservesRemoteCodexBeforeOtherWrites(t *testing.T) {
	root := configureFixture(t)
	path := filepath.Join(root, ".codex/config.toml")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	input := "[mcp_servers.hawp]\nurl = \"https://example.invalid/mcp\"\nenabled = false\n"
	if err := os.WriteFile(path, []byte(input), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := Configure(root, []string{"claude", "codex"}); err == nil {
		t.Fatal("replaced custom Codex configuration")
	}
	data, err := os.ReadFile(path)
	if err != nil || string(data) != input {
		t.Fatal("changed custom config", err)
	}
	if _, err := os.Stat(filepath.Join(root, ".mcp.json")); !os.IsNotExist(err) {
		t.Fatal("partial configuration write")
	}
}

func TestConfigureMissingBinaryAndInvalidSelection(t *testing.T) {
	for _, providers := range [][]string{nil, {"codxe"}, {"claude"}} {
		root := t.TempDir()
		if err := Configure(root, providers); err == nil {
			t.Fatal("accepted missing prerequisites")
		}
		entries, err := os.ReadDir(root)
		if err != nil || len(entries) != 0 {
			t.Fatal("wrote configuration", err)
		}
	}
}

func TestConfigureMigratesCustomCodexPolicy(t *testing.T) {
	root := configureFixture(t)
	path := filepath.Join(root, ".codex/config.toml")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte("# keep\n[mcp_servers.hawp]\ncommand='old'\nenabled=false\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := Configure(root, []string{"codex"}); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(path)
	if err != nil || !strings.Contains(string(data), "enabled=false\n") || !strings.HasPrefix(string(data), "# keep\n") {
		t.Fatalf("lost policy: %s %v", data, err)
	}
	info, err := os.Stat(path)
	if err != nil || info.Mode().Perm() != 0o600 {
		t.Fatal("changed permissions", err)
	}
}
