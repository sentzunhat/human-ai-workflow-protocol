package work

import (
	"io"
	"os"
	"path/filepath"
	"testing"
)

func TestNormalizeApplyMigrateFoldersUsesProductionFileCollector(t *testing.T) {
	root := t.TempDir()
	files := map[string]string{
		".hawp/work/BACKLOG.md": `# Backlog

## Active Work

| ID | Type | Title | Status | Plan File | Updated |
| --- | --- | --- | --- | --- | --- |
| ` + "`legacy-item`" + ` | task | legacy | inbox | [plan](active/legacy-item.md) | 2026-09-18 |

## Blocked / Parked

| ID | Type | Title | Reason | Detail | Updated |
| --- | --- | --- | --- | --- | --- |

## Recently Closed

| ID | Type | Title | Closed | Detail |
| --- | --- | --- | --- | --- |
`,
		".hawp/work/active/legacy-item.md": "# Legacy item\n\n**Plan file:** work/active/legacy-item.md\n",
	}
	for rel, content := range files {
		path := filepath.Join(root, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	code := Normalize(io.Discard, io.Discard, NormalizeOptions{
		RepoRoot:       root,
		Apply:          true,
		MigrateFolders: true,
		ForceDirty:     true,
	})
	if code != 0 {
		t.Fatalf("Normalize returned exit code %d", code)
	}
	if _, err := os.Stat(filepath.Join(root, ".hawp/work/active/legacy-item/plan.md")); err != nil {
		t.Fatalf("migrated plan missing: %v", err)
	}
}
