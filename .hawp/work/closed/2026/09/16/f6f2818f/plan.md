# Document TOCTOU in writeMCPJSON between Lstat and WriteFile

**Backlog ID (Legacy):** — (UUID-native item)
**UUID:** `f6f2818f`
**Type:** fix
**Reported:** 2026-09-16

---

## Input (verbatim)

> Document TOCTOU in writeMCPJSON between Lstat and WriteFile

## Intake Summary

`writeMCPJSON` performs an `os.Lstat` check followed by `os.WriteFile`. There is a TOCTOU (time-of-check/time-of-use) window between the two calls where the file could be replaced with a symlink. This was flagged in the Copilot review.

## Current Context

The function is a local CLI tool writing to a project config file. A symlink swap in this window would require an attacker with local filesystem write access. The risk was acknowledged but not mitigated in code; the decision was to document the limitation.

## Initial Analysis

**Directly verified:**

- `writeMCPJSON` called `os.Lstat` then `os.WriteFile` without atomic replacement.

**Inferred (not yet proven):**

- For a local tool with a trusted operator, documenting the limitation is an acceptable resolution.

**Likely scope:**

- `librarian/src/internal/application/mcp/` — `writeMCPJSON` function.

## Risk + Review Gate

**Risk:** low (local tool, trusted operator, no network exposure)
**Gate:** auto-implement on low

## Backlog + Plan Link

**Status now:** done
**Plan file:** work/closed/2026/09/16/f6f2818f/plan.md

## Verification

- No behavior change; documentation only.
- `go test ./...` passes; `go vet ./...` clean.
- Comment placed at the TOCTOU window in `writeMCPJSON` explaining the local threat model.

## Outcome

Added a comment in `writeMCPJSON` at the point between `os.Lstat` and `os.WriteFile` documenting the TOCTOU window and noting that it is accepted for a local CLI tool where the operator controls the filesystem. No code behavior changed.

## Close Checklist

- [x] TOCTOU window documented with a comment explaining the local-tool threat model.
- [x] `go test ./...` passes; `go vet ./...` clean.
- [x] Plan moved to `closed/2026/09/16/f6f2818f/`.
- [x] BACKLOG.md updated.
