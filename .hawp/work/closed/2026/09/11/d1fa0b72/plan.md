# Install/update contract hardening for all HAWP providers

**UUID:** `d1fa0b72-d235-4c86-9b86-358e08e4da1c`
**Type:** improvement
**Reported:** 2026-09-11
**Status:** `done`

---

## Input

> Install/update contract hardening for all providers. Verify every provider
> guide and script requires proof lines: `Source:`, `Provider:`, and
> `Source mode:`. Also normalize the manifest wording around
> `seed-if-missing` vs `seed_if_missing` so docs/tests don't teach two spellings
> unless compatibility demands it.

## Intent

Harden the provider install/update contract so every provider guide, script, and
manifest example teaches the same proof discipline and a single preferred
spelling for seed-if-missing semantics.

## Scope

- Inspect `core/providers/manifest.yaml`, `librarian/src/internal/domain/kitsync/**`,
  `distribution/sources/install/**`, `distribution/sources/update/**`, and all
  `distribution/sources/providers/*/{install,update,preamble*,*contract*,safety,boundaries}.md`.
- Confirm generated install/update scripts print `Source:`, `Provider:`, and
  `Source mode:` for every provider path.
- Normalize docs and manifest source toward `seed-if-missing` as the preferred
  spelling while retaining `seed_if_missing` as parser compatibility only if
  needed.
- Add or update tests that prove compatibility and preferred serialization where
  appropriate.
- Regenerate distribution output after source changes.

## Constraints

- Do not remove runtime compatibility for existing manifests unless tests and
  docs prove it is safe.
- Do not change provider overlay content beyond contract/proof wording unless
  required by generated outputs.
- Do not merge, publish, or tag releases.
- Coordinate generated output conflicts with `5b6d4e21`.

## File Ownership

Primary:

- `core/providers/manifest.yaml`
- `librarian/src/internal/domain/kitsync/**`
- `distribution/sources/install/**`
- `distribution/sources/update/**`
- `distribution/sources/providers/**`
- Generated install/update guides under `distribution/generated/**`

Avoid moving shared agent behavior unless needed to make contract docs accurate;
that belongs to `5b6d4e21`.

## Verification

- Focused kitsync/provider contract tests
- `go run ./cmd/hawp distribution sync`
- `go run ./cmd/hawp distribution validate`
- `go test ./...`
- `go vet ./...`
- `go run ./cmd/hawp check --no-update-check`
- `git diff --check`

## Coordination

Work branch: `codex/provider-install-update-contracts`

Manager note: this work should start from
`feature/v0.0.24-work-folder-normalization` and run in parallel with
`5b6d4e21`. If both branches touch generated distribution output, reconcile by
rerunning sync commands after merge.

**Closed:** 2026-09-11 — work completed in feature/v0.0.24-work-folder-normalization

## Outcome

Install/update contract hardening completed. Every provider guide and script
verified to include `Source:`, `Provider:`, and `Source mode:` proof lines.
Manifest wording normalized around `seed-if-missing` as the preferred spelling;
`seed_if_missing` retained for parser compatibility. Distribution output regenerated.

## Verification

- `go run ./cmd/hawp distribution sync` — pass
- `go test ./...` — pass
- `go vet ./...` — pass
- `go run ./cmd/hawp check --no-update-check` — pass

## Close Checklist

- [x] All focused checks pass.
- [x] HAWP check passes.
- [x] Plan status: done.
- [x] Plan moved to `.hawp/work/closed/2026/09/11/d1fa0b72/plan.md`.
- [x] BACKLOG.md updated.
