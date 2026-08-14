package context

import (
	"os"
	"path/filepath"
	"testing"
)

func writeFixture(t *testing.T, root string, files map[string]string) {
	t.Helper()
	for rel, content := range files {
		path := filepath.Join(root, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
}

func TestAdapterPreservesCorpusPolicies(t *testing.T) {
	root := t.TempDir()
	writeFixture(t, root, map[string]string{
		".hawp/kit/README.md":              "# Kit\n\nKit root.\n",
		".hawp/kit/usage/init.md":          "# Init\n",
		".hawp/work/BACKLOG.md":            "# Backlog\n",
		".hawp/work/active/TASK-001.md":    "# active\n",
		".hawp/work/active/README.md":      "# ignored readme\n",
		".hawp/work/future/TASK-002.md":    "# ignored folder\n",
		".hawp/work/decisions/decision.md": "# decision\n",
	})

	adapter := NewAdapter()
	kit, err := adapter.ReadKit(root, filepath.Join(root, ".hawp", "kit"))
	if err != nil {
		t.Fatal(err)
	}
	if len(kit.Files) != 2 {
		t.Fatalf("kit files = %d, want 2", len(kit.Files))
	}

	work, err := adapter.ReadWork(root, filepath.Join(root, ".hawp", "work"))
	if err != nil {
		t.Fatal(err)
	}
	if len(work.Files) != 2 {
		t.Fatalf("work files = %d, want 2 (active and decisions)", len(work.Files))
	}
	for _, file := range work.Files {
		if filepath.Base(file.RelPath) == "README.md" {
			t.Error("work adapter should skip README.md")
		}
	}
}
