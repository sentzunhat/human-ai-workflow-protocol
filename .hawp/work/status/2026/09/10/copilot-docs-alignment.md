# Status Report

## Intent

Summarize the 2026-09-10 work state, check whether HAWP records or docs drifted,
and align Copilot-facing docs using `.github/copilot-instructions.md` and the
HAWP change-review/reference-sync prompt.

## Current State

The architecture cleanup and docs consolidation work is mostly aligned. The main
drift found was record and propagation drift, not implementation drift:

- `8ddb06ea` is complete, but an active duplicate plan remained after the item
  was also closed.
- `multi-repo-context-d9b2f3a1` is complete, but an active duplicate plan
  remained after the item was also closed.
- Three docs tasks were removed from active backlog rows after being folded into
  the Copilot instruction update, but their active plan files had not been moved
  to closed records.
- `.github/copilot-instructions.md` had newer Copilot validation/checkpoint/MCP
  guidance that was not present in `core/providers/.github/copilot-instructions.md`,
  even though the GitHub provider update path refreshes downstream
  `.github/copilot-instructions.md` from that source pack.

## What Was Inspected

- `.hawp/kit/start-here.md`
- `.hawp/kit/usage/status-report.md`
- `.hawp/kit/references/docs-alignment.md`
- `.github/prompts/hawp-change-review-and-reference-sync.prompt.md`
- `.github/copilot-instructions.md`
- `core/providers/.github/copilot-instructions.md`
- `.github/skills/multi-repo-context.md`
- `core/providers/manifest.yaml`
- `.hawp/work/BACKLOG.md`
- active and closed plans for `8ddb06ea`, `multi-repo-context-d9b2f3a1`,
  `cli-rg-quoting-b9f3e1d2`, `session-checkpoint-c4a8d7e1`, and
  `validation-shortcut-e7c6b5d4`

## What Changed

- Added the multi-repo context skill cue to `.github/copilot-instructions.md`.
- Propagated the current Copilot instruction guidance to
  `core/providers/.github/copilot-instructions.md`.
- Moved the three consolidated docs-task plans from `.hawp/work/active/` to
  `.hawp/work/closed/2026/09/10/` with outcome, verification, and close
  checklist sections.
- Removed duplicate active copies for already-closed `8ddb06ea` and
  `multi-repo-context-d9b2f3a1`.
- Excluded the embedded `.hawp/.spaces/` agent-worktree mirrors from the link
  checker (`librarian/src/internal/application/links/check.go`) and added them
  to `.gitignore`, so `hawp check` no longer fails on the frozen snapshots'
  internal links.

## What Was Directly Verified

- HAWP MCP tools were available and used: `hawp_search` and
  `hawp_work_validate`.
- `cd librarian/src && go test ./...` passed.
- `cd librarian/src && go run ./cmd/hawp distribution sync` reported generated
  distribution outputs current and provider materialization current.
- `cd librarian/src && go run ./cmd/hawp work validate` passed after the record
  cleanup.
- `core/providers/manifest.yaml` maps the GitHub provider pack to downstream
  `.github/copilot-instructions.md` refresh behavior.
- `.github/copilot-instructions.md` contained validation, checkpoint, MCP, path,
  and agent-discipline guidance before this pass.
- `core/providers/.github/copilot-instructions.md` lacked that newer guidance
  before this pass.
- `8ddb06ea` and `multi-repo-context-d9b2f3a1` had closed plans and active
  duplicate plan files.

## What Remains Unproven

- Nothing blocking. `hawp check` now passes all three gates (kit validate,
  work validate, links check) after the `.spaces` exclusion.

## Decision

- `.hawp/.spaces/**` mirrors were excluded from link validation (kept as
  reference snapshots) rather than deleted. Rationale: they are 98MB embedded
  gitlink snapshots, not live docs; deleting them is a separate hygiene
  decision that loses reference data. Link-check noise is the only symptom, so
  fixing the symptom is the minimal change.
- The mirrors remain git-tracked gitlinks; `.gitignore` prevents new untracked
  content there but does not untrack the existing two gitlinks.

## Constraints

- Preserve unrelated dirty work.
- Keep docs edits limited to evidence-backed alignment.
- The one source change (links `skipDirs`) is a checker-scope change, not a
  runtime behavior change; it does not alter what `hawp check` validates about
  live docs.

## Verification (final state)

- `go build ./...` — OK
- `go vet` (links package) — OK
- `go test ./...` — pass
- `go run ./cmd/hawp links check` — 144 files, local links valid
- `go run ./cmd/hawp check` — all 3 validations pass
- `go run ./cmd/hawp work validate` — VALIDATION PASS
- `go run ./cmd/hawp distribution sync` — 0/40 updated (current)
