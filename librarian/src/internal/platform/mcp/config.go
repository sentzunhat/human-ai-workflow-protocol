package mcp

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func hawpBinaryPath(repoRoot string) string {
	return filepath.Join(repoRoot, ".hawp", "bin", "hawp")
}

// codexTOMLBlock returns the TOML snippet appended to .codex/config.toml for the
// hawp server. Codex does not resolve relative command paths or cwd relative to
// the project root for project-scoped configs, so both must be absolute.
func codexTOMLBlock(repoRoot string) string {
	bin := hawpBinaryPath(repoRoot)
	return "\n[mcp_servers.hawp]\ncommand = \"" + bin + "\"\nargs = [\"mcp\"]\ncwd = \"" + repoRoot + "\"\n"
}

// hawpServerEntry returns the MCP server config block written for JSON-RPC MCP
// providers that use JSON config files. We write an absolute command path so
// the config remains reliable regardless of the app's current working
// directory; each machine/project computes its own path at init/update time.
func hawpServerEntry(repoRoot string) map[string]any {
	return map[string]any{
		"command": hawpBinaryPath(repoRoot),
		"args":    []string{"mcp"},
	}
}

// WriteProviderConfigs writes (or merges) the hawp MCP server entry into the
// relevant provider config file for each named provider. Idempotent: if the
// hawp entry already exists it is not duplicated or overwritten.
//
// File-writing providers: claude (.mcp.json), cursor (.cursor/mcp.json), codex (.codex/config.toml).
// Continue prints a manual config block (no standard file location).
// github/Copilot prints a note (VS Code manages its own MCP config).
func WriteProviderConfigs(repoRoot string, providers []string) error {
	seen := map[string]bool{}
	for _, p := range providers {
		if seen[p] {
			continue
		}
		seen[p] = true

		switch p {
		case "claude":
			if err := writeMCPJSON(filepath.Join(repoRoot, ".mcp.json"), repoRoot); err != nil {
				return fmt.Errorf("claude MCP config: %w", err)
			}
			fmt.Println("MCP: wrote .mcp.json (Claude Code)")

		case "cursor":
			dir := filepath.Join(repoRoot, ".cursor")
			if err := os.MkdirAll(dir, 0o755); err != nil {
				return fmt.Errorf("cursor MCP config: %w", err)
			}
			if err := writeMCPJSON(filepath.Join(dir, "mcp.json"), repoRoot); err != nil {
				return fmt.Errorf("cursor MCP config: %w", err)
			}
			fmt.Println("MCP: wrote .cursor/mcp.json (Cursor)")

		case "continue":
			fmt.Println("MCP (Continue): add this block to .continue/config.yaml:")
			fmt.Println("  mcp:")
			fmt.Println("    servers:")
			fmt.Println("      - name: hawp")
			fmt.Println("        command: .hawp/bin/hawp")
			fmt.Println("        args: [mcp]")

		case "codex":
			codexDir := filepath.Join(repoRoot, ".codex")
			if err := os.MkdirAll(codexDir, 0o755); err != nil {
				return fmt.Errorf("codex MCP config: %w", err)
			}
			if err := writeCodexTOML(filepath.Join(codexDir, "config.toml"), repoRoot); err != nil {
				return fmt.Errorf("codex MCP config: %w", err)
			}
			fmt.Println("MCP: wrote .codex/config.toml (Codex)")
			fmt.Println("  Note: Codex only loads project MCP config for trusted projects.")
			fmt.Println("  Trust this repo in Codex settings, then start a fresh task/session.")
			fmt.Println("  CLI: `codex mcp list` confirms whether hawp is visible.")

		case "github":
			fmt.Println("MCP (github/Copilot): configure via VS Code MCP panel or .vscode/mcp.json — no file written.")

		case "all":
			// Expand "all" to the four main MCP-capable providers.
			sub := []string{"claude", "cursor", "continue", "codex"}
			if err := WriteProviderConfigs(repoRoot, sub); err != nil {
				return err
			}
		}
	}
	return nil
}

// writeCodexTOML appends the hawp MCP server entry to .codex/config.toml
// (creating the file if absent). Existing hawp blocks are replaced in place so
// relative or stale paths are upgraded to the current absolute-path form.
// We avoid a TOML library dependency by treating the file as plain text — the
// block format is fixed and isolated, so targeted replacement is sufficient.
func writeCodexTOML(path, repoRoot string) error {
	existing, err := os.ReadFile(path)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	content := upsertCodexTOML(string(existing), repoRoot)

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(content), 0o644)
}

func upsertCodexTOML(content, repoRoot string) string {
	block := strings.TrimPrefix(codexTOMLBlock(repoRoot), "\n")
	lines := strings.Split(content, "\n")
	start := -1
	end := -1
	for i, line := range lines {
		if line == "[mcp_servers.hawp]" {
			start = i
			end = len(lines)
			for j := i + 1; j < len(lines); j++ {
				if strings.HasPrefix(lines[j], "[") {
					end = j
					break
				}
			}
			break
		}
	}
	if start >= 0 {
		blockLines := strings.Split(strings.TrimSuffix(block, "\n"), "\n")
		updatedLines := append([]string{}, lines[:start]...)
		updatedLines = append(updatedLines, blockLines...)
		updatedLines = append(updatedLines, lines[end:]...)
		updated := strings.Join(updatedLines, "\n")
		if !strings.HasSuffix(updated, "\n") {
			updated += "\n"
		}
		return updated
	}
	if !strings.HasSuffix(content, "\n") && len(content) > 0 {
		content += "\n"
	}
	content += codexTOMLBlock(repoRoot)
	return content
}

// writeMCPJSON merges the hawp server entry into a JSON MCP config file.
// Creates the file when missing; patches it when present; upgrades any
// existing hawp entry to the current absolute-path form.
func writeMCPJSON(path, repoRoot string) error {
	config := map[string]any{}

	data, err := os.ReadFile(path)
	if err == nil {
		if err := json.Unmarshal(data, &config); err != nil {
			return fmt.Errorf("parse %s: %w", path, err)
		}
	} else if !os.IsNotExist(err) {
		return err
	}

	servers, _ := config["mcpServers"].(map[string]any)
	if servers == nil {
		servers = map[string]any{}
	}

	servers["hawp"] = hawpServerEntry(repoRoot)
	config["mcpServers"] = servers

	out, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}

	content := string(out)
	if !strings.HasSuffix(content, "\n") {
		content += "\n"
	}
	return os.WriteFile(path, []byte(content), 0o644)
}
