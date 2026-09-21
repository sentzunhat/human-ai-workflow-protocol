package validation

import (
	"path/filepath"
	"strings"

	"github.com/sentzunhat/hawp/librarian/src/internal/domain/work/markdown"
)

// Only active work is scanned — archives (closed/, evidence/, notes/,
// status/) may reference old paths and are not expected to have live links.
var (
	activeDirs      = []string{"active", "parked"}
	activeRootFiles = []string{"BACKLOG.md"}
)

// CheckDeadLinks verifies local markdown links in BACKLOG.md and active/ and
// parked/ plan files, including canonical folder-per-item plans.
func CheckDeadLinks(workDir string, source Source) DeadLinksCheck {
	var files []string
	for _, name := range activeRootFiles {
		full := filepath.Join(workDir, name)
		if source.Exists(full) {
			files = append(files, full)
		}
	}
	for _, dir := range activeDirs {
		collectMarkdownFiles(filepath.Join(workDir, dir), source, &files)
	}

	result := DeadLinksCheck{Scanned: len(files), Status: StatusPass}
	for _, file := range files {
		raw, err := source.ReadFile(file)
		if err != nil {
			warnf(source, "skipping unreadable file %s: %v", file, err)
			continue
		}
		content := markdown.BlankFences(string(raw))
		rel := source.ToRepoRelative(workDir, file)

		for _, link := range markdown.ExtractLinks(content) {
			if !markdown.IsLocalHref(link.Href) {
				continue
			}
			part := markdown.PathPart(link.Href)
			if part == "" {
				continue
			}
			target := filepath.Join(filepath.Dir(file), part)
			if !source.Exists(target) {
				result.Broken = append(result.Broken, BrokenLink{ID: rel, Link: link.Href})
			}
		}
	}

	if len(result.Broken) > 0 {
		result.Status = StatusFail
	}
	return result
}

func collectMarkdownFiles(root string, source Source, files *[]string) {
	entries, err := source.ReadDir(root)
	if err != nil {
		return
	}
	for _, entry := range entries {
		full := filepath.Join(root, entry.Name())
		if entry.IsDir() {
			collectMarkdownFiles(full, source, files)
			continue
		}
		if strings.HasSuffix(entry.Name(), ".md") {
			*files = append(*files, full)
		}
	}
}
