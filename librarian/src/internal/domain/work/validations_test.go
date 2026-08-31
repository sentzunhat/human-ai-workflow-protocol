package work

import (
	"os"
	"path/filepath"
	"testing"
)

// buildWorkDir creates a minimal .hawp/work fixture tree.
func buildWorkDir(t *testing.T, files map[string]string) string {
	t.Helper()
	workDir := t.TempDir()
	for rel, content := range files {
		full := filepath.Join(workDir, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return workDir
}

const closedPlanComplete = `# done thing

## Outcome (filled at close)

Done.

## Verification (filled at close)

- [x] It works (Evidence: test run output)

## Close Checklist

- [x] Outcome section filled
`

func TestBacklogConsistencyPassAndFail(t *testing.T) {
	workDir := buildWorkDir(t, map[string]string{
		"active/39bc92b6-5e3f-49ae-970d-6126ec2cfd82.md": "# plan",
		"closed/2026/07/03/TASK-086.md":                  closedPlanComplete,
		"parked/bee15107.md":                             "# parked",
	})
	backlog := &Backlog{
		Active: []BacklogRow{{ID: "39bc92b6"}}, // short-UUID row → full-UUID file
		Closed: []BacklogRow{{ID: "TASK-086"}},
		Parked: []BacklogRow{{ID: "bee15107", Detail: "[plan](parked/bee15107.md)"}},
	}
	result := CheckBacklogConsistency(workDir, backlog)
	if result.Status != StatusPass {
		t.Fatalf("status = %s, want PASS: %+v", result.Status, result)
	}

	// A row without a file and an orphaned file both fail.
	backlog.Active = append(backlog.Active, BacklogRow{ID: "deadbeef"})
	result = CheckBacklogConsistency(workDir, backlog)
	if result.Status != StatusFail || len(result.ActiveWork.Missing) != 1 {
		t.Fatalf("missing row not detected: %+v", result.ActiveWork)
	}

	orphanDir := buildWorkDir(t, map[string]string{
		"active/TASK-999.md": "# orphan",
	})
	result = CheckBacklogConsistency(orphanDir, &Backlog{})
	if result.Status != StatusFail || len(result.OrphanedFiles) != 1 {
		t.Fatalf("orphan not detected: %+v", result.OrphanedFiles)
	}
}

func TestBacklogConsistencyMatchesSlugClosedFolder(t *testing.T) {
	workDir := buildWorkDir(t, map[string]string{
		"closed/2026/08/25/manager-branch-kit-pattern/plan.md": closedPlanComplete,
	})
	backlog := &Backlog{
		Closed: []BacklogRow{{ID: "manager-branch-kit-pattern"}},
	}
	result := CheckBacklogConsistency(workDir, backlog)
	if result.Status != StatusPass {
		t.Fatalf("status = %s, want PASS: %+v", result.Status, result)
	}
	if result.RecentlyClosed.Found != 1 {
		t.Fatalf("closed found = %d, want 1", result.RecentlyClosed.Found)
	}
}

func TestClosedTaskCompleteness(t *testing.T) {
	workDir := buildWorkDir(t, map[string]string{
		"closed/2026/07/03/TASK-086.md":         closedPlanComplete,
		"closed/2026/07/03/TASK-087.md":         "# plan without sections",
		"closed/2026/04/01/TASK-001.md":         "# legacy plan without sections",
		"closed/2026/07/03/TASK-086-summary.md": "# supporting, skipped",
	})
	result := CheckClosedTaskCompleteness(workDir)
	if result.Total != 3 {
		t.Fatalf("total plans = %d, want 3", result.Total)
	}
	if result.Skipped != 1 {
		t.Errorf("skipped = %d, want 1 (summary suffix)", result.Skipped)
	}
	if len(result.Failing) != 1 || result.Failing[0].ID != "TASK-087" {
		t.Errorf("failing = %+v, want TASK-087", result.Failing)
	}
	if len(result.Warnings) != 1 || result.Warnings[0].ID != "TASK-001" {
		t.Errorf("warnings = %+v, want legacy TASK-001", result.Warnings)
	}
	if result.Status != StatusFail {
		t.Errorf("status = %s, want FAIL", result.Status)
	}

	passDir := buildWorkDir(t, map[string]string{
		"closed/2026/07/03/TASK-086.md": closedPlanComplete,
	})
	if got := CheckClosedTaskCompleteness(passDir); got.Status != StatusPass {
		t.Errorf("all-complete status = %s, want PASS", got.Status)
	}
}

func TestClosedTaskCompletenessAcceptsSlugPlanFolders(t *testing.T) {
	workDir := buildWorkDir(t, map[string]string{
		"closed/2026/08/25/releases-prerelease-fallback/plan.md": `# releases-prerelease-fallback

**Type:** investigation
**Status:** done

## Outcome

Done.

## Verification

- [x] ok (Evidence: output captured)

## Close Checklist

- [x] done
`,
	})
	result := CheckClosedTaskCompleteness(workDir)
	if result.Total != 1 {
		t.Fatalf("total plans = %d, want 1", result.Total)
	}
	if len(result.UntypedCurrent) != 0 {
		t.Fatalf("untyped current = %+v, want none", result.UntypedCurrent)
	}
	if result.Status != StatusPass {
		t.Fatalf("status = %s, want PASS", result.Status)
	}
}

func TestEvidenceIntegrity(t *testing.T) {
	workDir := buildWorkDir(t, map[string]string{
		"closed/2026/07/03/TASK-086.md":       "## Verification\n\n- [x] ok Evidence: ../evidence/2026/07/03/TASK-086-run.md\n- [x] bad Evidence: ../evidence/2026/07/03/missing.md\n",
		"evidence/2026/07/03/TASK-086-run.md": "# evidence",
	})
	files := CollectClosedPlanFiles(filepath.Join(workDir, "closed"))
	if len(files) != 1 {
		t.Fatalf("closed plan files = %d, want 1", len(files))
	}
	result := CheckEvidenceIntegrity(workDir, files)
	if result.Total != 2 || result.Valid != 1 || len(result.Broken) != 1 {
		t.Fatalf("evidence = %+v, want 2 total / 1 valid / 1 broken", result)
	}
	if result.Status != StatusWarn {
		t.Errorf("status = %s, want WARN", result.Status)
	}

	// A link escaping the evidence folder is skipped, not counted.
	escDir := buildWorkDir(t, map[string]string{
		"closed/2026/07/03/TASK-001.md": "Evidence: ../evidence/../../../etc/passwd.md\n",
	})
	escResult := CheckEvidenceIntegrity(escDir, CollectClosedPlanFiles(filepath.Join(escDir, "closed")))
	if escResult.Total != 0 {
		t.Errorf("escaping link counted: %+v", escResult)
	}
}

func TestCollectClosedPlanFilesIncludesFolderPlans(t *testing.T) {
	workDir := buildWorkDir(t, map[string]string{
		"closed/2026/08/25/flat.md":             closedPlanComplete,
		"closed/2026/08/25/folder-item/plan.md": closedPlanComplete,
	})
	files := CollectClosedPlanFiles(filepath.Join(workDir, "closed"))
	if len(files) != 2 {
		t.Fatalf("closed plan files = %d, want 2", len(files))
	}
}

func TestVerificationClarity(t *testing.T) {
	workDir := buildWorkDir(t, map[string]string{
		"closed/2026/07/03/TASK-086.md": `## Verification (filled at close)

- [x] proven claim (Evidence: test output)
- [x] ambiguous claim with no marker
- [ ] explicitly unproven claim
- [x] Research evidence for: skipped meta-claim

## Next Section
`,
	})
	files := CollectClosedPlanFiles(filepath.Join(workDir, "closed"))
	result := CheckVerificationClarity(files)
	if result.Total != 3 {
		t.Fatalf("total claims = %d, want 3 (meta-claim excluded)", result.Total)
	}
	if result.Proven != 1 || len(result.Ambiguous) != 1 || len(result.Unproven) != 1 {
		t.Fatalf("clarity = proven %d / ambiguous %d / unproven %d, want 1/1/1",
			result.Proven, len(result.Ambiguous), len(result.Unproven))
	}
	if result.Status != StatusWarn {
		t.Errorf("status = %s, want WARN", result.Status)
	}
}

func TestVerificationClarityCountsMultilineEvidenceAsProven(t *testing.T) {
	workDir := buildWorkDir(t, map[string]string{
		"closed/2026/08/25/folder-item/plan.md": `## Verification

- [x] benchmark evidence recorded
      Evidence: ../evidence/2026/08/25/folder-item-proof.md
- [x] explicit unproven follow-up
      explicitly unproven pending CI

## Close Checklist
`,
		"evidence/2026/08/25/folder-item-proof.md": "# proof",
	})
	files := CollectClosedPlanFiles(filepath.Join(workDir, "closed"))
	result := CheckVerificationClarity(files)
	if result.Total != 2 {
		t.Fatalf("total claims = %d, want 2", result.Total)
	}
	if result.Proven != 1 || len(result.Unproven) != 1 || len(result.Ambiguous) != 0 {
		t.Fatalf("clarity = proven %d / ambiguous %d / unproven %d, want 1/0/1",
			result.Proven, len(result.Ambiguous), len(result.Unproven))
	}
	if result.Status != StatusPass {
		t.Errorf("status = %s, want PASS", result.Status)
	}
}

func TestVerificationClarityUnprovenOnlyPasses(t *testing.T) {
	workDir := buildWorkDir(t, map[string]string{
		"closed/2026/08/25/folder-item/plan.md": `## Verification

- [x] proven claim (Evidence: test output)
- [ ] explicitly unproven claim

## Close Checklist
`,
	})
	files := CollectClosedPlanFiles(filepath.Join(workDir, "closed"))
	result := CheckVerificationClarity(files)
	if result.Total != 2 {
		t.Fatalf("total claims = %d, want 2", result.Total)
	}
	if result.Proven != 1 || len(result.Unproven) != 1 || len(result.Ambiguous) != 0 {
		t.Fatalf("clarity = proven %d / ambiguous %d / unproven %d, want 1/0/1",
			result.Proven, len(result.Ambiguous), len(result.Unproven))
	}
	if result.Status != StatusPass {
		t.Errorf("status = %s, want PASS", result.Status)
	}
	if result.Unproven[0].ID != "folder-item" {
		t.Fatalf("unproven id = %q, want folder-item", result.Unproven[0].ID)
	}
}

func TestDeadLinks(t *testing.T) {
	workDir := buildWorkDir(t, map[string]string{
		"BACKLOG.md":         "see [plan](active/TASK-001.md) and [gone](active/missing.md)\n",
		"active/TASK-001.md": "```\n[example in fence](nowhere.md)\n```\n",
	})
	result := CheckDeadLinks(workDir)
	if result.Scanned != 2 {
		t.Fatalf("scanned = %d, want 2", result.Scanned)
	}
	if len(result.Broken) != 1 || result.Broken[0].Link != "active/missing.md" {
		t.Fatalf("broken = %+v, want only active/missing.md (fenced link ignored)", result.Broken)
	}
	if result.Status != StatusFail {
		t.Errorf("status = %s, want FAIL", result.Status)
	}
}
