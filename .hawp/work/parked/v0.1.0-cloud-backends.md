---
work-item: v0.1.0-phase-3-cloud
type: feature
title: "v0.1.0: Cloud Backends (OpenAI + Anthropic) for Embeddings & LLM"
status: inbox
owner: unassigned
created: 2026-07-25
updated: 2026-07-25
---

# v0.1.0: Cloud Backends

## Mission

Add OpenAI and Anthropic as configurable backends for embeddings and LLM inference, enabling users to trade privacy (local) for quality/cost (cloud APIs).

---

## Context

v0.0.3 ships with local-only backends (ONNX + Ollama). v0.1.0 adds cloud options:

- **Phase 3.2c:** OpenAI embeddings (`text-embedding-3-small`, `text-embedding-3-large`)
- **Phase 3.2d:** Anthropic embeddings (stub for future; API not yet available)
- **Phase 3.3c:** OpenAI LLM (`gpt-3.5-turbo`, `gpt-4-turbo`)
- **Phase 3.3d:** Anthropic LLM (`claude-3-sonnet`, `claude-3-opus`)

**Timing:** Post-v0.0.3 (parallel tracks after v0.0.3 ship).

---

## Parallelization

All four phases can run in parallel:

```
Phase 3.2c (OpenAI Embed) ─────┐
Phase 3.2d (Anthropic Embed) ──┼─> v0.1.0 Cloud Backends
Phase 3.3c (OpenAI LLM) ───────┤
Phase 3.3d (Anthropic LLM) ────┘
```

**Effort estimate:** ~1 week (4 developers) or ~3-4 weeks (1 developer)

---

## Phase 3.2c: OpenAI Embeddings

### Design

```go
type OpenAIEmbedder struct {
    client *openai.Client
    model  string // text-embedding-3-small or text-embedding-3-large
}

func NewOpenAIEmbedder(apiKey, model string) (*OpenAIEmbedder, error) {
    // Validate API key, create client
}

func (e *OpenAIEmbedder) Embed(ctx context.Context, text string) ([]float32, error) {
    // Call /v1/embeddings
}
```

### Dependencies
- `github.com/sashabaranov/go-openai` (or official OpenAI SDK when available)

### Tests
- API key validation
- Single embed
- Batch embed (leverage OpenAI batch API for cost)
- Error handling (rate limits, invalid key, model not found)
- Dimension verification

### Configuration
```json
{
  "embeddings": {
    "backend": "openai",
    "model": "text-embedding-3-small",
    "api_key": "${OPENAI_API_KEY}"
  }
}
```

---

## Phase 3.2d: Anthropic Embeddings

### Design

```go
type AnthropicEmbedder struct {
    client *anthropic.Client
    model  string // claude-embed (when available)
}
```

### Status
**Stub for future:** Anthropic does not yet have a public embeddings API. Scaffold the interface; implementation deferred until API available.

### Tests
- Interface validation
- Placeholder tests
- Ready to implement when API ships

---

## Phase 3.3c: OpenAI LLM

### Design

```go
type OpenAILLMClient struct {
    client *openai.Client
    model  string // gpt-3.5-turbo or gpt-4-turbo
    temp   float32
}

func NewOpenAILLMClient(apiKey, model string, temperature float32) (*OpenAILLMClient, error) {
    // Validate API key, create client
}

func (c *OpenAILLMClient) Reshape(ctx context.Context, packedContext string, maxTokens int) (string, error) {
    // Call /v1/chat/completions with reshaping prompt
}
```

### Features
- Support both GPT-3.5-turbo (cheap, good) and GPT-4-turbo (expensive, best)
- Rate limiting (3500 RPM per tier)
- Token counting before API calls
- Fallback strategy (gpt-3.5-turbo default, gpt-4 optional)

### Tests
- API key validation
- Single reshape
- Batch reshape (via queue, respecting rate limits)
- Token budget validation
- Error handling (rate limits, quota exceeded, invalid key)

### Configuration
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

---

## Phase 3.3d: Anthropic LLM

### Design

```go
type AnthropicLLMClient struct {
    client *anthropic.Client
    model  string // claude-3-sonnet or claude-3-opus
    temp   float32
}

func NewAnthropicLLMClient(apiKey, model string, temperature float32) (*AnthropicLLMClient, error) {
    // Validate API key, create client
}

func (c *AnthropicLLMClient) Reshape(ctx context.Context, packedContext string, maxTokens int) (string, error) {
    // Call /messages with reshaping prompt
}
```

### Features
- Support both Claude 3 Sonnet (balanced) and Opus (best quality, slower)
- Token counting via `tokenizers` lib
- Streaming support (optional optimization)

### Tests
- API key validation
- Single reshape
- Batch reshape
- Token budget validation
- Error handling (rate limits, model not found, quota)

### Configuration
```json
{
  "llm": {
    "backend": "anthropic",
    "model": "claude-3-sonnet",
    "api_key": "${ANTHROPIC_API_KEY}",
    "temperature": 0.7,
    "max_tokens": 512
  }
}
```

---

## Additional v0.1.0 Features

### Cost Tracking
Track API spend per call; optionally warn users when exceeding budgets.

```go
type CostTracker struct {
    openaiSpend      float64 // $
    anthropicSpend   float64 // $
    budgetUSD        float64 // optional cap
}
```

### Rate Limiting
Implement per-backend rate limiters to respect API quotas.

```go
type RateLimiter interface {
    Wait(ctx context.Context, count int) error
}
```

### Token Budgeting
Allow users to set max tokens per search, capping LLM output cost.

### Error Recovery
Graceful fallback if primary backend unavailable (e.g., OpenAI quota exceeded → fallback to Ollama).

---

## Acceptance Criteria (v0.1.0)

- [ ] Phase 3.2c: OpenAI embeddings working (2 models, tests passing)
- [ ] Phase 3.2d: Anthropic embeddings stubbed (interface ready, tests pass)
- [ ] Phase 3.3c: OpenAI LLM working (2 models, rate limiting, tests passing)
- [ ] Phase 3.3d: Anthropic LLM working (2 models, tests passing)
- [ ] Cost tracking integrated
- [ ] Rate limiting working
- [ ] Token budgeting enforced
- [ ] Configuration updated (env vars + JSON)
- [ ] Tests: 30+ total (cloud backends)
- [ ] Documentation: setup guide, cost breakdown, API key management
- [ ] Example: AWS Lambda deployment with cloud backends

---

## Timeline

**After v0.0.3 ships (2026-07-30):**

- Week 1 (2026-07-30 to 2026-08-06): Phase 3.2c + 3.3c (OpenAI, highest priority)
- Week 2 (2026-08-06 to 2026-08-13): Phase 3.3d (Anthropic LLM) + cost tracking
- Week 3 (2026-08-13 to 2026-08-20): Polish, testing, documentation
- **Target ship:** 2026-08-20 (3 weeks post-v0.0.3)

---

## Dependencies

- OpenAI SDK (existing or `go-openai`)
- Anthropic SDK
- Token counter library (e.g., `js-tiktoken` Go port)

---

## Future (v0.2.0+)

- Azure OpenAI support
- HuggingFace Inference API
- Replicate
- Local ONNX LLM (when models available)
- Prompt optimization / compression
- Multi-turn conversations

---

## Notes

- **Private by default:** Local backends remain the default; cloud requires explicit config
- **No vendor lock-in:** All backends behind same interface; easy to switch
- **Cost transparency:** Always show estimated/actual costs to users
- **Fallback strategy:** Can configure automatic fallback if primary backend unavailable

---

**Status:** Ready to start after v0.0.3 ships.
