package none

import "context"

// NullEmbedder is a zero-cost, no-network Embedder for backend "none". It
// performs no computation and requires no model or server.
type NullEmbedder struct{}

// NewNullEmbedder creates a NullEmbedder. Never fails.
func NewNullEmbedder() *NullEmbedder { return &NullEmbedder{} }

func (n *NullEmbedder) Embed(_ context.Context, _ string) ([]float32, error) {
	return []float32{}, nil
}

func (n *NullEmbedder) EmbedBatch(_ context.Context, texts []string) ([][]float32, error) {
	out := make([][]float32, len(texts))
	for i := range out {
		out[i] = []float32{}
	}
	return out, nil
}

func (n *NullEmbedder) Dimension() int  { return 0 }
func (n *NullEmbedder) Backend() string { return "none" }
func (n *NullEmbedder) Model() string   { return "none" }
func (n *NullEmbedder) Close() error    { return nil }
