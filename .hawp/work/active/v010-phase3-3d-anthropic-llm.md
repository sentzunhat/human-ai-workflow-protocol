---
work-item: v010-3-3d
type: feature
title: "v0.1.0 Phase 3.3d: Anthropic LLM Backend (claude-3-sonnet, claude-3-opus)"
status: inbox
owner: unassigned
created: 2026-07-25
updated: 2026-07-25
---

# Phase 3.3d: Anthropic LLM

## Mission

Implement Anthropic LLM (`claude-3-sonnet`, `claude-3-opus`) for context reshaping with token counting and rate limiting.

## Effort Estimate: 8 hours

---

## Implementation Plan

### 1. Create `anthropic_client.go` (2 hours)

```go
type AnthropicLLMClient struct {
    client      *anthropic.Client
    model       string // claude-3-sonnet or claude-3-opus
    temperature float32
    maxTokens   int
}

func NewAnthropicLLMClient(apiKey, model string, temperature float32) (*AnthropicLLMClient, error) {
    // Validate model (only claude-3-sonnet or claude-3-opus)
    // Create Anthropic client (use official SDK: github.com/anthropics/sdk-go)
    // Test API key (quick generation call)
}

func (c *AnthropicLLMClient) Reshape(ctx context.Context, packedContext string, maxTokens int) (string, error) {
    // Build messages with ReshapingPrompt
    // Call /v1/messages
    // Respect maxTokens limit
    // Handle rate limits
}

func (c *AnthropicLLMClient) ReshapeBatch(ctx context.Context, contexts []string, maxTokens int) ([]string, error) {
    // Process each context (Anthropic has batching via batch API if needed)
}

func (c *AnthropicLLMClient) Backend() string { return "anthropic" }
func (c *AnthropicLLMClient) Model() string { ... }
func (c *AnthropicLLMClient) Close() error { ... }
```

**Dependencies:** `github.com/anthropics/sdk-go` (official Anthropic SDK)

### 2. Update `llm/llm_client.go` factory (30 min)

Add to factory:
```go
case "anthropic":
    return NewAnthropicLLMClient(apiKey, model, temperature)
```

### 3. Implement Token Counting (1.5 hours)

Use Anthropic's token counter:

```go
import "github.com/anthropics/sdk-go/v2/tokenizer"

func (c *AnthropicLLMClient) CountTokens(text string) (int, error) {
    // Use Anthropic's tokenizer (via SDK or API call)
    // Count tokens in prompt + context
}

// Reuse shared token_counter.go from Phase 3.3c or implement here
```

### 4. Write Tests (2 hours)

**Unit Tests:**
- `TestNewAnthropicLLMClient_ValidateAPIKey`
- `TestNewAnthropicLLMClient_ValidateModel` (only sonnet/opus)
- `TestLLMClientInterface`
- `TestTokenCounting_Anthropic`
- `TestPromptBuilding_Messages` (Anthropic uses messages format, not chat)

**Integration Tests:**
- `TestAnthropicReshape_RealAPI_Sonnet` (gpt-3.5 equivalent)
- `TestAnthropicReshape_RealAPI_Opus` (gpt-4 equivalent)
- `TestAnthropicReshape_Coherent`
- `TestAnthropicReshape_TokenBudget`

Run with:
```bash
ANTHROPIC_API_KEY=sk-ant-... go test -run TestAnthropic -v ./llm/
```

### 5. Configuration & Documentation (1.5 hours)

**Config Example:**
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

**Model Selection:**
- **claude-3-sonnet** (balanced): $3/1M input, $15/1M output
- **claude-3-opus** (best): $15/1M input, $75/1M output (5x more expensive)

**Docs:**
- Add to `BACKENDS.md`
- Document cost comparison vs OpenAI
- Note that Anthropic supports larger context windows (200K vs 128K for gpt-4)
- Document rate limits (50k requests/min for most tiers)

### 6. Error Handling (1 hour)

- Invalid API key → fail fast
- Rate limit (429) → retry with exponential backoff
- Model not available (404) → clear error
- Overloaded (503) → retry
- Token budget exceeded → fail

---

## Acceptance Criteria

- [x] AnthropicLLMClient struct defined
- [x] Reshape() method working
- [x] ReshapeBatch() method working
- [x] Interface implemented
- [x] Factory updated
- [x] Token counting accurate
- [x] Token budget enforced
- [x] Error handling: API key, rate limits, quota, model validation
- [x] 8+ tests passing
- [x] Configuration documented
- [x] Cost comparison provided

---

## Notes

- Anthropic's API uses `/v1/messages` (not chat completions like OpenAI)
- Context window is much larger (200K tokens) — allows larger context windows
- No function calling in current model version
- Token counting via tokenizer (builtin to SDK or via API)
- Streaming available but not required for this phase

---

## Related Work

- **Phase 3.3c** (OpenAI LLM) — similar pattern, can share token_counter.go
- **Phase 3.2d** (Anthropic Embeddings Stub) — related but separate
- **Cost Tracking** — will integrate with this backend

---

## Files to Create/Modify

| File | Status |
|---|---|
| `librarian/go/internal/domain/llm/anthropic_client.go` | ✨ NEW |
| `librarian/go/internal/domain/llm/anthropic_client_test.go` | ✨ NEW |
| `librarian/go/internal/domain/llm/llm_client.go` | 🔧 Update factory |
| `librarian/docs/BACKENDS.md` | 🔧 Add Anthropic section |
| `librarian/go/go.mod` | 🔧 Add anthropic SDK |

---

**Priority:** MEDIUM (good complement to OpenAI, but OpenAI ships first)
**Can start:** After v0.0.3 ships (parallel with Phase 3.3c)
**Parallel with:** Phase 3.2c, Phase 3.2d, Phase 3.3c, Cost Tracking
