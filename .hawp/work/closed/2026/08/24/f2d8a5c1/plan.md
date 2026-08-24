---
work-item: f2d8a5c1
type: fix
title: "`hawp update --check` performs the update instead of being read-only"
status: done
owner: unassigned
created: 2026-08-24
updated: 2026-08-24
closed: 2026-08-24
---

# `hawp update --check` — Read-Only Mode Bug

## Mission

Make `hawp update --check` a non-destructive read check: compare the running
binary version against the latest GitHub release and report whether an update
is available, without downloading or installing anything.

## Context

As of 0.0.4, `hawp update --check` triggers the full update flow rather than
being read-only. The correct behavior for a `--check` flag is to exit 0 when
up-to-date and exit 1 when an update is available, with no side effects.
`hawp update verify` already exists as a dedicated read-only subcommand; the
`--check` flag on `hawp update` should behave identically or alias it.

Observed: `hawp update --check` performs the update (same behavior as `hawp update`).
Expected: exits without downloading; prints availability status to stdout.

## Fix

Root cause: `hawp update --check` fell through to `runUpdateFull` because `--check`
was neither a documented flag nor a routing case. The flag was silently ignored
and the full update ran.

Fix: added a routing case in `run.go` before the bare `update` catch-all:

```go
case command == "update" && containsArg(args[1:], "--check"):
    return runUpdateVerify()
```

`runUpdateVerify` (already used by `hawp update verify`) is now the shared
implementation: read-only version check, exits 0 if up-to-date, exits 1 if an
update is available.

Also updated registry and help text to document `--check` as an official flag.

## Outcome

Fixed in `librarian/src/internal/platform/cli/run.go`: added routing case before
the bare `update` catch-all so `hawp update --check` dispatches to `runUpdateVerify()`
instead of `runUpdateFull`. Registry and help text updated to document `--check` as
an official flag. No logic changes to the verification path itself.

## Verification

- [x] `go build` clean (local dev build, 2026-08-24)
- [x] `go test ./...` — all packages pass (2026-08-24)
- [x] `hawp commands --json` reflects `--check` in `update` flags (2026-08-24)

## Close Checklist

- [x] Implementation complete on branch `feature/v0.0.5`
- [x] Tests pass (Go)
- [x] BACKLOG updated (moved to Recently Closed)
- [x] Plan file updated to `status: done`
