package normalization

import (
	"errors"
	"fmt"
	"io/fs"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// MigrationIO is the filesystem seam used by work-item folder migration.
// The normalization package owns the migration workflow while the parent
// package supplies the concrete filesystem operations.
type MigrationIO struct {
	CanonicalFolderID func(content, fallback string) string
	CollectFiles      func(dir string, skipReadme bool) []string
	ReadDir           func(path string) ([]fs.DirEntry, error)
	ReadFile          func(path string) ([]byte, error)
	WriteFile         func(path string, data []byte, perm fs.FileMode) error
	MkdirAll          func(path string, perm fs.FileMode) error
	Rename            func(oldPath, newPath string) error
	Remove            func(path string) error
	Stat              func(path string) (fs.FileInfo, error)
	ToRepoRelative    func(repoRoot, absolutePath string) string
}

var filesStemRe = regexp.MustCompile(`(?i)^(.*)-files$`)

func moveMarkdownFile(io MigrationIO, oldPath, newPath string, transform func(string) string) error {
	raw, err := io.ReadFile(oldPath)
	if err != nil {
		return err
	}
	next := string(raw)
	if transform != nil {
		next = transform(next)
	}
	if err := io.MkdirAll(filepath.Dir(newPath), 0o755); err != nil {
		return err
	}
	if err := io.WriteFile(newPath, []byte(next), 0o644); err != nil {
		return err
	}
	return io.Remove(oldPath)
}

func moveArtifactDir(io MigrationIO, oldPath, newPath string) error {
	if _, err := io.Stat(oldPath); err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil
		}
		return err
	}
	if _, err := io.Stat(newPath); err == nil {
		return fmt.Errorf("migration target already exists: %s", newPath)
	}
	if err := io.MkdirAll(filepath.Dir(newPath), 0o755); err != nil {
		return err
	}
	return io.Rename(oldPath, newPath)
}

func applyDirRename(io MigrationIO, dirPath, targetDir, workRoot string, touched map[string]struct{}, movedPlans *[]MovedPlan) error {
	if dirPath == targetDir {
		return nil
	}
	if _, err := io.Stat(targetDir); err == nil {
		return fmt.Errorf("migration target already exists: %s", targetDir)
	}
	if err := io.MkdirAll(filepath.Dir(targetDir), 0o755); err != nil {
		return err
	}
	if err := io.Rename(dirPath, targetDir); err != nil {
		return err
	}

	for _, newFile := range io.CollectFiles(targetDir, false) {
		relSuffix, err := filepath.Rel(targetDir, newFile)
		if err != nil {
			return err
		}
		oldFile := filepath.Join(dirPath, relSuffix)
		raw, err := io.ReadFile(newFile)
		if err != nil {
			return err
		}
		updated := RewriteMovedMarkdown(string(raw), oldFile, newFile)

		workRel := filepath.ToSlash(strings.TrimPrefix(strings.TrimPrefix(newFile, workRoot), string(filepath.Separator)))
		if strings.HasSuffix(newFile, string(filepath.Separator)+"plan.md") {
			updated = NormalizeMovedPlanContent(updated, workRel, oldFile, newFile)
			oldRel := filepath.ToSlash(strings.TrimPrefix(strings.TrimPrefix(oldFile, workRoot), string(filepath.Separator)))
			*movedPlans = append(*movedPlans, MovedPlan{OldRel: oldRel, NewRel: workRel})
		}
		if strings.HasSuffix(newFile, string(filepath.Separator)+"files.md") {
			baseDir := filepath.Dir(newFile)
			planRel := filepath.ToSlash(strings.TrimPrefix(strings.TrimPrefix(filepath.Join(baseDir, "plan.md"), workRoot), string(filepath.Separator)))
			updated = NormalizeMovedFilesContent(updated, planRel, oldFile, newFile)
		}

		if err := io.WriteFile(newFile, []byte(updated), 0o644); err != nil {
			return err
		}
		touched[newFile] = struct{}{}
	}
	return nil
}

func applyFlatPlanMove(io MigrationIO, scopeRoot, oldPath, targetDir, workRoot string, touched map[string]struct{}, movedPlans *[]MovedPlan) error {
	newPath := filepath.Join(targetDir, "plan.md")
	if _, err := io.Stat(newPath); err == nil {
		return fmt.Errorf("migration target already exists: %s", newPath)
	}

	workRel, err := filepath.Rel(workRoot, newPath)
	if err != nil {
		return err
	}
	if err := moveMarkdownFile(io, oldPath, newPath, func(content string) string {
		return NormalizeMovedPlanContent(content, workRel, oldPath, newPath)
	}); err != nil {
		return err
	}
	touched[newPath] = struct{}{}

	oldRel, err := filepath.Rel(workRoot, oldPath)
	if err != nil {
		return err
	}
	*movedPlans = append(*movedPlans, MovedPlan{OldRel: filepath.ToSlash(oldRel), NewRel: filepath.ToSlash(workRel)})

	stem := strings.TrimSuffix(filepath.Base(oldPath), filepath.Ext(oldPath))
	filesPath := filepath.Join(scopeRoot, stem+"-files.md")
	if _, err := io.Stat(filesPath); err == nil {
		target := filepath.Join(targetDir, "files.md")
		planRel := filepath.ToSlash(workRel)
		if err := moveMarkdownFile(io, filesPath, target, func(content string) string {
			return NormalizeMovedFilesContent(content, planRel, filesPath, target)
		}); err != nil {
			return err
		}
		touched[target] = struct{}{}
	}

	for _, name := range []string{"references", "evidence"} {
		oldDir := filepath.Join(scopeRoot, stem+"-"+name)
		newDir := filepath.Join(targetDir, name)
		if err := moveArtifactDir(io, oldDir, newDir); err != nil {
			return err
		}
		if _, err := io.Stat(newDir); err == nil {
			touched[newDir] = struct{}{}
		}
	}
	return nil
}

func applyRemainingSidecars(io MigrationIO, scopeRoot, workRoot string, touched map[string]struct{}) error {
	entries, err := io.ReadDir(scopeRoot)
	if err != nil {
		return nil
	}
	for _, entry := range entries {
		if entry.IsDir() || entry.Name() == "README.md" || !strings.HasSuffix(entry.Name(), ".md") {
			continue
		}
		stem := strings.TrimSuffix(entry.Name(), ".md")
		match := filesStemRe.FindStringSubmatch(stem)
		if match == nil {
			continue
		}
		itemStem := match[1]
		oldPath := filepath.Join(scopeRoot, entry.Name())
		targetDir := filepath.Join(scopeRoot, itemStem)
		target := filepath.Join(targetDir, "files.md")
		if _, err := io.Stat(target); err == nil {
			return fmt.Errorf("migration target already exists: %s", target)
		}
		planRel, err := filepath.Rel(workRoot, filepath.Join(targetDir, "plan.md"))
		if err != nil {
			return err
		}
		if err := moveMarkdownFile(io, oldPath, target, func(content string) string {
			return NormalizeMovedFilesContent(content, filepath.ToSlash(planRel), oldPath, target)
		}); err != nil {
			return err
		}
		touched[target] = struct{}{}
	}
	return nil
}

func rewriteBacklogPlanLinks(io MigrationIO, backlogPath string, movedPlans []MovedPlan, touched map[string]struct{}) error {
	if len(movedPlans) == 0 {
		return nil
	}
	raw, err := io.ReadFile(backlogPath)
	if err != nil {
		return err
	}
	content := string(raw)
	updated := content
	for _, move := range movedPlans {
		updated = strings.ReplaceAll(updated, "("+move.OldRel+")", "("+move.NewRel+")")
	}
	if updated == content {
		return nil
	}
	if err := io.WriteFile(backlogPath, []byte(updated), 0o644); err != nil {
		return err
	}
	touched[backlogPath] = struct{}{}
	return nil
}

// ApplyWorkItemFolderMigration performs the folder-per-item migration through
// the explicit filesystem seam and returns deterministic changed-file output.
func ApplyWorkItemFolderMigration(repoRoot string, io MigrationIO) (ApplyResult, error) {
	result := ApplyResult{}
	workRoot := filepath.Join(repoRoot, ".hawp", "work")
	touched := map[string]struct{}{}
	var movedPlans []MovedPlan

	for _, scope := range []string{"active", "parked"} {
		scopeRoot := filepath.Join(workRoot, scope)
		entries, err := io.ReadDir(scopeRoot)
		if err != nil {
			continue
		}

		for _, entry := range entries {
			if !entry.IsDir() || entry.Name() == "README.md" {
				continue
			}
			dirPath := filepath.Join(scopeRoot, entry.Name())
			planPath := filepath.Join(dirPath, "plan.md")
			raw, err := io.ReadFile(planPath)
			if err != nil {
				continue
			}
			canonicalID := io.CanonicalFolderID(string(raw), entry.Name())
			if canonicalID == "" || canonicalID == entry.Name() {
				continue
			}
			targetDir := filepath.Join(scopeRoot, canonicalID)
			if err := applyDirRename(io, dirPath, targetDir, workRoot, touched, &movedPlans); err != nil {
				return result, err
			}
		}

		entries, err = io.ReadDir(scopeRoot)
		if err != nil {
			continue
		}
		for _, entry := range entries {
			if entry.IsDir() || entry.Name() == "README.md" || !strings.HasSuffix(entry.Name(), ".md") {
				continue
			}
			stem := strings.TrimSuffix(entry.Name(), filepath.Ext(entry.Name()))
			if filesStemRe.MatchString(stem) {
				continue
			}
			oldPath := filepath.Join(scopeRoot, entry.Name())
			raw, err := io.ReadFile(oldPath)
			if err != nil {
				return result, err
			}
			canonicalID := io.CanonicalFolderID(string(raw), stem)
			if canonicalID == "" {
				canonicalID = stem
			}
			targetDir := filepath.Join(scopeRoot, canonicalID)
			if targetDir == filepath.Dir(oldPath) && filepath.Base(oldPath) == "plan.md" {
				continue
			}
			if err := applyFlatPlanMove(io, scopeRoot, oldPath, targetDir, workRoot, touched, &movedPlans); err != nil {
				return result, err
			}
		}

		if err := applyRemainingSidecars(io, scopeRoot, workRoot, touched); err != nil {
			return result, err
		}
	}

	backlogPath := filepath.Join(workRoot, "BACKLOG.md")
	if err := rewriteBacklogPlanLinks(io, backlogPath, movedPlans, touched); err != nil {
		return result, err
	}

	for path := range touched {
		result.ChangedFiles = append(result.ChangedFiles, io.ToRepoRelative(repoRoot, path))
	}
	sort.Strings(result.ChangedFiles)
	return result, nil
}
