package context

import "github.com/sentzunhat/hawp/librarian/src/internal/domain/work"

// FileLister walks a directory tree and returns every markdown file it finds.
type FileLister interface {
	// CollectFiles returns all .md files under dir; skipReadme excludes
	// README.md from the results. Missing directories yield nil.
	CollectFiles(dir string, skipReadme bool) []string
}

// BacklogParser parses raw BACKLOG.md content into structured metadata.
type BacklogParser interface {
	ParseBacklog(raw string) *work.Backlog
}

// ContextSource bundles every data operation EnrichKit and EnrichWork need
// from the environment without leaking infrastructure concerns into domain.
//
// A thin application-layer adapter (see application/index/context-source.go)
// implements this interface by delegating to the real file-system / markdown
// / repo-root packages.
type ContextSource struct {
	FileLister   FileLister
	BacklogParser BacklogParser
}
