# Bug / Task: Make clean-as-you-go and incremental decomposition always-on in Copilot

**Backlog ID:** TASK-045
**Type:** task
**Reported:** 2026-05-14
**Risk Level:** low

---

### Input (what was reported)

> please do that automatically applied by copilot
>
> not only documented but to always improve the code quality with breaking down large files small chunks at a time in concept that aligns what we have made here in this convo

---

### Context

A repo-local instruction was added under `.hawp/kit/instructions/`, but the request now asks to make this behavior always applied by Copilot in normal coding flow.

---

### Analysis

**Root cause (or most likely cause):**
The global Copilot overlay does not yet include an explicit always-on clean-as-you-go + incremental decomposition rule.

**Directly verified:**

- `.github/copilot-instructions.md` exists and is the primary integration overlay.
- `.github/instructions/intake.instructions.md` is global (`applyTo: "**"`) and can carry enforceable intake rules.

**Inferred (not yet proven):**

- Adding concise policy language to the always-loaded overlay and intake ambient rules will make agent behavior consistent with the requested approach.

**Scope - what else is affected:**

- `.github/copilot-instructions.md`
- `.github/instructions/intake.instructions.md`
- `.hawp/work/BACKLOG.md`
- `.hawp/work/active/TASK-045.md`

---

### Work Coordination

**Owner:** agent
**Implementation status:** not-started
**Overlapping files:**

- `.github/copilot-instructions.md`
- `.github/instructions/intake.instructions.md`
- `.hawp/work/BACKLOG.md`

**Parallel work risk:** low
**Can implement now:** yes

**Coordination note:**
This is a narrow policy-overlay update. No code-runtime files are touched.

**Path discipline:**

- Use exact repo-relative paths from repository root for all files in this plan.
- Basenames alone are unsafe for path-sensitive work unless the file is truly at repository root and explicitly marked as such.

**Repo-root proof (required before path-sensitive edits):**

```bash
pwd
<repo-root-abs>/human-ai-workflow-protocol

git rev-parse --show-toplevel
<repo-root-abs>/human-ai-workflow-protocol

git rev-parse --show-prefix


git status --short
```

---

### Options

#### Option A - Update global overlays only

Add concise always-on rules to `.github/copilot-instructions.md` and a matching operational rule in `.github/instructions/intake.instructions.md`.

#### Option B - Add another kit-only instruction

Keep rules in `.hawp/kit/instructions/` only and rely on manual interpretation.

---

### Recommended Fix

**Option chosen:** A
**Rationale:**
The user asked for automatic Copilot application. Global overlays are the direct control surface.

**Files to change:**

- `.github/copilot-instructions.md` - add always-on clean-code and decomposition baseline
- `.github/instructions/intake.instructions.md` - add operational rule for incremental decomposition within scope
- `.hawp/work/BACKLOG.md` - track status transitions

**What to verify after:**

- [ ] Overlay includes explicit always-on clean-as-you-go rule
- [ ] Intake rules allow bounded decomposition while preserving scope discipline
- [ ] Backlog/plan lifecycle reflects implementation and closure

---

### Implementation Notes

Keep wording strict enough to drive behavior but bounded to avoid broad refactor drift.

---

## Outcome (filled at close)

Implemented Option A by updating both always-on instruction surfaces.

- Added a new `Code Quality Baseline (Always On)` section in `.github/copilot-instructions.md`.
- Added Operating Rule 26 in `.github/instructions/intake.instructions.md` to enforce bounded clean-as-you-go and incremental decomposition during intake-driven work.
- Preserved scope guardrails to avoid broad unsanctioned refactors.

---

## Verification (filled at close)

- [x] Overlay includes explicit always-on clean-as-you-go rule. **Evidence:** `.github/copilot-instructions.md` contains `## Code Quality Baseline (Always On)` with four behavior rules.
- [x] Intake rules allow bounded decomposition while preserving scope discipline. **Evidence:** `.github/instructions/intake.instructions.md` contains `26. Apply clean-as-you-go by default.` plus bounded decomposition language.
- [x] Task lifecycle tracked through implementation. **Evidence:** `TASK-045` row present as `in-progress` in `.hawp/work/BACKLOG.md` prior to close transition.

---

## Close Checklist

**Before marking this task done in BACKLOG.md, verify ALL:**

- [x] Outcome section filled (what was actually implemented)
- [x] Verification section filled (all checks listed, each with direct evidence or "unproven" tag)
- [x] All evidence files referenced exist in `../evidence/YYYY/MM/DD/` or are noted as inline
- [x] Plan file will be moved to `../closed/YYYY/MM/DD/<ID>.md`
- [x] BACKLOG.md row moved from Active to Recently Closed (or marked done)
- [ ] Status report written (optional: only if non-trivial OR if something remains unproven OR if a decision/pattern emerged)
- [ ] Decision file created if applicable (only if this task resolves a design question)
- [x] Staged-path proof captured before commit:
  - [x] `git diff --name-status`
  - [x] `git diff --check`
  - [x] `git diff --cached --name-status`
  - [x] `git diff --cached --check`
  - [x] `git status --short`
- [ ] If basename-only paths appeared during path-sensitive work, mark report unsafe and restart after correction (fresh path-locked pass after two failures)

**Status:**

- [x] Plan written
- [x] Approved / auto-approved (low risk)
- [x] Implemented
- [x] Verified
- [x] Closed
