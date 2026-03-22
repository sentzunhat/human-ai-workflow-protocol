package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// OllamaLLMClient performs text generation via Ollama (local or remote API).
// Supports any model available on the Ollama server: mistral, llama2, neural-chat, etc.
type OllamaLLMClient struct {
	client    *http.Client
	baseURL   string // e.g., "http://localhost:11434"
	model     string // e.g., "mistral"
	maxTokens int
}

// ollamaGenerateRequest is the request format for Ollama generation API.
type ollamaGenerateRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

// ollamaGenerateResponse is the response format for Ollama generation API.
type ollamaGenerateResponse struct {
	Response string `json:"response"`
	Done     bool   `json:"done"`
}

// DefaultOllamaURL is the standard Ollama server URL (localhost:11434).
const DefaultOllamaLLMURL = "http://localhost:11434"

// DefaultOllamaModel is the recommended LLM model for Ollama.
const DefaultOllamaModel = "mistral"

// NewOllamaLLMClient creates a new Ollama LLM client for text generation.
// url: Ollama server URL (e.g., "http://localhost:11434")
// model: Model name available on Ollama (e.g., "mistral", "llama2", "neural-chat")
// Returns error if Ollama server is unreachable or model doesn't exist.
func NewOllamaLLMClient(url, model string) (*OllamaLLMClient, error) {
	if url == "" {
		url = DefaultOllamaLLMURL
	}
	if model == "" {
		model = DefaultOllamaModel
	}

	// Normalize URL (remove trailing slash)
	url = strings.TrimRight(url, "/")

	// Create HTTP client with timeout — LLM generation on CPU can take 2-5 min for large models.
	client := &http.Client{
		Timeout: 5 * time.Minute,
	}

	llmClient := &OllamaLLMClient{
		client:    client,
		baseURL:   url,
		model:     model,
		maxTokens: 512,
	}

	// Verify Ollama is reachable and model is available
	if err := llmClient.verifyModelAvailable(context.Background()); err != nil {
		return nil, fmt.Errorf("verify Ollama model: %w", err)
	}

	return llmClient, nil
}

// verifyModelAvailable checks if the Ollama server is reachable and the model is listed.
// Uses /api/tags (cheap list call) instead of a full generation to avoid slow startup.
func (c *OllamaLLMClient) verifyModelAvailable(ctx context.Context) error {
	// Use a short-timeout client just for this check
	checkClient := &http.Client{Timeout: 10 * time.Second}

	req, err := http.NewRequestWithContext(ctx, "GET", c.baseURL+"/api/tags", nil)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	resp, err := checkClient.Do(req)
	if err != nil {
		return fmt.Errorf("connect to Ollama at %s: %w", c.baseURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("Ollama returned %d: %s", resp.StatusCode, string(respBody))
	}

	// Parse tag list and confirm the model is present
	var tagsResp struct {
		Models []struct {
			Name string `json:"name"`
		} `json:"models"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&tagsResp); err != nil {
		return fmt.Errorf("decode tags response: %w", err)
	}

	for _, m := range tagsResp.Models {
		// Match on prefix so "qwen3.5:4b" matches "qwen3.5:4b" exactly or short names match
		if m.Name == c.model || strings.HasPrefix(m.Name, c.model+":") || strings.HasPrefix(c.model, strings.Split(m.Name, ":")[0]) {
			return nil
		}
	}

	return fmt.Errorf("model %q not found in Ollama (run: ollama pull %s)", c.model, c.model)
}

// Reshape returns reshaped context using Ollama LLM.
func (c *OllamaLLMClient) Reshape(ctx context.Context, packedContext string, maxTokens int) (string, error) {
	if packedContext == "" {
		return "", nil
	}

	results, err := c.ReshapeBatch(ctx, []string{packedContext}, maxTokens)
	if err != nil {
		return "", err
	}

	if len(results) > 0 {
		return results[0], nil
	}
	return "", fmt.Errorf("no result returned")
}

// ReshapeBatch reshapes multiple contexts using Ollama LLM.
func (c *OllamaLLMClient) ReshapeBatch(ctx context.Context, contexts []string, maxTokens int) ([]string, error) {
	if len(contexts) == 0 {
		return []string{}, nil
	}

	if maxTokens <= 0 {
		maxTokens = c.maxTokens
	}

	results := make([]string, len(contexts))

	// Generate for each context
	for i, context := range contexts {
		// Build prompt for context reshaping
		prompt := strings.ReplaceAll(ReshapingPrompt, "{maxTokens}", fmt.Sprintf("%d", maxTokens))
		fullPrompt := fmt.Sprintf("%s\n\nOriginal context:\n%s\n\nReshaped context:", prompt, context)

		reqBody := ollamaGenerateRequest{
			Model:  c.model,
			Prompt: fullPrompt,
			Stream: false,
		}

		body, err := json.Marshal(reqBody)
		if err != nil {
			return nil, fmt.Errorf("marshal request for context %d: %w", i, err)
		}

		req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/api/generate", strings.NewReader(string(body)))
		if err != nil {
			return nil, fmt.Errorf("create request for context %d: %w", i, err)
		}

		req.Header.Set("Content-Type", "application/json")

		resp, err := c.client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("generate for context %d: %w", i, err)
		}

		if resp.StatusCode != http.StatusOK {
			respBody, _ := io.ReadAll(resp.Body)
			resp.Body.Close()
			return nil, fmt.Errorf("Ollama returned %d for context %d: %s", resp.StatusCode, i, string(respBody))
		}

		var respData ollamaGenerateResponse
		if err := json.NewDecoder(resp.Body).Decode(&respData); err != nil {
			resp.Body.Close()
			return nil, fmt.Errorf("decode response for context %d: %w", i, err)
		}
		resp.Body.Close()

		results[i] = strings.TrimSpace(respData.Response)
	}

	return results, nil
}

// Backend returns the backend name.
func (c *OllamaLLMClient) Backend() string {
	return "ollama"
}

// Model returns the model name.
func (c *OllamaLLMClient) Model() string {
	return c.model
}

// Close releases HTTP resources.
func (c *OllamaLLMClient) Close() error {
	c.client.CloseIdleConnections()
	return nil
}
