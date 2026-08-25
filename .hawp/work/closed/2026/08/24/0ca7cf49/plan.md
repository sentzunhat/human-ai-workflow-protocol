# 0ca7cf49 — Code cleanup: stale llm-reshape doc traces

**Type:** fix  
**Status:** done  
**Closed:** 2026-08-24

## Outcome

3 stale `--llm-reshape` doc-comment references removed from Go source. `toJSONReferences` confirmed still live (used for `--context --format json`); ContextReshaper and NullEmbedder kept intact. `go vet ./...`, `go build ./cmd/hawp`, and `go test ./...` all clean after changes.

## Verification

- [x] `go vet ./...` — clean
- [x] `go build ./cmd/hawp` — clean
- [x] `go test ./...` (short) — all packages pass
- [x] grep for `llm-reshape`, `llmReshape`, `ragOutput` across `librarian/src/` — no actionable dead code remaining

## Close Checklist

- [x] Implementation complete (3 stale comments removed)
- [x] Build and tests verified
- [x] BACKLOG updated; plan moved to closed

## What was done

- `run.go`: removed stale reshape reference from `toJSONReferences` doc comment (function still live for --context --format json)
- `readme_generator.go`: removed `--llm-reshape` example from embedded README template
- `null_embedder.go`: removed stale --llm-reshape clause from NullEmbedder doc comment
- No dead functions found — ContextReshaper, toJSONReferences, NullEmbedder all still in active use
- `go vet ./...`, `go build ./cmd/hawp`, `go test ./...` all clean
