package context

import (
	"context"
	"fmt"

	"github.com/sentzunhat/hawp/librarian/go/internal/application/search"
	domainsearch "github.com/sentzunhat/hawp/librarian/go/internal/domain/search"
)

// DocumentReference tracks the source and position of a retrieved document.
// Used to maintain provenance through the RAG pipeline.
type DocumentReference struct {
	Source    string  // Document source/path (e.g., "README.md", "guide/setup.md")
	Title     string  // Human-readable title or heading
	Content   string  // Matched excerpt from this source (from its highest-relevance chunk)
	Relevance float32 // Confidence score (0.0 - 1.0)
	LineStart int     // Line number where chunk starts in source
	LineEnd   int     // Line number where chunk ends
}

// RAGPipelineOutput is the final result of a complete RAG pipeline run.
// It maintains references to the original sources while providing reshaped content.
type RAGPipelineOutput struct {
	Content     string              // Reshaped/improved content from the pipeline
	References  []DocumentReference // Provenance: which documents were used
	Pipeline    string              // Backend combination used (e.g., "onnx-ollama")
	KeyConcepts []string            // High-level themes identified by the pipeline
}

// RAGPipeline orchestrates the full retrieval-augmented generation flow:
// 1. Retrieve: Take a query and return relevant documents (via ContextBlock)
// 2. Reshape: Process documents through embeddings + LLM to restructure and improve clarity
// 3. Return: Reshaped content with references to original sources
//
// This interface allows v0.1.0+ to plug in IR retrieval (BM25, dense passage retrieval, etc.)
// without changing the API. In v0.0.3, Retrieve() is a no-op (accepts pre-formatted ContextBlock);
// only Reshape() is active.
type RAGPipeline interface {
	// Retrieve searches for documents matching the query and returns them as a ContextBlock.
	// In v0.0.3, this is unused (search happens in CLI before RAGPipeline is invoked).
	// In v0.1.0+, this could fetch results from BM25, dense retrieval, etc.
	Retrieve(ctx context.Context, query string, topK int) (ContextBlock, error)

	// Reshape takes a formatted ContextBlock and improves it via embeddings + LLM.
	// Returns restructured content with key concepts and source references.
	Reshape(ctx context.Context, block ContextBlock, maxTokens int) (RAGPipelineOutput, error)

	// Close releases all resources (embedder, LLM, DB connections, etc.).
	Close() error
}

// DefaultRAGPipeline implements RAGPipeline using local embeddings + LLM.
// It wraps a ContextReshaper and bridges between ContextBlock (search results)
// and RAGPipelineOutput (references + reshaped content).
type DefaultRAGPipeline struct {
	reshaper  *ContextReshaper
	retriever Retriever
}

// Retriever supplies ranked results for a context query. It is intentionally
// local to this capability so retrieval can be tested or replaced without
// coupling the RAG orchestration to a concrete index.
type Retriever interface {
	Retrieve(ctx context.Context, query string, topK int) ([]domainsearch.Result, error)
}

type localIndexRetriever struct {
	repoRoot string
}

func (r localIndexRetriever) Retrieve(ctx context.Context, query string, topK int) ([]domainsearch.Result, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	return search.Query(r.repoRoot, query, topK)
}

// NewDefaultRAGPipeline creates a RAG pipeline with the specified
// configuration, scoped to the index at repoRoot/.hawp/db/index.sqlite (used
// by Retrieve). It initializes embeddings and LLM backends.
//
// Example:
//
//	config := ReshapingConfig{
//	    EmbeddingsBackend: "onnx",
//	    EmbeddingsModel:   "all-MiniLM-L6-v2",
//	    LLMBackend:        "ollama",
//	    LLMModel:          "mistral",
//	}
//	pipeline, err := NewDefaultRAGPipeline(config, repoRoot)
//	if err != nil {
//	    log.Fatalf("Failed to create pipeline: %v", err)
//	}
//	defer pipeline.Close()
//
//	block, err := pipeline.Retrieve(ctx, "how does auth work", 10)
//	output, err := pipeline.Reshape(ctx, block, 2000)
func NewDefaultRAGPipeline(config ReshapingConfig, repoRoot string) (*DefaultRAGPipeline, error) {
	reshaper, err := NewContextReshaper(config)
	if err != nil {
		return nil, fmt.Errorf("initialize reshaper: %w", err)
	}

	return NewRAGPipeline(reshaper, localIndexRetriever{repoRoot: repoRoot})
}

// NewRAGPipeline composes an injected reshaper and retriever. The default
// constructor remains the compatibility entry point for the local index.
func NewRAGPipeline(reshaper *ContextReshaper, retriever Retriever) (*DefaultRAGPipeline, error) {
	if reshaper == nil {
		return nil, fmt.Errorf("reshaper is required")
	}
	if retriever == nil {
		return nil, fmt.Errorf("retriever is required")
	}
	return &DefaultRAGPipeline{reshaper: reshaper, retriever: retriever}, nil
}

// Retrieve runs query against the local index (internal/application/search —
// the same lexical+hybrid search `hawp search <query>` uses, so there is one
// search implementation, not a separate "RAG retrieval" system) and formats
// the results as a ContextBlock via FormatAsMarkdown.
//
// This is deliberately the same embeddings the rest of HAWP uses for
// search — retrieval-in-RAG and semantic search are the same operation
// (rank documents by embedding similarity to a query); there is no separate
// "IR model" layer to build. A specialized retrieval-tuned embedding model
// (e.g. one of MongoDB's mdbr-leaf models) would slot in as just another
// entry in embeddings.SupportedModels, not a new subsystem.
func (p *DefaultRAGPipeline) Retrieve(ctx context.Context, query string, topK int) (ContextBlock, error) {
	results, err := p.retriever.Retrieve(ctx, query, topK)
	if err != nil {
		return ContextBlock{}, fmt.Errorf("retrieve: %w", err)
	}
	return PrepareContext(results, query, defaultRetrieveMaxTokens), nil
}

// defaultRetrieveMaxTokens bounds FormatAsMarkdown's token budget when
// Retrieve builds the initial ContextBlock. Reshape can still be called
// with its own, separate maxTokens for the LLM step.
const defaultRetrieveMaxTokens = 4000

// Reshape takes a ContextBlock (formatted search results) and improves it via
// embeddings + LLM reshaping. Returns the reshaped content along with references
// to the original source documents and identified key concepts.
//
// The operation is fully context-aware (respects ctx cancellation).
// Returns error if the block is empty.
func (p *DefaultRAGPipeline) Reshape(ctx context.Context, block ContextBlock, maxTokens int) (RAGPipelineOutput, error) {
	if len(block.Results) == 0 {
		return RAGPipelineOutput{}, fmt.Errorf("empty context block")
	}

	// Run embeddings + LLM reshape on the block
	reshaped, err := p.reshaper.ReshapeWithMaxTokens(ctx, &block, maxTokens)
	if err != nil {
		return RAGPipelineOutput{}, fmt.Errorf("reshape failed: %w", err)
	}

	// Reuse the deduplicated References already computed by FormatAsMarkdown
	// (one entry per unique Source, highest-relevance Content kept). Callers
	// who build a ContextBlock by hand instead of via FormatAsMarkdown (e.g.
	// direct RAGPipeline API use) won't have References populated yet —
	// derive it from Results in that case so provenance is never silently lost.
	references := block.References
	if references == nil && len(block.Results) > 0 {
		references = deduplicateReferences(block.Results)
	}

	// Extract key concept texts for the output
	keyConceptTexts := make([]string, len(reshaped.KeyConcepts))
	for i, concept := range reshaped.KeyConcepts {
		keyConceptTexts[i] = concept.Text
	}

	return RAGPipelineOutput{
		Content:     reshaped.Content,
		References:  references,
		Pipeline:    reshaped.Pipeline,
		KeyConcepts: keyConceptTexts,
	}, nil
}

// Close releases all resources (embedder, LLM, etc.).
func (p *DefaultRAGPipeline) Close() error {
	if p.reshaper == nil {
		return nil
	}
	return p.reshaper.Close()
}
