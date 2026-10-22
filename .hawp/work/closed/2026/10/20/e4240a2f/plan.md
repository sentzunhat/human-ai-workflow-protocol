# Compact BACKLOG.md and archive closed work.

**Backlog ID (Legacy):** — (UUID-native item)
**UUID:** `e4240a2f-f989-42dd-8d58-17f2b320f88a`
**Type:** improvement
**Reported:** 2026-10-20

---

## Input (verbatim)

> Compact the oversized Recently Closed section of `.hawp/work/BACKLOG.md` before adding another closed work item. Preserve older closed history in the archive.

## Intake Summary

Reduce the Recently Closed table to the repository's stated 5–10-row cap while preserving older closed work in the archive section.

## Current Context

`.hawp/work/BACKLOG.md` currently lists 16 Recently Closed items despite stating that the table is limited to 5–10 items. Eight older rows can be archived without moving or deleting their closed plan files. Two active tasks from 2026-10-20 (this compaction and agent-copy cleanup) are expected to close, so the compacted table should retain eight existing rows and leave room for both new closures.

## Initial Analysis

**Directly verified:**

- The current Recently Closed table has 16 entries and states a 5–10-item limit.
- The eight oldest entries selected for archival link to existing closed plan files; all eight exact paths were checked.
- HAWP backlog guidance requires preserving closed history and allows archiving old Recently Closed rows.
- The agent-copy consolidation item `81d6a35c` and this compaction item are both active, so retaining eight existing recent rows leaves a total of ten after both close.

**Inferred (not yet proven):**

- Keeping the eight most recent existing rows and adding the two current closures will satisfy the repository's selected cap without changing any closed plan records.

**Likely scope:**

- Update `.hawp/work/BACKLOG.md` by moving eight old rows into a dated archived Recently Closed table and keep the eight newest existing rows.
- Update this plan and close it.

## Risk + Review Gate

**Risk:** low
**Gate:** auto-implement on low

## Backlog + Plan Link

**Status now:** closed
**Plan file:** work/closed/2026/10/20/e4240a2f/plan.md

## Next Step

- [x] Investigation recorded above (required before planning)
- [x] Write or update the plan file
- [x] Move backlog status accordingly

## Options

### Option A — Archive the eight oldest rows

Keep the eight newest existing Recently Closed rows and move the eight older rows into an archived Recently Closed table. Preserve their existing closed-plan links. This leaves room for the two 2026-10-20 items expected to close.

### Option B — Keep all rows in Recently Closed

Rejected: this would leave the table above its stated 5–10-row cap.

### Option C — Remove old rows without archiving them

Rejected: this would discard active-index history even though the closed plans exist.

## Recommended Fix

**Option chosen:** A

**Rationale:** Archiving old index rows preserves permanent history and keeps the Recently Closed table within the selected cap after the two current tasks close.

**Files to change:**

- `.hawp/work/BACKLOG.md` — move eight older rows into `Archived Recently-Closed (compacted 2026-10-20)`, keep eight newer existing rows, then add the two new closures.
- `.hawp/work/active/e4240a2f/plan.md` — record outcome and verification, then close.

**What to verify after:**

- [x] All eight archived plan links still point to existing closed plans.
- [x] Recently Closed contains ten or fewer rows after both tasks are closed.
- [x] HAWP validation and `git diff --check` pass.

## Work Coordination

**Owner:** maintainer
**Implementation status:** done
**Overlapping files:**

- `.hawp/work/BACKLOG.md`
- `.hawp/work/active/e4240a2f/plan.md`
- `.hawp/work/active/81d6a35c/plan.md`

**Parallel work risk:** low
**Can implement now:** yes

**Coordination note:** This compaction is a prerequisite to closing the agent-copy cleanup item. No closed plan file will be moved or deleted.

**Path discipline:** All paths are exact repo-relative paths.

**Repo-root proof (captured before archive edits):**

```text
$ pwd
<repo-root-abs>/human-ai-workflow-protocol
$ git rev-parse --show-toplevel
<repo-root-abs>/human-ai-workflow-protocol
$ git rev-parse --show-prefix

$ git status --short
 M .github/copilot-instructions.md
 D .github/prompts/hawp-backlog-alignment.prompt.md
 D .github/prompts/hawp-change-review-and-reference-sync.prompt.md
 D .github/prompts/hawp-commit-many-small.prompt.md
 D .github/prompts/hawp-commit-one-big.prompt.md
 D .github/prompts/hawp-conservative-docs-drift-cleanup.prompt.md
 D .github/prompts/hawp-docs-alignment-deterministic.prompt.md
 D .github/prompts/hawp-docs-alignment-simplicity.prompt.md
 D .github/prompts/hawp-intent-first-handoff.prompt.md
 D .github/prompts/hawp-status-report.prompt.md
 D .github/prompts/intake.prompt.md
 M .gitignore
 M .hawp/kit/usage/mcp/README.md
 M .hawp/kit/usage/mcp/github-copilot.md
 M .hawp/kit/usage/mcp/lessons.md
 M .hawp/kit/usage/mcp/verification.md
 M .hawp/work/BACKLOG.md
 D .hawp/work/active/1607037f/plan.md
 D .hawp/work/active/2eea565c/plan.md
 D .hawp/work/active/429e075e/plan.md
 D .hawp/work/active/7394f90c/plan.md
 D .hawp/work/active/96c45719/plan.md
 D .hawp/work/active/b4300599/plan.md
 M .hawp/work/active/c8d9e1cb/plan.md
 D .hawp/work/parked/usage-tracking/plan.md
 M core/providers/.claude/rules/hawp-core.md
 M core/providers/.codex/AGENTS.md.seed
 M core/providers/.continue/rules/hawp-01-core.md
 M core/providers/.cursor/rules/hawp-core.mdc
 M core/providers/.github/copilot-instructions.md
 M core/providers/doc-verification.md
 M core/providers/shared/behaviors/hawp-core.md
 M distribution/generated/claude/install/development.md
 M distribution/generated/claude/install/main.md
 M distribution/generated/claude/update/development.md
 M distribution/generated/claude/update/main.md
 M distribution/generated/codex/install/development.md
 M distribution/generated/codex/install/main.md
 M distribution/generated/codex/update/development.md
 M distribution/generated/codex/update/main.md
 M distribution/generated/continue/install/development.md
 M distribution/generated/continue/install/main.md
 M distribution/generated/continue/update/development.md
 M distribution/generated/continue/update/main.md
 M distribution/generated/cursor/install/development.md
 M distribution/generated/cursor/install/main.md
 M distribution/generated/cursor/update/development.md
 M distribution/generated/cursor/update/main.md
 M distribution/generated/github/install/development.md
 M distribution/generated/github/install/main.md
 M distribution/generated/github/update/development.md
 M distribution/generated/github/update/main.md
 M distribution/sources/shared/install.md
 M distribution/sources/shared/update.md
 M librarian/src/internal/platform/mcp/configure/config.go
 M librarian/src/internal/platform/mcp/configure/configure.go
 M librarian/src/internal/platform/mcp/configure/configure_test.go
 M librarian/src/internal/platform/mcp/server/server.go
?? .github/agents/hawp.agent.md
?? .github/skills/commit-many-small/
?? .github/skills/commit-one-big/
?? .github/skills/conservative-docs-drift-safe-cleanup/
?? .github/skills/docs-alignment-deterministic/
?? .github/skills/docs-alignment-simplicity-pass/
?? .github/skills/hawp-backlog-alignment/
?? .github/skills/hawp-change-review-and-reference-sync/
?? .github/skills/intake/
?? .github/skills/intent-first-handoff/
?? .github/skills/status-report/
?? .hawp/work/active/3c137bad/
?? .hawp/work/active/3cf52223/
?? .hawp/work/active/81d6a35c/
?? .hawp/work/active/e4240a2f/
?? .hawp/work/active/9dce1d95/
?? .hawp/work/active/b152c56c/
?? .hawp/work/active/e1dc28ac/
?? .hawp/work/closed/2026/10/16/
?? .hawp/work/closed/2026/10/17/
?? .hawp/work/parked/1607037f/
?? .hawp/work/parked/7394f90c/
?? .hawp/work/parked/96c45719/
?? .hawp/work/parked/b4300599/
?? librarian/src/internal/platform/mcp/configure/config_agent.go
?? librarian/src/internal/platform/mcp/server/prompts.go
?? librarian/src/internal/platform/mcp/server/prompts_test.go
?? librarian/src/internal/platform/mcp/server/resources.go
?? librarian/src/internal/platform/mcp/server/resources_test.go
```

## Implementation Notes

Move rows only; preserve each archived closed-plan file and its link. Keep the remaining Recently Closed entries and the two new closures within the 10-row cap.

## Outcome

Moved the eight oldest rows from Recently Closed into a new dated archive table, preserving their existing closed-plan links. The eight newest prior rows remain, leaving room for this compaction and the agent-copy cleanup closures while keeping the section at its 10-row cap.

## Verification

- [x] All eight archived closed-plan paths exist. **Evidence:** each path was checked before editing `.hawp/work/BACKLOG.md`.
- [x] Recently Closed contains 10 rows after both closures. **Evidence:** the final `awk` count returned `10`.
- [x] `git diff --check -- .hawp/work/BACKLOG.md` passed.
- [x] Final `hawp-hawp_work_validate`: kit validation passed (3 checks, 0 issues); work validation passed (0 issues, 1 warning); local links passed (150 Markdown files); overall `hawp check` passed all 3 validations.
- [x] Detailed `cd librarian/src && go run ./cmd/hawp work validate`: backlog consistency, closed-task completeness (255/255), evidence integrity (9/9), and dead links passed; verification clarity warned (309/390 evidence-backed, 71 explicitly unproven, 10 ambiguous). Overall work validation passed.

## Close Checklist

- [x] Outcome section filled.
- [x] Verification section filled with direct evidence.
- [x] Evidence is inline; no separate evidence file is needed.
- [x] Plan moved to `.hawp/work/closed/2026/10/20/e4240a2f/plan.md`.
- [x] BACKLOG row moved from Active Work to Recently Closed.
- [x] Status report not needed for this bounded index cleanup.
- [x] No separate decision record is needed.
- [x] Staged-path proof not applicable; no commit was requested.

**Status:**

- [x] Plan written.
- [x] Auto-approved (low risk).
- [x] Implemented.
- [x] Verified.
- [x] Closed.
