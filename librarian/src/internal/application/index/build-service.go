package index

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	domaincontext "github.com/sentzunhat/hawp/librarian/src/internal/domain/context"
	domainindex "github.com/sentzunhat/hawp/librarian/src/internal/domain/index"
)

// BuildResult is the enriched document corpus for the requested scope,
// ready for Slice 2's chunking/embedding step (fbf12a93) to consume.
type BuildResult struct {
	Scope     domainindex.DocumentScope
	Documents []domaincontext.Document
}

type BuildService struct {
	RepoRoot string
}

func NewBuildService(repoRoot string) BuildService {
	return BuildService{RepoRoot: repoRoot}
}

// Execute walks the requested corpus (or both) and returns every
// document enriched with its folder/record context.
func (service BuildService) Execute(scope domainindex.DocumentScope) (BuildResult, error) {
	result := BuildResult{Scope: scope}

	if scope == domainindex.ScopeAll || scope == domainindex.ScopeKit {
		docs, err := domaincontext.EnrichKit(service.RepoRoot, filepath.Join(service.RepoRoot, ".hawp", "kit"))
		if err != nil {
			return BuildResult{}, fmt.Errorf("kit enrichment: %w", err)
		}
		result.Documents = append(result.Documents, docs...)
	}

	if scope == domainindex.ScopeAll || scope == domainindex.ScopeWork {
		docs, err := domaincontext.EnrichWork(service.RepoRoot, filepath.Join(service.RepoRoot, ".hawp", "work"))
		if err != nil {
			return BuildResult{}, fmt.Errorf("work enrichment: %w", err)
		}
		result.Documents = append(result.Documents, docs...)
	}

	return result, nil
}

// String renders a summary report: counts per corpus/role, and for work
// documents, counts per type. Raw content is never printed here — use
// Export for the full enriched corpus.
func (result BuildResult) String() string {
	byCorpusRole := map[string]int{}
	byType := map[string]int{}
	withMetadata, total := 0, len(result.Documents)

	for _, doc := range result.Documents {
		byCorpusRole[string(doc.Corpus)+"/"+doc.Role]++
		if doc.Type != "" {
			byType[doc.Type]++
			withMetadata++
		}
	}

	lines := fmt.Sprintf("index:build\n===========\nscope: %s\ndocuments: %d\n\n", result.Scope, total)
	lines += "By corpus/role:\n"
	for _, key := range sortedKeys(byCorpusRole) {
		lines += fmt.Sprintf("  %s: %d\n", key, byCorpusRole[key])
	}
	if len(byType) > 0 {
		lines += "\nWork records with resolved backlog metadata:\n"
		for _, key := range sortedKeys(byType) {
			lines += fmt.Sprintf("  %s: %d\n", key, byType[key])
		}
		lines += fmt.Sprintf("  (%d of %d work documents matched a backlog row)\n", withMetadata, countWork(result.Documents))
	}
	return lines
}

func countWork(docs []domaincontext.Document) int {
	n := 0
	for _, doc := range docs {
		if doc.Corpus == domaincontext.CorpusWork {
			n++
		}
	}
	return n
}

func sortedKeys(m map[string]int) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// Export writes the full enriched document corpus as indented JSON —
// the exact input Slice 2's chunking/embedding step will consume.
func (result BuildResult) Export(path string) error {
	out, err := json.MarshalIndent(result.Documents, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, append(out, '\n'), 0o644)
}
