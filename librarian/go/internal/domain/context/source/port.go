// Package source defines the input boundary for context corpus enrichment.
package source

import domainwork "github.com/sentzunhat/hawp/librarian/go/internal/domain/work"

// File is one corpus document acquired by a source adapter.
// RelPath is relative to the corpus root; RepoPath is relative to the repo.
type File struct {
	RelPath  string
	RepoPath string
	Content  string
}

// KitCorpus is the acquired kit input passed to context enrichment.
type KitCorpus struct {
	Files []File
}

// WorkCorpus is the acquired work input and its parsed backlog metadata.
type WorkCorpus struct {
	Files   []File
	Backlog *domainwork.Backlog
}

// CorpusSource acquires repository content for the context capability.
type CorpusSource interface {
	ReadKit(repoRoot, kitPath string) (KitCorpus, error)
	ReadWork(repoRoot, workPath string) (WorkCorpus, error)
}
