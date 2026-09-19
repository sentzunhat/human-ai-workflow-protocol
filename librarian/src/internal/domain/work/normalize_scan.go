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

const (
	SectionActive  = normalization.SectionActive
	SectionBlocked = normalization.SectionBlocked
	SectionClosed  = normalization.SectionClosed
	SectionOther   = normalization.SectionOther
)

func ParseNormalizeBacklog(path string) (*NormalizeBacklog, error) {
	return normalization.ParseNormalizeBacklog(path)
}

func ScanPlanFiles(workRoot string) *PlanScan { return normalization.ScanPlanFiles(workRoot) }

func readBacklogID(content string) string  { return normalization.ReadBacklogID(content) }
func walkPlanMarkdown(dir string) []string { return normalization.WalkPlanMarkdown(dir) }
