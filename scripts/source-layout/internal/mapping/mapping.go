// Package mapping defines the repository-specific, reviewed destination policy.
package mapping

import (
	"path"
	"strings"
)

// Revision invalidates saved plans when ownership decisions change.
const Revision = "2026-09-11-normalize-alignment"

// This is the reviewed mapping, not a heuristic based on filenames being unused.
// Files outside these groups are inventoried and retained, including assets/tests.
func Destination(p string) (string, string) {
	// These pairs own self-contained helpers/tests. Rendering and reshaping
	// remain together because they still share private reference helpers.
	if path.Dir(p) == "internal/application/context" {
		folder := ""
		switch path.Base(p) {
		case "dedup.go", "dedup_test.go":
			folder = "dedup"
		}
		if folder != "" {
			return path.Join(path.Dir(p), folder, path.Base(p)), "context capability: " + folder
		}
	}
	whole := []struct{ from, to, reason string }{
		{"internal/domain/embeddings/", "internal/domain/providers/embeddings/", "embedding consumer port"},
		{"internal/domain/llm/", "internal/domain/providers/llm/", "LLM consumer port"},
		{"internal/infrastructure/sqlite/", "internal/infrastructure/repositories/index/", "index persistence adapter"},
		{"internal/infrastructure/githubrelease/", "internal/infrastructure/clients/githubrelease/", "release service client"},
		{"internal/infrastructure/download/", "internal/infrastructure/clients/download/", "HTTP download client"},
		{"tests/infrastructure/githubrelease/", "tests/infrastructure/clients/githubrelease/", "release client tests follow owner"},
	}
	for _, r := range whole {
		if strings.HasPrefix(p, r.from) {
			return r.to + strings.TrimPrefix(p, r.from), r.reason
		}
	}
	if strings.HasPrefix(p, "internal/platform/mcp/") && strings.Count(p, "/") == 3 {
		folder := "server"
		if strings.HasPrefix(path.Base(p), "config") {
			folder = "configure"
		}
		return path.Join(path.Dir(p), folder, path.Base(p)), "MCP " + folder + " and adjacent tests"
	}
	if strings.HasPrefix(p, "internal/application/work/") && strings.Count(p, "/") == 3 {
		folder := ""
		switch path.Base(p) {
		case "intake.go", "draft.go", "draft_test.go":
			folder = "intake"
		case "validate.go":
			folder = "validation"
		case "normalize.go":
			folder = "normalize"
		}
		if folder != "" {
			return path.Join(path.Dir(p), folder, path.Base(p)), "work use case: " + folder
		}
	}
	if p == "internal/domain/integration_test.go" || p == "internal/domain/benchmark_test.go" {
		return path.Join("internal/infrastructure/models", path.Base(p)), "live model tests follow adapters; build tags retained"
	}
	return p, "retain existing cohesive owner or module support file"
}
