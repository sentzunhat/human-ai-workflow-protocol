// Package models provides factory functions for constructing embedding and LLM
// backend implementations. Callers depend on the domain interfaces; this package
// selects the correct infrastructure adapter based on the backend name.
package models

import (
	"fmt"

	embeddings "github.com/sentzunhat/hawp/librarian/src/internal/domain/providers/embeddings"
	llm "github.com/sentzunhat/hawp/librarian/src/internal/domain/providers/llm"
	"github.com/sentzunhat/hawp/librarian/src/internal/infrastructure/models/none"
	"github.com/sentzunhat/hawp/librarian/src/internal/infrastructure/models/ollama"
	onnxmodel "github.com/sentzunhat/hawp/librarian/src/internal/infrastructure/models/onnx"
)

// NewEmbedder creates an Embedder for the given backend and model.
// Supports: "onnx", "ollama", "none".
func NewEmbedder(backend, model string) (embeddings.Embedder, error) {
	return NewEmbedderWithURL(backend, model, "")
}

// NewEmbedderWithURL creates an Embedder with an optional URL override.
// url is honored for "ollama" and ignored for "onnx" and "none".
func NewEmbedderWithURL(backend, model, url string) (embeddings.Embedder, error) {
	switch backend {
	case "onnx":
		return onnxmodel.NewONNXEmbedder(model)
	case "ollama":
		return ollama.NewOllamaEmbedder(url, model)
	case "none":
		return none.NewNullEmbedder(), nil
	default:
		return nil, fmt.Errorf("unsupported embedding backend: %s", backend)
	}
}

// NewLLMClient creates a LLMClient for the given backend and model.
// Supports: "onnx", "ollama", "none".
func NewLLMClient(backend, model string) (llm.LLMClient, error) {
	return NewLLMClientWithURL(backend, model, "")
}

// NewLLMClientWithURL creates a LLMClient with an optional URL override.
// url is honored for "ollama" and ignored for "onnx" and "none".
func NewLLMClientWithURL(backend, model, url string) (llm.LLMClient, error) {
	switch backend {
	case "onnx":
		return onnxmodel.NewONNXLLMClient(model)
	case "ollama":
		return ollama.NewOllamaLLMClient(url, model)
	case "none":
		return none.NewNullLLMClient(), nil
	default:
		return nil, fmt.Errorf("unsupported LLM backend: %s", backend)
	}
}
