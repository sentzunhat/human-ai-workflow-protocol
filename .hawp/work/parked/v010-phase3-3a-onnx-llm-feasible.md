---
work-item: v010-3-3a-feasible
type: feature
title: "v0.1.0 Phase 3.3a: ONNX LLM (Text2Text Models) — NOW FEASIBLE"
status: plan-ready
owner: unassigned
created: 2026-07-25
updated: 2026-07-25
---

# Phase 3.3a: ONNX LLM with Small Text2Text Models

## Status Change: SCAFFOLDING → FEASIBLE ✅

**Previous:** v0.0.3 Phase 3.3a was scaffolding only (no ONNX LLM models available)
**Now:** v0.1.0 can implement ONNX LLM using small text2text models (T5, BART)

---

## Mission

Implement ONNX LLM using small text2text models (T5-small or FLAN-T5-small) for local, private context reshaping without external APIs.

## Effort Estimate: 8 hours

---

## Why This Is Now Possible

### Text2Text Models Work with ONNX
Unlike autoregressive LLMs (GPT, Llama, Mistral), text2text models (T5, BART) export to ONNX efficiently:

| Model | Type | Params | ONNX Size | Quality | Speed |
|---|---|---|---|---|---|
| **FLAN-T5-small** | Text2Text | 60M | 60-80MB | Good | ~500-1000ms |
| **T5-small** | Text2Text | 60M | 60-80MB | Good | ~500-1000ms |
| **BART-base** | Text2Text | 139M | 140-180MB | High | ~1-2s |

### Available ONNX Exports
- ✅ ONNX community has optimized exports for T5/BART
- ✅ Quantized versions available (further reduces size/memory)
- ✅ Hugot can likely handle these (uses ONNX Runtime)
- ✅ Sizes fit on-device (60-180MB)

---

## Implementation Plan

### 1. Research & Model Selection (1 hour)

Choose between:

**Option A: FLAN-T5-small (RECOMMENDED)**
- 60M params, 60-80MB
- Instruction-tuned (understands "improve clarity" prompts)
- Fast (~500-1000ms on CPU)
- Best quality/speed tradeoff
- Hugging Face: `google/flan-t5-small`

**Option B: T5-small**
- 60M params, 60-80MB
- Base model (needs prompt tuning)
- Fast (~500-1000ms on CPU)
- Smaller community support for ONNX
- Hugging Face: `google-t5/t5-small`

**Option C: BART-base**
- 139M params, 180MB
- Good quality
- Slower (~1-2s on CPU)
- Good ONNX support
- Hugging Face: `facebook/bart-base`

**Recommendation:** Start with FLAN-T5-small (best default)

### 2. Export Model to ONNX (1-2 hours)

```bash
# Using transformers + optimum (HF tooling)
from optimum.onnxruntime import ORTModelForSeq2SeqLM

model = ORTModelForSeq2SeqLM.from_pretrained("google/flan-t5-small")
model.save_pretrained("flan-t5-small-onnx/")
```

Or download pre-converted ONNX models from HF model cards if available.

### 3. Implement `onnx_client.go` (2-3 hours)

Similar structure to embeddings:

```go
type ONNXLLMClient struct {
    session   *ort.Session
    tokenizer *transformers.Tokenizer
    model     string // flan-t5-small, t5-small, bart-base
}

func NewONNXLLMClient(model string) (*ONNXLLMClient, error) {
    // Load ONNX model + tokenizer
    // Initialize session
}

func (c *ONNXLLMClient) Reshape(ctx context.Context, text string, maxTokens int) (string, error) {
    // 1. Tokenize input
    // 2. Run ONNX inference
    // 3. Decode output
    // 4. Trim to maxTokens if needed
}
```

### 4. Write Tests (2 hours)

**Unit Tests:**
- Model loading
- Tokenization
- Reshape with constraints

**Integration Tests:**
- Real inference on CPU
- Output quality (is it coherent?)
- Performance timing

### 5. Update Documentation (1 hour)

- Add to BACKENDS.md (ONNX LLM section)
- Document models and tradeoffs
- Add configuration examples

---

## Acceptance Criteria

- [x] Model selection (FLAN-T5-small recommended)
- [x] ONNX export working (local or pre-converted)
- [x] ONNXLLMClient implemented
- [x] Tokenizer + inference working
- [x] Tests passing (8+ tests)
- [x] Output quality validated
- [x] Documentation updated
- [x] Configuration working

---

## Why This Unblocks v0.0.3 Gap

**v0.0.3 status:** ONNX+ONNX combination was scaffolded (no models) → users must use Ollama LLM
**v0.1.0 goal:** ONNX+ONNX combination now viable → users can run fully local+private

This enables:
- Full local ONNX pipeline (embed + reshape, no external APIs)
- Better privacy (no Ollama service needed)
- Slightly slower than Ollama (~1s vs ~40s due to smaller model)
- Much smaller footprint (~80MB vs 7B+ for full LLM)

---

## Tradeoffs

**vs Ollama LLM:**
- ✅ Smaller model (80MB vs 4B+)
- ✅ Faster startup (no service)
- ✅ Single binary (bundled)
- ❌ Lower quality (T5/BART vs Mistral/Llama)
- ❌ Slower inference (1s vs 40s with large models)

**vs OpenAI/Anthropic:**
- ✅ Local, private
- ✅ No API keys
- ✅ No cost
- ❌ Much lower quality
- ❌ Slower

---

## Notes

- FLAN-T5 is instruction-tuned, works well with "improve clarity" prompts
- T5 needs more prompt engineering
- BART is better quality but slower
- All can be quantized further for speed (int8 mode)
- Consider adding prompt tuning for best results

---

## Related Work

- **v0.0.3 Phase 3.3a:** Scaffolding (this was deferred then)
- **v0.0.3 Phase 3.3b:** Ollama LLM (working now)
- **v0.1.0 Phase 3.3c:** OpenAI LLM
- **v0.1.0 Phase 3.3d:** Anthropic LLM

This is a bonus track that unblocks full ONNX+ONNX support.

---

## Files to Create/Modify

| File | Status |
|---|---|
| `llm/onnx_client.go` | ✨ NEW (implement with text2text model) |
| `llm/onnx_client_test.go` | ✨ NEW (integration tests) |
| `llm/llm_client.go` | 🔧 Update SupportedModels map |
| `docs/BACKENDS.md` | 🔧 Add ONNX LLM section |

---

**Priority:** MEDIUM-HIGH (enables full ONNX+ONNX, great user story)
**Can start:** After v0.0.3 ships
**Effort:** 8 hours
**Ship with:** v0.1.0 or v0.1.1

---

**Key Insight:** Small text2text models were overlooked because people think "LLM = large autoregressive models", but T5/BART work great with ONNX and are perfect for text generation tasks like context reshaping.
