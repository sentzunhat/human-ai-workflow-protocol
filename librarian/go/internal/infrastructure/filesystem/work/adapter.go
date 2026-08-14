// Package work provides the filesystem adapter for work validation input.
package work

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/sentzunhat/hawp/librarian/go/internal/domain/work/source"
	"github.com/sentzunhat/hawp/librarian/go/internal/infrastructure/markdown"
)

// Adapter acquires a read-only work validation snapshot from disk.
type Adapter struct{}

func NewAdapter() Adapter { return Adapter{} }

func (Adapter) Read(workDir string) (source.Snapshot, error) {
	backlog, err := os.ReadFile(filepath.Join(workDir, "BACKLOG.md"))
	if err != nil {
		return source.Snapshot{}, err
	}

	snapshot := source.Snapshot{BacklogContent: string(backlog), ExistingPaths: map[string]struct{}{}}
	for _, path := range markdown.CollectFiles(workDir, false) {
		raw, readErr := os.ReadFile(path)
		if readErr != nil {
			continue
		}
		rel, relErr := filepath.Rel(workDir, path)
		if relErr != nil {
			rel = path
		}
		file := source.File{Path: path, RelPath: strings.ReplaceAll(rel, "\\", "/"), Content: string(raw)}
		for _, link := range markdown.ExtractLinks(markdown.BlankFences(file.Content)) {
			file.Links = append(file.Links, source.Link{Href: link.Href})
		}
		snapshot.Files = append(snapshot.Files, file)
		snapshot.ExistingPaths[path] = struct{}{}
	}
	return snapshot, nil
}
