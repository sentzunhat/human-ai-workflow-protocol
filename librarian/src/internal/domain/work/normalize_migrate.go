package work

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/sentzunhat/hawp/librarian/src/internal/infrastructure/markdown"
	"github.com/sentzunhat/hawp/librarian/src/internal/infrastructure/repo"
)

var (
	filesStemRe       = regexp.MustCompile(`(?i)^(.*)-files$`)
	workItemLineRe    = regexp.MustCompile(`(?m)^\*\*Work Item:\*\*\s+.*$`)
	planFileLineRe    = regexp.MustCompile(`(?m)^\*\*Plan file:\*\*\s+.*$`)
	markdownLinkRefRe = regexp.MustCompile(`\[([^\]]*)\]\(([^)]+)\)`)
)

type movedPlan struct {
	oldRel string
	newRel string
}

func canonicalFolderID(content, fallback string) string {
	if m := uuidFieldRe.FindStringSubmatch(content); m != nil {
		raw := strings.ToLower(strings.TrimSpace(m[1]))
		if fullUUIDRe.MatchString(raw) {
			return raw[:8]
		}
		if shortUUIDRe.MatchString(raw) {
			return raw
		}
	}
	if id := readBacklogID(content); id != "" {
		return id
	}
	if id := ExtractIDFromFilename(fallback); id != "" {
		return id
	}
	base := filepath.Base(fallback)
	if strings.HasSuffix(strings.ToLower(base), ".md") {
		base = strings.TrimSuffix(base, filepath.Ext(base))
	}
	if base == "" {
		return fallback
	}
	return base
}

func rewriteMovedMarkdown(content, oldPath, newPath string) string {
	blanked := markdown.BlankFences(content)
	matches := markdownLinkRefRe.FindAllStringSubmatchIndex(blanked, -1)
	if len(matches) == 0 {
		return content
	}

	var out strings.Builder
	last := 0
	for _, idx := range matches {
		hrefStart, hrefEnd := idx[4], idx[5]
		href := content[hrefStart:hrefEnd]
		nextHref := rewriteLocalHref(href, oldPath, newPath)
		out.WriteString(content[last:hrefStart])
		out.WriteString(nextHref)
		last = hrefEnd
	}
	out.WriteString(content[last:])
	return out.String()
}

func rewriteLocalHref(href, oldPath, newPath string) string {
	if !markdown.IsLocalHref(href) {
		return href
	}
	pathPart, anchor, hasAnchor := strings.Cut(href, "#")
	target := filepath.Clean(filepath.Join(filepath.Dir(oldPath), filepath.FromSlash(pathPart)))
	rel, err := filepath.Rel(filepath.Dir(newPath), target)
	if err != nil {
		return href
	}
	rewritten := filepath.ToSlash(rel)
	if hasAnchor {
		rewritten += "#" + anchor
	}
	return rewritten
}

func normalizeMovedPlanContent(content, workRel, oldPath, newPath string) string {
	updated := rewriteMovedMarkdown(content, oldPath, newPath)
	if planFileLineRe.MatchString(updated) {
		updated = planFileLineRe.ReplaceAllString(updated, "**Plan file:** work/"+filepath.ToSlash(workRel))
	}
	return updated
}

func normalizeMovedFilesContent(content, planRel, oldPath, newPath string) string {
	updated := rewriteMovedMarkdown(content, oldPath, newPath)
	if workItemLineRe.MatchString(updated) {
		updated = workItemLineRe.ReplaceAllString(updated, "**Work Item:** .hawp/work/"+filepath.ToSlash(planRel))
	}
	return updated
}

func moveMarkdownFile(oldPath, newPath string, transform func(string) string) error {
	raw, err := os.ReadFile(oldPath)
	if err != nil {
		return err
	}
	next := string(raw)
	if transform != nil {
		next = transform(next)
	}
	if err := os.MkdirAll(filepath.Dir(newPath), 0o755); err != nil {
		return err
	}
	if err := os.WriteFile(newPath, []byte(next), 0o644); err != nil {
		return err
	}
	return os.Remove(oldPath)
}

func moveArtifactDir(oldPath, newPath string) error {
	if _, err := os.Stat(oldPath); err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	if _, err := os.Stat(newPath); err == nil {
		return fmt.Errorf("migration target already exists: %s", newPath)
	}
	if err := os.MkdirAll(filepath.Dir(newPath), 0o755); err != nil {
		return err
	}
	return os.Rename(oldPath, newPath)
}

func applyDirRename(dirPath, targetDir, workRoot string, touched map[string]struct{}, movedPlans *[]movedPlan) error {
	if dirPath == targetDir {
		return nil
	}
	if _, err := os.Stat(targetDir); err == nil {
		return fmt.Errorf("migration target already exists: %s", targetDir)
	}
	if err := os.MkdirAll(filepath.Dir(targetDir), 0o755); err != nil {
		return err
	}
	if err := os.Rename(dirPath, targetDir); err != nil {
		return err
	}

	for _, newFile := range markdown.CollectFiles(targetDir, false) {
		relSuffix, err := filepath.Rel(targetDir, newFile)
		if err != nil {
			return err
		}
		oldFile := filepath.Join(dirPath, relSuffix)
		raw, err := os.ReadFile(newFile)
		if err != nil {
			return err
		}
		updated := rewriteMovedMarkdown(string(raw), oldFile, newFile)

		workRel := filepath.ToSlash(strings.TrimPrefix(strings.TrimPrefix(newFile, workRoot), string(filepath.Separator)))
		planRel := workRel
		if strings.HasSuffix(newFile, string(filepath.Separator)+"plan.md") {
			updated = normalizeMovedPlanContent(updated, workRel, oldFile, newFile)
			oldRel := filepath.ToSlash(strings.TrimPrefix(strings.TrimPrefix(oldFile, workRoot), string(filepath.Separator)))
			*movedPlans = append(*movedPlans, movedPlan{oldRel: oldRel, newRel: workRel})
		}
		if strings.HasSuffix(newFile, string(filepath.Separator)+"files.md") {
			baseDir := filepath.Dir(newFile)
			planRel = filepath.ToSlash(strings.TrimPrefix(strings.TrimPrefix(filepath.Join(baseDir, "plan.md"), workRoot), string(filepath.Separator)))
			updated = normalizeMovedFilesContent(updated, planRel, oldFile, newFile)
		}

		if err := os.WriteFile(newFile, []byte(updated), 0o644); err != nil {
			return err
		}
		touched[newFile] = struct{}{}
	}
	return nil
}

func applyFlatPlanMove(scopeRoot, oldPath, targetDir, canonicalID, workRoot string, touched map[string]struct{}, movedPlans *[]movedPlan) error {
	newPath := filepath.Join(targetDir, "plan.md")
	if _, err := os.Stat(newPath); err == nil {
		return fmt.Errorf("migration target already exists: %s", newPath)
	}

	workRel, err := filepath.Rel(workRoot, newPath)
	if err != nil {
		return err
	}
	if err := moveMarkdownFile(oldPath, newPath, func(content string) string {
		return normalizeMovedPlanContent(content, workRel, oldPath, newPath)
	}); err != nil {
		return err
	}
	touched[newPath] = struct{}{}

	oldRel, err := filepath.Rel(workRoot, oldPath)
	if err != nil {
		return err
	}
	*movedPlans = append(*movedPlans, movedPlan{
		oldRel: filepath.ToSlash(oldRel),
		newRel: filepath.ToSlash(workRel),
	})

	stem := strings.TrimSuffix(filepath.Base(oldPath), filepath.Ext(oldPath))
	filesPath := filepath.Join(scopeRoot, stem+"-files.md")
	if _, err := os.Stat(filesPath); err == nil {
		target := filepath.Join(targetDir, "files.md")
		planRel := filepath.ToSlash(workRel)
		if err := moveMarkdownFile(filesPath, target, func(content string) string {
			return normalizeMovedFilesContent(content, planRel, filesPath, target)
		}); err != nil {
			return err
		}
		touched[target] = struct{}{}
	}

	for _, name := range []string{"references", "evidence"} {
		oldDir := filepath.Join(scopeRoot, stem+"-"+name)
		newDir := filepath.Join(targetDir, name)
		if err := moveArtifactDir(oldDir, newDir); err != nil {
			return err
		}
		if _, err := os.Stat(newDir); err == nil {
			touched[newDir] = struct{}{}
		}
	}

	_ = canonicalID
	return nil
}

func applyRemainingSidecars(scopeRoot, workRoot string, touched map[string]struct{}) error {
	entries, err := os.ReadDir(scopeRoot)
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
		if _, err := os.Stat(target); err == nil {
			return fmt.Errorf("migration target already exists: %s", target)
		}
		planRel, err := filepath.Rel(workRoot, filepath.Join(targetDir, "plan.md"))
		if err != nil {
			return err
		}
		if err := moveMarkdownFile(oldPath, target, func(content string) string {
			return normalizeMovedFilesContent(content, filepath.ToSlash(planRel), oldPath, target)
		}); err != nil {
			return err
		}
		touched[target] = struct{}{}
	}
	return nil
}

func rewriteBacklogPlanLinks(backlogPath string, movedPlans []movedPlan, touched map[string]struct{}) error {
	if len(movedPlans) == 0 {
		return nil
	}
	raw, err := os.ReadFile(backlogPath)
	if err != nil {
		return err
	}
	content := string(raw)
	updated := content
	for _, move := range movedPlans {
		updated = strings.ReplaceAll(updated, "("+move.oldRel+")", "("+move.newRel+")")
	}
	if updated == content {
		return nil
	}
	if err := os.WriteFile(backlogPath, []byte(updated), 0o644); err != nil {
		return err
	}
	touched[backlogPath] = struct{}{}
	return nil
}

// ApplyWorkItemFolderMigration migrates active/parked work-item plans and
// sidecar artifacts into folder-per-item layout, preserving relative links.
func ApplyWorkItemFolderMigration(repoRoot string) (ApplyResult, error) {
	result := ApplyResult{}
	workRoot := filepath.Join(repoRoot, ".hawp", "work")
	touched := map[string]struct{}{}
	var movedPlans []movedPlan

	for _, scope := range []string{"active", "parked"} {
		scopeRoot := filepath.Join(workRoot, scope)
		entries, err := os.ReadDir(scopeRoot)
		if err != nil {
			continue
		}

		for _, entry := range entries {
			if !entry.IsDir() || entry.Name() == "README.md" {
				continue
			}
			dirPath := filepath.Join(scopeRoot, entry.Name())
			planPath := filepath.Join(dirPath, "plan.md")
			raw, err := os.ReadFile(planPath)
			if err != nil {
				continue
			}
			canonicalID := canonicalFolderID(string(raw), entry.Name())
			if canonicalID == "" || canonicalID == entry.Name() {
				continue
			}
			targetDir := filepath.Join(scopeRoot, canonicalID)
			if err := applyDirRename(dirPath, targetDir, workRoot, touched, &movedPlans); err != nil {
				return result, err
			}
		}

		entries, err = os.ReadDir(scopeRoot)
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
			raw, err := os.ReadFile(oldPath)
			if err != nil {
				return result, err
			}
			canonicalID := canonicalFolderID(string(raw), stem)
			if canonicalID == "" {
				canonicalID = stem
			}
			targetDir := filepath.Join(scopeRoot, canonicalID)
			if targetDir == filepath.Dir(oldPath) && filepath.Base(oldPath) == "plan.md" {
				continue
			}
			if err := applyFlatPlanMove(scopeRoot, oldPath, targetDir, canonicalID, workRoot, touched, &movedPlans); err != nil {
				return result, err
			}
		}

		if err := applyRemainingSidecars(scopeRoot, workRoot, touched); err != nil {
			return result, err
		}
	}

	backlogPath := filepath.Join(workRoot, "BACKLOG.md")
	if err := rewriteBacklogPlanLinks(backlogPath, movedPlans, touched); err != nil {
		return result, err
	}

	for path := range touched {
		result.ChangedFiles = append(result.ChangedFiles, repo.ToRepoRelative(repoRoot, path))
	}
	sort.Strings(result.ChangedFiles)
	return result, nil
}
