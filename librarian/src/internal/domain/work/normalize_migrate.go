package work

import (
	"path/filepath"

	"github.com/sentzunhat/hawp/librarian/src/internal/domain/work/normalization"
)

func canonicalFolderID(content, fallback string) string {
	return normalization.CanonicalFolderID(content, fallback, readBacklogID(content), ExtractIDFromFilename(fallback))
}

// PreviewWorkItemFolderMigration runs the real migration against an isolated
// temp copy of .hawp/work and returns the files that would change.
func (w *WorkSource) PreviewWorkItemFolderMigration(repoRoot string) (ApplyResult, error) {
	tempRoot, err := w.MkdirTemp("", "hawp-work-migrate-preview-")
	if err != nil {
		return ApplyResult{}, err
	}
	defer w.RemoveAll(tempRoot)

	srcWorkRoot := filepath.Join(repoRoot, ".hawp", "work")
	dstWorkRoot := filepath.Join(tempRoot, ".hawp", "work")
	if err := w.MkdirAll(dstWorkRoot, 0o755); err != nil {
		return ApplyResult{}, err
	}
	if err := copyTree(w, srcWorkRoot, dstWorkRoot); err != nil {
		return ApplyResult{}, err
	}

	return w.ApplyWorkItemFolderMigration(tempRoot)
}

// ApplyWorkItemFolderMigration keeps the parent domain/work API stable while
// normalization owns migration orchestration behind an explicit filesystem
// seam.
func (w *WorkSource) ApplyWorkItemFolderMigration(repoRoot string) (ApplyResult, error) {
	return normalization.ApplyWorkItemFolderMigration(repoRoot, migrationIO(w))
}

func migrationIO(w *WorkSource) normalization.MigrationIO {
	return normalization.MigrationIO{
		CanonicalFolderID:      canonicalFolderID,
		CollectFiles:           w.CollectFiles,
		ReadDir:                w.ReadDir,
		ReadFile:               w.ReadFile,
		WriteFile:              w.WriteFile,
		MkdirAll:               w.MkdirAll,
		Rename:                 w.Rename,
		Remove:                 w.Remove,
		Stat:                   w.Stat,
		RejectSymlinkAncestors: w.RejectSymlinkAncestors,
		ToRepoRelative:         w.ToRepoRelative,
	}
}

func copyTree(w *WorkSource, srcRoot, dstRoot string) error {
	entries, err := w.ReadDir(srcRoot)
	if err != nil {
		return err
	}
	for _, entry := range entries {
		srcPath := filepath.Join(srcRoot, entry.Name())
		dstPath := filepath.Join(dstRoot, entry.Name())
		if entry.IsDir() {
			if err := w.MkdirAll(dstPath, 0o755); err != nil {
				return err
			}
			if err := copyTree(w, srcPath, dstPath); err != nil {
				return err
			}
			continue
		}
		data, err := w.ReadFile(srcPath)
		if err != nil {
			return err
		}
		if err := w.WriteFile(dstPath, data, 0o644); err != nil {
			return err
		}
	}
	return nil
}
