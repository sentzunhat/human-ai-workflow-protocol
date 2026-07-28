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

// ChunkBySection splits a document by ## section headings (for plans/guides),
// returning chunks that preserve semantic boundaries.
func ChunkBySection(content string) []string {
	// Split by ## headings; for each section, chunk by paragraphs if too long.
	// Target: 150-200 words per chunk, but respect section boundaries.

	var chunks []string
	sections := regexp.MustCompile(`(?m)^## `).Split(content, -1)

	for _, section := range sections {
		if strings.TrimSpace(section) == "" {
			continue
		}
		// If the section is short, keep it whole
		wordCount := len(strings.Fields(section))
		if wordCount <= 250 {
			chunks = append(chunks, strings.TrimSpace(section))
			continue
		}

		// Large section: split by paragraphs (blank-line delimited)
		paras := regexp.MustCompile(`\n\n+`).Split(section, -1)
		var current strings.Builder
		currentWords := 0

		for _, para := range paras {
			para = strings.TrimSpace(para)
			if para == "" {
				continue
			}
			paraWords := len(strings.Fields(para))

			// If current chunk + this paragraph fits, add it
			if currentWords+paraWords <= 200 {
				if currentWords > 0 {
					current.WriteString("\n\n")
				}
				current.WriteString(para)
				currentWords += paraWords
				continue
			}

			// Otherwise, flush current and start fresh
			if currentWords > 0 {
				chunks = append(chunks, strings.TrimSpace(current.String()))
			}
			current.Reset()
			current.WriteString(para)
			currentWords = paraWords
		}

		// Flush any remaining
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
