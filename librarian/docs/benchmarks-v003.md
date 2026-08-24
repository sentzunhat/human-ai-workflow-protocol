# v0.0.3 Benchmarks — ONNX & Ollama

**Date:** 2026-07-26
**Machine:** Apple M1 Max (CPU-only, no GPU acceleration)
**Method:** `go test -bench` against real running backends (ONNX Runtime in-process, Ollama at localhost:11434) — no mocks.

---

## Embeddings

| Backend | Model | Dim | Single embed (ns/op → ms) | Batch (5 texts) |
|---|---|---|---|---|
| ONNX | all-MiniLM-L6-v2 | 384 | 103,969,180 ns ≈ **104ms** | ~100ms |
| ONNX | bge-base-en-v1.5 | 768 | 638,918,347 ns ≈ **639ms** | ~950ms |
| Ollama | all-minilm | 384 | 34,197,347 ns ≈ **34ms** | ~161ms |
| Ollama | nomic-embed-text | 768 | 34,982,805 ns ≈ **35ms** | ~164ms |
| Ollama | mxbai-embed-large | 1024 | 31,048,736 ns ≈ **31ms** | ~175ms |

**Bigger model requested:** `mxbai-embed-large` (1024-dim) via Ollama — the largest embedding
dimension exercised in this benchmark. It is *not slower* than the smaller Ollama models; Ollama's
embedding latency is dominated by HTTP/model-load overhead, not dimension, at this scale.

**ONNX bigger model attempted and reverted:** `bge-large-en-v1.5` (1024-dim ONNX) was added to
`SupportedModels` and benchmarked, but the download **deadlocked** inside the `go-huggingface`
dependency (`hugot@v0.7.5` → `github.com/gomlx/go-huggingface@v0.3.5`'s `internal/downloader.Manager`)
— a genuine Go runtime "all goroutines are asleep" deadlock, not a slow-network timeout. Reproduced
2026-07-26. Reverted from `SupportedModels` rather than ship a model that hangs on first use; see the
code comment in `internal/domain/embeddings/onnx_embedder.go` for the exact failure signature.

### Takeaway

Ollama embeddings are **3–20x faster** than ONNX at every dimension tested, confirming the finding
from the 2026-07-25 session. ONNX's advantage is zero network dependency (fully offline), which is
why it remains the v0.0.3 default for `embeddings.backend`.

---

## LLM (reshaping)

| Backend | Model | Params | Single reshape (maxTokens=100) |
|---|---|---|---|
| Ollama | mistral | 7B (Q4_K_M) | 12,447,584,917 ns ≈ **12.4s** |

Only Mistral was benchmarked directly (qwen3.5:4b, used in the 2026-07-25 session, is no longer
pulled on this machine — `BenchmarkAll_Ollama_LLM_Qwen` skips gracefully when the model is absent,
by design). 12.4s for a 100-token reshape on CPU is consistent with the 2026-07-25 baseline
(2.5–10s for shorter prompts, up to 120s for larger ones).

**ONNX LLM:** working, verified 2026-07-27 (manual `-tags ORT` + native-libraries build,
macOS arm64) — SmolLM2-360M-Instruct (genai-converted variant), single-context reshape
in **~1.1s**. Not usable on this repo's default `go test`/`go build` (no `-tags ORT`):
`NewONNXLLMClient` creates the ORT session before touching the network, so a default
build fails fast with hugot's own "to enable ORT" error rather than downloading a model
that can't run. Getting here required 8 distinct fixes (wrong HF repo format, missing
external-data sidecar file, wrong prompt format, rpath vs DYLD_LIBRARY_PATH, directory
vs file path for the library option, and more) — full blow-by-blow in
`v0.1.0_VISION.md`'s "ONNX Text2Text Model (ORT backend)" section, including the exact
local build command. Known gap: multi-context `ReshapeBatch()` calls can still fail with
a "max length reached" error not yet root-caused; single-context `Reshape()` (what the
RAG pipeline actually uses) is reliable.

---

## End-to-end CLI proof (not a `go test` benchmark — real invocation)

```
$ time hawp search "semantic search" --context --llm-reshape --max-tokens 800
...
95.85s user 1.50s system 264% cpu 36.741 total
```

36.7s wall time for the full pipeline: ONNX embeddings (bge-base-en-v1.5, concept extraction over
several sentences) + Ollama Mistral reshape (800-token budget). This is the number a user actually
experiences running `--llm-reshape` end to end, as opposed to the isolated per-call benchmarks above.

---

## Test suite timing (context, not a benchmark)

Running the full `internal/...` suite against live Ollama (no mocks skipped) takes **~300s**,
dominated by the `internal/domain` package's integration tests and benchmarks, which call real
Mistral/embedding models. `go test -short ./internal/...` does **not** currently skip these — see
the "Minor / deferred" note in `SHIP_READY_VERIFICATION.md`.

---

## Raw logs

- `/tmp/bench_embeddings2.log` — ONNX MiniLM/BGE + Ollama MiniLM/Nomic
- `/tmp/bench_mxbai.log` — Ollama mxbai-embed-large (1024-dim)
- `/tmp/bench_mistral.log` — Ollama Mistral LLM reshape
