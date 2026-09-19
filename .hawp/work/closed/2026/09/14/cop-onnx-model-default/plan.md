# Improvement: ONNX LLM default model should be Phi-3-mini, not SmolLM2-360M

**UUID:** `cop-onnx-model-default`
**Type:** improvement
**Severity:** Low
**Reported:** 2026-09-13
**Status:** `inbox`
**Source:** Copilot PR #40 review — `librarian/src/internal/infrastructure/models/onnx/llm.go`

---

## Finding

The v0.0.24 CHANGELOG gate says SmolLM2-360M scored 1/10 and is not a
recommended default, yet `DefaultLLMModel = "SmolLM2-360M-Instruct"` is still
the only model in `SupportedLLMModels`.

Phi-3-mini was identified as the validated model but is not present in the
registry. This means the advertised ONNX gate cannot be reproduced through
the shipped code.

## Fix Plan

1. Add `Phi-3-mini` (or the specific HuggingFace repo used in benchmarks) to
   `SupportedLLMModels`.
2. Change `DefaultLLMModel` to the Phi-3-mini entry.
3. Keep `SmolLM2-360M-Instruct` in the registry (with a note) for backwards
   compatibility with any existing downloaded models.
4. Update the CHANGELOG to reference the correct default and to document the
   `ExternalDataFile` sidecar download requirement (already noted in code).

## Evidence

Benchmark run: `benchmark/runs/2026-09-11-reshape-token-savings-onnx-phi3-mini.md`

## Files

- `librarian/src/internal/infrastructure/models/onnx/llm.go`
- `librarian/src/CHANGELOG.md`

## Verification

- `DefaultLLMModel` is `"Phi-3-mini-4k-instruct"` in `llm.go`.
- `SupportedLLMModels` includes the Phi-3-mini entry with correct HFRepo and model files.
- `go test ./...` passes; `go vet ./...` clean.

## Outcome

Added Phi-3-mini-4k-instruct to `SupportedLLMModels` and set it as `DefaultLLMModel`. SmolLM2-360M-Instruct remains in the registry for users who explicitly request it. Release notes and benchmarks are now consistent with the shipped default. `go test ./...` and `go vet ./...` pass.

## Close Checklist

- [x] `DefaultLLMModel = "Phi-3-mini-4k-instruct"` in `llm.go`.
- [x] Phi-3-mini entry present in `SupportedLLMModels`.
- [x] `go test ./...` passes; `go vet ./...` clean.
- [x] Plan status: done.
- [x] Plan moved to `.hawp/work/closed/2026/09/14/cop-onnx-model-default/plan.md`.
- [x] BACKLOG.md updated.
