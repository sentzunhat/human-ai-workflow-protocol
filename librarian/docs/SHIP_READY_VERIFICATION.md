# v0.0.3 Ship Readiness — Audit

**Date:** 2026-07-26 (revised, second pass — fixes applied and re-verified)
**Status:** ✅ **Blockers resolved.** All findings from the 2026-07-26 audit are fixed and re-verified
against real running backends (Ollama at localhost:11434, ONNX Runtime in-process — no mocks for the
CLI proof).

> **History.** The 2026-07-25 revision of this file claimed "ALL CHECKS PASSED" based on checking file
> existence, not reachability. That was wrong and was retracted in the first 2026-07-26 revision, which
> found the headline feature (`--llm-reshape`) was unreachable from the CLI, plus 7 failing tests and
> two display/config bugs. This revision documents the fixes for every finding in that audit, each
> verified by execution, not by re-reading source.

---

## Method

Every claim below was produced by running a command or reading the named file. Labels: **Confirmed**
(directly executed/observed) or **Likely** (evidence-backed but interpretive).

---

## Fix 1 — Phase 3 reshaping wired into the CLI (was Finding 1, Blocker)

**Confirmed.** `internal/platform/cli/run.go`:
- `runSearch` now parses `--llm-reshape` (implies `--context`).
- On the flag, `tryReshape()` loads `appcontext.LoadContextConfig`, constructs
  `appcontext.NewContextReshaper`, and calls `.Reshape()`.
- Any failure (config error, backend unreachable, reshape error) prints an explicit `Warning:` line to
  stderr and falls back to the unreshaped context block — never silent.

**Executed proof — real reshape:**
```
$ time hawp search "semantic search" --context --llm-reshape --max-tokens 800
# Search Results: "semantic search" (LLM-reshaped, pipeline: onnx-ollama)
**Key concepts:** Results, ONNX, Input, CLI, Go-first
...
36.741 total
```
Real ONNX embeddings (bge-base-en-v1.5) + real Ollama Mistral reshape. Pipeline label correct.

**Executed proof — graceful fallback (backend unreachable):**
```
$ HAWP_OLLAMA_URL="http://localhost:9999" hawp search "semantic search" --context --llm-reshape --max-tokens 500
Warning: --llm-reshape unavailable (initialize LLM: verify Ollama model: connect to Ollama at
http://localhost:9999: ...: connection refused); using unreshaped context. Configure via
~/.hawp/config/context.json or HAWP_LLM_BACKEND/HAWP_EMBEDDINGS_BACKEND.
# Search Results: "semantic search"
**Results:** 2 | **Tokens:** 496/500
...
```
Exit 0, but the failure is visible on stderr and the fallback output is the honest unreshaped block —
this is the correct behavior for a documented-but-unavailable feature, contrasted with the prior
silent-ignore behavior.

---

## Fix 1b — `ReshapingConfig` ignored the configured Ollama URL (new finding, found while verifying Fix 1)

**Confirmed, Severity: High.** While proving Fix 1's fallback path, the first attempt (pointing
`HAWP_OLLAMA_URL` at an unreachable port) unexpectedly succeeded — because `ReshapingConfig` had no
URL field, so `NewContextReshaper` always called `NewEmbedder`/`NewLLMClient` (no-URL variants), which
default to `localhost:11434` regardless of config. A user pointing HAWP at a remote Ollama host would
have silently talked to `localhost` instead.

**Fix:** added `EmbeddingsURL`/`LLMURL` to `ReshapingConfig`; `reshaper.go` now calls
`NewEmbedderWithURL`/`NewLLMClientWithURL`; `run.go`'s `tryReshape` passes `cfg.Backends.Ollama.URL`
through. Re-ran the fallback proof above after this fix — it now correctly fails against the bad URL.

---

## Fix 1c — `DefaultConfig()` defaulted to a nonexistent ONNX LLM model (new finding, found while verifying Fix 1)

**Confirmed, Severity: Blocker if left as-is.** `DefaultConfig()` set `LLM.Backend: "onnx"`,
`LLM.Model: "TinyLlama-1.1B"`. `llm.SupportedModels` is an empty map (ONNX LLM has no implemented
models — documented as scaffolding). Every `--llm-reshape` call under the untouched default config
would have failed with `unsupported ONNX model`. `TestDefaultConfig`/`TestMergeConfig` were asserting
this broken default as correct, masking the bug.

**Fix:** default LLM backend changed to `ollama`/`mistral` (matches `BACKENDS.md`'s documented default
combination: "ONNX + Ollama (default, recommended)"). Tests updated to match, with a comment
explaining why ONNX cannot be the LLM default yet.

---

## Fix 2 — Documentation now matches reality (was Finding 2, Blocker)

**Confirmed.** `CONTEXT_RESHAPING.md`, `CHANGELOG.md`, and `readme_generator.go`'s `--llm-reshape`
claims are now true — verified by the executed proof in Fix 1. Additionally found and fixed:
`CONTEXT_RESHAPING.md` claimed the default ONNX embeddings model was `all-MiniLM-L6-v2` (384-dim);
the actual default is `bge-base-en-v1.5` (768-dim). Corrected, and the model comparison table was
extended with the exact tradeoffs (see `BENCHMARKS_v003.md`).

---

## Fix 3 — 8 of 10 Ollama LLM tests were failing (was Finding 3, reported as "7"; actual count was 8)

**Confirmed.** Root cause verified: `verifyModelAvailable` calls `GET /api/tags`, but every mock server
in `ollama_client_test.go` only handled `POST /api/generate` and 404'd everything else. Rewrote all 10
mock servers to branch on `GET /api/tags` (serving a matching model-list response) before falling
through to their existing `/api/generate` behavior.

**Executed proof:**
```
$ go test ./internal/domain/llm/... -v
--- PASS: TestNewOllamaLLMClient
--- PASS: TestOllamaLLMClientInterface
--- PASS: TestOllamaUnreachable
--- PASS: TestOllamaEmptyContext
--- PASS: TestOllamaReshapeSingleContext
--- PASS: TestOllamaReshapeBatch
--- PASS: TestOllamaEmptyBatch
--- PASS: TestOllamaServerError
--- PASS: TestOllamaContext
--- PASS: TestOllamaDefaultModel
ok  	.../internal/domain/llm	0.402s
```

Also added `TestOllamaEmbedRequestUsesPromptField` to `internal/domain/embeddings` — the regression
test recommended in `CODE_REVIEW_v003.md` Finding 2, which locks in the `Prompt`-not-`Input` field
contract so a silent revert of that earlier fix would be caught by unit tests, not just integration.

---

## Fix 4 — Token budget display (was Finding 4)

**Confirmed.** `format.go` hardcoded the budget operand to `0`. Added a `Budget int` field to
`ContextBlock`, set it in `FormatAsMarkdown`, and used it in `String()`. Reproduced in the executed
fallback proof above: `**Results:** 2 | **Tokens:** 496/500` (previously would have read `496/0`).

---

## Fix 5 — `aes256` silent downgrade (was Finding 5)

**Confirmed.** Both `EncryptKey`/`DecryptKey` (`encryption.go`) and `ContextConfig.Validate()`
(`config.go`) now reject `aes256` with an explicit error instead of silently falling back to base64.
No existing test depended on the fallback behavior (verified by grep before changing).

---

## Fix 6 — Reshape's empty-block guard never actually fired (new finding, found while re-running tests with real Ollama)

**Confirmed, Severity: Medium.** `Reshape()` guarded against empty input via
`block == nil || block.String() == ""`. `block.String()` always renders a `# Title` header even with
zero `Results`, so this condition can never be true for a non-nil block — `TestReshapeEmptyBlock`'s
"empty content" case was failing once Ollama was actually reachable (it had been masked for an
unknown period by a `t.Skipf` that fired whenever Ollama was unavailable). Fixed the guard to check
`len(block.Results) == 0` — the actual signal for "nothing to reshape."

---

## Benchmarks (requested: ONNX + Ollama, including a bigger model)

Full results in `BENCHMARKS_v003.md`. Summary:

| Backend | Model | Dim | Latency |
|---|---|---|---|
| ONNX | all-MiniLM-L6-v2 | 384 | ~104ms |
| ONNX | bge-base-en-v1.5 | 768 | ~639ms |
| Ollama | all-minilm | 384 | ~34ms |
| Ollama | nomic-embed-text | 768 | ~35ms |
| Ollama | mxbai-embed-large (bigger) | 1024 | ~31ms |
| Ollama LLM | mistral (100 tokens) | 7B | ~12.4s |

**Bigger ONNX model attempted and reverted:** `bge-large-en-v1.5` was added to `SupportedModels` and
benchmarked, but its download deadlocked inside the `go-huggingface` dependency (Go runtime-detected
deadlock, not a timeout). Reverted rather than ship a model that hangs on first use — see
`onnx_embedder.go`'s comment and `BENCHMARKS_v003.md` for the full trace. The "bigger model" ask was
satisfied via Ollama's `mxbai-embed-large` (1024-dim) instead, which has no such issue.

---

## Verified test state

```
$ go test ./internal/...          # full suite, real Ollama, no mocks skipped
ok  	.../internal/application/context    (was FAIL: 3 tests — now fixed, see Fix 1c/6)
ok  	.../internal/domain                 298s (real Mistral/embedding calls)
ok  	.../internal/domain/embeddings       (+1 new regression test)
ok  	.../internal/domain/llm              (was FAIL: 8 tests — now fixed, see Fix 3)
ok  	.../internal/platform/cli
... all other packages ok
```

`go vet ./...` — clean. `gofmt` applied to every file touched in this pass.

---

## What's still open (not blockers, tracked separately)

- `go test -short` does not actually skip the live-backend integration tests in
  `internal/domain/{integration,benchmark}_test.go` — they don't check `testing.Short()`. Not a
  correctness issue (they pass), but means "short mode" isn't actually short (~300s). Worth a follow-up.
- `bge-large-en-v1.5` ONNX download deadlock is an upstream `go-huggingface`/`hugot` issue, not
  something to fix in this repo. Tracked as a known gap, not attempted further here.
- `run.go` is ~1150 lines mixing flag parsing, search, ranking, and reshaping. Splitting is justified
  under the kit's "split large files for a reason" rule but is not a ship blocker.

---

## Verdict

All blockers from the 2026-07-26 first-pass audit are fixed and re-verified by execution (not just
by re-reading the diff). Two additional real bugs were found and fixed while verifying the primary
fix (URL passthrough, broken default config) — both would have caused user-visible failures if left
in place. Full test suite green. CLI proof executed against real backends, both success and
graceful-failure paths.

**No new version has been tagged or published as part of this pass**, per instruction — this is a
verification-and-fix pass only.
