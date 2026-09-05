package mcp

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
	"github.com/sentzunhat/hawp/librarian/src/internal/infrastructure/tomlconfig"
)

func mergeCodexTOML(data []byte, root string) ([]byte, error) {
	var doc map[string]any
	if err := toml.Unmarshal(data, &doc); err != nil {
		return nil, fmt.Errorf("invalid Codex TOML; preserved configuration")
	}
	var defaults map[string]any
	if err := toml.Unmarshal([]byte(codexTOMLBlock(root)), &defaults); err != nil {
		return nil, fmt.Errorf("cannot encode repository path for Codex")
	}
	want := defaults["mcp_servers"].(map[string]any)["hawp"].(map[string]any)
	var server map[string]any
	if raw, exists := doc["mcp_servers"]; exists {
		servers, ok := raw.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("mcp_servers must be a table")
		}
		if raw, exists := servers["hawp"]; exists {
			server, ok = raw.(map[string]any)
			if !ok {
				return nil, fmt.Errorf("HAWP server must be a table")
			}
		}
	}
	if server == nil {
		out := append(append([]byte(nil), data...), []byte(codexTOMLBlock(root))...)
		var check map[string]any
		if err := toml.Unmarshal(out, &check); err != nil {
			return nil, fmt.Errorf("Codex table layout requires manual merge")
		}
		return out, nil
	}
	if _, exists := server["url"]; exists {
		return nil, fmt.Errorf("remote HAWP server requires manual merge")
	}
	if kind, exists := server["type"]; exists && kind != "stdio" {
		return nil, fmt.Errorf("non-stdio HAWP server requires manual merge")
	}
	if command, exists := server["command"]; exists {
		if _, ok := command.(string); !ok {
			return nil, fmt.Errorf("HAWP command must be a string")
		}
	}
	args, err := codexLaunchArgs(server["args"], root)
	if err != nil {
		return nil, err
	}
	fields := map[string]any{"command": want["command"], "args": args, "cwd": root}
	for key, value := range want {
		if _, exists := server[key]; !exists {
			if _, overridden := fields[key]; !overridden {
				fields[key] = value
			}
		}
	}
	return tomlconfig.Rewrite(data, []string{"mcp_servers", "hawp"}, fields)
}

func writeCodexTOML(path, root string) error {
	info, err := os.Lstat(path)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	if err == nil && !info.Mode().IsRegular() {
		return fmt.Errorf("refusing non-regular Codex configuration")
	}
	data, err := os.ReadFile(path)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	out, err := mergeCodexTOML(data, root)
	if err != nil {
		return err
	}
	if bytes.Equal(data, out) {
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, out, 0o644)
}
