# Check index existence before opening SQLite

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
**UUID:** `80c677fb-9958-49ec-9232-05794b91e7dc`
**Type:** bug
**Reported:** 2026-09-23

## Input (verbatim)

> Check index existence before opening SQLite

## Intake Summary

The Copilot review identifies incorrect first-run and error handling in `hawp search
embed`. The command calls `sqlite.Open` before checking whether the index exists; that
constructor can create a missing database and parent directory, after which the command
reports a schema error or misleadingly treats unrelated failures as "index not found".

## Current Context

The command should distinguish a missing index from permission, corruption, and other
open failures. It must not create a database during a read/embedding preflight, and it
must preserve the existing successful embedding path.

## Initial Analysis

**Directly verified:**

- `librarian/src/internal/platform/cli/model/search-embed/command.go:29-35` opens the
  resolved path before an existence check and returns nil for every open error.
- The supplied review identifies `sqlite.Open` as creating the missing path.

**Inferred (not yet proven):**

- `os.Stat`/`os.Lstat` before `sqlite.Open` can provide deterministic missing-index
  guidance; after existence is confirmed, open errors should be returned.

**Likely scope:**

- `librarian/src/internal/platform/cli/model/search-embed/command.go`, focused command
  tests or filesystem seam tests, and the SQLite adapter contract.

## Root Cause

Database construction is being used as an existence probe, and all constructor errors
are collapsed into a successful "not found" message.

## Options Considered

1. Change `sqlite.Open` globally to stop creating databases. This risks breaking index
   creation callers and expands the change beyond this command's read preflight.
2. Check the path before opening, then propagate open errors after existence is confirmed.
   This is the recommended scoped fix.

## Recommended Fix

Use a pre-open existence check. Print the first-run guidance and return nil only for a
confirmed missing path; return a contextual error for permission, directory, or
corruption/open failures. Close the database on all successful-open paths.

## Risk + Review Gate

**Risk:** medium (CLI error-contract change)
**Gate:** user authorized implementation in the request; no merge authorization.

## Backlog + Plan Link

**Status now:** done
**Plan file:** work/active/80c677fb/plan.md

## Verification

- Test a missing index, an existing usable index, and an open failure.
- Run focused search-embed tests plus all Go tests and vet.
- Confirm the missing-index path does not create `.hawp/db/index.sqlite`.

## Next Step

- [x] Investigation recorded above
- [x] Write or update the plan file
- [x] Create the corresponding fix work item and implement sequentially

## Outcome

Done. `search embed` performs a read-only existence check before opening the
SQLite index, returns first-run guidance only for a confirmed missing path,
and preserves contextual open errors.
