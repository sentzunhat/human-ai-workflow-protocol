# Consolidating copied agent definitions into the current agent set

**Backlog ID (Legacy):** — (UUID-native item)
**UUID:** `81d6a35c-8192-4235-9421-5f9384cad92a`
**Type:** improvement
**Reported:** 2026-10-20

---

## Input (verbatim)

> hey thank you so much now with that done can you please upgrade and merge or combine the agents copy to the agtents we have and the clean up the copy folder pls

## Intake Summary

Compare the copied agent definitions with the active `.github/agents/` set, retain the more current repository-specific versions, and remove the redundant copy directory.

## Current Context

The copied folder contains four files: the Docs Alignment Auditor agent and schema, plus the Master Builder agent and schema. The active folder contains corresponding agent/schema files and a dedicated HAWP agent.

## Initial Analysis

**Directly verified:**

- The copied Docs Alignment Auditor describes generic `src/`, `docs/`, Node/package paths and standards that do not match this Go/HAWP repository. The active version instead names `librarian/src/`, `.hawp/kit/`, `core/.hawp/kit/`, and `distribution/sources/`.
- The copied Master Builder has an older explicit tool list without the active version's HAWP MCP tool group. The active version includes the repository-specific identity and HAWP workflow.
- The copied and active Master Builder JSON schemas have identical contents.
- The active `.github/agents/` contains `hawp.agent.md`, which has no counterpart in the copy.
- Repository search found no consumer of `.github/agents copy/` outside the copied files' schema references.
- HAWP search returned no matching active item; the checked active backlog rows contain no agent-copy consolidation task.

**Inferred (not yet proven):**

- The copied definitions are an older backup rather than a source of upgrades; copying them into the active folder would regress repository-specific paths and tool access.

**Likely scope:**

- Keep all existing files in `.github/agents/` unchanged.
- Remove only the four confirmed files in `.github/agents copy/`.
- Update this work item and `.hawp/work/BACKLOG.md`.

## Risk + Review Gate

**Risk:** low
**Gate:** auto-implement on low

## Backlog + Plan Link

**Status now:** closed
**Plan file:** work/closed/2026/10/20/81d6a35c/plan.md

## Next Step

- [x] Investigation recorded above (required before planning)
- [x] Write or update the plan file
- [x] Move backlog status accordingly

## Options

### Option A — Keep the active agents and remove the stale copy

Retain the current repository-specific agents and schemas unchanged, then delete only the four files in `.github/agents copy/`.

### Option B — Merge the copied files into the active set

Rejected: the copied docs auditor and schema use unrelated generic/Node.js paths, the copied Master Builder has an older tool list without HAWP access, and its schema duplicates the active schema.

## Recommended Fix

**Option chosen:** A

**Rationale:** The active definitions already contain the repository-specific guidance and tools. Merging the copy would risk replacing current behavior with stale instructions. No content needs to be carried over.

**Files to change:**

- `.github/agents copy/docs-alignment.agent.md` — remove stale duplicate.
- `.github/agents copy/docs-alignment.schema.json` — remove stale duplicate.
- `.github/agents copy/master-builder.agent.md` — remove stale duplicate.
- `.github/agents copy/master-builder.schema.json` — remove duplicate schema.
- `.hawp/work/BACKLOG.md` — update status and close this work item.
- `.hawp/work/active/81d6a35c/plan.md` — record the plan, outcome, and verification.

**What to verify after:**

- [x] All five current files under `.github/agents/`, including `hawp.agent.md`, remain present.
- [x] `.github/agents copy/` and its four duplicate files are gone.
- [x] HAWP validation and `git diff --check` pass.

## Work Coordination

**Owner:** maintainer
**Implementation status:** done
**Overlapping files:**

- `.github/agents/`
- `.github/agents copy/`
- `.hawp/work/BACKLOG.md`
- `.hawp/work/active/81d6a35c/plan.md`

**Parallel work risk:** low
**Can implement now:** yes

**Coordination note:** No active backlog item or external consumer was found for the copy folder. Existing working-tree changes are unrelated and will remain untouched.

**Path discipline:** All paths in this plan are exact repo-relative paths.

**Repo-root proof (captured before agent cleanup):**

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
?? ".github/agents copy/"
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

Delete the four known duplicate files individually. Do not overwrite or rewrite any active agent. Preserve all unrelated working-tree changes.

## Outcome

Kept the current repository-specific Docs Alignment Auditor, Master Builder, their schemas, and the HAWP agent unchanged. The copied definitions had no unique content worth merging and were older or identical to the active files. Deleted the four duplicate files and removed the now-empty `.github/agents copy/` directory.

## Verification

- [x] The copy directory is absent; `test ! -e '.github/agents copy'` passed.
- [x] All five current files under `.github/agents/` remain present.
- [x] `git diff --check -- .hawp/work/BACKLOG.md` passed.
- [x] Final `hawp-hawp_work_validate`: kit validation passed (3 checks, 0 issues); work validation passed (0 issues, 1 warning); local links passed (150 Markdown files); overall `hawp check` passed all 3 validations.
- [x] Detailed `cd librarian/src && go run ./cmd/hawp work validate`: backlog consistency, closed-task completeness (255/255), evidence integrity (9/9), and dead links passed; verification clarity warned (309/390 evidence-backed, 71 explicitly unproven, 10 ambiguous). Overall work validation passed.

## Close Checklist

- [x] Outcome section filled.
- [x] Verification section filled with direct evidence.
- [x] Evidence is inline; no separate evidence file is needed.
- [x] Plan moved to `.hawp/work/closed/2026/10/20/81d6a35c/plan.md`.
- [x] BACKLOG row moved from Active Work to Recently Closed.
- [x] Status report not needed for this bounded cleanup.
- [x] No separate decision record is needed.
- [ ] Staged-path proof not applicable; no commit was requested.

**Status:**

- [x] Plan written.
- [x] Auto-approved (low risk).
- [x] Implemented.
- [x] Verified.
- [x] Closed.
