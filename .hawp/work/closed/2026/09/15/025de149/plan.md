# intake: add symlink guard for BACKLOG.md and active plan dir writes

**Backlog ID (Legacy):** — (UUID-native item)
**UUID:** `025de149-4a22-4398-8a0f-852a2acb3929`
**Type:** fix
**Reported:** 2026-09-15

---

## Input (verbatim)

> intake: add symlink guard for BACKLOG.md and active plan dir writes

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
**Plan file:** work/active/025de149/plan.md

## Verification

- `go test ./...` passes; `go vet ./...` clean.
- `RejectSymlinkAncestors(workDir, backlogPath)` called before reading or writing
  BACKLOG.md; non-regular backlog file rejected with a clear error.
- `RejectSymlinkAncestors(workDir, planDir)` called before `MkdirAll(planDir)`.

## Outcome

Added `filesystem.RejectSymlinkAncestors` guards in `intake/intake.go`:

1. Before reading BACKLOG.md — blocks a symlinked `.hawp/work/BACKLOG.md` from
   redirecting writes outside the repository.
2. Lstat check rejects a non-regular BACKLOG.md (symlink, directory, device).
3. Before `MkdirAll(planDir)` — blocks a symlinked `active/` subtree.

Mirrors the same protection already applied in `doc.go`, `config.go`, and
`config_codex.go`.

## Close Checklist

- [x] Symlink guard for backlog path added and verified.
- [x] Symlink guard for active plan dir added and verified.
- [x] Non-regular backlog file rejected before read/write.
- [x] `go test ./...` passes; `go vet ./...` clean.
- [x] Plan moved to `closed/2026/09/15/025de149/`.
- [x] BACKLOG.md updated.
