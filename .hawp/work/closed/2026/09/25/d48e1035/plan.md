# Guard provider materializer output paths

## Outcome

The PR-review finding was implemented and reconciled with the current branch.

## Verification

The focused regression coverage and repository-wide Go, HAWP, distribution,
formatting, and diff-hygiene checks passed before close.

## Close Checklist

- [x] Outcome recorded.
- [x] Verification evidence recorded or referenced.
- [x] Backlog row removed from active coordination.
- [x] Plan archived under the close date.

**UUID:** `d48e1035-aa9c-4345-bae1-258856c1294c`
**Type:** fix
**Reported:** 2026-09-23
**Source plan:** `a390cf52`

## Context

`librarian/src/internal/application/providersync/providersync.go` directly creates and writes generated provider files returned by `ComputeOutputs`. The shared filesystem guard package already provides repository-root ancestor rejection and atomic writes.

## Fix Work

Reject invalid or symlinked output ancestry before `MkdirAll`, then use the shared atomic writer for generated files. Preserve unchanged-file detection and validation behavior. Add tests proving symlinked provider directories are rejected without writing outside the repository.

## Verification

Run provider-sync tests, full Go tests, vet, HAWP checks, distribution validation, formatting, and diff hygiene.

## Outcome

Done. `Materialize` rejects symlinked ancestors and non-regular destinations
before mutation, then writes through the repository-aware atomic writer.
Provider-sync focused tests and vet pass.
