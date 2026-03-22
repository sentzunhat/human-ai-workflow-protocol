---
work-item: v010-3-2c
type: feature
title: "v0.1.0 Phase 3.2c: OpenAI Embeddings Backend"
status: inbox
owner: unassigned
created: 2026-07-25
updated: 2026-07-25
---

# Phase 3.2c: OpenAI Embeddings

## Mission

Implement OpenAI embeddings (`text-embedding-3-small`, `text-embedding-3-large`) as a configurable backend.

## Effort Estimate: 6 hours

---

## Implementation Plan

### 1. Create `openai_embedder.go` (1.5 hours)

```go
type OpenAIEmbedder struct {
    client *openai.Client
    model  string // text-embedding-3-small or text-embedding-3-large
}

func NewOpenAIEmbedder(apiKey, model string) (*OpenAIEmbedder, error) {
    // Validate model (only support 3-small, 3-large)
    // Create OpenAI client
    // Test API key (quick embedding call)
}

func (e *OpenAIEmbedder) Embed(ctx context.Context, text string) ([]float32, error) {
    // Call /v1/embeddings
    // Handle rate limits (429)
}

func (e *OpenAIEmbedder) EmbedBatch(ctx context.Context, texts []string) ([][]float32, error) {
    // Batch up to 100 texts per call (OpenAI limit)
    // Handle rate limits
}

func (e *OpenAIEmbedder) Dimension() int { ... } // 1536 for both models
func (e *OpenAIEmbedder) Backend() string { return "openai" }
func (e *OpenAIEmbedder) Model() string { ... }
func (e *OpenAIEmbedder) Close() error { ... }
```

**Dependencies:** `github.com/sashabaranov/go-openai`

### 2. Update `embeddings/embedder.go` factory (30 min)

Add to `NewEmbedder()` factory:
```go
case "openai":
    return NewOpenAIEmbedder(apiKey, model)
```

### 3. Write Tests (2 hours)

**Unit Tests (no API calls, mocks only):**
- `TestNewOpenAIEmbedder_ValidateAPIKey` — catch empty/invalid keys early
- `TestNewOpenAIEmbedder_ValidateModel` — reject unsupported models
- `TestEmbedderInterface` — interface conformance
- `TestBatchSplitting` — texts > 100 split correctly

**Integration Tests (require API key, skip if unavailable):**
- `TestOpenAIEmbed_RealAPI` — single embed against real API
- `TestOpenAIEmbed_Dimension` — verify 1536-dim output
- `TestOpenAIEmbed_Similarity` — semantic signal (like Ollama tests)
- `TestOpenAIEmbed_RateLimitHandling` — (optional; use 429 mock)

Run with:
```bash
OPENAI_API_KEY=sk-... go test -run TestOpenAI -v ./embeddings/
```

### 4. Configuration & Documentation (1 hour)

**Config Example:**
```json
{
  "embeddings": {
    "backend": "openai",
    "model": "text-embedding-3-small",
    "api_key": "${OPENAI_API_KEY}"
  }
}
```

**CLI Usage:**
```bash
hawp search "query" --embeddings-backend openai --embeddings-model text-embedding-3-small
```

**Docs:**
- Add to `BACKENDS.md` (OpenAI section)
- Document models: small (cheap) vs large (better)
- Document costs: ~$0.02 per 1M tokens
- Document rate limits: 3500 TPM (tier-dependent)

### 5. Code Review Checklist

- [ ] Implements Embedder interface
- [ ] API key validation in constructor (fail fast)
- [ ] Batch size split (max 100 texts per call)
- [ ] Error handling (rate limits, invalid key, model not found)
- [ ] Tests pass (unit + integration)
- [ ] No hardcoded API keys in tests (use env vars)
- [ ] Documentation matches code

---

## Acceptance Criteria

- [x] OpenAIEmbedder struct defined
- [x] Embed() method working (single)
- [x] EmbedBatch() method working (respects 100-text limit)
- [x] Interface implemented (Dimension, Backend, Model, Close)
- [x] Factory function updated
- [x] 8+ tests passing (unit + integration)
- [x] Configuration documented
- [x] Error handling: rate limits, invalid API key, model validation

---

## Notes

- Model dimension is always 1536 (both 3-small and 3-large)
- 3-small is ~13x cheaper than 3-large but quality is comparable for most tasks
- Rate limits are per API key; may need retry logic for production
- Consider adding rate limiter (Phase 3.2c+) for multi-user scenarios

---

## Related Work

- **Phase 3.3c** (OpenAI LLM) — can start in parallel
- **Cost Tracking** — will integrate with this backend later
- **Rate Limiting** — will wrap this backend later

---

## Files to Create/Modify

| File | Status |
|---|---|
| `librarian/go/internal/domain/embeddings/openai_embedder.go` | ✨ NEW |
| `librarian/go/internal/domain/embeddings/openai_embedder_test.go` | ✨ NEW |
| `librarian/go/internal/domain/embeddings/embedder.go` | 🔧 Update factory |
| `librarian/docs/BACKENDS.md` | 🔧 Add OpenAI section |
| `librarian/go/go.mod` | 🔧 Add openai dependency |

---

**Priority:** HIGH (most users want OpenAI first)
**Can start:** Immediately after v0.0.3 ships
**Parallel with:** Phase 3.3c, Phase 3.2d, Phase 3.3d, Cost Tracking
