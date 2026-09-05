# Status Checkpoint — 2026-09-09 — v0.0.24 Architecture Work

**Branch:** `feature/v0.0.24-work-folder-normalization`
**Parent item:** [47c793d6](../../active/47c793d6/plan.md)

---

## What Is Complete (this session)

### e7a5e294 — CLI decomposition (closed 2026-09-08)

All 11 CLI command families migrated from flat root files to nested
`cli/<family>/` packages with operation sub-packages. `run.go` is the sole
routing file. Full test suite, vet, and HAWP checks passed at each slice.
Closed to `closed/2026/09/08/e7a5e294/plan.md`.

Commits (relevant recent ones):
- `1d53f96` — link/boundary cleanup
- `fadb50d` — distribution + providers families; root cleanup
- `0978bad` — update family + init
- `462fa83` — close e7a5e294, CHANGELOG, BACKLOG housekeeping

### b61770a6 slice 1 — Ollama + None adapters (commit `8cfda10`, 2026-09-08)

- `infrastructure/models/ollama/` — `OllamaEmbedder` and `OllamaLLMClient` (HTTP)
- `infrastructure/models/none/` — `NullEmbedder` and `NullLLMClient`
- `infrastructure/models/factory.go` — routes "ollama"/"none"/"onnx" to adapters
- Domain factories removed; application callers updated to `inframodels.*`

### b61770a6 slice 2 — ONNX adapters (commit `283d640`, 2026-09-09)

- `infrastructure/models/onnx/` — `ONNXEmbedder` and `ONNXLLMClient` (hugot/ORT)
- `domain/embeddings/` — now owns only `Embedder` interface + `EmbeddingResult` + `DefaultModel` const
- `domain/llm/` — now owns only `LLMClient` interface + `ReshapingPrompt` const
- Factory fully infra-internal: no domain→infra calls remain in domain packages
- Integration and benchmark tests updated to import infrastructure/models/onnx

**Verification at each slice:** `go test ./...`, `go vet ./...`, HAWP check
(all 3 validations: kit, work, links) — all pass.

---

## Current Infrastructure Layout

```
infrastructure/
  models/
    factory.go        ← routes backends; only place selecting adapters
    ollama/           ← OllamaEmbedder, OllamaLLMClient, tests
    onnx/             ← ONNXEmbedder, ONNXLLMClient, tests
    none/             ← NullEmbedder, NullLLMClient
  sqlite/             ← index storage (unchanged)
  repositories/       ← (planned; empty now)
  ...
```

---

## What Is Next (ordered)

### 1. b61770a6 step 3 — `internal/bootstrap/` (next immediate slice)

Create a minimal bootstrap package that:
- Imports `infrastructure/models` factory and all infrastructure adapters
- Imports application ports/services
- Constructs concrete default services (embedder, LLM client, repos)
- Is imported ONLY by the main package / CLI entrypoint — never by domain, application, or infrastructure

Scope: investigate which entrypoints currently do their own `inframodels.NewEmbedder` calls and what the wiring looks like in `application/context/reshaper.go`, `application/search/service.go`, etc.

### 2. b61770a6 step 4 — `infrastructure/repositories/usage/`

Move `domain/usage/store.go` (imports `database/sql`, `modernc.org/sqlite`, `os`) to `infrastructure/repositories/usage/`. Introduce a port interface in `domain/usage/` if one is missing. Update consumers. The SQL logic belongs in infrastructure per the package map.

### 3. 742aa60b — Work domain separation

After the model/usage adapters are fully moved, begin separating `domain/work/` (currently mixes pure rules with `os` reads/writes) into:
- `domain/work/` — shared IDs/values only
- `domain/work/backlog/`, `intake/`, `validation/`, `normalize/` — pure rule sub-packages
- `application/work/intake/`, `validation/`, `normalize/` — use-case orchestration
- `infrastructure/workfs/` — concrete filesystem reads/writes

### 4. Release gate — PR to development → main

Once `b61770a6` and enough of `742aa60b` are complete, open:
- `feature/v0.0.24-work-folder-normalization` → `development` PR
- `development` → `main` PR

After merge to main: `git add core/.hawp/bin/hawp && git commit` (binary must not be committed on the feature branch — the deletion is still unstaged and must stay that way until the PR merges).

---

## Standards and Guidelines Check

| Concern | Status |
|---|---|
| Domain imports no infra | ✓ verified — `go vet ./...` catches cycles; all clean |
| No child→parent imports | ✓ factory is the only place that imports all sub-packages |
| Tests move with adapters | ✓ all four test files moved and adapted |
| One slice per turn | ✓ two slices across two separate commits |
| Full checks at slice completion | ✓ `go test ./...` + `go vet ./...` + HAWP check at each commit |
| Binary not committed on feature branch | ✓ `core/.hawp/bin/hawp` deletion remains unstaged |
| HAWP work records up to date | ✓ e7a5e294 closed, b61770a6 progress recorded, BACKLOG updated |
| Package map followed | ✓ names match `infrastructure/models/ollama`, `onnx`, `none` from package-boundaries.md |

---

## Open Constraints

- No PR or merge in this session.
- Binary commit only after PR merges to main.
- `a3df8a9c` (reshape) is parked behind topology work — do not resume until b61770a6 and at least slice 1 of 742aa60b are complete.
