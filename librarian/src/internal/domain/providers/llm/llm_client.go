package llm

import "context"

// LLMClient is the interface all LLM backends implement.
// An LLM client performs text generation (context reshaping) using local or remote models.
type LLMClient interface {
	// Reshape takes packed context and returns reshaped context (improved clarity/structure).
	Reshape(ctx context.Context, packedContext string, maxTokens int) (string, error)

	// ReshapeBatch processes multiple contexts efficiently.
	ReshapeBatch(ctx context.Context, contexts []string, maxTokens int) ([]string, error)

	// Backend returns the backend name for logging and config.
	Backend() string

	// Model returns the specific model name in use.
	Model() string

	// Close releases any resources.
	Close() error
}

// ReshapingPrompt is the standard system prompt for context reshaping.
const ReshapingPrompt = `You are a context reshaping assistant. Your job is to improve the clarity and structure of technical documentation for AI consumption.

Reshape this technical context for optimal readability by:
1. Re-prioritizing information by importance
2. Removing redundancy
3. Structuring hierarchically
4. Highlighting key concepts
5. Improving clarity where needed

Keep the total output under {maxTokens} tokens. Preserve all critical information.`


