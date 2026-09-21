package migration

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Refuse occupied destinations, including ignored files and symlink ancestors.
// Apply is a single-writer operation; stop editors/other migrations before use.
func checkDestinations(root string, p plan) error {
	base := filepath.Join(root, sourceRoot)
	for _, e := range p.Files {
		if !safeRelative(e.Destination) {
			return fmt.Errorf("unsafe destination: %s", e.Destination)
		}
		full := filepath.Join(base, filepath.FromSlash(e.Destination))
		for current := full; current != root; current = filepath.Dir(current) {
			info, err := os.Lstat(current)
			if os.IsNotExist(err) {
				continue
			}
			if err != nil {
				return err
			}
			if info.Mode()&os.ModeSymlink != 0 {
				return fmt.Errorf("symlink destination: %s", e.Destination)
			}
			if current == full {
				if e.Source != e.Destination || !info.Mode().IsRegular() {
					return fmt.Errorf("occupied destination: %s", e.Destination)
				}
			} else if !info.IsDir() {
				return fmt.Errorf("non-directory destination parent: %s", e.Destination)
			}
		}
	}
	// Collect absent destinations and reject any that would be gitignored.
	// inventory uses git ls-files --exclude-standard, which omits ignored
	// untracked files; an ignored destination would disappear from the
	// post-apply state check and roll the migration back.
	var absent []string
	for _, e := range p.Files {
		if !safeRelative(e.Destination) {
			continue // already rejected above
		}
		full := filepath.Join(base, filepath.FromSlash(e.Destination))
		if _, err := os.Lstat(full); os.IsNotExist(err) {
			absent = append(absent, filepath.Join(sourceRoot, e.Destination))
		}
	}
	if len(absent) > 0 {
		cmd := exec.Command("git", "-C", root, "check-ignore", "--stdin", "-z")
		cmd.Stdin = strings.NewReader(strings.Join(absent, "\x00") + "\x00")
		out, _ := cmd.Output()
		if len(bytes.TrimSpace(out)) > 0 {
			ignored := strings.TrimRight(string(out), "\x00")
			return fmt.Errorf("destination matches .gitignore and would be excluded from the post-apply inventory check: %s", ignored)
		}
	}
	return nil
}

// Generated review artifacts must not modify the inventoried source tree or
// follow a symlink. An existing artifact is replaced only on an explicit flag.
func writeArtifact(root, name string, data []byte) error {
	root, err := filepath.EvalSymlinks(root)
	if err != nil {
		return err
	}
	full, err := filepath.Abs(name)
	if err != nil {
		return err
	}
	parent, err := filepath.EvalSymlinks(filepath.Dir(full))
	if err != nil {
		return err
	}
	full = filepath.Join(parent, filepath.Base(full))
	rel, err := filepath.Rel(filepath.Join(root, sourceRoot), full)
	if err != nil {
		return err
	}
	if rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return fmt.Errorf("review artifacts must be outside %s", sourceRoot)
	}
	info, err := os.Lstat(full)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	if err == nil {
		if !info.Mode().IsRegular() {
			return fmt.Errorf("non-regular artifact: %s", name)
		}
		// Unlink before writing so a hard-linked artifact cannot redirect the
		// write to its paired source inode; WriteFile truncates in place.
		if err := os.Remove(full); err != nil {
			return fmt.Errorf("remove existing artifact: %w", err)
		}
	}
	return os.WriteFile(full, data, 0o644)
}
