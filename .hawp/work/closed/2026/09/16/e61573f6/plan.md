# Remove double expandConfigProviders call in Configure()

**Backlog ID (Legacy):** — (UUID-native item)
**UUID:** `e61573f6`
**Type:** fix
**Reported:** 2026-09-16

---

## Input (verbatim)

> Remove double expandConfigProviders call in Configure()

## Intake Summary

`Configure()` called `expandConfigProviders` twice: once to get the list, then internally again when it called `WriteProviderConfigs`. This was identified in the Copilot review.

## Current Context

`Configure()` computed an expanded provider list, then passed it to `WriteProviderConfigs`, which called `expandConfigProviders` again on its own. The duplication was harmless but wasteful and could cause issues if expansion had side effects.

## Initial Analysis

**Directly verified:**

- `Configure()` called `expandConfigProviders` and then `WriteProviderConfigs`, which called it again internally.

**Inferred (not yet proven):**

- Splitting to an unexported `writeProviderConfigs` that accepts an already-expanded list eliminates the double call.

**Likely scope:**

- `librarian/src/internal/application/provider/configure.go`.

## Risk + Review Gate

**Risk:** low
**Gate:** auto-implement on low

## Backlog + Plan Link

**Status now:** done
**Plan file:** work/closed/2026/09/16/e61573f6/plan.md

## Verification

- `go test ./...` passes; `go vet ./...` clean.
- `Configure()` now calls unexported `writeProviderConfigs` with the already-expanded list; `expandConfigProviders` is invoked exactly once per `Configure()` call.

## Outcome

Split `WriteProviderConfigs` into a public wrapper (that calls `expandConfigProviders`) and an unexported `writeProviderConfigs` (that accepts an already-expanded list). `Configure()` now calls the unexported variant with the list it already holds, eliminating the double expansion.

## Close Checklist

- [x] Double expansion eliminated; `expandConfigProviders` called exactly once per `Configure()`.
- [x] `go test ./...` passes; `go vet ./...` clean.
- [x] Plan moved to `closed/2026/09/16/e61573f6/`.
- [x] BACKLOG.md updated.
