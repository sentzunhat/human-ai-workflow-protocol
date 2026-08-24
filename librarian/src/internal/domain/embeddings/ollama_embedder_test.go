package embeddings

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestOllamaEmbedRequestUsesPromptField locks in the request contract with
// Ollama's /api/embeddings endpoint: it expects a "prompt" field, not
// "input" (that field belongs to /api/embed instead). A prior regression
// used the wrong field name and broke live embedding calls while unit tests
// still passed, because nothing asserted on request shape. See
// librarian/docs/code-review-v003.md Finding 2.
func TestOllamaEmbedRequestUsesPromptField(t *testing.T) {
	var gotBody map[string]interface{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&gotBody); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ollamaEmbedResponse{Embedding: make([]float32, 384)})
	}))
	defer server.Close()

	embedder, err := NewOllamaEmbedder(server.URL, "all-minilm")
	if err != nil {
		t.Fatalf("failed to create embedder: %v", err)
	}
	defer embedder.Close()

	if _, ok := gotBody["prompt"]; !ok {
		t.Errorf("request body should have a %q field, got keys: %v", "prompt", gotBody)
	}
	if _, ok := gotBody["input"]; ok {
		t.Errorf("request body should not have an %q field (that's /api/embed, not /api/embeddings)", "input")
	}
}

func TestNewOllamaEmbedder(t *testing.T) {
	// Create a mock Ollama server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" || r.URL.Path != "/api/embeddings" {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}

		var req ollamaEmbedRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		resp := ollamaEmbedResponse{
			Embedding: make([]float32, 384), // Mock embedding
		}
		for i := range resp.Embedding {
			resp.Embedding[i] = float32(i) / 384.0
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	embedder, err := NewOllamaEmbedder(server.URL, "all-minilm")
	if err != nil {
		t.Fatalf("NewOllamaEmbedder failed: %v", err)
	}
	defer embedder.Close()

	if embedder.Model() != "all-minilm" {
		t.Errorf("model should be all-minilm, got %s", embedder.Model())
	}
	if embedder.Backend() != "ollama" {
		t.Errorf("backend should be ollama, got %s", embedder.Backend())
	}
	if embedder.Dimension() != 384 {
		t.Errorf("dimension should be 384, got %d", embedder.Dimension())
	}
}

func TestOllamaEmbedderInterface(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := ollamaEmbedResponse{
			Embedding: make([]float32, 768),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	embedder, err := NewOllamaEmbedder(server.URL, "nomic-embed-text")
	if err != nil {
		t.Fatalf("failed to create embedder: %v", err)
	}
	defer embedder.Close()

	// Test that embedder implements interface
	var _ Embedder = embedder
}

func TestOllamaUnreachable(t *testing.T) {
	_, err := NewOllamaEmbedder("http://localhost:9999", "some-model")
	if err == nil {
		t.Error("should fail when Ollama is unreachable")
	}
}

func TestOllamaEmptyText(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := ollamaEmbedResponse{
			Embedding: make([]float32, 384),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	embedder, err := NewOllamaEmbedder(server.URL, "all-minilm")
	if err != nil {
		t.Fatalf("failed to create embedder: %v", err)
	}
	defer embedder.Close()

	// Empty text should return zero vector
	vec, err := embedder.Embed(context.Background(), "")
	if err != nil {
		t.Fatalf("embedding empty text failed: %v", err)
	}

	if len(vec) != embedder.Dimension() {
		t.Errorf("zero vector should have length %d, got %d", embedder.Dimension(), len(vec))
	}

	for _, v := range vec {
		if v != 0 {
			t.Errorf("zero vector should have all zeros, got %f", v)
		}
	}
}

func TestOllamaEmbedSingleText(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req ollamaEmbedRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		// Mock: return embedding based on text length
		embedding := make([]float32, 384)
		for i := 0; i < len(embedding); i++ {
			embedding[i] = float32(len(req.Prompt)) / float32(i+1)
		}

		resp := ollamaEmbedResponse{Embedding: embedding}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	embedder, err := NewOllamaEmbedder(server.URL, "all-minilm")
	if err != nil {
		t.Fatalf("failed to create embedder: %v", err)
	}
	defer embedder.Close()

	text := "This is a test sentence for Ollama embedding."
	vec, err := embedder.Embed(context.Background(), text)
	if err != nil {
		t.Fatalf("embedding failed: %v", err)
	}

	if len(vec) != embedder.Dimension() {
		t.Errorf("embedding dimension mismatch: expected %d, got %d", embedder.Dimension(), len(vec))
	}

	// Check that vector has non-zero values
	hasNonZero := false
	for _, v := range vec {
		if v != 0 {
			hasNonZero = true
			break
		}
	}
	if !hasNonZero {
		t.Error("embedding vector should not be all zeros")
	}
}

func TestOllamaEmbedBatch(t *testing.T) {
	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		if r.Method != "POST" {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req ollamaEmbedRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		embedding := make([]float32, 768)
		for i := 0; i < len(embedding); i++ {
			embedding[i] = float32(callCount) / float32(i+1)
		}

		resp := ollamaEmbedResponse{Embedding: embedding}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	embedder, err := NewOllamaEmbedder(server.URL, "nomic-embed-text")
	if err != nil {
		t.Fatalf("failed to create embedder: %v", err)
	}
	defer embedder.Close()

	texts := []string{
		"First text about embeddings.",
		"Second text about Ollama.",
		"Third text about API calls.",
	}

	embeddings, err := embedder.EmbedBatch(context.Background(), texts)
	if err != nil {
		t.Fatalf("batch embedding failed: %v", err)
	}

	if len(embeddings) != len(texts) {
		t.Errorf("batch size mismatch: expected %d, got %d", len(texts), len(embeddings))
	}

	// Verify each embedding has correct dimension
	for i, vec := range embeddings {
		if len(vec) != embedder.Dimension() {
			t.Errorf("embedding[%d] dimension mismatch: expected %d, got %d",
				i, embedder.Dimension(), len(vec))
		}
	}

	// Verify API was called for each text
	if callCount != len(texts)+1 { // +1 for verification call
		t.Errorf("expected %d API calls, got %d", len(texts)+1, callCount)
	}
}

func TestOllamaEmptyBatch(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Return valid embedding for verification, but error for other calls
		resp := ollamaEmbedResponse{
			Embedding: make([]float32, 384),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	embedder, err := NewOllamaEmbedder(server.URL, "all-minilm")
	if err != nil {
		t.Fatalf("failed to create embedder: %v", err)
	}
	defer embedder.Close()

	embeddings, err := embedder.EmbedBatch(context.Background(), []string{})
	if err != nil {
		t.Fatalf("empty batch should not error: %v", err)
	}

	if len(embeddings) != 0 {
		t.Errorf("empty batch should return empty slice, got %d", len(embeddings))
	}
}

func TestOllamaServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "model not found", http.StatusNotFound)
	}))
	defer server.Close()

	_, err := NewOllamaEmbedder(server.URL, "nonexistent-model")
	if err == nil {
		t.Error("should fail when model not found")
	}
}

func TestOllamaContext(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := ollamaEmbedResponse{
			Embedding: make([]float32, 384),
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	embedder, err := NewOllamaEmbedder(server.URL, "all-minilm")
	if err != nil {
		t.Fatalf("failed to create embedder: %v", err)
	}
	defer embedder.Close()

	// Test with cancelled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err = embedder.Embed(ctx, "test")
	if err == nil {
		t.Error("should fail with cancelled context")
	}
}

func TestModelDimensions(t *testing.T) {
	if len(ModelDimensions) == 0 {
		t.Error("ModelDimensions should not be empty")
	}

	expected := map[string]int{
		"nomic-embed-text":  768,
		"mxbai-embed-large": 1024,
		"all-minilm":        384,
	}

	for model, dim := range expected {
		if ModelDimensions[model] != dim {
			t.Errorf("ModelDimensions[%s] should be %d, got %d", model, dim, ModelDimensions[model])
		}
	}
}
