package work

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/sentzunhat/hawp/librarian/src/internal/infrastructure/repo"
)

var supportingSuffixes = []string{
	"-summary", "-status", "-status-report", "-checkpoint", "-evidence",
}

var (
	outcomeHeadingRe        = regexp.MustCompile(`(?i)^#{1,6}\s*Outcome\b`)
	verificationHeadingRe   = regexp.MustCompile(`(?i)^#{1,6}\s*Verification\b`)
	closeChecklistHeadingRe = regexp.MustCompile(`(?i)^#{1,6}\s*Close Checklist\b`)
	pathDateRe              = regexp.MustCompile(`/(\d{4})/(\d{2})/(\d{2})/`)
)

type closedFileEntry struct {
	filePath string
	date     string // "YYYY-MM-DD" from the path, or ""
}

type classification struct {
	kind   string // "plan" | "supporting" | "legacy-untyped" | "current-untyped"
	id     string
	reason string
}

// CheckClosedTaskCompleteness verifies that closed plan files carry the
// Outcome, Verification, and Close Checklist sections. Files dated before the
// legacy cutoff only warn; supporting files (summaries, evidence, archives)
// are skipped.
func CheckClosedTaskCompleteness(workDir string) ClosedTaskCheck {
	result := ClosedTaskCheck{Status: StatusPass}

	for _, entry := range collectClosedFilesGeneric(filepath.Join(workDir, "closed")) {
		filename := filepath.Base(entry.filePath)
		// Folder-per-item layout: plan.md in a uuid-named dir — use parent dir as ID
		effectiveName := filename
		if filename == "plan.md" {
			parts := strings.Split(strings.ReplaceAll(entry.filePath, "\\", "/"), "/")
			if len(parts) >= 2 {
				effectiveName = parts[len(parts)-2]
			}
		}

		content, err := os.ReadFile(entry.filePath)
		if err != nil {
			warnf("skipping unreadable closed plan %s: %v", entry.filePath, err)
			continue
		}

		class := classifyClosedFile(effectiveName, entry.date, string(content))
		date := entry.date
		if date == "" {
			date = "unknown"
		}
		finding := FileFinding{ID: class.id, Date: date, Reason: class.reason, FilePath: entry.filePath}

		switch class.kind {
		case "supporting":
			result.Skipped++
			result.SupportingSkipped = append(result.SupportingSkipped, finding)
			continue
		case "legacy-untyped":
			result.UntypedLegacy = append(result.UntypedLegacy, finding)
			continue
		case "current-untyped":
			result.UntypedCurrent = append(result.UntypedCurrent, finding)
			continue
		}

		result.Total++

		var missing []string
		hasOutcome, hasVerification, hasChecklist := findRequiredHeadings(string(content))
		if hasOutcome {
			result.WithOutcome++
		} else {
			missing = append(missing, "Outcome")
		}
		if hasVerification {
			result.WithVerification++
		} else {
			missing = append(missing, "Verification")
		}
		if hasChecklist {
			result.WithCloseChecklist++
		} else {
			missing = append(missing, "Close Checklist")
		}

		if len(missing) > 0 {
			finding.Sections = missing
			if entry.date == "" || entry.date < repo.LegacyClosedCutoff {
				result.Warnings = append(result.Warnings, finding)
			} else {
				result.Failing = append(result.Failing, finding)
			}
		}
	}

	if len(result.Failing) > 0 || len(result.UntypedCurrent) > 0 {
		result.Status = StatusFail
	} else if len(result.Warnings) > 0 || len(result.UntypedLegacy) > 0 {
		result.Status = StatusWarn
	}
	return result
}

func classifyClosedFile(filename, date, content string) classification {
	nameWithoutExt := strings.TrimSuffix(filename, ".md")
	nameLower := strings.ToLower(nameWithoutExt)
	id := ExtractIDFromFilename(nameWithoutExt)
	isLegacy := date == "" || date < repo.LegacyClosedCutoff

	// Shared batch-close record: **Closes:** line names every backlog row it covers
	if strings.Contains(content, "**Closes:**") {
		return classification{"plan", nameWithoutExt, ""}
	}

	if strings.HasPrefix(nameLower, "backlog") {
		return classification{"supporting", nameWithoutExt, "matches BACKLOG supporting-file pattern"}
	}
	if strings.Contains(nameLower, "archive") {
		return classification{"supporting", nameWithoutExt, "matches archive supporting-file pattern"}
	}

	if id != "" {
		suffix := ""
		if pos := strings.Index(nameLower, strings.ToLower(id)); pos >= 0 {
			suffix = nameLower[pos+len(id):]
		}
		for _, kw := range supportingSuffixes {
			if suffix == kw || strings.HasSuffix(suffix, kw) {
				reason := suffix
				if reason == "" {
					reason = "none"
				}
				return classification{"supporting", nameWithoutExt, "matches supporting suffix pattern (" + reason + ")"}
			}
		}
		if strings.Contains(suffix, "-archive") {
			return classification{"supporting", nameWithoutExt, "matches archive supporting-file pattern"}
		}
		return classification{"plan", id, ""}
	}

	if isSlugID(nameLower) && hasPlanMetadata(content) {
		return classification{"plan", nameLower, ""}
	}

	if isLegacy {
		return classification{"legacy-untyped", nameWithoutExt, "legacy file without TASK-/BUG-style ID"}
	}
	return classification{"current-untyped", nameWithoutExt, "current file without TASK-/BUG-style ID"}
}

func isSlugID(value string) bool {
	matched, _ := regexp.MatchString(`^[a-z0-9]+(?:-[a-z0-9]+)+$`, value)
	return matched
}

func hasPlanMetadata(content string) bool {
	return strings.Contains(content, "**Type:**") && strings.Contains(content, "**Status:**")
}

func findRequiredHeadings(content string) (outcome, verification, checklist bool) {
	for _, line := range strings.Split(content, "\n") {
		trimmed := strings.TrimSpace(line)
		if !outcome && outcomeHeadingRe.MatchString(trimmed) {
			outcome = true
			continue
		}
		if !verification && verificationHeadingRe.MatchString(trimmed) {
			verification = true
			continue
		}
		if !checklist && closeChecklistHeadingRe.MatchString(trimmed) {
			checklist = true
		}
	}
	return
}

// collectClosedFilesGeneric walks closed/ at any depth, extracting the date
// from a /YYYY/MM/DD/ path segment when present. README.md files are skipped.
func collectClosedFilesGeneric(closedDir string) []closedFileEntry {
	var out []closedFileEntry
	var walk func(dir string)
	walk = func(dir string) {
		entries, err := os.ReadDir(dir)
		if err != nil {
			return
		}
		for _, entry := range entries {
			full := filepath.Join(dir, entry.Name())
			if entry.IsDir() {
				walk(full)
				continue
			}
			if !strings.HasSuffix(entry.Name(), ".md") || entry.Name() == "README.md" {
				continue
			}
			date := ""
			if m := pathDateRe.FindStringSubmatch(strings.ReplaceAll(full, "\\", "/")); m != nil {
				date = m[1] + "-" + m[2] + "-" + m[3]
			}
			out = append(out, closedFileEntry{filePath: full, date: date})
		}
	}
	walk(closedDir)
	return out
}
