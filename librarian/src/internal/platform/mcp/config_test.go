package mcp

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestWriteCodexTOMLCreatesFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".codex/config.toml")

	if err := writeCodexTOML(path); err != nil {
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
	if !strings.Contains(content, `command = ".hawp/bin/hawp"`) {
		t.Errorf(".codex/config.toml missing command, got:\n%s", content)
	}
	if !strings.Contains(content, `cwd = "."`) {
		t.Errorf(".codex/config.toml missing cwd, got:\n%s", content)
	}
}

func TestWriteCodexTOMLIdempotent(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".codex/config.toml")

	if err := writeCodexTOML(path); err != nil {
		t.Fatal(err)
	}
	if err := writeCodexTOML(path); err != nil {
		t.Fatal(err)
	}

	data, _ := os.ReadFile(path)
	count := strings.Count(string(data), "[mcp_servers.hawp]")
	if count != 1 {
		t.Errorf("expected 1 [mcp_servers.hawp] block, got %d:\n%s", count, string(data))
	}
}

func TestWriteCodexTOMLAppendsToExisting(t *testing.T) {
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, ".codex"), 0o755); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(dir, ".codex/config.toml")

	existing := "[model]\nname = \"o4-mini\"\n"
	if err := os.WriteFile(path, []byte(existing), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := writeCodexTOML(path); err != nil {
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
	dir := t.TempDir()
	path := filepath.Join(dir, ".mcp.json")

	if err := writeMCPJSON(path); err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), `"hawp"`) {
		t.Errorf("expected hawp key in .mcp.json, got:\n%s", string(data))
	}
}

func TestWriteMCPJSONIdempotent(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".mcp.json")

	if err := writeMCPJSON(path); err != nil {
		t.Fatal(err)
	}
	before, _ := os.ReadFile(path)
	if err := writeMCPJSON(path); err != nil {
		t.Fatal(err)
	}
	after, _ := os.ReadFile(path)
	if string(before) != string(after) {
		t.Errorf("writeMCPJSON not idempotent:\nbefore:\n%s\nafter:\n%s", before, after)
	}
}
