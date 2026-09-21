package onnx

import (
	"context"
	"math"
	"os"
	"path/filepath"
	"testing"

	domainembeddings "github.com/sentzunhat/hawp/librarian/src/internal/domain/providers/embeddings"
)

// skipWithoutModel skips the test when the named embedding model is not
// already downloaded. NewONNXEmbedder downloads the model on first use
// (potentially hundreds of MB); CI runners don't have the models, so tests
// that require a real session must opt out when the model is absent.
func skipWithoutModel(t *testing.T, model string) {
	t.Helper()
	home, _ := os.UserHomeDir()
	path := filepath.Join(home, ".hawp", "models", "embedding", model)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Skipf("skip: model %q not found at %s; run `hawp model pull %s` first", model, path, model)
	}
}

func TestNewONNXEmbedder(t *testing.T) {
	skipWithoutModel(t, DefaultEmbedModel)
	embedder, err := NewONNXEmbedder("")
	if err != nil {
		t.Fatalf("NewONNXEmbedder with default model failed: %v", err)
	}
	if embedder.Model() != DefaultEmbedModel {
		t.Errorf("default model should be %s, got %s", DefaultEmbedModel, embedder.Model())
	}
	if embedder.Backend() != "onnx" {
		t.Errorf("backend should be onnx, got %s", embedder.Backend())
	}

	embedder, err = NewONNXEmbedder("bge-base-en-v1.5")
	if err != nil {
		t.Fatalf("NewONNXEmbedder with bge model failed: %v", err)
	}
	if embedder.Dimension() != 768 {
		t.Errorf("bge dimension should be 768, got %d", embedder.Dimension())
	}

	embedder, err = NewONNXEmbedder("all-MiniLM-L6-v2")
	if err != nil {
		t.Fatalf("NewONNXEmbedder with MiniLM model failed: %v", err)
	}
	if embedder.Dimension() != 384 {
		t.Errorf("MiniLM dimension should be 384, got %d", embedder.Dimension())
	}
}

func TestUnsupportedEmbedModel(t *testing.T) {
	_, err := NewONNXEmbedder("nonexistent-model")
	if err == nil {
		t.Error("unsupported model should return error")
	}
}

func TestEmbedderInterface(t *testing.T) {
	skipWithoutModel(t, "bge-base-en-v1.5")
	embedder, err := NewONNXEmbedder("bge-base-en-v1.5")
	if err != nil {
		t.Fatalf("failed to create embedder: %v", err)
	}
	defer embedder.Close()

	var _ domainembeddings.Embedder = embedder

	if embedder.Backend() != "onnx" {
		t.Errorf("Backend() should return onnx, got %s", embedder.Backend())
	}
	if embedder.Model() != "bge-base-en-v1.5" {
		t.Errorf("Model() should return bge-base-en-v1.5, got %s", embedder.Model())
	}
	if embedder.Dimension() != 768 {
		t.Errorf("Dimension() should return 768, got %d", embedder.Dimension())
	}
}

func TestEmptyText(t *testing.T) {
	skipWithoutModel(t, "bge-base-en-v1.5")
	embedder, err := NewONNXEmbedder("bge-base-en-v1.5")
	if err != nil {
		t.Fatalf("failed to create embedder: %v", err)
	}
	defer embedder.Close()

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

func TestEmbedBatchEmpty(t *testing.T) {
	skipWithoutModel(t, "all-MiniLM-L6-v2")
	embedder, err := NewONNXEmbedder("all-MiniLM-L6-v2")
	if err != nil {
		t.Fatalf("failed to create embedder: %v", err)
	}
	defer embedder.Close()

	embeddings, err := embedder.EmbedBatch(context.Background(), []string{})
	if err != nil {
		t.Fatalf("embed batch empty failed: %v", err)
	}
	if len(embeddings) != 0 {
		t.Errorf("empty batch should return empty slice, got %d embeddings", len(embeddings))
	}
}

func TestSupportedEmbedModels(t *testing.T) {
	if len(SupportedEmbedModels) < 2 {
		t.Errorf("should have at least 2 supported models, got %d", len(SupportedEmbedModels))
	}

	if info, ok := SupportedEmbedModels["bge-base-en-v1.5"]; !ok {
		t.Error("BGE model not in SupportedEmbedModels")
	} else {
		if info.Dimension != 768 {
			t.Errorf("BGE dimension should be 768, got %d", info.Dimension)
		}
		if info.HFRepo == "" {
			t.Error("BGE HFRepo should not be empty")
		}
	}

	if info, ok := SupportedEmbedModels["all-MiniLM-L6-v2"]; !ok {
		t.Error("MiniLM model not in SupportedEmbedModels")
	} else {
		if info.Dimension != 384 {
			t.Errorf("MiniLM dimension should be 384, got %d", info.Dimension)
		}
		if info.HFRepo == "" {
			t.Error("MiniLM HFRepo should not be empty")
		}
	}

	// mdbr-leaf-ir: MongoDB's retrieval model (verified working 2026-07-27).
	if info, ok := SupportedEmbedModels["mdbr-leaf-ir"]; !ok {
		t.Error("mdbr-leaf-ir model not in SupportedEmbedModels")
	} else {
		if info.Dimension != 384 {
			t.Errorf("mdbr-leaf-ir dimension should be 384, got %d", info.Dimension)
		}
		if info.HFRepo == "" {
			t.Error("mdbr-leaf-ir HFRepo should not be empty")
		}
		if info.ExternalDataFile == "" {
			t.Error("mdbr-leaf-ir ExternalDataFile should not be empty — its weights are split into a sidecar file")
		}
	}
}

func TestGetEmbedModelPath(t *testing.T) {
	home := "/home/user"
	model := "bge-base-en-v1.5"
	path := GetEmbedModelPath(home, model)
	if path != "/home/user/.hawp/models/embedding/bge-base-en-v1.5" {
		t.Errorf("model path incorrect: %s", path)
	}
}

func TestEmbedSingleText(t *testing.T) {
	skipWithoutModel(t, "all-MiniLM-L6-v2")
	embedder, err := NewONNXEmbedder("all-MiniLM-L6-v2")
	if err != nil {
		t.Fatalf("failed to create embedder: %v", err)
	}
	defer embedder.Close()

	text := "This is a test sentence for embedding."
	vec, err := embedder.Embed(context.Background(), text)
	if err != nil {
		t.Fatalf("embedding failed: %v", err)
	}
	if len(vec) != embedder.Dimension() {
		t.Errorf("embedding dimension mismatch: expected %d, got %d", embedder.Dimension(), len(vec))
	}

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

func TestEmbedBatchMultipleTexts(t *testing.T) {
	skipWithoutModel(t, "all-MiniLM-L6-v2")
	embedder, err := NewONNXEmbedder("all-MiniLM-L6-v2")
	if err != nil {
		t.Fatalf("failed to create embedder: %v", err)
	}
	defer embedder.Close()

	texts := []string{
		"Hello, world!",
		"This is another test.",
		"Embeddings are useful for semantic search.",
	}
	embeddings, err := embedder.EmbedBatch(context.Background(), texts)
	if err != nil {
		t.Fatalf("batch embedding failed: %v", err)
	}
	if len(embeddings) != len(texts) {
		t.Errorf("batch size mismatch: expected %d, got %d", len(texts), len(embeddings))
	}
	for i, vec := range embeddings {
		if len(vec) != embedder.Dimension() {
			t.Errorf("embedding[%d] dimension mismatch: expected %d, got %d",
				i, embedder.Dimension(), len(vec))
		}
		hasNonZero := false
		for _, v := range vec {
			if v != 0 {
				hasNonZero = true
				break
			}
		}
		if !hasNonZero {
			t.Errorf("embedding[%d] should not be all zeros", i)
		}
	}
}

func TestDifferentTextsProduceDifferentEmbeddings(t *testing.T) {
	skipWithoutModel(t, "all-MiniLM-L6-v2")
	embedder, err := NewONNXEmbedder("all-MiniLM-L6-v2")
	if err != nil {
		t.Fatalf("failed to create embedder: %v", err)
	}
	defer embedder.Close()

	text1 := "I love sunny days at the beach."
	text2 := "The cat is sleeping under the table."

	vec1, err := embedder.Embed(context.Background(), text1)
	if err != nil {
		t.Fatalf("embedding text1 failed: %v", err)
	}
	vec2, err := embedder.Embed(context.Background(), text2)
	if err != nil {
		t.Fatalf("embedding text2 failed: %v", err)
	}

	similarity := cosineSimilarity(vec1, vec2)
	if similarity > 0.9 {
		t.Errorf("different texts should have lower similarity, got %f", similarity)
	}
}

func cosineSimilarity(a, b []float32) float32 {
	if len(a) != len(b) || len(a) == 0 {
		return 0
	}
	var dotProduct, normA, normB float32
	for i := range a {
		dotProduct += a[i] * b[i]
		normA += a[i] * a[i]
		normB += b[i] * b[i]
	}
	if normA == 0 || normB == 0 {
		return 0
	}
	return dotProduct / float32(math.Sqrt(float64(normA)*float64(normB)))
}
