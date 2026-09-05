# Status Report

## Intent

Capture the post-PR #39 cleanup checkpoint and support the next branch-level
decision for HAWP docs/work organization.

## Current State

The `v0.0.23` release lane is closed on `main`. The cleanup branch was promoted
to a `v0.0.24` patch branch for repo-local documentation, UUID-folder
validation, and older-repo work normalization:
`feature/v0.0.24-work-folder-normalization`.

The active cleanup item is now tracked as UUID `e5fca9c7-ba0a-424e-b0a5-ef725786a99e`
at `.hawp/work/active/e5fca9c7/plan.md`.

## What Was Inspected

- Recent commit history
- `.hawp/work/BACKLOG.md`
- `.hawp/work/STATUS.md`
- `core/.hawp/work/README.md`
- `.hawp/kit/references/docs-alignment.md`
- `.hawp/work/parked/89cf7a85/plan.md`
- Recursive diff between `.hawp/kit` and `core/.hawp/kit`
- Recursive diff between `.hawp/work` and `core/.hawp/work`

## What Changed

- Added an active cleanup plan under `.hawp/work/active/`.
- Updated `.hawp/work/BACKLOG.md` to show the current cleanup branch and active
  maintenance item.
- Refreshed `.hawp/work/STATUS.md`, replacing the stale June closeout pointer
  with the current post-release state and next queue.
- Revised `README.md` to lead with the builder/open-source positioning while
  keeping benchmark claims tied to local evidence artifacts.
- Moved the cleanup plan from a descriptive active folder into the canonical
  UUID-folder layout.
- Advanced CLI version metadata and checked-in binary to `0.0.24`.
- Added conservative apply cleanup for stale completed Active Work rows and
  unreferenced active/parked duplicates when a closed canonical record exists.
- Added Mochila agent instructions at
  `.hawp/work/active/e5fca9c7/mochila-agent-instructions.md`.

## What Was Directly Verified

- `git log` shows `94a66001` as current `main`, immediately after `05aa28ce`.
- `.hawp/kit` and `core/.hawp/kit` have no recursive diff output.
- `.hawp/work` and `core/.hawp/work` differ, as expected, because root work
  state contains repository history while core work state is scaffold source.
- `go run ./cmd/hawp work validate` passed with 0 issues and 0 warnings.
- `go run ./cmd/hawp kit validate` passed with 0 issues.
- The README improvement input on 2026-09-05 provided candidate messaging and
  practical next targets for the older-repo compatibility lane.
- After README edits, `go run ./cmd/hawp work validate` still passed with 0
  issues and 0 warnings.
- After README edits, `go run ./cmd/hawp kit validate` still passed with 0
  issues.
- The initial descriptive active folder was a repo-local workflow drift issue,
  not a CLI scaffold issue: `hawp work new` already emits short-UUID folders.
- The validator now fails active rows that use non-canonical new-work IDs while
  preserving legacy `TASK-NNN` and numeric-row compatibility.
- Read-only validation against `mochila-archive-viewer` reached the intended
  older-repo compatibility path: numeric active rows remain tolerated, while
  stale done-active rows and incomplete closed records are reported.
- Temporary-copy Mochila apply proof passed after default normalize plus
  folder migration: final validation reported 0 issues and 0 warnings.
- Temporary-copy Tekit apply proof passed after default normalize plus folder
  migration: final validation reported 0 issues and 2 tolerated warnings.
- Full Go suite passed with `go test ./...`.
- Checked-in binary reports `0.0.24`.

## What Remains Unproven

- `mochila-archive-viewer` has unrelated dirty changes, so any migration there
  should happen on its own branch after preserving that state.
- Target repositories were not modified; only temporary copies were used for
  apply proof.

## Constraints

- Preserve repo-owned `.hawp/work/**` history.
- Do not rewrite old closed archives for cosmetic consistency.
- Do not hand-edit generated distribution guides.

## Help Wanted

Review whether `89cf7a85` should become the next active maintenance lane or
remain parked until there is a target older repo to test against.

## Suggested Next Step

Review the `v0.0.24` patch diff, then use the Mochila agent instructions to
apply and inspect the target-repo cleanup on its own branch without committing.
