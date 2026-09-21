package work

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	domainwork "github.com/sentzunhat/hawp/librarian/src/internal/domain/work"
	"github.com/sentzunhat/hawp/librarian/src/internal/domain/work/normalization"
	"github.com/sentzunhat/hawp/librarian/src/internal/infrastructure/filesystem"
	"github.com/sentzunhat/hawp/librarian/src/internal/infrastructure/repo"
	reposwork "github.com/sentzunhat/hawp/librarian/src/internal/infrastructure/repositories/work"
)

// PreviewDuplicateLinks reports cross-reference edits without modifying files.
func PreviewDuplicateLinks(repoRoot string) (domainwork.ApplyResult, error) {
	return linkDuplicatePlans(repoRoot, false)
}

// ApplyDuplicateLinks preserves every plan and adds bidirectional references
// only when an unreferenced working record has one archived counterpart.
func ApplyDuplicateLinks(repoRoot string) (domainwork.ApplyResult, error) {
	return linkDuplicatePlans(repoRoot, true)
}

func linkDuplicatePlans(repoRoot string, apply bool) (domainwork.ApplyResult, error) {
	result := domainwork.ApplyResult{}
	workRoot := filepath.Join(repoRoot, ".hawp", "work")
	if err := filesystem.RejectSymlinksInTree(repoRoot, workRoot); err != nil {
		return result, fmt.Errorf("unsafe work tree: %w", err)
	}
	backlog, err := reposwork.ReadBacklog(filepath.Join(workRoot, "BACKLOG.md"))
	if err != nil {
		return result, err
	}
	referenced := domainwork.IDSet(append(append([]domainwork.BacklogRow{}, backlog.Active...), backlog.Parked...))
	scan := normalization.ScanPlanFiles(workRoot, normalization.ScanSource{ReadFile: defaultWorkSource.ReadFile, ReadDir: defaultWorkSource.ReadDir, Stat: defaultWorkSource.Stat})
	closed := map[string][]string{}
	for _, file := range scan.Files {
		rel, err := filepath.Rel(workRoot, file.Path)
		if err != nil {
			return result, err
		}
		if strings.HasPrefix(filepath.ToSlash(rel), "closed/") {
			closed[file.ID] = append(closed[file.ID], file.Path)
		}
	}
	edits := map[string]string{}
	for _, file := range scan.Files {
		archives := closed[file.ID]
		if len(archives) == 0 || matchesAnyID(referenced, file.ID) {
			continue
		}
		rel, err := filepath.Rel(workRoot, file.Path)
		if err != nil {
			return result, err
		}
		rel = filepath.ToSlash(rel)
		if !strings.HasPrefix(rel, "active/") && !strings.HasPrefix(rel, "parked/") {
			continue
		}
		if len(archives) != 1 {
			result.ReviewFiles = append(result.ReviewFiles, repo.ToRepoRelative(repoRoot, file.Path))
			continue
		}
		archive := archives[0]
		regular := true
		for _, path := range []string{file.Path, archive} {
			info, err := defaultWorkSource.Lstat(path)
			if err != nil {
				return result, err
			}
			if !info.Mode().IsRegular() {
				regular = false
			}
		}
		if !regular {
			result.ReviewFiles = append(result.ReviewFiles, repo.ToRepoRelative(repoRoot, file.Path))
			continue
		}
		for _, pair := range [][2]string{{file.Path, archive}, {archive, file.Path}} {
			content, ok := edits[pair[0]]
			if !ok {
				raw, err := defaultWorkSource.ReadFile(pair[0])
				if err != nil {
					return result, err
				}
				content = string(raw)
			}
			target, err := filepath.Rel(filepath.Dir(pair[0]), pair[1])
			if err != nil {
				return result, err
			}
			link := fmt.Sprintf("- [Related work record](<%s>)", filepath.ToSlash(target))
			next := addRelatedRecordLink(content, link)
			if next != content {
				edits[pair[0]] = next
			}
		}
	}
	paths := make([]string, 0, len(edits))
	for path := range edits {
		paths = append(paths, path)
	}
	sort.Strings(paths)
	for _, path := range paths {
		if apply {
			info, err := defaultWorkSource.Stat(path)
			if err != nil {
				return result, err
			}
			if err := filesystem.AtomicWriteFile(repoRoot, path, []byte(edits[path]), info.Mode().Perm()); err != nil {
				return result, err
			}
		}
		result.ChangedFiles = append(result.ChangedFiles, repo.ToRepoRelative(repoRoot, path))
	}
	return result, nil
}

func addRelatedRecordLink(content, link string) string {
	if strings.Contains(content, link) {
		return content
	}
	const heading = "## Related Work Records\n"
	if start := strings.Index(content, "\n"+heading); start >= 0 {
		body := start + 1 + len(heading)
		end := len(content)
		if next := strings.Index(content[body:], "\n## "); next >= 0 {
			end = body + next
		}
		return strings.TrimRight(content[:end], "\n") + "\n" + link + "\n" + content[end:]
	}
	return ensureBlankLine(content) + heading + "\n" + link + "\n"
}

func matchesAnyID(ids map[string]struct{}, id string) bool {
	if _, ok := ids[id]; ok {
		return true
	}
	for known := range ids {
		if domainwork.IDsMatch(known, id) {
			return true
		}
	}
	return false
}

func ensureBlankLine(content string) string {
	switch {
	case strings.HasSuffix(content, "\n\n"):
		return content
	case strings.HasSuffix(content, "\n"):
		return content + "\n"
	default:
		return content + "\n\n"
	}
}
