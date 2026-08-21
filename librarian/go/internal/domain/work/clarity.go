package work

import (
	"regexp"
	"strings"

	"github.com/sentzunhat/hawp/librarian/go/internal/domain/work/source"
)

var (
	verificationSectionRe = regexp.MustCompile(`^##\s+Verification\b`)
	claimPrefixRe         = regexp.MustCompile(`^[\s\-\[\]x ]+`)
	unprovenRe            = regexp.MustCompile(`(?i)\b(?:explicitly )?unproven\b`)
)

// CheckVerificationClarity scans Verification sections in closed plan files for
// checklist claims, classifying each as proven (Evidence:), explicitly
// unproven, or ambiguous.
func CheckVerificationClarity(closedFiles []source.File) VerificationCheck {
	result := VerificationCheck{Status: StatusPass}

	for _, file := range closedFiles {
		fileName := strings.TrimSuffix(filenameFromPath(file.RelPath), ".md")

		section := extractVerificationSection(file.Content)
		if section == "" {
			continue
		}

		for index, line := range strings.Split(section, "\n") {
			if !strings.Contains(line, "- [x]") && !strings.Contains(line, "- [ ]") {
				continue
			}
			claim := claimPrefixRe.ReplaceAllString(line, "")
			if len(claim) > 100 {
				claim = claim[:100]
			}
			if strings.HasPrefix(claim, "Research evidence for:") ||
				strings.HasPrefix(claim, "Update the original verification checklist line with Evidence:") {
				continue
			}

			result.Total++
			entry := Claim{ID: fileName, Claim: claim, FilePath: file.Path, LineNumber: index + 1}

			switch {
			case strings.Contains(line, "Evidence:"):
				result.Proven++
			case strings.Contains(line, "NOT YET VERIFIED") || unprovenRe.MatchString(line):
				result.Unproven = append(result.Unproven, entry)
			default:
				result.Ambiguous = append(result.Ambiguous, entry)
			}
		}
	}

	if len(result.Unproven) > 0 || len(result.Ambiguous) > 0 {
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
