// Package source defines the filesystem boundary for the kit capability.
package source

// Entry is a file or directory discovered in a kit workspace.
type Entry struct {
	Path    string
	RelPath string
	Name    string
	IsDir   bool
}

// Link is a parsed Markdown link. Offset identifies the opening bracket in
// the source content so updates can be applied without reparsing in domain.
type Link struct {
	Href   string
	Offset int
}

// File is a Markdown document available to kit rules.
type File struct {
	Path    string
	RelPath string
	Content string
	Links   []Link
}

// Snapshot is the read-only kit workspace input consumed by domain rules.
type Snapshot struct {
	Entries []Entry
	Files   []File
}

// Rename is one filesystem rename planned by the domain kit rules.
type Rename struct {
	From string
	To   string
}

// LinkUpdate is one Markdown href rewrite planned by the domain kit rules.
type LinkUpdate struct {
	File  string
	From  string
	To    string
	Start int
	End   int
}

// Workspace owns kit acquisition and mutation. Application services compose
// it; domain rules only consume Snapshot and emit mutation plans.
type Workspace interface {
	Read(kitPath string) (Snapshot, error)
	ApplyRenames(renames []Rename) (conflictFrom, conflictTo string, err error)
	ApplyLinkUpdates(updates []LinkUpdate) (int, error)
}
