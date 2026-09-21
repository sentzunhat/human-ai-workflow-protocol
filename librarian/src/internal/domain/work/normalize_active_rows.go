package work

import "github.com/sentzunhat/hawp/librarian/src/internal/domain/work/normalization"

// ApplyCompletedActiveRowCleanup removes completed Active Work rows only when
// they already point at an existing closed plan.
func (w *WorkSource) ApplyCompletedActiveRowCleanup(repoRoot string) (ApplyResult, error) {
	return normalization.ApplyCompletedActiveRowCleanup(repoRoot, normalization.ActiveSource{
		ScanSource:     normalization.ScanSource{ReadDir: w.ReadDir, ReadFile: w.ReadFile, Stat: w.Stat},
		ToRepoRelative: w.ToRepoRelative,
		EvalSymlinks:   w.EvalSymlinks, Lstat: w.Lstat, WriteFile: w.WriteFile,
	})
}
