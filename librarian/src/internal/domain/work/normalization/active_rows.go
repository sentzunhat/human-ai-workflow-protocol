package normalization

import (
	"fmt"
	"path/filepath"
	"strings"
)

// ApplyCompletedActiveRowCleanup removes completed Active Work rows only when
// they already point at an existing closed plan. This handles older repos that
// copied the row to Recently Closed but left a stale done/wont-fix row active.
func ApplyCompletedActiveRowCleanup(repoRoot string, source ActiveSource) (ApplyResult, error) {
	result := ApplyResult{}
	workRoot := filepath.Join(repoRoot, ".hawp", "work")
	backlogPath := filepath.Join(workRoot, "BACKLOG.md")
	backlog, err := ParseNormalizeBacklog(backlogPath, source.ScanSource)
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
		resolved := resolveWithinWorkRoot(workRoot, planPath)
		if resolved == "" {
			continue
		}
		// Verify the resolved (cleaned) path is under <workRoot>/closed/.
		// The raw planPath prefix check is insufficient: closed/../active/x/plan.md
		// passes the string test but resolves into the active tree.
		closedRoot := filepath.Join(workRoot, "closed")
		closedResolved, err := source.EvalSymlinks(closedRoot)
		if err != nil {
			continue
		}
		resolvedTarget, err := source.EvalSymlinks(resolved)
		if err != nil {
			continue
		}
		relative, err := filepath.Rel(closedResolved, resolvedTarget)
		if err != nil || relative == ".." || strings.HasPrefix(relative, ".."+string(filepath.Separator)) {
			continue
		}
		if info, err := source.Lstat(resolved); err == nil && info.Mode().IsRegular() {
			removeLines[row.LineNumber] = struct{}{}
		}
	}
	if len(removeLines) == 0 {
		return result, nil
	}

	// Guard against a symlinked BACKLOG.md redirecting the write outside the repo.
	backlogTarget, err := source.EvalSymlinks(backlogPath)
	if err != nil {
		return result, err
	}
	workTarget, err := source.EvalSymlinks(workRoot)
	if err != nil {
		return result, err
	}
	relative, err := filepath.Rel(workTarget, backlogTarget)
	if err != nil || relative != "BACKLOG.md" {
		return result, fmt.Errorf("refusing to rewrite %s: symlinked path escapes work root", backlogPath)
	}
	if fi, err := source.Lstat(backlogPath); err != nil {
		return result, err
	} else if !fi.Mode().IsRegular() {
		return result, fmt.Errorf("refusing to rewrite %s: not a regular file", backlogPath)
	}
	raw, err := source.ReadFile(backlogPath)
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
	if err := source.WriteFile(backlogPath, []byte(strings.Join(kept, "\n")), 0o644); err != nil {
		return result, err
	}
	result.ChangedFiles = append(result.ChangedFiles, source.ToRepoRelative(repoRoot, backlogPath))
	return result, nil
}
