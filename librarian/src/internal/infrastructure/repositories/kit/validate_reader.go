// Package kit provides infrastructure-layer adapters for kit domain operations.
package kit

import (
	"os"

	domainkit "github.com/sentzunhat/hawp/librarian/src/internal/domain/kit"
	"github.com/sentzunhat/hawp/librarian/src/internal/infrastructure/markdown"
	"github.com/sentzunhat/hawp/librarian/src/internal/infrastructure/repo"
)

// kitSource wraps infrastructure implementations so the domain layer
// can be exercised without filesystem concerns (see validate_reader_test.go).
type kitSource struct{}

func (k *kitSource) FileLister(dir string, skipReadme bool) []string {
	return markdown.CollectFiles(dir, skipReadme)
}

func (k *kitSource) BlankFences(content string) string {
	return markdown.BlankFences(content)
}

func (k *kitSource) Exists(path string) bool {
	return repo.Exists(path)
}

func (k *kitSource) ToRepoRelative(root, rel string) string {
	return repo.ToRepoRelative(root, rel)
}

// CheckInternalLinks walks kitPath, reads each markdown file with os.ReadFile,
// and delegates link validation to the domain.
func CheckInternalLinks(kitPath string) []domainkit.Issue {
	src := &domainkit.KitSource{
		FileLister:   (*kitSource)(nil).FileLister,
		BlankFences:  (*kitSource)(nil).BlankFences,
		Exists:       (*kitSource)(nil).Exists,
		ToRepoRelative: (*kitSource)(nil).ToRepoRelative,
	}
	return src.CheckInternalLinks(kitPath, os.ReadFile)
}

// Validate runs all kit checks against kitPath using os.ReadFile for file content.
func Validate(kitPath string) ([]domainkit.Issue, int) {
	src := &domainkit.KitSource{
		FileLister:   (*kitSource)(nil).FileLister,
		BlankFences:  (*kitSource)(nil).BlankFences,
		Exists:       (*kitSource)(nil).Exists,
		ToRepoRelative: (*kitSource)(nil).ToRepoRelative,
	}
	return src.Validate(kitPath, os.ReadFile)
}
