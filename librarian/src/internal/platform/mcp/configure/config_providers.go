package mcp

import "fmt"

// Validate the entire selection before configuration writes. "all" covers
// every supported provider integration, including the repo-local VS Code
// configuration used by GitHub/Copilot.
func expandConfigProviders(providers []string) ([]string, error) {
	var result []string
	seen := map[string]bool{}
	for _, provider := range providers {
		names := []string{provider}
		switch provider {
		case "all":
			names = []string{"claude", "cursor", "continue", "codex", "github"}
		case "claude", "cursor", "continue", "codex", "github":
		default:
			return nil, fmt.Errorf("unknown MCP provider %q; use claude, cursor, continue, codex, github, or all", provider)
		}
		for _, name := range names {
			if !seen[name] {
				seen[name] = true
				result = append(result, name)
			}
		}
	}
	return result, nil
}
