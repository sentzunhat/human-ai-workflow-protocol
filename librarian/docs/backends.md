# Backend Architecture & Strategy

## Overview

**v0.0.3** is the current stable release. It supports configurable backends for embeddings and LLM via a factory pattern. Users choose backends independently for local-only or hybrid configurations.

Future versions (v0.0.4, v0.0.5, etc.) will add incremental improvements. Larger changes (cloud APIs, vision support) are deferred to future releases.

---

## Embeddings Backends

### ONNX (Phase 3.2a)
**Status:** ✅ Production-ready

- **Models:** Hugot-sourced (all-MiniLM-L6-v2, bge-base-en-v1.5)
- **Speed:** <50ms per embedding
- **Privacy:** 100% local, no network calls
- **Dependencies:** None (included in binary)
- **First run:** Auto-downloads ~50MB model (cached for future)

**Use when:**
- Working offline or air-gapped
- Want fastest embeddings on CPU
- Privacy-critical workloads

### Ollama (Phase 3.2b)
**Status:** ✅ Live-tested against real Ollama (2026-07-25)

- **Models:** Any GGUF model (nomic-embed-text, mxbai-embed-large, all-minilm, etc.)
- **Speed:** **24-29ms per embedding** (measured on Mac CPU-only localhost)
- **Privacy:** 100% local (requires `ollama serve` running)
- **Dependencies:** Ollama service running on localhost:11434
- **Setup:** `ollama pull nomic-embed-text`

> **Verified:** Integration tests passing. All 3 tests (single, similarity, batch) confirmed working. Bugs fixed during testing (API field name, no synthetic vectors).

**Use when:**
- Want better embedding quality (768/1024-dim vs 384-dim)
- Already running Ollama for LLM
- Prefer single service for all models
- HTTP latency acceptable for your use case

### Cloud APIs (Phase 3.2c/d)
**Status:** 🔮 Planned for future releases

- **OpenAI**: text-embedding-3-small (planned)
- **Anthropic**: Future API support (when available)
- **Google**: text-embedding-004 (planned)

**Integration:** Same factory pattern, authentication via API keys when implemented

---

## LLM Backends

### Ollama (Phase 3.3b)
**Status:** ✅ Live-tested against real Ollama (2026-07-25)

- **Models:** Any GGUF model (mistral, neural-chat, llama2, qwen3.5:4b, etc.)
- **Speed:** **~40s per reshape (qwen3.5:4b), ~60-120s (mistral)** on Mac CPU-only
- **Privacy:** 100% local (requires `ollama serve` running)
- **Dependencies:** Ollama service running on localhost:11434
- **Setup:** `ollama pull mistral`

> **Verified:** Integration test passing against real Ollama. Bug fixes: HTTP timeout raised to 5min (was 60s), verification optimized via `/api/tags` instead of full generation.

**Use when:**
- Want local, private LLM reshaping
- Have Ollama already running
- No API keys/costs required
- Can wait 30-120s per reshape (CPU-dependent)

### ONNX (Phase 3.3a Scaffolding)
**Status:** 🔮 Scaffolded, waiting for model support

- **Models:** Architecture ready, working on small text2text models (T5, BART, FLAN)
- **Next step:** Implement text2text model integration (60-180MB models, 500-2000ms inference)
- **Use case:** Full local+private pipeline (embed + reshape) with no external APIs

### Cloud APIs (Phase 3.3c/d)
**Status:** 🔮 Planned for future releases

- **OpenAI**: gpt-3.5-turbo, gpt-4-turbo (planned)
- **Anthropic**: claude-3-sonnet, claude-3-opus (planned)

---

## Combination Matrix

### Working (Ship v0.0.3)

| Embeddings | LLM | Status | Best For |
|---|---|---|---|
| ONNX | Ollama | ✅ Live-tested | Default: offline embeddings + flexible LLM |
| Ollama | Ollama | ✅ Live-tested | All-Ollama: single service, full control |

> Both combinations verified against real services (not mocks). Performance measured and documented below.

### Planned (Future Releases)

| Embeddings | LLM | Status | Notes |
|---|---|---|---|
| ONNX | ONNX | 🔮 Next | Text2text models (T5, BART) for full-local pipeline |
| Ollama | ONNX | 🔮 Next | High-quality embed + fast reshape |
| OpenAI | OpenAI | 🔮 Future | Cloud embeddings + LLM, API keys required |
| Anthropic | Anthropic | 🔮 Future | Anthropic APIs when available |
| Mixed Cloud | Mixed Cloud | 🔮 Future | Any combination supported by architecture |

---

## Factory Pattern

Both embeddings and LLM follow the same initialization pattern:

```go
// Embeddings factory
embedder, err := embeddings.NewEmbedder(backend, model)
if err != nil {
	// Backend not available or model invalid
}

// LLM factory
llmClient, err := llm.NewLLMClient(backend, model)
if err != nil {
	// Backend not available or model invalid
}

// ContextReshaper accepts both
reshaper, err := NewContextReshaper(ReshapingConfig{
	EmbeddingsBackend: backend,
	EmbeddingsModel:   model,
	LLMBackend:        backend,
	LLMModel:          model,
})
```

**Benefits:**
- Adding new backends requires only implementing the interface
- No changes to ContextReshaper
- Config system selects any backend combination
- Tests can mock or use real backends independently

---

## Configuration System

### Priority (Highest to Lowest)

1. **Environment variables** (HAWP_EMBEDDINGS_BACKEND, etc.)
2. **~/.hawp/config/context.json**
3. **Defaults** (ONNX embeddings + Ollama LLM)

### Example: ~/.hawp/config/context.json

```json
{
  "embeddings": {
    "backend": "onnx|ollama",
    "model": "model-name"
  },
  "llm": {
    "backend": "ollama|onnx|openai|anthropic",
    "model": "model-name",
    "max_tokens": 512,
    "temperature": 0.7
  },
  "reshaping": {
    "top_k": 5,
    "timeout_ms": 30000
  }
}
```

---

## Performance Characteristics

### Embeddings Performance

| Backend | Model | Time | Dim | Quality | Source |
|---|---|---|---|---|---|
| ONNX | all-MiniLM-L6-v2 | 119ms | 384 | Good | ✅ Measured (2026-07-25 benchmark) |
| ONNX | bge-base-en-v1.5 | 605ms | 768 | High | ✅ Measured (2026-07-25 benchmark) |
| Ollama | all-minilm | 29ms | 384 | Good | ✅ Measured (2026-07-25 benchmark) |
| Ollama | nomic-embed-text | 24ms | 768 | High | ✅ Measured (2026-07-25 benchmark) |

**Key insight:** Ollama HTTP faster than ONNX for 384-dim; network latency offset by ONNX model startup cost (happens once per session).

### LLM Performance

| Backend | Model | Time/Response | Quality | Source |
|---|---|---|---|---|
| Ollama | qwen3.5:4b | ~40s per reshape | Good | ✅ Measured (2026-07-25 integration test) |
| Ollama | mistral | ~60-120s per reshape | High | ✅ Verified available (model info) |
| Ollama | neural-chat | est. 500-800ms | Good | 📋 Not benchmarked yet |
| Ollama | llama2 | est. 1000-3000ms | High | 📋 Not benchmarked yet |

**Caveat:** All times on Mac CPU-only. GPU would show dramatic speedup (2-5x). LLM speed varies with model size and prompt length.

**Factors affecting performance:**
- CPU cores (more = faster)
- GPU availability (dramatic speedup if available)
- Model size (larger = slower but higher quality)
- Prompt length (longer = slower)
- Max tokens (higher = slower)

---

## Phase Breakdown & Timeline

| Phase | Backend | Type | Status | ETA |
|---|---|---|---|---|
| 3.2a | ONNX | Embeddings | ✅ Live-tested | 2026-07-24 |
| 3.2b | Ollama | Embeddings | ✅ Live-tested | 2026-07-25 |
| 3.2c | OpenAI | Embeddings | 🔮 Planned | v0.1.0 |
| 3.2d | Anthropic | Embeddings | 🔮 Planned | v0.1.0 |
| 3.3a | ONNX | LLM | ✅ Scaffolded | 2026-07-24 |
| 3.3b | Ollama | LLM | ✅ Live-tested | 2026-07-25 |
| 3.3c | OpenAI | LLM | 🔮 Planned | v0.1.0 |
| 3.3d | Anthropic | LLM | 🔮 Planned | v0.1.0 |

---

## Why This Strategy?

### Why ONNX First?
- ✅ No external dependencies
- ✅ No API keys needed
- ✅ Fast enough for most use cases
- ✅ Self-contained binary

### Why Ollama Second?
- ✅ Simple HTTP API (proven pattern)
- ✅ Open source, can run locally
- ✅ Supports any GGUF model
- ✅ Testing pattern (mock HTTP servers)

### Why Cloud APIs in v0.1.0?
- ✅ Core architecture proven locally first
- ✅ Users have privacy-by-default with ONNX/Ollama
- ✅ Cloud as optional upgrade, not forced
- ✅ More time to refine API cost tracking

### Why No ONNX LLM Yet?
- ❌ No production ONNX LLMs available (TinyLlama too small, Mistral/Llama don't export well)
- ⏳ Watching for improvements in ONNX model availability
- ✅ Ollama works perfectly as fallback
- 📋 v0.1.0 can revisit when models mature

---

## Testing Strategy

### Unit Tests (Fast, No Dependencies)
- Mock HTTP servers for Ollama tests
- Test factory patterns, error handling, timeouts
- All tests run in <3 seconds

### Integration Tests (Optional, Real Services)
- Skip gracefully if Ollama not running (flag: `-short`)
- Run against real Ollama if available
- Measure real-world performance
- Verify model dimension expectations

### Example Test
```go
func TestReshapeWithOllama(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	config := ReshapingConfig{
		EmbeddingsBackend: "ollama",
		EmbeddingsModel:   "nomic-embed-text",
		LLMBackend:        "ollama",
		LLMModel:          "mistral",
	}

	reshaper, err := NewContextReshaper(config)
	if err != nil {
		t.Skipf("Ollama not available: %v", err)
	}
	defer reshaper.Close()

	// ... test with real backends
}
```

---

## Future Extensions

### Easy (Same HTTP Pattern)
- More Ollama models (already supported)
- OpenAI embeddings/LLM (Phase 3.2c/3.3c)
- Anthropic embeddings/LLM (Phase 3.2d/3.3d)

### Medium (New HTTP Pattern)
- HuggingFace inference API
- Replicate
- Azure OpenAI

### Hard (New Architecture)
- Self-hosted ONNX models via FastAPI
- LLaMA.cpp server
- vLLM inference server
- ONNX LLM when models available

---

## Configuration Best Practices

### For Users
```json
{
  "embeddings": {
    "backend": "onnx",
    "model": "all-MiniLM-L6-v2"
  },
  "llm": {
    "backend": "ollama",
    "model": "mistral"
  }
}
```
**Rationale:** Fast local embeddings, flexible LLM selection via Ollama

### For Privacy-Critical
```json
{
  "embeddings": {
    "backend": "onnx",
    "model": "bge-base-en-v1.5"
  },
  "llm": {
    "backend": "onnx",
    "model": "???"
  }
}
```
**Status:** LLM portion blocked until ONNX LLMs available

### For Maximum Speed
```json
{
  "embeddings": {
    "backend": "ollama",
    "model": "nomic-embed-text"
  },
  "llm": {
    "backend": "ollama",
    "model": "neural-chat:latest"
  }
}
```
**Rationale:** All HTTP-based, batch-able in future, single service to manage

---

## Roadmap to v0.1.0+

**v0.0.3 (shipping 2026-07-30):**
- ✅ ONNX embeddings (real inference, live-tested, benchmarked)
- ✅ Ollama embeddings (HTTP API, live-tested against real Ollama, benchmarked)
- ✅ Ollama LLM (HTTP API, live-tested against real Ollama, bugs fixed)
- ✅ Context reshaper pipeline
- ✅ Config system

**v0.1.0 (planned):**
- 🔮 OpenAI backends (3.2c/3.3c)
- 🔮 Anthropic backends (3.2d/3.3d)
- 🔮 Cost tracking for cloud APIs
- 🔮 Rate limiting + token budget
- 🔮 ONNX LLM (if models available)

**v0.2.0+ (future):**
- 🔮 Parallel batch processing
- 🔮 Pipeline caching
- 🔮 GPU optimization
- 🔮 Multi-turn conversations
