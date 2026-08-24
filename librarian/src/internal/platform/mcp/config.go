package mcp

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

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
// Supported: claude, cursor (JSON .mcp.json/.cursor/mcp.json).
// Continue, codex, and github print a manual-config message instead.
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

		case "codex", "github":
			fmt.Printf("MCP (%s): no native MCP support yet — use `hawp` CLI commands directly.\n", p)

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
