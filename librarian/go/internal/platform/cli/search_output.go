package cli

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	appcontext "github.com/sentzunhat/hawp/librarian/go/internal/application/context"
	"github.com/sentzunhat/hawp/librarian/go/internal/domain/search"
	"github.com/sentzunhat/hawp/librarian/go/internal/infrastructure/repo"
)

// prepareSearchContext keeps CLI context output on the same deduplication and
// formatting path used by application context retrieval.
func prepareSearchContext(results []search.Result, query string, maxTokens int) appcontext.ContextBlock {
	return appcontext.PrepareContext(results, query, maxTokens)
}

// toJSONReferences converts document references into the stable context-output
// shape consumed by downstream retrieval and prompt tooling.
func toJSONReferences(refs []appcontext.DocumentReference) []map[string]interface{} {
	out := make([]map[string]interface{}, len(refs))
	for i, r := range refs {
		out[i] = map[string]interface{}{
			"source":     r.Source,
			"title":      r.Title,
			"content":    r.Content,
			"relevance":  r.Relevance,
			"line_start": r.LineStart,
			"line_end":   r.LineEnd,
		}
	}
	return out
}

// tryReshapeViaRAGPipeline loads context config and runs the optional reshape
// step. Any failure is explicit and leaves the caller's unreshaped block intact.
func tryReshapeViaRAGPipeline(block appcontext.ContextBlock, maxTokens int) *appcontext.RAGPipelineOutput {
	home, err := os.UserHomeDir()
	hawpHome := ""
	if err == nil {
		hawpHome = filepath.Join(home, ".hawp")
	}

	root, _ := repo.FindBacklogRepoRoot(mustGetwd())
	cfg, err := appcontext.LoadContextConfig(hawpHome, root)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: --llm-reshape config error: %v; using unreshaped context.\n", err)
		return nil
	}

	pipeline, err := appcontext.NewDefaultRAGPipeline(appcontext.ReshapingConfig{
		EmbeddingsBackend: cfg.Embeddings.Backend,
		EmbeddingsModel:   cfg.Embeddings.Model,
		EmbeddingsURL:     cfg.Backends.Ollama.URL,
		LLMBackend:        cfg.LLM.Backend,
		LLMModel:          cfg.LLM.Model,
		LLMURL:            cfg.Backends.Ollama.URL,
		MaxTokens:         maxTokens,
		Temperature:       cfg.LLM.Temperature,
	}, root)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: --llm-reshape unavailable (%v); using unreshaped context. "+
			"Configure via ~/.hawp/config/context.json or HAWP_LLM_BACKEND/HAWP_EMBEDDINGS_BACKEND.\n", err)
		return nil
	}
	defer pipeline.Close()

	result, err := pipeline.Reshape(context.Background(), block, maxTokens)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: --llm-reshape failed (%v); using unreshaped context.\n", err)
		return nil
	}
	return &result
}

// renderReshapedWithReferences formats a RAG result for terminal output while
// retaining pipeline provenance and source references.
func renderReshapedWithReferences(title string, output *appcontext.RAGPipelineOutput) string {
	var sb strings.Builder

	if output.Pipeline == "none-none" {
		fmt.Fprintf(&sb, "# %s\n\n", title)
	} else {
		fmt.Fprintf(&sb, "# %s (LLM-reshaped, pipeline: %s)\n\n", title, output.Pipeline)
	}

	if len(output.KeyConcepts) > 0 {
		fmt.Fprintf(&sb, "**Key concepts:** %s\n\n", strings.Join(output.KeyConcepts, ", "))
	}

	sb.WriteString(output.Content)
	sb.WriteString("\n")

	if output.Pipeline != "none-none" && len(output.References) > 0 {
		sb.WriteString("\n---\n\n## References\n\n")
		for _, ref := range output.References {
			fmt.Fprintf(&sb, "**Reference:** %s (%d%% relevant)\n", ref.Source, int(ref.Relevance*100))
			if ref.LineStart > 0 || ref.LineEnd > 0 {
				fmt.Fprintf(&sb, "    Lines: %d-%d\n", ref.LineStart, ref.LineEnd)
			}
			sb.WriteString("\n")
		}
	}

	return sb.String()
}
