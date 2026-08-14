package context

import (
	"context"
	"testing"

	"github.com/sentzunhat/hawp/librarian/go/internal/domain/embeddings"
	"github.com/sentzunhat/hawp/librarian/go/internal/domain/llm"
	"github.com/sentzunhat/hawp/librarian/go/internal/domain/search"
)

type fakeRetriever struct {
	results []search.Result
}

func (r fakeRetriever) Retrieve(context.Context, string, int) ([]search.Result, error) {
	return r.results, nil
}

type fakeEmbedder struct{}

func (fakeEmbedder) Embed(context.Context, string) ([]float32, error) { return []float32{1}, nil }
func (fakeEmbedder) EmbedBatch(context.Context, []string) ([][]float32, error) {
	return nil, nil
}
func (fakeEmbedder) Dimension() int  { return 1 }
func (fakeEmbedder) Backend() string { return "fake" }
func (fakeEmbedder) Model() string   { return "fake" }
func (fakeEmbedder) Close() error    { return nil }

var _ embeddings.Embedder = fakeEmbedder{}

type fakeLLM struct {
	maxTokens int
}

func (l *fakeLLM) Reshape(_ context.Context, _ string, maxTokens int) (string, error) {
	l.maxTokens = maxTokens
	return "reshaped", nil
}
func (l *fakeLLM) ReshapeBatch(context.Context, []string, int) ([]string, error) { return nil, nil }
func (l *fakeLLM) Backend() string                                               { return "fake" }
func (l *fakeLLM) Model() string                                                 { return "fake" }
func (l *fakeLLM) Close() error                                                  { return nil }

var _ llm.LLMClient = (*fakeLLM)(nil)

func TestRAGPipelineRetrieveUsesInjectedRetrieverAndPreparationPolicy(t *testing.T) {
	reshaper, err := NewContextReshaperWithClients(
		ReshapingConfig{EmbeddingsBackend: "none", LLMBackend: "none"},
		fakeEmbedder{},
		&fakeLLM{},
	)
	if err != nil {
		t.Fatalf("NewContextReshaperWithClients() error = %v", err)
	}

	pipeline, err := NewRAGPipeline(reshaper, fakeRetriever{results: []search.Result{
		{Source: "duplicate.md", Content: "duplicate", Relevance: 0.6, Embedding: []float32{1, 0}},
		{Source: "duplicate-copy.md", Content: "duplicate", Relevance: 0.9, Embedding: []float32{1, 0}},
		{Source: "unique.md", Content: "unique", Relevance: 0.8, Embedding: []float32{0, 1}},
	}})
	if err != nil {
		t.Fatalf("NewRAGPipeline() error = %v", err)
	}

	block, err := pipeline.Retrieve(context.Background(), "query", 3)
	if err != nil {
		t.Fatalf("Retrieve() error = %v", err)
	}
	if len(block.Results) != 2 {
		t.Fatalf("Retrieve() results = %d, want 2 after deduplication", len(block.Results))
	}
	if block.Results[0].Source != "unique.md" {
		t.Errorf("Retrieve() sorts prepared results by relevance; first source = %q, want unique.md", block.Results[0].Source)
	}
}

func TestRAGPipelineReshapeHonorsInvocationMaxTokens(t *testing.T) {
	client := &fakeLLM{}
	reshaper, err := NewContextReshaperWithClients(
		ReshapingConfig{EmbeddingsBackend: "none", LLMBackend: "fake", MaxTokens: 512},
		fakeEmbedder{},
		client,
	)
	if err != nil {
		t.Fatalf("NewContextReshaperWithClients() error = %v", err)
	}
	pipeline, err := NewRAGPipeline(reshaper, fakeRetriever{})
	if err != nil {
		t.Fatalf("NewRAGPipeline() error = %v", err)
	}

	_, err = pipeline.Reshape(context.Background(), ContextBlock{Results: []FormattedResult{{Content: "context", Source: "test.md"}}}, 123)
	if err != nil {
		t.Fatalf("Reshape() error = %v", err)
	}
	if client.maxTokens != 123 {
		t.Errorf("LLM maxTokens = %d, want invocation limit 123", client.maxTokens)
	}
}
