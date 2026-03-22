package work

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var evidenceLinkRe = regexp.MustCompile(`Evidence:\s*(?:link to )?\.\./evidence/([\w/.-]+\.md)`)

// CollectClosedPlanFiles gathers all .md files in the strict
// closed/YYYY/MM/DD/ layout (including README.md, matching the TS walker).
func CollectClosedPlanFiles(closedDir string) []string {
	var files []string
	years, err := os.ReadDir(closedDir)
	if err != nil {
		return files
	}
	for _, year := range years {
		if year.Name() == "README.md" || !year.IsDir() {
			continue
		}
		months, err := os.ReadDir(filepath.Join(closedDir, year.Name()))
		if err != nil {
			continue
		}
		for _, month := range months {
			if !month.IsDir() {
				continue
			}
			days, err := os.ReadDir(filepath.Join(closedDir, year.Name(), month.Name()))
			if err != nil {
				continue
			}
			for _, day := range days {
				if !day.IsDir() {
					continue
				}
				entries, err := os.ReadDir(filepath.Join(closedDir, year.Name(), month.Name(), day.Name()))
				if err != nil {
					continue
				}
				for _, entry := range entries {
					if strings.HasSuffix(entry.Name(), ".md") {
						files = append(files, filepath.Join(closedDir, year.Name(), month.Name(), day.Name(), entry.Name()))
					}
				}
			}
		}
	}
	return files
}

// CheckEvidenceIntegrity verifies that `Evidence: ../evidence/...` links in
// closed plans resolve inside the evidence folder.
func CheckEvidenceIntegrity(workDir string, closedFiles []string) EvidenceCheck {
	result := EvidenceCheck{Status: StatusPass}
	evidenceRoot := filepath.Join(workDir, "evidence")

	for _, filePath := range closedFiles {
		content, err := os.ReadFile(filePath)
		if err != nil {
			warnf("skipping unreadable closed plan %s: %v", filePath, err)
			continue
		}
		fileName := strings.TrimSuffix(filepath.Base(filePath), ".md")

		for _, line := range strings.Split(string(content), "\n") {
			m := evidenceLinkRe.FindStringSubmatch(line)
			if m == nil {
				continue
			}
			relativePath := m[1]
			fullPath := filepath.Clean(filepath.Join(evidenceRoot, relativePath))
			if !strings.HasPrefix(fullPath, evidenceRoot+string(filepath.Separator)) {
				warnf("evidence link escapes evidence folder, skipping: %s", relativePath)
				continue
			}
			result.Total++
			if _, err := os.Stat(fullPath); err == nil {
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
