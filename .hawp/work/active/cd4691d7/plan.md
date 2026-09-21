# Reject symlink ancestors in kitsync repository writes

**Backlog ID (Legacy):** — (UUID-native item)
**UUID:** `cd4691d7-d43e-48ff-a009-8eb8f984a846`
**Type:** bug
**Reported:** 2026-09-20

---

## Input (verbatim)

> Reject symlink ancestors in kitsync repository writes

## Intake Summary

The kitsync path checks added for manifest traversal prove only lexical
containment. Kit and provider writes still call `MkdirAll`, `CreateTemp`, and
`Rename` through pre-existing repository-controlled symlink ancestors. A
symlinked `.hawp/kit` or provider destination directory can therefore redirect
release content outside `repoRoot` when the user runs update/sync.

## Current Context

The current dirty branch contains the just-completed lexical containment fix in
`librarian/src/internal/domain/kitsync/apply.go`; this item is a distinct
physical-containment follow-up and overlaps those uncommitted files. It must be
implemented only after that lane is committed or explicitly reconciled.

## Initial Analysis

**Directly verified:**

- `librarian/src/internal/domain/kitsync/apply.go:16-59` rejects absolute and
  `..` manifest paths but does not resolve or reject symlink ancestors.
- `SyncKit` at `librarian/src/internal/domain/kitsync/apply.go:62-68` derives a
  fixed lexical destination without a physical-containment guard.
- `copyTree` and `copyFile` at
  `librarian/src/internal/domain/kitsync/apply.go:154-185` and `:292-315`
  create destination directories and temporary files before renaming them.
- `librarian/src/internal/infrastructure/repositories/kitsync/filecopy.go:18-43`
  delegates directly to symlink-following OS operations.
- `librarian/src/internal/infrastructure/filesystem/symlink_guard.go:10-50`
  already provides an ancestor-rejection policy used by other write paths, but
  kitsync does not call it.

**Inferred (not yet proven):**

- Final destination symlinks are normally replaced by rename rather than
  followed; the confirmed escape requires a symlink ancestor.
- Archive and manifest lexical traversal checks remain effective and are not
  invalidated by this finding.
- A user must invoke update/sync in the malicious or compromised repository.

**Likely scope:**

- `librarian/src/internal/application/kitsync/kitsync.go`
- `librarian/src/internal/domain/kitsync/apply.go`
- `librarian/src/internal/domain/kitsync/filecopier.go`
- `librarian/src/internal/infrastructure/repositories/kitsync/filecopy.go`
- Focused domain/application tests for kit, install, update, refresh, and
  seed-if-missing destinations.

## Root Cause

Destination authorization is completed before filesystem resolution. Once a
lexically valid path reaches the copier, no operation knows the trusted root,
so the OS follows an attacker-planted parent symlink and creates the temporary
file in the external directory. Atomic rename prevents partial writes but does
not restore containment.

## Options Considered

1. Call `RejectSymlinkAncestors` once for each resolved destination before
   copying. Small and consistent with existing code, but retains a documented
   check/use race.
2. Extend the `FileCopier` boundary with root-aware guarded destination
   operations, applying ancestor and target-type checks immediately before
   mutation. This keeps policy at the concrete filesystem seam and is the
   recommended portable fix.
3. Implement descriptor-relative no-follow operations per platform. Strongest
   race resistance, but substantially larger and should be justified separately
   if the portable guard is insufficient for the local single-user threat model.

## Recommended Fix

- Carry the authorized `repoRoot` into destination mutation operations without
  moving concrete filesystem behavior back into the domain package.
- Reject symlink ancestors for `.hawp/kit` and every provider destination;
  reject non-regular existing final targets before replacement.
- Re-check as close as practical to `MkdirAll`, `CreateTemp`, and `Rename`, and
  document any remaining TOCTOU limitation.
- Preserve the `FileCopier` port and the new manifest path-validation API.

## Verification Plan

- Add malicious-repository tests with `.hawp`, `.hawp/kit`, provider roots,
  and nested destination ancestors symlinked to an external temporary directory.
- Cover `SyncKit`, provider install/update, refresh, directory copy, and
  seed-if-missing behavior; assert no outside file is created or changed.
- Preserve ordinary copy behavior and the new absolute/`..` manifest tests.
- Run focused kitsync tests, `go test ./...`, `go vet ./...`, source-layout
  checks, `git diff --check`, `go run ./cmd/hawp check`, and
  `go run ./cmd/hawp work validate` from `librarian/src`.

## Risk + Review Gate

**Risk:** high — security-sensitive filesystem boundary with overlapping edits
**Gate:** explicit review plus overlap reconciliation before implementation

**Overlap check:** blocked for implementation until current uncommitted changes
to `librarian/src/internal/domain/kitsync/apply.go` and
`librarian/src/internal/domain/kitsync/kitsync_test.go` are committed or
explicitly incorporated. Planning is safe now.

## Backlog + Plan Link

**Status now:** plan-ready
**Plan file:** work/active/cd4691d7/plan.md

## Next Step

- [x] Investigation recorded above
- [x] Plan written with options, recommendation, and verification scope
- [x] Backlog moved to `plan-ready`
- [ ] Obtain explicit approval and reconcile overlapping changes before implementation
