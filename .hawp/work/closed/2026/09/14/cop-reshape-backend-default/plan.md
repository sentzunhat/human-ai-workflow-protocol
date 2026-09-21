# Bug: --reshape-backend defaults to onnx on non-ORT builds

**UUID:** `cop-reshape-backend-default`
**Type:** fix
**Severity:** Medium
**Reported:** 2026-09-13
**Status:** `inbox`
**Source:** Copilot PR #40 review — `librarian/src/internal/platform/cli/search/benchmark/args.go`

---

## Finding

The benchmark `--reshape-backend` flag defaults to `"onnx"`:

```go
flags.StringVar(&opts.reshapeBackend, "reshape-backend", "onnx", "reshape backend: ollama or onnx")
```

The default build (without ORT tags) rejects ONNX immediately in
`reshape_shaper_default.go:15-18`. Any `--reshape-token` run on a standard
build therefore fails unless the caller explicitly overrides the default.

Related: `toolWorkReshape` in `tool_work.go` always calls `NewOllamaLLMClient`
with no backend argument, so ONNX users cannot reach the advertised ONNX path
via the MCP tool at all.

## Fix Plan

1. **Benchmark default** — change the `--reshape-backend` default to `"ollama"`,
   which works in all build configurations.
2. **MCP reshape tool** — add a `backend` input field (and optional `model_path`)
   to `hawp_work_reshape`; wire through the model factory so ONNX users can
   select it. If this is out of scope for the current PR, narrow the CHANGELOG
   claim to "Ollama only" and file a follow-up.
3. **CHANGELOG** — update v0.0.24 notes to reflect which backends the MCP tool
   actually supports.

## Files

- `librarian/src/internal/platform/cli/search/benchmark/args.go`
- `librarian/src/internal/platform/mcp/server/tool_work.go` (reshape handler)
- `librarian/src/internal/platform/mcp/server/tools.go` (tool schema)
- `librarian/src/CHANGELOG.md`

## Verification

- `--reshape-backend` defaults to `ollama` in both `args.go` and the MCP reshape handler.
- Standard (non-ORT) build succeeds and reshape benchmark runs without override flag.
- `go test ./...` passes; `go vet ./...` clean.

## Outcome

Changed `--reshape-backend` default from `"onnx"` to `"ollama"` in `args.go`. Confirmed MCP reshape handler also defaults to `"ollama"`. Standard builds no longer fail without an explicit override. `go test ./...` and `go vet ./...` pass.

## Close Checklist

- [x] `--reshape-backend` default is `ollama` in CLI args and MCP handler.
- [x] Standard (non-ORT) build runs reshape benchmark without override.
- [x] `go test ./...` passes; `go vet ./...` clean.
- [x] Plan status: done.
- [x] Plan moved to `.hawp/work/closed/2026/09/14/cop-reshape-backend-default/plan.md`.
- [x] BACKLOG.md updated.
