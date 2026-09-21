# doc: quote YAML front-matter title to prevent newline injection

**Backlog ID (Legacy):** — (UUID-native item)
**UUID:** `cd5c215b-66af-420e-8ecb-b62930bc5c6c`
**Type:** fix
**Reported:** 2026-09-15

---

## Input (verbatim)

> doc: quote YAML front-matter title to prevent newline injection

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
**Plan file:** work/active/cd5c215b/plan.md

## Verification

- `go test ./...` passes; `go vet ./...` clean.
- `template()` now renders `title: "My title"` (double-quoted); a title
  containing `\n` is escaped by `%q`, preventing front-matter field injection.

## Outcome

Changed `title: %s` to `title: %q` in `doc/doc.go`'s `template()` function.
A title containing a newline character would previously have been written as a
raw newline inside the YAML scalar, splitting the front-matter block. `%q`
escapes special characters so the value stays within the quoted field.

## Close Checklist

- [x] YAML injection vector closed — title is now double-quoted with Go escaping.
- [x] `go test ./...` passes; `go vet ./...` clean.
- [x] Plan moved to `closed/2026/09/15/cd5c215b/`.
- [x] BACKLOG.md updated.
