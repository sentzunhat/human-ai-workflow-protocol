package embeddings

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync/atomic"

	"github.com/knights-analytics/hugot"
)

// ONNXEmbedder performs embedding via ONNX Runtime (local, fast, private).
// Supports BGE-base-en-v1.5 (768-dim) and all-MiniLM-L6-v2 (384-dim).
type ONNXEmbedder struct {
	model     string
	dimension int
	modelPath string
	session   *hugot.Session
}

// ModelInfo holds metadata for supported ONNX embedding models.
type ModelInfo struct {
	Name             string
	Dimension        int
	HFRepo           string // Hugging Face repo for download
	ONNXFile         string // ONNX file to use in the repo
	ExternalDataFile string // external weights sidecar (e.g. "onnx/model.onnx_data"), "" if none
}

// SupportedModels lists all available ONNX embedding models.
var SupportedModels = map[string]ModelInfo{
	"bge-base-en-v1.5": {
		Name:      "bge-base-en-v1.5",
		Dimension: 768,
		HFRepo:    "BAAI/bge-base-en-v1.5",
		ONNXFile:  "onnx/model.onnx",
	},
	"all-MiniLM-L6-v2": {
		Name:      "all-MiniLM-L6-v2",
		Dimension: 384,
		HFRepo:    "sentence-transformers/all-MiniLM-L6-v2",
		ONNXFile:  "onnx/model.onnx",
	},
	// mdbr-leaf-ir: MongoDB's retrieval-tuned embedding model (BERT, 384-dim,
	// standard optimum ONNX export — plain feature-extraction, so it works
	// with the existing Go/XLA backend same as the models above; no CGO/ORT
	// needed, unlike ONNX LLM generation). Fine-tuned specifically for
	// retrieval quality rather than general-purpose similarity — this is
	// what "RAG retrieval" means in HAWP: a better embedding model choice,
	// not a separate IR subsystem. See v0.1.0_VISION.md.
	"mdbr-leaf-ir": {
		Name:             "mdbr-leaf-ir",
		Dimension:        384,
		HFRepo:           "MongoDB/mdbr-leaf-ir",
		ONNXFile:         "onnx/model.onnx",
		ExternalDataFile: "onnx/model.onnx_data",
	},
	// bge-large-en-v1.5 (1024-dim) was attempted as a "bigger model" option
	// but is NOT included: downloading it deadlocks inside the go-huggingface
	// downloader dependency (hugot@v0.7.5 -> github.com/gomlx/go-huggingface
	// @v0.3.5's internal/downloader.Manager) — a Go runtime "all goroutines
	// are asleep" deadlock, reproduced 2026-07-26, not a slow-network timeout.
	// Do not add it until that upstream issue is understood/fixed; shipping
	// a model that hangs on first use would repeat the exact class of bug
	// this release is fixing elsewhere.
}

// DefaultModel is the recommended embedding model (best quality).
const DefaultModel = "bge-base-en-v1.5"

// pipelineCounter ensures unique pipeline names across embedder instances
var pipelineCounter atomic.Int64

// NewONNXEmbedder creates a new ONNX embedder for the given model.
// Model should be "bge-base-en-v1.5" (default) or "all-MiniLM-L6-v2".
// On first use, automatically downloads the model to ~/.hawp/models/embedding/.
func NewONNXEmbedder(model string) (*ONNXEmbedder, error) {
	if model == "" {
		model = DefaultModel
	}

	// Validate model name
	info, ok := SupportedModels[model]
	if !ok {
		return nil, fmt.Errorf("unsupported ONNX model: %s; supported: %v", model, supportedModelNames())
	}

	// Resolve model directory (~/.hawp/models/embedding/{model}/)
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("get home directory: %w", err)
	}

	modelsDir := filepath.Join(home, ".hawp", "models", "embedding")
	if err := os.MkdirAll(modelsDir, 0755); err != nil {
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
			return nil, fmt.Errorf("download model %s from %s: %w", model, info.HFRepo, err)
		}
	}

	// Create hugot session for feature extraction (embedding)
	ctx := context.Background()
	session, err := hugot.NewGoSession(ctx)
	if err != nil {
		return nil, fmt.Errorf("start hugot session: %w", err)
	}

	embedder := &ONNXEmbedder{
		model:     model,
		dimension: info.Dimension,
		modelPath: modelPath,
		session:   session,
	}

	return embedder, nil
}

// Embed returns the embedding vector for a single text.
func (e *ONNXEmbedder) Embed(ctx context.Context, text string) ([]float32, error) {
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

// EmbedBatch returns embedding vectors for multiple texts (more efficient).
func (e *ONNXEmbedder) EmbedBatch(ctx context.Context, texts []string) ([][]float32, error) {
	if len(texts) == 0 {
		return [][]float32{}, nil
	}

	// Create pipeline with unique name for this embedding run
	pipelineID := pipelineCounter.Add(1)
	config := hugot.FeatureExtractionConfig{
		ModelPath: e.modelPath,
		Name:      fmt.Sprintf("hawp-embed-%s-%d", e.model, pipelineID),
	}

	pipeline, err := hugot.NewPipeline(e.session, config)
	if err != nil {
		return nil, fmt.Errorf("load embedding model at %s: %w", e.modelPath, err)
	}

	result, err := pipeline.RunPipeline(ctx, texts)
	if err != nil {
		return nil, fmt.Errorf("run embedding pipeline: %w", err)
	}

	if result == nil || len(result.Embeddings) != len(texts) {
		return nil, fmt.Errorf("pipeline returned %d embeddings for %d texts",
			len(result.Embeddings), len(texts))
	}

	return result.Embeddings, nil
}

// Dimension returns the vector dimension for this model.
func (e *ONNXEmbedder) Dimension() int {
	return e.dimension
}

// Backend returns the backend name.
func (e *ONNXEmbedder) Backend() string {
	return "onnx"
}

// Model returns the model name.
func (e *ONNXEmbedder) Model() string {
	return e.model
}

// Close releases ONNX session resources.
func (e *ONNXEmbedder) Close() error {
	if e.session != nil {
		e.session.Destroy()
	}
	return nil
}

// GetModelPath returns the path where a model should be downloaded to.
// Path: ~/.hawp/models/embedding/{model-name}/
func GetModelPath(home, model string) string {
	return filepath.Join(home, ".hawp", "models", "embedding", model)
}

// supportedModelNames returns a slice of supported model names.
func supportedModelNames() []string {
	names := make([]string, 0, len(SupportedModels))
	for name := range SupportedModels {
		names = append(names, name)
	}
	return names
}
