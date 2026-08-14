package work

import (
	"path/filepath"
	"strings"

	"github.com/sentzunhat/hawp/librarian/go/internal/domain/work/source"
)

// Only active work is scanned — archives (closed/, evidence/, notes/,
// status/) may reference old paths and are not expected to have live links.
var (
	activeDirs      = []string{"active", "parked"}
	activeRootFiles = []string{"BACKLOG.md"}
)

// CheckDeadLinks verifies local markdown links in BACKLOG.md and the flat
// active/ and parked/ plan files from an acquired work snapshot.
func CheckDeadLinks(snapshot source.Snapshot) DeadLinksCheck {
	var files []source.File
	for _, file := range snapshot.Files {
		if file.RelPath == activeRootFiles[0] || isFlatActiveFile(file.RelPath) {
			files = append(files, file)
		}
	}

	result := DeadLinksCheck{Scanned: len(files), Status: StatusPass}
	for _, file := range files {
		for _, link := range file.Links {
			if !isLocalHref(link.Href) {
				continue
			}
			pathPart := pathPart(link.Href)
			if pathPart == "" {
				continue
			}
			target := filepath.Clean(filepath.Join(filepath.Dir(file.Path), pathPart))
			if _, ok := snapshot.ExistingPaths[target]; !ok {
				result.Broken = append(result.Broken, BrokenLink{ID: file.RelPath, Link: link.Href})
			}
		}
	}

	if len(result.Broken) > 0 {
		result.Status = StatusFail
	}
	return result
}

func isFlatActiveFile(relPath string) bool {
	for _, dir := range activeDirs {
		if strings.HasPrefix(relPath, dir+"/") && !strings.Contains(strings.TrimPrefix(relPath, dir+"/"), "/") && strings.HasSuffix(relPath, ".md") {
			return true
		}
	}
	return false
}

func isLocalHref(href string) bool {
	return href != "" && !strings.HasPrefix(href, "http") && !strings.HasPrefix(href, "/") && !strings.HasPrefix(href, "#")
}

func pathPart(href string) string {
	part, _, _ := strings.Cut(href, "#")
	return part
}
