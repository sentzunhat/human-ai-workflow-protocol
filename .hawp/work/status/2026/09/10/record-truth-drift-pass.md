# Record-truth drift pass — v0.0.24 lane

**Date:** 2026-09-10
**Branch:** `feature/v0.0.24-work-folder-normalization`
**Type:** record correction + validation
**Parent commit (pre):** `1bc8e3e8`

---

## Trigger

Next-compounding-action review surfaced that two active plans describe code
that is not in the committed tree. Correcting record truth to match verified
reality is the precondition for either cutting the `0.0.24` release or
continuing implementation, so this pass ran first.

## Findings (direct evidence)

### `a3df8a9c` — reshape support (overstated)

- **Confirmed present:** `librarian/src/internal/domain/work/draft.go`
  (committed as `d3559556`), defines `Draft` + completeness `Validate()`.
- **Confirmed absent:** `librarian/src/internal/application/work/draft.go`
  (the `RequestShaper` port + service) — `ls` → no such file.
- **Confirmed absent:** `librarian/src/internal/application/work/draft_test.go`
  — `ls` → no such file.
- **Confirmed unused:** `grep -rln "work.Draft\|RequestShaper"` across
  `librarian/src/internal/` returned no references outside the domain value.

The plan's "Draft contract verification" section and a checked `[x]` box had
claimed the application-side contract and tests were implemented and passing.
That did not match the committed tree (working tree clean; `git log` for the
application files empty).

### `47c793d6` — CLI decomposition (stale box)

- `[ ] Work architecture layer moves: 742aa60b` — but
  `closed/2026/09/10/742aa60b/plan.md` exists; the item is closed. Box marked
  done to match reality.

## What changed

- `.hawp/work/active/a3df8a9c/plan.md` — checklist split so only the landed
  domain `Draft` value is checked; application port + tests are now an explicit
  open box; verification section rewritten to state what exists, what does not,
  and to withdraw the uncommitted test-coverage claims.
- `.hawp/work/active/47c793d6/plan.md` — `742aa60b` box marked done (closed
  2026-09-10).

## Verification

Run from `librarian/src` (repo root `<repo-root-abs>`):

- `go build ./...` — pass
- `go vet ./...` — pass
- `go test ./...` — pass (full suite, no failures)
- `go run ./cmd/hawp work validate` — VALIDATION PASS (5/5 checks, 0 issues)
- `go run ./cmd/hawp kit validate` — 3 checks pass
- `go run ./cmd/hawp links check` — 144 files valid
- `go run ./cmd/hawp check --no-update-check` — all 3 validations passed

## Out of scope / not claimed

- No source code changes in this pass (the domain `Draft` value is pre-existing
  and untouched).
- The application `RequestShaper` port, service, and test suite remain open
  work under `a3df8a9c`; nothing was implemented here.
- No release tag was cut; `0.0.24` does not yet exist (highest tag `0.0.23`).
  Cutting the release is a separate decision.
- Branch not pushed to origin (no upstream tracking).
