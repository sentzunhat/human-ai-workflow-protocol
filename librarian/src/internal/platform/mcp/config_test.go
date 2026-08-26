package mcp

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestWriteCodexTOMLCreatesFile(t *testing.T) {
	repoRoot := t.TempDir()
	path := filepath.Join(repoRoot, ".codex/config.toml")

	if err := writeCodexTOML(path, repoRoot); err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	content := string(data)
	if !strings.Contains(content, "[mcp_servers.hawp]") {
		t.Errorf(".codex/config.toml missing [mcp_servers.hawp], got:\n%s", content)
	}
	expectedBin := filepath.Join(repoRoot, ".hawp", "bin", "hawp")
	if !strings.Contains(content, `command = "`+expectedBin+`"`) {
		t.Errorf(".codex/config.toml missing absolute command path, got:\n%s", content)
	}
	if !strings.Contains(content, `cwd = "`+repoRoot+`"`) {
		t.Errorf(".codex/config.toml missing absolute cwd, got:\n%s", content)
	}
}

func TestWriteCodexTOMLIdempotent(t *testing.T) {
	repoRoot := t.TempDir()
	path := filepath.Join(repoRoot, ".codex/config.toml")

	if err := writeCodexTOML(path, repoRoot); err != nil {
		t.Fatal(err)
	}
	if err := writeCodexTOML(path, repoRoot); err != nil {
		t.Fatal(err)
	}

	data, _ := os.ReadFile(path)
	count := strings.Count(string(data), "[mcp_servers.hawp]")
	if count != 1 {
		t.Errorf("expected 1 [mcp_servers.hawp] block, got %d:\n%s", count, string(data))
	}
}

func TestWriteCodexTOMLUpgradesExistingRelativeBlock(t *testing.T) {
	repoRoot := t.TempDir()
	path := filepath.Join(repoRoot, ".codex/config.toml")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	existing := "[mcp_servers.hawp]\ncommand = \".hawp/bin/hawp\"\nargs = [\"mcp\"]\ncwd = \".\"\n"
	if err := os.WriteFile(path, []byte(existing), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := writeCodexTOML(path, repoRoot); err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	content := string(data)
	expectedBin := filepath.Join(repoRoot, ".hawp", "bin", "hawp")
	if !strings.Contains(content, `command = "`+expectedBin+`"`) {
		t.Errorf("expected codex config to upgrade command path, got:\n%s", content)
	}
	if !strings.Contains(content, `cwd = "`+repoRoot+`"`) {
		t.Errorf("expected codex config to upgrade cwd, got:\n%s", content)
	}
	if strings.Contains(content, `command = ".hawp/bin/hawp"`) || strings.Contains(content, `cwd = "."`) {
		t.Errorf("expected stale relative codex config to be removed, got:\n%s", content)
	}
}

func TestWriteCodexTOMLAppendsToExisting(t *testing.T) {
	repoRoot := t.TempDir()
	if err := os.MkdirAll(filepath.Join(repoRoot, ".codex"), 0o755); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(repoRoot, ".codex/config.toml")

	existing := "[model]\nname = \"o4-mini\"\n"
	if err := os.WriteFile(path, []byte(existing), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := writeCodexTOML(path, repoRoot); err != nil {
		t.Fatal(err)
	}

	data, _ := os.ReadFile(path)
	content := string(data)
	if !strings.Contains(content, "[model]") {
		t.Error("writeCodexTOML must preserve existing content")
	}
	if !strings.Contains(content, "[mcp_servers.hawp]") {
		t.Error("writeCodexTOML must add hawp block")
	}
}

func TestWriteMCPJSONCreatesFile(t *testing.T) {
	repoRoot := t.TempDir()
	path := filepath.Join(repoRoot, ".mcp.json")

	if err := writeMCPJSON(path, repoRoot); err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	content := string(data)
	if !strings.Contains(content, `"hawp"`) {
		t.Errorf("expected hawp key in .mcp.json, got:\n%s", string(data))
	}
	expectedBin := filepath.Join(repoRoot, ".hawp", "bin", "hawp")
	if !strings.Contains(content, `"command": "`+expectedBin+`"`) {
		t.Errorf(".mcp.json missing absolute command path, got:\n%s", content)
	}
}

func TestWriteMCPJSONUpgradesExistingRelativeCommand(t *testing.T) {
	repoRoot := t.TempDir()
	path := filepath.Join(repoRoot, ".mcp.json")
	existing := `{
  "mcpServers": {
    "hawp": {
      "command": ".hawp/bin/hawp",
      "args": ["mcp"]
    }
  }
}
`
	if err := os.WriteFile(path, []byte(existing), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := writeMCPJSON(path, repoRoot); err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	content := string(data)
	expectedBin := filepath.Join(repoRoot, ".hawp", "bin", "hawp")
	if !strings.Contains(content, `"command": "`+expectedBin+`"`) {
		t.Errorf("expected .mcp.json to upgrade command path, got:\n%s", content)
	}
	if strings.Contains(content, `"command": ".hawp/bin/hawp"`) {
		t.Errorf("expected stale relative command to be removed, got:\n%s", content)
	}
}

func TestWriteMCPJSONIdempotent(t *testing.T) {
	repoRoot := t.TempDir()
	path := filepath.Join(repoRoot, ".mcp.json")

	if err := writeMCPJSON(path, repoRoot); err != nil {
		t.Fatal(err)
	}
	before, _ := os.ReadFile(path)
	if err := writeMCPJSON(path, repoRoot); err != nil {
		t.Fatal(err)
	}
	after, _ := os.ReadFile(path)
	if string(before) != string(after) {
		t.Errorf("writeMCPJSON not idempotent:\nbefore:\n%s\nafter:\n%s", before, after)
	}
}
