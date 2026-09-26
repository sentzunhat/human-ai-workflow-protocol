package query

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"unicode/utf8"

	appcontext "github.com/sentzunhat/hawp/librarian/src/internal/application/context"
	appcontext_layout1 "github.com/sentzunhat/hawp/librarian/src/internal/application/context/dedup"
	appsearch "github.com/sentzunhat/hawp/librarian/src/internal/application/search"
	"github.com/sentzunhat/hawp/librarian/src/internal/bootstrap"
	domainsearch "github.com/sentzunhat/hawp/librarian/src/internal/domain/search"
	"github.com/sentzunhat/hawp/librarian/src/internal/infrastructure/repo"
)

func Run(args []string, cwd string) error {
	opts, err := parseArgs(args)
	if err != nil {
		return err
	}
	query, limit := opts.query, opts.limit
	wantContext, wantSemantic := opts.context, opts.semantic
	format, maxTokens := opts.format, opts.maxTokens
	verbose, hybridRatio := opts.verbose, opts.hybridRatio

	root, err := repo.FindBacklogRepoRoot(cwd)
	if err != nil {
		fmt.Println("Not in a HAWP repo; no index to search.")
		return nil
	}

	execution, err := bootstrap.NewSearchService().Execute(root, appsearch.QueryOptions{
		Query:       query,
		Limit:       limit,
		Semantic:    wantSemantic,
		HybridRatio: float32(hybridRatio),
	})
	if err != nil {
		var missingIndex appsearch.IndexNotFoundError
		if errors.As(err, &missingIndex) {
			fmt.Printf("Index not found at %s. Run `hawp search index` first.\n", missingIndex.Path)
			return nil
		}
		return err
	}

	if wantSemantic && !execution.HasVectors {
		fmt.Println("No vectors found. Run `hawp search embed` first to enable semantic search.")
		return nil
	}

	results := execution.Rows
	if wantSemantic && results == nil {
		fmt.Printf("Semantic search failed for %q — check that your embedding backend is running.\n", query)
		return nil
	}

	if len(results) == 0 {
		fmt.Printf("No results found for %q\n", query)
		return nil
	}

	if !wantContext {
		fmt.Printf("Search results for %q (%d found):\n\n", query, len(results))
		for i, result := range results {
			fmt.Printf("[%d] %s | %s (chunk %v)\n",
				i+1,
				getStr(result, "path"),
				getStr(result, "folder_role"),
				getInt(result, "chunk_idx"),
			)
			fmt.Printf("    Context: %s\n", getStr(result, "folder_context"))
			text := getStr(result, "text")
			text = truncatePreview(text, 150)
			fmt.Printf("    %q\n\n", text)
		}
		return nil
	}

	searchResults := appsearch.RowsToResults(results, execution.HasVectors)

	deduped, droppedByDedup := appcontext_layout1.ContentJaccardDedup(searchResults, 0.70)

	avgChunkTokens := 0
	if len(searchResults) > 0 {
		total := 0
		for _, r := range searchResults {
			total += (len(r.Content) + 3) / 4
		}
		avgChunkTokens = total / len(searchResults)
	}

	capped := make([]domainsearch.Result, 0, len(deduped))
	runningTokens := 0
	for _, r := range deduped {
		chunkEst := (len(r.Content) + 3) / 4
		if len(capped) > 0 && runningTokens+chunkEst > maxTokens {
			break
		}
		capped = append(capped, r)
		runningTokens += chunkEst
	}

	if verbose {
		savedTokens := droppedByDedup * avgChunkTokens
		fmt.Fprintf(os.Stderr, "context: %d chunks, ~%d tokens (saved ~%d tokens via dedup)\n",
			len(capped), runningTokens, savedTokens)
	}

	block := appcontext.FormatAsMarkdown(capped, query, maxTokens)

	switch format {
	case "json":
		jsonBlock := map[string]interface{}{
			"title":        block.Title,
			"query":        block.Query,
			"result_count": block.ResultCount,
			"token_count":  block.TokenCount,
			"results":      block.Results,
			"references":   toJSONReferences(block.References),
			"metadata":     block.Metadata,
		}
		out, err := json.MarshalIndent(jsonBlock, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal JSON: %w", err)
		}
		fmt.Println(string(out))
	default: // markdown
		fmt.Println(block.String())
	}

	return nil
}

func toJSONReferences(refs []appcontext.DocumentReference) []map[string]interface{} {
	out := make([]map[string]interface{}, len(refs))
	for i, r := range refs {
		out[i] = map[string]interface{}{
			"source":     r.Source,
			"title":      r.Title,
			"content":    r.Content,
			"relevance":  r.Relevance,
			"line_start": r.LineStart,
			"line_end":   r.LineEnd,
		}
	}
	return out
}

func getStr(m map[string]interface{}, key string) string {
	if v, ok := m[key]; ok && v != nil {
		if s, ok := v.(*string); ok && s != nil {
			return *s
		}
		return fmt.Sprintf("%v", v)
	}
	return ""
}

func getInt(m map[string]interface{}, key string) int64 {
	if v, ok := m[key]; ok {
		if i, ok := v.(int64); ok {
			return i
		}
	}
	return 0
}

func truncatePreview(text string, maxBytes int) string {
	if len(text) <= maxBytes {
		return text
	}
	end := maxBytes
	for end > 0 && !utf8.RuneStart(text[end]) {
		end--
	}
	return text[:end] + "..."
}
