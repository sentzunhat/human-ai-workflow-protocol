# Constrain project search index paths to the repository root

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
**UUID:** `0c8d6391-4b07-483f-a8a6-8487a8b3e539`
**Type:** bug
**Reported:** 2026-09-20

---

## Input (verbatim)

> Constrain project search index paths to the repository root

## Intake Summary

The tracked project configuration `.hawp/config/search.json` can replace the
default index path list. Those strings are joined to the repository root and
used without a containment check. A contributor-controlled `../secret.md` or
outside directory can therefore be read when a user runs `hawp search index`,
persisted into `.hawp/db/index.sqlite`, and returned later through CLI or MCP
search.

## Current Context

This is a local confused-deputy boundary: repository configuration selects
what the user's HAWP process reads. Documentation says paths are relative to
the repository root, but the implementation does not enforce that contract.

## Initial Analysis

**Directly verified:**

- `librarian/src/internal/application/search/config.go:43-66` accepts project
  `index.paths` without path validation.
- `librarian/src/internal/platform/cli/index/ingest/corpus.go:18-33` joins each
  configured value to `repoRoot` and dispatches it without a `filepath.Rel`
  containment check.
- `librarian/src/internal/platform/cli/index/ingest/corpus.go:113-164` reads a
  selected file or recursively walks a selected directory.
- `librarian/src/internal/application/index/ingest-service.go:84-163` chunks
  the resulting content and persists it in the repository index.
- `librarian/src/internal/platform/mcp/server/tool_search.go:76-128` returns
  matching stored chunk content to MCP callers.

**Inferred (not yet proven):**

- No automatic transmission to the repository supplier was established; the
  user must run indexing and a later consumer must query matching content.
- The `.md` suffix restriction reduces the default path-traversal surface but
  does not restore the repository boundary.

**Likely scope:**

- `librarian/src/internal/application/search/config.go`
- `librarian/src/internal/platform/cli/index/ingest/corpus.go`
- `librarian/src/internal/platform/cli/index/ingest/corpus_test.go`
- Adjacent command/config tests needed to preserve home-level and project-level
  valid relative path behavior.

## Root Cause

The application treats the documented "relative to repository root" rule as a
caller obligation. `filepath.Join` cleans traversal components but does not
prove that the result remains below `repoRoot`, and no later layer restores
that missing invariant before `os.Stat`, `os.ReadFile`, or `filepath.Walk`.

## Options Considered

1. Validate configured path strings only when loading JSON. Simple, but it
   couples repository-specific policy to a config loader also used for home
   defaults and does not protect future callers.
2. Resolve and validate every path at the corpus-ingest boundary. This keeps
   the security invariant beside the filesystem access and protects all config
   sources. **Recommended.**
3. Remove custom index paths. Strongest restriction, but breaks a documented
   feature unnecessarily.

## Recommended Fix

- Add one repository-root path resolver used by `buildCorpusFromRepo` before
  any filesystem access.
- Reject empty and absolute configured paths, then use cleaned absolute paths
  plus `filepath.Rel` to reject `..` escapes on each supported platform.
- Keep symlink resolution and non-regular-file policy in `225708ae`; share a
  helper only if the resulting ownership remains clear.
- Return an actionable error naming the rejected configured value without
  reading it.

## Verification Plan

- Add regression tests for `../outside.md`, an outside directory, absolute
  paths, and cleaned nested traversal.
- Preserve tests for valid `.hawp/kit`, `.hawp/work`, custom directories, and
  single Markdown files.
- Run focused ingest/config tests, `go test ./...`, `go vet ./...`,
  `git diff --check`, `go run ./cmd/hawp check`, and
  `go run ./cmd/hawp work validate` from `librarian/src`.

## Risk + Review Gate

**Risk:** high — security boundary change in a shared indexing path
**Gate:** explicit review before implementation

**Overlap check:** no current product-source edits overlap this planned search
scope. Preserve unrelated dirty changes.

## Backlog + Plan Link

**Status now:** done
**Plan file:** work/active/0c8d6391/plan.md

## Next Step

- [x] Investigation recorded above
- [x] Plan written with options, recommendation, and verification scope
- [x] Backlog moved to `plan-ready`
- [x] Obtain explicit approval before implementation

## Outcome

Done. `buildCorpusFromRepo` resolves each configured path relative to the
repository and rejects empty, absolute, and cleaned traversal paths before
filesystem access. Regression coverage covers outside files, directories,
and valid corpus paths.
