package embeddings

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// OllamaEmbedder performs embedding via Ollama (local or remote API).
// Supports any model available on the Ollama server: nomic-embed-text, mxbai-embed-large, etc.
type OllamaEmbedder struct {
	client    *http.Client
	baseURL   string // e.g., "http://localhost:11434"
	model     string // e.g., "nomic-embed-text"
	dimension int
}

// ollamaEmbedRequest is the request format for Ollama /api/embeddings endpoint.
// Note: /api/embeddings uses "prompt" (not "input" which is for /api/embed).
type ollamaEmbedRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
}

// ollamaEmbedResponse is the response format for Ollama embedding API.
type ollamaEmbedResponse struct {
	Embedding []float32 `json:"embedding"`
}

// DefaultOllamaURL is the standard Ollama server URL (localhost:11434).
const DefaultOllamaURL = "http://localhost:11434"

// ModelDimensions maps Ollama models to their embedding dimensions.
var ModelDimensions = map[string]int{
	"nomic-embed-text":  768,
	"mxbai-embed-large": 1024,
	"all-minilm":        384,
	"orca-mini":         3072,
	"neural-chat":       4096,
}

// NewOllamaEmbedder creates a new Ollama embedder.
// url: Ollama server URL (e.g., "http://localhost:11434" or "http://ollama-server:11434")
// model: Model name available on Ollama (e.g., "nomic-embed-text")
// Returns error if Ollama server is unreachable or model doesn't exist.
func NewOllamaEmbedder(url, model string) (*OllamaEmbedder, error) {
	if url == "" {
		url = DefaultOllamaURL
	}
	if model == "" {
		return nil, fmt.Errorf("model name required")
	}

	// Normalize URL (remove trailing slash)
	url = strings.TrimRight(url, "/")

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	embedder := &OllamaEmbedder{
		client:  client,
		baseURL: url,
		model:   model,
	}

	// Verify Ollama is reachable and get model info
	dimension, err := embedder.verifyModelAvailable(context.Background())
	if err != nil {
		return nil, fmt.Errorf("verify Ollama model: %w", err)
	}

	embedder.dimension = dimension
	return embedder, nil
}

// verifyModelAvailable checks if model exists and returns its dimension.
func (e *OllamaEmbedder) verifyModelAvailable(ctx context.Context) (int, error) {
	// Try to embed a short test string to verify model is available
	reqBody := ollamaEmbedRequest{
		Model:  e.model,
		Prompt: "test",
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return 0, fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", e.baseURL+"/api/embeddings", strings.NewReader(string(body)))
	if err != nil {
		return 0, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := e.client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("connect to Ollama at %s: %w", e.baseURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return 0, fmt.Errorf("Ollama returned %d: %s", resp.StatusCode, string(respBody))
	}

	var respData ollamaEmbedResponse
	if err := json.NewDecoder(resp.Body).Decode(&respData); err != nil {
		return 0, fmt.Errorf("decode response: %w", err)
	}

	// Return dimension from test embedding
	if len(respData.Embedding) == 0 {
		return 0, fmt.Errorf("no embedding in response")
	}

	return len(respData.Embedding), nil
}

// Embed returns the embedding vector for a single text via Ollama API.
func (e *OllamaEmbedder) Embed(ctx context.Context, text string) ([]float32, error) {
	if text == "" {
		return make([]float32, e.dimension), nil
	}

	embeddings, err := e.EmbedBatch(ctx, []string{text})
	if err != nil {
		return nil, err
	}

	if len(embeddings) > 0 {
		return embeddings[0], nil
	}
	return nil, fmt.Errorf("no embedding returned")
}

// EmbedBatch returns embedding vectors for multiple texts via Ollama API.
func (e *OllamaEmbedder) EmbedBatch(ctx context.Context, texts []string) ([][]float32, error) {
	if len(texts) == 0 {
		return [][]float32{}, nil
	}

	embeddings := make([][]float32, len(texts))

	// Call Ollama for each text (Ollama API doesn't support batch in single call)
	for i, text := range texts {
		reqBody := ollamaEmbedRequest{
			Model:  e.model,
			Prompt: text,
		}

		body, err := json.Marshal(reqBody)
		if err != nil {
			return nil, fmt.Errorf("marshal request for text %d: %w", i, err)
		}

		req, err := http.NewRequestWithContext(ctx, "POST", e.baseURL+"/api/embeddings", strings.NewReader(string(body)))
		if err != nil {
			return nil, fmt.Errorf("create request for text %d: %w", i, err)
		}

		req.Header.Set("Content-Type", "application/json")

		resp, err := e.client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("embed text %d: %w", i, err)
		}

		if resp.StatusCode != http.StatusOK {
			respBody, _ := io.ReadAll(resp.Body)
			resp.Body.Close()
			return nil, fmt.Errorf("Ollama returned %d for text %d: %s", resp.StatusCode, i, string(respBody))
		}

		var respData ollamaEmbedResponse
		if err := json.NewDecoder(resp.Body).Decode(&respData); err != nil {
			resp.Body.Close()
			return nil, fmt.Errorf("decode response for text %d: %w", i, err)
		}
		resp.Body.Close()

		embeddings[i] = respData.Embedding
	}

	return embeddings, nil
}

// Dimension returns the vector dimension for this model.
func (e *OllamaEmbedder) Dimension() int {
	return e.dimension
}

// Backend returns the backend name.
func (e *OllamaEmbedder) Backend() string {
	return "ollama"
}

// Model returns the model name.
func (e *OllamaEmbedder) Model() string {
	return e.model
}

// Close releases HTTP resources.
func (e *OllamaEmbedder) Close() error {
	e.client.CloseIdleConnections()
	return nil
}
