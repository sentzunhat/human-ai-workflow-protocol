package embeddings

import (
	"context"
	"math"
	"testing"
)

func TestNewONNXEmbedder(t *testing.T) {
	// Test default model
	embedder, err := NewONNXEmbedder("")
	if err != nil {
		t.Fatalf("NewONNXEmbedder with default model failed: %v", err)
	}
	if embedder.Model() != DefaultModel {
		t.Errorf("default model should be %s, got %s", DefaultModel, embedder.Model())
	}
	if embedder.Backend() != "onnx" {
		t.Errorf("backend should be onnx, got %s", embedder.Backend())
	}

	// Test explicit model
	embedder, err = NewONNXEmbedder("bge-base-en-v1.5")
	if err != nil {
		t.Fatalf("NewONNXEmbedder with bge model failed: %v", err)
	}
	if embedder.Dimension() != 768 {
		t.Errorf("bge dimension should be 768, got %d", embedder.Dimension())
	}

	// Test lighter model
	embedder, err = NewONNXEmbedder("all-MiniLM-L6-v2")
	if err != nil {
		t.Fatalf("NewONNXEmbedder with MiniLM model failed: %v", err)
	}
	if embedder.Dimension() != 384 {
		t.Errorf("MiniLM dimension should be 384, got %d", embedder.Dimension())
	}
}

func TestUnsupportedModel(t *testing.T) {
	_, err := NewONNXEmbedder("nonexistent-model")
	if err == nil {
		t.Error("unsupported model should return error")
	}
}

func TestEmbedderInterface(t *testing.T) {
	embedder, err := NewONNXEmbedder("bge-base-en-v1.5")
	if err != nil {
		t.Fatalf("failed to create embedder: %v", err)
	}
	defer embedder.Close()

	// Test that embedder implements interface
	var _ Embedder = embedder

	// Test backend and model names
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
	embedder, err := NewONNXEmbedder("bge-base-en-v1.5")
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

	// Check all values are zero
	for _, v := range vec {
		if v != 0 {
			t.Errorf("zero vector should have all zeros, got %f", v)
		}
	}
}

func TestEmbedBatchEmpty(t *testing.T) {
	embedder, err := NewONNXEmbedder("all-MiniLM-L6-v2")
	if err != nil {
		t.Fatalf("failed to create embedder: %v", err)
	}
	defer embedder.Close()

	// Empty batch should return empty slice
	embeddings, err := embedder.EmbedBatch(context.Background(), []string{})
	if err != nil {
		t.Fatalf("embed batch empty failed: %v", err)
	}

	if len(embeddings) != 0 {
		t.Errorf("empty batch should return empty slice, got %d embeddings", len(embeddings))
	}
}

func TestSupportedModels(t *testing.T) {
	if len(SupportedModels) < 2 {
		t.Errorf("should have at least 2 supported models, got %d", len(SupportedModels))
	}

	// Check BGE model
	if info, ok := SupportedModels["bge-base-en-v1.5"]; !ok {
		t.Error("BGE model not in SupportedModels")
	} else {
		if info.Dimension != 768 {
			t.Errorf("BGE dimension should be 768, got %d", info.Dimension)
		}
		if info.HFRepo == "" {
			t.Error("BGE HFRepo should not be empty")
		}
	}

	// Check MiniLM model
	if info, ok := SupportedModels["all-MiniLM-L6-v2"]; !ok {
		t.Error("MiniLM model not in SupportedModels")
	} else {
		if info.Dimension != 384 {
			t.Errorf("MiniLM dimension should be 384, got %d", info.Dimension)
		}
		if info.HFRepo == "" {
			t.Error("MiniLM HFRepo should not be empty")
		}
	}

	// Check MongoDB's retrieval-tuned model (verified working live 2026-07-27
	// — plain feature-extraction ONNX, no CGO/ORT needed, unlike ONNX LLM).
	if info, ok := SupportedModels["mdbr-leaf-ir"]; !ok {
		t.Error("mdbr-leaf-ir model not in SupportedModels")
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

func TestGetModelPath(t *testing.T) {
	home := "/home/user"
	model := "bge-base-en-v1.5"

	path := GetModelPath(home, model)

	if path != "/home/user/.hawp/models/embedding/bge-base-en-v1.5" {
		t.Errorf("model path incorrect: %s", path)
	}
}

// TestEmbedSingleText tests actual embedding inference on a single text.
// Requires the model to be downloaded first (happens in NewONNXEmbedder).
func TestEmbedSingleText(t *testing.T) {
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

	// Vector should not be all zeros
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

// TestEmbedBatchMultipleTexts tests batch embedding inference.
func TestEmbedBatchMultipleTexts(t *testing.T) {
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

// TestDifferentTextsProduceDifferentEmbeddings verifies that semantically
// different texts produce different embeddings.
func TestDifferentTextsProduceDifferentEmbeddings(t *testing.T) {
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

	// Compute cosine similarity to check they're different
	similarity := cosineSimilarity(vec1, vec2)
	if similarity > 0.9 {
		t.Errorf("different texts should have lower similarity, got %f", similarity)
	}
}

// cosineSimilarity computes cosine similarity between two vectors.
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
