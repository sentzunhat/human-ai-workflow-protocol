package embeddings

import "context"

// NullEmbedder is a zero-cost, no-network Embedder for backend "none". It
// performs no computation and requires no model or server — used when a
// caller wants structured search/reference output (deduplicated results,
// source references) without running any embedding model at all, e.g. a
// user with no ONNX model downloaded and no Ollama server running who still
// wants --context output.
type NullEmbedder struct{}

// NewNullEmbedder creates a NullEmbedder. Never fails — no model or network
// dependency to fail on.
func NewNullEmbedder() *NullEmbedder {
	return &NullEmbedder{}
}

// Embed returns an empty vector; there is no model to run.
func (n *NullEmbedder) Embed(ctx context.Context, text string) ([]float32, error) {
	return []float32{}, nil
}

// EmbedBatch returns one empty vector per input text (preserving length, so
// callers that check len(embeddings) == len(texts) don't error).
func (n *NullEmbedder) EmbedBatch(ctx context.Context, texts []string) ([][]float32, error) {
	out := make([][]float32, len(texts))
	for i := range out {
		out[i] = []float32{}
	}
	return out, nil
}

// Dimension returns 0 — there is no vector space.
func (n *NullEmbedder) Dimension() int { return 0 }

// Backend returns "none".
func (n *NullEmbedder) Backend() string { return "none" }

// Model returns "none".
func (n *NullEmbedder) Model() string { return "none" }

// Close is a no-op — there are no resources to release.
func (n *NullEmbedder) Close() error { return nil }
