// Package work is the application service for `hawp work validate`.
package work

import (
	"fmt"
	"io"
	"path/filepath"

	domainwork "github.com/sentzunhat/hawp/librarian/go/internal/domain/work"
	filesystemwork "github.com/sentzunhat/hawp/librarian/go/internal/infrastructure/filesystem/work"
)

// Validate parses the backlog and runs all five checks against workDir
// (a .hawp/work directory).
func Validate(workDir string) (*domainwork.Report, error) {
	snapshot, err := filesystemwork.NewAdapter().Read(workDir)
	if err != nil {
		return nil, fmt.Errorf("could not parse BACKLOG.md: %w", err)
	}
	backlog := domainwork.ParseBacklogContent(snapshot.BacklogContent)

	closedFiles := domainwork.CollectClosedPlanFiles(filepath.Join(workDir, "closed"))

	report := &domainwork.Report{
		Backlog:      domainwork.CheckBacklogConsistency(workDir, backlog),
		Completeness: domainwork.CheckClosedTaskCompleteness(snapshot),
		Evidence:     domainwork.CheckEvidenceIntegrity(workDir, closedFiles),
		Clarity:      domainwork.CheckVerificationClarity(closedFiles),
		DeadLinks:    domainwork.CheckDeadLinks(snapshot),
	}
	report.Summarize()
	return report, nil
}

// Render writes the validation report and returns the exit code
// (0 for PASS/WARN, 1 for FAIL).
func Render(out io.Writer, workDir string, r *domainwork.Report) int {
	line := func(format string, args ...any) { fmt.Fprintf(out, format+"\n", args...) }
	rule := func() { line("----------------------------------------------------------------------") }

	line("Validating: %s", workDir)
	line("")

	line("1. BACKLOG CONSISTENCY [%s]", r.Backlog.Status)
	rule()
	line("  Active: %d/%d  Closed: %d/%d  Parked: %d/%d",
		r.Backlog.ActiveWork.Found, r.Backlog.ActiveWork.Total,
		r.Backlog.RecentlyClosed.Found, r.Backlog.RecentlyClosed.Total,
		r.Backlog.ParkedWork.Found, r.Backlog.ParkedWork.Total)
	for _, id := range r.Backlog.ActiveWork.Missing {
		line("  ✗ active row without plan file: %s", id)
	}
	for _, id := range r.Backlog.RecentlyClosed.Missing {
		line("  ✗ closed row without plan file: %s", id)
	}
	for _, id := range r.Backlog.ParkedWork.Missing {
		line("  ✗ parked row without plan file: %s", id)
	}
	for _, f := range r.Backlog.OrphanedFiles {
		line("  ✗ orphaned active file: %s", f)
	}
	for _, f := range r.Backlog.OrphanedParked {
		line("  ✗ orphaned parked file: %s", f)
	}
	line("")

	line("2. CLOSED TASK COMPLETENESS [%s]", r.Completeness.Status)
	rule()
	line("  Plans: %d (skipped %d supporting)  Outcome: %d/%d  Verification: %d/%d  Close Checklist: %d/%d",
		r.Completeness.Total, r.Completeness.Skipped,
		r.Completeness.WithOutcome, r.Completeness.Total,
		r.Completeness.WithVerification, r.Completeness.Total,
		r.Completeness.WithCloseChecklist, r.Completeness.Total)
	for _, f := range r.Completeness.Failing {
		line("  ✗ %s missing %v (%s) [%s]", f.ID, f.Sections, f.Date, f.FilePath)
	}
	for _, f := range r.Completeness.UntypedCurrent {
		line("  ✗ %s: %s (%s)", f.ID, f.Reason, f.Date)
	}
	if n := len(r.Completeness.Warnings); n > 0 {
		line("  ! %d legacy plan(s) missing sections (tolerated)", n)
	}
	if n := len(r.Completeness.UntypedLegacy); n > 0 {
		line("  ! %d legacy untyped file(s) (tolerated)", n)
	}
	line("")

	line("3. EVIDENCE INTEGRITY [%s]", r.Evidence.Status)
	rule()
	line("  Evidence links: %d/%d valid", r.Evidence.Valid, r.Evidence.Total)
	for _, b := range r.Evidence.Broken {
		line("  ! %s -> %s", b.ID, b.Link)
	}
	line("")

	line("4. VERIFICATION CLARITY [%s]", r.Clarity.Status)
	rule()
	line("  Evidence-backed: %d/%d  Explicitly unproven: %d  Ambiguous: %d",
		r.Clarity.Proven, r.Clarity.Total, len(r.Clarity.Unproven), len(r.Clarity.Ambiguous))
	line("")

	line("5. DEAD LINKS [%s]", r.DeadLinks.Status)
	rule()
	line("  Scanned %d file(s), %d broken", r.DeadLinks.Scanned, len(r.DeadLinks.Broken))
	for _, b := range r.DeadLinks.Broken {
		line("  ✗ %s -> %s", b.ID, b.Link)
	}
	line("")

	line("SUMMARY")
	rule()
	line("  ✓ Checks passed: %d", r.Passed)
	line("  ✗ Issues found:  %d", r.Failed)
	line("  ! Warnings:      %d", r.Warnings)
	line("")
	line("Result: VALIDATION %s", r.Overall)

	if r.Overall == domainwork.StatusFail {
		return 1
	}
	return 0
}
