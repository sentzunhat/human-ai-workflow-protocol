package mcp

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"

	appcontext "github.com/sentzunhat/hawp/librarian/src/internal/application/context"
	appsearch "github.com/sentzunhat/hawp/librarian/src/internal/application/search"
	domainsearch "github.com/sentzunhat/hawp/librarian/src/internal/domain/search"
	"github.com/sentzunhat/hawp/librarian/src/internal/infrastructure/sqlite"
)

const contextRadius = 40 // lines above/below source line for context window

func searchToolDef() map[string]any {
	return map[string]any{
		"name":        "hawp_search",
		"description": "Search indexed HAWP kit and work documents. By default returns structured results with source path, relevance, chunk content, and line positions. Pass context:true to instead receive a single deduplicated, token-capped markdown block ready for LLM injection (Jaccard dedup + greedy token cap applied before formatting). Run `hawp search index` then `hawp search embed` first to build the index.",
		"inputSchema": map[string]any{
			"type": "object",
			"properties": map[string]any{
				"query": map[string]any{
					"type":        "string",
					"description": "Search query",
				},
				"limit": map[string]any{
					"type":        "integer",
					"description": "Max results to return (default 5)",
				},
				"context": map[string]any{
					"type":        "boolean",
					"description": "When true, return a pre-shaped markdown context block (deduplicated + token-capped) instead of raw results",
				},
				"max_tokens": map[string]any{
					"type":        "integer",
					"description": "Token budget for the context block when context:true (default 2000)",
				},
			},
			"required": []string{"query"},
		},
	}
}

func toolSearch(args json.RawMessage, repoRoot string) rpcResponse {
	var a struct {
		Query     string `json:"query"`
		Limit     int    `json:"limit"`
		Context   bool   `json:"context"`
		MaxTokens int    `json:"max_tokens"`
	}
	if err := json.Unmarshal(args, &a); err != nil {
		return toolErr("invalid args: " + err.Error())
	}
	if a.Query == "" {
		return toolErr("query is required")
	}
	if a.Limit <= 0 {
		a.Limit = 5
	}
	if a.MaxTokens <= 0 {
		a.MaxTokens = 2000
	}

	if a.Context {
		return toolSearchContext(a.Query, a.Limit, a.MaxTokens, repoRoot)
	}

	dbPath := filepath.Join(repoRoot, ".hawp", "db", "index.sqlite")
	db, err := sqlite.Open(dbPath)
	if err != nil {
		return toolErr(fmt.Sprintf("index not found at %s; run `hawp search index` first", dbPath))
	}
	defer db.Close()

	rows, err := db.QueryChunksLexical(a.Query, a.Limit*3)
	if err != nil {
		return toolErr("search failed: " + err.Error())
	}
	if len(rows) == 0 {
		return jsonResult(SearchResponse{Query: a.Query, Results: []SearchResult{}})
	}

	hasVectors, _ := db.HasVectors()
	if hasVectors {
		rows = appsearch.HybridRank(rows, a.Query, db, a.Limit, 0)
	} else if len(rows) > a.Limit {
		rows = rows[:a.Limit]
	}

	results := make([]SearchResult, 0, len(rows))
	for _, r := range rows {
		path := strVal(r, "path")
		content := strVal(r, "text")
		lineStart := int(intVal(r, "line_start"))
		lineEnd := int(intVal(r, "line_end"))

		var relevance float32 = 0.95
		if hasVectors {
			relevance = float32(floatVal(r, "_hybrid_score"))
		}

		sourceLine := findSourceLine(content, a.Query, lineStart)
		window := LineRange{
			Start: max1(sourceLine - contextRadius),
			End:   sourceLine + contextRadius,
		}

		results = append(results, SearchResult{
			Source:    path,
			Relevance: relevance,
			Content:   content,
			Lines: LineInfo{
				Range:  LineRange{Start: lineStart, End: lineEnd},
				Source: sourceLine,
			},
			Context: ContextInfo{Window: window},
		})
	}

	return jsonResult(SearchResponse{Query: a.Query, Results: results})
}

// toolSearchContext runs dedup + dynamic token cap + markdown format and returns
// a single pre-shaped context block ready for LLM injection.
func toolSearchContext(query string, limit, maxTokens int, repoRoot string) rpcResponse {
	results, err := appsearch.Query(repoRoot, query, limit)
	if err != nil {
		return toolErr("search failed: " + err.Error())
	}
	if len(results) == 0 {
		return jsonResult(ContextSearchResponse{Query: query, Budget: maxTokens})
	}

	deduped, dropped := appcontext.ContentJaccardDedup(results, 0.70)

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

	block := appcontext.FormatAsMarkdown(capped, query, maxTokens)
	return jsonResult(ContextSearchResponse{
		Query:         query,
		Content:       block.String(),
		TokenCount:    block.TokenCount,
		Budget:        maxTokens,
		ChunksUsed:    block.ResultCount,
		ChunksDropped: dropped,
	})
}

// findSourceLine returns the 1-indexed line in the source file where the first
// query term appears within chunkText. Falls back to chunkStartLine when no
// term matches (e.g. semantic-only match).
func findSourceLine(chunkText, query string, chunkStartLine int) int {
	terms := strings.Fields(strings.ToLower(query))
	for i, line := range strings.Split(chunkText, "\n") {
		lower := strings.ToLower(line)
		for _, term := range terms {
			if strings.Contains(lower, term) {
				return chunkStartLine + i
			}
		}
	}
	return chunkStartLine
}

func max1(n int) int {
	if n < 1 {
		return 1
	}
	return n
}

func strVal(m map[string]interface{}, key string) string {
	v, ok := m[key]
	if !ok || v == nil {
		return ""
	}
	if s, ok := v.(*string); ok && s != nil {
		return *s
	}
	return fmt.Sprintf("%v", v)
}

func intVal(m map[string]interface{}, key string) int64 {
	if v, ok := m[key]; ok {
		if i, ok := v.(int64); ok {
			return i
		}
	}
	return 0
}

func floatVal(m map[string]interface{}, key string) float64 {
	if v, ok := m[key]; ok {
		if f, ok := v.(float64); ok {
			return f
		}
	}
	return 0
}
