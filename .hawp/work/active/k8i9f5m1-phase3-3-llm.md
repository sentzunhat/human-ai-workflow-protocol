---
work-item: k8i9f5m1-p3-3
type: feature
title: "v0.0.3 Phase 3.3: LLM Backends (ONNX, Ollama, OpenAI, Anthropic)"
status: in-progress
owner: unassigned
created: 2026-07-24
updated: 2026-07-25
---

# Phase 3.3: LLM Backends

## Mission

Implement 4 configurable LLM backends so users can choose their LLM for context reshaping: local ONNX (default), Ollama, OpenAI, or Anthropic.

**Can run in parallel with Phase 3.2 (Embedding Backends) once Phase 3.1 (Config System) is complete.**

---

## Context

**Background:**
- Phase 3.1 provides config system
- Phase 3.2 implements embedding backends (parallel track)
- Phase 3.3 implements LLM backends (this, parallel track)
- Phase 3.4 combines both for context reshaping

**LLM Models (Recommended for Context Reshaping):**

1. **ONNX Local (Primary)**
   - Option A: TinyLlama (1.1B params) — fastest, lighter
   - Option B: Mistral-7B (7B params, quantized) — better quality
   - Trade-off: Speed vs quality; both run locally

2. **Ollama (Secondary)**
   - User chooses model: mistral, llama2, neural-chat, etc.
   - Local execution (like ONNX but via Ollama server)

3. **OpenAI (Premium)**
   - gpt-4-turbo (best quality)
   - gpt-3.5-turbo (cheaper, good quality)
   - User's choice

4. **Anthropic (Premium)**
   - claude-3-sonnet (balanced speed/quality)
   - claude-3-opus (best quality, slower)
   - User's choice

---

## Design

### LLM Backend Interface

```go
// LLMClient is the interface all LLM backends implement
type LLMClient interface {
    // Reshape takes packed context and returns reshaped context
    Reshape(ctx context.Context, packedContext string, maxTokens int) (string, error)
    
    // ReshapeBatch processes multiple contexts efficiently
    ReshapeBatch(ctx context.Context, contexts []string, maxTokens int) ([]string, error)
    
    // Backend returns the backend name for logging
    Backend() string
    
    // Model returns the model name in use
    Model() string
    
    // Close releases resources
    Close() error
}

// Factory function
func NewLLMClient(cfg LLMConfig) (LLMClient, error)
```

### Backend Implementations

#### Phase 3.3a: ONNX LLM (Primary, ~1.5-2 days)

```go
type ONNXLLMClient struct {
    session *ort.Session
    model   string // "TinyLlama-1.1B" or "Mistral-7B-v0.1-Q4"
    temp    float32
    maxToks int
}

func NewONNXLLMClient(model string, temperature float32, maxTokens int) (*ONNXLLMClient, error) {
    // 1. Check if model is downloaded (~/.hawp/models/)
    // 2. If not, download via hugot
    // 3. Load ONNX Runtime session
    // 4. Return client
}

func (c *ONNXLLMClient) Reshape(ctx context.Context, packedContext string, maxTokens int) (string, error) {
    // 1. Prepare prompt: "Reshape this context for clarity: {packedContext}"
    // 2. Tokenize
    // 3. Run inference
    // 4. Decode output
    // 5. Return reshaped context
}
```

**Challenge:** ONNX LLM requires quantized models (fit in memory). We'll start with TinyLlama (smallest), offer Mistral-7B-Q4 as option.

**Dependencies:** hugot (already integrated)

---

#### Phase 3.3b: Ollama LLM (~8 hours, parallel)

```go
type OllamaLLMClient struct {
    client *http.Client
    url    string // "http://localhost:11434"
    model  string // "mistral", "llama2", "neural-chat", etc.
}

func NewOllamaLLMClient(url, model string) (*OllamaLLMClient, error) {
    // 1. Verify Ollama is running
    // 2. Check if model is available
    // 3. Return client
}

func (c *OllamaLLMClient) Reshape(ctx context.Context, packedContext string, maxTokens int) (string, error) {
    // POST /api/generate with context
    // Stream response and collect output
    // Return reshaped context
}
```

**Tests:**
- Verify Ollama running (mock server)
- Generate text via API
- Handle streaming response
- Respect max tokens

---

#### Phase 3.3c: OpenAI LLM (~6 hours, parallel)

```go
type OpenAILLMClient struct {
    client *openai.Client
    model  string // "gpt-4-turbo", "gpt-3.5-turbo"
    temp   float32
}

func NewOpenAILLMClient(apiKey, model string, temperature float32) (*OpenAILLMClient, error) {
    // 1. Create OpenAI client
    // 2. Verify API key
    // 3. Return client
}

func (c *OpenAILLMClient) Reshape(ctx context.Context, packedContext string, maxTokens int) (string, error) {
    // Call /v1/chat/completions
    // Prompt: "Reshape this context for clarity: {packedContext}"
    // Return response
}
```

**Dependencies:** OpenAI SDK

**Tests:**
- Verify API key (mock)
- Generate text via API
- Handle rate limiting
- Respect max tokens

---

#### Phase 3.3d: Anthropic LLM (~6 hours, parallel)

```go
type AnthropicLLMClient struct {
    client *anthropic.Client
    model  string // "claude-3-sonnet", "claude-3-opus"
    temp   float32
}

func NewAnthropicLLMClient(apiKey, model string, temperature float32) (*AnthropicLLMClient, error) {
    // 1. Create Anthropic client
    // 2. Verify API key
    // 3. Return client
}

func (c *AnthropicLLMClient) Reshape(ctx context.Context, packedContext string, maxTokens int) (string, error) {
    // Call /messages
    // Prompt: "Reshape this context for clarity: {packedContext}"
    // Return response
}
```

**Dependencies:** Anthropic SDK

**Tests:**
- Verify API key (mock)
- Generate text via API
- Handle streaming (if available)
- Respect max tokens

---

## Reshaping Prompt Template

All backends use this consistent prompt:

```
System: You are a context reshaping assistant. Your job is to improve the clarity and structure of technical documentation for AI consumption.

User: Reshape this technical context for optimal readability by:
1. Re-prioritizing information by importance
2. Removing redundancy
3. Structuring hierarchically
4. Highlighting key concepts
5. Improving clarity where needed

Keep the total output under {maxTokens} tokens. Preserve all critical information.

Original context:
{packedContext}

Reshaped context:
```

---

## Parallelization Strategy

**After Phase 3.1 (Config complete):**
- 3.3a (ONNX): Primary, 1.5-2 days
- 3.3b (Ollama): Parallel, ~8 hours
- 3.3c (OpenAI): Parallel, ~6 hours
- 3.3d (Anthropic): Parallel, ~6 hours

Can run fully in parallel (4 developers) or sequentially (1 developer, ~2 days).

**Note:** Phase 3.2 and 3.3 are independent — can run fully parallel tracks.

---

## Acceptance Criteria

- [~] ONNX LLM — scaffolding only; `SupportedModels` intentionally empty; no production ONNX LLM models available yet ⚠️
- [x] Ollama LLM — HTTP wiring done, **live tested against ollama serve** ✅
- [ ] OpenAI LLM (gpt-3.5-turbo and gpt-4-turbo) — deferred to v0.1.0
- [ ] Anthropic LLM (claude-3-sonnet, claude-3-opus) — deferred to v0.1.0
- [x] All implement LLMClient interface consistently
- [x] Config system correctly routes to chosen backend
- [x] Consistent reshaping prompt across all backends
- [x] Tests: 20+ unit/mock tests passing; integration tests with real Ollama passing ✅
- [x] Documentation: supported models + quality/speed trade-offs

**Evidence:** `integration_test.go` — Ollama LLM reshape test passing against real `ollama serve` (qwen3.5:4b, ~40s). Bug fixes: HTTP timeout increased to 5min for CPU inference; verification optimized to use `/api/tags` (fast) instead of full generation. ONNX LLM documented as scaffolding only (no models) — workaround via Ollama.

---

## Implementation Steps (Per Backend)

1. **Define backend struct** (30 min)
2. **Implement LLMClient interface** (1 hour)
3. **API communication** (1.5-2 hours per backend)
4. **Prompt engineering** (1 hour — test prompt quality)
5. **Error handling** (45 min)
6. **Tests** (1.5 hours)
7. **Documentation** (30 min)

---

## Effort Estimate

| Backend | Effort | Can Parallelize |
|---------|--------|-----------------|
| ONNX (primary) | 1.5-2 days | Yes |
| Ollama | 8 hours | Yes (parallel with ONNX) |
| OpenAI | 6 hours | Yes (parallel with ONNX) |
| Anthropic | 6 hours | Yes (parallel with ONNX) |
| **Total** | 2-3 days | **2 days parallel** |

---

## Dependencies

- ✅ Phase 3.1 (Config System) must be complete
- hugot (already integrated, for ONNX)
- OpenAI SDK (add dependency)
- Anthropic SDK (add dependency)

---

## Unblocks

- Phase 3.4: Context Reshaping (uses LLM backends from 3.3 + embeddings from 3.2)

---

## Success Metrics

✅ **All backends work:**
- ONNX: Local inference with TinyLlama or Mistral
- Ollama: Inference via local/remote Ollama server
- OpenAI: API calls with rate limiting
- Anthropic: API calls with proper auth

✅ **Consistent behavior:**
- All use same reshaping prompt
- All respect token budget
- All return reshaped context string

✅ **Ready for Phase 3.4:**
- Config selects any backend
- LLM reshaping ready to integrate with embeddings

---

## Quality Considerations

**Prompt Quality:**
- Test reshaping on representative contexts
- Measure: "Is output more readable than input?"
- Iterate on prompt if needed

**Speed Trade-offs:**
- ONNX: Fastest (local)
- Ollama: Fast (local server)
- OpenAI: Slower (API latency)
- Anthropic: Similar to OpenAI

**Cost:**
- ONNX/Ollama: Free (user's compute)
- OpenAI: ~$0.001-0.002 per context
- Anthropic: ~$0.003-0.015 per context

---

**Status:** Ready to implement (after Phase 3.1).

Can run in parallel with Phase 3.2: two independent teams working simultaneously.
