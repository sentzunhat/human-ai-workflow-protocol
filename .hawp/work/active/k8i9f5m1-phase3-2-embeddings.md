---
work-item: k8i9f5m1-p3-2
type: feature
title: "v0.0.3 Phase 3.2: Embedding Backends (ONNX, Ollama, OpenAI, Anthropic)"
status: in-progress
owner: unassigned
created: 2026-07-24
updated: 2026-07-25
---

# Phase 3.2: Embedding Backends

## Mission

Implement 4 configurable embedding backends so users can choose their embedding source: local ONNX (default), Ollama, OpenAI, or Anthropic.

**Can run in parallel once Phase 3.1 (Config System) is complete.**

---

## Context

**Background:**
- Phase 3.1 provides config system telling us which backend to use
- Phase 3.2 implements the actual backends
- Embeddings are used in Phase 3.4 to identify key concepts in packed context

**Embedding Models (Recommended 2 for Context Reshaping):**

1. **BGE-base-en-v1.5** (768-dim) — Primary
   - Semantic quality: 95%+ on MTEB benchmarks
   - Speed: 110s for full repo (2,445 chunks)
   - Size: Reasonable for local use
   - Use case: Identify key semantic concepts

2. **all-MiniLM-L6-v2** (384-dim) — Secondary
   - Semantic quality: 90%+ (good, slightly lower)
   - Speed: 140s for full repo
   - Size: Lighter than BGE
   - Use case: Quick re-ranking, fallback for resource-constrained

---

## Design

### Embedding Backend Interface

```go
// Embedder is the interface all embedding backends implement
type Embedder interface {
    // Embed returns a vector for the given text
    Embed(ctx context.Context, text string) ([]float32, error)
    
    // EmbedBatch returns vectors for multiple texts (more efficient)
    EmbedBatch(ctx context.Context, texts []string) ([][]float32, error)
    
    // Dimension returns the vector dimension (768 for BGE, 384 for MiniLM)
    Dimension() int
    
    // Backend returns the backend name for logging
    Backend() string
    
    // Close releases resources
    Close() error
}

// Factory function to create embedders based on config
func NewEmbedder(cfg EmbeddingsConfig) (Embedder, error)
```

### Backend Implementations

#### Phase 3.2a: ONNX Embeddings (Primary, ~1-1.5 days)

```go
type ONNXEmbedder struct {
    session *ort.Session
    model   string // "bge-base-en-v1.5" or "all-MiniLM-L6-v2"
}

func NewONNXEmbedder(model string) (*ONNXEmbedder, error) {
    // 1. Check if model is downloaded (~/.hawp/models/)
    // 2. If not, download via hugot (like `hawp model pull`)
    // 3. Load ONNX Runtime session
    // 4. Return embedder
}

func (e *ONNXEmbedder) Embed(ctx context.Context, text string) ([]float32, error) {
    // 1. Tokenize text
    // 2. Run ONNX inference
    // 3. Return embeddings
}

func (e *ONNXEmbedder) EmbedBatch(ctx context.Context, texts []string) ([][]float32, error) {
    // Batch inference for efficiency
}
```

**Dependencies:** Use existing hugot integration (already in v0.0.1)

**Tests:**
- Load model correctly (or download on first use)
- Embed single text
- Embed batch of texts
- Verify dimension (768 for BGE, 384 for MiniLM)
- Handle invalid models gracefully

---

#### Phase 3.2b: Ollama Embeddings (~8 hours, parallel)

```go
type OllamaEmbedder struct {
    client *http.Client
    url    string // "http://localhost:11434"
    model  string // "nomic-embed-text", "mxbai-embed-large", etc.
}

func NewOllamaEmbedder(url, model string) (*OllamaEmbedder, error) {
    // 1. Verify Ollama is running at url
    // 2. Check if model is available
    // 3. Return embedder
}

func (e *OllamaEmbedder) Embed(ctx context.Context, text string) ([]float32, error) {
    // POST /api/embed with text
    // Parse response, return embeddings
}

func (e *OllamaEmbedder) EmbedBatch(ctx context.Context, texts []string) ([][]float32, error) {
    // Batch via Ollama API
}
```

**Tests:**
- Verify Ollama running (with mock server)
- Embed single text via HTTP
- Embed batch via HTTP
- Handle connection errors gracefully

---

#### Phase 3.2c: OpenAI Embeddings (~6 hours, parallel)

```go
type OpenAIEmbedder struct {
    client *openai.Client
    model  string // "text-embedding-3-small" (default), "text-embedding-3-large"
}

func NewOpenAIEmbedder(apiKey string, model string) (*OpenAIEmbedder, error) {
    // 1. Create OpenAI client with API key
    // 2. Verify API key (optional: test embed)
    // 3. Return embedder
}

func (e *OpenAIEmbedder) Embed(ctx context.Context, text string) ([]float32, error) {
    // Call OpenAI /v1/embeddings
    // Rate limiting: max 3500 requests/min
    // Return embeddings
}

func (e *OpenAIEmbedder) EmbedBatch(ctx context.Context, texts []string) ([][]float32, error) {
    // Batch up to 2048 texts per API call
    // Handle rate limiting
}
```

**Dependencies:** `github.com/sashabaranov/go-openai` (or OpenAI's official SDK)

**Tests:**
- Verify API key (mock)
- Embed single text
- Embed batch (test rate limiting)
- Handle API errors

---

#### Phase 3.2d: Anthropic Embeddings (~6 hours, parallel)

```go
type AnthropicEmbedder struct {
    client *anthropic.Client
    model  string // "claude-embed" (when available)
}

func NewAnthropicEmbedder(apiKey string, model string) (*AnthropicEmbedder, error) {
    // 1. Create Anthropic client
    // 2. Verify API key
    // 3. Return embedder
}

func (e *AnthropicEmbedder) Embed(ctx context.Context, text string) ([]float32, error) {
    // Call Anthropic embedding API (when available)
}
```

**Note:** Anthropic doesn't currently have a public embeddings API, but we stub it for future use.

---

## Parallelization Strategy

**Week 1 (after Phase 3.1):**
- 3.2a (ONNX): Primary developer, 1-1.5 days
- 3.2b (Ollama): Second developer, ~8 hours
- 3.2c (OpenAI): Third developer, ~6 hours
- 3.2d (Anthropic): Can be stubbed/mocked, ~6 hours

All 4 run in parallel (4 developers) or sequentially (1 developer, ~2 days).

**Dependencies between 3.2a/b/c/d:** None. Each is independent once Phase 3.1 is done.

---

## Acceptance Criteria

- [x] ONNX embedder with BGE-base-en-v1.5 + all-MiniLM-L6-v2 support — real inference tested ✅
- [x] Ollama embedder (local or remote API) — HTTP wiring done, **live tested against ollama serve** ✅
- [ ] OpenAI embedder (text-embedding-3-small as default) — deferred to v0.1.0
- [ ] Anthropic embedder (stubbed for future use) — deferred to v0.1.0
- [x] All implement Embedder interface consistently
- [x] Config system (Phase 3.1) correctly routes to chosen backend
- [x] Tests: 20+ unit/mock tests passing; integration tests with real Ollama passing ✅
- [x] Documentation: supported models + expected dimensions

**Evidence:** `integration_test.go` — 3 tests for Ollama embeddings (single, similarity, batch), all passing against real `ollama serve`. Bug fixes for API field names. Performance: 24-29ms per embed (Ollama HTTP) vs 119-605ms (ONNX, model-dependent).

---

## Implementation Steps (Per Backend)

1. **Define backend struct** (30 min)
2. **Implement Embedder interface** (1 hour)
3. **API communication** (1-1.5 hours per backend)
4. **Error handling** (45 min)
5. **Tests** (1.5 hours)
6. **Documentation** (30 min)

---

## Effort Estimate

| Backend | Effort | Can Parallelize |
|---------|--------|-----------------|
| ONNX (primary) | 1.5 days | Yes (start first) |
| Ollama | 8 hours | Yes (after Phase 3.1) |
| OpenAI | 6 hours | Yes (after Phase 3.1) |
| Anthropic | 6 hours | Yes (after Phase 3.1) |
| **Total** | 2-3 days | **2 days parallel** |

---

## Dependencies

- ✅ Phase 3.1 (Config System) must be complete
- ✅ hugot already integrated (ONNX models)
- Optional: OpenAI SDK, Anthropic SDK (add as dependencies)

---

## Unblocks

- Phase 3.3: LLM backends (embeddings + LLM together in Phase 3.4)
- Phase 3.4: Context reshaping (uses embeddings from 3.2)

---

## Success Metrics

✅ **All backends work:**
- ONNX: Embed texts using local models
- Ollama: Embed texts using local/remote Ollama API
- OpenAI: Embed texts using OpenAI API (with rate limiting)
- Anthropic: Stub ready for future use

✅ **Consistent interface:**
- All backends implement Embedder interface
- Dimension reported correctly
- Batch and single operations work

✅ **Ready for Phase 3.4:**
- Config system can select any backend
- Embeddings ready to use for context reshaping

---

## Notes

- **Model Selection:** BGE-base-en-v1.5 is best for semantic search + context analysis
- **Dimension Mismatch:** When users switch backends, output dimensions may change (768 vs 384 vs 1536) — Phase 3.4 handles this transparently
- **Rate Limiting:** OpenAI has rate limits (3500 req/min), batch calls to reduce calls
- **Future:** If Anthropic adds embeddings API, implementation is straightforward (same pattern)

---

**Status:** Ready to implement (after Phase 3.1).

Can run in parallel: 4 backends = 4 developers, or 1 developer spending 2 days sequentially.
