# Bug / Task: Add clean-as-you-go and structural decomposition guidance for core work

**Backlog ID:** TASK-043
**Type:** task
**Reported:** 2026-05-14
**Risk Level:** low

---

### Input (what was reported)

> add the following coding guideline to always clean up the code and to leave it better than before ... the meaning is to always have a rule to keep the code clean and leave it better than before and to also break down larger code files into smaller ones and organize them into folder like components can have multiple component folders and so on and so forth please add to work items if needed

---

### Context

The repo already tracks workflow and policy changes through HAWP work items. The request asks for an explicit durable rule set: leave touched code cleaner, and split/organize larger files with sensible boundaries.

---

### Analysis

**Root cause (or most likely cause):**
A concise, durable repo-local instruction covering clean-as-you-go and pragmatic decomposition is missing from `.hawp/kit/instructions/`.

**Directly verified:**
- `.hawp/kit/instructions/` currently contains `da-file-tracking.md` only.
- `TASK-041` exists and is focused on commit standards, not code-cleanliness/decomposition policy.

**Inferred (not yet proven):**
- A dedicated instruction file in `.hawp/kit/instructions/` is the right scope to capture this policy without over-expanding global protocol docs.

**Scope - what else is affected:**
- `.hawp/work/BACKLOG.md`
- `.hawp/work/active/TASK-043.md`
- `.hawp/kit/instructions/` (new instruction file)

---

### Work Coordination

**Owner:** agent
**Implementation status:** not-started
**Overlapping files:**

- `.hawp/work/BACKLOG.md`
- `.hawp/kit/instructions/`

**Parallel work risk:** low
**Can implement now:** yes

**Coordination note:**
No active task currently targets this new instruction file path. Existing unrelated working-tree changes are in `shared_standards/public/guidelines/` and will remain untouched.

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
 M shared_standards/public/guidelines/architecture.md
 M shared_standards/public/guidelines/code-style.md
```

---

### Options

#### Option A - Add a dedicated concise instruction file

Create one focused instruction document under `.hawp/kit/instructions/` and keep backlog/plan linkage minimal.

#### Option B - Expand existing global intake instructions

Add the new rules to broad `.github/instructions/*.md` policy files.

---

### Recommended Fix

**Option chosen:** A
**Rationale:**
The request is policy-focused and repo-local. A dedicated file keeps intent explicit, avoids bloating global intake rules, and matches the proposal direction from the provided material.

**Files to change:**

- `.hawp/kit/instructions/clean-code-and-structure.md` - add concise policy guidance
- `.hawp/work/BACKLOG.md` - track task status transitions
- `.hawp/work/active/TASK-043.md` - maintain plan and outcome/verification

**What to verify after:**

- [ ] New instruction file exists and is concise/actionable
- [ ] Guidance explicitly covers leave-better-than-found and reasoned file decomposition
- [ ] Backlog row reflects lifecycle state accurately

---

### Implementation Notes

Keep language practical and bounded: encourage incremental cleanup and decomposition, but forbid broad rewrites without explicit task/approval.

---

## Outcome (filled at close)

Implemented Option A.

- Added `.hawp/kit/instructions/clean-code-and-structure.md` as a concise repo-local instruction covering:
  - leave-code-better incremental cleanup
  - reason-based large-file decomposition
  - ownership-first folder organization
  - approval gate for high-impact structure changes
  - early verification after first structural edit
  - separation of naming cleanup from runtime/schema migration changes
- Tracked this request through a dedicated work item (`TASK-043`) instead of expanding `TASK-041`, because `TASK-041` remains scoped to commit standards.

---

## Verification (filled at close)

- [x] New instruction file exists and is concise/actionable. **Evidence:** file confirmed at `.hawp/kit/instructions/clean-code-and-structure.md`.
- [x] Guidance explicitly covers leave-better-than-found and reasoned decomposition. **Evidence:** sections `Leave It Better Than You Found It` and `Split Large Files For A Reason` in `.hawp/kit/instructions/clean-code-and-structure.md`.
- [x] Backlog row reflects lifecycle through implementation. **Evidence:** `TASK-043` present as `in-progress` in `.hawp/work/BACKLOG.md` prior to close transition.

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
