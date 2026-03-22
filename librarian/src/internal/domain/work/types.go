// Package work implements the HAWP work-record validations: backlog
// consistency, closed-task completeness, evidence integrity, verification
// clarity, and dead links. Ported from librarian/scripts/hawp/work-validate.
package work

// CheckStatus is PASS, FAIL, or WARN.
type CheckStatus string

const (
	StatusPass CheckStatus = "PASS"
	StatusFail CheckStatus = "FAIL"
	StatusWarn CheckStatus = "WARN"
)

// BacklogRow is one parsed table row from BACKLOG.md.
type BacklogRow struct {
	ID     string
	Type   string
	Title  string
	Status string
	Detail string
}

// Backlog holds the rows of the three tracked sections.
type Backlog struct {
	Active []BacklogRow
	Closed []BacklogRow
	Parked []BacklogRow
}

// SectionCount tracks row-to-file resolution for one backlog section.
type SectionCount struct {
	Total   int
	Found   int
	Missing []string
}

// BacklogCheck is the backlog-consistency result.
type BacklogCheck struct {
	ActiveWork     SectionCount
	RecentlyClosed SectionCount
	ParkedWork     SectionCount
	OrphanedFiles  []string
	OrphanedParked []string
	Status         CheckStatus
}

// FileFinding names a closed record plus context for a report line.
type FileFinding struct {
	ID       string
	Sections []string
	Date     string
	FilePath string
	Reason   string
}

// ClosedTaskCheck is the closed-task-completeness result.
type ClosedTaskCheck struct {
	Total              int
	Skipped            int
	WithOutcome        int
	WithVerification   int
	WithCloseChecklist int
	Failing            []FileFinding
	Warnings           []FileFinding
	SupportingSkipped  []FileFinding
	UntypedLegacy      []FileFinding
	UntypedCurrent     []FileFinding
	Status             CheckStatus
}

// EvidenceCheck is the evidence-integrity result.
type EvidenceCheck struct {
	Total  int
	Valid  int
	Broken []BrokenLink
	Status CheckStatus
}

// BrokenLink is one unresolved link with its source record.
type BrokenLink struct {
	ID   string
	Link string
}

// Claim is a checklist line in a Verification section.
type Claim struct {
	ID         string
	Claim      string
	FilePath   string
	LineNumber int
}

// VerificationCheck is the verification-clarity result.
type VerificationCheck struct {
	Total     int
	Proven    int
	Unproven  []Claim
	Ambiguous []Claim
	Status    CheckStatus
}

// DeadLinksCheck is the dead-links result over active work files.
type DeadLinksCheck struct {
	Scanned int
	Broken  []BrokenLink
	Status  CheckStatus
}

// Report is the combined validation outcome.
type Report struct {
	Backlog      BacklogCheck
	Completeness ClosedTaskCheck
	Evidence     EvidenceCheck
	Clarity      VerificationCheck
	DeadLinks    DeadLinksCheck
	Passed       int
	Failed       int
	Warnings     int
	Overall      CheckStatus // PASS or FAIL
}

// Summarize fills the counters and overall status from the check statuses.
func (r *Report) Summarize() {
	statuses := []CheckStatus{
		r.Backlog.Status, r.Completeness.Status, r.Evidence.Status,
		r.Clarity.Status, r.DeadLinks.Status,
	}
	for _, s := range statuses {
		switch s {
		case StatusPass:
			r.Passed++
		case StatusFail:
			r.Failed++
		}
	}
	for _, s := range []CheckStatus{r.Completeness.Status, r.Evidence.Status, r.Clarity.Status} {
		if s == StatusWarn {
			r.Warnings++
		}
	}
	r.Overall = StatusPass
	if r.Failed > 0 {
		r.Overall = StatusFail
	}
}
