package work

import (
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/sentzunhat/hawp/librarian/src/internal/infrastructure/filesystem"
)

// defaultWorkSource for domain/work tests.
var defaultWorkSource = &WorkSource{
	Exists:                 fileExists,
	ToRepoRelative:         func(_, p string) string { return p },
	CollectFiles:           collectFiles,
	ReadDir:                os.ReadDir,
	ReadFile:               os.ReadFile,
	WriteFile:              os.WriteFile,
	MkdirAll:               os.MkdirAll,
	MkdirTemp:              os.MkdirTemp,
	Rename:                 os.Rename,
	Remove:                 os.Remove,
	RemoveAll:              os.RemoveAll,
	Stat:                   os.Stat,
	Lstat:                  os.Lstat,
	EvalSymlinks:           filepath.EvalSymlinks,
	RejectSymlinkAncestors: filesystem.RejectSymlinkAncestors,
}

// collectFiles returns all regular files under dir, recursively.
func collectFiles(dir string, _ bool) []string {
	var result []string
	_ = filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			result = append(result, path)
		}
		return nil
	})
	return result
}

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
	source := &WorkSource{
		Exists:         fileExists,
		ToRepoRelative: func(_, p string) string { return p },
		CollectFiles:   collectFiles,
	}
	return EvaluateRules(root, workRoot, ".hawp/work/BACKLOG.md", backlog, scan, source)
}

var defaultWorkSourceForTests = &WorkSource{
	Exists:         fileExists,
	ToRepoRelative: func(_, p string) string { return p },
	CollectFiles:   collectFiles,
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
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
	if len(ops) != 0 {
		t.Fatalf("legacy numeric repo produced operations: %v", rulesOf(ops))
	}
}

func TestRulesAllowLinkedLegacySlugsOutsideActiveRows(t *testing.T) {
	root := buildRepoFixture(t, map[string]string{
		".hawp/work/BACKLOG.md": cleanBacklogHeader + backlogFooter +
			"| release-benchmark-backfill | planning | benchmark backfill | parked | [plan](parked/release-benchmark-backfill/plan.md) | 2026-08-31 |\n" +
			"| multi-repo-context-d9b2f3a1 | improvement | multi repo context | 2026-09-10 | [plan](closed/2026/09/10/multi-repo-context-d9b2f3a1/plan.md) |\n",
		".hawp/work/active/.keep":                                          "",
		".hawp/work/parked/release-benchmark-backfill/plan.md":             "# parked\n",
		".hawp/work/closed/2026/09/10/multi-repo-context-d9b2f3a1/plan.md": closedPlanComplete,
	})

	ops := detect(t, root)
	if hasRule(ops, "A3") {
		t.Fatalf("linked parked/recently closed legacy slugs should not trigger A3: %v", rulesOf(ops))
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
	result, err := defaultWorkSource.ApplyClosedRecordNormalization(root)
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
	again, err := defaultWorkSource.ApplyClosedRecordNormalization(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(again.ChangedFiles) != 0 {
		t.Errorf("second apply changed files: %v", again.ChangedFiles)
	}
}

func TestApplyClosedRecordNormalizationRejectsSymlinkedRecord(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("symlink permissions are not reliable on Windows")
	}
	root := buildRepoFixture(t, map[string]string{
		".hawp/work/closed/2026/07/01/TASK-020.md": "# closed record without sections\n",
	})
	path := filepath.Join(root, ".hawp/work/closed/2026/07/01/TASK-020.md")
	outside := filepath.Join(root, "outside.md")
	if err := os.WriteFile(outside, []byte("# outside\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Remove(path); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(outside, path); err != nil {
		t.Skipf("symlinks unavailable: %v", err)
	}
	if _, err := defaultWorkSource.ApplyClosedRecordNormalization(root); err == nil {
		t.Fatal("expected symlinked closed record to be rejected")
	}
	content, err := os.ReadFile(outside)
	if err != nil {
		t.Fatal(err)
	}
	if string(content) != "# outside\n" {
		t.Fatalf("outside record was changed: %q", content)
	}
}

func TestApplyWorkItemFolderMigrationRejectsSymlinkedScope(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("symlink permissions are not reliable on Windows")
	}
	root := buildRepoFixture(t, map[string]string{
		".hawp/work/BACKLOG.md": cleanBacklogHeader + backlogFooter,
	})
	outside := filepath.Join(root, "outside-active")
	if err := os.Mkdir(outside, 0o755); err != nil {
		t.Fatal(err)
	}
	active := filepath.Join(root, ".hawp/work/active")
	if err := os.Symlink(outside, active); err != nil {
		t.Skipf("symlinks unavailable: %v", err)
	}
	if _, err := defaultWorkSource.ApplyWorkItemFolderMigration(root); err == nil {
		t.Fatal("expected symlinked active root to be rejected")
	}
}

func TestApplyCompletedActiveRowCleanupRemovesDoneRowsPointingToClosedPlans(t *testing.T) {
	root := buildRepoFixture(t, map[string]string{
		".hawp/work/BACKLOG.md": cleanBacklogHeader +
			"| 048 | task | done stale | done | [plan](closed/2026/09/01/048.md) | 2026-09-01 |\n" +
			"| 049 | task | still active | in-progress | [plan](active/049.md) | 2026-09-01 |\n" +
			"| 050 | task | done but missing closed plan | done | [plan](closed/2026/09/01/050.md) | 2026-09-01 |\n" + backlogFooter +
			"| 048 | task | done stale | 2026-09-01 | [plan](closed/2026/09/01/048.md) |\n",
		".hawp/work/active/049.md":            "# active\n",
		".hawp/work/closed/2026/09/01/048.md": closedPlanComplete,
	})
	source := &WorkSource{Exists: fileExists, ToRepoRelative: func(_, p string) string { return p }, ReadDir: os.ReadDir, ReadFile: os.ReadFile, Stat: os.Stat, Lstat: os.Lstat, EvalSymlinks: filepath.EvalSymlinks, WriteFile: os.WriteFile}
	result, err := source.ApplyCompletedActiveRowCleanup(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.ChangedFiles) != 1 {
		t.Fatalf("changed files = %v, want BACKLOG.md", result.ChangedFiles)
	}
	backlog, err := os.ReadFile(filepath.Join(root, ".hawp/work/BACKLOG.md"))
	if err != nil {
		t.Fatal(err)
	}
	content := string(backlog)
	if strings.Contains(content, "| 048 | task | done stale | done |") {
		t.Fatalf("stale active row still present:\n%s", content)
	}
	if !strings.Contains(content, "| 049 | task | still active | in-progress |") {
		t.Fatalf("active row removed unexpectedly:\n%s", content)
	}
	if !strings.Contains(content, "| 050 | task | done but missing closed plan | done |") {
		t.Fatalf("missing-plan done row should be preserved:\n%s", content)
	}
	if strings.Count(content, "| 048 | task | done stale | 2026-09-01 |") != 1 {
		t.Fatalf("recently closed row not preserved exactly once:\n%s", content)
	}
}

func TestApplyCompletedActiveRowCleanupRejectsTraversalAndSymlinkedBacklog(t *testing.T) {
	root := buildRepoFixture(t, map[string]string{
		".hawp/work/BACKLOG.md": cleanBacklogHeader +
			"| 048 | task | traversal | done | [plan](closed/../active/048/plan.md) | 2026-09-01 |\n" + backlogFooter,
		".hawp/work/active/048/plan.md": "# active\n",
	})
	source := &WorkSource{Exists: fileExists, ToRepoRelative: func(_, p string) string { return p }, ReadDir: os.ReadDir, ReadFile: os.ReadFile, Stat: os.Stat, Lstat: os.Lstat, EvalSymlinks: filepath.EvalSymlinks, WriteFile: os.WriteFile}
	result, err := source.ApplyCompletedActiveRowCleanup(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.ChangedFiles) != 0 {
		t.Fatalf("traversal row was removed: %v", result.ChangedFiles)
	}

	backlogPath := filepath.Join(root, ".hawp/work/BACKLOG.md")
	outside := filepath.Join(root, "outside.md")
	outsideContent := cleanBacklogHeader +
		"| 049 | task | outside | done | [plan](closed/2026/09/01/049.md) | 2026-09-01 |\n" + backlogFooter
	if err := os.WriteFile(outside, []byte(outsideContent), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(root, "outside-closed", "2026", "09", "01"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(root, ".hawp/work/closed/2026/09/01"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, ".hawp/work/closed/2026/09/01/049.md"), []byte(closedPlanComplete), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Remove(backlogPath); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(outside, backlogPath); err != nil {
		t.Skipf("symlinks unavailable: %v", err)
	}
	if _, err := source.ApplyCompletedActiveRowCleanup(root); err == nil {
		t.Fatal("expected symlinked BACKLOG.md to be rejected")
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

	result, err := defaultWorkSource.ApplyWorkItemFolderMigration(root)
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

	again, err := defaultWorkSource.ApplyWorkItemFolderMigration(root)
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

	result, err := defaultWorkSource.ApplyWorkItemFolderMigration(root)
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

	preview, err := defaultWorkSource.PreviewWorkItemFolderMigration(root)
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

	applied, err := defaultWorkSource.ApplyWorkItemFolderMigration(root)
	if err != nil {
		t.Fatal(err)
	}
	previewNorm := normalizeMigratePaths(preview.ChangedFiles)
	applyNorm := normalizeMigratePaths(applied.ChangedFiles)
	if strings.Join(previewNorm, "\n") != strings.Join(applyNorm, "\n") {
		t.Fatalf("preview changes %v do not match apply changes %v", previewNorm, applyNorm)
	}
}

// normalizeMigratePaths strips the work-root prefix so absolute paths
// from different temp roots become comparable relative identifiers.
func normalizeMigratePaths(paths []string) []string {
	out := make([]string, len(paths))
	for i, p := range paths {
		idx := strings.Index(p, string(filepath.Separator)+".hawp"+string(filepath.Separator)+"work")
		if idx >= 0 {
			out[i] = p[idx:]
		} else {
			out[i] = p
		}
	}
	return out
}
