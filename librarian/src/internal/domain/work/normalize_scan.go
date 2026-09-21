package work

import "github.com/sentzunhat/hawp/librarian/src/internal/domain/work/normalization"

// Compatibility aliases keep the normalization application API stable while
// scan parsing lives in its capability-local package.
type NormalizeSection = normalization.NormalizeSection
type NormalizeRow = normalization.NormalizeRow
type NormalizeBacklog = normalization.NormalizeBacklog
type PlanFileRecord = normalization.PlanFileRecord
type PlanScan = normalization.PlanScan
type FixOperation = normalization.FixOperation
type BlockedInfo = normalization.BlockedInfo
type ResearchItem = normalization.ResearchItem
type ApplyResult = normalization.ApplyResult
type FixPlan = normalization.FixPlan
type SyncPlanStep = normalization.SyncPlanStep
type DetectionReport = normalization.DetectionReport

func BuildDetectionReport(scannedAt, backlogPath string, filesScanned, itemsAnalyzed int, operations []FixOperation) DetectionReport {
	return normalization.BuildDetectionReport(scannedAt, backlogPath, filesScanned, itemsAnalyzed, operations)
}

func RenderTextReport(report DetectionReport) string { return normalization.RenderTextReport(report) }
func RenderJSONReport(report DetectionReport) (string, error) {
	return normalization.RenderJSONReport(report)
}
func RenderJSONValue(value any) (string, error) { return normalization.RenderJSONValue(value) }

func EvaluateRules(repoRoot, workRoot, backlogPath string, backlog *NormalizeBacklog, scan *PlanScan, source *WorkSource) []FixOperation {
	return normalization.EvaluateRules(repoRoot, workRoot, backlogPath, backlog, scan, normalization.Source{Exists: source.Exists, ToRepoRelative: source.ToRepoRelative})
}

func hasUnprovenChecklistMarker(content string) bool {
	return normalization.HasUnprovenChecklistMarker(content)
}
func hasHeadingNamed(content, heading string) bool {
	return normalization.HasHeadingNamed(content, heading)
}
func AmbiguousVerificationClaims(content string) []string {
	return normalization.AmbiguousClaims(content)
}

const (
	SectionActive  = normalization.SectionActive
	SectionBlocked = normalization.SectionBlocked
	SectionClosed  = normalization.SectionClosed
	SectionOther   = normalization.SectionOther
)

func readBacklogID(content string) string { return normalization.ReadBacklogID(content) }
