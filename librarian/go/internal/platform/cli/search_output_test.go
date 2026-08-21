package cli

import (
	"context"
	"reflect"
	"testing"

	appcontext "github.com/sentzunhat/hawp/librarian/go/internal/application/context"
	"github.com/sentzunhat/hawp/librarian/go/internal/domain/search"
)

func TestPrepareSearchContextMatchesApplicationPreparation(t *testing.T) {
	results := []search.Result{
		{Source: "guide.md", Title: "Guide", Content: "First result.", Relevance: 0.9, Embedding: []float32{1, 0}},
		{Source: "guide.md", Title: "Guide", Content: "Duplicate result.", Relevance: 0.8, Embedding: []float32{1, 0}},
		{Source: "plan.md", Title: "Plan", Content: "Independent result.", Relevance: 0.7, Embedding: []float32{0, 1}},
	}

	want := appcontext.PrepareContext(results, "context query", 250)
	got := prepareSearchContext(results, "context query", 250)

	if !reflect.DeepEqual(got, want) {
		t.Fatalf("CLI context preparation diverged from application preparation (-want +got):\nwant: %#v\ngot:  %#v", want, got)
	}
}

// cliTestRetriever is a fake Retriever that returns a fixed result set.
// It makes the pipeline path deterministic in tests without a real index.
type cliTestRetriever struct{ results []search.Result }

func (r *cliTestRetriever) Retrieve(_ context.Context, _ string, _ int) ([]search.Result, error) {
	return r.results, nil
}

// cliTestEmbedder and cliTestLLM satisfy the domain interfaces required by
// NewContextReshaperWithClients. The reshaper is not exercised in
// TestCLIAndPipelineContextAreEquivalent (both backends are "none"), so these
// stubs only need to compile, not produce meaningful output.
type cliTestEmbedder struct{}

func (cliTestEmbedder) Embed(_ context.Context, _ string) ([]float32, error) {
	return []float32{1}, nil
}
func (cliTestEmbedder) EmbedBatch(_ context.Context, _ []string) ([][]float32, error) {
	return nil, nil
}
func (cliTestEmbedder) Dimension() int  { return 1 }
func (cliTestEmbedder) Backend() string { return "fake" }
func (cliTestEmbedder) Model() string   { return "fake" }
func (cliTestEmbedder) Close() error    { return nil }

type cliTestLLM struct{}

func (cliTestLLM) Reshape(_ context.Context, _ string, _ int) (string, error) {
	return "reshaped", nil
}
func (cliTestLLM) ReshapeBatch(_ context.Context, _ []string, _ int) ([]string, error) {
	return nil, nil
}
func (cliTestLLM) Backend() string { return "fake" }
func (cliTestLLM) Model() string   { return "fake" }
func (cliTestLLM) Close() error    { return nil }

// TestCLIAndPipelineContextAreEquivalent proves that the CLI context path
// (prepareSearchContext → PrepareContext) and the pipeline path
// (pipeline.Retrieve → PrepareContext) produce identical ContextBlocks for
// the same inputs. This is the regression guard that breaks if either path
// diverges: e.g. if the pipeline were changed to call FormatAsMarkdown
// directly (bypassing deduplication), or the CLI were changed to bypass
// PrepareContext.
//
// The token budget is appcontext.DefaultRetrieveMaxTokens so both paths share
// the same value: the pipeline's Retrieve always uses that constant, and the
// CLI path is called with it explicitly here to match.
func TestCLIAndPipelineContextAreEquivalent(t *testing.T) {
	results := []search.Result{
		{Source: "a.md", Title: "Alpha", Content: "Alpha content.", Relevance: 0.9, Embedding: []float32{1, 0}},
		// Near-duplicate of a.md — deduplication should collapse it.
		{Source: "a-copy.md", Title: "Alpha copy", Content: "Alpha content.", Relevance: 0.85, Embedding: []float32{1, 0}},
		{Source: "b.md", Title: "Beta", Content: "Beta content.", Relevance: 0.7, Embedding: []float32{0, 1}},
	}
	query := "equivalence check"

	// CLI path: prepareSearchContext delegates to PrepareContext.
	cliBlock := prepareSearchContext(results, query, appcontext.DefaultRetrieveMaxTokens)

	// Pipeline path: Retrieve runs the same PrepareContext logic after fetching
	// from the injected retriever.
	reshaper, err := appcontext.NewContextReshaperWithClients(
		appcontext.ReshapingConfig{EmbeddingsBackend: "none", LLMBackend: "none"},
		cliTestEmbedder{},
		cliTestLLM{},
	)
	if err != nil {
		t.Fatalf("NewContextReshaperWithClients: %v", err)
	}
	pipeline, err := appcontext.NewRAGPipeline(reshaper, &cliTestRetriever{results: results})
	if err != nil {
		t.Fatalf("NewRAGPipeline: %v", err)
	}

	pipelineBlock, err := pipeline.Retrieve(context.Background(), query, len(results))
	if err != nil {
		t.Fatalf("pipeline.Retrieve: %v", err)
	}

	if !reflect.DeepEqual(cliBlock, pipelineBlock) {
		t.Fatalf("CLI and pipeline context blocks diverged:\nCLI:      %#v\nPipeline: %#v",
			cliBlock, pipelineBlock)
	}
}
