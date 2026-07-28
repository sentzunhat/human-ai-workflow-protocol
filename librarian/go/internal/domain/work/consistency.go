package work

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	dateDirRe  = regexp.MustCompile(`^\d{2,4}$`)
	yearDirRe  = regexp.MustCompile(`^\d{4}$`)
	linkPathRe = regexp.MustCompile(`\[.*?\]\((.*?)\)`)
)

func warnf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, "[validate] warning: "+format+"\n", args...)
}

// CheckBacklogConsistency verifies that active/closed/parked rows resolve to
// real plan files and that no plan files are orphaned from the backlog.
func CheckBacklogConsistency(workDir string, backlog *Backlog) BacklogCheck {
	result := BacklogCheck{Status: StatusPass}

	result.ActiveWork.Total = len(backlog.Active)
	for _, row := range backlog.Active {
		if findActiveFile(workDir, row.ID) {
			result.ActiveWork.Found++
		} else {
			result.ActiveWork.Missing = append(result.ActiveWork.Missing, row.ID)
		}
	}

	result.RecentlyClosed.Total = len(backlog.Closed)
	closedDir := filepath.Join(workDir, "closed")
	for _, row := range backlog.Closed {
		if findClosedFile(closedDir, row.ID) {
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
		if _, err := os.Stat(filepath.Join(workDir, linkPath)); err == nil {
			result.ParkedWork.Found++
		} else {
			result.ParkedWork.Missing = append(result.ParkedWork.Missing, row.ID)
		}
	}

	activeIDs := idSet(backlog.Active)
	collectOrphanedActive(filepath.Join(workDir, "active"), "active", activeIDs, &result.OrphanedFiles)

	parkedDir := filepath.Join(workDir, "parked")
	if entries, err := os.ReadDir(parkedDir); err == nil {
		parkedIDs := idSet(backlog.Parked)
		for _, entry := range entries {
			if !strings.HasSuffix(entry.Name(), ".md") {
				continue
			}
			id := ExtractIDFromFilename(strings.TrimSuffix(entry.Name(), ".md"))
			if id != "" && !matchesAnyID(parkedIDs, id) {
				result.OrphanedParked = append(result.OrphanedParked, "parked/"+entry.Name())
			}
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

// findActiveFile supports flat active/<ID>.md, flat short-UUID prefix
// matches, and date-nested active/YYYY/MM/DD/<ID>.md layouts.
func findActiveFile(workDir, id string) bool {
	activeDir := filepath.Join(workDir, "active")
	if _, err := os.Stat(filepath.Join(activeDir, id+".md")); err == nil {
		return true
	}
	entries, err := os.ReadDir(activeDir)
	if err != nil {
		return false
	}
	for _, entry := range entries {
		if !strings.HasSuffix(entry.Name(), ".md") {
			continue
		}
		fileID := ExtractIDFromFilename(strings.TrimSuffix(entry.Name(), ".md"))
		if fileID != "" && IDsMatch(id, fileID) {
			return true
		}
	}
	for _, entry := range entries {
		if !entry.IsDir() || !yearDirRe.MatchString(entry.Name()) {
			continue
		}
		if findInDateTree(filepath.Join(activeDir, entry.Name()), id) {
			return true
		}
	}
	return false
}

func findInDateTree(yearPath, id string) bool {
	months, err := os.ReadDir(yearPath)
	if err != nil {
		return false
	}
	for _, month := range months {
		if !month.IsDir() {
			continue
		}
		days, err := os.ReadDir(filepath.Join(yearPath, month.Name()))
		if err != nil {
			continue
		}
		for _, day := range days {
			if !day.IsDir() {
				continue
			}
			files, err := os.ReadDir(filepath.Join(yearPath, month.Name(), day.Name()))
			if err != nil {
				continue
			}
			for _, file := range files {
				name := file.Name()
				if name == id+".md" {
					return true
				}
				if strings.HasSuffix(name, ".md") &&
					strings.Contains(strings.ToLower(name), strings.ToLower(id)) {
					return true
				}
			}
		}
	}
	return false
}

// findClosedFile searches closed/YYYY/MM/DD/ for an exact, containing, or
// short-UUID-prefix match.
func findClosedFile(closedDir, id string) bool {
	years, err := os.ReadDir(closedDir)
	if err != nil {
		return false
	}
	for _, year := range years {
		if year.Name() == "README.md" || !year.IsDir() {
			continue
		}
		months, err := os.ReadDir(filepath.Join(closedDir, year.Name()))
		if err != nil {
			warnf("error while searching closed plans for %s: %v", id, err)
			continue
		}
		for _, month := range months {
			if !month.IsDir() {
				continue
			}
			days, err := os.ReadDir(filepath.Join(closedDir, year.Name(), month.Name()))
			if err != nil {
				continue
			}
			for _, day := range days {
				if !day.IsDir() {
					continue
				}
				files, err := os.ReadDir(filepath.Join(closedDir, year.Name(), month.Name(), day.Name()))
				if err != nil {
					continue
				}
				for _, file := range files {
					name := file.Name()
					if !strings.HasSuffix(name, ".md") {
						continue
					}
					if name == id+".md" ||
						strings.Contains(strings.ToLower(name), strings.ToLower(id)) {
						return true
					}
					fileID := ExtractIDFromFilename(strings.TrimSuffix(name, ".md"))
					if fileID != "" && IDsMatch(id, fileID) {
						return true
					}
				}
			}
		}
	}
	return false
}

// collectOrphanedActive flags active plan files whose ID is not in the
// backlog, recursing only into date-shaped subdirectories.
func collectOrphanedActive(dir, relPrefix string, activeIDs map[string]struct{}, out *[]string) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return
	}
	for _, entry := range entries {
		if entry.IsDir() {
			if dateDirRe.MatchString(entry.Name()) {
				collectOrphanedActive(filepath.Join(dir, entry.Name()), relPrefix+"/"+entry.Name(), activeIDs, out)
			}
			continue
		}
		if !strings.HasSuffix(entry.Name(), ".md") {
			continue
		}
		id := ExtractIDFromFilename(strings.TrimSuffix(entry.Name(), ".md"))
		if id != "" && !matchesAnyID(activeIDs, id) {
			*out = append(*out, relPrefix+"/"+entry.Name())
		}
	}
}

func extractLinkPath(detail string) string {
	if m := linkPathRe.FindStringSubmatch(detail); m != nil {
		return m[1]
	}
	return ""
}
