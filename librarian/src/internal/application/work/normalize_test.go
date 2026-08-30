package work

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNormalizeDryRunMigrateFoldersPreviewsChanges(t *testing.T) {
	root := t.TempDir()
	files := map[string]string{
		".hawp/work/BACKLOG.md": `# Backlog

## Active Work

| ID | Type | Title | Status | Plan File | Updated |
| --- | --- | --- | --- | --- | --- |
| legacy-item | task | legacy | inbox | [plan](active/legacy-item.md) | 2026-08-30 |

## Blocked / Parked

| ID | Type | Title | Reason | Detail | Updated |
| --- | --- | --- | --- | --- | --- |

## Recently Closed

| ID | Type | Title | Closed | Detail |
| --- | --- | --- | --- | --- |
`,
		".hawp/work/active/legacy-item.md": `# Legacy Item

**Plan file:** work/active/legacy-item.md
`,
		".hawp/work/parked/.keep": "",
		".hawp/work/closed/.keep": "",
	}
	for rel, content := range files {
		full := filepath.Join(root, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	var out bytes.Buffer
	var errOut bytes.Buffer
	code := Normalize(&out, &errOut, NormalizeOptions{
		RepoRoot:       root,
		MigrateFolders: true,
	})
	if code != 0 {
		t.Fatalf("Normalize() code = %d, stderr = %s", code, errOut.String())
	}
	if !strings.Contains(out.String(), "HAWP Work Folder Migration Dry-Run") {
		t.Fatalf("expected migration dry-run header, got %q", out.String())
	}
	if !strings.Contains(out.String(), ".hawp/work/active/legacy-item/plan.md") {
		t.Fatalf("expected migrated plan path in output, got %q", out.String())
	}
	if _, err := os.Stat(filepath.Join(root, ".hawp/work/active/legacy-item.md")); err != nil {
		t.Fatalf("dry-run should not mutate source repo, stat err = %v", err)
	}
}
