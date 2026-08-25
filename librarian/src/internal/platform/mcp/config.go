package mcp

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// codexTOMLBlock is the TOML snippet appended to .codex/config.toml for the hawp server.
const codexTOMLBlock = `
[mcp_servers.hawp]
command = ".hawp/bin/hawp"
args = ["mcp"]
cwd = "."
`

// hawpServerEntry is the MCP server config block written for all providers
// that support JSON-RPC MCP (Claude Code, Cursor).
var hawpServerEntry = map[string]any{
	"command": ".hawp/bin/hawp",
	"args":    []string{"mcp"},
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
			if err := writeMCPJSON(filepath.Join(repoRoot, ".mcp.json")); err != nil {
				return fmt.Errorf("claude MCP config: %w", err)
			}
			fmt.Println("MCP: wrote .mcp.json (Claude Code)")

		case "cursor":
			dir := filepath.Join(repoRoot, ".cursor")
			if err := os.MkdirAll(dir, 0o755); err != nil {
				return fmt.Errorf("cursor MCP config: %w", err)
			}
			if err := writeMCPJSON(filepath.Join(dir, "mcp.json")); err != nil {
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
			if err := writeCodexTOML(filepath.Join(codexDir, "config.toml")); err != nil {
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
// (creating the file if absent). Idempotent: skips when an [mcp_servers.hawp]
// block is already present. We avoid a TOML library dependency by treating the
// file as plain text — the block format is fixed, so a substring check is sufficient.
func writeCodexTOML(path string) error {
	existing, err := os.ReadFile(path)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	if strings.Contains(string(existing), "[mcp_servers.hawp]") {
		return nil // already configured
	}

	content := string(existing)
	if !strings.HasSuffix(content, "\n") && len(content) > 0 {
		content += "\n"
	}
	content += codexTOMLBlock

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(content), 0o644)
}

// writeMCPJSON merges the hawp server entry into a JSON MCP config file.
// Creates the file when missing; patches it when present; no-ops when the
// hawp entry already exists.
func writeMCPJSON(path string) error {
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

	if _, exists := servers["hawp"]; exists {
		return nil // already configured
	}

	servers["hawp"] = hawpServerEntry
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
