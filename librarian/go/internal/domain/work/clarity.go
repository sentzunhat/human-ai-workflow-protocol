package work

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	verificationSectionRe = regexp.MustCompile(`^##\s+Verification\b`)
	claimPrefixRe         = regexp.MustCompile(`^[\s\-\[\]x ]+`)
	unprovenRe            = regexp.MustCompile(`(?i)\b(?:explicitly )?unproven\b`)
)

// CheckVerificationClarity scans Verification sections in closed plans for
// checklist claims, classifying each as proven (Evidence:), explicitly
// unproven, or ambiguous.
func CheckVerificationClarity(closedFiles []string) VerificationCheck {
	result := VerificationCheck{Status: StatusPass}

	for _, filePath := range closedFiles {
		content, err := os.ReadFile(filePath)
		if err != nil {
			warnf("skipping unreadable closed plan %s: %v", filePath, err)
			continue
		}
		fileName := strings.TrimSuffix(filepath.Base(filePath), ".md")

		section := extractVerificationSection(string(content))
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
			entry := Claim{ID: fileName, Claim: claim, FilePath: filePath, LineNumber: index + 1}

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
