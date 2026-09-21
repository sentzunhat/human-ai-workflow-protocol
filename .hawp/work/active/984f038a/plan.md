# Reject Markdown symlinks in links clean --apply

**Backlog ID (Legacy):** — (UUID-native item)
**UUID:** `984f038a-3887-49f0-8bb4-773af577a3e8`
**Type:** bug
**Reported:** 2026-09-20

---

## Input (verbatim)

> Markdown link cleanup can rewrite an external file through a repository symlink

## Intake Summary

`hawp links clean --apply` collects `.md` entries using `os.ReadDir` and
writes any file where broken local links are detected. The collection step
does not filter out symlinks, so a repository-controlled `.md` symlink
pointing to an external readable document can be collected and, if a link
inside it is classified as broken, rewritten by `os.WriteFile`.

## Current Context

The apply flag is required, reducing the default-invocation risk. Only files
where a broken local Markdown link is detected are rewritten, and only with
the repaired version of that file's content. The attacker must therefore
control both the in-repository symlink and an external file that contains a
Markdown link that the checker classifies as broken.

## Initial Analysis

**Directly verified:**

- `librarian/src/internal/application/links/check.go:131-150` iterates
  `os.ReadDir` results, selecting entries where `!entry.IsDir()` and the
  name ends `.md`; no `entry.Type().IsRegular()` or symlink check is present.
- `librarian/src/internal/application/links/check.go:254-299` calls
  `os.WriteFile` on the collected path in apply mode.
- `librarian/src/internal/platform/cli/links/clean/command.go:11-27` is
  the apply-capable CLI entry point.
- No symlink containment check appears between collection and write.

**Inferred (not yet proven):**

- The write preserves the file's content with only broken-link repairs; it
  does not inject arbitrary content.
- `--apply` must be explicitly passed; dry-run is the default.
- Impact is low because the attacker-controlled content written back is
  derived from a broken-link repair of the same file's Markdown.

**Likely scope:**

- `librarian/src/internal/application/links/check.go`
- Focused test for a symlinked `.md` entry under a scanned root.

## Root Cause

The file collector uses `!IsDir()` as a regular-file proxy, which also
admits symlinks to non-directory targets. No subsequent guard prevents the
apply path from writing to the resolved external target.

## Recommended Fix

- In the collection loop, use `entry.Type().IsRegular()` (or equivalent
  `Lstat`-based check) to accept only regular files; reject symlinks
  explicitly.
- Optionally verify that the canonical path remains beneath the scanned root
  before collection or writing.

## Verification Plan

- Add a test that places a `.md` symlink pointing to an external file with a
  broken link inside it, calls the apply path, and asserts the external file
  is unchanged.
- Run `go test ./...`, `go vet ./...`, `git diff --check`, and
  `go run ./cmd/hawp check` from `librarian/src`.

## Risk + Review Gate

**Risk:** low — apply-flag-gated path with constrained write content
**Gate:** standard review

## Backlog + Plan Link

**Status now:** plan-ready
**Plan file:** work/active/984f038a/plan.md

## Next Step

- [x] Investigation recorded above
- [x] Plan written with options, recommendation, and verification scope
- [x] Backlog moved to `plan-ready`
- [ ] Obtain explicit approval before implementation
