package providersync

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestComputeOutputsRendersEveryTarget(t *testing.T) {
	root := t.TempDir()
	shared := filepath.Join(root, "core", "providers", "shared", "behaviors")
	if err := os.MkdirAll(shared, 0o755); err != nil {
		t.Fatal(err)
	}

	seen := map[string]bool{}
	for _, target := range MaterializationTargets {
		if seen[target.Behavior] {
			continue
		}
		seen[target.Behavior] = true
		path := filepath.Join(shared, target.Behavior+".md")
		if err := os.WriteFile(path, []byte("# "+target.Behavior+"\n\nbody\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	outputs, err := ComputeOutputs(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(outputs) != len(MaterializationTargets) {
		t.Fatalf("outputs = %d, want %d", len(outputs), len(MaterializationTargets))
	}

	for _, output := range outputs {
		if !strings.HasPrefix(output.Content, "---\n") {
			t.Fatalf("frontmatter missing for %s", output.OutputPath)
		}
		if !strings.Contains(output.Content, strings.TrimSpace(GeneratedBanner)) {
			t.Fatalf("banner missing for %s", output.OutputPath)
		}
		if !strings.HasSuffix(output.Content, "\n") {
			t.Fatalf("trailing newline missing for %s", output.OutputPath)
		}
	}
}
