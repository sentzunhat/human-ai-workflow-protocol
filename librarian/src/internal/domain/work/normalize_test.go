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

func TestRulesAcceptLegacySubsectionsAndPlainPlanPaths(t *testing.T) {
	root := buildRepoFixture(t, map[string]string{
		".hawp/work/BACKLOG.md": `# Backlog

## Active Work

### Release Readiness

| # | Status | Title | Plan File | Next action |
| --- | --- | --- | --- | --- |
| 049 | in-progress | legacy item | active/049.md | next |

## Blocked / Parked

| # | Status | Title | Detail | Next action |
| --- | --- | --- | --- | --- |
| 040 | parked | parked legacy item | parked/040.md | later |

## Recently Closed

| # | Title | Closed | Plan File |
| --- | --- | --- | --- |
| 042 | closed legacy item | 2026-07-27 | closed/2026/07/27/042.md |
`,
		".hawp/work/active/049.md":            "# active plan\n",
		".hawp/work/parked/040.md":            "# parked plan\n",
		".hawp/work/closed/2026/07/27/042.md": closedPlanComplete,
	})
	ops := detect(t, root)
	if len(ops) != 3 {
		t.Fatalf("legacy repo produced %d operations, want 3 type warnings: %v", len(ops), rulesOf(ops))
	}
	for _, op := range ops {
		if op.RuleID != "B1" {
			t.Fatalf("legacy repo produced non-type operation: %+v", op)
		}
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

func TestApplyWorkItemFolderMigrationMovesFlatPlanAndSidecar(t *testing.T) {
	root := buildRepoFixture(t, map[string]string{
		".hawp/work/BACKLOG.md": cleanBacklogHeader +
			"| `legacy-item` | task | legacy | inbox | [plan](active/legacy-item.md) | 2026-08-30 |\n" +
			"| `v0.1.0-cloud-backends` | task | parked | parked | [plan](parked/v0.1.0-cloud-backends.md) | 2026-08-30 |\n" + backlogFooter,
		".hawp/work/active/legacy-item.md": `# Legacy Item

**Plan file:** work/active/legacy-item.md

See [notes](../notes/context.md).
`,
		".hawp/work/active/legacy-item-files.md": `# Files

**Work Item:** .hawp/work/active/legacy-item.md

See [notes](../notes/context.md).
`,
		".hawp/work/parked/v0.1.0-cloud-backends.md": `# Parked

**Plan file:** work/parked/v0.1.0-cloud-backends.md
`,
		".hawp/work/notes/context.md": "notes",
		".hawp/work/closed/.keep":     "",
	})

	result, err := ApplyWorkItemFolderMigration(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.ChangedFiles) == 0 {
		t.Fatal("expected migration changes")
	}

	planPath := filepath.Join(root, ".hawp/work/active/legacy-item/plan.md")
	plan, err := os.ReadFile(planPath)
	if err != nil {
		t.Fatalf("migrated active plan missing: %v", err)
	}
	if strings.Contains(string(plan), "work/active/legacy-item.md") {
		t.Errorf("plan file path not rewritten:\n%s", string(plan))
	}
	if !strings.Contains(string(plan), "work/active/legacy-item/plan.md") {
		t.Errorf("plan file path missing new location:\n%s", string(plan))
	}
	if !strings.Contains(string(plan), "[notes](../../notes/context.md)") {
		t.Errorf("relative link not rewritten for moved plan:\n%s", string(plan))
	}

	filesPath := filepath.Join(root, ".hawp/work/active/legacy-item/files.md")
	filesContent, err := os.ReadFile(filesPath)
	if err != nil {
		t.Fatalf("migrated files sidecar missing: %v", err)
	}
	if !strings.Contains(string(filesContent), ".hawp/work/active/legacy-item/plan.md") {
		t.Errorf("files.md work-item path not rewritten:\n%s", string(filesContent))
	}
	if !strings.Contains(string(filesContent), "[notes](../../notes/context.md)") {
		t.Errorf("relative link not rewritten for moved files.md:\n%s", string(filesContent))
	}

	if _, err := os.Stat(filepath.Join(root, ".hawp/work/active/legacy-item.md")); !os.IsNotExist(err) {
		t.Errorf("old flat plan should be removed, stat err = %v", err)
	}
	if _, err := os.Stat(filepath.Join(root, ".hawp/work/active/legacy-item-files.md")); !os.IsNotExist(err) {
		t.Errorf("old flat sidecar should be removed, stat err = %v", err)
	}

	parkedPlan, err := os.ReadFile(filepath.Join(root, ".hawp/work/parked/v0.1.0-cloud-backends/plan.md"))
	if err != nil {
		t.Fatalf("migrated parked plan missing: %v", err)
	}
	if !strings.Contains(string(parkedPlan), "work/parked/v0.1.0-cloud-backends/plan.md") {
		t.Errorf("parked plan file path not rewritten:\n%s", string(parkedPlan))
	}

	backlog, err := os.ReadFile(filepath.Join(root, ".hawp/work/BACKLOG.md"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(backlog), "(active/legacy-item/plan.md)") {
		t.Errorf("backlog active link not rewritten:\n%s", string(backlog))
	}
	if !strings.Contains(string(backlog), "(parked/v0.1.0-cloud-backends/plan.md)") {
		t.Errorf("backlog parked link not rewritten:\n%s", string(backlog))
	}

	again, err := ApplyWorkItemFolderMigration(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(again.ChangedFiles) != 0 {
		t.Errorf("second migration changed files: %v", again.ChangedFiles)
	}
}

func TestApplyWorkItemFolderMigrationRenamesFolderToUUID(t *testing.T) {
	root := buildRepoFixture(t, map[string]string{
		".hawp/work/BACKLOG.md": cleanBacklogHeader +
			"| `12345678` | task | uuid item | inbox | [plan](active/legacy-slug/plan.md) | 2026-08-30 |\n" + backlogFooter,
		".hawp/work/active/legacy-slug/plan.md": `# UUID Item

**UUID:** ` + "`12345678-abcd-4abc-8def-1234567890ab`" + `
**Plan file:** work/active/legacy-slug/plan.md
`,
		".hawp/work/active/legacy-slug/files.md": `# Files

**Work Item:** .hawp/work/active/legacy-slug/plan.md
`,
		".hawp/work/parked/.keep": "",
		".hawp/work/closed/.keep": "",
	})

	result, err := ApplyWorkItemFolderMigration(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.ChangedFiles) == 0 {
		t.Fatal("expected rename changes")
	}

	if _, err := os.Stat(filepath.Join(root, ".hawp/work/active/12345678/plan.md")); err != nil {
		t.Fatalf("uuid target folder missing: %v", err)
	}
	if _, err := os.Stat(filepath.Join(root, ".hawp/work/active/legacy-slug")); !os.IsNotExist(err) {
		t.Errorf("legacy folder should be removed, stat err = %v", err)
	}

	plan, err := os.ReadFile(filepath.Join(root, ".hawp/work/active/12345678/plan.md"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(plan), "work/active/12345678/plan.md") {
		t.Errorf("uuid plan path not rewritten:\n%s", string(plan))
	}

	filesContent, err := os.ReadFile(filepath.Join(root, ".hawp/work/active/12345678/files.md"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(filesContent), ".hawp/work/active/12345678/plan.md") {
		t.Errorf("uuid files.md work-item path not rewritten:\n%s", string(filesContent))
	}
}

func TestPreviewWorkItemFolderMigrationMatchesApplyWithoutMutatingSource(t *testing.T) {
	root := buildRepoFixture(t, map[string]string{
		".hawp/work/BACKLOG.md": cleanBacklogHeader +
			"| `legacy-item` | task | legacy | inbox | [plan](active/legacy-item.md) | 2026-08-30 |\n" + backlogFooter,
		".hawp/work/active/legacy-item.md": `# Legacy Item

**Plan file:** work/active/legacy-item.md
`,
		".hawp/work/parked/.keep": "",
		".hawp/work/closed/.keep": "",
	})

	preview, err := PreviewWorkItemFolderMigration(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(preview.ChangedFiles) == 0 {
		t.Fatal("expected preview to report migration changes")
	}
	if _, err := os.Stat(filepath.Join(root, ".hawp/work/active/legacy-item.md")); err != nil {
		t.Fatalf("preview should not mutate source repo, stat err = %v", err)
	}
	if _, err := os.Stat(filepath.Join(root, ".hawp/work/active/legacy-item/plan.md")); !os.IsNotExist(err) {
		t.Fatalf("preview should not create migrated folder in source repo, stat err = %v", err)
	}

	applied, err := ApplyWorkItemFolderMigration(root)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Join(preview.ChangedFiles, "\n") != strings.Join(applied.ChangedFiles, "\n") {
		t.Fatalf("preview changes %v do not match apply changes %v", preview.ChangedFiles, applied.ChangedFiles)
	}
}
