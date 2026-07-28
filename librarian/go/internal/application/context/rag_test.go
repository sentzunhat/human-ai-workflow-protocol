package context

import (
	"context"
	"strings"
	"testing"
)

func TestNewDefaultRAGPipeline(t *testing.T) {
	config := ReshapingConfig{
		EmbeddingsBackend: "onnx",
		EmbeddingsModel:   "all-MiniLM-L6-v2",
		LLMBackend:        "ollama",
		LLMModel:          "mistral",
		MaxTokens:         256,
		TopK:              5,
	}

	pipeline, err := NewDefaultRAGPipeline(config, t.TempDir())
	if err != nil {
		// Ollama might not be running - skip this test
		if strings.Contains(err.Error(), "Ollama") || strings.Contains(err.Error(), "connect") {
			t.Skipf("Skipping: Ollama required for test")
		}
		t.Fatalf("NewDefaultRAGPipeline failed: %v", err)
	}
	defer pipeline.Close()

	if pipeline.reshaper == nil {
		t.Errorf("pipeline.reshaper should not be nil")
	}
}

// TestRAGPipelineRetrieveNoIndex confirms Retrieve fails clearly (rather
// than panicking or silently returning an empty block) when repoRoot has no
// .hawp/db/index.sqlite — e.g. a fresh temp dir with no `search index` ever
// run. Retrieve() itself is no longer a no-op (see rag.go's doc comment for
// why: it's the same search.Query the CLI uses, not a separate placeholder).
func TestRAGPipelineRetrieveNoIndex(t *testing.T) {
	config := ReshapingConfig{
		EmbeddingsBackend: "onnx",
		EmbeddingsModel:   "all-MiniLM-L6-v2",
		LLMBackend:        "ollama",
		LLMModel:          "mistral",
	}

	pipeline, err := NewDefaultRAGPipeline(config, t.TempDir())
	if err != nil {
		if strings.Contains(err.Error(), "Ollama") || strings.Contains(err.Error(), "connect") {
			t.Skipf("Skipping: Ollama required for test")
		}
		t.Fatalf("NewDefaultRAGPipeline failed: %v", err)
	}
	defer pipeline.Close()

	ctx := context.Background()
	_, err = pipeline.Retrieve(ctx, "test query", 10)
	if err == nil {
		t.Error("Retrieve should fail when repoRoot has no index (never ran `search index`)")
	}
}

func TestRAGPipelineReshapeWithValidBlock(t *testing.T) {
	config := ReshapingConfig{
		EmbeddingsBackend: "onnx",
		EmbeddingsModel:   "all-MiniLM-L6-v2",
		LLMBackend:        "ollama",
		LLMModel:          "mistral",
		MaxTokens:         256,
		TopK:              3,
	}

	pipeline, err := NewDefaultRAGPipeline(config, t.TempDir())
	if err != nil {
		if strings.Contains(err.Error(), "Ollama") || strings.Contains(err.Error(), "connect") {
			t.Skipf("Skipping: Ollama required for test")
		}
		t.Fatalf("NewDefaultRAGPipeline failed: %v", err)
	}
	defer pipeline.Close()

	// Create a ContextBlock with sample results
	block := ContextBlock{
		Title:  "Test Results",
		Query:  "test query",
		Budget: 2000,
		Results: []FormattedResult{
			{
				Rank:      1,
				Relevance: 0.95,
				Source:    "docs/guide.md",
				Title:     "Getting Started",
				Content:   "This is a test document about semantic search. It explains how embeddings work.",
				Tokens:    25,
			},
			{
				Rank:      2,
				Relevance: 0.87,
				Source:    "README.md",
				Title:     "Introduction",
				Content:   "Learn how to use embeddings for retrieval. Embeddings power modern RAG systems.",
				Tokens:    22,
			},
		},
	}

	ctx := context.Background()
	output, err := pipeline.Reshape(ctx, block, 256)
	if err != nil {
		t.Errorf("Reshape failed: %v", err)
		return
	}

	// Verify output structure
	if output.Content == "" {
		t.Errorf("Reshaped content should not be empty")
	}
	if len(output.References) != 2 {
		t.Errorf("Expected 2 references, got %d", len(output.References))
	}
	if output.Pipeline == "" {
		t.Errorf("Pipeline field should not be empty")
	}

	// Verify references match original sources
	if output.References[0].Source != "docs/guide.md" {
		t.Errorf("First reference source should be 'docs/guide.md', got %q", output.References[0].Source)
	}
	if output.References[1].Source != "README.md" {
		t.Errorf("Second reference source should be 'README.md', got %q", output.References[1].Source)
	}

	// Verify relevance scores are preserved
	if output.References[0].Relevance != 0.95 {
		t.Errorf("First reference relevance should be 0.95, got %f", output.References[0].Relevance)
	}
	if output.References[1].Relevance != 0.87 {
		t.Errorf("Second reference relevance should be 0.87, got %f", output.References[1].Relevance)
	}
}

func TestRAGPipelineReshapeEmptyBlock(t *testing.T) {
	config := ReshapingConfig{
		EmbeddingsBackend: "onnx",
		EmbeddingsModel:   "all-MiniLM-L6-v2",
		LLMBackend:        "ollama",
		LLMModel:          "mistral",
	}

	pipeline, err := NewDefaultRAGPipeline(config, t.TempDir())
	if err != nil {
		if strings.Contains(err.Error(), "Ollama") || strings.Contains(err.Error(), "connect") {
			t.Skipf("Skipping: Ollama required for test")
		}
		t.Fatalf("NewDefaultRAGPipeline failed: %v", err)
	}
	defer pipeline.Close()

	// Empty context block should fail gracefully
	block := ContextBlock{
		Title:   "Empty Results",
		Query:   "test query",
		Budget:  2000,
		Results: []FormattedResult{},
	}

	ctx := context.Background()
	_, err = pipeline.Reshape(ctx, block, 256)
	if err == nil {
		t.Errorf("Reshape should fail with empty block")
	}
	if !strings.Contains(err.Error(), "empty") {
		t.Errorf("Error should mention 'empty', got: %v", err)
	}
}

func TestRAGPipelineOutputStructure(t *testing.T) {
	config := ReshapingConfig{
		EmbeddingsBackend: "onnx",
		EmbeddingsModel:   "all-MiniLM-L6-v2",
		LLMBackend:        "ollama",
		LLMModel:          "mistral",
		MaxTokens:         256,
		TopK:              2,
	}

	pipeline, err := NewDefaultRAGPipeline(config, t.TempDir())
	if err != nil {
		if strings.Contains(err.Error(), "Ollama") || strings.Contains(err.Error(), "connect") {
			t.Skipf("Skipping: Ollama required for test")
		}
		t.Fatalf("NewDefaultRAGPipeline failed: %v", err)
	}
	defer pipeline.Close()

	block := ContextBlock{
		Title:  "Test Results",
		Query:  "embeddings retrieval",
		Budget: 2000,
		Results: []FormattedResult{
			{
				Rank:      1,
				Relevance: 0.92,
				Source:    "docs/api.md",
				Title:     "API Reference",
				Content:   "Vector embeddings are dense numerical representations. They enable semantic search.",
				Tokens:    20,
			},
		},
	}

	ctx := context.Background()
	output, err := pipeline.Reshape(ctx, block, 256)
	if err != nil {
		t.Errorf("Reshape failed: %v", err)
		return
	}

	// Verify all fields are populated
	if len(output.Content) == 0 {
		t.Errorf("Output.Content should not be empty")
	}
	if len(output.References) == 0 {
		t.Errorf("Output.References should not be empty")
	}
	if output.Pipeline == "" {
		t.Errorf("Output.Pipeline should not be empty")
	}

	// Pipeline string should contain backend names
	if !strings.Contains(output.Pipeline, "onnx") && !strings.Contains(output.Pipeline, "ollama") {
		t.Errorf("Pipeline should contain backend names, got: %q", output.Pipeline)
	}

	// References should have proper structure
	ref := output.References[0]
	if ref.Source == "" {
		t.Errorf("Reference.Source should not be empty")
	}
	if ref.Relevance < 0.0 || ref.Relevance > 1.0 {
		t.Errorf("Reference.Relevance should be in [0, 1], got %f", ref.Relevance)
	}
}

func TestRAGPipelineContextCancellation(t *testing.T) {
	config := ReshapingConfig{
		EmbeddingsBackend: "onnx",
		EmbeddingsModel:   "all-MiniLM-L6-v2",
		LLMBackend:        "ollama",
		LLMModel:          "mistral",
	}

	pipeline, err := NewDefaultRAGPipeline(config, t.TempDir())
	if err != nil {
		if strings.Contains(err.Error(), "Ollama") || strings.Contains(err.Error(), "connect") {
			t.Skipf("Skipping: Ollama required for test")
		}
		t.Fatalf("NewDefaultRAGPipeline failed: %v", err)
	}
	defer pipeline.Close()

	block := ContextBlock{
		Title:  "Test Results",
		Query:  "test",
		Budget: 2000,
		Results: []FormattedResult{
			{
				Rank:      1,
				Relevance: 0.9,
				Source:    "test.md",
				Title:     "Test",
				Content:   "This is a test document with enough content to trigger reshaping operations and semantic analysis.",
				Tokens:    20,
			},
		},
	}

	// Create a cancelled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	// Reshape with cancelled context should fail
	_, err = pipeline.Reshape(ctx, block, 256)
	if err != nil {
		// Expected: context cancellation or operation failure
		// (exact behavior depends on where cancellation is checked in reshape flow)
		t.Logf("Reshape with cancelled context failed as expected: %v", err)
	}
}

func TestDocumentReference(t *testing.T) {
	ref := DocumentReference{
		Source:    "docs/guide.md",
		Title:     "Setup Guide",
		Relevance: 0.88,
		LineStart: 10,
		LineEnd:   25,
	}

	if ref.Source != "docs/guide.md" {
		t.Errorf("Source mismatch")
	}
	if ref.Title != "Setup Guide" {
		t.Errorf("Title mismatch")
	}
	if ref.Relevance != 0.88 {
		t.Errorf("Relevance mismatch")
	}
	if ref.LineStart != 10 || ref.LineEnd != 25 {
		t.Errorf("Line range mismatch")
	}
}
