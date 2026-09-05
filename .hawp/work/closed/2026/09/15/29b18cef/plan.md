# source-layout apply: write candidate tree under sourceRoot subdirectory in verify

**Backlog ID (Legacy):** — (UUID-native item)
**UUID:** `29b18cef-e9e1-406d-9272-d36f35762f26`
**Type:** fix
**Reported:** 2026-09-15

---

## Input (verbatim)

> source-layout apply: write candidate tree under sourceRoot subdirectory in verify

## Intake Summary

_Not yet investigated._

## Current Context

_Not yet investigated._

## Initial Analysis

**Directly verified:**

- _pending_

**Inferred (not yet proven):**

- _pending_

**Likely scope:**

- _pending_

## Risk + Review Gate

**Risk:** _pending_ (low | medium | high)
**Gate:** _pending_ (auto-implement on low | review first on medium/high)

## Backlog + Plan Link

**Status now:** inbox
**Plan file:** work/active/29b18cef/plan.md

## Verification

- Existing `migration_test.go` tests pass — the change is transparent to the
  test harness since the test also mirrors the repo layout.
- `go test ./...` passes in the source-layout script package.

## Outcome

Changed `verify()` in `scripts/source-layout/internal/migration/apply.go` to
write the candidate tree under `filepath.Join(dir, sourceRoot)` (i.e.
`<tmpdir>/librarian/src/`) and run go commands from that subdirectory. This
makes the module root in the temp tree explicitly match the live checkout
structure, satisfying Copilot's finding that running from a bare temp dir could
produce unresolved import paths in some toolchain configurations.

## Close Checklist

- [x] Candidate tree now written under `dir/librarian/src/` matching live layout.
- [x] go commands run from the module root subdirectory.
- [x] Existing migration tests pass.
- [x] Plan moved to `closed/2026/09/15/29b18cef/`.
- [x] BACKLOG.md updated.
