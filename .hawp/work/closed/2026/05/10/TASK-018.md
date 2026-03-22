# Work Intake — Plan

## Bug / Task: Surface legacy untyped closed work in validator output

**Backlog ID:** TASK-018
**Type:** improvement
**Reported:** 2026-05-10
**Risk Level:** low

---

### Input (what was reported)

> TASK-016 looks good and appears backwards compatible.
>
> One follow-up check before we consider the validator stable:
>
> Please verify how the validator classifies legacy closed files with no TASK-/BUG-style ID, such as:
>
> 2026-04-26-hawp-adr-template-review.md
>
> Expected behavior:
>
> - Do not require modern close sections if it is before 2026-05-10.
> - Do not silently hide it as an archive/supporting file unless it clearly matches support-file patterns like status, summary, checkpoint, archive, or BACKLOG.
> - Prefer reporting it as INFO or WARN legacy untyped work, so old work remains visible to the user/librarian.
>
> Do not rewrite or rename the file.
> Do not start UUID migration.
> Do not add indexing/search/SQLite.
>
> make it backwards companitble for all the folders in a generic way please if possible

---

### Context

`checkClosedTaskCompleteness` currently treats any file without an extractable `TYPE-NNN` ID as supporting and skips it. That hides legacy untyped work like `2026-04-26-hawp-adr-template-review.md` from report output.

---

### Analysis

**Root cause (or most likely cause):**
Classification conflates two cases into one skip path: `(1) explicit supporting file` and `(2) untyped legacy work`.

**Directly verified:**

- File exists at `closed/2026/04/26/2026-04-26-hawp-adr-template-review.md`.
- Validator output currently reports only `1 supporting file skipped` and does not list this file in WARN/FAIL.

**Inferred (not yet proven):**
Adding a separate `untypedLegacy` classification and reporting it as WARN will preserve visibility and avoid false strict failures.

**Scope — what else is affected:**

- `types.ts` result model for completeness
- `closed-task-completeness.ts` classification logic
- `reporter.ts` output sections

---

### Work Coordination

**Owner:** agent
**Implementation status:** in-progress
**Overlapping files:**

- `librarian/scripts/validate-hawp-workflow/validations/closed-task-completeness.ts`
- `librarian/scripts/validate-hawp-workflow/types.ts`
- `librarian/scripts/validate-hawp-workflow/reporter.ts`

**Parallel work risk:** low
**Can implement now:** yes

**Coordination note:**
No active HAWP item currently targets these files.

---

### Options

#### Option A — Treat untyped files as supporting

Keep current behavior. Legacy untyped files remain hidden.

#### Option B — Add explicit untyped legacy classification and report

Keep strict checks for modern typed plan files, keep explicit supporting skips, and report untyped legacy closed files as WARN/INFO-visible items.

---

### Recommended Fix

**Option chosen:** B
**Rationale:**
Meets visibility requirement while preserving backward compatibility and strictness for modern work.

**Files to change:**

- `librarian/scripts/validate-hawp-workflow/types.ts` — add untyped legacy reporting fields
- `librarian/scripts/validate-hawp-workflow/validations/closed-task-completeness.ts` — split classification paths
- `librarian/scripts/validate-hawp-workflow/reporter.ts` — print untyped legacy section

**What to verify after:**

- [ ] `2026-04-26-hawp-adr-template-review.md` appears in a visible WARN/INFO section
- [ ] Explicit supporting files still counted as skipped
- [ ] `npm run typecheck` passes
- [ ] `npm run validate:workflow` passes

---

### Implementation Notes

Prefer generic rules:

- supporting skip only for clear patterns (`backlog`, `status-report`, `summary`, `checkpoint`, `archive`)
- untyped files before cutoff are tolerated but visible
- untyped files on/after cutoff should be FAIL to keep strict current/future policy

---

## Outcome (filled at close)

Implemented classifier split and visibility improvements for legacy untyped closed work:

1. Added explicit classification buckets in `closed-task-completeness.ts`:
	- `supporting` (skip only on explicit support patterns)
	- `legacy-untyped` (before cutoff, visible WARN)
	- `current-untyped` (on/after cutoff, FAIL)
	- `plan` (typed files requiring modern close sections)
2. Replaced rigid closed-folder traversal with recursive generic traversal under `closed/`.
3. Added date extraction from any `/YYYY/MM/DD/` path segment.
4. Updated result model (`types.ts`) to track:
	- `supportingSkipped[]`
	- `untypedLegacy[]`
	- `untypedCurrent[]`
5. Updated reporter output to show:
	- `[FAIL] Untyped closed files (modern)`
	- `[WARN] Legacy untyped closed files`
	- `[INFO] Supporting files skipped by pattern`

Result: `2026-04-26-hawp-adr-template-review.md` is now visible as legacy WARN instead of being silently skipped.

---

## Verification (filled at close)

- [x] Legacy untyped file is visible in output: **Evidence:** `../evidence/2026/05/10/TASK-018-legacy-untyped-classification.md`
- [x] File `2026-04-26-hawp-adr-template-review.md` no longer hidden as skipped supporting file: **Evidence:** `../evidence/2026/05/10/TASK-018-legacy-untyped-classification.md`
- [x] `npm run typecheck` passes: **Evidence:** inline command output (tsc --noEmit succeeded)
- [x] `npm run validate:workflow` passes with warnings only: **Evidence:** `../evidence/2026/05/10/TASK-018-legacy-untyped-classification.md`

---

## Close Checklist

- [x] Outcome section filled (what was actually implemented)
- [x] Verification section filled (all checks listed, each with direct evidence or "unproven" tag)
- [x] All evidence files referenced exist in `../evidence/YYYY/MM/DD/` or are noted as inline
- [x] Plan file will be moved to `../closed/YYYY/MM/DD/<ID>.md`
- [x] BACKLOG.md row moved from Active to Recently Closed (or marked done)
- [ ] Status report written (optional: only if non-trivial OR if something remains unproven OR if a decision/pattern emerged)
- [ ] Decision file created if applicable (only if this task resolves a design question)

**Status:**

- [x] Plan written
- [x] Approved / auto-approved (low risk)
- [x] Implemented
- [x] Verified
- [x] Closed
