package context

import (
	"fmt"
	"sort"
	"strings"

	"github.com/sentzunhat/hawp/librarian/src/internal/domain/search"
)

// ContextBlock is formatted search results ready for LLM injection.
type ContextBlock struct {
	Title       string              // e.g., "Search Results for 'vector embedding'"
	Query       string              // Original query
	ResultCount int                 // Number of results included
	TokenCount  int                 // Approximate token usage
	Budget      int                 // Token budget passed to FormatAsMarkdown
	Results     []FormattedResult   // Ordered, deduplicated results
	References  []DocumentReference // Deduplicated source documents (by Source)
	Metadata    map[string]string   // Query timestamp, filters, etc
}

// FormattedResult is one result formatted for readability.
type FormattedResult struct {
	Rank      int     // Position (1, 2, 3, ...)
	Relevance float32 // Confidence (0.0 - 1.0)
	Source    string  // Document source/path
	Title     string  // Chunk title/heading
	Content   string  // Actual text content
	Tokens    int     // Approximate tokens in this result
	LineStart int     // Line number where chunk starts in source
	LineEnd   int     // Line number where chunk ends in source
}

// FormatAsMarkdown converts deduplicated results into LLM-ready markdown.
// Respects token budget by truncating content as needed.
// Also collects and deduplicates document references by Source.
func FormatAsMarkdown(results []search.Result, query string, maxTokens int) ContextBlock {
	block := ContextBlock{
		Title:    fmt.Sprintf("Search Results: \"%s\"", query),
		Query:    query,
		Budget:   maxTokens,
		Metadata: make(map[string]string),
	}

	// Sort by relevance descending
	sorted := make([]search.Result, len(results))
	copy(sorted, results)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Relevance > sorted[j].Relevance
	})

	// Reserve space for the compact header, then add results until the budget is
	// exhausted based on the rendered representation rather than a fixed surcharge.
	usedTokens := estimateTokens(renderHeader(block.Title, len(sorted)))

	for rank, result := range sorted {
		if usedTokens >= maxTokens {
			break // Stop adding if we exceed budget
		}

		candidate := FormattedResult{
			Rank:      rank + 1,
			Relevance: result.Relevance,
			Source:    result.Source,
			Title:     result.Title,
			Content:   result.Content,
			LineStart: result.LineStart,
			LineEnd:   result.LineEnd,
		}
		referenceTokens := estimateTokens(renderReferenceLine(candidate)) + 1
		if usedTokens+referenceTokens >= maxTokens {
			break // No room for even a minimal result
		}

		resultTokens := estimateTokens(result.Content)
		availableTokens := maxTokens - usedTokens - referenceTokens

		if resultTokens > availableTokens {
			if availableTokens <= 5 {
				break
			}
			// Truncate this result to fit available space
			truncated := truncateToTokens(result.Content, availableTokens-5) // -5 for truncation markers
			resultTokens = estimateTokens(truncated)
			result.Content = truncated
			candidate.Content = truncated
		}

		candidate.Tokens = resultTokens
		block.Results = append(block.Results, candidate)

		usedTokens += estimateTokens(renderFormattedResult(candidate))
	}

	block.References = deduplicateReferences(block.Results)
	block.ResultCount = len(block.Results)
	block.TokenCount = estimateTokens(block.String())
	block.Metadata["result_count"] = fmt.Sprintf("%d", block.ResultCount)
	block.Metadata["token_count"] = fmt.Sprintf("%d", block.TokenCount)
	block.Metadata["budget"] = fmt.Sprintf("%d", maxTokens)

	return block
}

// deduplicateReferences collapses FormattedResults into one DocumentReference
// per unique Source, keeping the highest-relevance chunk's Relevance and
// Content, sorted by Relevance descending. Shared by FormatAsMarkdown (the
// normal path) and by ContextReshaper.Reshape as a fallback for callers who
// construct a ContextBlock by hand (e.g. RAGPipeline consumers) without
// going through FormatAsMarkdown first.
func deduplicateReferences(results []FormattedResult) []DocumentReference {
	refMap := make(map[string]DocumentReference)
	for _, result := range results {
		if ref, exists := refMap[result.Source]; exists {
			if result.Relevance > ref.Relevance {
				ref.Relevance = result.Relevance
				ref.Content = result.Content
				ref.LineStart = result.LineStart
				ref.LineEnd = result.LineEnd
				refMap[result.Source] = ref
			}
		} else {
			refMap[result.Source] = DocumentReference{
				Source:    result.Source,
				Title:     result.Title,
				Content:   result.Content,
				Relevance: result.Relevance,
				LineStart: result.LineStart,
				LineEnd:   result.LineEnd,
			}
		}
	}

	refs := make([]DocumentReference, 0, len(refMap))
	for _, ref := range refMap {
		refs = append(refs, ref)
	}
	sort.Slice(refs, func(i, j int) bool {
		return refs[i].Relevance > refs[j].Relevance
	})
	return refs
}

// String renders the context block as formatted markdown: for small/sparse
// result sets it emits only the interleaved provenance+content body; for
// denser sets it prepends the compact title header before the same body.
func (cb ContextBlock) String() string {
	var sb strings.Builder
	sb.WriteString(renderHeader(cb.Title, len(cb.Results)))
	sb.WriteString(formatResultsInline(cb.Results))
	return sb.String()
}

// formatResultsInline renders just the body — each result's `Ref:`
// line interleaved immediately above its own content, rather than collected
// into one list at the end — so a reader (human or LLM) sees which source a
// chunk came from right where it's used, without having to cross-reference
// a footnote list. No title/summary header; callers that want one (like
// ContextBlock.String) prepend it themselves. Also used directly as the
// "none"-pipeline passthrough content in ContextReshaper.Reshape, so a
// header wouldn't get double-rendered by the RAG output wrapper.
func formatResultsInline(results []FormattedResult) string {
	var sb strings.Builder
	for _, result := range results {
		sb.WriteString(renderFormattedResult(result))
	}
	return sb.String()
}

func renderHeader(title string, resultCount int) string {
	if title == "" || resultCount <= 5 {
		return ""
	}
	return fmt.Sprintf("# %s\n\n", title)
}

func renderReferenceLine(result FormattedResult) string {
	return fmt.Sprintf("Ref: %s\n", formatReferenceLocation(result.Source, result.LineStart, result.LineEnd))
}

func renderFormattedResult(result FormattedResult) string {
	return renderReferenceLine(result) + result.Content + "\n\n"
}

func formatReferenceLocation(source string, lineStart, lineEnd int) string {
	location := source
	switch {
	case lineStart > 0 && lineEnd > 0 && lineEnd >= lineStart:
		location = fmt.Sprintf("%s:%d-%d", source, lineStart, lineEnd)
	case lineStart > 0:
		location = fmt.Sprintf("%s:%d", source, lineStart)
	}
	return location
}

// estimateTokens estimates token count using ~4 chars per token heuristic.
func estimateTokens(text string) int {
	// Rough estimate: ~4 characters per token for English
	return (len(text) + 3) / 4
}

// truncateToTokens truncates text to approximately the given token count.
func truncateToTokens(text string, maxTokens int) string {
	maxChars := maxTokens * 4

	if len(text) <= maxChars {
		return text
	}

	// Truncate at character boundary
	truncated := text[:maxChars]

	// Try to truncate at word boundary for better readability
	if lastSpace := strings.LastIndex(truncated, " "); lastSpace > maxChars/2 {
		truncated = truncated[:lastSpace]
	}

	return strings.TrimRight(truncated, " \n\t") + "..."
}
