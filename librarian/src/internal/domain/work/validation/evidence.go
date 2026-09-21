package validation

import (
	"path/filepath"
	"regexp"
	"strings"
)

var evidenceLinkRe = regexp.MustCompile(`Evidence:\*{0,2}\s*(?:inline or )?(?:link to )?((?:\.\./evidence/|\.hawp/work/evidence/)[\w/.-]+\.md(?:#[\w.-]+)?)`)

// CollectClosedPlanFiles gathers closed-plan markdown files from both the
// legacy flat layout (`closed/YYYY/MM/DD/<id>.md`) and the current
// folder-per-item layout (`closed/YYYY/MM/DD/<id>/plan.md`).
func CollectClosedPlanFiles(closedDir string, source Source) []string {
	var files []string
	years, err := source.ReadDir(closedDir)
	if err != nil {
		return files
	}
	for _, year := range years {
		if year.Name() == "README.md" || !year.IsDir() {
			continue
		}
		months, err := source.ReadDir(filepath.Join(closedDir, year.Name()))
		if err != nil {
			continue
		}
		for _, month := range months {
			if !month.IsDir() {
				continue
			}
			days, err := source.ReadDir(filepath.Join(closedDir, year.Name(), month.Name()))
			if err != nil {
				continue
			}
			for _, day := range days {
				if !day.IsDir() {
					continue
				}
				entries, err := source.ReadDir(filepath.Join(closedDir, year.Name(), month.Name(), day.Name()))
				if err != nil {
					continue
				}
				for _, entry := range entries {
					if strings.HasSuffix(entry.Name(), ".md") {
						files = append(files, filepath.Join(closedDir, year.Name(), month.Name(), day.Name(), entry.Name()))
						continue
					}
					if !entry.IsDir() {
						continue
					}
					planPath := filepath.Join(closedDir, year.Name(), month.Name(), day.Name(), entry.Name(), "plan.md")
					if _, err := source.Stat(planPath); err == nil {
						files = append(files, planPath)
					}
				}
			}
		}
	}
	return files
}

func closedPlanID(filePath string) string {
	base := filepath.Base(filePath)
	if base == "plan.md" {
		return filepath.Base(filepath.Dir(filePath))
	}
	return strings.TrimSuffix(base, ".md")
}

// CheckEvidenceIntegrity verifies both legacy `../evidence/...` references
// and current repo-relative `.hawp/work/evidence/...` references in closed
// plans, resolving both forms inside the work evidence folder.
func CheckEvidenceIntegrity(workDir string, closedFiles []string, source Source) EvidenceCheck {
	result := EvidenceCheck{Status: StatusPass}
	evidenceRoot := filepath.Join(workDir, "evidence")

	for _, filePath := range closedFiles {
		content, err := source.ReadFile(filePath)
		if err != nil {
			warnf(source, "skipping unreadable closed plan %s: %v", filePath, err)
			continue
		}
		fileName := closedPlanID(filePath)

		for _, line := range strings.Split(string(content), "\n") {
			m := evidenceLinkRe.FindStringSubmatch(line)
			if m == nil {
				continue
			}
			reference := m[1]
			pathPart, _, _ := strings.Cut(reference, "#")
			relativePath := ""
			switch {
			case strings.HasPrefix(pathPart, "../evidence/"):
				relativePath = strings.TrimPrefix(pathPart, "../evidence/")
			case strings.HasPrefix(pathPart, ".hawp/work/evidence/"):
				relativePath = strings.TrimPrefix(pathPart, ".hawp/work/evidence/")
			default:
				continue
			}
			fullPath := filepath.Clean(filepath.Join(evidenceRoot, filepath.FromSlash(relativePath)))
			if !strings.HasPrefix(fullPath, evidenceRoot+string(filepath.Separator)) {
				warnf(source, "evidence link escapes evidence folder, skipping: %s", reference)
				continue
			}
			result.Total++
			if _, err := source.Stat(fullPath); err == nil {
				result.Valid++
			} else {
				result.Broken = append(result.Broken, BrokenLink{ID: fileName, Link: reference})
			}
		}
	}

	if len(result.Broken) > 0 {
		result.Status = StatusWarn
	}
	return result
}
