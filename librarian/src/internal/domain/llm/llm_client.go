package llm

import (
	"context"
	"fmt"
)

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

// ModelInfo holds metadata for supported LLM models.
type ModelInfo struct {
	Name             string
	HFRepo           string // Hugging Face repo for download
	ONNXFile         string // ONNX file to use in the repo (passed as opts.OnnxFilePath)
	ExternalDataFile string // External weights sidecar (opts.ExternalDataPath), "" if none
	Params           int64  // Number of parameters (for reference)
}

// SupportedModels lists all available LLM models with ONNX export.
// SmolLM2-360M-Instruct was chosen over FLAN-T5-small after benchmarking
// (2026-07-26): SmolLM2 is a modern (2026), instruction-tuned generative
// model trained on 4T tokens — meaningfully higher quality than FLAN-T5's
// older T5-based architecture at a similar CPU latency (measured ~120ms per
// 50-token extraction), for ~10x the download size (819MB vs 60-80MB).
// See librarian/docs/BENCHMARKS_v003.md and v0.1.0_VISION.md.
// NOTE: hugot's ORT generative pipeline requires a genai_config.json —
// the ONNX Runtime GenAI model-builder format, NOT the standard
// optimum/Transformers.js ONNX export (HuggingFaceTB/SmolLM2-360M-Instruct
// itself lacks genai_config.json and does not work here). This repo is a
// community genai-format conversion of the same model, int4-quantized.
var SupportedModels = map[string]ModelInfo{
	"SmolLM2-360M-Instruct": {
		Name:             "SmolLM2-360M-Instruct",
		HFRepo:           "homen3/SmolLM2-360M-Instruct-ort-genai-int4-cpu",
		ONNXFile:         "model.onnx",
		ExternalDataFile: "model.onnx.data", // model weights split into a sidecar file; without downloading this too, ORT fails at session-init with "External data path does not exist"
		Params:           360_000_000,
	},
}

// DefaultModel is the recommended ONNX LLM model for context reshaping.
const DefaultModel = "SmolLM2-360M-Instruct"

// NewLLMClient creates an LLM client based on the backend name and model.
// Supports: "onnx" (local), "ollama" (local/remote API), "none" (zero-cost
// passthrough, see NullLLMClient).
// Phase 3.3c/d will add OpenAI, Anthropic backends.
func NewLLMClient(backend, model string) (LLMClient, error) {
	switch backend {
	case "onnx":
		return NewONNXLLMClient(model)
	case "ollama":
		return NewOllamaLLMClient("", model) // Use default Ollama URL
	case "none":
		return NewNullLLMClient(), nil
	default:
		return nil, fmt.Errorf("unsupported LLM backend: %s", backend)
	}
}

// NewLLMClientWithURL creates an LLM client with a custom URL for remote backends.
// Used for Ollama, OpenAI, Anthropic where URL may be non-default. url is
// ignored for "onnx" and "none".
func NewLLMClientWithURL(backend, model, url string) (LLMClient, error) {
	switch backend {
	case "onnx":
		return NewONNXLLMClient(model)
	case "ollama":
		return NewOllamaLLMClient(url, model)
	case "none":
		return NewNullLLMClient(), nil
	default:
		return nil, fmt.Errorf("unsupported LLM backend: %s", backend)
	}
}
