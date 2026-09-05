package work

import "github.com/sentzunhat/hawp/librarian/src/internal/domain/work/model"

// Compatibility aliases preserve the historical domain/work API while the
// reusable work values live in their capability-neutral model package.
type CheckStatus = model.CheckStatus
type BacklogRow = model.BacklogRow
type Backlog = model.Backlog
type SectionCount = model.SectionCount
type BacklogCheck = model.BacklogCheck
type FileFinding = model.FileFinding
type ClosedTaskCheck = model.ClosedTaskCheck
type EvidenceCheck = model.EvidenceCheck
type BrokenLink = model.BrokenLink
type Claim = model.Claim
type VerificationCheck = model.VerificationCheck
type DeadLinksCheck = model.DeadLinksCheck
type Report = model.Report

const (
	StatusPass = model.StatusPass
	StatusFail = model.StatusFail
	StatusWarn = model.StatusWarn
)
