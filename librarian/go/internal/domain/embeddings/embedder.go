package embeddings

import (
	"context"
	"fmt"
)

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

// EmbeddingResult holds a text and its embedding vector.
type EmbeddingResult struct {
	Text      string
	Embedding []float32
	Tokens    int
}

// NewEmbedder creates an embedder based on the backend name and model.
// Supports: "onnx" (local), "ollama" (local/remote API), "none" (zero-cost
// no-op, see NullEmbedder).
// Phase 3.2c/d will add OpenAI, Anthropic backends.
func NewEmbedder(backend, model string) (Embedder, error) {
	switch backend {
	case "onnx":
		return NewONNXEmbedder(model)
	case "ollama":
		return NewOllamaEmbedder("", model) // Use default Ollama URL
	case "none":
		return NewNullEmbedder(), nil
	default:
		return nil, fmt.Errorf("unsupported embedding backend: %s", backend)
	}
}

// NewEmbedderWithURL creates an embedder with a custom URL for remote backends.
// Used for Ollama, OpenAI, Anthropic where URL may be non-default. url is
// ignored for "onnx" and "none".
func NewEmbedderWithURL(backend, model, url string) (Embedder, error) {
	switch backend {
	case "onnx":
		return NewONNXEmbedder(model)
	case "ollama":
		return NewOllamaEmbedder(url, model)
	case "none":
		return NewNullEmbedder(), nil
	default:
		return nil, fmt.Errorf("unsupported embedding backend: %s", backend)
	}
}
