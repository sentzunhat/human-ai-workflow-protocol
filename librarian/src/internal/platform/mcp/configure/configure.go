package mcp

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Configure performs config-only setup. Preflight every selected configuration
// before writing any; filesystem failures during writes are not transactional.
func Configure(root string, providers []string) error {
	names, err := expandConfigProviders(providers)
	if err != nil {
		return err
	}
	if len(names) == 0 {
		return fmt.Errorf("select at least one --provider")
	}
	for _, path := range []string{hawpBinaryPath(root), filepath.Join(root, ".hawp", "work", "BACKLOG.md")} {
		info, err := os.Stat(path)
		if err != nil {
			return err
		}
		if !info.Mode().IsRegular() {
			return fmt.Errorf("expected a regular file: %s", path)
		}
	}
	for _, name := range names {
		var path string
		var entry map[string]any
		switch name {
		case "codex":
			path = filepath.Join(root, ".codex", "config.toml")
		case "claude":
			path, entry = filepath.Join(root, ".mcp.json"), claudeServerEntry(root)
		case "cursor":
			path, entry = filepath.Join(root, ".cursor", "mcp.json"), cursorServerEntry(root)
		default:
			continue
		}
		info, err := os.Lstat(path)
		if os.IsNotExist(err) {
			continue
		}
		if err != nil {
			return err
		}
		if !info.Mode().IsRegular() {
			return fmt.Errorf("refusing non-regular configuration: %s", path)
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		if name == "codex" {
			if _, err := mergeCodexTOML(data, root); err != nil {
				return fmt.Errorf("%s: %w", path, err)
			}
		} else {
			if len(strings.TrimSpace(string(data))) == 0 {
				return fmt.Errorf("empty configuration requires review: %s", path)
			}
			if _, err := mergeMCPJSON(data, entry); err != nil {
				return fmt.Errorf("%s: %w", path, err)
			}
		}
	}
	return writeProviderConfigs(root, names)
}
