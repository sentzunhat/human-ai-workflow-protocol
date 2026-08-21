package work

import (
	"regexp"
	"strings"

	"github.com/sentzunhat/hawp/librarian/go/internal/domain/work/source"
)

var linkPathRe = regexp.MustCompile(`\[.*?\]\((.*?)\)`)

// CheckBacklogConsistency verifies that active/closed/parked rows resolve to
// real plan files and that no plan files are orphaned from the backlog.
func CheckBacklogConsistency(snapshot source.Snapshot, backlog *Backlog) BacklogCheck {
	result := BacklogCheck{Status: StatusPass}

	result.ActiveWork.Total = len(backlog.Active)
	for _, row := range backlog.Active {
		if findActiveFileInFiles(snapshot.Files, row.ID) {
			result.ActiveWork.Found++
		} else {
			result.ActiveWork.Missing = append(result.ActiveWork.Missing, row.ID)
		}
	}

	result.RecentlyClosed.Total = len(backlog.Closed)
	for _, row := range backlog.Closed {
		if findClosedFileInFiles(snapshot.Files, row.ID) {
			result.RecentlyClosed.Found++
		} else {
			result.RecentlyClosed.Missing = append(result.RecentlyClosed.Missing, row.ID)
		}
	}

	result.ParkedWork.Total = len(backlog.Parked)
	for _, row := range backlog.Parked {
		linkPath := extractLinkPath(row.Detail)
		if linkPath == "" {
			result.ParkedWork.Missing = append(result.ParkedWork.Missing, row.ID)
			continue
		}
		if hasFileWithRelPath(snapshot.Files, linkPath) {
			result.ParkedWork.Found++
		} else {
			result.ParkedWork.Missing = append(result.ParkedWork.Missing, row.ID)
		}
	}

	activeIDs := idSet(backlog.Active)
	result.OrphanedFiles = collectOrphanedActiveFromFiles(snapshot.Files, activeIDs)

	parkedIDs := idSet(backlog.Parked)
	for _, f := range snapshot.Files {
		relPath := strings.ReplaceAll(f.RelPath, "\\", "/")
		if !strings.HasPrefix(relPath, "parked/") || !strings.HasSuffix(relPath, ".md") {
			continue
		}
		afterPrefix := relPath[len("parked/"):]
		if strings.Contains(afterPrefix, "/") {
			continue // not a flat parked file
		}
		name := strings.TrimSuffix(filenameFromPath(relPath), ".md")
		id := ExtractIDFromFilename(name)
		if id != "" && !matchesAnyID(parkedIDs, id) {
			result.OrphanedParked = append(result.OrphanedParked, relPath)
		}
	}

	if len(result.ActiveWork.Missing) > 0 || len(result.RecentlyClosed.Missing) > 0 ||
		len(result.ParkedWork.Missing) > 0 || len(result.OrphanedFiles) > 0 ||
		len(result.OrphanedParked) > 0 {
		result.Status = StatusFail
	}
	return result
}

func idSet(rows []BacklogRow) map[string]struct{} {
	ids := make(map[string]struct{}, len(rows))
	for _, row := range rows {
		ids[row.ID] = struct{}{}
	}
	return ids
}

// findActiveFileInFiles reports whether any file in the snapshot has a RelPath
// under active/ that matches id by exact name, ID extraction, or short-UUID prefix.
func findActiveFileInFiles(files []source.File, id string) bool {
	for _, f := range files {
		relPath := strings.ReplaceAll(f.RelPath, "\\", "/")
		if !strings.HasPrefix(relPath, "active/") || !strings.HasSuffix(relPath, ".md") {
			continue
		}
		name := strings.TrimSuffix(filenameFromPath(relPath), ".md")
		if name == id {
			return true
		}
		fileID := ExtractIDFromFilename(name)
		if fileID != "" && IDsMatch(id, fileID) {
			return true
		}
	}
	return false
}

// findClosedFileInFiles reports whether any file in the snapshot has a RelPath
// under closed/ that matches id by exact name, substring, or ID extraction.
func findClosedFileInFiles(files []source.File, id string) bool {
	for _, f := range files {
		relPath := strings.ReplaceAll(f.RelPath, "\\", "/")
		if !strings.HasPrefix(relPath, "closed/") || !strings.HasSuffix(relPath, ".md") {
			continue
		}
		name := strings.TrimSuffix(filenameFromPath(relPath), ".md")
		if name == id || strings.Contains(strings.ToLower(name), strings.ToLower(id)) {
			return true
		}
		fileID := ExtractIDFromFilename(name)
		if fileID != "" && IDsMatch(id, fileID) {
			return true
		}
	}
	return false
}

// hasFileWithRelPath reports whether any file in the snapshot has the given
// relative path (normalised to forward slashes).
func hasFileWithRelPath(files []source.File, relPath string) bool {
	target := strings.ReplaceAll(relPath, "\\", "/")
	for _, f := range files {
		if strings.ReplaceAll(f.RelPath, "\\", "/") == target {
			return true
		}
	}
	return false
}

// collectOrphanedActiveFromFiles returns relative paths of active plan files
// whose IDs are not present in the backlog.
func collectOrphanedActiveFromFiles(files []source.File, activeIDs map[string]struct{}) []string {
	var orphans []string
	for _, f := range files {
		relPath := strings.ReplaceAll(f.RelPath, "\\", "/")
		if !strings.HasPrefix(relPath, "active/") || !strings.HasSuffix(relPath, ".md") {
			continue
		}
		name := strings.TrimSuffix(filenameFromPath(relPath), ".md")
		id := ExtractIDFromFilename(name)
		if id != "" && !matchesAnyID(activeIDs, id) {
			orphans = append(orphans, relPath)
		}
	}
	return orphans
}

func extractLinkPath(detail string) string {
	if m := linkPathRe.FindStringSubmatch(detail); m != nil {
		return m[1]
	}
	return ""
}
