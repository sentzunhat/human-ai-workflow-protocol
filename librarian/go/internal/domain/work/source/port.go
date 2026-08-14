// Package source defines the read-only workspace boundary for work validation.
package source

// File is a Markdown record acquired from a work workspace.
type File struct {
	Path    string
	RelPath string
	Content string
	Links   []Link
}

// Link is one Markdown link already parsed by the source adapter.
type Link struct{ Href string }

// Snapshot is the acquired input for work validation rules.
type Snapshot struct {
	BacklogContent string
	Files          []File
	ExistingPaths  map[string]struct{}
}

// Workspace acquires work records. Application services compose it; domain
// rules consume the acquired data without owning filesystem access.
type Workspace interface {
	Read(workDir string) (Snapshot, error)
}
