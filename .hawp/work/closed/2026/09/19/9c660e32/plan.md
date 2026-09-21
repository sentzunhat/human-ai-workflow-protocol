# Remove filesystem operations from domain work

**Backlog ID (Legacy):** — (UUID-native item)
**UUID:** `9c660e32-a2ac-4bdd-831c-2c7c18995676`
**Type:** infrastructure
**Reported:** 2026-09-19

---

## Input (verbatim)

> Move concrete os and filesystem traversal out of internal/domain/work and its subpackages into application or infrastructure adapters. Preserve parent compatibility where practical, inject readers writers listers and path policies, verify the work package has no production os imports, then audit domain/kit and domain/kitsync for the next boundary.

## Intake Summary

The folder/package split is substantially improved, but production domain code
still imports `os` or performs concrete filesystem traversal. This item tracks
the stricter ports-and-adapters cleanup separately from package organization.

## Current Context

`domain/work` now has reusable `backlog`, `identity`, `intake`, `markdown`,
`model`, `normalization`, `table`, and `validation` packages. The parent retains
compatibility aliases and adapters. Application callers already construct
`WorkSource` from infrastructure functions, but several domain implementations
still bypass that seam directly.

## Initial Analysis

**Directly verified:**

- Production `os` imports remain in 15 `domain/work` files, concentrated in
  validation, normalization, and parent migration/closed-record adapters.
- `domain/work/validation` directly reads directories/files and writes warnings
  to stderr.
- `domain/work/normalization` directly reads plan/backlog files and performs
  active-row cleanup; migration path manipulation has an injected I/O seam but
  the concrete adapter still lives in the domain parent.
- The source package graph has no imports from `domain/work` to infrastructure,
  but concrete standard-library filesystem access is still present.
- Validation policy now consumes an injected `Source`; the application layer
  supplies the concrete `os.ReadDir`, `os.ReadFile`, and `os.Stat` functions.
- `domain/work/validation` has zero production `os` imports after this slice.
- Adjacent production domain folders with the same pattern are `context`,
  `kit`, `kitsync`, `providersync`, `provision`, and `distribution`; `kit` and
  `kitsync` are the next likely candidates after this item.

**Inferred (not yet proven):**

- The current package layout may create a false sense of adapter separation
  while domain code still owns filesystem policy and effects.
- Removing direct filesystem access will require source contracts or snapshot
  inputs and may require compatibility changes in parent APIs.

**Likely scope:**

- `domain/work/validation`: injected read/stat/list/warning source.
- `domain/work/normalization`: injected scan and mutation source; move preview
  copy and concrete OS adapter to application/infrastructure.
- Parent compatibility wrappers: retain only where they do not reintroduce
  filesystem effects; migrate internal callers to capability packages.
- Follow-up audit: `domain/kit` then `domain/kitsync`.

## Risk + Review Gate

**Risk:** medium — changes several domain contracts and filesystem seams.
**Gate:** review each bounded extraction with focused tests, then run the full
repository checks before removing compatibility paths.

## Backlog + Plan Link

**Status now:** done
**Plan file:** work/closed/2026/09/19/9c660e32/plan.md

## Next Step

- [x] Investigation recorded above (required before planning)
- [x] Write or update the plan file
- [x] Create a separate tracked item for the strict boundary cleanup
- [x] Move backlog status to analyzing
- [x] Extract validation filesystem source
- [x] Update application work validation to inject infrastructure-backed source
- [x] Verify `domain/work/validation` has no production `os` imports
- [x] Extract normalization scan/cleanup filesystem source
- [x] Extract normalization read-only scan source
- [x] Extract normalization active-row cleanup mutation source
- [x] Extract closed-record normalization filesystem mutation source
- [x] Remove direct `os` types from the normalization migration seam
- [x] Move migration preview and concrete OS adapter out of domain
- [x] Verify zero production `os` imports under `internal/domain/work`
- [x] Audit `domain/kit` and `domain/kitsync`

## Completion Evidence (2026-09-19)

## Verification

- `go test ./...` passes from `librarian/src`.
- `go vet ./...` passes from `librarian/src`.
- `git diff --check` passes.
- Production `os` import scan is empty for `domain/work`, `domain/kit`, and `domain/kitsync`.

## Outcome

Removed the remaining concrete filesystem choices from the audited domain packages. Application and infrastructure layers now provide directory readers, file readers, normalization sources, provider detection, and copier contracts.

## Close Checklist

- [x] Zero production `os` imports under `internal/domain/work`.
- [x] `domain/kit` validation uses an injected directory reader.
- [x] `domain/kitsync` uses injected detection and `io/fs`-based copier contracts.
- [x] Full repository verification passed.
- [x] Plan moved to `closed/2026/09/19/9c660e32/`.
- [x] `BACKLOG.md` updated.

- Removed the final concrete filesystem callbacks from `domain/work` compatibility adapters; application callers now pass injected normalization and validation sources.
- Added injected directory reading to `domain/kit` validation and moved provider detection in `domain/kitsync` behind the existing copier seam.
- Replaced `os`-specific domain copier types with `io`/`io/fs` contracts; concrete implementations remain in infrastructure.
- Confirmed zero production `os` imports under `internal/domain/work`, `internal/domain/kit`, and `internal/domain/kitsync`.
- Passed focused tests, `go test ./...`, `go vet ./...`, `git diff --check`, and `go run ./cmd/hawp check --no-update-check`.
