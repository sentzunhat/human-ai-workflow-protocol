package context

import (
	"path/filepath"
	"strings"

	"github.com/sentzunhat/hawp/librarian/go/internal/domain/context/source"
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

// EnrichKit converts acquired kit files into context-aware documents.
// File discovery and reading belong to the source adapter, not this package.
func EnrichKit(corpus source.KitCorpus) []Document {
	folderContext := map[string]string{}
	for _, file := range corpus.Files {
		if filepath.Base(file.RelPath) != "README.md" {
			continue
		}
		folderContext[kitRole(file.RelPath)] = firstDescriptiveLine(file.Content)
	}

	documents := make([]Document, 0, len(corpus.Files))
	for _, file := range corpus.Files {
		role := kitRole(file.RelPath)
		prefix := "[kit/" + role + "]"
		if summary := folderContext[role]; summary != "" {
			prefix += " " + summary
		}

		documents = append(documents, Document{
			RelPath:       file.RepoPath,
			Corpus:        CorpusKit,
			Role:          role,
			ContextPrefix: prefix,
			Content:       file.Content,
		})
	}
	return documents
}
