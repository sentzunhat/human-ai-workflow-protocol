---
work-item: i6g8d3k9
type: feature
title: "Context Packing (Slice 4): Search results → LLM-ready context"
status: inbox
owner: unassigned
created: 2026-07-23
updated: 2026-07-23
---

# Slice 4: Context Packing

## Mission

Convert raw search results into structured, LLM-ready context blocks that can be directly used for prompt injection into Claude/OpenAI/other LLM APIs.

This is the "magic" that makes search useful for agentic loops: users ask questions, search finds relevant code/docs, and results are automatically formatted as context for LLMs to reason about.

---

## Context

**What Works (Slices 1-3):**
- ✅ Indexing 2,445 document chunks (HAWP kit + work docs)
- ✅ Lexical search (FTS5): <1ms, 70% quality
- ✅ Semantic search (ONNX): 100ms, 95% quality
- ✅ Hybrid ranking: 15-20ms, 96% quality

**What's Missing (This Slice):**
- ❌ Formatting search results for LLM consumption
- ❌ Context packing (combining results into coherent blocks)
- ❌ Integration with LLM APIs (OpenAI, Anthropic, etc.)
- ❌ Agentic loops (iterative refinement of queries)

**This Slice Delivers:**
- ✅ Context packing: Multi-result synthesis
- ✅ LLM-ready format: Structured markdown/JSON
- ✅ Integration points: OpenAI, Anthropic APIs
- ✅ Foundation for agentic loops in v0.2.0

---

## Design

### Context Packing Pipeline

```
Search Results (N chunks)
    ↓
Deduplication (remove similar results)
    ↓
Grouping (by topic/document)
    ↓
Ordering (most relevant first)
    ↓
Formatting (markdown/JSON)
    ↓
Truncation (fit within token budget)
    ↓
LLM-Ready Context Block
```

### Data Structures

```go
// ContextBlock is a formatted search result ready for LLM injection
type ContextBlock struct {
    Title        string            // e.g., "Search Results for 'vector embedding'"
    ResultCount  int               // Number of results included
    TokenCount   int               // Approximate token usage
    Results      []FormattedResult // Ordered, deduplicated results
    Metadata     map[string]string // Query, timestamp, filters
}

// FormattedResult is one result formatted for readability
type FormattedResult struct {
    Rank        int       // Position (1, 2, 3, ...)
    Relevance   float64   // Confidence (0.0 - 1.0)
    Source      string    // Document source/path
    Title       string    // Chunk title/context
    Content     string    // Actual text content
    TokenCount  int       // Approximate tokens in this result
}
```

### Example Output (Markdown Format)

```markdown
# Search Results: "vector embedding ONNX"

**Results:** 5 chunks | **Tokens:** ~1,200

## Result 1: BGE Model Integration (96% relevant)
**Source:** librarian/go/internal/domain/embed-service.go:45-67
**Path:** .hawp/kit/patterns/vector-embedding.md

The embedding service uses ONNX Runtime to load pre-quantized models.
BGE-base-en-v1.5 (768-dim, 110s embedding time) is the default...

## Result 2: Embedding Pipeline (94% relevant)
**Source:** benchmark/queries.go:128-145
...

## Context Summary
- Total chunks: 5
- Quality: High (avg 94%)
- Topics: embedding, onnx, performance, indexing
```

### Example Output (JSON Format)

```json
{
  "title": "Search Results for 'vector embedding ONNX'",
  "query": "vector embedding ONNX",
  "result_count": 5,
  "token_count": 1200,
  "results": [
    {
      "rank": 1,
      "relevance": 0.96,
      "source": "librarian/go/internal/domain/embed-service.go",
      "content": "The embedding service uses ONNX Runtime...",
      "tokens": 250
    },
    ...
  ]
}
```

---

## Implementation Plan

### Phase 1: Deduplication & Grouping (1 day)

**Goal:** Remove similar results, group by source

```go
// Deduplication: remove near-duplicates (cosine > 0.95)
func DeduplicateResults(results []SearchResult) []SearchResult {
    // Compare embeddings, filter out duplicates
}

// Grouping: organize by source document
func GroupBySource(results []SearchResult) map[string][]SearchResult {
    // Group by filename, path, or document ID
}
```

**Deliverable:** `context_dedup.go`

### Phase 2: Formatting & Truncation (1 day)

**Goal:** Format results for LLM consumption, respect token budgets

```go
// Format: Convert raw results to readable markdown/JSON
func FormatAsMarkdown(results []SearchResult, maxTokens int) ContextBlock {
    // Format results, truncate if necessary
}

// TokenCounter: Estimate tokens in formatted output
func EstimateTokens(text string) int {
    // Rough estimate: ~4 chars per token
}
```

**Deliverable:** `context_format.go`, `token_counter.go`

### Phase 3: Context Reshaping via Configurable LLM + Embeddings (DEFERRED TO v0.0.3, ~2-3 days)

**Goal:** Reshape packed context using LLM inference + configurable embeddings for better LLM consumption

**Architecture:**

1. **Embedding Backends (configurable):**
   - Primary: ONNX local (fast, private, no cost)
     - Model 1: BGE-base-en-v1.5 (768-dim, 110s) — best for semantic search
     - Model 2: all-MiniLM-L6-v2 (384-dim, 140s) — lighter, good for mobile
   - Optional: Ollama (local or remote embedding API)
   - Optional: OpenAI (text-embedding-3-small or text-embedding-3-large)
   - Optional: Anthropic (when available, via API)

2. **LLM Inference Backends (configurable):**
   - Primary: ONNX local inference (fast, private, no cost)
   - Optional: Ollama (local LLM API)
   - Optional: OpenAI ChatGPT API
   - Optional: Anthropic Claude API

3. **Context Reshaping Pipeline:**
   - Step 1: Embed packed context chunks (using selected embedding backend)
   - Step 2: Identify key concepts via embedding similarity
   - Step 3: Send reshaped context to LLM for final restructuring
   - Step 4: Return improved context (better prioritization + readability)

4. **Secure Configuration & Key Handling:**
   - Keys NOT saved/available by default (user opt-in)
   - Config sources (in order):
     1. Environment file: `.env` (user-created, gitignored)
     2. Context config: `~/.hawp/config/context.json` (base64-encoded keys)
     3. Runtime: Keys decrypted only when sending, never logged/cached
   - Encryption: Base64 initially, AES-256 later
   - Per-backend config: embeddings backend ≠ LLM backend (independent)

```go
// Fully configurable backends
type ContextReshapingConfig struct {
    // Embeddings backend
    EmbeddingsBackend string // "onnx" (default) | "ollama" | "openai" | "anthropic"
    EmbeddingsModel   string // "bge-base-en-v1.5" (default) or user choice
    EmbeddingsURL     string // For Ollama: http://localhost:11434
    
    // LLM backend
    LLMBackend        string // "onnx" (default) | "ollama" | "openai" | "anthropic"
    LLMModel          string // Model name per backend
    LLMTemperature    float32
    LLMMaxTokens      int
    
    // Secure key storage
    OpenAIKey         string // (encrypted in config)
    AnthropicKey      string // (encrypted in config)
    OllamaURL         string // (for remote Ollama)
    
    // Output budget
    MaxTokens         int    // Total output token budget
}

// Reshaping pipeline (step 1: embeddings)
func ReshapeContextWithEmbeddings(ctx ContextBlock, cfg ContextReshapingConfig) ContextBlock {
    embedder := SelectEmbedder(cfg.EmbeddingsBackend, cfg.EmbeddingsModel)
    // 1. Embed all chunks in packed context
    // 2. Identify key concepts via cosine similarity
    // 3. Re-prioritize results by semantic importance
    // 4. Restructure headings based on concept hierarchy
    // 5. Return intermediate reshaped context
}

// Reshaping pipeline (step 2: LLM)
func ReshapeContextWithLLM(ctx ContextBlock, cfg ContextReshapingConfig) ContextBlock {
    llm := SelectLLM(cfg.LLMBackend, cfg.LLMModel)
    // 1. Send intermediate context to LLM
    // 2. LLM refines structure, removes redundancy, adds clarity
    // 3. Return final reshaped context
}
```

**Deliverables:** 
- `context_reshape.go` (embedding + LLM pipeline)
- `config_system.go` (configurable backends + secure storage)
- `embedding_backends.go` (ONNX, Ollama, OpenAI, Anthropic loaders)
- `llm_backends.go` (ONNX, Ollama, OpenAI, Anthropic clients)
- Tests: 15+ for config, 10+ for embeddings, 10+ for LLM, 5+ for reshaping

### Phase 4: CLI Integration (1 day)

**Goal:** Expose context packing via CLI

```bash
# Search and return LLM-ready context
hawp search --context "vector embedding" --format markdown --max-tokens 2000

# Output context to file
hawp search --context "query" --output context.md

# JSON format for API integration
hawp search --context "query" --format json
```

**Deliverable:** CLI command updates

### Phase 5: Testing & Documentation (1 day)

**Goal:** Comprehensive tests and usage guides

- Unit tests for deduplication, formatting, truncation
- Integration tests with mock LLM APIs
- Documentation: how to use context packing
- Examples: Claude, ChatGPT, Claude Code integration

**Deliverable:** Tests, docs, examples

---

## Acceptance Criteria

### Definition of Done

- [ ] Deduplication removes >90% of near-duplicates
- [ ] Formatting preserves all relevant information
- [ ] Token counting accurate to within 10%
- [ ] Truncation respects token budget
- [ ] OpenAI integration tested with real API calls
- [ ] Anthropic integration tested with real API calls
- [ ] CLI commands work end-to-end
- [ ] All tests pass (18+ existing + 20+ new)
- [ ] Documentation complete with examples

---

## Effort Estimate

| Phase | Work | Effort |
|-------|------|--------|
| Deduplication & Grouping | Design + implement | 1 day |
| Formatting & Truncation | Design + implement | 1 day |
| LLM API Integration | OpenAI + Anthropic + generic | 1.5 days |
| CLI Integration | Wire up commands | 1 day |
| Testing & Documentation | Tests + docs + examples | 1 day |
| **Total** | | **5-6 days** |

---

## Dependencies

- ✅ Search system complete (Slices 1-3)
- ✅ All tests passing
- ✅ Binary version 0.0.1+ released
- ⏳ v0.0.1 released (current work)

---

## Success Metrics

✅ **Context Packing Works:**
- Search results → deduplicated
- Deduplicated → formatted (markdown/JSON)
- Formatted → truncated (token budget respected)
- Truncated → ready for LLM API

✅ **LLM Integration Works:**
- Claude API: context injected correctly
- OpenAI API: system prompt formatted correctly
- Results used in real LLM calls

✅ **User Value:**
- Users run: `hawp search --context "query"`
- Get back: LLM-ready context in 50ms
- Can paste directly into ChatGPT/Claude
- Foundation for agentic loops

---

## Timeline

**Week 1 (v0.0.2):**
- ✅ v0.0.1 released (completed 2026-07-24)
- ✅ Phase 1: Deduplication (completed 2026-07-23)
- ✅ Phase 2: Formatting (completed 2026-07-23)

**Week 2 (v0.0.2) — TODAY:**
- 🔄 Phase 4: CLI Integration (1 day) — **STARTING NOW**
- 🔄 Phase 5: Testing & Docs (1 day)
- 📦 **Tag v0.0.2 with context packing (Fri 2026-07-26)**
- ✅ Test auto-update from v0.0.1 → v0.0.2

**Week 3 (v0.0.3):**
- 📝 Phase 3: Context Reshaping via LLM (2 days)
  - ONNX embeddings (primary, local)
  - Ollama/ChatGPT/Anthropic backends (configurable)
  - Secure API key storage + config system
- ✅ Testing & docs
- 📦 **Tag v0.0.3 with context reshaping**

**After v0.0.3:**
- v0.1.0: Agentic loops (iterative refinement)
- v0.1.0: Provider overlays (IDE integration)

---

## Notes

- **Foundation:** Search system is proven and production-ready
- **Value:** Context packing makes search useful for LLMs
- **Extension:** Agentic loops build on context packing
- **API-agnostic:** Works with Claude, ChatGPT, open-source LLMs
- **User-friendly:** Simple CLI, no configuration needed

---

**Status:** ✅ READY TO PLAN → IMPLEMENT

This is Slice 4 of the core feature set. Highest priority after v0.0.1 release.
