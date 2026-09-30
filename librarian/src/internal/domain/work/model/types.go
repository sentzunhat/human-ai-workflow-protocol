// Package model contains reusable work-domain values and validation results.
// It has no filesystem or orchestration dependencies so capability packages
// can depend on these values without importing the compatibility parent.
package model

type CheckStatus string

const (
	StatusPass CheckStatus = "PASS"
	StatusFail CheckStatus = "FAIL"
	StatusWarn CheckStatus = "WARN"
)

type BacklogRow struct{ ID, Type, Title, Status, Detail string }
type Backlog struct{ Active, Closed, Parked []BacklogRow }
type SectionCount struct {
	Total   int
	Found   int
	Missing []string
}
type BacklogCheck struct {
	ActiveWork, RecentlyClosed, ParkedWork                 SectionCount
	OrphanedFiles, OrphanedParked, NonCanonicalActiveItems []string
	Status                                                 CheckStatus
}
type FileFinding struct {
	ID                     string
	Sections               []string
	Date, FilePath, Reason string
}
type ClosedTaskCheck struct {
	Total, Skipped, WithOutcome, WithVerification, WithCloseChecklist   int
	Failing, Warnings, SupportingSkipped, UntypedLegacy, UntypedCurrent []FileFinding
	Status                                                              CheckStatus
}
type EvidenceCheck struct {
	Total, Valid int
	Broken       []BrokenLink
	Status       CheckStatus
}
type BrokenLink struct{ ID, Link string }
type Claim struct {
	ID, Claim, FilePath string
	LineNumber          int
}
type VerificationCheck struct {
	Total, Proven       int
	Unproven, Ambiguous []Claim
	Status              CheckStatus
}
type DeadLinksCheck struct {
	Scanned int
	Broken  []BrokenLink
	Status  CheckStatus
}
type Report struct {
	Backlog                  BacklogCheck
	Completeness             ClosedTaskCheck
	Evidence                 EvidenceCheck
	Clarity                  VerificationCheck
	DeadLinks                DeadLinksCheck
	Passed, Failed, Warnings int
	Overall                  CheckStatus
}

func (r *Report) Summarize() {
	r.Passed, r.Failed, r.Warnings = 0, 0, 0
	statuses := []CheckStatus{r.Backlog.Status, r.Completeness.Status, r.Evidence.Status, r.Clarity.Status, r.DeadLinks.Status}
	for _, status := range statuses {
		switch status {
		case StatusPass:
			r.Passed++
		case StatusFail:
			r.Failed++
		}
	}
	for _, status := range []CheckStatus{r.Completeness.Status, r.Evidence.Status, r.Clarity.Status} {
		if status == StatusWarn {
			r.Warnings++
		}
	}
	r.Overall = StatusPass
	if r.Failed > 0 {
		r.Overall = StatusFail
	}
}
