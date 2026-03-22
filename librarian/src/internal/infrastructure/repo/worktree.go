package repo

import (
	"os/exec"
	"strings"
)

// HasDirtyWorktree reports whether the git worktree at repoRoot has
// uncommitted changes. Errors count as dirty so mutating commands fail
// closed.
func HasDirtyWorktree(repoRoot string) bool {
	cmd := exec.Command("git", "status", "--short")
	cmd.Dir = repoRoot
	out, err := cmd.Output()
	if err != nil {
		return true
	}
	return strings.TrimSpace(string(out)) != ""
}
