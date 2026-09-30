package work

import (
	"os"

	"github.com/sentzunhat/hawp/librarian/src/internal/domain/work/normalization"
	"github.com/sentzunhat/hawp/librarian/src/internal/domain/work/validation"
)

// These test-only shims preserve coverage for the historical convenience
// functions while production callers use injected capability sources.
func testValidationSource() validation.Source {
	return validation.Source{Exists: func(path string) bool { _, err := os.Stat(path); return err == nil }, ReadDir: os.ReadDir, ReadFile: os.ReadFile, Stat: os.Stat, ToRepoRelative: func(_, path string) string { return path }}
}

func CheckBacklogConsistency(workDir string, backlog *Backlog) BacklogCheck {
	return validation.CheckBacklogConsistency(workDir, backlog, testValidationSource())
}
func CheckClosedTaskCompleteness(workDir string) ClosedTaskCheck {
	return validation.CheckClosedTaskCompleteness(workDir, testValidationSource())
}
func CollectClosedPlanFiles(closedDir string) []string {
	return validation.CollectClosedPlanFiles(closedDir, testValidationSource())
}
func CheckEvidenceIntegrity(workDir string, files []string) EvidenceCheck {
	return validation.CheckEvidenceIntegrity(workDir, files, testValidationSource())
}
func CheckVerificationClarity(files []string) VerificationCheck {
	return validation.CheckVerificationClarity(files, testValidationSource())
}

func ParseNormalizeBacklog(path string) (*NormalizeBacklog, error) {
	return normalization.ParseNormalizeBacklog(path, normalization.ScanSource{ReadFile: os.ReadFile, ReadDir: os.ReadDir, Stat: os.Stat})
}
func ScanPlanFiles(root string) *PlanScan {
	return normalization.ScanPlanFiles(root, normalization.ScanSource{ReadFile: os.ReadFile, ReadDir: os.ReadDir, Stat: os.Stat})
}
