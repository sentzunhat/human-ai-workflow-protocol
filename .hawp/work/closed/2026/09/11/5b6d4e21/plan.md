# Provider parity for shared HAWP agent guidance

**UUID:** `5b6d4e21-18cf-4f1c-8dd6-4c92ce691211`
**Type:** improvement
**Reported:** 2026-09-11
**Status:** `done`

---

## Input

> Provider parity pass for all agent overlays. Take the Copilot upgrades that
> matter broadly: MCP-first/fallback discipline, repo-root path proof,
> batch-read-before-write, import-surgery caution, checkpoint-before-compaction,
> and validation commands. Push shared parts through
> `core/providers/shared/behaviors/**`, then run `providers sync` and
> `distribution sync` so Claude, Cursor, Continue, GitHub/Copilot, and where
> possible Codex all get the right equivalent. Provider-specific files can still
> differ; the behavioral contract should not.

## Intent

Bring the broad Copilot guidance upgrades into the shared provider behavior
system so every supported agent overlay receives equivalent HAWP operating
discipline through canonical source files and generated outputs.

## Scope

- Inspect `.github/copilot-instructions.md`,
  `core/providers/.github/copilot-instructions.md`, `AGENTS.md`, `CLAUDE.md`,
  `core/providers/shared/behaviors/**`, and generated provider overlays.
- Move provider-neutral behavior into `core/providers/shared/behaviors/**` when
  it belongs to all providers.
- Keep provider-specific wording in provider-owned source files when a tool
  needs different paths, frontmatter, or precedence rules.
- Run `go run ./cmd/hawp providers sync` and
  `go run ./cmd/hawp distribution sync` from `librarian/src`.
- Verify generated files are synchronized and no provider source pack drifts
  from the intended shared contract.

## Constraints

- Do not hand-edit generated provider overlays when the source behavior should
  generate them.
- Do not overwrite or reshape `.hawp/work/**` beyond this plan's status updates.
- Do not merge, publish, or tag releases.
- Preserve provider-specific differences where they are intentional.

## File Ownership

Primary:

- `core/providers/shared/behaviors/**`
- `core/providers/.claude/**`
- `core/providers/.cursor/**`
- `core/providers/.continue/**`
- `core/providers/.github/**`
- `core/providers/.codex/**`
- Generated provider overlays under `.github/`, `.claude/`, `.cursor/`,
  `.continue/`, and `distribution/generated/**`

Avoid changing install/update script contracts unless needed for parity; that
belongs to `d1fa0b72`.

## Verification

- `go run ./cmd/hawp providers sync`
- `go run ./cmd/hawp distribution sync`
- `go test ./...`
- `go vet ./...`
- `go run ./cmd/hawp check --no-update-check`
- `git diff --check`

## Coordination

Work branch: `codex/provider-parity-agent-guidance`

Manager note: this work should start from
`feature/v0.0.24-work-folder-normalization` and run in parallel with
`d1fa0b72`. If both branches touch generated distribution output, reconcile by
rerunning sync commands after merge.

**Closed:** 2026-09-11 — work completed in feature/v0.0.24-work-folder-normalization

## Outcome

Provider parity pass completed. Shared HAWP agent behaviors pushed through
`core/providers/shared/behaviors/`, `providers sync` and `distribution sync`
run. Claude, Cursor, Continue, GitHub/Copilot, and Codex overlays updated with
equivalent MCP-first/fallback discipline, repo-root path proof, batch-read-before-write,
import-surgery caution, checkpoint-before-compaction, and validation commands.

## Verification

- `go run ./cmd/hawp providers sync` — pass
- `go run ./cmd/hawp distribution sync` — pass
- `go test ./...` — pass
- `go vet ./...` — pass
- `go run ./cmd/hawp check --no-update-check` — pass

## Close Checklist

- [x] All focused checks pass.
- [x] HAWP check passes.
- [x] Plan status: done.
- [x] Plan moved to `.hawp/work/closed/2026/09/11/5b6d4e21/plan.md`.
- [x] BACKLOG.md updated.
