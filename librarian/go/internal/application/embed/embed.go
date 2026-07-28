// Package embed wraps hugot (github.com/knights-analytics/hugot) — an
// open-source "transformers.js for Go" — to run Hugging Face ONNX models
// locally: pulling a model by repo ID and running feature-extraction
// (embedding) inference through hugot's pure-Go backend (GoMLX), with no
// cgo and no external ONNX Runtime binary required. This is the
// foundation `hawp embed`/`hawp generate` build on; text generation
// needs hugot's ORT/cgo backend and is tracked separately (see
// 748609a8's Recommended Fix) rather than wired in here.
package embed

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/knights-analytics/hugot"
)

// DefaultModel is the pinned embedding model already provisioned by
// `hawp init` (internal/domain/provision) — reused here as the default
// so `hawp embed` works out of the box without a separate pull step.
const DefaultModel = "Xenova/all-MiniLM-L6-v2"

// defaultOnnxFile is the specific ONNX export to load: the repo has
// several variants (model.onnx, model_quantized.onnx, ...) under onnx/,
// and hugot's downloader errors if more than one exists without an
// explicit choice.
const defaultOnnxFile = "onnx/model_quantized.onnx"

// PullModel downloads modelRepo (a "org/name" Hugging Face repo ID) into
// destDir/<org_name>/, verified to contain exactly one ONNX file (or the
// one named onnxFile, if given) plus its tokenizer. Returns the local
// model directory hugot pipelines load from.
func PullModel(ctx context.Context, modelRepo, onnxFile, destDir string) (string, error) {
	opts := hugot.NewDownloadOptions()
	opts.OnnxFilePath = onnxFile
	return hugot.DownloadModel(ctx, modelRepo, destDir, opts)
}

// Embed runs the feature-extraction pipeline over texts using the model
// at modelPath (as returned by PullModel), via hugot's pure-Go backend.
// Returns one embedding vector per input text, in order.
func Embed(ctx context.Context, modelPath string, texts []string) ([][]float32, error) {
	session, err := hugot.NewGoSession(ctx)
	if err != nil {
		return nil, fmt.Errorf("start hugot session: %w", err)
	}
	defer session.Destroy()

	config := hugot.FeatureExtractionConfig{
		ModelPath: modelPath,
		Name:      "hawp-embed",
	}
	pipeline, err := hugot.NewPipeline(session, config)
	if err != nil {
		return nil, fmt.Errorf("load embedding model at %s: %w", modelPath, err)
	}

	result, err := pipeline.RunPipeline(ctx, texts)
	if err != nil {
		return nil, fmt.Errorf("run embedding pipeline: %w", err)
	}

	vectors := make([][]float32, len(result.Embeddings))
	copy(vectors, result.Embeddings)
	return vectors, nil
}

// DefaultModelDir is where PullModel puts DefaultModel under a ~/.hawp/
// models root, matching the naming DownloadModel itself derives
// (org/name -> org_name).
func DefaultModelDir(modelsRoot string) string {
	return filepath.Join(modelsRoot, "Xenova_all-MiniLM-L6-v2")
}

// PullDefaultModel downloads DefaultModel (the pinned all-MiniLM-L6-v2
// embedding model) into destDir if not already present, returning its
// local directory.
func PullDefaultModel(ctx context.Context, destDir string) (string, error) {
	return PullModel(ctx, DefaultModel, defaultOnnxFile, destDir)
}
