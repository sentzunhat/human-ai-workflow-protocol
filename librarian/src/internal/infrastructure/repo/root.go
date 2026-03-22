// Package repo locates the HAWP repository root and converts paths for
// display. Mirrors librarian/scripts/lib (findUpward, findBacklogRepoRoot,
// toRepoRelative).
package repo

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

// LegacyClosedCutoff: closed files on or after this date require Outcome,
// Verification, and Close Checklist sections; earlier files are legacy.
const LegacyClosedCutoff = "2026-05-10"

// FindUpward walks from startDir toward the filesystem root until predicate
// matches a directory. Returns "" when nothing matches within maxDepth.
func FindUpward(startDir string, predicate func(dir string) bool, maxDepth int) string {
	current, err := filepath.Abs(startDir)
	if err != nil {
		return ""
	}
	for depth := 0; depth < maxDepth; depth++ {
		if predicate(current) {
			return current
		}
		parent := filepath.Dir(current)
		if parent == current {
			break
		}
		current = parent
	}
	return ""
}

// FindBacklogRepoRoot finds the ancestor directory containing
// .hawp/work/BACKLOG.md.
func FindBacklogRepoRoot(startDir string) (string, error) {
	root := FindUpward(startDir, func(dir string) bool {
		_, err := os.Stat(filepath.Join(dir, ".hawp", "work", "BACKLOG.md"))
		return err == nil
	}, 12)
	if root == "" {
		return "", errors.New("could not locate repo root containing .hawp/work/BACKLOG.md")
	}
	return root, nil
}

// ToRepoRelative converts an absolute path to a POSIX-style path relative to
// repoRoot; falls back to the input when it cannot be made relative.
func ToRepoRelative(repoRoot, absolutePath string) string {
	rel, err := filepath.Rel(repoRoot, absolutePath)
	if err != nil {
		return absolutePath
	}
	return strings.ReplaceAll(rel, "\\", "/")
}

// Exists reports whether the path exists (any file type).
func Exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
