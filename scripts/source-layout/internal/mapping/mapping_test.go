package mapping

import "testing"

func TestDestination(t *testing.T) {
	for _, tc := range []struct{ source, target string }{
		{"internal/application/context/config.go", "internal/application/context/config.go"},
		{"internal/application/context/config_test.go", "internal/application/context/config_test.go"},
		{"internal/application/context/dedup.go", "internal/application/context/dedup/dedup.go"},
		{"internal/application/context/dedup_test.go", "internal/application/context/dedup/dedup_test.go"},
		{"internal/application/context/format.go", "internal/application/context/format.go"},
		{"internal/application/context/reshaper.go", "internal/application/context/reshaper.go"},
		{"internal/platform/mcp/config_codex.go", "internal/platform/mcp/configure/config_codex.go"},
		{"internal/platform/cli/work/new/command.go", "internal/platform/cli/work/new/command.go"},
	} {
		t.Run(tc.source, func(t *testing.T) {
			got, reason := Destination(tc.source)
			if got != tc.target || reason == "" {
				t.Fatalf("got %q (%q), want %q", got, reason, tc.target)
			}
			if again, _ := Destination(got); again != got {
				t.Fatalf("unstable destination: %s -> %s", got, again)
			}
		})
	}
}
