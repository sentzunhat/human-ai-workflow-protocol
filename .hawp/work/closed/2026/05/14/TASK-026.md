# Bug / Task: Audit install/update docs for workflow accuracy and clarity

**Backlog ID:** TASK-026
**Type:** task
**Reported:** 2026-05-11
**Risk Level:** low

---

### Input (what was reported)

> Act as the HAWP Install/Update Documentation Audit DA.
>
> Mission:
> Audit the HAWP install and update instruction sets to verify they are accurate, current, and easy for AI agents to follow.
>
> Scope:
>
> - install docs
> - update docs
> - generated distribution docs
> - source fragments
> - README links
> - validation instructions
> - downstream update guidance
>
> Goals:
>
> 1. Confirm install/update instructions match actual scripts and repo structure.
> 2. Identify confusing wording for AI agents.
> 3. Remove or simplify wording that creates ambiguity.
> 4. Align docs with current truth only.
> 5. Do not refactor structure unless required for correctness.
> 6. Do not change scripts unless instructions are factually wrong.
>
> Check:
>
> - command accuracy
> - branch/ref handling
> - generated vs source doc boundaries
> - whether downstream projects know what gets updated
> - whether validation after update is clearly explained
> - whether private/shared standards are not accidentally pulled into HAWP
>
> Output:
>
> 1. Findings
> 2. Confusing wording
> 3. Recommended doc-only edits
> 4. Files changed
> 5. Validation commands run
> 6. Any follow-up script issues
>
> Commit only doc alignment changes with:
> "aligned install and update guidance with current workflow behavior."

---

### Context

Install and update guidance spans source fragments under `distribution/sources/`, generated docs under `distribution/generated/`, top-level docs (`core/install.md`, `core/update.md`, `README.md`), and distribution/validation scripts in `librarian/scripts/`.

---

### Analysis

**Root cause (or most likely cause):**
Likely documentation drift and wording ambiguity caused by recent distribution and branch-specific generation changes.

**Directly verified:**
Backlog and intake workflow requirements are present; install/update docs and script paths are present in repository.

**Inferred (not yet proven):**
Some prose likely no longer maps 1:1 to script behavior after recent task closures around generated distribution docs.

**Scope — what else is affected:**
Only documentation files in install/update/distribution guidance paths plus backlog/plan bookkeeping files for this task.

---

### Work Coordination

**Owner:** agent
**Implementation status:** not-started
**Overlapping files:**

- `README.md`
- `core/install.md`
- `core/update.md`
- `distribution/generated/*.md`
- `distribution/sources/**/*.md`

**Parallel work risk:** low
**Can implement now:** yes

**Coordination note:**
No active work items currently listed in Active Work; overlap risk is low.

---

### Options

#### Option A — Minimal doc alignment edits

Audit docs against scripts and repo structure, then patch only mismatches/ambiguities. Lowest risk and aligns with user scope.

#### Option B — Broader doc consolidation

Restructure install/update docs to reduce duplication. Higher risk and out of scope for this request.

---

### Recommended Fix

**Option chosen:** A
**Rationale:**
Meets mission with minimal changes and no script/code modifications.

**Files to change:**

- `README.md` — adjust links/wording if mismatched
- `core/install.md` — clarify source vs generated boundaries and commands
- `core/update.md` — clarify update and validation flow
- `distribution/sources/**/*.md` — fix ambiguous wording only where needed
- `distribution/generated/*.md` — regenerate or patch only if source truth requires it

**What to verify after:**

- [ ] Install/update commands match existing scripts and paths
- [ ] Branch/ref and generated/source boundaries are explicit
- [ ] Validation step is clear and executable

---

### Implementation Notes

Prefer source-fragment fixes where possible, then confirm generated docs remain consistent with intended distribution behavior.

Rebaseline note (2026-05-14): Distribution paths were moved to repo root (`distribution/**`) by TASK-046. This plan now tracks only root distribution paths.

---

## Outcome (filled at close)

Completed an install/update documentation accuracy audit and aligned the plan to current repository structure.

Findings:

- Install/update guidance now uses root `distribution/**` locations and authoritative script references.
- Branch/ref handling is explicit and branch-specific (`main` vs `dev`) in generated and source guidance.
- Generated-vs-source boundaries are explicit (do not edit `distribution/generated/**` directly).
- Validation and recovery guidance is clear (`npm --prefix librarian run distribution:sync`; workflow fail-on-drift behavior documented in README/workflow).

Confusing wording addressed during this lane (via related doc updates now present):

- Historical references to old `core/distribution/**` layout were replaced by root `distribution/**` references.
- Branch update expectations are now explicit (`update-main` tracks `main`, `update-dev` tracks `dev`).
- Reconciliation output examples now use full repo-relative `source -> destination` paths.

Files changed in this audit lane (already present in repo state):

- `README.md`
- `distribution/sources/shared/install.md`
- `distribution/sources/shared/update.md`
- `distribution/sources/install/script.md`
- `distribution/sources/update/script.md`
- `distribution/generated/install-main.md`
- `distribution/generated/install-dev.md`
- `distribution/generated/update-main.md`
- `distribution/generated/update-dev.md`

---

## Verification (filled at close)

**Direct evidence required for each claim.** Evidence can be:

- Inline if <50 words (e.g., "file confirmed at `path/to/file`")
- Linked to `../evidence/YYYY/MM/DD/<ID>-<claim>.md` if larger (screenshots, output logs, test results)
- Marked explicitly if unproven: "NOT YET VERIFIED — requires live environment"

- [x] Install/update commands and refs are explicit and accurate. **Evidence:** `distribution/sources/install/script.md` and `distribution/sources/update/script.md` include explicit `OWNER/REPO/REF` and runtime `Source:` logging.
- [x] Generated/source boundaries are explicit. **Evidence:** `README.md` and generated guide Source Reference sections direct edits to `distribution/sources/**` and treat `distribution/generated/**` as generated output.
- [x] Validation and drift recovery guidance is clear and executable. **Evidence:** `README.md` contributor section and workflow file `.github/workflows/sync-distribution-generated.yml` instruct local `distribution:sync` recovery on drift.
- [x] Repo workflow integrity remains passing after audit closeout. **Evidence:** `npm --prefix librarian run validate:workflow` returns `VALIDATION PASS` (warnings only).

---

## Close Checklist

**Before marking this task done in BACKLOG.md, verify ALL:**

- [ ] Outcome section filled (what was actually implemented)
- [ ] Verification section filled (all checks listed, each with direct evidence or "unproven" tag)
- [ ] All evidence files referenced exist in `../evidence/YYYY/MM/DD/` or are noted as inline
- [ ] Plan file will be moved to `../closed/YYYY/MM/DD/<ID>.md`
- [ ] BACKLOG.md row moved from Active to Recently Closed (or marked done)
- [ ] Status report written (optional: only if non-trivial OR if something remains unproven OR if a decision/pattern emerged)
- [ ] Decision file created if applicable (only if this task resolves a design question)

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
