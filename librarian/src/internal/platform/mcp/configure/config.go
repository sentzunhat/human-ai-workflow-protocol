package mcp

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/sentzunhat/hawp/librarian/src/internal/infrastructure/filesystem"
)

func hawpBinaryPath(repoRoot string) string {
	name := "hawp"
	if runtime.GOOS == "windows" {
		name += ".exe"
	}
	return filepath.Join(repoRoot, ".hawp", "bin", name)
}

// codexTOMLBlock returns the TOML snippet appended to .codex/config.toml for the
// hawp server. Codex does not resolve relative command paths or cwd relative to
// the project root for project-scoped configs, so both must be absolute.
func codexTOMLBlock(repoRoot string) string {
	bin := hawpBinaryPath(repoRoot)
	return fmt.Sprintf("\n[mcp_servers.hawp]\ncommand = %q\nargs = [\"mcp\", \"--repo-root\", %q]\ncwd = %q\nenabled = true\nstartup_timeout_sec = 30\ntool_timeout_sec = 60\nenabled_tools = [\"hawp_search\", \"hawp_usage\", \"hawp_work_intake\", \"hawp_work_new\", \"hawp_work_validate\"]\n", bin, repoRoot, repoRoot)
}

// claudeServerEntry returns the MCP server config block written for Claude
// Code. We write an absolute command path so the config remains reliable
// regardless of the app's current working directory; each machine/project
// computes its own path at init/update time.
func claudeServerEntry(repoRoot string) map[string]any {
	return map[string]any{
		"command": hawpBinaryPath(repoRoot),
		"args":    []string{"mcp", "--repo-root", repoRoot},
	}
}

// Cursor launches the native executable with an explicit repository root,
// independent of its spawn directory and without requiring a shell wrapper.
func cursorServerEntry(repoRoot string) map[string]any {
	return map[string]any{
		"type":    "stdio",
		"command": hawpBinaryPath(repoRoot),
		"args":    []string{"mcp", "--repo-root", repoRoot},
	}
}

// WriteProviderConfigs writes (or merges) the hawp MCP server entry into the
// relevant provider config file for each named provider. Existing HAWP launch
// fields are upgraded; JSON provider-specific settings are preserved.
//
// File-writing providers: claude (.mcp.json), cursor (.cursor/mcp.json), codex (.codex/config.toml).
// Continue prints a manual config block (no standard file location).
// github/Copilot prints a note (VS Code manages its own MCP config).
// WriteProviderConfigs validates and writes provider configs for each named provider.
// It expands shorthand names (e.g. "all") before writing. Callers that have already
// validated and expanded the list should call writeProviderConfigs directly to avoid
// the redundant expansion.
func WriteProviderConfigs(repoRoot string, providers []string) error {
	expanded, err := expandConfigProviders(providers)
	if err != nil {
		return err
	}
	return writeProviderConfigs(repoRoot, expanded)
}

// writeProviderConfigs writes configs for an already-validated, already-expanded provider list.
// Writes are not transactional: if a provider write fails, earlier providers already written
// are not rolled back. The error message names the failed provider so the caller can see
// which configs were written before the failure.
func writeProviderConfigs(repoRoot string, providers []string) error {
	var written []string
	wrapErr := func(provider string, err error) error {
		if len(written) == 0 {
			return fmt.Errorf("%s MCP config: %w", provider, err)
		}
		return fmt.Errorf("%s MCP config: %w (already written: %s)", provider, err, strings.Join(written, ", "))
	}
	for _, p := range providers {

		switch p {
		case "claude":
			if err := writeMCPJSON(filepath.Join(repoRoot, ".mcp.json"), claudeServerEntry(repoRoot)); err != nil {
				return wrapErr("claude", err)
			}
			written = append(written, "claude")
			fmt.Println("MCP: wrote .mcp.json (Claude Code)")

		case "cursor":
			dir := filepath.Join(repoRoot, ".cursor")
			if err := filesystem.RejectSymlinkAncestors(repoRoot, dir); err != nil {
				return wrapErr("cursor", err)
			}
			if err := os.MkdirAll(dir, 0o755); err != nil {
				return wrapErr("cursor", err)
			}
			if err := writeMCPJSON(filepath.Join(dir, "mcp.json"), cursorServerEntry(repoRoot)); err != nil {
				return wrapErr("cursor", err)
			}
			written = append(written, "cursor")
			fmt.Println("MCP: wrote .cursor/mcp.json (Cursor)")
			fmt.Println("  Cursor expects `type: stdio` and an absolute command path.")
			fmt.Println("  Open Customize → MCPs and enable the workspace server.")
			fmt.Println("  If tools still do not appear, open a new chat or reload Cursor.")

		case "continue":
			fmt.Println("MCP (Continue): add this block to .continue/config.yaml:")
			fmt.Println("  mcp:")
			fmt.Println("    servers:")
			fmt.Println("      - name: hawp")
			fmt.Printf("        command: %q\n", hawpBinaryPath(repoRoot))
			fmt.Printf("        args: [mcp, --repo-root, %q]\n", repoRoot)

		case "codex":
			codexDir := filepath.Join(repoRoot, ".codex")
			if err := filesystem.RejectSymlinkAncestors(repoRoot, codexDir); err != nil {
				return wrapErr("codex", err)
			}
			if err := os.MkdirAll(codexDir, 0o755); err != nil {
				return wrapErr("codex", err)
			}
			if err := writeCodexTOML(filepath.Join(codexDir, "config.toml"), repoRoot); err != nil {
				return wrapErr("codex", err)
			}
			written = append(written, "codex")
			fmt.Println("MCP: wrote .codex/config.toml (Codex)")
			fmt.Println("  Note: Codex only loads project MCP config for trusted projects.")
			fmt.Println("  Trust this repo in Codex settings, then start a fresh task/session.")
			fmt.Println("  CLI: `codex mcp list` confirms whether hawp is visible.")

		case "github":
			fmt.Println("MCP (github/Copilot): configure via VS Code MCP panel or .vscode/mcp.json — no file written.")

		}
	}
	return nil
}
