package mcp

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestMCPJSONMigrationPreservesSettings(t *testing.T) {
	input := []byte(`{"custom":9007199254740993,"mcpServers":{"other":{"command":"keep"},"hawp":{"command":"old","args":["mcp"],"env":{"LOCAL":"keep"},"timeout":123,"disabled":true}}}`)
	entry := cursorServerEntry(t.TempDir())
	out, err := mergeMCPJSON(input, entry)
	if err != nil {
		t.Fatal(err)
	}
	var config map[string]json.RawMessage
	if err := json.Unmarshal(out, &config); err != nil {
		t.Fatal(err)
	}
	if string(config["custom"]) != "9007199254740993" {
		t.Fatalf("unrelated number changed: %s", out)
	}
	var servers map[string]map[string]json.RawMessage
	if err := json.Unmarshal(config["mcpServers"], &servers); err != nil {
		t.Fatal(err)
	}
	for key, want := range map[string]string{"timeout": "123", "disabled": "true"} {
		if string(servers["hawp"][key]) != want {
			t.Fatalf("lost %s: %s", key, out)
		}
	}
	if !bytes.Contains(servers["hawp"]["env"], []byte(`"keep"`)) || string(servers["other"]["command"]) != `"keep"` {
		t.Fatalf("lost settings: %s", out)
	}
	for key, value := range entry {
		want, _ := json.Marshal(value)
		var compact bytes.Buffer
		if err := json.Compact(&compact, servers["hawp"][key]); err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(compact.Bytes(), want) {
			t.Fatalf("launch field %s not migrated: %s", key, out)
		}
	}
	again, err := mergeMCPJSON(out, entry)
	if err != nil || !bytes.Equal(out, again) {
		t.Fatalf("not idempotent: %v", err)
	}
}

func TestMCPJSONInvalidMigrationDoesNotWrite(t *testing.T) {
	for _, input := range []string{
		"", " ", "null", "[]", "{", `{"mcpServers":null}`,
		`{"mcpServers":[]}`, `{"mcpServers":{"hawp":null}}`,
		`{"mcpServers":{"hawp":"custom"}}`,
		`{"mcpServers":{"hawp":{"url":"https://example.invalid/mcp"}}}`,
		`{"mcpServers":{"hawp":{"type":"sse"}}}`,
	} {
		t.Run(input, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "config.json")
			if err := os.WriteFile(path, []byte(input), 0o600); err != nil {
				t.Fatal(err)
			}
			if err := writeMCPJSON(path, claudeServerEntry(t.TempDir())); err == nil {
				t.Fatal("accepted invalid or incompatible config")
			}
			got, err := os.ReadFile(path)
			if err != nil || string(got) != input {
				t.Fatalf("changed rejected configuration: %q, %v", got, err)
			}
		})
	}
}
