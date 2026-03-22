package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

// tagsHandler writes the /api/tags response NewOllamaLLMClient's
// verifyModelAvailable checks against, listing exactly one model.
func tagsHandler(w http.ResponseWriter, model string) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"models":[{"name":%q}]}`, model+":latest")
}

func TestNewOllamaLLMClient(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" && r.URL.Path == "/api/tags" {
			tagsHandler(w, "mistral")
			return
		}
		if r.Method != "POST" || r.URL.Path != "/api/generate" {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}

		var req ollamaGenerateRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		resp := ollamaGenerateResponse{
			Response: "Generated text response.",
			Done:     true,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client, err := NewOllamaLLMClient(server.URL, "mistral")
	if err != nil {
		t.Fatalf("NewOllamaLLMClient failed: %v", err)
	}
	defer client.Close()

	if client.Model() != "mistral" {
		t.Errorf("model should be mistral, got %s", client.Model())
	}
	if client.Backend() != "ollama" {
		t.Errorf("backend should be ollama, got %s", client.Backend())
	}
}

func TestOllamaLLMClientInterface(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" && r.URL.Path == "/api/tags" {
			tagsHandler(w, "llama2")
			return
		}
		resp := ollamaGenerateResponse{
			Response: "Test response",
			Done:     true,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client, err := NewOllamaLLMClient(server.URL, "llama2")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	defer client.Close()

	// Test that client implements interface
	var _ LLMClient = client
}

func TestOllamaUnreachable(t *testing.T) {
	_, err := NewOllamaLLMClient("http://localhost:9999", "mistral")
	if err == nil {
		t.Error("should fail when Ollama is unreachable")
	}
}

func TestOllamaEmptyContext(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" && r.URL.Path == "/api/tags" {
			tagsHandler(w, "neural-chat")
			return
		}
		resp := ollamaGenerateResponse{
			Response: "Response",
			Done:     true,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client, err := NewOllamaLLMClient(server.URL, "neural-chat")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	defer client.Close()

	// Empty context should return empty string
	result, err := client.Reshape(context.Background(), "", 100)
	if err != nil {
		t.Fatalf("reshaping empty context failed: %v", err)
	}

	if result != "" {
		t.Errorf("empty context should return empty string, got %q", result)
	}
}

func TestOllamaReshapeSingleContext(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" && r.URL.Path == "/api/tags" {
			tagsHandler(w, "mistral")
			return
		}
		if r.Method != "POST" {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req ollamaGenerateRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		// Mock: return generated text
		generated := "This is a reshaped version of the input context with better structure and clarity."
		resp := ollamaGenerateResponse{
			Response: generated,
			Done:     true,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client, err := NewOllamaLLMClient(server.URL, "mistral")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	defer client.Close()

	ctx := context.Background()
	inputContext := "This is the original context that needs to be reshaped for better clarity."
	result, err := client.Reshape(ctx, inputContext, 200)
	if err != nil {
		t.Fatalf("reshape failed: %v", err)
	}

	if result == "" {
		t.Error("reshaped context should not be empty")
	}

	// Result should contain some text
	if len(result) < 10 {
		t.Errorf("reshaped context too short: %s", result)
	}
}

func TestOllamaReshapeBatch(t *testing.T) {
	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" && r.URL.Path == "/api/tags" {
			tagsHandler(w, "llama2")
			return
		}
		callCount++
		if r.Method != "POST" {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req ollamaGenerateRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		// Mock response based on call count
		generated := "Reshaped response " + string(rune(callCount))
		resp := ollamaGenerateResponse{
			Response: generated,
			Done:     true,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client, err := NewOllamaLLMClient(server.URL, "llama2")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	defer client.Close()

	ctx := context.Background()
	contexts := []string{
		"First context for reshaping.",
		"Second context with more details.",
		"Third context for final test.",
	}

	results, err := client.ReshapeBatch(ctx, contexts, 100)
	if err != nil {
		t.Fatalf("batch reshape failed: %v", err)
	}

	if len(results) != len(contexts) {
		t.Errorf("batch size mismatch: expected %d, got %d", len(contexts), len(results))
	}

	// Verify each result has content
	for i, result := range results {
		if result == "" {
			t.Errorf("result[%d] should not be empty", i)
		}
	}

	// Verify /api/generate was called once per context (verification uses /api/tags, not counted here)
	if callCount != len(contexts) {
		t.Errorf("expected %d generate calls, got %d", len(contexts), callCount)
	}
}

func TestOllamaEmptyBatch(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" && r.URL.Path == "/api/tags" {
			tagsHandler(w, "mistral")
			return
		}
		resp := ollamaGenerateResponse{
			Response: "Should not be called",
			Done:     true,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client, err := NewOllamaLLMClient(server.URL, "mistral")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	defer client.Close()

	results, err := client.ReshapeBatch(context.Background(), []string{}, 100)
	if err != nil {
		t.Fatalf("empty batch should not error: %v", err)
	}

	if len(results) != 0 {
		t.Errorf("empty batch should return empty slice, got %d", len(results))
	}
}

func TestOllamaServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "model not found", http.StatusNotFound)
	}))
	defer server.Close()

	_, err := NewOllamaLLMClient(server.URL, "nonexistent-model")
	if err == nil {
		t.Error("should fail when model not found")
	}
}

func TestOllamaContext(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" && r.URL.Path == "/api/tags" {
			tagsHandler(w, "mistral")
			return
		}
		resp := ollamaGenerateResponse{
			Response: "Generated response",
			Done:     true,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client, err := NewOllamaLLMClient(server.URL, "mistral")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	defer client.Close()

	// Test with cancelled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err = client.Reshape(ctx, "test", 100)
	if err == nil {
		t.Error("should fail with cancelled context")
	}
}

func TestOllamaDefaultModel(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" && r.URL.Path == "/api/tags" {
			tagsHandler(w, DefaultOllamaModel)
			return
		}

		var req ollamaGenerateRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		// Verify default model is used
		if req.Model != DefaultOllamaModel {
			http.Error(w, "wrong model", http.StatusBadRequest)
			return
		}

		resp := ollamaGenerateResponse{
			Response: "Response",
			Done:     true,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	// Create client with empty model string (should use default)
	client, err := NewOllamaLLMClient(server.URL, "")
	if err != nil {
		t.Fatalf("failed to create client with default model: %v", err)
	}
	defer client.Close()

	if client.Model() != DefaultOllamaModel {
		t.Errorf("should use default model %s, got %s", DefaultOllamaModel, client.Model())
	}
}
