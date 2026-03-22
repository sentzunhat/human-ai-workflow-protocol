## Improvement: Automate closed-record section scaffolding and date-folder reconciliation

**Backlog ID:** TASK-068
**Type:** improvement
**Reported:** 2026-05-17
**Risk Level:** medium
**Status:** done

---

### Input (what was reported)

Confirm whether automation can add missing close sections/templates and detect/move stray closed records into correct date-based folders to keep `.hawp/work/` organized by rule.

---

### Context

The backlog-upgrade command already supports dry-run detection with CLI adapter boundaries. Apply mode and write operations were not implemented when this item started.

---

### Analysis

**Root cause (or most likely cause):**
Automation stopped at reporting. It did not apply deterministic fixes for missing sections or misplaced closed files.

**Directly verified:**

- `librarian/scripts/backlog-upgrade/script.ts` previously returned usage error for non-dry-run mode.
- `librarian/scripts/backlog-upgrade/CLI.md` previously documented apply mode as future slice.

**Inferred (not yet proven):**

- Applying safe operations from TASK-067 classifications would automate most cleanup while preserving historical content.

**Scope — what else is affected:**

- `librarian/scripts/backlog-upgrade/script.ts`
- `librarian/scripts/backlog-upgrade/cli.ts`
- `librarian/scripts/backlog-upgrade/CLI.md`
- `librarian/scripts/backlog-upgrade/__tests__/**`
- `.hawp/work/closed/2026/05/17/TASK-066.md`

---

### Work Coordination

**Owner:** agent
**Implementation status:** done
**Overlapping files:**

- `librarian/scripts/backlog-upgrade/**`
- `.hawp/work/closed/**` (when apply mode is exercised)

**Parallel work risk:** medium
**Can implement now:** yes

**Coordination note:**
This executed after TASK-067 produced the exact action matrix, then fed TASK-066 manual fallback completion.

---

### Recommended Fix

1. Implement apply-mode operation pipeline in script layer with deterministic action classes.
2. Add safeguards:
   - dry-run default,
   - dirty-tree checks with opt-out flag,
   - post-apply validator summary.
3. Implement templating/scaffolding for missing close sections without inventing evidence content.
4. Implement date-folder reconciliation for explicitly classifiable stray files.
5. Add tests for apply behavior, idempotence, and no-op safety.

Implemented slice:

- Closed-record scaffolding added for Outcome, Verification, and Close Checklist.
- Filename-based Backlog ID inference added for inferable legacy closed files.
- Apply-mode dirty-tree guard added.
- Date-folder reconciliation added for date-prefixed filename mismatch cases.
- CLI help and contract updated to match implemented behavior.

---

## Outcome (filled at close)

- Implemented apply-mode execution in `librarian/scripts/backlog-upgrade/script.ts`.
- Added safeguards:
  - clean working tree required in apply mode unless `--force-dirty`,
  - optional post-run validator summary with existing workflow validator orchestration,
  - ambiguous legacy files without inferable ID skipped with notice.
- Implemented closed-record normalization actions:
  - add missing Backlog ID from inferable filename IDs,
  - scaffold missing `Outcome`, `Verification`, and `Close Checklist` sections,
  - reconcile date-prefixed closed filenames into matching date folders when mismatch exists.
- Updated CLI help/contract docs to reflect implemented apply behavior.
- Executed apply mode against repo closed records; 20 records normalized and downstream TASK-066 fallback completed.

## Verification (filled at close)

- [x] Focused script tests pass: **Evidence:** `node --import tsx --test scripts/backlog-upgrade/__tests__/script.test.ts` passed including apply-mode dirty-tree guard, scaffolding, and date-folder reconciliation tests.
- [x] TypeScript checks pass: **Evidence:** `npm run -s typecheck` in `librarian/` completed with no errors.
- [x] Workflow remains consistent after apply run: **Evidence:** `npm run -s validate:workflow` reports `Result: VALIDATION PASS` with `Warnings: 0`.
- [x] Apply run performed on audited set: **Evidence:** `npx tsx scripts/backlog-upgrade/index.ts backlog upgrade --apply --validate --force-dirty` output reported `20 closed-record file(s) normalized.`

## Close Checklist

- [x] Outcome section filled
- [x] Verification section filled (all claims have direct evidence or "unproven" tag)
- [x] Evidence files created if large/complex
- [x] Plan file moved to `closed/YYYY/MM/DD/`
- [x] BACKLOG.md updated
- [x] Status report written (if non-trivial / unproven / decision-bearing)
- [x] Decision file created (if applicable)

**Status:**

- [x] Plan written
- [x] Approved / auto-approved (low risk)
- [x] Implemented
- [x] Verified
- [x] Closed
