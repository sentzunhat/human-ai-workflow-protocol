package mcp

import "fmt"

// Validate the entire selection before configuration writes. Preserve the
// historical meaning of all: GitHub remains an explicitly requested advisory.
func expandConfigProviders(providers []string) ([]string, error) {
	var result []string
	seen := map[string]bool{}
	for _, provider := range providers {
		names := []string{provider}
		switch provider {
		case "all":
			names = []string{"claude", "cursor", "continue", "codex"}
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
