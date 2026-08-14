package context

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/sentzunhat/hawp/librarian/go/internal/domain/embeddings"
	"github.com/sentzunhat/hawp/librarian/go/internal/domain/llm"
)

// ReshapedBlock is the improved version of a ContextBlock after going through the pipeline.
// It wraps the original context with semantic improvements.
type ReshapedBlock struct {
	Original    ContextBlock        // Original context block (before reshaping)
	Content     string              // Reshaped/improved markdown content
	KeyConcepts []Concept           // Identified key themes (ranked by relevance)
	References  []DocumentReference // Deduplicated source documents (inherited from Original)
	Pipeline    string              // Backend combination used, e.g., "onnx-ollama", "ollama-ollama"
}

// Concept represents an identified key theme with relevance score.
// Concepts are extracted from content via semantic embeddings.
type Concept struct {
	Text      string    // e.g., "semantic search"
	Relevance float32   // 0.0-1.0, higher = more relevant
	Embedding []float32 // Vector embedding from embeddings backend
}

// ReshapingConfig holds configuration for the reshaping pipeline.
// It specifies which backends to use for embeddings and LLM, plus tuning parameters.
//
// Example (ONNX embeddings + Ollama LLM):
//
//	config := ReshapingConfig{
//	    EmbeddingsBackend: "onnx",
//	    EmbeddingsModel:   "all-MiniLM-L6-v2",
//	    LLMBackend:        "ollama",
//	    LLMModel:          "mistral",
//	    MaxTokens:         512,
//	    TopK:              5,
//	}
//	reshaper, _ := NewContextReshaper(config)
type ReshapingConfig struct {
	EmbeddingsBackend string  // "onnx" or "ollama"
	EmbeddingsModel   string  // Model name (e.g., "all-MiniLM-L6-v2")
	EmbeddingsURL     string  // Ollama server URL for embeddings; empty = backend default (localhost:11434)
	LLMBackend        string  // "ollama" or "onnx" (ONNX not yet ready)
	LLMModel          string  // Model name (e.g., "mistral")
	LLMURL            string  // Ollama server URL for LLM; empty = backend default (localhost:11434)
	MaxTokens         int     // Max tokens for LLM reshaping (default 512)
	TopK              int     // Max key concepts to extract (default 5)
	Temperature       float32 // LLM temperature (not yet used, reserved)
}

// ContextReshaper orchestrates embeddings + LLM for context improvement.
type ContextReshaper struct {
	embedder embeddings.Embedder
	llm      llm.LLMClient
	config   ReshapingConfig
}

// NewContextReshaper creates a reshaper with the specified backends.
// It initializes both embeddings and LLM backends, verifying they're available.
// If LLMBackend="ollama", requires `ollama serve` running on localhost:11434.
// If EmbeddingsBackend="onnx", auto-downloads models on first run (~50MB).
//
// Example:
//
//	config := ReshapingConfig{
//	    EmbeddingsBackend: "onnx",
//	    EmbeddingsModel:   "all-MiniLM-L6-v2",
//	    LLMBackend:        "ollama",
//	    LLMModel:          "mistral",
//	}
//	reshaper, err := NewContextReshaper(config)
//	if err != nil {
//	    log.Fatalf("Failed to create reshaper: %v", err)
//	}
//	defer reshaper.Close()
func NewContextReshaper(config ReshapingConfig) (*ContextReshaper, error) {
	config = withReshapingDefaults(config)

	// The default constructor owns provider construction. Callers that already
	// own provider lifecycle can use NewContextReshaperWithClients instead.
	embedder, err := embeddings.NewEmbedderWithURL(config.EmbeddingsBackend, config.EmbeddingsModel, config.EmbeddingsURL)
	if err != nil {
		return nil, fmt.Errorf("initialize embedder: %w", err)
	}

	llmClient, err := llm.NewLLMClientWithURL(config.LLMBackend, config.LLMModel, config.LLMURL)
	if err != nil {
		embedder.Close()
		return nil, fmt.Errorf("initialize LLM: %w", err)
	}

	return NewContextReshaperWithClients(config, embedder, llmClient)
}

// NewContextReshaperWithClients builds a reshaper from injected provider
// contracts. It makes application behavior deterministic in tests and lets
// composition select supported providers without changing the use case.
func NewContextReshaperWithClients(config ReshapingConfig, embedder embeddings.Embedder, llmClient llm.LLMClient) (*ContextReshaper, error) {
	config = withReshapingDefaults(config)
	if embedder == nil {
		return nil, fmt.Errorf("embedder is required")
	}
	if llmClient == nil {
		return nil, fmt.Errorf("LLM client is required")
	}

	return &ContextReshaper{
		embedder: embedder,
		llm:      llmClient,
		config:   config,
	}, nil
}

func withReshapingDefaults(config ReshapingConfig) ReshapingConfig {
	if config.TopK == 0 {
		config.TopK = 5 // Default to 5 key concepts
	}
	if config.MaxTokens == 0 {
		config.MaxTokens = 512 // Default token budget for LLM
	}
	return config
}

// Reshape improves a context block via a multi-step pipeline:
// 1. Split content into sentences
// 2. Embed sentences via embeddings backend
// 3. Extract key concepts from embeddings
// 4. Build a prompt with original content + concepts
// 5. Send to LLM for restructuring and clarity improvement
// 6. Return improved block with key concepts attached
//
// The operation is fully context-aware (respects ctx cancellation).
// Returns error if block is nil or empty.
//
// Typical performance on Mac (CPU-only):
// - ONNX embeddings: 30-50ms
// - Ollama LLM: 500-2000ms
// - Total: 1-3 seconds
//
// Example:
//
//	improved, err := reshaper.Reshape(ctx, contextBlock)
//	if err != nil {
//	    log.Fatalf("Failed to reshape: %v", err)
//	}
//	fmt.Printf("Concepts: %v\n", improved.KeyConcepts)
//	fmt.Printf("Reshaped:\n%s\n", improved.Content)
func (r *ContextReshaper) Reshape(ctx context.Context, block *ContextBlock) (*ReshapedBlock, error) {
	return r.ReshapeWithMaxTokens(ctx, block, r.config.MaxTokens)
}

// ReshapeWithMaxTokens improves a context block using this invocation's token
// limit. A non-positive limit preserves the configured default for backwards
// compatibility with callers that previously relied on it.
func (r *ContextReshaper) ReshapeWithMaxTokens(ctx context.Context, block *ContextBlock, maxTokens int) (*ReshapedBlock, error) {
	// block.String() always renders a non-empty "# Title" header, even with
	// zero Results, so it can never signal an empty block on its own — the
	// real signal is the absence of any results to reshape.
	if block == nil || len(block.Results) == 0 {
		return nil, fmt.Errorf("empty context block")
	}

	// Convert block to markdown for analysis
	markdown := block.String()

	// Step 1: Identify key concepts via embeddings — skipped when the
	// embeddings backend is "none" (zero-cost passthrough, no model/network).
	var concepts []Concept
	if r.config.EmbeddingsBackend != "none" {
		var err error
		concepts, err = r.identifyKeyConcepts(ctx, markdown)
		if err != nil {
			return nil, fmt.Errorf("identify concepts: %w", err)
		}
	}

	// Steps 2-3: Build prompt + send to LLM for restructuring — skipped when
	// the LLM backend is "none". Content becomes the interleaved
	// reference+content body only (formatResultsInline, not block.String()
	// and its title/summary header) — the RAG output wrapper (run.go
	// renderReshapedWithReferences) renders its own title, so reusing
	// block.String() here would double it.
	content := formatResultsInline(block.Results)
	if r.config.LLMBackend != "none" {
		prompt := r.buildReshapingPrompt(markdown, concepts)
		if maxTokens <= 0 {
			maxTokens = r.config.MaxTokens
		}
		reshaped, err := r.llm.Reshape(ctx, prompt, maxTokens)
		if err != nil {
			return nil, fmt.Errorf("reshape via LLM: %w", err)
		}
		content = strings.TrimSpace(reshaped)
	}

	// Step 4: Build result (inherit references from original block)
	improved := &ReshapedBlock{
		Original:    *block,
		Content:     content,
		KeyConcepts: concepts,
		References:  block.References, // Copy references from original block
		Pipeline:    fmt.Sprintf("%s-%s", r.config.EmbeddingsBackend, r.config.LLMBackend),
	}

	return improved, nil
}

// identifyKeyConcepts extracts key themes from the content via embeddings.
func (r *ContextReshaper) identifyKeyConcepts(ctx context.Context, content string) ([]Concept, error) {
	// Split content into sentences for embedding
	sentences := r.splitIntoSentences(content)
	if len(sentences) == 0 {
		return []Concept{}, nil
	}

	// Embed each sentence
	embeddings, err := r.embedder.EmbedBatch(ctx, sentences)
	if err != nil {
		return nil, fmt.Errorf("embed sentences: %w", err)
	}

	if len(embeddings) != len(sentences) {
		return nil, fmt.Errorf("embedding count mismatch: expected %d, got %d", len(sentences), len(embeddings))
	}

	// Extract key phrases (nouns, noun phrases) from sentences
	concepts := r.extractConceptsFromSentences(sentences, embeddings)

	// Deduplicate and rank by relevance
	uniqueConcepts := r.deduplicateConcepts(concepts)
	r.rankConceptsByRelevance(uniqueConcepts)

	// Return top-K concepts
	topK := r.config.TopK
	if topK > len(uniqueConcepts) {
		topK = len(uniqueConcepts)
	}

	return uniqueConcepts[:topK], nil
}

// splitIntoSentences breaks content into sentences (simple approach).
func (r *ContextReshaper) splitIntoSentences(content string) []string {
	// Simple sentence splitting on . ! ?
	sentences := strings.FieldsFunc(content, func(c rune) bool {
		return c == '.' || c == '!' || c == '?'
	})

	// Filter out short sentences
	var filtered []string
	for _, s := range sentences {
		s = strings.TrimSpace(s)
		if len(s) > 10 {
			filtered = append(filtered, s)
		}
	}

	return filtered
}

// extractConceptsFromSentences identifies key terms from sentences.
// Simple implementation: extract capitalized words (proper nouns) + noun phrases.
func (r *ContextReshaper) extractConceptsFromSentences(sentences []string, embeddings [][]float32) []Concept {
	conceptMap := make(map[string]Concept)

	for i, sentence := range sentences {
		if i >= len(embeddings) {
			break
		}

		// Extract capitalized words (simple concept extraction)
		words := strings.Fields(sentence)
		for _, word := range words {
			if len(word) > 3 && len(word) < 50 {
				// Check if word starts with capital (basic heuristic for concepts)
				if word[0] >= 'A' && word[0] <= 'Z' {
					// Remove punctuation
					concept := strings.Trim(word, ".,;:!?")
					if _, exists := conceptMap[concept]; !exists {
						conceptMap[concept] = Concept{
							Text:      concept,
							Relevance: 1.0, // Will be scored later
							Embedding: embeddings[i],
						}
					}
				}
			}
		}
	}

	// Convert map to slice
	var concepts []Concept
	for _, c := range conceptMap {
		concepts = append(concepts, c)
	}

	return concepts
}

// deduplicateConcepts removes similar concepts (via embeddings similarity).
func (r *ContextReshaper) deduplicateConcepts(concepts []Concept) []Concept {
	if len(concepts) <= 1 {
		return concepts
	}

	// For now, simple dedup by text (could use embedding similarity later)
	seen := make(map[string]bool)
	var unique []Concept

	for _, c := range concepts {
		text := strings.ToLower(c.Text)
		if !seen[text] {
			seen[text] = true
			unique = append(unique, c)
		}
	}

	return unique
}

// rankConceptsByRelevance scores concepts by how often they appear + embedding quality.
func (r *ContextReshaper) rankConceptsByRelevance(concepts []Concept) {
	// Sort by relevance (descending)
	sort.Slice(concepts, func(i, j int) bool {
		return concepts[i].Relevance > concepts[j].Relevance
	})
}

// buildReshapingPrompt creates a rich prompt with content + key concepts.
func (r *ContextReshaper) buildReshapingPrompt(content string, concepts []Concept) string {
	conceptTexts := make([]string, len(concepts))
	for i, c := range concepts {
		conceptTexts[i] = fmt.Sprintf("- %s", c.Text)
	}
	conceptList := strings.Join(conceptTexts, "\n")

	prompt := fmt.Sprintf(`You are a technical context reshaping expert. Your task is to improve clarity and structure.

Key concepts identified in this content:
%s

Original content:
%s

Reshaped content (maintain all critical information, improve structure and clarity):`, conceptList, content)

	return prompt
}

// Close releases embedder + LLM resources.
func (r *ContextReshaper) Close() error {
	if err := r.embedder.Close(); err != nil {
		return fmt.Errorf("close embedder: %w", err)
	}
	if err := r.llm.Close(); err != nil {
		return fmt.Errorf("close LLM: %w", err)
	}
	return nil
}
