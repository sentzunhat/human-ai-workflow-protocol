## Task: Promote public standards and audit private-shareable workflow guidance

**Backlog ID:** TASK-069
**Type:** standards-update
**Reported:** 2026-06-04
**Risk Level:** low

---

### Input (what was reported)

> USING THE standards and docs folders please update the core hawp folder with new public standards and review the private if they can be shared no secret sauce about internal operations just workflows and review and audit file and make it a new work items if needed and update the core .hawp folder then we will do an update of hawp on dev after these changes and we get teh new standards and docs publics

---

### Context

The request is to absorb new public-safe standards from shared source folders into `core/.hawp`, review private lane content for workflow-only shareability, and capture audit evidence plus follow-up work items where adaptation is needed.

### Analysis

**Root cause (or most likely cause):**
The source standards tree contains newer public docs content that is not fully represented in `core/.hawp/kit/standards/**`, while private lane files need a fresh shareability review to avoid leaking internal implementation details.

**Directly verified:**

- Source public standards contain `public/standards/docs/hawp-install-update-safety.md` and `public/standards/docs/README.md`.
- `core/.hawp/kit/standards/**` currently has no `docs/` category.
- Existing `core/.hawp/kit/standards/guidelines/testing.md` and `core/.hawp/kit/standards/guidelines/documentation.md` already align with source copies.
- Private files reviewed include auth, providers, product, internal-runtime, and security lanes.

**Inferred (not yet proven):**

- The docs safety standard can be absorbed directly as public-safe operational guidance.
- Private files contain some reusable process patterns, but should be routed to adaptation work rather than direct absorption.

**Scope - what else is affected:**

- `core/.hawp/kit/standards/README.md`
- `core/.hawp/kit/standards/docs/`
- `.hawp/work/BACKLOG.md`
- `.hawp/work/evidence/2026/06/04/`

---

### Work Coordination

**Owner:** agent
**Implementation status:** in-progress
**Overlapping files:**

- `core/.hawp/kit/standards/README.md`
- `core/.hawp/kit/standards/docs/README.md`
- `core/.hawp/kit/standards/docs/hawp-install-update-safety.md`
- `.hawp/work/BACKLOG.md`
- `.hawp/work/evidence/2026/06/04/TASK-069-standards-public-private-audit.md`

**Parallel work risk:** low
**Can implement now:** yes

**Coordination note:**
No active overlap detected in backlog; this item can implement directly with bounded docs-only changes.

**Path discipline:**

- Use exact repo-relative paths from repository root for all files in this plan.
- Basenames alone are unsafe for path-sensitive work unless the file is truly at repository root and explicitly marked as such.

**Repo-root proof (required before path-sensitive edits):**

```bash
pwd
git rev-parse --show-toplevel
git rev-parse --show-prefix
git status --short
```

Recorded output (machine-local prefix redacted):

```text
<repo-root-abs>
<repo-root-abs>

```

---

### Options

#### Option A - Direct absorb of new public docs standards + private audit artifact

Add missing public docs standards to `core/.hawp/kit/standards/docs/`, update standards index, and write an audit artifact classifying private files as non-absorb or adapt-required.

#### Option B - Full re-sync of all public standards folders

Re-copy all public standards categories into core and then re-review all private lanes.

---

### Recommended Fix

**Option chosen:** A
**Rationale:**
This is the smallest safe change set that satisfies the request while preserving prior absorbed standards and avoiding unnecessary churn.

**Files to change:**

- `core/.hawp/kit/standards/README.md` - add docs category and canonical map entries
- `core/.hawp/kit/standards/docs/README.md` - add category index
- `core/.hawp/kit/standards/docs/hawp-install-update-safety.md` - add public-safe standard
- `.hawp/work/evidence/2026/06/04/TASK-069-standards-public-private-audit.md` - add audit and classification results
- `.hawp/work/BACKLOG.md` - add follow-up item if adaptation work is needed

**What to verify after:**

- [ ] New docs standard files exist in core standards tree.
- [ ] Standards index references new docs category and file.
- [ ] Private lane audit evidence is captured with direct vs adaptation classification.
- [ ] Backlog includes a follow-up work item for private workflow adaptation.

---

### Implementation Notes

Keep private-lane recommendations process-focused and avoid copying route inventories, token internals, provider strategy, or internal runtime instructions.

---

## Outcome (filled at close)

- Added a new public standards category at `core/.hawp/kit/standards/docs/`.
- Added `core/.hawp/kit/standards/docs/hawp-install-update-safety.md` with reusable install/update safety and boundary guidance.
- Added `core/.hawp/kit/standards/docs/README.md` as the category index.
- Updated `core/.hawp/kit/standards/README.md` to include docs category entries in both section index and canonical map.
- Created audit evidence at `.hawp/work/evidence/2026/06/04/TASK-069-standards-public-private-audit.md`.
- Opened follow-up work item `TASK-070` for private workflow-only adaptation review.

## Verification (filled at close)


### Evidence Follow-Up

- [ ] Research evidence for: New docs standards files were added under `core/.hawp/kit/standards/docs/`.
- [ ] Update the original verification checklist line with Evidence: ... or explicit unproven wording.
- [ ] Research evidence for: Core standards index references the new docs category and canonical standard.
- [ ] Update the original verification checklist line with Evidence: ... or explicit unproven wording.
- [ ] Research evidence for: Private/public review classification is documented with direct absorb vs adapt-required outputs.
- [ ] Update the original verification checklist line with Evidence: ... or explicit unproven wording.
- [ ] Research evidence for: Follow-up work item created for private workflow adaptation.
- [ ] Update the original verification checklist line with Evidence: ... or explicit unproven wording.

- [x] New docs standards files were added under `core/.hawp/kit/standards/docs/`. — unproven: evidence not recorded at close (annotated 2026-07-20)
  - **Evidence:** `core/.hawp/kit/standards/docs/README.md`, `core/.hawp/kit/standards/docs/hawp-install-update-safety.md`.
- [x] Core standards index references the new docs category and canonical standard. — unproven: evidence not recorded at close (annotated 2026-07-20)
  - **Evidence:** `core/.hawp/kit/standards/README.md`.
- [x] Private/public review classification is documented with direct absorb vs adapt-required outputs. — unproven: evidence not recorded at close (annotated 2026-07-20)
  - **Evidence:** `.hawp/work/evidence/2026/06/04/TASK-069-standards-public-private-audit.md`.
- [x] Follow-up work item created for private workflow adaptation. — unproven: evidence not recorded at close (annotated 2026-07-20)
  - **Evidence:** `.hawp/work/BACKLOG.md` contains `TASK-070`.

## Close Checklist

- [x] Outcome section filled
- [x] Verification section filled (all claims have direct evidence or "unproven" tag)
- [x] Evidence files created if large/complex
- [x] Plan file moved to closed/YYYY/MM/DD/
- [x] BACKLOG.md row moved from Active to Recently Closed (or marked done)
- [ ] Status report written (optional)
- [ ] Decision file created if applicable
- [ ] Staged-path proof captured before commit

**Status:**

- [x] Plan written
- [x] Approved / auto-approved (low risk)
- [x] Implemented
- [x] Verified
- [x] Closed
