package work

import (
	"path/filepath"
	"testing"

	"github.com/sentzunhat/hawp/librarian/go/internal/domain/work/source"
)

const closedPlanComplete = `# done thing

## Outcome (filled at close)

Done.

## Verification (filled at close)

- [x] It works (Evidence: test run output)

## Close Checklist

- [x] Outcome section filled
`

// TestBacklogConsistencyPassAndFail tests the snapshot-based consistency check
// without touching the filesystem (non-repository path).
func TestBacklogConsistencyPassAndFail(t *testing.T) {
	snapshot := source.Snapshot{
		Files: []source.File{
			{RelPath: "active/39bc92b6-5e3f-49ae-970d-6126ec2cfd82.md"},
			{RelPath: "closed/2026/07/03/TASK-086.md"},
			{RelPath: "parked/bee15107.md"},
		},
	}
	backlog := &Backlog{
		Active: []BacklogRow{{ID: "39bc92b6"}}, // short-UUID row → full-UUID file
		Closed: []BacklogRow{{ID: "TASK-086"}},
		Parked: []BacklogRow{{ID: "bee15107", Detail: "[plan](parked/bee15107.md)"}},
	}
	result := CheckBacklogConsistency(snapshot, backlog)
	if result.Status != StatusPass {
		t.Fatalf("status = %s, want PASS: %+v", result.Status, result)
	}

	// A row without a matching file is reported as missing.
	backlog.Active = append(backlog.Active, BacklogRow{ID: "deadbeef"})
	result = CheckBacklogConsistency(snapshot, backlog)
	if result.Status != StatusFail || len(result.ActiveWork.Missing) != 1 {
		t.Fatalf("missing row not detected: %+v", result.ActiveWork)
	}

	// An active plan file whose ID is not in the backlog is reported as orphaned.
	orphanSnapshot := source.Snapshot{
		Files: []source.File{
			{RelPath: "active/TASK-999.md"},
		},
	}
	result = CheckBacklogConsistency(orphanSnapshot, &Backlog{})
	if result.Status != StatusFail || len(result.OrphanedFiles) != 1 {
		t.Fatalf("orphan not detected: %+v", result.OrphanedFiles)
	}
}

func TestClosedTaskCompleteness(t *testing.T) {
	snapshot := source.Snapshot{Files: []source.File{
		{Path: "/work/closed/2026/07/03/TASK-086.md", RelPath: "closed/2026/07/03/TASK-086.md", Content: closedPlanComplete},
		{Path: "/work/closed/2026/07/03/TASK-087.md", RelPath: "closed/2026/07/03/TASK-087.md", Content: "# plan without sections"},
		{Path: "/work/closed/2026/04/01/TASK-001.md", RelPath: "closed/2026/04/01/TASK-001.md", Content: "# legacy plan without sections"},
		{Path: "/work/closed/2026/07/03/TASK-086-summary.md", RelPath: "closed/2026/07/03/TASK-086-summary.md", Content: "# supporting, skipped"},
	}}
	result := CheckClosedTaskCompleteness(snapshot)
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

	if got := CheckClosedTaskCompleteness(source.Snapshot{Files: []source.File{
		{Path: "/work/closed/2026/07/03/TASK-086.md", RelPath: "closed/2026/07/03/TASK-086.md", Content: closedPlanComplete},
	}}); got.Status != StatusPass {
		t.Errorf("all-complete status = %s, want PASS", got.Status)
	}
}

// TestEvidenceIntegrity tests the snapshot-based evidence check without
// touching the filesystem (non-repository path).
func TestEvidenceIntegrity(t *testing.T) {
	snapshot := source.Snapshot{
		ClosedFiles: []source.File{
			{
				Path:    "/work/closed/2026/07/03/TASK-086.md",
				RelPath: "closed/2026/07/03/TASK-086.md",
				Content: "## Verification\n\n- [x] ok Evidence: ../evidence/2026/07/03/TASK-086-run.md\n- [x] bad Evidence: ../evidence/2026/07/03/missing.md\n",
			},
		},
		Files: []source.File{
			{RelPath: "evidence/2026/07/03/TASK-086-run.md"},
			// missing.md is intentionally absent
		},
	}
	result := CheckEvidenceIntegrity(snapshot)
	if result.Total != 2 || result.Valid != 1 || len(result.Broken) != 1 {
		t.Fatalf("evidence = %+v, want 2 total / 1 valid / 1 broken", result)
	}
	if result.Status != StatusWarn {
		t.Errorf("status = %s, want WARN", result.Status)
	}

	// A link containing ".." is rejected and not counted.
	escSnapshot := source.Snapshot{
		ClosedFiles: []source.File{
			{
				Path:    "/work/closed/2026/07/03/TASK-001.md",
				RelPath: "closed/2026/07/03/TASK-001.md",
				Content: "Evidence: ../evidence/../../../etc/passwd.md\n",
			},
		},
	}
	escResult := CheckEvidenceIntegrity(escSnapshot)
	if escResult.Total != 0 {
		t.Errorf("escaping link counted: %+v", escResult)
	}
}

// TestVerificationClarity tests the snapshot-based clarity check without
// touching the filesystem (non-repository path).
func TestVerificationClarity(t *testing.T) {
	closedFiles := []source.File{
		{
			Path:    "/work/closed/2026/07/03/TASK-086.md",
			RelPath: "closed/2026/07/03/TASK-086.md",
			Content: `## Verification (filled at close)

- [x] proven claim (Evidence: test output)
- [x] ambiguous claim with no marker
- [ ] explicitly unproven claim
- [x] Research evidence for: skipped meta-claim

## Next Section
`,
		},
	}
	result := CheckVerificationClarity(closedFiles)
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

func TestDeadLinks(t *testing.T) {
	root := t.TempDir()
	backlog := filepath.Join(root, "BACKLOG.md")
	plan := filepath.Join(root, "active", "TASK-001.md")
	result := CheckDeadLinks(source.Snapshot{
		Files: []source.File{
			{Path: backlog, RelPath: "BACKLOG.md", Links: []source.Link{{Href: "active/TASK-001.md"}, {Href: "active/missing.md"}}},
			{Path: plan, RelPath: "active/TASK-001.md"},
		},
		ExistingPaths: map[string]struct{}{backlog: {}, plan: {}},
	})
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
