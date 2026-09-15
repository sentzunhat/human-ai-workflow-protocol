package work

import (
	"os"
	"path/filepath"
	"strings"
)

// ApplyCompletedActiveRowCleanup removes completed Active Work rows only when
// they already point at an existing closed plan. This handles older repos that
// copied the row to Recently Closed but left a stale done/wont-fix row active.
func (w *WorkSource) ApplyCompletedActiveRowCleanup(repoRoot string) (ApplyResult, error) {
	result := ApplyResult{}
	workRoot := filepath.Join(repoRoot, ".hawp", "work")
	backlogPath := filepath.Join(workRoot, "BACKLOG.md")
	backlog, err := ParseNormalizeBacklog(backlogPath)
	if err != nil {
		return result, err
	}

	removeLines := map[int]struct{}{}
	for _, row := range backlog.Rows {
		if row.Section != SectionActive {
			continue
		}
		status := strings.ToLower(strings.TrimSpace(row.Status))
		if status != "done" && status != "wont-fix" {
			continue
		}
		planPath := strings.TrimSpace(row.PlanPath)
		if !strings.HasPrefix(filepath.ToSlash(planPath), "closed/") {
			continue
		}
		if resolved := resolveWithinWorkRoot(workRoot, planPath); resolved != "" {
			if _, err := os.Stat(resolved); err == nil {
				removeLines[row.LineNumber] = struct{}{}
			}
		}
	}
	if len(removeLines) == 0 {
		return result, nil
	}

	raw, err := os.ReadFile(backlogPath)
	if err != nil {
		return result, err
	}
	var kept []string
	for i, line := range strings.Split(string(raw), "\n") {
		if _, ok := removeLines[i+1]; ok {
			continue
		}
		kept = append(kept, line)
	}
	if err := os.WriteFile(backlogPath, []byte(strings.Join(kept, "\n")), 0o644); err != nil {
		return result, err
	}
	result.ChangedFiles = append(result.ChangedFiles, w.ToRepoRelative(repoRoot, backlogPath))
	return result, nil
}
