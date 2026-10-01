package validation

import (
	"path/filepath"
	"regexp"
	"strings"

	"github.com/sentzunhat/hawp/librarian/src/internal/domain/work/identity"
)

var (
	dateDirRe      = regexp.MustCompile(`^\d{2,4}$`)
	yearDirRe      = regexp.MustCompile(`^\d{4}$`)
	linkPathRe     = regexp.MustCompile(`\[.*?\]\((.*?)\)`)
	canonicalIDRe  = regexp.MustCompile(`^(TASK|BUG)-\d+$`)
	numericRowIDRe = regexp.MustCompile(`^\d+$`)
)

// CheckBacklogConsistency verifies that active/closed/parked rows resolve to
// real plan files and that no plan files are orphaned from the backlog.
func CheckBacklogConsistency(workDir string, backlog *Backlog, source Source) BacklogCheck {
	result := BacklogCheck{Status: StatusPass}

	result.ActiveWork.Total = len(backlog.Active)
	for _, row := range backlog.Active {
		if isNonCanonicalActiveID(row.ID) {
			result.NonCanonicalActiveItems = append(result.NonCanonicalActiveItems, row.ID)
		}
		if findActiveFile(workDir, row.ID, source) {
			result.ActiveWork.Found++
		} else {
			result.ActiveWork.Missing = append(result.ActiveWork.Missing, row.ID)
		}
	}

	result.RecentlyClosed.Total = len(backlog.Closed)
	closedDir := filepath.Join(workDir, "closed")
	for _, row := range backlog.Closed {
		if findClosedFile(closedDir, row.ID, source) {
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
		if _, err := source.Stat(filepath.Join(workDir, linkPath)); err == nil {
			result.ParkedWork.Found++
		} else {
			result.ParkedWork.Missing = append(result.ParkedWork.Missing, row.ID)
		}
	}

	activeIDs := idSet(backlog.Active)
	collectOrphanedActive(filepath.Join(workDir, "active"), "active", activeIDs, &result.OrphanedFiles, source)

	parkedDir := filepath.Join(workDir, "parked")
	if entries, err := source.ReadDir(parkedDir); err == nil {
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
		len(result.OrphanedParked) > 0 || len(result.NonCanonicalActiveItems) > 0 {
		result.Status = StatusFail
	}
	return result
}

func isNonCanonicalActiveID(id string) bool {
	if ExtractShortUUID(id) != "" || identity.IsFullUUID(id) {
		return false
	}
	if canonicalIDRe.MatchString(id) || identity.IsNumericID(id) {
		return false
	}
	return true
}

// IDSet returns a set of row IDs from the given backlog rows.
func IDSet(rows []BacklogRow) map[string]struct{} {
	ids := make(map[string]struct{}, len(rows))
	for _, row := range rows {
		ids[row.ID] = struct{}{}
	}
	return ids
}

func idSet(rows []BacklogRow) map[string]struct{} { return IDSet(rows) }

// findActiveFile supports flat active/<ID>.md, folder-per-item
// active/<ID>/plan.md, flat short-UUID prefix matches, and date-nested
// active/YYYY/MM/DD/<ID>.md layouts.
func findActiveFile(workDir, id string, source Source) bool {
	activeDir := filepath.Join(workDir, "active")
	if _, err := source.Stat(filepath.Join(activeDir, id+".md")); err == nil {
		return true
	}
	// Folder-per-item: active/{id}/plan.md
	if _, err := source.Stat(filepath.Join(activeDir, id, "plan.md")); err == nil {
		return true
	}
	entries, err := source.ReadDir(activeDir)
	if err != nil {
		return false
	}
	for _, entry := range entries {
		if entry.IsDir() {
			if yearDirRe.MatchString(entry.Name()) {
				if findInDateTree(filepath.Join(activeDir, entry.Name()), id, source) {
					return true
				}
				continue
			}
			// Non-date directory: folder-per-item layout
			dirID := ExtractIDFromFilename(entry.Name())
			if dirID != "" && IDsMatch(id, dirID) {
				if _, err := source.Stat(filepath.Join(activeDir, entry.Name(), "plan.md")); err == nil {
					return true
				}
			}
			continue
		}
		if !strings.HasSuffix(entry.Name(), ".md") {
			continue
		}
		fileID := ExtractIDFromFilename(strings.TrimSuffix(entry.Name(), ".md"))
		if fileID != "" && IDsMatch(id, fileID) {
			return true
		}
	}
	return false
}

func findInDateTree(yearPath, id string, source Source) bool {
	months, err := source.ReadDir(yearPath)
	if err != nil {
		return false
	}
	for _, month := range months {
		if !month.IsDir() {
			continue
		}
		days, err := source.ReadDir(filepath.Join(yearPath, month.Name()))
		if err != nil {
			continue
		}
		for _, day := range days {
			if !day.IsDir() {
				continue
			}
			files, err := source.ReadDir(filepath.Join(yearPath, month.Name(), day.Name()))
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

// recordListsID returns true when the file at path contains the id wrapped in
// backticks (the convention for shared batch-close records like v003-ship-audit).
func recordListsID(path, id string, source Source) bool {
	content, err := source.ReadFile(path)
	if err != nil {
		return false
	}
	needle := "`" + id + "`"
	return strings.Contains(strings.ToLower(string(content)), strings.ToLower(needle))
}

// findClosedFile searches closed/YYYY/MM/DD/ for an exact, containing, or
// short-UUID-prefix match. Also supports folder-per-item layout
// (closed/YYYY/MM/DD/{id}/plan.md) and shared batch-close records whose
// plan.md content lists the id in backticks.
func findClosedFile(closedDir, id string, source Source) bool {
	years, err := source.ReadDir(closedDir)
	if err != nil {
		return false
	}
	for _, year := range years {
		if year.Name() == "README.md" || !year.IsDir() {
			continue
		}
		months, err := source.ReadDir(filepath.Join(closedDir, year.Name()))
		if err != nil {
			warnf(source, "error while searching closed plans for %s: %v", id, err)
			continue
		}
		for _, month := range months {
			if !month.IsDir() {
				continue
			}
			days, err := source.ReadDir(filepath.Join(closedDir, year.Name(), month.Name()))
			if err != nil {
				continue
			}
			for _, day := range days {
				if !day.IsDir() {
					continue
				}
				dayPath := filepath.Join(closedDir, year.Name(), month.Name(), day.Name())
				entries, err := source.ReadDir(dayPath)
				if err != nil {
					continue
				}
				for _, entry := range entries {
					entryPath := filepath.Join(dayPath, entry.Name())
					if entry.IsDir() {
						// Folder-per-item: closed/YYYY/MM/DD/{id}/plan.md
						if strings.EqualFold(entry.Name(), id) {
							if _, err := source.Stat(filepath.Join(entryPath, "plan.md")); err == nil {
								return true
							}
						}
						dirID := ExtractIDFromFilename(entry.Name())
						if dirID != "" && IDsMatch(id, dirID) {
							if _, err := source.Stat(filepath.Join(entryPath, "plan.md")); err == nil {
								return true
							}
						}
						// Shared batch-close: plan.md content may list the id
						planPath := filepath.Join(entryPath, "plan.md")
						if _, err := source.Stat(planPath); err == nil && recordListsID(planPath, id, source) {
							return true
						}
						continue
					}
					name := entry.Name()
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
// backlog. Recurses into date-shaped subdirectories; treats non-date
// directories as folder-per-item work items and flags them if orphaned.
func collectOrphanedActive(dir, relPrefix string, activeIDs map[string]struct{}, out *[]string, source Source) {
	entries, err := source.ReadDir(dir)
	if err != nil {
		return
	}
	for _, entry := range entries {
		if entry.IsDir() {
			if dateDirRe.MatchString(entry.Name()) {
				collectOrphanedActive(filepath.Join(dir, entry.Name()), relPrefix+"/"+entry.Name(), activeIDs, out, source)
			} else {
				// Non-date dir: folder-per-item work item
				dirID := ExtractIDFromFilename(entry.Name())
				if dirID != "" && !matchesAnyID(activeIDs, dirID) {
					*out = append(*out, relPrefix+"/"+entry.Name()+"/plan.md")
				}
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
