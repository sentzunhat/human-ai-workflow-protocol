# Work Intake - Plan Template

---

## Bug / Task: Design HAWP fix-up and upgrade command model

**Backlog ID:** TASK-021
**Type:** task
**Reported:** 2026-05-11
**Risk Level:** low

---

### Input (what was reported)

> Create a design note for future HAWP librarian fix-up and upgrade commands. Define command model, safety, dry-run/apply behavior, boundaries for mechanical vs AI-drafted edits, and recommend a V1 implementation slice.

---

### Context

The validator already detects workflow drift for local and external `.hawp` roots. The next stage needs safe, reviewable fix-up/upgrade command behavior that preserves context and keeps files as source of truth.

---

### Analysis

**Root cause (or most likely cause):**
There is no explicit command design note yet for automated or semi-automated fix-up/upgrade flows. Without it, future implementation may drift on safety and evidence discipline.

**Directly verified:**
- `.hawp/work/BACKLOG.md` exists and currently tracks active items.
- Request provides explicit command candidates, goals, and constraints.

**Inferred (not yet proven):**
- Future CLI command implementation will need shared patch planning primitives derived from validator findings.

**Scope - what else is affected:**
- New design note content in `.hawp/work/notes/`.
- Backlog and this plan file for intake tracking.

---

### Work Coordination

**Owner:** agent
**Implementation status:** done
**Overlapping files:**

- `.hawp/work/BACKLOG.md`
- `.hawp/work/active/TASK-021.md`

**Parallel work risk:** low
**Can implement now:** yes

**Coordination note:**
No overlap with active code implementation files is expected; only work-tracking and notes files are touched.

---

### Options

#### Option A - Add one focused design note in work notes

Write a single note that defines command surfaces, safety gates, AI boundaries, and example outputs. Fast to review and directly actionable for V1 planning.

#### Option B - Add command-level mini-specs in separate files

Write one file per command family (`validate`, `fix-up`, `upgrade`, `backlog`). More modular but heavier to maintain before first implementation.

---

### Recommended Fix

**Option chosen:** A
**Rationale:**
Produces a compact, implementation-oriented baseline with clear safety constraints and allows later splitting by command as needed.

**Files to change:**

- `.hawp/work/notes/2026/05/11/hawp-librarian-fixup-upgrade-command-model.md` - add the new design note.
- `.hawp/work/BACKLOG.md` - maintain intake status transitions and close tracking.
- `.hawp/work/active/TASK-021.md` - complete outcome/verification at close.

**What to verify after:**

- [ ] Design note includes all requested output sections.
- [ ] Safety model explicitly encodes dry-run-first and human review before apply.
- [ ] AI drafting constraints state evidence-only and unproven verification marking.

---

### Implementation Notes

Keep the note implementation-facing (commands, inputs, outputs, refusal cases, and guardrails) and avoid prescribing runtime database dependencies.

---

## Outcome (filled at close)

Created a new design note at `.hawp/work/notes/2026/05/11/hawp-librarian-fixup-upgrade-command-model.md` that defines:

- proposed command names and scope
- a layered safety model and guardrails
- explicit `--dry-run` and `--apply` behavior
- boundaries for mechanical fix-up vs AI-drafted sections
- refusal/no-auto-fix cases
- an example dry-run report
- a pragmatic V1 implementation slice

Also updated backlog lifecycle for TASK-021 through intake states and prepared close artifacts.

---

## Verification (filled at close)

- [x] Design note includes all requested output sections. **Evidence:** file confirmed at `.hawp/work/notes/2026/05/11/hawp-librarian-fixup-upgrade-command-model.md`.
- [x] Safety model encodes dry-run-first and human review before apply. **Evidence:** "Safety Model" and "Dry-run and Apply Behavior" sections in `.hawp/work/notes/2026/05/11/hawp-librarian-fixup-upgrade-command-model.md`.
- [x] AI drafting constraints enforce evidence-only inputs and unproven verification marking. **Evidence:** "What Can Be AI-Drafted" section in `.hawp/work/notes/2026/05/11/hawp-librarian-fixup-upgrade-command-model.md`.
- [x] Backlog transition and closure link updated. **Evidence:** `.hawp/work/BACKLOG.md` includes TASK-021 in Recently Closed after close step.

---

## Close Checklist

**Before marking this task done in BACKLOG.md, verify ALL:**

- [x] Outcome section filled (what was actually implemented)
- [x] Verification section filled (all checks listed, each with direct evidence or "unproven" tag)
- [x] All evidence files referenced exist in `../evidence/YYYY/MM/DD/` or are noted as inline
- [x] Plan file moved to `../closed/YYYY/MM/DD/TASK-021.md`
- [x] BACKLOG.md row moved from Active to Recently Closed (or marked done)
- [x] Status report written (optional: only if non-trivial OR if something remains unproven OR if a decision/pattern emerged)
- [x] Decision file created if applicable (not applicable for this task)

**Status:**

- [x] Plan written
- [x] Approved / auto-approved (low risk)
- [x] Implemented
- [x] Verified
- [x] Closed
