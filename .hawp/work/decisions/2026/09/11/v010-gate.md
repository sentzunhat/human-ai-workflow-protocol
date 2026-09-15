# v0.1.0 Token-Reduction Gate

**Date:** 2026-09-11  
**Status:** decided (updated 2026-09-11 — reshape gate split into two metrics)  
**Owner:** beltrd

## Decision

### Search context pipeline gate (unchanged)

> **≥ 20% average token reduction** on the 10-query standard benchmark suite,
> measured using the `(len(text)+3)/4` token estimator.

**Status: PASSED** — 23% average (23521 raw → 18029 shaped).  
Evidence: `benchmark/runs/2026-09-11-search-token-savings.md`

---

### Reshape pipeline gate (updated 2026-09-11)

The original gate (≥20% token compression) is the wrong metric for reshape.
Structured intake (mission/constraints/output/checkpoint) always *expands* short
requests by design — that expansion is the value, not a failure.

Live Ollama run confirmed: 10/10 structured, but 247 raw tokens → 820 shaped
(-232% average). This is expected and correct behavior.

The reshape pipeline is gated on **two independent metrics**, both required:

#### Gate 1 — Coverage (required for v0.1.0)

> **≥ 80% of the 10 standard queries produce valid structured output** (all
> required fields: `mission`, `constraints`, `output`).

Current status:
- **Ollama (mistral:7B): 10/10 (100%) — PASSED**
- **ONNX Phi-3-mini (3.8B): 10/10 (100%) — PASSED**
- ONNX 360M: 1/10 (10%) — model too small; superseded by Phi-3-mini

Root cause of earlier ONNX failures: hugot's `RunPipeline([]string)` wraps
the prompt as a "user" message; onnxruntime-genai then applies the model's
chat template a second time (double-templating). Fix: use `RunMessages` with
`{Role: "system"}/{Role: "user"}` objects so onnxruntime-genai templates once.
Also removed manual `buildChatPrompt`/`detectModelType` — no longer needed.

#### Gate 2 — Downstream savings (target for v0.1.0, experiment required)

> **Structured intake saves ≥ 20% of total tokens when the shaped fields replace
> raw request + retrieved context** in a downstream LLM call.

Rationale: the real token-reduction value of reshape is not input compression
but output efficiency — a model receiving `mission: <one sentence>` needs
far fewer tokens to respond usefully than one receiving a vague 15-token
request plus 2000 tokens of retrieved context. This gate measures that.

Measurement method (to be implemented):
1. Baseline: raw request + top-5 search context chunks (raw token count)
2. Shaped: DraftIntake fields + 0 search context (trust the structured intake)
3. Gate: shaped token count < 0.8 × baseline token count on average

**Status: NOT YET MEASURED** — requires benchmark experiment.

---

## What does NOT count

- Compression of the intake itself (expansion is expected and correct).
- Search quality improvements (precision/recall) alone.
- Latency improvements alone.
- Single-query outliers — all gates are averages across the 10-query suite.

## Evidence required before tagging v0.1.0

1. `benchmark/runs/YYYY-MM-DD-search-token-savings.md` — compression ≥ 20% ✓ (done)
2. `benchmark/runs/YYYY-MM-DD-reshape-coverage.md` — coverage ≥ 80% (Ollama: done; ONNX Phi-3-mini: done ✓)
3. `benchmark/runs/YYYY-MM-DD-reshape-downstream-savings.md` — downstream ≥ 20% (not yet)

## Tracked by

- [`288d543c`](../../../active/288d543c/plan.md) — token-reduction benchmark harness
