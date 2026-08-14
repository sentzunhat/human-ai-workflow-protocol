// Package context provides the filesystem adapter for the context corpus.
package context

import (
	"os"
	"path/filepath"

	"github.com/sentzunhat/hawp/librarian/go/internal/domain/context/source"
	domainwork "github.com/sentzunhat/hawp/librarian/go/internal/domain/work"
	"github.com/sentzunhat/hawp/librarian/go/internal/infrastructure/markdown"
	"github.com/sentzunhat/hawp/librarian/go/internal/infrastructure/repo"
)

// Adapter acquires kit and work files from a repository filesystem.
type Adapter struct{}

func NewAdapter() Adapter {
	return Adapter{}
}

// ReadKit preserves the previous kit behavior: all markdown files, including
// README.md files, are collected; unreadable files are skipped.
func (Adapter) ReadKit(repoRoot, kitPath string) (source.KitCorpus, error) {
	files := markdown.CollectFiles(kitPath, false)
	corpus := source.KitCorpus{Files: make([]source.File, 0, len(files))}
	for _, path := range files {
		raw, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		corpus.Files = append(corpus.Files, source.File{
			RelPath:  repo.ToRepoRelative(kitPath, path),
			RepoPath: repo.ToRepoRelative(repoRoot, path),
			Content:  string(raw),
		})
	}
	return corpus, nil
}

// ReadWork preserves the previous work behavior: backlog parsing is required,
// known work folders are scanned recursively, README.md files are skipped,
// and unreadable documents are ignored.
func (Adapter) ReadWork(repoRoot, workPath string) (source.WorkCorpus, error) {
	backlog, err := domainwork.ParseBacklog(filepath.Join(workPath, "BACKLOG.md"))
	if err != nil {
		return source.WorkCorpus{}, err
	}

	var corpus source.WorkCorpus
	corpus.Backlog = backlog
	for _, role := range []string{"active", "closed", "parked", "decisions", "evidence", "notes", "status"} {
		dir := filepath.Join(workPath, role)
		for _, path := range markdown.CollectFiles(dir, true) {
			raw, readErr := os.ReadFile(path)
			if readErr != nil {
				continue
			}
			corpus.Files = append(corpus.Files, source.File{
				RelPath:  repo.ToRepoRelative(workPath, path),
				RepoPath: repo.ToRepoRelative(repoRoot, path),
				Content:  string(raw),
			})
		}
	}
	return corpus, nil
}
