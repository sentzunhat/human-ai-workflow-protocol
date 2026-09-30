package none

import "context"

// NullLLMClient is a zero-cost, no-network LLMClient for backend "none".
// Reshape and ReshapeBatch return their input unchanged.
type NullLLMClient struct{}

// NewNullLLMClient creates a NullLLMClient. Never fails.
func NewNullLLMClient() *NullLLMClient { return &NullLLMClient{} }

func (n *NullLLMClient) Reshape(_ context.Context, content string, _ int) (string, error) {
	return content, nil
}

func (n *NullLLMClient) ReshapeBatch(_ context.Context, contents []string, _ int) ([]string, error) {
	return contents, nil
}

func (n *NullLLMClient) Backend() string { return "none" }
func (n *NullLLMClient) Model() string   { return "none" }
func (n *NullLLMClient) Close() error    { return nil }
