package work

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/sentzunhat/hawp/librarian/go/internal/infrastructure/markdown"
	"github.com/sentzunhat/hawp/librarian/go/internal/infrastructure/repo"
)

// Only active work is scanned — archives (closed/, evidence/, notes/,
// status/) may reference old paths and are not expected to have live links.
var (
	activeDirs      = []string{"active", "parked"}
	activeRootFiles = []string{"BACKLOG.md"}
)

// CheckDeadLinks verifies local markdown links in BACKLOG.md and the flat
// active/ and parked/ plan files.
func CheckDeadLinks(workDir string) DeadLinksCheck {
	var files []string
	for _, name := range activeRootFiles {
		full := filepath.Join(workDir, name)
		if repo.Exists(full) {
			files = append(files, full)
		}
	}
	for _, dir := range activeDirs {
		entries, err := os.ReadDir(filepath.Join(workDir, dir))
		if err != nil {
			continue
		}
		for _, entry := range entries {
			if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".md") {
				files = append(files, filepath.Join(workDir, dir, entry.Name()))
			}
		}
	}

	result := DeadLinksCheck{Scanned: len(files), Status: StatusPass}
	for _, file := range files {
		raw, err := os.ReadFile(file)
		if err != nil {
			warnf("skipping unreadable file %s: %v", file, err)
			continue
		}
		content := markdown.BlankFences(string(raw))
		rel := repo.ToRepoRelative(workDir, file)

		for _, link := range markdown.ExtractLinks(content) {
			if !markdown.IsLocalHref(link.Href) {
				continue
			}
			pathPart := markdown.PathPart(link.Href)
			if pathPart == "" {
				continue
			}
			target := filepath.Join(filepath.Dir(file), pathPart)
			if !repo.Exists(target) {
				result.Broken = append(result.Broken, BrokenLink{ID: rel, Link: link.Href})
			}
		}
	}

	if len(result.Broken) > 0 {
		result.Status = StatusFail
	}
	return result
}
