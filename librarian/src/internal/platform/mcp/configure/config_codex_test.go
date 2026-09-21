package mcp

import (
	"bytes"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/pelletier/go-toml/v2"
)

func TestCodexMigrationPreservesCustomSettingsAndComments(t *testing.T) {
	root := t.TempDir()
	input := "# keep header\r\nmodel = 'custom'\r\n[ mcp_servers . \"hawp\" ] # keep table note\r\n" +
		"command = 'old' # keep command note\r\nargs = [\r\n 'mcp', '--no-update-check',\r\n]\r\n" +
		"enabled = false # never enable silently\r\nstartup_timeout_sec = 123\r\n" +
		"enabled_tools = ['hawp_work_validate']\r\ncustom = '''line one\r\n[mcp_servers.fake]\r\nline three'''\r\n" +
		"[mcp_servers.hawp.env]\r\nLOCAL_ONLY = 'keep'\r\n[mcp_servers.other]\r\ncommand = 'untouched'\r\n"
	out, err := mergeCodexTOML([]byte(input), root)
	if err != nil {
		t.Fatal(err)
	}
	for _, exact := range []string{"# keep header\r\nmodel = 'custom'", "[ mcp_servers . \"hawp\" ] # keep table note\r\n", "# keep command note\r\n", "enabled = false # never enable silently\r\n", "startup_timeout_sec = 123\r\n", "custom = '''line one\r\n[mcp_servers.fake]\r\nline three'''\r\n", "[mcp_servers.hawp.env]\r\nLOCAL_ONLY = 'keep'\r\n[mcp_servers.other]\r\ncommand = 'untouched'\r\n"} {
		if !strings.Contains(string(out), exact) {
			t.Fatalf("lost bytes %q in %s", exact, out)
		}
	}
	var doc map[string]any
	if err := toml.Unmarshal(out, &doc); err != nil {
		t.Fatal(err)
	}
	server := doc["mcp_servers"].(map[string]any)["hawp"].(map[string]any)
	if server["enabled"] != false || server["command"] != hawpBinaryPath(root) || server["cwd"] != root {
		t.Fatalf("bad migration: %s", out)
	}
	if !reflect.DeepEqual(server["enabled_tools"], []any{"hawp_search", "hawp_usage", "hawp_work_intake", "hawp_work_new", "hawp_work_validate", "hawp_work_doc", "hawp_work_reshape"}) {
		t.Fatalf("Codex tool allowlist was not upgraded: %s", out)
	}
	if !reflect.DeepEqual(server["args"], []any{"mcp", "--no-update-check", "--repo-root", root}) {
		t.Fatalf("lost args: %s", out)
	}
	again, err := mergeCodexTOML(out, root)
	if err != nil || !bytes.Equal(out, again) {
		t.Fatalf("not idempotent: %v", err)
	}
}

func TestCodexMigrationRefusesWithoutWrites(t *testing.T) {
	for _, input := range []string{
		"[broken", "[mcp_servers.hawp]\nx=1\nx=2\n",
		"[mcp_servers.hawp]\nurl='https://example.invalid'\n",
		"mcp_servers.hawp.command='old'\n",
		"mcp_servers = { hawp = { command = 'old' } }\n",
		"[mcp_servers.hawp]\ncommand=5\n",
		"[mcp_servers.hawp]\nargs=['other']\n",
		"[mcp_servers.hawp]\nargs=['mcp', 'configure']\n",
		"[mcp_servers.hawp]\nargs=['mcp', '--repo-root']\n",
		"[mcp_servers.hawp]\nargs=['mcp', # preserve this\n]\n",
	} {
		path := filepath.Join(t.TempDir(), "config.toml")
		if err := os.WriteFile(path, []byte(input), 0o600); err != nil {
			t.Fatal(err)
		}
		if err := writeCodexTOML(path, t.TempDir()); err == nil {
			t.Fatalf("accepted %q", input)
		}
		got, err := os.ReadFile(path)
		if err != nil || string(got) != input {
			t.Fatalf("changed rejected file: %v", err)
		}
	}
}

func TestCodexMigrationExistingRootAndEmptyTable(t *testing.T) {
	for _, input := range []string{"[mcp_servers.hawp]", "[mcp_servers.hawp]\nargs=['mcp','--repo-root=old']\n", "[mcp_servers.hawp]\nargs=['mcp','--repo-root','old']\n"} {
		out, err := mergeCodexTOML([]byte(input), "new root")
		if err != nil {
			t.Fatal(err)
		}
		again, err := mergeCodexTOML(out, "new root")
		if err != nil || !bytes.Equal(out, again) {
			t.Fatalf("not idempotent: %s, %v", out, err)
		}
	}
}
