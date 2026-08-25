# 0ca7cf49 — Code cleanup: stale llm-reshape doc traces

**Type:** fix  
**Status:** done  
**Closed:** 2026-08-24

## What was done

- `run.go`: removed stale reshape reference from `toJSONReferences` doc comment (function still live for --context --format json)
- `readme_generator.go`: removed `--llm-reshape` example from embedded README template
- `null_embedder.go`: removed stale --llm-reshape clause from NullEmbedder doc comment
- No dead functions found — ContextReshaper, toJSONReferences, NullEmbedder all still in active use
- `go vet ./...`, `go build ./cmd/hawp`, `go test ./...` all clean
