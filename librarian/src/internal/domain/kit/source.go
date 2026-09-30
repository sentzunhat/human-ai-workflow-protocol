package kit

import "io/fs"

// KitSource is a thin interface for data operations that PlanFileRenames,
// CheckInternalLinks, and their callers need without touching infrastructure
// concerns directly.
//
// Application-layer adapters delegate to the real filesystem, markdown parsing,
// and repo-path helpers so the domain stays pure.
type KitSource struct {
	// ReadDir lists entries under a directory.
	ReadDir func(dir string) ([]fs.DirEntry, error)
	// FileLister returns all .md files under kitPath; skipReadme excludes
	// README.md from results. Missing directories yield nil.
	FileLister func(kitPath string, skipReadme bool) []string

	// BlankFences masks fenced code blocks so link extraction only finds
	// references in actual prose.
	BlankFences func(content string) string

	// ReadFile reads a single file's bytes.
	ReadFile func(path string) ([]byte, error)

	// Exists reports whether the path points to an existing file or directory.
	Exists func(path string) bool

	// ToRepoRelative returns a repo-relative path for kitPath relative to root.
	ToRepoRelative func(root, rel string) string
}
