package embeddings

import "context"

// Embedder is the interface all embedding backends implement.
// An embedder converts text to numerical vectors for semantic search.
type Embedder interface {
	// Embed returns a vector for a single text.
	Embed(ctx context.Context, text string) ([]float32, error)

	// EmbedBatch returns vectors for multiple texts (more efficient).
	EmbedBatch(ctx context.Context, texts []string) ([][]float32, error)

	// Dimension returns the vector dimension (e.g., 768 for BGE, 384 for MiniLM).
	Dimension() int

	// Backend returns the backend name for logging and config.
	Backend() string

	// Model returns the specific model name in use.
	Model() string

	// Close releases any resources (e.g., ONNX session).
	Close() error
}

// DefaultModel is the recommended embedding model (best quality).
// Matches onnx.DefaultEmbedModel — kept here so application packages that
// reference embeddings.DefaultModel do not need to import infrastructure.
const DefaultModel = "bge-base-en-v1.5"

// EmbeddingResult holds a text and its embedding vector.
type EmbeddingResult struct {
	Text      string
	Embedding []float32
	Tokens    int
}
