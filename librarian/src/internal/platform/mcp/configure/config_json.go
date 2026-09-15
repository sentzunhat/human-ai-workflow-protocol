package mcp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Preserve unrelated settings as raw JSON, including large numeric values.
// Only the supplied HAWP launch fields belong to this migration.
func mergeMCPJSON(data []byte, entry map[string]any) ([]byte, error) {
	config := map[string]json.RawMessage{}
	if len(data) > 0 {
		if err := decodeConfigObject(data, &config); err != nil {
			return nil, fmt.Errorf("config: %w", err)
		}
	}
	servers := map[string]json.RawMessage{}
	if raw, ok := config["mcpServers"]; ok {
		if err := decodeConfigObject(raw, &servers); err != nil {
			return nil, fmt.Errorf("mcpServers: %w", err)
		}
	}
	hawp := map[string]json.RawMessage{}
	if raw, ok := servers["hawp"]; ok {
		if err := decodeConfigObject(raw, &hawp); err != nil {
			return nil, fmt.Errorf("mcpServers.hawp: %w", err)
		}
	}
	// A remote server cannot be converted to stdio by merely changing command.
	if _, ok := hawp["url"]; ok {
		return nil, fmt.Errorf("mcpServers.hawp has a remote URL; review manually")
	}
	if raw, ok := hawp["type"]; ok {
		var kind string
		if err := json.Unmarshal(raw, &kind); err != nil || kind != "stdio" {
			return nil, fmt.Errorf("mcpServers.hawp has a non-stdio type; review manually")
		}
	}
	for key, value := range entry {
		raw, err := json.Marshal(value)
		if err != nil {
			return nil, err
		}
		hawp[key] = raw
	}
	var err error
	if servers["hawp"], err = json.Marshal(hawp); err != nil {
		return nil, err
	}
	if config["mcpServers"], err = json.Marshal(servers); err != nil {
		return nil, err
	}
	out, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return nil, err
	}
	return append(out, '\n'), nil
}

func decodeConfigObject(data []byte, target *map[string]json.RawMessage) error {
	if err := json.Unmarshal(data, target); err != nil {
		return err
	}
	if *target == nil {
		return fmt.Errorf("expected a JSON object, not null")
	}
	return nil
}

func writeMCPJSON(path string, entry map[string]any) error {
	data, err := os.ReadFile(path)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	if err == nil && len(bytes.TrimSpace(data)) == 0 {
		return fmt.Errorf("parse %s: empty configuration; review manually", path)
	}
	out, err := mergeMCPJSON(data, entry)
	if err != nil {
		return fmt.Errorf("parse %s: %w", path, err)
	}
	if bytes.Equal(data, out) {
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	// Reject symlinks and non-regular files to prevent redirected writes.
	if fi, err := os.Lstat(path); err == nil && !fi.Mode().IsRegular() {
		return fmt.Errorf("refusing to write %s: not a regular file", path)
	}
	return os.WriteFile(path, out, 0o644)
}
