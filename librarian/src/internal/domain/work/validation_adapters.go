package work

import "github.com/sentzunhat/hawp/librarian/src/internal/domain/work/validation"

func IDSet(rows []BacklogRow) map[string]struct{} { return validation.IDSet(rows) }

func (w *WorkSource) CheckDeadLinks(workDir string) DeadLinksCheck {
	return validation.CheckDeadLinks(workDir, validation.Source{
		Exists: w.Exists, ReadDir: w.ReadDir, ReadFile: w.ReadFile,
		ToRepoRelative: w.ToRepoRelative,
	})
}
