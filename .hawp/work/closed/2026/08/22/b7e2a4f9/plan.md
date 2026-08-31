---
work-item: b7e2a4f9
type: refactor
title: "Rename librarian/go/ → src/ and update all references"
status: plan-ready
owner: unassigned
created: 2026-08-22
updated: 2026-08-22
---

# Rename librarian/go/ → src/

## Mission

Rename the `librarian/go/` directory to `src/` to simplify the project
layout and reduce redundant path nesting. Update every reference across
docs, CI, Makefile, CHANGELOG, and provider overlays.

## Why

`librarian/go/` made sense when multiple language implementations were
anticipated. With Go as the only implementation, the double nesting
(`librarian/go/cmd/`, `librarian/go/internal/`) is noise. `src/` is
the standard Go project convention for a single-language repo.

## Scope

### Files / paths to update

| Location | What changes |
|---|---|
| Directory itself | `git mv librarian/go src` |
| `.github/workflows/release-librarian-go.yml` | `working-directory: librarian/go` → `src` |
| `.github/workflows/release-librarian-go.yml` | `librarian/go/bin/dist/*` → `src/bin/dist/*` |
| `.github/workflows/release-librarian-go.yml` | changelog path `librarian/go/CHANGELOG.md` → `src/CHANGELOG.md` |
| `CLAUDE.md` | `npm --prefix librarian run …` commands remain (they call the librarian JS tooling, not Go) |
| `CLAUDE.md` | Project layout table row for `librarian/go/` |
| `librarian/docs/` | All paths referencing `librarian/go/` |
| `librarian/scripts/` | Any path references |
| `distribution/` | Any path references |
| `.hawp/kit/` | Any install/update guides referencing `librarian/go/` |
| `core/providers/` | Provider overlays referencing the path |
| Memory / status files | Update memory index to reflect new path |

### What does NOT change

- `librarian/` directory itself (the JS tooling lives there)
- `npm --prefix librarian` scripts (JS, not Go)
- Go module path inside `src/go.mod` (keep as-is unless explicitly changing module path)
- Binary names (`hawp`, `hawp-*`) — these are user-facing and stable

## Plan

1. `git mv librarian/go src`
2. Update `.github/workflows/release-librarian-go.yml` paths
3. Audit remaining references: `grep -r "librarian/go" . --include="*.md" --include="*.yml" --include="*.yaml" --include="*.json" --include="*.go"`
4. Fix each hit
5. Run `npm --prefix librarian run validate` to confirm no broken links
6. Run `go build ./...` from `src/` to confirm Go still builds
7. Commit on development branch; open PR to main

## Verification

- [ ] explicitly unproven in this record: `git mv` succeeds, directory is `src/`
- [x] GitHub Actions workflow references updated (4 workflow files). Evidence: see Outcome section above.
- [x] `npm --prefix librarian run check:markdown-links` passes. Evidence: see follow-up verification notes in this plan.
- [x] `go build ./...` passes from `src/`. Evidence: see follow-up verification notes in this plan.
- [x] No stale `librarian/go` hits in audited file types. Evidence: see Outcome section above.
- [x] `npm --prefix librarian run validate` (Go-backed): kit ✓, work ✓, links ✓. Evidence: see follow-up verification notes in this plan.

## Outcome

`librarian/go/` renamed to `librarian/src/` via `git mv`. All CI workflows, docs, and
CLAUDE.md updated. npm scripts for hawp validation (kit:validate, work:validate,
hawp:check, check:markdown-links, work:normalize:cli) now route to
`../.hawp/bin/hawp` — TS validators retired. Distribution/providers scripts
remain in TypeScript (not yet ported). All checks pass.

## Verification

- `go build ./...` from `src/`: clean
- `npm --prefix librarian run kit:validate`: PASS (Go)
- `npm --prefix librarian run work:validate`: PASS (Go)
- `npm --prefix librarian run check:markdown-links`: PASS (Go)

## Close Checklist

- [x] Plan file complete with Outcome and Verification
- [x] Work item moved to closed/
