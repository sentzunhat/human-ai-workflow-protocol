package context

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/sentzunhat/hawp/librarian/go/internal/infrastructure/markdown"
	"github.com/sentzunhat/hawp/librarian/go/internal/infrastructure/repo"
)

// kitRole classifies a kit-relative path by its top-level segment, e.g.
// "usage/init.md" -> "usage"; "start-here.md" (no subfolder) -> "root".
func kitRole(kitRelPath string) string {
	segments := strings.Split(filepath.ToSlash(kitRelPath), "/")
	if len(segments) <= 1 {
		return "root"
	}
	return segments[0]
}

// firstDescriptiveLine returns the first non-empty, non-heading line of
// content — typically a folder README's one-line description right
// after its title — or "" if none is found.
func firstDescriptiveLine(content string) string {
	for _, line := range strings.Split(content, "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}
		return trimmed
	}
	return ""
}

// EnrichKit walks kitPath and returns every markdown document (including
// README.md files) tagged with its folder role and a context prefix.
// Folders with a README.md contribute their first descriptive line as
// shared context for every other document in that folder.
func EnrichKit(repoRoot, kitPath string) ([]Document, error) {
	files := markdown.CollectFiles(kitPath, false)

	folderContext := map[string]string{}
	for _, file := range files {
		if filepath.Base(file) != "README.md" {
			continue
		}
		raw, err := os.ReadFile(file)
		if err != nil {
			continue
		}
		role := kitRole(repo.ToRepoRelative(kitPath, file))
		folderContext[role] = firstDescriptiveLine(string(raw))
	}

	documents := make([]Document, 0, len(files))
	for _, file := range files {
		raw, err := os.ReadFile(file)
		if err != nil {
			continue
		}
		kitRel := repo.ToRepoRelative(kitPath, file)
		role := kitRole(kitRel)

		prefix := "[kit/" + role + "]"
		if summary := folderContext[role]; summary != "" {
			prefix += " " + summary
		}

		documents = append(documents, Document{
			RelPath:       repo.ToRepoRelative(repoRoot, file),
			Corpus:        CorpusKit,
			Role:          role,
			ContextPrefix: prefix,
			Content:       string(raw),
		})
	}
	return documents, nil
}
