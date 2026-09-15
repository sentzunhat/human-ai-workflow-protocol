package migration

import (
	"fmt"
	"os"
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
	if err == nil && !info.Mode().IsRegular() {
		return fmt.Errorf("non-regular artifact: %s", name)
	}
	return os.WriteFile(full, data, 0o644)
}
