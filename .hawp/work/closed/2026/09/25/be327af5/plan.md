# Reject symlinked .gitignore in MCP provider configuration

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

**Backlog ID (Legacy):** — (UUID-native item)
**UUID:** `be327af5-2f1c-4b63-b7e6-502dc4a5ae00`
**Type:** bug
**Reported:** 2026-09-20

---

## Input (verbatim)

> MCP configuration writes to .gitignore without checking for symlinks

## Intake Summary

Every file-writing MCP provider configuration path calls
`ensureGitignoreEntry`, which reads and writes `repoRoot/.gitignore` using
`os.ReadFile` and `os.WriteFile` without checking whether the file or any
of its ancestors is a symlink. A repository author who commits `.gitignore`
as a symlink can redirect the append to an arbitrary external file.

## Current Context

Other provider configuration paths (JSON/TOML) include regular-file checks
or ancestor guards, but the shared `.gitignore` helper is not covered. The
injected text is limited to one of the fixed HAWP ignore entries, which
constrains but does not eliminate the integrity impact on an external file.

## Initial Analysis

**Directly verified:**

- `librarian/src/internal/platform/mcp/configure/config.go:172-193`
  implements `ensureGitignoreEntry` using `os.ReadFile` and `os.WriteFile`
  on the target path.
- No `os.Lstat`, `symlink_guard.RejectSymlinkAncestors`, or regular-file
  check is present in this function.
- `librarian/src/internal/platform/mcp/configure/config.go:89-159`
  shows that all file-writing configure paths reach this helper.
- `librarian/src/internal/infrastructure/filesystem/symlink_guard.go:10-50`
  provides `RejectSymlinkAncestors` which is already used by other write
  paths in the codebase.

**Inferred (not yet proven):**

- Severity is low because the appended text is one of a small fixed set of
  HAWP ignore entries; arbitrary content cannot be injected.
- Exploitation requires a repository author to commit `.gitignore` as a
  symlink and an operator to run `hawp mcp configure`.

**Likely scope:**

- `librarian/src/internal/platform/mcp/configure/config.go`
- Focused test for a symlinked `.gitignore`.

## Root Cause

`ensureGitignoreEntry` was written as a convenience helper and was not
treated as a security-sensitive write boundary. It performs no symlink or
regular-file checks before reading or writing the target.

## Recommended Fix

- Before reading or writing `.gitignore`, call `RejectSymlinkAncestors` on
  the file path and its parent, and use `os.Lstat` to reject a non-regular
  existing `.gitignore`.
- If no `.gitignore` exists, create it with a non-following open to prevent
  a race.

## Verification Plan

- Add a test that creates `.gitignore` as a symlink to an external file,
  calls the configure path, and asserts it errors without modifying the
  external target.
- Run `go test ./...`, `go vet ./...`, `git diff --check`, and
  `go run ./cmd/hawp check` from `librarian/src`.

## Risk + Review Gate

**Risk:** low — targeted helper with fixed injection payload; no configuration
state is changed by the fix
**Gate:** standard review

## Backlog + Plan Link

**Status now:** done
**Plan file:** work/active/be327af5/plan.md

## Next Step

- [x] Investigation recorded above
- [x] Plan written with options, recommendation, and verification scope
- [x] Backlog moved to `plan-ready`
- [x] Obtain explicit approval before implementation

## Outcome

Done. `ensureGitignoreEntry` rejects symlinked ancestors before reading and
uses the managed atomic writer with `Lstat`-based regular-file validation.
The focused regression proves a symlinked `.gitignore` cannot modify its
external target.
