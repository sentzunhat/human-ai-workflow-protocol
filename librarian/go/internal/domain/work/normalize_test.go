package work

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// buildRepoFixture creates a repo root with a .hawp/work tree.
func buildRepoFixture(t *testing.T, files map[string]string) string {
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

func detect(t *testing.T, root string) []FixOperation {
	t.Helper()
	workRoot := filepath.Join(root, ".hawp", "work")
	backlog, err := ParseNormalizeBacklog(filepath.Join(workRoot, "BACKLOG.md"))
	if err != nil {
		t.Fatal(err)
	}
	scan := ScanPlanFiles(workRoot)
	return EvaluateRules(root, workRoot, ".hawp/work/BACKLOG.md", backlog, scan)
}

func rulesOf(ops []FixOperation) []string {
	var rules []string
	for _, op := range ops {
		rules = append(rules, op.RuleID)
	}
	return rules
}

func hasRule(ops []FixOperation, rule string) bool {
	for _, op := range ops {
		if op.RuleID == rule {
			return true
		}
	}
	return false
}

const cleanBacklogHeader = `# Backlog

## Active Work

| ID | Type | Title | Status | Plan File | Updated |
| --- | --- | --- | --- | --- | --- |
`

const backlogFooter = `
## Blocked / Parked

| ID | Type | Title | Reason | Detail | Updated |
| --- | --- | --- | --- | --- | --- |

## Recently Closed

| ID | Type | Title | Closed | Detail |
| --- | --- | --- | --- | --- |
`

func TestRulesCleanRepo(t *testing.T) {
	root := buildRepoFixture(t, map[string]string{
		".hawp/work/BACKLOG.md": cleanBacklogHeader +
			"| TASK-001 | task | thing | in-progress | [plan](active/TASK-001.md) | 2026-07-20 |\n" + backlogFooter,
		".hawp/work/active/TASK-001.md": "# plan",
		".hawp/work/parked/.keep":       "",
		".hawp/work/closed/.keep":       "",
	})
	if ops := detect(t, root); len(ops) != 0 {
		t.Fatalf("clean repo produced operations: %v", rulesOf(ops))
	}
}

func TestRowRules(t *testing.T) {
	root := buildRepoFixture(t, map[string]string{
		".hawp/work/BACKLOG.md": cleanBacklogHeader +
			"| TASK-002 | | missing type | inbox | [plan](active/TASK-002.md) | 2026-07-20 |\n" +
			"| `weird_id!` | | odd id, no type, no plan, bad date | inbox | | July 2026 |\n" +
			"| TASK-003 | task | done row | done | [plan](active/TASK-003.md) | 2026-07-20 |\n" +
			"| TASK-004 | task | escaping link | inbox | [plan](../../etc/passwd.md) | 2026-07-20 |\n" +
			"| TASK-005 | task | duplicate plans | inbox | [plan](active/TASK-005.md) | 2026-07-20 |\n" + backlogFooter,
		".hawp/work/active/TASK-002.md": "# plan",
		".hawp/work/active/TASK-003.md": "# plan",
		".hawp/work/active/TASK-005.md": "# plan",
		".hawp/work/parked/TASK-005.md": "# duplicate",
		".hawp/work/closed/.keep":       "",
	})
	ops := detect(t, root)

	for _, want := range []string{"A1", "A2", "A3", "A6", "B1", "B2", "B3"} {
		if !hasRule(ops, want) {
			t.Errorf("missing rule %s in %v", want, rulesOf(ops))
		}
	}
	// A1 fires for TASK-002 (inferable); B1 for weird_id (not inferable).
	for _, op := range ops {
		if op.RuleID == "A1" && op.ItemID != "TASK-002" {
			t.Errorf("A1 fired for %s", op.ItemID)
		}
		if op.RuleID == "A6" && op.ItemID != "TASK-003" {
			t.Errorf("A6 fired for %s", op.ItemID)
		}
	}
}

func TestClosedRecordRules(t *testing.T) {
	incomplete := "# closed\n\n**Backlog ID:** TASK-010\n\nno sections here\n"
	ambiguous := `# closed

**Backlog ID:** TASK-011

## Outcome

Done.

## Verification

- [x] works on my machine
- [x] proven (Evidence: log)

## Close Checklist

- [x] done
`
	root := buildRepoFixture(t, map[string]string{
		".hawp/work/BACKLOG.md": cleanBacklogHeader + backlogFooter +
			"| TASK-010 | task | incomplete | 2026-07-01 | [plan](closed/2026/07/01/TASK-010.md) |\n" +
			"| TASK-011 | task | ambiguous | 2026-07-01 | [plan](closed/2026/07/01/TASK-011.md) |\n",
		".hawp/work/closed/2026/07/01/TASK-010.md": incomplete,
		".hawp/work/closed/2026/07/01/TASK-011.md": ambiguous,
		".hawp/work/active/.keep":                  "",
		".hawp/work/parked/.keep":                  "",
	})
	ops := detect(t, root)

	if !hasRule(ops, "A4") || !hasRule(ops, "A5") {
		t.Errorf("expected A4+A5 for incomplete closed record, got %v", rulesOf(ops))
	}
	found := false
	for _, op := range ops {
		if op.RuleID == "B7" {
			found = true
			if op.ItemID != "TASK-011" {
				t.Errorf("B7 fired for %s", op.ItemID)
			}
			if len(op.Blocked.Candidates) != 1 || !strings.Contains(op.Blocked.Candidates[0], "works on my machine") {
				t.Errorf("B7 candidates = %v", op.Blocked.Candidates)
			}
		}
	}
	if !found {
		t.Error("B7 not fired for ambiguous claims")
	}
}

func TestStructuralRules(t *testing.T) {
	root := buildRepoFixture(t, map[string]string{
		".hawp/work/BACKLOG.md": "# Backlog\n\n## Active Work\n\n| ID | Type | Title | Status |\n| --- | --- | --- | --- |\n",
	})
	ops := detect(t, root)
	if !hasRule(ops, "A8") {
		t.Errorf("A8 not fired for missing sections: %v", rulesOf(ops))
	}
	if !hasRule(ops, "B5") {
		t.Errorf("B5 not fired for missing directories: %v", rulesOf(ops))
	}
}

func TestUnprovenMarkerSkipsEvidenceLines(t *testing.T) {
	content := `## Verification

- [x] quotes validator output "Explicitly unproven: 1" (Evidence: run log)
`
	if hasUnprovenChecklistMarker(content) {
		t.Error("Evidence-backed line quoting 'unproven' must not trigger B4")
	}
	if !hasUnprovenChecklistMarker("## Verification\n\n- [ ] NOT YET VERIFIED: deploy step\n") {
		t.Error("real unproven marker must trigger B4")
	}
}

func TestApplyClosedRecordNormalization(t *testing.T) {
	root := buildRepoFixture(t, map[string]string{
		".hawp/work/closed/2026/07/01/TASK-020.md":                 "# closed record without sections\n",
		".hawp/work/closed/misplaced/2026-07-02-TASK-021-thing.md": "# date-prefixed, wrong folder\n",
	})
	result, err := ApplyClosedRecordNormalization(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.ChangedFiles) != 2 {
		t.Fatalf("changed = %v, want 2 files", result.ChangedFiles)
	}

	normalized, err := os.ReadFile(filepath.Join(root, ".hawp/work/closed/2026/07/01/TASK-020.md"))
	if err != nil {
		t.Fatal(err)
	}
	content := string(normalized)
	for _, section := range []string{"**Backlog ID:** TASK-020", "## Outcome", "## Verification", "## Close Checklist"} {
		if !strings.Contains(content, section) {
			t.Errorf("normalized record missing %q", section)
		}
	}

	// The misplaced date-prefixed file moved into closed/2026/07/02/.
	moved := filepath.Join(root, ".hawp/work/closed/2026/07/02/2026-07-02-TASK-021-thing.md")
	if _, err := os.Stat(moved); err != nil {
		t.Errorf("date-prefixed file not reconciled to %s", moved)
	}

	// Re-running is idempotent.
	again, err := ApplyClosedRecordNormalization(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(again.ChangedFiles) != 0 {
		t.Errorf("second apply changed files: %v", again.ChangedFiles)
	}
}

func TestEvidenceFollowUpQueue(t *testing.T) {
	content := `# plan

## Outcome

x

## Verification

- [x] ambiguous thing

## Close Checklist

- [x] done
`
	next, claims := normalizeClosedRecord(content, "/x/closed/2026/07/01/TASK-030.md")
	if len(claims) != 1 || claims[0] != "ambiguous thing" {
		t.Fatalf("claims = %v", claims)
	}
	if !strings.Contains(next, "### Evidence Follow-Up") ||
		!strings.Contains(next, "Research evidence for: ambiguous thing") {
		t.Errorf("follow-up subsection missing:\n%s", next)
	}
	// Idempotent: same claim is not queued twice.
	final, claimsAgain := normalizeClosedRecord(next, "/x/closed/2026/07/01/TASK-030.md")
	if len(claimsAgain) != 0 {
		t.Errorf("re-run added claims again: %v", claimsAgain)
	}
	if strings.Count(final, "Research evidence for: ambiguous thing") != 1 {
		t.Errorf("duplicate research entries:\n%s", final)
	}
}
