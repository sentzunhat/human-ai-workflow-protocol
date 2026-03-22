# Task: Decouple backlog-upgrade scripts from CLI-first architecture

**Backlog ID:** TASK-049
**Type:** task
**Reported:** 2026-05-14
**Risk Level:** medium

---

### Input (what was reported)

> the librarian shoud not be started or scafolded for cli usage please keep the scripts simples and without cli in mind also if that is a new work item please update it and then do the implementation of clear scripts that are extensible for future use on the cli if the script seems to belong in the cli add a file called CLI.md where it says how will the command will be

---

### Context

`librarian/scripts/backlog-upgrade/cli.ts` currently mixes argument parsing, command execution, process exit handling, and output file writes. This makes script reuse harder when consumed as a library/API and makes future CLI extension less clear.

---

### Analysis

**Root cause (or most likely cause):**
The current architecture is CLI-first (`runCLI` drives everything and calls `process.exit`), while the reusable execution logic is not exposed as a clean script API.

**Directly verified:**
- `librarian/scripts/backlog-upgrade/cli.ts` contains parse + execute + process exit behavior.
- `librarian/scripts/backlog-upgrade/index.ts` invokes `runCLI` directly.
- `librarian/scripts/backlog-upgrade/__tests__/cli.test.ts` only validates argument parsing and does not cover script-first execution.

**Inferred (not yet proven):**
Separating execution logic into script modules and keeping CLI as a thin adapter will improve maintainability and future command expansion.

**Scope — what else is affected:**
- `librarian/scripts/backlog-upgrade/cli.ts`
- `librarian/scripts/backlog-upgrade/index.ts`
- `librarian/scripts/backlog-upgrade/__tests__/cli.test.ts`
- new docs artifact for command contract (`CLI.md`) if CLI-facing behavior remains

---

### Work Coordination

**Owner:** agent
**Implementation status:** done
**Overlapping files:**

- `librarian/scripts/backlog-upgrade/cli.ts`
- `librarian/scripts/backlog-upgrade/index.ts`
- `librarian/scripts/backlog-upgrade/__tests__/cli.test.ts`

**Parallel work risk:** medium
**Can implement now:** yes

**Coordination note:**
Overlap exists with TASK-028 in the same backlog-upgrade lane, but this request is a direct user-prioritized architecture adjustment and can be implemented in the same lane safely.

**Path discipline:**

- Use exact repo-relative paths from repository root for all files in this plan.
- Basenames alone are unsafe for path-sensitive work unless the file is truly at repository root and explicitly marked as such.

**Repo-root proof (required before path-sensitive edits):**

```bash
$ pwd
<repo-root-abs>/human-ai-workflow-protocol

$ git rev-parse --show-toplevel
<repo-root-abs>/human-ai-workflow-protocol

$ git rev-parse --show-prefix


$ git status --short
 M .github/copilot-instructions.md
 M .hawp/work/BACKLOG.md
 D .hawp/work/active/TASK-026.md
 D .hawp/work/active/TASK-027.md
 M .hawp/work/active/TASK-028.md
 D .hawp/work/active/TASK-030.md
 D .hawp/work/active/TASK-031.md
 D .hawp/work/active/TASK-032.md
 D .hawp/work/active/TASK-037.md
 D .hawp/work/active/TASK-038.md
 D .hawp/work/active/TASK-040.md
 D .hawp/work/active/TASK-041.md
 D .hawp/work/closed/2026/05/12/0007.md
 D .hawp/work/closed/2026/05/12/0008-CLARIFICATION-exact-paths.md
 D .hawp/work/closed/2026/05/12/0008-install-update-distribution-review.md
 D .hawp/work/closed/2026/05/12/HAWP-BACKLOG-VALIDATE-PLAN.md
 M .hawp/work/closed/2026/05/12/TASK-030-files.md
 D .hawp/work/closed/2026/05/13/0007.md
 D .hawp/work/closed/2026/05/13/HAWP-BACKLOG-VALIDATE-PLAN.md
 M .hawp/work/closed/2026/05/13/TASK-030-files.md
 M .hawp/work/closed/2026/05/13/TASK-033.md
 M .hawp/work/closed/2026/05/13/TASK-038.md
 M .hawp/work/closed/2026/05/14/TASK-044.md
 M librarian/scripts/backlog-upgrade/cli.ts
 M librarian/scripts/backlog-upgrade/models/backlog-fix-plan.ts
?? .hawp/work/closed/2026/05/14/TASK-026.md
?? .hawp/work/closed/2026/05/14/TASK-027.md
?? .hawp/work/closed/2026/05/14/TASK-030.md
?? .hawp/work/closed/2026/05/14/TASK-031.md
?? .hawp/work/closed/2026/05/14/TASK-032.md
?? .hawp/work/closed/2026/05/14/TASK-037.md
?? .hawp/work/closed/2026/05/14/TASK-040.md
?? .hawp/work/closed/2026/05/14/TASK-047.md
?? .hawp/work/closed/2026/05/14/TASK-048.md
?? .hawp/work/evidence/2026/05/14/
?? .hawp/work/notes/2026/05/12/
?? .hawp/work/notes/2026/05/13/
?? .hawp/work/status/2026/05/14/TASK-047-status.md
?? .hawp/work/status/2026/05/14/TASK-048-status.md
?? librarian/scripts/backlog-upgrade/__tests__/detection.test.ts
?? librarian/scripts/backlog-upgrade/detection/
?? librarian/scripts/backlog-upgrade/output/
```

---

### Options

#### Option A — Keep `cli.ts` and extract script core

Create a script-first execution module (no process exit, no argv assumptions), then keep `cli.ts` only as parse/adapter glue.

Trade-offs: minimal churn, clear boundary, easiest migration path.

#### Option B — Remove `cli.ts` and make `index.ts` do everything

Push parsing and execution into one file while keeping an exported helper API.

Trade-offs: less files, but blurs adapter vs logic boundary and is harder to extend cleanly.

---

### Recommended Fix

**Option chosen:** A
**Rationale:**
It cleanly separates reusable script execution from CLI concerns and keeps a small adapter layer for future command evolution.

**Files to change:**

- `librarian/scripts/backlog-upgrade/script.ts` — new script-first API for backlog-upgrade execution
- `librarian/scripts/backlog-upgrade/cli.ts` — reduce to argument parsing + adapter (no core logic)
- `librarian/scripts/backlog-upgrade/index.ts` — keep node entrypoint, delegate to adapter
- `librarian/scripts/backlog-upgrade/__tests__/cli.test.ts` — keep parser tests aligned to new adapter exports
- `librarian/scripts/backlog-upgrade/__tests__/script.test.ts` — add script-first behavior tests
- `librarian/scripts/backlog-upgrade/CLI.md` — CLI contract and future command shape

**What to verify after:**

- [ ] Typecheck passes for backlog-upgrade modules
- [ ] Existing CLI parse tests still pass
- [ ] New script tests pass
- [ ] Dry-run execution still prints report and supports output/export writes
- [ ] CLI adapter still returns expected help/version and exit codes

---

### Implementation Notes

- Script API should return structured results (status code + outputs), not exit the process.
- CLI adapter should be the only layer that uses `process.exit` and raw argv semantics.
- Keep behavior backward-compatible for `backlog upgrade` arguments.

---

## Outcome (filled at close)

- Added script-first execution module at `librarian/scripts/backlog-upgrade/script.ts` with reusable API:
  - `runBacklogUpgradeScript(options)` returns structured result data (`exitCode`, output text, stderr lines, notices)
  - dry-run detection/report execution no longer depends on process-level CLI control flow
- Refactored `librarian/scripts/backlog-upgrade/cli.ts` into adapter responsibilities:
  - argument parsing and command-shape validation
  - help/version text generation
  - delegation to `runBacklogUpgradeScript`
  - return exit codes instead of forcing process termination inside parser flow
- Updated executable boundary in `librarian/scripts/backlog-upgrade/index.ts`:
  - keeps process-level behavior limited to setting `process.exitCode`
- Added script behavior tests at `librarian/scripts/backlog-upgrade/__tests__/script.test.ts`:
  - dry-run output behavior
  - `--output` and `--export-plan` file write behavior
  - apply-mode usage error behavior in this slice
- Added CLI contract doc at `librarian/scripts/backlog-upgrade/CLI.md` with command shape and extension pattern.

---

## Verification (filled at close)

- [x] TypeScript compilation passes for `librarian`. **Evidence:** `npm --prefix librarian run typecheck` (pass).
- [x] Existing parser/detector tests still pass. **Evidence:** `cd librarian && npx --yes tsx --test scripts/backlog-upgrade/__tests__/cli.test.ts scripts/backlog-upgrade/__tests__/detection.test.ts` (pass).
- [x] New script-level tests pass. **Evidence:** `cd librarian && npx --yes tsx --test scripts/backlog-upgrade/__tests__/script.test.ts` (pass).
- [x] Combined focused suite passes. **Evidence:** `cd librarian && npx --yes tsx --test scripts/backlog-upgrade/__tests__/cli.test.ts scripts/backlog-upgrade/__tests__/detection.test.ts scripts/backlog-upgrade/__tests__/script.test.ts` (9/9 pass).

---

## Close Checklist

- [x] Outcome section filled
- [x] Verification section filled (all claims have direct evidence or "unproven" tag)
- [x] All evidence files referenced exist in `../evidence/YYYY/MM/DD/` or are noted as inline
- [x] Plan file will be moved to `../closed/YYYY/MM/DD/TASK-049.md`
- [x] BACKLOG.md row moved from Active to Recently Closed (or marked done)
- [x] Status report written (optional)
- [x] Decision file created if applicable (not applicable for this task)
- [x] Staged-path proof captured before commit:
  - [x] `git diff --name-status`
  - [x] `git diff --check`
  - [x] `git diff --cached --name-status`
  - [x] `git diff --cached --check`
  - [x] `git status --short`
- [x] If basename-only paths appeared during path-sensitive work, mark report unsafe and restart after correction

**Status:**

- [x] Plan written
- [x] Approved / auto-approved (explicit user request to implement)
- [x] Implemented
- [x] Verified
- [x] Closed
