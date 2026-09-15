package filesystem

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// RejectSymlinkAncestors walks every directory component of target, starting
// at root, and returns an error if any component (including root itself) is a
// symlink. This prevents MkdirAll or WriteFile from following a symlink planted
// inside an otherwise-trusted directory tree and redirecting writes outside it.
//
// root must be an ancestor of target (or equal to it). Components that do not
// yet exist are skipped — MkdirAll is expected to create them.
func RejectSymlinkAncestors(root, target string) error {
	// Check root itself first — callers often trust root without verifying it.
	if fi, err := os.Lstat(root); err == nil {
		if fi.Mode()&os.ModeSymlink != 0 {
			return fmt.Errorf("path %s is a symlink; refusing to follow", root)
		}
	}

	rel, err := filepath.Rel(root, target)
	if err != nil {
		return fmt.Errorf("resolve path: %w", err)
	}
	parts := strings.Split(rel, string(filepath.Separator))
	current := root
	for _, part := range parts {
		if part == "." {
			continue
		}
		current = filepath.Join(current, part)
		fi, err := os.Lstat(current)
		if err != nil {
			if os.IsNotExist(err) {
				break // not yet created — MkdirAll will create it safely
			}
			return fmt.Errorf("stat path component %s: %w", current, err)
		}
		if fi.Mode()&os.ModeSymlink != 0 {
			return fmt.Errorf("path component %s is a symlink; refusing to follow", current)
		}
	}
	return nil
}
