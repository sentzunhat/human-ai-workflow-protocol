// Package index defines domain models and operations for document indexing.
package index

import (
	"regexp"
	"strings"
)

// Document is a document to be indexed (kit guide, work plan, evidence file, etc).
type Document struct {
	ID         int64   // assigned by DB
	Category   string  // "kit" or "work"
	Type       string  // "guide", "plan", "backlog-row", "evidence"
	Path       string  // file source
	FolderRole string  // "kit/start-here", "work/active", etc.
	Content    string  // full text
}

// DocumentMetadata is work-item-specific metadata (optional, for category='work').
type DocumentMetadata struct {
	DocumentID int64
	WorkUUID   string
	Status     string
	Owner      *string
	RiskLevel  *string
	ReportedAt *string
	ClosedAt   *string
}

// Chunk is a searchable unit within a document.
type Chunk struct {
	DocumentID    int64   // FK to documents.id
	ChunkIdx      int     // 0, 1, 2, ... within the doc
	Text          string  // raw text (what gets embedded)
	FolderContext string  // metadata prefix (returned with search results)
	MetadataJSON  *string // structured metadata as JSON
}

// ChunkRange is a chunk with its line range in the source document (1-indexed).
type ChunkRange struct {
	Text      string
	StartLine int
	EndLine   int
}

// ChunkBySectionWithLines splits content the same way as ChunkBySection but
// also tracks the 1-indexed start/end line of each chunk in the original
// content, so callers can store precise file positions in the index.
func ChunkBySectionWithLines(content string) []ChunkRange {
	texts := chunkTexts(content)
	results := make([]ChunkRange, 0, len(texts))
	searchFrom := 0
	for _, text := range texts {
		idx := strings.Index(content[searchFrom:], text)
		if idx < 0 {
			results = append(results, ChunkRange{Text: text})
			continue
		}
		absStart := searchFrom + idx
		startLine := strings.Count(content[:absStart], "\n") + 1
		endLine := startLine + strings.Count(text, "\n")
		results = append(results, ChunkRange{Text: text, StartLine: startLine, EndLine: endLine})
		searchFrom = absStart + len(text)
	}
	return results
}

// ChunkBySection splits a document by ## section headings (for plans/guides),
// returning chunks that preserve semantic boundaries.
func ChunkBySection(content string) []string {
	chunks := ChunkBySectionWithLines(content)
	texts := make([]string, len(chunks))
	for i, c := range chunks {
		texts[i] = c.Text
	}
	return texts
}

// chunkTexts is the shared splitting logic used by both exported functions.
func chunkTexts(content string) []string {
	var chunks []string
	sections := regexp.MustCompile(`(?m)^## `).Split(content, -1)

	for _, section := range sections {
		if strings.TrimSpace(section) == "" {
			continue
		}
		wordCount := len(strings.Fields(section))
		if wordCount <= 250 {
			chunks = append(chunks, strings.TrimSpace(section))
			continue
		}

		paras := regexp.MustCompile(`\n\n+`).Split(section, -1)
		var current strings.Builder
		currentWords := 0

		for _, para := range paras {
			para = strings.TrimSpace(para)
			if para == "" {
				continue
			}
			paraWords := len(strings.Fields(para))

			if currentWords+paraWords <= 200 {
				if currentWords > 0 {
					current.WriteString("\n\n")
				}
				current.WriteString(para)
				currentWords += paraWords
				continue
			}

			if currentWords > 0 {
				chunks = append(chunks, strings.TrimSpace(current.String()))
			}
			current.Reset()
			current.WriteString(para)
			currentWords = paraWords
		}

		if currentWords > 0 {
			chunks = append(chunks, strings.TrimSpace(current.String()))
		}
	}

	return chunks
}

// DeterministicUUID generates a stable UUID for a kit document based on its path.
// For now, returns the path hash; will use UUID v5 when formalized.
func DeterministicUUID(path string) string {
	// Simplified: use path as-is for now; later: UUID v5(namespace, path)
	// Example: ".hawp/kit/start-here.md" → "kit-start-here"
	base := strings.TrimPrefix(path, ".hawp/")
	base = strings.TrimSuffix(base, ".md")
	base = strings.ReplaceAll(base, "/", "-")
	return base
}

// BuildFolderContext creates the context prefix for a chunk.
func BuildFolderContext(doc Document, metadata *DocumentMetadata) string {
	var sb strings.Builder
	sb.WriteString("# Context: ")
	sb.WriteString(doc.FolderRole)
	sb.WriteString(" (")
	sb.WriteString(doc.Type)
	sb.WriteString(") ")
	sb.WriteString(doc.Path)
	sb.WriteString("\n")

	if metadata != nil {
		sb.WriteString("Status: ")
		sb.WriteString(metadata.Status)
		if metadata.WorkUUID != "" {
			sb.WriteString(" | UUID: ")
			sb.WriteString(metadata.WorkUUID)
		}
		if metadata.ClosedAt != nil {
			sb.WriteString(" | Closed: ")
			sb.WriteString(*metadata.ClosedAt)
		}
		sb.WriteString("\n")
	}
	sb.WriteString("---\n")
	return sb.String()
}
