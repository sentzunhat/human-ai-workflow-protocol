package work

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/sentzunhat/hawp/librarian/go/internal/domain/work/source"
)

var (
	legacyIDInPathRe   = regexp.MustCompile(`(?i)(TASK|BUG)-\d+`)
	closedDateInPathRe = regexp.MustCompile(`/closed/(\d{4})/(\d{2})/(\d{2})/`)
	fileDatePrefixRe   = regexp.MustCompile(`^(\d{4})-(\d{2})-(\d{2})-`)
	backlogIDLineRe    = regexp.MustCompile(`(?i)\*\*Backlog ID:\*\*`)
	h1Re               = regexp.MustCompile(`^#\s+`)
	multiBlankRe       = regexp.MustCompile(`\n{3,}`)
	researchEntryRe    = regexp.MustCompile(`(?m)^\s*-\s+\[ \]\s+Research evidence for: (.+)$`)
	verifyHeadBodyRe   = regexp.MustCompile(`(?ms)(^##\s+Verification\b[^\n]*\n)(.*?)(?:^##\s+|\z)`)
)

// ApplyResult summarizes an apply-mode normalization run.
type ApplyResult struct {
	ChangedFiles  []string
	SkippedFiles  []string
	ResearchQueue []ResearchItem
}

func inferBacklogIDFromPath(path string) string {
	return strings.ToUpper(legacyIDInPathRe.FindString(path))
}

func inferClosedDateFromPath(path string) string {
	if m := closedDateInPathRe.FindStringSubmatch(strings.ReplaceAll(path, "\\", "/")); m != nil {
		return m[1] + "-" + m[2] + "-" + m[3]
	}
	return ""
}

func inferFileDatePrefix(path string) string {
	if m := fileDatePrefixRe.FindStringSubmatch(filepath.Base(path)); m != nil {
		return m[1] + "-" + m[2] + "-" + m[3]
	}
	return ""
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

func appendSection(content, heading, body string) string {
	if hasHeadingNamed(content, heading) {
		return content
	}
	return ensureBlankLine(content) + "## " + heading + "\n\n" + body + "\n"
}

func insertBacklogID(content, backlogID string) string {
	if backlogIDLineRe.MatchString(content) {
		return content
	}
	lines := strings.Split(content, "\n")
	for i, line := range lines {
		if h1Re.MatchString(line) {
			inserted := append(append(append([]string{}, lines[:i+1]...),
				"", "**Backlog ID:** "+backlogID, ""), lines[i+1:]...)
			return multiBlankRe.ReplaceAllString(strings.Join(inserted, "\n"), "\n\n")
		}
	}
	return "**Backlog ID:** " + backlogID + "\n\n" + content
}

// ensureEvidenceFollowUp appends a "### Evidence Follow-Up" subsection with
// research entries for each ambiguous claim not already queued.
func ensureEvidenceFollowUp(content string) (string, []string) {
	ambiguous := AmbiguousVerificationClaims(content)
	if len(ambiguous) == 0 {
		return content, nil
	}
	loc := verifyHeadBodyRe.FindStringSubmatchIndex(content)
	if loc == nil {
		return content, nil
	}
	heading := content[loc[2]:loc[3]]
	body := content[loc[4]:loc[5]]

	existing := map[string]struct{}{}
	for _, m := range researchEntryRe.FindAllStringSubmatch(body, -1) {
		existing[strings.TrimSpace(m[1])] = struct{}{}
	}

	var pending []string
	for _, claim := range ambiguous {
		claim = strings.TrimSpace(claim)
		if claim == "" {
			continue
		}
		if _, ok := existing[claim]; !ok {
			pending = append(pending, claim)
		}
	}
	if len(pending) == 0 {
		return content, nil
	}

	var entries []string
	for _, claim := range pending {
		entries = append(entries,
			"- [ ] Research evidence for: "+claim+"\n- [ ] Update the original verification checklist line with Evidence: ... or explicit unproven wording.")
	}
	newEntries := strings.Join(entries, "\n")

	const subsection = "### Evidence Follow-Up"
	trimmedBody := strings.TrimRight(body, " \t\n")
	var nextBody string
	if strings.Contains(body, subsection) {
		nextBody = trimmedBody + "\n" + newEntries + "\n"
	} else {
		nextBody = trimmedBody + "\n\n" + subsection + "\n\n" + newEntries + "\n"
	}

	sectionStart, bodyEnd := loc[2], loc[5]
	return content[:sectionStart] + heading + nextBody + content[bodyEnd:], pending
}

// normalizeClosedRecord scaffolds missing sections and evidence follow-ups.
func normalizeClosedRecord(content, filePath string) (string, []string) {
	updated := content
	if id := inferBacklogIDFromPath(filePath); id != "" {
		updated = insertBacklogID(updated, id)
	}
	updated = appendSection(updated, "Outcome", "_Legacy normalization scaffold added._")
	updated = appendSection(updated, "Verification", "_Legacy normalization scaffold added._")
	updated = appendSection(updated, "Close Checklist", "- [ ] Legacy normalization scaffold added.")
	return ensureEvidenceFollowUp(updated)
}

// reconcileClosedRecordPath moves date-prefixed files into the matching
// closed/YYYY/MM/DD/ folder when they live elsewhere.
func reconcileClosedRecordPath(repoRoot, absolutePath string) (string, bool, error) {
	fileDate := inferFileDatePrefix(absolutePath)
	if fileDate == "" || inferClosedDateFromPath(absolutePath) == fileDate {
		return absolutePath, false, nil
	}
	parts := strings.SplitN(fileDate, "-", 3)
	target := filepath.Join(repoRoot, ".hawp", "work", "closed", parts[0], parts[1], parts[2], filepath.Base(absolutePath))
	if target == absolutePath {
		return absolutePath, false, nil
	}
	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		return absolutePath, false, err
	}
	if _, err := os.Stat(target); err == nil {
		return absolutePath, false, nil
	}
	if err := os.Rename(absolutePath, target); err != nil {
		return absolutePath, false, err
	}
	return target, true, nil
}

// ApplyClosedRecordNormalization normalizes every closed record in place.
func ApplyClosedRecordNormalization(repoRoot string) (ApplyResult, error) {
	result := ApplyResult{}
	closedRoot := filepath.Join(repoRoot, ".hawp", "work", "closed")
	touched := map[string]struct{}{}

	for _, absolutePath := range walkPlanMarkdown(closedRoot) {
		currentPath, moved, err := reconcileClosedRecordPath(repoRoot, absolutePath)
		if err != nil {
			return result, err
		}
		if moved {
			touched[currentPath] = struct{}{}
		}

		raw, err := os.ReadFile(currentPath)
		if err != nil {
			return result, err
		}
		current := string(raw)
		next, addedClaims := normalizeClosedRecord(current, currentPath)

		if next == current {
			if !backlogIDLineRe.MatchString(current) && inferBacklogIDFromPath(absolutePath) == "" {
				result.SkippedFiles = append(result.SkippedFiles, absolutePath)
			}
			continue
		}
		if err := os.WriteFile(currentPath, []byte(next), 0o644); err != nil {
			return result, err
		}
		touched[currentPath] = struct{}{}

		for _, claim := range addedClaims {
			itemID := inferBacklogIDFromPath(currentPath)
			if itemID == "" {
				itemID = strings.ToUpper(strings.TrimSuffix(filepath.Base(currentPath), ".md"))
			}
			result.ResearchQueue = append(result.ResearchQueue, ResearchItem{
				ItemID: itemID, Claim: claim, FilePath: currentPath, LineNumber: 0,
				RecommendedAction: "Gather supporting proof for this verification claim, then replace the original checklist entry with an Evidence: citation or mark it explicitly unproven.",
			})
		}
	}

	for path := range touched {
		result.ChangedFiles = append(result.ChangedFiles, path)
	}
	return result, nil
}

// BuildResearchQueue lists ambiguous verification claims across all closed
// records for dry-run reporting.
func BuildResearchQueue(repoRoot string) []ResearchItem {
	closedRoot := filepath.Join(repoRoot, ".hawp", "work", "closed")
	closedFiles := closedFilesFromPaths(walkPlanMarkdown(closedRoot))
	clarity := CheckVerificationClarity(closedFiles)
	items := make([]ResearchItem, 0, len(clarity.Ambiguous))
	for _, claim := range clarity.Ambiguous {
		items = append(items, ResearchItem{
			ItemID: claim.ID, Claim: claim.Claim, FilePath: claim.FilePath, LineNumber: claim.LineNumber,
			RecommendedAction: "Research concrete supporting evidence for this verification claim, then update the checklist line with Evidence: ... or mark it explicitly unproven.",
		})
	}
	return items
}

// closedFilesFromPaths reads each path and wraps it in a source.File for
// use by validators that consume the snapshot-based API. This adapter exists
// inside the normalization layer, which already owns direct filesystem access.
func closedFilesFromPaths(paths []string) []source.File {
	files := make([]source.File, 0, len(paths))
	for _, path := range paths {
		raw, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		files = append(files, source.File{Path: path, RelPath: path, Content: string(raw)})
	}
	return files
}
