package work

import (
	"path/filepath"
	"regexp"
	"strings"

	"github.com/sentzunhat/hawp/librarian/src/internal/domain/work/normalization"
	"github.com/sentzunhat/hawp/librarian/src/internal/domain/work/validation"
)

var (
	legacyIDInPathRe   = regexp.MustCompile(`(?i)(TASK|BUG)-\d+`)
	closedDateInPathRe = regexp.MustCompile(`/closed/(\d{4})/(\d{2})/(\d{2})/`)
	backlogIDLineRe    = regexp.MustCompile(`(?i)\*\*Backlog ID:\*\*`)
	h1Re               = regexp.MustCompile(`^#\s+`)
	multiBlankRe       = regexp.MustCompile(`\n{3,}`)
	researchEntryRe    = regexp.MustCompile(`(?m)^\s*-\s+\[ \]\s+Research evidence for: (.+)$`)
	verifyHeadBodyRe   = regexp.MustCompile(`(?ms)(^##\s+Verification\b[^\n]*\n)(.*?)(?:^##\s+|\z)`)
)

func inferBacklogIDFromPath(path string) string {
	return strings.ToUpper(legacyIDInPathRe.FindString(path))
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
	return normalization.NormalizeClosedRecord(content, filePath)
}

// ApplyClosedRecordNormalization normalizes every closed record in place.
func (w *WorkSource) ApplyClosedRecordNormalization(repoRoot string) (ApplyResult, error) {
	return normalization.ApplyClosedRecordNormalization(repoRoot, normalization.ClosedSource{
		ScanSource: normalization.ScanSource{ReadDir: w.ReadDir, ReadFile: w.ReadFile, Stat: w.Stat},
		MkdirAll:   w.MkdirAll, Rename: w.Rename, WriteFile: w.WriteFile,
		RejectSymlinkAncestors: w.RejectSymlinkAncestors,
	})
}

// BuildResearchQueue lists ambiguous verification claims across all closed
// records for dry-run reporting.
func BuildResearchQueue(repoRoot string, source *WorkSource) []ResearchItem {
	closedRoot := filepath.Join(repoRoot, ".hawp", "work", "closed")
	clarity := validation.CheckVerificationClarity(normalization.WalkPlanMarkdown(closedRoot, normalization.ScanSource{ReadDir: source.ReadDir}), validation.Source{ReadFile: source.ReadFile})
	items := make([]ResearchItem, 0, len(clarity.Ambiguous))
	for _, claim := range clarity.Ambiguous {
		items = append(items, ResearchItem{
			ItemID: claim.ID, Claim: claim.Claim, FilePath: claim.FilePath, LineNumber: claim.LineNumber,
			RecommendedAction: "Research concrete supporting evidence for this verification claim, then update the checklist line with Evidence: ... or mark it explicitly unproven.",
		})
	}
	return items
}
