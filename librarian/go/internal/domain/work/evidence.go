package work

import (
	"regexp"
	"strings"

	"github.com/sentzunhat/hawp/librarian/go/internal/domain/work/source"
)

var evidenceLinkRe = regexp.MustCompile(`Evidence:\s*(?:link to )?\.\./evidence/([\w/.-]+\.md)`)

// CheckEvidenceIntegrity verifies that `Evidence: ../evidence/...` links in
// closed plans resolve to files present in the snapshot.
func CheckEvidenceIntegrity(snapshot source.Snapshot) EvidenceCheck {
	result := EvidenceCheck{Status: StatusPass}

	for _, file := range snapshot.ClosedFiles {
		fileName := strings.TrimSuffix(filenameFromPath(file.RelPath), ".md")

		for _, line := range strings.Split(file.Content, "\n") {
			m := evidenceLinkRe.FindStringSubmatch(line)
			if m == nil {
				continue
			}
			relativePath := m[1]
			// Reject any link that tries to escape the evidence folder.
			if strings.Contains(relativePath, "..") {
				continue
			}
			result.Total++
			evidenceRelPath := "evidence/" + strings.ReplaceAll(relativePath, "\\", "/")
			if hasFileWithRelPath(snapshot.Files, evidenceRelPath) {
				result.Valid++
			} else {
				result.Broken = append(result.Broken, BrokenLink{ID: fileName, Link: "../evidence/" + relativePath})
			}
		}
	}

	if len(result.Broken) > 0 {
		result.Status = StatusWarn
	}
	return result
}
