package work

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

// buildDuplicateFixture creates a repo root with a .hawp/work tree for duplicate-link tests.
func buildDuplicateFixture(t *testing.T, files map[string]string) string {
	t.Helper()
	root := t.TempDir()
	for rel, content := range files {
		full := filepath.Join(root, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return root
}

const duplicateBacklogHeader = `# Backlog

## Active Work

| ID | Type | Title | Status | Plan File | Updated |
| --- | --- | --- | --- | --- | --- |
`

const duplicateBacklogFooter = `
## Blocked / Parked

| ID | Type | Title | Reason | Detail | Updated |
| --- | --- | --- | --- | --- | --- |

## Recently Closed

| ID | Type | Title | Closed | Detail |
| --- | --- | --- | --- | --- |
`

const duplicateClosedPlanComplete = `# done thing

## Outcome (filled at close)

Done.

## Verification (filled at close)

- [x] It works (Evidence: test run output)

## Close Checklist

- [x] Outcome section filled
`

func TestDuplicateLinksPreservesDifferentContentAndArtifacts(t *testing.T) {
	root := buildDuplicateFixture(t, map[string]string{
		".hawp/work/BACKLOG.md":                    duplicateBacklogHeader + duplicateBacklogFooter,
		".hawp/work/active/BUG-901.md":             "# working changes\n",
		".hawp/work/closed/2026/09/05/BUG-901.md":  "# archive\n",
		".hawp/work/active/abcd1234/plan.md":       "# plan\n\n**UUID:** abcd1234\n",
		".hawp/work/active/abcd1234/evidence.txt":  "unique supporting evidence",
		".hawp/work/closed/2026/09/05/abcd1234.md": "# plan\n\n**UUID:** abcd1234\n",
		".hawp/work/parked/BUG-902.md":             "# identical\n",
		".hawp/work/closed/2026/09/05/BUG-902.md":  "# identical\n",
	})
	preview, err := PreviewDuplicateLinks(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(preview.ChangedFiles) != 6 || len(preview.ReviewFiles) != 0 {
		t.Fatalf("unexpected preview: %+v", preview)
	}
	if _, err := os.Stat(filepath.Join(root, ".hawp/work/parked/BUG-902.md")); err != nil {
		t.Fatal("preview mutated files")
	}
	result, err := ApplyDuplicateLinks(root)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(preview, result) {
		t.Fatalf("preview/apply mismatch: %+v / %+v", preview, result)
	}
	for _, path := range []string{"active/BUG-901.md", "active/abcd1234/plan.md", "active/abcd1234/evidence.txt"} {
		if _, err := os.Stat(filepath.Join(root, ".hawp/work", path)); err != nil {
			t.Fatalf("lost %s: %v", path, err)
		}
	}
	if _, err := os.Stat(filepath.Join(root, ".hawp/work/parked/BUG-902.md")); err != nil {
		t.Fatal("identical unreferenced copy not preserved")
	}
	for _, path := range result.ChangedFiles {
		data, err := os.ReadFile(filepath.Join(root, path))
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(string(data), "## Related Work Records") {
			t.Fatalf("missing link in %s", path)
		}
	}
	data, err := os.ReadFile(filepath.Join(root, ".hawp/work/active/abcd1234/evidence.txt"))
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "unique supporting evidence" {
		t.Fatal("supporting evidence changed")
	}
	again, err := ApplyDuplicateLinks(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(again.ChangedFiles) != 0 {
		t.Fatalf("cleanup not idempotent: %+v", again)
	}
}

func TestDuplicateLinksLeaveAmbiguousArchivesUntouched(t *testing.T) {
	root := buildDuplicateFixture(t, map[string]string{
		".hawp/work/BACKLOG.md":                   duplicateBacklogHeader + duplicateBacklogFooter,
		".hawp/work/active/BUG-901.md":            "# original\n",
		".hawp/work/closed/2026/09/04/BUG-901.md": "# first archive\n",
		".hawp/work/closed/2026/09/05/BUG-901.md": "# second archive\n",
	})
	result, err := ApplyDuplicateLinks(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.ChangedFiles) != 0 || len(result.ReviewFiles) != 1 {
		t.Fatalf("unexpected result: %+v", result)
	}
	data, err := os.ReadFile(filepath.Join(root, ".hawp/work/active/BUG-901.md"))
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "# original\n" {
		t.Fatal("ambiguous working plan changed")
	}
}

func TestApplyDuplicateLinksPreservesUnreferencedWorkingCopy(t *testing.T) {
	root := buildDuplicateFixture(t, map[string]string{
		".hawp/work/BACKLOG.md": duplicateBacklogHeader +
			"| BUG-250 | bug | reopened item | in-progress | [plan](active/BUG-250.md) | 2026-09-01 |\n" + duplicateBacklogFooter +
			"| BUG-249 | bug | closed duplicate | 2026-09-01 | [plan](closed/2026/09/01/BUG-249.md) |\n",
		".hawp/work/active/BUG-249.md":            duplicateClosedPlanComplete,
		".hawp/work/active/BUG-250.md":            "# active copy\n\n**Backlog ID:** BUG-250\n",
		".hawp/work/closed/2026/09/01/BUG-249.md": duplicateClosedPlanComplete,
		".hawp/work/closed/2026/09/01/BUG-250.md": duplicateClosedPlanComplete,
	})
	result, err := ApplyDuplicateLinks(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.ChangedFiles) != 2 {
		t.Fatalf("changed files = %v, want two cross-referenced plans", result.ChangedFiles)
	}
	if _, err := os.Stat(filepath.Join(root, ".hawp/work/active/BUG-249.md")); err != nil {
		t.Fatalf("working copy should be preserved: %v", err)
	}
	if _, err := os.Stat(filepath.Join(root, ".hawp/work/active/BUG-250.md")); err != nil {
		t.Fatalf("referenced active copy should be preserved: %v", err)
	}
}
