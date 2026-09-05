package work

import "io/fs"

// WorkSource provides infrastructure dependencies to domain/work functions.
// Mirrors the KitSource pattern: domain defines function-typed fields;
// application constructs the source from infrastructure implementations.
type WorkSource struct {
	Exists         func(path string) bool
	ToRepoRelative func(repoRoot, absolutePath string) string
	CollectFiles   func(dir string, skipReadme bool) []string

	// Migration filesystem capabilities are supplied by the application layer.
	// Keeping them here preserves the parent WorkSource API while preventing
	// domain/work from choosing a concrete operating-system implementation.
	ReadDir      func(path string) ([]fs.DirEntry, error)
	ReadFile     func(path string) ([]byte, error)
	WriteFile    func(path string, data []byte, perm fs.FileMode) error
	MkdirAll     func(path string, perm fs.FileMode) error
	MkdirTemp    func(dir, pattern string) (string, error)
	Rename       func(oldPath, newPath string) error
	Remove       func(path string) error
	RemoveAll    func(path string) error
	Stat         func(path string) (fs.FileInfo, error)
	Lstat        func(path string) (fs.FileInfo, error)
	EvalSymlinks func(path string) (string, error)
}
