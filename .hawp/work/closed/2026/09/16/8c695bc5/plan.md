# Improve WriteProviderConfigs partial-failure error message

**Backlog ID (Legacy):** — (UUID-native item)
**UUID:** `8c695bc5`
**Type:** fix
**Reported:** 2026-09-16

---

## Input (verbatim)

> Improve WriteProviderConfigs partial-failure error message

## Intake Summary

`WriteProviderConfigs` returned a generic error on partial failure without indicating which providers were successfully written before the failure. This made diagnosis difficult.

## Current Context

When `WriteProviderConfigs` failed partway through writing multiple provider configs, the error message named only the failed provider, not the ones already written. On a partial failure, users would not know which configs were applied.

## Initial Analysis

**Directly verified:**

- `WriteProviderConfigs` returned an error that named the failing provider but not the previously-written ones.

**Inferred (not yet proven):**

- Collecting successfully written providers and including them in the error message improves diagnosability.

**Likely scope:**

- `librarian/src/internal/application/provider/configure.go` — `WriteProviderConfigs` / `writeProviderConfigs`.

## Risk + Review Gate

**Risk:** low
**Gate:** auto-implement on low

## Backlog + Plan Link

**Status now:** done
**Plan file:** work/closed/2026/09/16/8c695bc5/plan.md

## Verification

- `go test ./...` passes; `go vet ./...` clean.
- Partial-failure error message now names the providers successfully written before the failure.

## Outcome

Updated `writeProviderConfigs` to accumulate the names of successfully written providers. On failure, the error message now includes both the failing provider and the list of those already written, enabling users to know the partial state.

## Close Checklist

- [x] Error message names successfully written providers on partial failure.
- [x] `go test ./...` passes; `go vet ./...` clean.
- [x] Plan moved to `closed/2026/09/16/8c695bc5/`.
- [x] BACKLOG.md updated.
