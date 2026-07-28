package llm

import "context"

// NullLLMClient is a zero-cost, no-network LLMClient for backend "none". It
// performs no generation and requires no model or server — Reshape and
// ReshapeBatch return their input unchanged. Used when a caller wants
// structured output (search results + references) without requiring any
// LLM at all: no Ollama server, no ONNX model, no API key. This is a
// deliberate mode, not an error fallback — see ContextReshaper.Reshape,
// which skips calling into the LLM entirely when the configured backend is
// "none" rather than routing a no-op prompt through it.
type NullLLMClient struct{}

// NewNullLLMClient creates a NullLLMClient. Never fails — no model or
// network dependency to fail on.
func NewNullLLMClient() *NullLLMClient {
	return &NullLLMClient{}
}

// Reshape returns content unchanged.
func (n *NullLLMClient) Reshape(ctx context.Context, content string, maxTokens int) (string, error) {
	return content, nil
}

// ReshapeBatch returns contents unchanged.
func (n *NullLLMClient) ReshapeBatch(ctx context.Context, contents []string, maxTokens int) ([]string, error) {
	return contents, nil
}

// Backend returns "none".
func (n *NullLLMClient) Backend() string { return "none" }

// Model returns "none".
func (n *NullLLMClient) Model() string { return "none" }

// Close is a no-op — there are no resources to release.
func (n *NullLLMClient) Close() error { return nil }
