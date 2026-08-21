package work

import (
	"regexp"
	"strings"

	"github.com/sentzunhat/hawp/librarian/go/internal/domain/work/source"
)

const legacyClosedCutoff = "2026-05-10"

var supportingSuffixes = []string{
	"-summary", "-status", "-status-report", "-checkpoint", "-evidence",
}

var (
	outcomeHeadingRe        = regexp.MustCompile(`(?i)^#{1,6}\s*Outcome\b`)
	verificationHeadingRe   = regexp.MustCompile(`(?i)^#{1,6}\s*Verification\b`)
	closeChecklistHeadingRe = regexp.MustCompile(`(?i)^#{1,6}\s*Close Checklist\b`)
	pathDateRe              = regexp.MustCompile(`/(\d{4})/(\d{2})/(\d{2})/`)
)

type classification struct {
	kind   string // "plan" | "supporting" | "legacy-untyped" | "current-untyped"
	id     string
	reason string
}

// CheckClosedTaskCompleteness verifies that closed plan files in snapshot
// carry the Outcome, Verification, and Close Checklist sections. Files dated
// before the legacy cutoff only warn; supporting files are skipped.
func CheckClosedTaskCompleteness(snapshot source.Snapshot) ClosedTaskCheck {
	result := ClosedTaskCheck{Status: StatusPass}

	for _, file := range snapshot.Files {
		if !isClosedMarkdownFile(file.RelPath) {
			continue
		}

		filename := filenameFromPath(file.RelPath)
		dateFromPath := dateFromClosedPath(file.RelPath)
		class := classifyClosedFile(filename, dateFromPath)
		date := dateFromPath
		if date == "" {
			date = "unknown"
		}
		finding := FileFinding{ID: class.id, Date: date, Reason: class.reason, FilePath: file.Path}

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
		hasOutcome, hasVerification, hasChecklist := findRequiredHeadings(file.Content)
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
			if dateFromPath == "" || dateFromPath < legacyClosedCutoff {
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

func classifyClosedFile(filename, date string) classification {
	nameWithoutExt := strings.TrimSuffix(filename, ".md")
	nameLower := strings.ToLower(nameWithoutExt)
	id := ExtractIDFromFilename(nameWithoutExt)
	isLegacy := date == "" || date < legacyClosedCutoff

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

	if isLegacy {
		return classification{"legacy-untyped", nameWithoutExt, "legacy file without TASK-/BUG-style ID"}
	}
	return classification{"current-untyped", nameWithoutExt, "current file without TASK-/BUG-style ID"}
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

func isClosedMarkdownFile(relPath string) bool {
	relPath = strings.ReplaceAll(relPath, "\\", "/")
	return strings.HasPrefix(relPath, "closed/") &&
		strings.HasSuffix(relPath, ".md") &&
		filenameFromPath(relPath) != "README.md"
}

func filenameFromPath(relPath string) string {
	relPath = strings.ReplaceAll(relPath, "\\", "/")
	if index := strings.LastIndex(relPath, "/"); index >= 0 {
		return relPath[index+1:]
	}
	return relPath
}

func dateFromClosedPath(relPath string) string {
	if m := pathDateRe.FindStringSubmatch("/" + strings.ReplaceAll(relPath, "\\", "/")); m != nil {
		return m[1] + "-" + m[2] + "-" + m[3]
	}
	return ""
}
