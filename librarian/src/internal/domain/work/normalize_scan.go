package work

import "github.com/sentzunhat/hawp/librarian/src/internal/domain/work/normalization"

// Compatibility aliases keep the normalization application API stable while
// scan parsing lives in its capability-local package.
type NormalizeSection = normalization.NormalizeSection
type NormalizeRow = normalization.NormalizeRow
type NormalizeBacklog = normalization.NormalizeBacklog
type PlanFileRecord = normalization.PlanFileRecord
type PlanScan = normalization.PlanScan

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
