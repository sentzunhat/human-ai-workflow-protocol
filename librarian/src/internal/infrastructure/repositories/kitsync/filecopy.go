// Package kitsync provides infrastructure adapters for the kitsync file-copy engine.
// All os/io filesystem operations live here; domain packages delegate to this layer.
package kitsync

import (
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/sentzunhat/hawp/librarian/src/internal/infrastructure/filesystem"
)

// FileCopier handles low-level filesystem operations used by the kitsync engine.
type FileCopier struct{}

// NewFileCopier creates a new FileCopier backed by os/io calls.
func NewFileCopier() *FileCopier { return &FileCopier{} }

func (f *FileCopier) RejectSymlinkAncestors(root, target string) error {
	return filesystem.RejectSymlinkAncestors(root, target)
}

// MkdirAll calls os.MkdirAll with mode 0o755.
func (f *FileCopier) MkdirAll(dir string) error { return os.MkdirAll(dir, 0o755) }

// ReadDir calls os.ReadDir.
func (f *FileCopier) ReadDir(dir string) ([]fs.DirEntry, error) { return os.ReadDir(dir) }

// Stat calls os.Stat.
func (f *FileCopier) Stat(path string) (fs.FileInfo, error) { return os.Stat(path) }

// IsNotExist returns true if the error is os.ErrNotExist.
func (f *FileCopier) IsNotExist(err error) bool { return os.IsNotExist(err) }

// Open calls os.Open.
func (f *FileCopier) Open(path string) (io.ReadCloser, error) { return os.Open(path) }

// CreateTemp creates a temporary file in dir with the given pattern.
func (f *FileCopier) CreateTemp(dir, pattern string) (io.WriteCloser, string, error) {
	fm, err := os.CreateTemp(dir, pattern)
	if err != nil {
		return nil, "", err
	}
	return fm, fm.Name(), nil
}

// Rename calls os.Rename.
func (f *FileCopier) Rename(src, dst string) error { return os.Rename(src, dst) }

// Remove calls os.Remove.
func (f *FileCopier) Remove(path string) error { return os.Remove(path) }

// DirIsDir returns true if the entry is a directory.
func dirIsDir(entry os.DirEntry) bool {
	if entry.IsDir() {
		return true
	}
	info, err := entry.Info()
	if err != nil {
		return false
	}
	return info.IsDir()
}

// copyTree copies every file from srcDir into destDir, optionally filtered by
// an fnmatch-style pattern on the base filename. Returns count of files copied.
func (f *FileCopier) copyTree(srcDir, destDir, pattern string) (int, error) {
	entries, err := f.ReadDir(srcDir)
	if err != nil {
		return 0, err
	}
	if err := f.MkdirAll(destDir); err != nil {
		return 0, err
	}

	written := 0
	for _, entry := range entries {
		srcPath := filepath.Join(srcDir, entry.Name())
		destPath := filepath.Join(destDir, entry.Name())
		if dirIsDir(entry) {
			count, err := f.copyTree(srcPath, destPath, pattern)
			if err != nil {
				return written, err
			}
			written += count
			continue
		}

		if pattern != "" {
			if matched, _ := filepath.Match(pattern, entry.Name()); !matched {
				continue
			}
		}
		if err := f.copyFile(srcPath, destPath); err != nil {
			return written, err
		}
		written++
	}
	return written, nil
}

// seedTree copies only files that do not already exist at their destination.
func (f *FileCopier) seedTree(srcDir, destDir, pattern string) (int, error) {
	entries, err := f.ReadDir(srcDir)
	if err != nil {
		return 0, err
	}
	if err := f.MkdirAll(destDir); err != nil {
		return 0, err
	}

	written := 0
	for _, entry := range entries {
		srcPath := filepath.Join(srcDir, entry.Name())
		destPath := filepath.Join(destDir, entry.Name())
		if dirIsDir(entry) {
			count, err := f.seedTree(srcPath, destPath, pattern)
			if err != nil {
				return written, err
			}
			written += count
			continue
		}
		if pattern != "" {
			if matched, _ := filepath.Match(pattern, entry.Name()); !matched {
				continue
			}
		}
		if _, statErr := f.Stat(destPath); !f.IsNotExist(statErr) {
			continue // already exists
		}
		if err := f.copyFile(srcPath, destPath); err != nil {
			return written, err
		}
		written++
	}
	return written, nil
}

// copyFile copies srcPath to destPath using a temp file for atomicity.
func (f *FileCopier) copyFile(srcPath, destPath string) error {
	src, err := f.Open(srcPath)
	if err != nil {
		return err
	}
	defer src.Close()

	dir := filepath.Dir(destPath)
	if err := f.MkdirAll(dir); err != nil {
		return err
	}
	temp, tempPath, err := f.CreateTemp(dir, ".kitsync-*")
	if err != nil {
		return err
	}
	defer f.Remove(tempPath)

	if _, err = io.Copy(temp, src); err != nil {
		temp.Close()
		return err
	}
	if err := temp.Close(); err != nil {
		return err
	}
	return f.Rename(tempPath, destPath)
}
