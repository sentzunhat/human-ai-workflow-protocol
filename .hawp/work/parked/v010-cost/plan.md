---
work-item: v010-cost
type: feature
title: "v0.1.0: Cost Tracking + Rate Limiting (Cross-Cutting)"
status: inbox
owner: unassigned
created: 2026-07-25
updated: 2026-07-25
---

# Cost Tracking + Rate Limiting

## Mission

Add cost tracking and rate limiting to all cloud backends (OpenAI, Anthropic). Allow users to track spend, set budgets, and avoid hitting API quotas.

## Effort Estimate: 10 hours

---

## Scope

### Cost Tracking

Track spend per backend, per call, and aggregate:

```go
type CostTracker struct {
    OpenAISpend     float64 // $
    AnthropicSpend  float64 // $
    BudgetUSD       float64 // optional cap
    CallLog         []CallRecord
}

type CallRecord struct {
    Timestamp   time.Time
    Backend     string // "openai" or "anthropic"
    Operation   string // "embed" or "reshape"
    TokensUsed  int
    CostUSD     float64
}

func (ct *CostTracker) LogCall(backend, op string, tokens int, cost float64) error {
    // Add to log
    // Check if over budget (return error if yes)
}

func (ct *CostTracker) Total() float64 { ... }
func (ct *CostTracker) Remaining() float64 { ... }
func (ct *CostTracker) Reset() { ... }
```

### Rate Limiting

Respect API rate limits per backend:

```go
type RateLimiter interface {
    Wait(ctx context.Context, tokens int) error
}

type OpenAIRateLimiter struct {
    tokensPerMinute int // 3500 TPM (default)
    // Sliding window implementation
}

type AnthropicRateLimiter struct {
    tokensPerMinute int // 50000+ TPM (tier-dependent)
}

func (r *RateLimiter) Wait(ctx context.Context, tokens int) error {
    // Block until tokens are available
    // Return error if deadline exceeded
}
```

---

## Implementation Plan

### 1. Create `cost_tracker.go` (2 hours)

```go
// Structured cost tracking
type CostTracker struct {
    mu sync.RWMutex
    openaiCost float64
    anthropicCost float64
    budgetUSD float64
    history []CallRecord
}

// Configurable cost matrices per model
var OpenAICosts = map[string]struct{
    InputPer1MTokUSD float64
    OutputPer1MTokUSD float64
}{
    "text-embedding-3-small": {0.00002, 0},      // input only for embeddings
    "text-embedding-3-large": {0.00013, 0},
    "gpt-3.5-turbo": {0.0005, 0.0015},           // input, output
    "gpt-4-turbo": {0.01, 0.03},
}

var AnthropicCosts = map[string]struct{ ... }{
    "claude-3-sonnet": {0.003, 0.015},
    "claude-3-opus": {0.015, 0.075},
}

func (ct *CostTracker) EstimateCost(backend, model string, inputTokens, outputTokens int) float64 {
    // Return estimated cost before API call
}

func (ct *CostTracker) LogCost(backend, model string, inputTokens, outputTokens int) error {
    // Calculate actual cost, log it, check budget
    if ct.openaiCost + ct.anthropicCost > ct.budgetUSD {
        return fmt.Errorf("budget exceeded")
    }
}
```

### 2. Create `rate_limiter.go` (2.5 hours)

```go
type RateLimiter struct {
    backend string
    maxTPM  int  // tokens per minute
    bucket  int  // token bucket implementation
    mu      sync.Mutex
    lastRefill time.Time
}

func NewRateLimiter(backend string, maxTPM int) *RateLimiter {
    // OpenAI: 3500 TPM (default tier)
    // Anthropic: 50000+ TPM (varies by tier)
}

func (r *RateLimiter) Wait(ctx context.Context, tokens int) error {
    r.mu.Lock()
    defer r.mu.Unlock()
    
    // Token bucket: refill at start of each minute
    now := time.Now()
    if now.Sub(r.lastRefill) >= time.Minute {
        r.bucket = r.maxTPM
        r.lastRefill = now
    }
    
    // Check if tokens available
    if r.bucket >= tokens {
        r.bucket -= tokens
        return nil
    }
    
    // Wait for refill or context deadline
    timeToWait := r.nextRefillTime().Sub(now)
    select {
    case <-time.After(timeToWait):
        r.bucket -= tokens
        return nil
    case <-ctx.Done():
        return ctx.Err()
    }
}
```

### 3. Integrate with Backends (2 hours)

Wrap each backend client with cost tracking + rate limiting:

```go
type TrackedOpenAIEmbedder struct {
    embedder *OpenAIEmbedder
    tracker  *CostTracker
    limiter  *RateLimiter
}

func (t *TrackedOpenAIEmbedder) Embed(ctx context.Context, text string) ([]float32, error) {
    // 1. Estimate tokens
    tokens := estimateTokens(text)
    
    // 2. Check budget
    if t.tracker.Over Budget() {
        return nil, fmt.Errorf("cost budget exceeded")
    }
    
    // 3. Wait for rate limit
    if err := t.limiter.Wait(ctx, tokens); err != nil {
        return nil, err
    }
    
    // 4. Call actual backend
    vec, err := t.embedder.Embed(ctx, text)
    if err != nil {
        return nil, err
    }
    
    // 5. Log cost
    t.tracker.LogCost("openai", "text-embedding-3-small", tokens, 0)
    return vec, nil
}
```

### 4. Configuration (1.5 hours)

**Extended config format:**

```json
{
  "cost": {
    "budget_usd": 100.0,
    "tracking_enabled": true,
    "warn_threshold_usd": 80.0
  },
  "rate_limits": {
    "openai_tpm": 3500,
    "anthropic_tpm": 50000
  }
}
```

**Environment variables:**
```bash
HAWP_COST_BUDGET=100.0
HAWP_COST_WARN_THRESHOLD=80.0
HAWP_OPENAI_TPM=3500
HAWP_ANTHROPIC_TPM=50000
```

### 5. CLI Integration (1 hour)

New flags:
```bash
hawp search "query" \
  --cost-budget 50.0 \
  --cost-warn-threshold 40.0 \
  --show-cost-estimate
```

Output example:
```
Estimated cost: $0.045 (3 API calls)
Current spend: $12.34 (budget: $50.00)
Remaining: $37.66
```

### 6. Tests (1 hour)

**Unit Tests:**
- `TestCostTracking_Calculation` — verify cost calculations
- `TestCostTracking_BudgetCheck` — budget enforcement
- `TestRateLimiter_TokenBucket` — rate limiting logic
- `TestRateLimiter_RespectDeadline` — context cancellation

**Integration Tests (with mocks):**
- `TestTrackedOpenAI_LogsCost` — cost logged after call
- `TestTrackedOpenAI_StopsOnBudget` — error when budget hit

---

## Acceptance Criteria

- [x] CostTracker struct defined
- [x] Cost calculations accurate for all models
- [x] Budget enforcement working
- [x] RateLimiter (token bucket) working
- [x] Backends wrapped with tracking
- [x] Configuration system updated
- [x] CLI flags added
- [x] Tests passing (unit + integration)
- [x] Documentation: cost breakdown, budget setup

---

## Notes

- Cost matrix should be updatable (prices change; allow config override)
- Rate limits are tier-dependent; allow users to set their own limits
- Consider offering cost estimation before expensive operations
- Log format should be exportable (CSV/JSON) for billing/analysis
- Token counting must be accurate (use SDK tokenizers, not estimates)

---

## Future Enhancements

- Per-day/per-month budget rollover
- Cost alerts (email/Slack integration)
- Batch API for cost savings (OpenAI batch endpoint)
- Fallback strategy (if expensive model over budget, try cheaper alternative)

---

## Dependencies

- Cost tracking depends on accurate token counting from Phase 3.3c/d
- Rate limiting is independent, can be implemented in parallel

---

## Files to Create/Modify

| File | Status |
|---|---|
| `librarian/go/internal/domain/cost_tracker.go` | ✨ NEW |
| `librarian/go/internal/domain/rate_limiter.go` | ✨ NEW |
| `librarian/go/internal/domain/embeddings/openai_embedder.go` | 🔧 Wrap with tracker |
| `librarian/go/internal/domain/llm/openai_client.go` | 🔧 Wrap with tracker |
| `librarian/go/internal/domain/llm/anthropic_client.go` | 🔧 Wrap with tracker |
| `librarian/go/cmd/hawp/main.go` | 🔧 Add cost flags |
| `librarian/docs/BACKENDS.md` | 🔧 Add cost section |
| `librarian/docs/COST_TRACKING.md` | ✨ NEW (guide) |

---

**Priority:** MEDIUM (nice to have, but important for production use with cloud APIs)
**Can start:** After v0.1.0 core backends land (Phase 3.2c/d, 3.3c/d)
**Or parallel with:** All v0.1.0 phases if time allows
