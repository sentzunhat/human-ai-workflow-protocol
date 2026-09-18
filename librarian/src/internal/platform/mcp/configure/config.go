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

// vscodeServerEntry uses the workspace MCP shape consumed by VS Code's
// GitHub/Copilot integration. Keep this repo-local so init never writes to a
// user's global VS Code profile.
func vscodeServerEntry(repoRoot string) map[string]any {
	return cursorServerEntry(repoRoot)
}

// WriteProviderConfigs writes (or merges) the hawp MCP server entry into the
// relevant provider config file for each named provider. Existing HAWP launch
// fields are upgraded; JSON provider-specific settings are preserved.
//
// File-writing providers: claude (.mcp.json), cursor (.cursor/mcp.json),
// codex (.codex/config.toml), and github (.vscode/mcp.json).
// Continue prints a manual config block because its working config is the
// user-scoped ~/.continue/config.yaml file.
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
			if err := ensureGitignoreEntry(repoRoot, ".mcp.json"); err != nil {
				return wrapErr("claude", err)
			}
			if err := writeMCPJSON(filepath.Join(repoRoot, ".mcp.json"), claudeServerEntry(repoRoot)); err != nil {
				return wrapErr("claude", err)
			}
			written = append(written, "claude")
			fmt.Println("MCP: wrote .mcp.json (Claude Code)")

		case "cursor":
			if err := ensureGitignoreEntry(repoRoot, ".cursor/mcp.json"); err != nil {
				return wrapErr("cursor", err)
			}
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
			fmt.Println("  mcpServers:")
			fmt.Println("    - name: hawp")
			fmt.Printf("      command: %q\n", hawpBinaryPath(repoRoot))
			fmt.Printf("      args: [mcp, --repo-root, %q]\n", repoRoot)

		case "github":
			if err := ensureGitignoreEntry(repoRoot, ".vscode/mcp.json"); err != nil {
				return wrapErr("github", err)
			}
			dir := filepath.Join(repoRoot, ".vscode")
			if err := filesystem.RejectSymlinkAncestors(repoRoot, dir); err != nil {
				return wrapErr("github", err)
			}
			if err := os.MkdirAll(dir, 0o755); err != nil {
				return wrapErr("github", err)
			}
			if err := writeServerJSON(filepath.Join(dir, "mcp.json"), vscodeServerEntry(repoRoot), "servers"); err != nil {
				return wrapErr("github", err)
			}
			written = append(written, "github")
			fmt.Println("MCP: wrote .vscode/mcp.json (GitHub/Copilot)")
			fmt.Println("  Open the VS Code MCP panel and enable the workspace server.")

		case "codex":
			if err := ensureGitignoreEntry(repoRoot, ".codex/config.toml"); err != nil {
				return wrapErr("codex", err)
			}
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

		}
	}
	return nil
}

// ensureGitignoreEntry protects generated provider configuration from being
// committed. Entries are repo-relative, exact, and idempotent; existing
// comments and unrelated rules are preserved.
func ensureGitignoreEntry(repoRoot, entry string) error {
	path := filepath.Join(repoRoot, ".gitignore")
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		data = nil
	} else if err != nil {
		return err
	}
	for _, line := range strings.Split(strings.ReplaceAll(string(data), "\r\n", "\n"), "\n") {
		if strings.TrimSpace(line) == entry {
			return nil
		}
	}
	content := string(data)
	if content != "" && !strings.HasSuffix(content, "\n") {
		content += "\n"
	}
	content += entry + "\n"
	return os.WriteFile(path, []byte(content), 0o644)
}
