package kitsync

import (
	"io"
	"io/fs"
)

// FileCopier abstracts filesystem operations used by the kitsync engine.
// The domain layer delegates all os/io calls through this interface;
// infrastructure/repositories/kitsync/filecopy.go supplies the concrete
// implementation and its test double.
type FileCopier interface {
	// MkdirAll creates a directory including parents (mode 0o755).
	MkdirAll(dir string) error
	// ReadDir lists directory entries.
	ReadDir(dir string) ([]fs.DirEntry, error)
	// Stat returns file info.
	Stat(path string) (fs.FileInfo, error)
	// IsNotExist reports whether err means the path does not exist.
	IsNotExist(err error) bool
	// Open opens a file for reading.
	Open(path string) (io.ReadCloser, error)
	// CreateTemp creates a temporary file in dir with pattern;
	// returns the file handle and its path.
	CreateTemp(dir, pattern string) (io.WriteCloser, string, error)
	// Rename renames a file (used for atomic write via temp file).
	Rename(src, dst string) error
	// Remove removes a file.
	Remove(path string) error
}
