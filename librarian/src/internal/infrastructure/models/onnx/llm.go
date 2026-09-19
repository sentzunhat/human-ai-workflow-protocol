package onnx

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/knights-analytics/hugot"
	"github.com/knights-analytics/hugot/backends"
	"github.com/knights-analytics/hugot/options"
	"github.com/knights-analytics/hugot/pipelines"

	domainllm "github.com/sentzunhat/hawp/librarian/src/internal/domain/providers/llm"
)

// ONNXLLMClient performs text generation via ONNX Runtime (local, fast, private).
// Supports SmolLM2-360M-Instruct (see SupportedLLMModels).
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

// LLMModelInfo holds metadata for supported ONNX LLM models.
type LLMModelInfo struct {
	Name             string
	HFRepo           string // Hugging Face repo for download
	ONNXFile         string // ONNX file to use in the repo (passed as opts.OnnxFilePath)
	ExternalDataFile string // External weights sidecar (opts.ExternalDataPath), "" if none
	Params           int64  // Number of parameters (for reference)
}

// SupportedLLMModels lists all available LLM models with ONNX export.
// NOTE: hugot's ORT generative pipeline requires a genai_config.json —
// the ONNX Runtime GenAI model-builder format, NOT the standard
// optimum/Transformers.js ONNX export. Both entries below use community
// genai-format conversions that include genai_config.json.
var SupportedLLMModels = map[string]LLMModelInfo{
	// Phi-3-mini-4k-instruct: benchmarked at 10/10 coverage for HAWP intake
	// shaping (2026-09-11). Recommended default. Requires ORT GenAI build.
	// External data sidecar required — without model.onnx.data, ORT fails
	// at session-init with "External data path does not exist".
	"Phi-3-mini-4k-instruct": {
		Name:             "Phi-3-mini-4k-instruct",
		HFRepo:           "microsoft/Phi-3-mini-4k-instruct-ort-genai-int4-cpu",
		ONNXFile:         "model.onnx",
		ExternalDataFile: "model.onnx.data",
		Params:           3_800_000_000,
	},
	// SmolLM2-360M-Instruct: original default (v0.0.14–v0.0.23), scored 1/10
	// on HAWP intake shaping. Kept for backwards compatibility with downloaded
	// models; not recommended for new installs.
	"SmolLM2-360M-Instruct": {
		Name:             "SmolLM2-360M-Instruct",
		HFRepo:           "homen3/SmolLM2-360M-Instruct-ort-genai-int4-cpu",
		ONNXFile:         "model.onnx",
		ExternalDataFile: "model.onnx.data",
		Params:           360_000_000,
	},
}

// DefaultLLMModel is the recommended ONNX LLM model for context reshaping.
// Phi-3-mini-4k-instruct benchmarked at 10/10 coverage (2026-09-11).
const DefaultLLMModel = "Phi-3-mini-4k-instruct"

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
// Model may be:
//   - A short name from SupportedLLMModels (e.g. "SmolLM2-360M-Instruct") —
//     the model is located at ~/.hawp/models/llm/{name}/ and downloaded on
//     first use.
//   - An absolute directory path to an already-present model directory —
//     used directly without any download or SupportedLLMModels lookup. This
//     lets the benchmark CLI and manual tests point at models in non-standard
//     locations (e.g. ~/.hawp/models/llm/homen3_SmolLM2-360M-Instruct-ort-genai-int4-cpu).
//
// Requires a `-tags ORT` build with the native libraries described in this
// file's package doc comment; without them, this fails immediately with a
// clear error (from hugot itself, or from the ORT library-path check below)
// rather than downloading a model that can't run.
func NewONNXLLMClient(model string) (*ONNXLLMClient, error) {
	if model == "" {
		model = DefaultLLMModel
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

	// If model is an absolute path to an existing directory, use it directly
	// (no SupportedLLMModels lookup, no download).
	if filepath.IsAbs(model) {
		if fi, statErr := os.Stat(model); statErr == nil && fi.IsDir() {
			return &ONNXLLMClient{
				model:     filepath.Base(model),
				modelPath: model,
				session:   session,
				maxTokens: 512,
			}, nil
		}
		session.Destroy()
		return nil, fmt.Errorf("model path %s does not exist or is not a directory", model)
	}

	info, ok := SupportedLLMModels[model]
	if !ok {
		session.Destroy()
		return nil, fmt.Errorf("unsupported ONNX model: %s; supported: %v", model, llmSupportedModelNames())
	}

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

	return &ONNXLLMClient{
		model:     model,
		modelPath: modelPath,
		session:   session,
		maxTokens: 512,
	}, nil
}

// estimateTokenCount roughly estimates token count using the same ~4
// chars-per-token heuristic used elsewhere in HAWP.
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

// reshapeOne reshapes a single context by passing system + user content as
// separate messages, letting onnxruntime-genai apply the model's built-in
// chat template (avoids double-templating if we build the string ourselves).
func (c *ONNXLLMClient) reshapeOne(ctx context.Context, text string, maxTokens int) (string, error) {
	instruction := strings.ReplaceAll(domainllm.ReshapingPrompt, "{maxTokens}", fmt.Sprintf("%d", maxTokens))
	return c.generateFromMessages(ctx, instruction, text, maxTokens)
}

// IntakeReshape generates text with a caller-supplied system prompt placed in
// the system turn, letting the model's built-in chat template handle formatting.
func (c *ONNXLLMClient) IntakeReshape(ctx context.Context, systemPrompt, userContent string, maxTokens int) (string, error) {
	return c.generateFromMessages(ctx, systemPrompt, userContent, maxTokens)
}

// generateFromMessages passes system + user content as separate Message objects
// to hugot's RunMessages, which hands them to onnxruntime-genai for chat-template
// application. This is correct: passing a pre-formatted prompt string through
// RunPipeline causes double-templating (hugot wraps it as a "user" message and
// then onnxruntime-genai applies the model's template a second time), confirmed
// 2026-09-11 as the root cause of all "max length reached" failures on Phi-3.
func (c *ONNXLLMClient) generateFromMessages(ctx context.Context, systemPrompt, userContent string, maxTokens int) (string, error) {
	if maxTokens <= 0 {
		maxTokens = c.maxTokens
	}
	// maxTokens + 1024 gives onnxruntime-genai a comfortable total-sequence
	// budget: ~1024 tokens covers the prompt overhead across all supported
	// model types, and maxTokens is the target for the generated response.
	totalMaxLength := maxTokens + 1024

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

	messages := [][]backends.Message{{
		{Role: "system", Content: systemPrompt},
		{Role: "user", Content: userContent},
	}}

	output, err := pipeline.RunMessages(ctx, messages)
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
func (c *ONNXLLMClient) Backend() string { return "onnx" }

// Model returns the model name.
func (c *ONNXLLMClient) Model() string { return c.model }

// Close releases ONNX session resources.
func (c *ONNXLLMClient) Close() error {
	if c.session != nil {
		c.session.Destroy()
	}
	return nil
}

// GetLLMModelPath returns the path where an LLM model should be downloaded to.
// Path: ~/.hawp/models/llm/{model-name}/
func GetLLMModelPath(home, model string) string {
	return filepath.Join(home, ".hawp", "models", "llm", model)
}

func llmSupportedModelNames() []string {
	names := make([]string, 0, len(SupportedLLMModels))
	for name := range SupportedLLMModels {
		names = append(names, name)
	}
	return names
}
