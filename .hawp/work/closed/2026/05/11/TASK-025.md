# Work Intake — Plan

## Bug / Task: Regenerate branch-specific install/update distribution files

**Backlog ID:** TASK-025
**Type:** task
**Reported:** 2026-05-11
**Risk Level:** low

---

### Input (what was reported)

> it did not install it like copy paste and install update it just did something different,
>
> can you please update the distribution folders files to be able to have a fully generated update and install branches files ready for distribution please

---

### Context

The distribution guides under `core/distribution/generated/` must be copy/paste-ready for install/update on both `main` and `dev` branches and must be regenerated from source fragments + authoritative scripts.

### Analysis

**Root cause (or most likely cause):**
Prior work appears to have focused on a docs-only path outside this repository context, so the user did not receive a repo-local regeneration pass from the distribution pipeline in this workspace.

**Directly verified:**
- `core/distribution/sources/` and `core/distribution/generated/` exist in this repository.
- Build script exists at `librarian/scripts/distribution/build/index.ts` and composes the four generated outputs from source fragments plus extracted script blocks.
- NPM scripts include `distribution:build` and `distribution:validate`.

**Inferred (not yet proven):**
- Generated files may still need a fresh rebuild/validation pass to guarantee branch-ready copy/paste artifacts after recent content updates.

**Scope — what else is affected:**
- `core/distribution/generated/install-main.md`
- `core/distribution/generated/install-dev.md`
- `core/distribution/generated/update-main.md`
- `core/distribution/generated/update-dev.md`
- Potentially source fragments if build/validation reveals drift.

### Work Coordination

**Owner:** agent
**Implementation status:** not-started
**Overlapping files:**

- `.hawp/work/BACKLOG.md`
- `.hawp/work/active/TASK-025.md`
- `core/distribution/generated/*.md`

**Parallel work risk:** low
**Can implement now:** yes

**Coordination note:**
No active backlog rows currently overlap distribution sources/generated outputs.

### Options

#### Option A — Regenerate from existing source pipeline

Run `distribution:build` and `distribution:validate`; commit resulting generated deltas only. Fastest and lowest risk.

#### Option B — Manually edit generated files

Patch generated outputs by hand. Higher drift risk from source-of-truth and likely to regress on next build.

### Recommended Fix

**Option chosen:** A
**Rationale:**
The repo already has a deterministic composition pipeline. Regenerating from source guarantees branch-specific copy/paste script sections stay aligned with authoritative install/update scripts.

**Files to change:**

- `core/distribution/generated/*.md` — rebuilt output if drift is detected.
- `.hawp/work/BACKLOG.md` — status transitions and closure entry.
- `.hawp/work/active/TASK-025.md` — outcome + verification at close.

**What to verify after:**

- [ ] Build completes and reports generated file state.
- [ ] Validation passes for all generated files.
- [ ] `install-dev/main` and `update-dev/main` outputs include correct branch-specific `REF` values.

### Implementation Notes

If validation fails, apply minimal source-fragment corrections and rebuild before closure.

---

## Outcome (filled at close)

- Executed the repo-local distribution pipeline to regenerate all branch-specific install/update guides from source fragments and authoritative scripts.
- `npm run distribution:build` reported all generated outputs as unchanged, confirming they are already in sync with current source content.
- `npm run distribution:validate` passed, confirming no composition drift.
- Confirmed generated guides include correct branch-pinned `REF` assignments (`main` vs `dev`) in copy/paste command blocks.

---

## Verification (filled at close)

**Direct evidence required for each claim.** Evidence can be:

- Inline if <50 words
- Linked to `../evidence/YYYY/MM/DD/<ID>-<claim>.md` if larger
- Marked explicitly if unproven: "NOT YET VERIFIED — requires live environment"

- [x] Generated distribution outputs are current. **Evidence:** `npm run distribution:build` output showed `unchanged` for all four files and `0/4 file(s) updated`.
- [x] Distribution validation passes. **Evidence:** `npm run distribution:validate` output reported `distribution validation passed: generated outputs are current`.
- [x] Branch pinning is correct for all generated guides. **Evidence:** checked `REF="main"`/`REF="dev"` matches in `core/distribution/generated/install-main.md`, `core/distribution/generated/install-dev.md`, `core/distribution/generated/update-main.md`, `core/distribution/generated/update-dev.md`.

---

## Close Checklist

- [x] Outcome section filled
- [x] Verification section filled (all checks listed, each with direct evidence or "unproven" tag)
- [x] All evidence files referenced exist in `../evidence/YYYY/MM/DD/` or are noted as inline
- [x] Plan file will be moved to `../closed/YYYY/MM/DD/TASK-025.md`
- [x] BACKLOG.md row moved from Active to Recently Closed (or marked done)
- [ ] Status report written (optional: non-trivial/unproven/decision-bearing)
- [ ] Decision file created if applicable

**Status:**

- [x] Plan written
- [x] Approved / auto-approved (low risk)
- [x] Implemented
- [x] Verified
- [x] Closed
