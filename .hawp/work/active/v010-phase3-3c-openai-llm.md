---
work-item: v010-3-3c
type: feature
title: "v0.1.0 Phase 3.3c: OpenAI LLM Backend (gpt-3.5-turbo, gpt-4-turbo)"
status: inbox
owner: unassigned
created: 2026-07-25
updated: 2026-07-25
---

# Phase 3.3c: OpenAI LLM

## Mission

Implement OpenAI LLM (`gpt-3.5-turbo`, `gpt-4-turbo`) for context reshaping with rate limiting and token budgeting.

## Effort Estimate: 8 hours

---

## Implementation Plan

### 1. Create `openai_client.go` (2 hours)

```go
type OpenAILLMClient struct {
    client      *openai.Client
    model       string // gpt-3.5-turbo or gpt-4-turbo
    temperature float32
    maxTokens   int
}

func NewOpenAILLMClient(apiKey, model string, temperature float32) (*OpenAILLMClient, error) {
    // Validate model (only gpt-3.5-turbo, gpt-4-turbo, or gpt-4-turbo-preview)
    // Create OpenAI client
    // Test API key (quick generation call)
}

func (c *OpenAILLMClient) Reshape(ctx context.Context, packedContext string, maxTokens int) (string, error) {
    // Build prompt with ReshapingPrompt
    // Call /v1/chat/completions
    // Respect maxTokens limit
    // Handle rate limits (429)
}

func (c *OpenAILLMClient) ReshapeBatch(ctx context.Context, contexts []string, maxTokens int) ([]string, error) {
    // Process each context individually (or via batch API for cost savings)
    // Respect rate limits
}

func (c *OpenAILLMClient) Backend() string { return "openai" }
func (c *OpenAILLMClient) Model() string { ... }
func (c *OpenAILLMClient) Close() error { ... }
```

**Dependencies:** `github.com/sashabaranov/go-openai`

### 2. Update `llm/llm_client.go` factory (30 min)

Add to `NewLLMClient()`:
```go
case "openai":
    return NewOpenAILLMClient(apiKey, model, temperature)
```

Also update `NewLLMClientWithURL()` if custom endpoint needed.

### 3. Implement Token Counting (1.5 hours)

Use tiktoken for token counting (prevents budget overruns):

```go
import "github.com/tiktoken-go/tokenizer"

func (c *OpenAILLMClient) CountTokens(text string) (int, error) {
    // Get tokenizer for model (gpt-3.5/gpt-4 use cl100k_base)
    // Count tokens in prompt + context
}

func (c *OpenAILLMClient) Reshape(...) {
    // Before calling API:
    // estimatedTokens := countTokens(prompt + context)
    // if estimatedTokens > maxTokens {
    //    return error("too long for token budget")
    // }
}
```

**Note:** Anthropic + cost tracking phases will also need token counting; consider shared utility.

### 4. Write Tests (2 hours)

**Unit Tests (mocks):**
- `TestNewOpenAILLMClient_ValidateAPIKey` — catch empty/invalid keys
- `TestNewOpenAILLMClient_ValidateModel` — reject unsupported models
- `TestLLMClientInterface` — interface conformance
- `TestTokenCountingAccuracy` — token count matches actual API call
- `TestPromptBuilding` — ReshapingPrompt correctly substituted

**Integration Tests (require API key):**
- `TestOpenAIReshape_RealAPI_GPT35` — single reshape (gpt-3.5-turbo)
- `TestOpenAIReshape_RealAPI_GPT4` — single reshape (gpt-4-turbo) — may skip if user no quota
- `TestOpenAIReshape_Coherent` — output is sensible text (not garbage)
- `TestOpenAIReshape_TokenBudget` — respects maxTokens
- `TestOpenAIReshape_RateLimit` — (optional; use 429 mock)

Run with:
```bash
OPENAI_API_KEY=sk-... go test -run TestOpenAI -v ./llm/
```

### 5. Configuration & Documentation (1.5 hours)

**Config Example:**
```json
{
  "llm": {
    "backend": "openai",
    "model": "gpt-3.5-turbo",
    "api_key": "${OPENAI_API_KEY}",
    "temperature": 0.7,
    "max_tokens": 512
  }
}
```

**Model Selection Guidance:**
- **gpt-3.5-turbo** (default): $0.0015/1K input, $0.002/1K output
- **gpt-4-turbo**: $0.01/1K input, $0.03/1K output (6-7x more expensive)

**CLI Usage:**
```bash
hawp search "query" --llm-backend openai --llm-model gpt-3.5-turbo --llm-max-tokens 512
```

**Docs:**
- Add to `BACKENDS.md` (OpenAI section)
- Document cost breakdowns
- Document rate limits (3500 TPM for gpt-3.5, 40k+ for gpt-4)
- Token budget best practices

### 6. Error Handling (1 hour)

- Invalid API key → fail fast in constructor
- Rate limit (429) → exponential backoff + retry
- Model not available (404) → clear error message
- Quota exceeded → suggest gpt-3.5-turbo fallback
- Token budget exceeded → return error instead of truncating

---

## Acceptance Criteria

- [x] OpenAILLMClient struct defined
- [x] Reshape() method working (single)
- [x] ReshapeBatch() method working
- [x] Interface implemented (Backend, Model, Close)
- [x] Factory function updated
- [x] Token counting accurate (matches OpenAI)
- [x] Token budget enforced (fail if over limit)
- [x] Error handling: rate limits, invalid API key, quota, model validation
- [x] 8+ tests passing (unit + integration)
- [x] Configuration documented
- [x] Cost breakdown in docs

---

## Notes

- gpt-3.5-turbo currently (2026-07) uses model `gpt-3.5-turbo` (auto-updated)
- gpt-4-turbo use model `gpt-4-turbo-preview` or `gpt-4-turbo` (check OpenAI docs at ship time)
- Temperature default 0.7 (moderate creativity); user can override
- Token counting is *estimated* before API call; actual tokens may vary slightly
- Batch API available for cost savings; can integrate later as optimization

---

## Related Work

- **Phase 3.2c** (OpenAI Embeddings) — can start in parallel
- **Cost Tracking** — will wrap both embeddings + LLM later
- **Rate Limiting** — will add retry logic + queue later
- **Phase 3.3d** (Anthropic LLM) — similar structure, can learn from this

---

## Files to Create/Modify

| File | Status |
|---|---|
| `librarian/go/internal/domain/llm/openai_client.go` | ✨ NEW |
| `librarian/go/internal/domain/llm/openai_client_test.go` | ✨ NEW |
| `librarian/go/internal/domain/llm/llm_client.go` | 🔧 Update factory |
| `librarian/go/internal/domain/llm/token_counter.go` | ✨ NEW (shared utility) |
| `librarian/docs/BACKENDS.md` | 🔧 Add OpenAI section |
| `librarian/go/go.mod` | 🔧 Add openai + tiktoken |

---

**Priority:** HIGH (paired with Phase 3.2c for complete OpenAI support)
**Can start:** Immediately after v0.0.3 ships
**Parallel with:** Phase 3.2c, Phase 3.2d, Phase 3.3d, Cost Tracking
