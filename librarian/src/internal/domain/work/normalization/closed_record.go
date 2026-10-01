package normalization

import (
	"regexp"
	"strings"
)

var (
	legacyIDInPathRe = regexp.MustCompile(`(?i)(TASK|BUG)-\d+`)
	backlogIDLineRe  = regexp.MustCompile(`(?i)\*\*Backlog ID:\*\*`)
	h1Re             = regexp.MustCompile(`^#\s+`)
	multiBlankRe     = regexp.MustCompile(`\n{3,}`)
	researchEntryRe  = regexp.MustCompile(`(?m)^\s*-\s+\[ \]\s+Research evidence for: (.+)$`)
	verifyHeadBodyRe = regexp.MustCompile(`(?ms)(^##\s+Verification\b[^\n]*\n)(.*?)(?:^##\s+|\z)`)
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
			inserted := append(append(append([]string{}, lines[:i+1]...), "", "**Backlog ID:** "+backlogID, ""), lines[i+1:]...)
			return multiBlankRe.ReplaceAllString(strings.Join(inserted, "\n"), "\n\n")
		}
	}
	return "**Backlog ID:** " + backlogID + "\n\n" + content
}

func ensureEvidenceFollowUp(content string) (string, []string) {
	ambiguous := AmbiguousVerificationClaims(content)
	if len(ambiguous) == 0 {
		return content, nil
	}
	loc := verifyHeadBodyRe.FindStringSubmatchIndex(content)
	if loc == nil {
		return content, nil
	}
	heading, body := content[loc[2]:loc[3]], content[loc[4]:loc[5]]
	existing := map[string]struct{}{}
	for _, m := range researchEntryRe.FindAllStringSubmatch(body, -1) {
		existing[strings.TrimSpace(m[1])] = struct{}{}
	}
	var pending []string
	for _, claim := range ambiguous {
		claim = strings.TrimSpace(claim)
		if claim != "" {
			if _, ok := existing[claim]; !ok {
				pending = append(pending, claim)
			}
		}
	}
	if len(pending) == 0 {
		return content, nil
	}
	var entries []string
	for _, claim := range pending {
		entries = append(entries, "- [ ] Research evidence for: "+claim+"\n- [ ] Update the original verification checklist line with Evidence: ... or explicit unproven wording.")
	}
	newEntries := strings.Join(entries, "\n")
	const subsection = "### Evidence Follow-Up"
	trimmedBody := strings.TrimRight(body, " \t\n")
	if strings.Contains(body, subsection) {
		body = trimmedBody + "\n" + newEntries + "\n"
	} else {
		body = trimmedBody + "\n\n" + subsection + "\n\n" + newEntries + "\n"
	}
	return content[:loc[2]] + heading + body + content[loc[5]:], pending
}

// NormalizeClosedRecord scaffolds missing sections and evidence follow-ups.
func NormalizeClosedRecord(content, filePath string) (string, []string) {
	updated := content
	if id := inferBacklogIDFromPath(filePath); id != "" {
		updated = insertBacklogID(updated, id)
	}
	updated = appendSection(updated, "Outcome", "_Legacy normalization scaffold added._")
	updated = appendSection(updated, "Verification", "_Legacy normalization scaffold added._")
	updated = appendSection(updated, "Close Checklist", "- [ ] Legacy normalization scaffold added.")
	return ensureEvidenceFollowUp(updated)
}
