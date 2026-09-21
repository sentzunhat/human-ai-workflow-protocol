package normalization

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"regexp"
	"strings"
)

// ClosedSource supplies the filesystem capabilities needed to normalize
// closed records. The application layer owns the concrete implementation.
type ClosedSource struct {
	ScanSource
	MkdirAll               func(path string, perm fs.FileMode) error
	Rename                 func(oldPath, newPath string) error
	WriteFile              func(path string, data []byte, perm fs.FileMode) error
	RejectSymlinkAncestors func(root, target string) error
}

var closedFileDatePrefixRe = regexp.MustCompile(`^(\d{4})-(\d{2})-(\d{2})-`)

func closedDateFromPath(path string) string {
	m := regexp.MustCompile(`/closed/(\d{4})/(\d{2})/(\d{2})/`).FindStringSubmatch(strings.ReplaceAll(path, "\\", "/"))
	if m == nil {
		return ""
	}
	return m[1] + "-" + m[2] + "-" + m[3]
}

func reconcileClosedRecordPath(repoRoot, absolutePath string, source ClosedSource) (string, bool, error) {
	m := closedFileDatePrefixRe.FindStringSubmatch(filepath.Base(absolutePath))
	if m == nil {
		return absolutePath, false, nil
	}
	fileDate := m[1] + "-" + m[2] + "-" + m[3]
	if closedDateFromPath(absolutePath) == fileDate {
		return absolutePath, false, nil
	}
	target := filepath.Join(repoRoot, ".hawp", "work", "closed", m[1], m[2], m[3], filepath.Base(absolutePath))
	if target == absolutePath {
		return absolutePath, false, nil
	}
	closedRoot := filepath.Join(repoRoot, ".hawp", "work", "closed")
	if err := source.RejectSymlinkAncestors(closedRoot, target); err != nil {
		return absolutePath, false, err
	}
	if err := source.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		return absolutePath, false, err
	}
	if _, err := source.Stat(target); err == nil {
		return absolutePath, false, nil
	}
	if err := source.Rename(absolutePath, target); err != nil {
		return absolutePath, false, err
	}
	return target, true, nil
}

// ApplyClosedRecordNormalization normalizes closed records and reconciles
// date-prefixed records into their canonical closed/YYYY/MM/DD folder.
func ApplyClosedRecordNormalization(repoRoot string, source ClosedSource) (ApplyResult, error) {
	result := ApplyResult{}
	closedRoot := filepath.Join(repoRoot, ".hawp", "work", "closed")
	if err := source.RejectSymlinkAncestors(repoRoot, closedRoot); err != nil {
		return result, fmt.Errorf("refusing symlinked closed work root: %w", err)
	}
	touched := map[string]struct{}{}

	for _, absolutePath := range WalkPlanMarkdown(closedRoot, source.ScanSource) {
		if err := source.RejectSymlinkAncestors(closedRoot, absolutePath); err != nil {
			return result, fmt.Errorf("refusing symlinked closed record: %w", err)
		}
		currentPath, moved, err := reconcileClosedRecordPath(repoRoot, absolutePath, source)
		if err != nil {
			return result, err
		}
		if err := source.RejectSymlinkAncestors(closedRoot, currentPath); err != nil {
			return result, fmt.Errorf("refusing symlinked closed record: %w", err)
		}
		if moved {
			touched[currentPath] = struct{}{}
		}

		raw, err := source.ReadFile(currentPath)
		if err != nil {
			return result, err
		}
		current := string(raw)
		next, addedClaims := NormalizeClosedRecord(current, currentPath)
		if next == current {
			if !strings.Contains(current, "**Backlog ID:**") && inferBacklogIDFromPath(absolutePath) == "" {
				result.SkippedFiles = append(result.SkippedFiles, absolutePath)
			}
			continue
		}
		if err := source.WriteFile(currentPath, []byte(next), 0o644); err != nil {
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
