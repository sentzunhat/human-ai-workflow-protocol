package llm

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync/atomic"

	"github.com/knights-analytics/hugot"
	"github.com/knights-analytics/hugot/backends"
	"github.com/knights-analytics/hugot/options"
	"github.com/knights-analytics/hugot/pipelines"
)

// ONNXLLMClient performs text generation via ONNX Runtime (local, fast, private).
// Supports SmolLM2-360M-Instruct (see SupportedModels).
//
// Generative models require hugot's CGO ORT backend — its portable Go/XLA
// backend (what HAWP's embeddings use) cannot run them at all, by hugot's
// own design (backends/model.go: `if options.Backend != "ORT" { return error }`).
// A default `go build` (no `-tags ORT`) makes NewORTSession itself return a
// clear "to enable ORT, run `go build -tags ORT`" error immediately, before
// any download — see hugot's hugot_ort_disabled.go. Built with `-tags ORT`,
// it additionally needs three native libraries at link/run time:
// libonnxruntime, libonnxruntime-genai, and libtokenizers (static). See
// onnxRuntimeLibraryPath's doc comment for how HAWP locates them.
type ONNXLLMClient struct {
	model     string
	modelPath string
	session   *hugot.Session
	maxTokens int
}

// pipelineCounter ensures unique pipeline names across client instances
var pipelineCounter atomic.Int64

// onnxRuntimeLibraryPath resolves the DIRECTORY containing the plain ONNX
// Runtime shared library (libonnxruntime.{dylib,so}, or onnxruntime.dll on
// Windows) hugot's ORT session needs at init time.
// hugot.options.WithOnnxLibraryPath takes a directory, not the library file
// itself — it joins the directory with the platform-default library
// filename internally (confirmed 2026-07-27: passing the file path directly
// fails with "<path> is not a directory").
// HAWP does not bundle this library today (see v0.1.0_VISION.md's ONNX LLM
// section) — for now it looks in ~/.hawp/native/lib/ (matching the manual
// local setup used to validate ORT support), overridable via
// HAWP_ONNX_LIBRARY_DIR for anyone testing a different location. Returns ""
// if not found, so callers fall back to hugot's own platform default search.
func onnxRuntimeLibraryPath() string {
	if p := os.Getenv("HAWP_ONNX_LIBRARY_DIR"); p != "" {
		return p
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}

	name := "libonnxruntime.so"
	switch runtime.GOOS {
	case "darwin":
		name = "libonnxruntime.dylib"
	case "windows":
		name = "onnxruntime.dll"
	}

	dir := filepath.Join(home, ".hawp", "native", "lib")
	if _, err := os.Stat(filepath.Join(dir, name)); err != nil {
		return ""
	}
	return dir
}

// NewONNXLLMClient creates a new ONNX LLM client for the given model.
// Model should be "SmolLM2-360M-Instruct" (default) or another entry in
// SupportedModels. On first use, automatically downloads the model to
// ~/.hawp/models/llm/.
//
// Requires a `-tags ORT` build with the native libraries described in this
// file's package doc comment; without them, this fails immediately with a
// clear error (from hugot itself, or from the ORT library-path check below)
// rather than downloading a model that can't run.
func NewONNXLLMClient(model string) (*ONNXLLMClient, error) {
	if model == "" {
		model = DefaultModel
	}

	// Validate model name
	info, ok := SupportedModels[model]
	if !ok {
		return nil, fmt.Errorf("unsupported ONNX model: %s; supported: %v", model, supportedModelNames())
	}

	// Create the ORT session FIRST, before touching the network. Without a
	// `-tags ORT` build this call itself returns hugot's own clear "to
	// enable ORT, run `go build -tags ORT`" error, instantly and for free —
	// so selecting the "onnx" LLM backend on a default build never silently
	// downloads an ~1GB+ model that can't actually run. Only once this
	// succeeds (a real ORT build with the native libraries present) do we
	// proceed to download.
	ctx := context.Background()
	var sessionOpts []options.WithOption
	if libPath := onnxRuntimeLibraryPath(); libPath != "" {
		sessionOpts = append(sessionOpts, options.WithOnnxLibraryPath(libPath))
	}
	session, err := hugot.NewORTSession(ctx, sessionOpts...)
	if err != nil {
		return nil, fmt.Errorf("start hugot ORT session: %w", err)
	}

	// Resolve model directory (~/.hawp/models/llm/{model}/)
	home, err := os.UserHomeDir()
	if err != nil {
		session.Destroy()
		return nil, fmt.Errorf("get home directory: %w", err)
	}

	modelsDir := filepath.Join(home, ".hawp", "models", "llm")
	if err := os.MkdirAll(modelsDir, 0755); err != nil {
		session.Destroy()
		return nil, fmt.Errorf("create models directory: %w", err)
	}

	modelPath := filepath.Join(modelsDir, model)

	// Download model if not already present
	if _, err := os.Stat(modelPath); os.IsNotExist(err) {
		opts := hugot.NewDownloadOptions()
		opts.OnnxFilePath = info.ONNXFile
		opts.ExternalDataPath = info.ExternalDataFile
		modelPath, err = hugot.DownloadModel(context.Background(), info.HFRepo, modelsDir, opts)
		if err != nil {
			session.Destroy()
			return nil, fmt.Errorf("download model %s from %s: %w", model, info.HFRepo, err)
		}
	}

	client := &ONNXLLMClient{
		model:     model,
		modelPath: modelPath,
		session:   session,
		maxTokens: 512,
	}

	return client, nil
}

// supportedModelNames returns a slice of supported model names.
func supportedModelNames() []string {
	names := make([]string, 0, len(SupportedModels))
	for name := range SupportedModels {
		names = append(names, name)
	}
	return names
}

// estimateTokenCount roughly estimates token count using the same ~4
// chars-per-token heuristic used elsewhere in HAWP (e.g.
// application/context/format.go's estimateTokens). It's deliberately an
// overestimate-friendly approximation, not a real tokenizer call — used
// here only to size a generation budget, where erring high costs nothing
// but erring low causes a hard failure.
func estimateTokenCount(text string) int {
	return (len(text) + 3) / 4
}

// Reshape returns reshaped context using ONNX LLM inference.
func (c *ONNXLLMClient) Reshape(ctx context.Context, packedContext string, maxTokens int) (string, error) {
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

// ReshapeBatch reshapes multiple contexts.
//
// This runs one single-prompt inference call per context rather than a
// single batched multi-prompt call. That used to be batched (one
// RunPipeline call for the whole slice), which is the more efficient shape
// on paper, but real multi-context batches could fail with "generation
// stopped: max length reached" even under a very generous shared budget
// (tested up to 2x prompt+2x maxTokens) — root-caused 2026-07-27 to
// onnxruntime-genai's batch generation not honoring per-sequence EOS
// independently within a batch, so a longer sequence in the batch can force
// every sequence to keep generating past its own natural stop. Looping
// single-prompt calls avoids that failure mode entirely, since each call
// gets its own EOS handling — confirmed reliable (what Reshape() and the
// RAG pipeline already used, at the tighter single-prompt budget below).
// See v0.1.0_VISION.md for the batched approach this replaced.
func (c *ONNXLLMClient) ReshapeBatch(ctx context.Context, contexts []string, maxTokens int) ([]string, error) {
	if len(contexts) == 0 {
		return []string{}, nil
	}

	if maxTokens <= 0 {
		maxTokens = c.maxTokens
	}

	results := make([]string, len(contexts))
	for i, text := range contexts {
		result, err := c.reshapeOne(ctx, text, maxTokens)
		if err != nil {
			return nil, fmt.Errorf("reshape context %d/%d: %w", i+1, len(contexts), err)
		}
		results[i] = result
	}

	return results, nil
}

// reshapeOne runs a single ChatML-wrapped prompt through one inference call.
func (c *ONNXLLMClient) reshapeOne(ctx context.Context, text string, maxTokens int) (string, error) {
	// Wrapped in ChatML — the format SmolLM2-Instruct's own
	// chat_template.jinja expects (<|im_start|>role\ncontent<|im_end|>).
	// Confirmed 2026-07-27: a plain text prompt (no special tokens) makes
	// the model treat the whole thing as a document to continue, and it
	// just echoes the input back character-for-character rather than
	// following the instruction.
	instruction := strings.ReplaceAll(ReshapingPrompt, "{maxTokens}", fmt.Sprintf("%d", maxTokens))
	prompt := fmt.Sprintf(
		"<|im_start|>system\n%s<|im_end|>\n<|im_start|>user\n%s<|im_end|>\n<|im_start|>assistant\n",
		instruction, text)

	// For the ORT generative backend, MaxLength is the TOTAL sequence length
	// (prompt + generation), not just newly generated tokens — confirmed
	// 2026-07-27: passing just maxTokens fails with "input_ids size (143) +
	// current sequence length (0) exceeds max length (60)" as soon as the
	// prompt alone exceeds maxTokens. The ChatML special tokens
	// (<|im_start|>, <|im_end|>, role names) tokenize denser than
	// estimateTokenCount's ~4-chars-per-token prose heuristic assumes, so
	// this pads generously (25%) rather than tightly — erring high only
	// costs a slightly larger generation buffer, erring low is a hard
	// failure ("max length reached" mid-generation).
	promptTokens := estimateTokenCount(prompt)
	totalMaxLength := promptTokens + promptTokens/4 + maxTokens + 32

	pipelineID := pipelineCounter.Add(1)
	config := hugot.TextGenerationConfig{
		ModelPath: c.modelPath,
		Name:      fmt.Sprintf("hawp-llm-%s-%d", c.model, pipelineID),
		Options: []backends.PipelineOption[*pipelines.TextGenerationPipeline]{
			pipelines.WithMaxLength(totalMaxLength),
		},
	}

	pipeline, err := hugot.NewPipeline(c.session, config)
	if err != nil {
		return "", fmt.Errorf("load LLM model at %s: %w", c.modelPath, err)
	}

	output, err := pipeline.RunPipeline(ctx, []string{prompt})
	if err != nil {
		return "", fmt.Errorf("run LLM pipeline: %w", err)
	}
	if output == nil || len(output.Responses) != 1 {
		got := 0
		if output != nil {
			got = len(output.Responses)
		}
		return "", fmt.Errorf("pipeline returned %d responses for 1 prompt", got)
	}

	// output.Responses[0] is only the newly generated continuation (hugot
	// accumulates token deltas per sequence, not prompt+continuation), so no
	// prompt-stripping is needed — just trim incidental whitespace.
	return strings.TrimSpace(output.Responses[0]), nil
}

// Backend returns the backend name.
func (c *ONNXLLMClient) Backend() string {
	return "onnx"
}

// Model returns the model name.
func (c *ONNXLLMClient) Model() string {
	return c.model
}

// Close releases ONNX session resources.
func (c *ONNXLLMClient) Close() error {
	if c.session != nil {
		c.session.Destroy()
	}
	return nil
}

// GetModelPath returns the path where a model should be downloaded to.
// Path: ~/.hawp/models/llm/{model-name}/
func GetModelPath(home, model string) string {
	return filepath.Join(home, ".hawp", "models", "llm", model)
}
