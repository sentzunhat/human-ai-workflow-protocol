package mcp

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/sentzunhat/hawp/librarian/src/internal/infrastructure/filesystem"
)

// Configure performs config-only setup. The public wrapper preflights every
// selected configuration before writing any; filesystem failures during writes
// are not transactional.
func Configure(root string, providers []string) error {
	return WriteProviderConfigs(root, providers)
}

// preflightProviderConfigs validates every requested provider configuration
// before writeProviderConfigs begins its non-transactional writes.
func preflightProviderConfigs(root string, names []string) error {
	if len(names) == 0 {
		return fmt.Errorf("select at least one --provider")
	}
	for _, path := range []string{hawpBinaryPath(root), filepath.Join(root, ".hawp", "work", "BACKLOG.md")} {
		if err := filesystem.RejectSymlinkAncestors(root, path); err != nil {
			return fmt.Errorf("unsafe prerequisite %s: %w", path, err)
		}
		info, err := os.Lstat(path)
		if err != nil {
			return err
		}
		if !info.Mode().IsRegular() {
			return fmt.Errorf("expected a regular file: %s", path)
		}
	}
	needsGitignore := false
	for _, name := range names {
		if name != "continue" {
			needsGitignore = true
			break
		}
	}
	if needsGitignore {
		gitignorePath := filepath.Join(root, ".gitignore")
		if err := filesystem.RejectSymlinkAncestors(root, gitignorePath); err != nil {
			return fmt.Errorf("unsafe configuration prerequisite %s: %w", gitignorePath, err)
		}
		info, err := os.Lstat(gitignorePath)
		if err == nil && !info.Mode().IsRegular() {
			return fmt.Errorf("refusing non-regular configuration prerequisite: %s", gitignorePath)
		}
		if err != nil && !os.IsNotExist(err) {
			return err
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
		case "github":
			path, entry = filepath.Join(root, ".vscode", "mcp.json"), vscodeServerEntry(root)
		default:
			continue
		}
		if err := filesystem.RejectSymlinkAncestors(root, path); err != nil {
			return fmt.Errorf("unsafe configuration %s: %w", path, err)
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
		} else if name == "github" {
			if len(strings.TrimSpace(string(data))) == 0 {
				return fmt.Errorf("empty configuration requires review: %s", path)
			}
			if _, err := mergeServerJSON(data, entry, "servers"); err != nil {
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
	return nil
}
