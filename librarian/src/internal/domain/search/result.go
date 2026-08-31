package search

// Result represents a single search result with all metadata needed for context packing.
type Result struct {
	// Content
	ChunkID string // Unique identifier for this chunk
	Content string // The actual text content
	Source  string // Document source/path (e.g., "README.md", "guide/setup.md")
	Title   string // Human-readable title or heading

	// Relevance
	Relevance     float32 // Confidence score (0.0 - 1.0)
	LexicalRank   float32 // Relative lexical rank used for hybrid scoring
	SemanticScore float32 // Cosine similarity from embedding search

	// Embedding (for deduplication)
	Embedding []float32 // Vector embedding for similarity comparison

	// Metadata
	LineStart int // Line number where chunk starts in source
	LineEnd   int // Line number where chunk ends
	Priority  int // Ordering hint (lower = higher priority)
}

// Score calculates the hybrid relevance score.
func (r Result) Score() float32 {
	return HybridScore(r.LexicalRank, r.SemanticScore, 0.3, 0.7)
}

// Copy returns a shallow copy of the result.
func (r Result) Copy() Result {
	return Result{
		ChunkID:       r.ChunkID,
		Content:       r.Content,
		Source:        r.Source,
		Title:         r.Title,
		Relevance:     r.Relevance,
		LexicalRank:   r.LexicalRank,
		SemanticScore: r.SemanticScore,
		Embedding:     r.Embedding, // Shallow copy of slice
		LineStart:     r.LineStart,
		LineEnd:       r.LineEnd,
		Priority:      r.Priority,
	}
}
