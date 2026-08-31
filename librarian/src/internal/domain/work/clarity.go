package work

import (
	"os"
	"regexp"
	"strings"
)

var (
	verificationSectionRe = regexp.MustCompile(`^##\s+Verification\b`)
	verificationItemRe    = regexp.MustCompile(`^\s*-\s+\[(?:x| )\]`)
	claimPrefixRe         = regexp.MustCompile(`^[\s\-\[\]x ]+`)
	unprovenRe            = regexp.MustCompile(`(?i)\b(?:explicitly )?unproven\b`)
)

type checklistItem struct {
	StartLine int
	Lines     []string
}

// CheckVerificationClarity scans Verification sections in closed plans for
// checklist claims, classifying each as proven (Evidence:), explicitly
// unproven, or ambiguous. Explicitly unproven claims are tracked for context
// but do not warn on their own because the uncertainty is already labeled.
func CheckVerificationClarity(closedFiles []string) VerificationCheck {
	result := VerificationCheck{Status: StatusPass}

	for _, filePath := range closedFiles {
		content, err := os.ReadFile(filePath)
		if err != nil {
			warnf("skipping unreadable closed plan %s: %v", filePath, err)
			continue
		}
		fileName := closedPlanID(filePath)

		section := extractVerificationSection(string(content))
		if section == "" {
			continue
		}

		for _, item := range extractChecklistItems(section) {
			line := item.Lines[0]
			claim := claimPrefixRe.ReplaceAllString(line, "")
			if len(claim) > 100 {
				claim = claim[:100]
			}
			if strings.HasPrefix(claim, "Research evidence for:") ||
				strings.HasPrefix(claim, "Update the original verification checklist line with Evidence:") {
				continue
			}

			result.Total++
			entry := Claim{ID: fileName, Claim: claim, FilePath: filePath, LineNumber: item.StartLine}
			body := strings.Join(item.Lines, "\n")

			switch {
			case strings.Contains(body, "Evidence:"):
				result.Proven++
			case strings.Contains(body, "NOT YET VERIFIED") || unprovenRe.MatchString(body):
				result.Unproven = append(result.Unproven, entry)
			default:
				result.Ambiguous = append(result.Ambiguous, entry)
			}
		}
	}

	if len(result.Ambiguous) > 0 {
		result.Status = StatusWarn
	}
	return result
}

// extractVerificationSection returns the lines under "## Verification" up to
// the next "## " heading, or "" when the section is absent or empty.
func extractVerificationSection(content string) string {
	var section []string
	inVerification := false
	for _, line := range strings.Split(content, "\n") {
		if verificationSectionRe.MatchString(strings.TrimSpace(line)) {
			inVerification = true
			continue
		}
		if inVerification {
			if strings.HasPrefix(line, "## ") {
				break
			}
			section = append(section, line)
		}
	}
	return strings.Join(section, "\n")
}

func extractChecklistItems(section string) []checklistItem {
	var items []checklistItem
	var current *checklistItem

	for index, line := range strings.Split(section, "\n") {
		if verificationItemRe.MatchString(line) {
			if current != nil {
				items = append(items, *current)
			}
			current = &checklistItem{StartLine: index + 1, Lines: []string{line}}
			continue
		}
		if current != nil {
			current.Lines = append(current.Lines, line)
		}
	}

	if current != nil {
		items = append(items, *current)
	}

	return items
}
