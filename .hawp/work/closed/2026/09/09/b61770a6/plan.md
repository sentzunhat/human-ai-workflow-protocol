# Move model and usage adapters out of domain packages

**Backlog ID (Legacy):** — (UUID-native item)
**UUID:** `b61770a6-2623-4808-9a03-ba3571fbd214`
**Type:** improvement
**Reported:** 2026-09-07

---

## Input (verbatim)

> Decide package boundaries across layers; plan folder-within-folder organization that stays simple, repeatable, and avoids duplication. Planning only now; pace usage and keep implementation slices small.

## Investigation

Confirmed: domain/embeddings and domain/llm contain net/http and ONNX/session
implementations plus concrete factories. domain/usage/store.go imports database/sql
and os. Search default wiring also constructs concrete adapters in application.

This is a distinct implementation subitem of
[47c793d6](../../../../../active/47c793d6/plan.md), following the
[shared package map](../../../../../active/47c793d6/package-boundaries.md).

## Decision And Plan

1. Move Ollama embedding/LLM implementations and HTTP tests into planned
   infrastructure/models/ollama. Keep the existing consumer ports/contracts stable.
2. Move ONNX and explicit none implementations into models/onnx and models/none;
   preserve build tags, model behavior, session closure, errors, and defaults.
3. Introduce internal/bootstrap for explicit service wiring shared by entrypoints;
   move concrete factory/default construction there. Update callers atomically;
   do not retain domain factories importing infrastructure, or let bootstrap
   import CLI/MCP packages that call it.
4. Move usage SQL implementation to infrastructure/repositories/usage; keep report
   data and policy in domain/usage and introduce application/usage only for real
   use-case orchestration. No generic provider registry, model change, or new
   dependency. Naming rule: infrastructure folder names reflect the domain they
   serve, not the mechanism — see [package-boundaries.md](../../../../../active/47c793d6/package-boundaries.md).

Alternative: retain flat packages with renamed files. Rejected for these scoped
owners because it leaves navigation/ownership concerns unresolved. Do not nest
already-cohesive small packages or duplicate rules between layers.

## File Ownership And Coordination

Own domain/embeddings, domain/llm, domain/usage and planned models adapters,
`librarian/src/internal/infrastructure/repositories/usage/`, and internal/bootstrap.
Application/search/context construction and CLI/MCP composition callsites require
serialized edits after the other items' current slices. Tests move with adapters.
All shorthand directories above are under `librarian/src/internal` unless stated.
Before implementation, enumerate exact source/test/destination paths for the
selected slice. No broad simultaneous moves; preserve the current dirty checkpoint.

## Risk And Execution Gate

Risk: medium (Go package/import topology). Status: done. Owner: Codex.
Implementation completed in bounded slices on 2026-09-09. No release, merge, push,
or binary refresh was included. The remaining work-domain migration is tracked by
`742aa60b`.

## Verification / Acceptance

Domain ports import no HTTP/SQL/model runtime. Mock HTTP and injected-service
error/closure tests pass. Default build passes without model downloads. Preserve
ORT build-tag behavior; run compilation where configured, label unavailable
native-library runtime checks as unproven. Full Go suite/vet/HAWP pass.
Initial planning did not execute models. Implementation checks are recorded below;
HAWP validation verifies record integrity, not package migration correctness.

## Usage Budget

One cohesive slice per turn. Read only owned files, direct dependencies and tests.
Short updates; update this plan rather than creating a report for each move.
Focused tests during edits, full checks once at slice completion. No recurring
usage monitor, reset redemption, model switch, or paid action requested.

## Progress

- [x] Step 1 — Ollama + None adapters moved to `infrastructure/models/ollama/` and
  `infrastructure/models/none/`; factory created at `infrastructure/models/factory.go`
  (commit `8cfda10`, 2026-09-08)
- [x] Step 2 — ONNX adapters moved to `infrastructure/models/onnx/`; domain/embeddings
  and domain/llm now own only interfaces and constants; factory fully infra-internal
  (commit `283d640`, 2026-09-09)
- [x] Step 3 slice 1 — Added `internal/bootstrap/` and moved model-factory wiring for
  the CLI embedding service behind an injected application factory (working tree,
  2026-09-09). Focused tests and vet pass. Bootstrap remains outside `domain/` because
  it imports infrastructure by design.
- [x] Step 3 — Routed index, search, and context model construction through
  `internal/bootstrap` with narrow injected factories; application/domain production
  packages no longer import `infrastructure/models` (working tree, 2026-09-09).
- [x] Step 4 — Moved the SQLite usage adapter and config file I/O to
  `infrastructure/repositories/usage/`; kept usage ports, values, and report policy
  in `domain/usage/` and updated CLI/MCP callers (working tree, 2026-09-09).

## Outcome

Model implementations and factories now live under infrastructure, concrete
composition lives under `internal/bootstrap`, usage SQLite/config I/O lives under
`infrastructure/repositories/usage`, and application search/context/index services
receive narrow injected factories. Domain production packages no longer import
model, SQL, or filesystem adapters for this item.

## Verification

- [x] `go test ./...` passed with live integration/benchmark tests excluded by
  explicit build tags.
- [x] `go vet ./...` passed.
- [x] `git diff --check` passed.
- [x] Production boundary audit found no `infrastructure/models` imports under
  application or domain packages.

## Close Checklist

- [x] Plan moved to `closed/2026/09/09/b61770a6/`.
- [x] Existing status checkpoint retained and updated.
- [x] Unrelated `core/.hawp/bin/hawp` deletion preserved.
- [x] Remaining work-domain migration identified without applying a broad rewrite.
