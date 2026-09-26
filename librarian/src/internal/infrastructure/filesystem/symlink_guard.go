package filesystem

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
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
	fi, err := os.Lstat(root)
	if err != nil {
		return fmt.Errorf("stat root %s: %w", root, err)
	}
	if fi.Mode()&os.ModeSymlink != 0 {
		return fmt.Errorf("path %s is a symlink; refusing to follow", root)
	}

	rel, err := filepath.Rel(root, target)
	if err != nil {
		return fmt.Errorf("resolve path: %w", err)
	}
	if rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return fmt.Errorf("target %s is outside root %s; refusing to inspect", target, root)
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
			if isMacOSSystemAlias(current) {
				continue
			}
			return fmt.Errorf("path component %s is a symlink; refusing to follow", current)
		}
	}
	return nil
}

func isMacOSSystemAlias(path string) bool {
	if runtime.GOOS != "darwin" {
		return false
	}
	// macOS exposes these stable compatibility aliases into /private. They are
	// part of the platform's normal temporary/configuration path layout, not a
	// repository-owned redirection that callers can plant.
	switch path {
	case "/etc", "/tmp", "/var":
		return true
	default:
		return false
	}
}

// RejectSymlinksInTree verifies target is contained by root and rejects any
// symlink at or below target. Mutating commands can use this as an up-front
// fail-closed check before traversing a repository-owned tree.
//
// A missing target is allowed so callers that create the tree on first use can
// preserve their existing behavior. Existing ancestors are still checked by
// RejectSymlinkAncestors before that decision is made.
func RejectSymlinksInTree(root, target string) error {
	if err := RejectSymlinkAncestors(root, target); err != nil {
		return err
	}

	info, err := os.Lstat(target)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("stat tree root %s: %w", target, err)
	}
	if !info.IsDir() {
		return nil
	}

	return filepath.WalkDir(target, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return fmt.Errorf("walk path %s: %w", path, walkErr)
		}
		if entry.Type()&os.ModeSymlink != 0 {
			return fmt.Errorf("path component %s is a symlink; refusing to follow", path)
		}
		return nil
	})
}

// RejectSymlinksInPath rejects symlinks in every existing path component
// leading to target and anywhere below target when target is a directory.
// Unlike RejectSymlinksInTree, it does not require the caller to already know
// a trusted ancestor; it anchors the check at the filesystem volume root.
func RejectSymlinksInPath(target string) error {
	abs, err := filepath.Abs(target)
	if err != nil {
		return fmt.Errorf("resolve absolute path %s: %w", target, err)
	}
	volume := filepath.VolumeName(abs)
	root := string(filepath.Separator)
	if volume != "" {
		root = volume + string(filepath.Separator)
	}
	return RejectSymlinksInTree(root, abs)
}

// AtomicWriteFile writes data through a temporary sibling and atomically
// renames it into place after rechecking containment. Renaming replaces a
// substituted final-file symlink instead of following it.
func AtomicWriteFile(root, path string, data []byte, perm os.FileMode) error {
	if err := RejectSymlinkAncestors(root, path); err != nil {
		return err
	}
	dir := filepath.Dir(path)
	tmp, err := os.CreateTemp(dir, ".hawp-write-*")
	if err != nil {
		return err
	}
	tmpPath := tmp.Name()
	defer os.Remove(tmpPath)
	if err := tmp.Chmod(perm); err != nil {
		tmp.Close()
		return err
	}
	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	if err := RejectSymlinkAncestors(root, path); err != nil {
		return err
	}
	return os.Rename(tmpPath, path)
}

// RejectSymlinkedWorkRoot verifies a HAWP work tree from the repository level
// rather than treating .hawp/work as an independently trusted root. This makes
// a symlinked .hawp ancestor visible to every work-record mutation.
func RejectSymlinkedWorkRoot(workRoot string) error {
	repoRoot := filepath.Dir(filepath.Dir(workRoot))
	if err := RejectSymlinkAncestors(repoRoot, workRoot); err != nil {
		return fmt.Errorf("work root: %w", err)
	}
	return nil
}
