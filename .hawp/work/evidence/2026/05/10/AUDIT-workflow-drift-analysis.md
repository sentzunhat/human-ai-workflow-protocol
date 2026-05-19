# HAWP Workflow Audit — Drift & Simplification Analysis

**Date:** 2026-05-10
**Scope:** Files-as-source-of-truth operating model (kit/, work/, instructions)
**Mission:** Find drift, collision, and duplication points; simplify task flow for future librarian indexing

---

## Executive Summary

The workflow is well-structured and largely sound, but **lacks atomic close patterns** and **evidence discipline** at scale. Task flow is clear conceptually but **handoff points slip** because:

1. **Plan files and status reports duplicate effort** without clear rules for which is source.
2. **Evidence storage is diffuse** (in plan files, status reports, separate evidence/, or notes/).
3. **Active → Closed transitions are multi-step with no checksum** — items can get lost between backlog update, plan file move, evidence linking.
4. **Backlog compaction policy exists but lacks enforcement** — "create a task first" is advisory, not automated.
5. **Decisions live in a separate tree** from work items, making traceability weak.
6. **Parked items have no re-triage rule** — they can age without review.
7. **Overlapping-file detection** relies on manual plan-file notation with no validation gate.

These are not urgent bugs but are the exact **friction points that make indexing and search unsafe** without careful schema design.

---

## Part 1: Current Drift & Slip Points

### 1.1 — Plan File vs. Status Report Duplication

**Observed Drift:**

Plan files (in `work/active/<ID>.md`) contain:

- Input, Context, Analysis
- Root cause, verification, scope
- Work Coordination (owner, status, overlap check)
- Options and Recommended Fix
- Outcome (added at close)
- Verification section
- Final Status checklist

Status reports (in `work/status/YYYY/MM/DD/`) also capture:

- Intent, Current State, What Changed
- What Was Verified, What Remains Unproven
- Constraints, Help Wanted, Suggested Next Step

**The Problem:**

- When closing a work item, should you update the plan file's "Outcome" section OR write a separate status report?
- Current rule is vague: "save a status report if non-trivial or if context matters."
- An agent can't safely determine which is the source of truth for "what was verified at close."
- A librarian would see both, not know which one is current, and have to reconcile manually.

**Directly Verified:**

- TASK-008 closed with plan-file "Outcome" section but no status report in `work/status/2026/05/02/`.
- TASK-007 closed with plan-file "Verification" checklist but no separate status report.
- TASK-006 closed with plan-file "Outcome" section and reference to separate evidence, but the evidence file is not present in `work/evidence/2026/05/02/`.

**Evidence for Inference:**

- Kit guidance says status reports are optional companions, not required structure.
- No rule exists for "when a plan file alone is sufficient vs. when you must also write a status report."

---

### 1.2 — Evidence Storage is Diffuse

**Observed Drift:**

Evidence can live in five places:

1. Directly in plan file (Verification section, links to screenshots/outputs)
2. In `work/evidence/YYYY/MM/DD/` files
3. In `work/status/YYYY/MM/DD/` files
4. In `work/notes/YYYY/MM/DD/` files
5. Not documented anywhere (referenced but not stored)

**The Problem:**

- TASK-006 plan file says "Confirmed new files exist" but no linked evidence files.
- No rule enforces "if you claim verification, you must store evidence here."
- No rule enforces "evidence goes in plan file OR in evidence/ folder, not both."
- When a librarian indexes and finds a verification claim, it won't know whether to trust it or look for supporting evidence.

**Directly Verified:**

- `work/evidence/` folder exists but is nearly empty (only dated subfolders, no files in them).
- Most plan files store Verification as a checklist within the plan file itself.
- No referenced evidence file paths in closed work.

---

### 1.3 — Active → Closed Transition is Multi-Step with No Atomic Close

**Observed Drift:**

Closing a task requires:

1. Update plan file with Outcome and Verification
2. Move plan file from `work/active/<ID>.md` to `work/closed/YYYY/MM/DD/<ID>.md`
3. Update BACKLOG.md: remove from Active, update Recently Closed row
4. (Optional) Write status report to `work/status/YYYY/MM/DD/`
5. (Optional) Write evidence links or files to `work/evidence/YYYY/MM/DD/`

**The Problem:**

- This is 3–5 operations. Any can slip.
- No checksum or manifest to verify all steps completed.
- A plan file could be moved to closed/ but Backlog not updated. A new agent wouldn't know the item is closed.
- Recently Closed can diverge from actual closed/ folder contents.
- No step validates "all evidence links in the plan file exist in evidence/ folder."

**Directly Verified:**

- BACKLOG.md Recently Closed table stops at 10 items.
- `work/closed/2026/05/` contains 3 date folders (01, 02, 03).
- No manifest file tracks which items should be in Recently Closed.

---

### 1.4 — Backlog Compaction Policy Lacks Enforcement

**Observed Drift:**

The policy rule:

> "If BACKLOG.md already has many Done rows, create or recommend a work item titled 'Compact BACKLOG.md and archive closed work.' before adding more Done rows."

**The Problem:**

- What is "many"? 5? 20? 100?
- Who decides when to trigger? The agent? A human?
- It's advisory, not enforced. If an agent misses the signal, backlog grows silently.
- No threshold check; no automation.

**Directly Verified:**

- TASK-008 was the most recent compaction task (closed 2026-05-02).
- It successfully capped Recently Closed at 10 items.
- No warning system exists to flag when another compaction is needed.

---

### 1.5 — Overlapping-File Collision Detection is Manual

**Observed Drift:**

Plan files have a "Parallel work risk" table:

```
Overlapping files:
- path/to/file

Parallel work risk: low | medium | high
Can implement now: yes | no | only after approval
```

**The Problem:**

- This only works if **all** active items have written plan files and filled in their overlap table.
- If item A has a plan but item B is still in `inbox`, they won't see each other.
- If two agents work in parallel, one might finish its plan before the other has even started.
- The backlog rule "check backlog before starting" is procedural, not enforced.
- No automated check; purely human discipline.

**Directly Verified:**

- Currently, BACKLOG.md shows "No active work" so no collision risk.
- But TASK-008 noted "no overlapping items" in its coordination section because it was doing cleanup.
- A concurrent new task would have had to manually check BACKLOG.md to find existing work.

---

### 1.6 — Decisions Live in a Separate Tree

**Observed Drift:**

ADRs and project decisions go to `work/decisions/YYYY/MM/DD/` and are **not tracked in BACKLOG.md**. There is no backlog ID that says "decision made about X; see decisions/2026/05/02/Y.md."

**The Problem:**

- A work item can reference a decision, but the decision is not cross-indexed in the backlog.
- If you're indexing, you have to search both trees: backlog for work items, decisions/ for decisions.
- Hard to reconstruct "which decision led to which work item?"
- A librarian can't build a dependency graph without duplicating effort.

**Directly Verified:**

- `work/decisions/` exists with dated subfolders but only a README file in 2026/.
- TASK-006 mentions creating an ADR (`adr-root-hawp-work-for-source-repo.md`) but the decision file itself is not present in `work/decisions/2026/05/02/`.
- No backlog row links to the decision; there's only a plan-file reference.

---

### 1.7 — Parked Items Lack Re-Triage Rule

**Observed Drift:**

Work can be moved to `work/parked/` with backlog status `parked`. No rule specifies when or how to re-evaluate parked items.

**The Problem:**

- Parked items could age indefinitely without review.
- No expiration or trigger rule: "if parked > N days, re-triage" or "review parked every sprint."
- Cognitive load: every agent has to remember which items are parked and why.

**Directly Verified:**

- `work/parked/` folder exists but is empty (no current parked items).
- BACKLOG.md mentions a "Blocked / Parked" section but it's empty.
- No examples of what makes an item parked vs. blocked vs. deferred.

---

### 1.8 — Work Item Type Classification is Loose

**Observed Drift:**

Backlog IDs use type prefixes: `TASK-`, `BUG-`, `improvement-`.

**The Problem:**

- No clear rule for which to use.
- "Task" vs. "Improvement" distinction is not defined.
- Could lead to inconsistent classification that makes indexing harder.
- Numbering sequences per type are not enforced (TASK-001, TASK-002… vs. TASK-008, BUG-005).

**Directly Verified:**

- BACKLOG.md uses both lowercase `task` and lowercase `bug` in the Type column.
- TASK-008 is labeled `task`, TASK-007 is labeled `improvement`, BUG-005 is labeled `bug`.
- No schema rule prevents new classifications.

---

### 1.9 — Kit Patterns and Reviews Are Not Linked to Work

**Observed Drift:**

`kit/patterns/`, `kit/examples/`, and `kit/reviews/` exist but are never referenced from backlog items or work files.

**The Problem:**

- A work item that creates a reusable pattern doesn't mark itself as pattern-creating.
- Conversely, when starting new work, patterns are not automatically surfaced.
- A librarian would see work items but wouldn't know which ones created reusable wisdom.

**Directly Verified:**

- `kit/reviews/` and `kit/patterns/` directories exist but are not explored.
- Plan files don't have a "Generates Reusable Pattern?" field or link.

---

## Part 2: Recommended Simplified Task Flow

The current flow (`intake → analyze → plan → review gate → implement → verify → close`) is sound. Simplify by making these three changes:

### 2.1 — Unified Close Checklist (Atomic Handoff)

Replace the ad-hoc "5 optional steps" with a single **mandatory close checklist**:

```markdown
## Close Checklist

Before marking done in BACKLOG.md, verify ALL of:

- [ ] Plan file in work/active/<ID>.md has:
  - [ ] Root cause documented
  - [ ] Outcome section filled (what was implemented)
  - [ ] Verification section filled (all checks listed and checked)
  - [ ] Each verification claim has direct evidence (see Evidence Discipline below)
- [ ] Backlog row updated:
  - [ ] Moved from Active/Blocked/Parked to Recently Closed
  - [ ] Updated date set to close date
  - [ ] Status set to `done`
- [ ] Plan file moved:
  - [ ] From work/active/<ID>.md to work/closed/YYYY/MM/DD/<ID>.md
- [ ] Evidence linked (if non-trivial):
  - [ ] All claims in Verification have supporting files in work/evidence/YYYY/MM/DD/
  - [ ] File names are descriptive (e.g., TASK-010-verification-screenshots.md)
- [ ] Status report written (if item was non-trivial OR context matters):
  - [ ] File saved to work/status/YYYY/MM/DD/<ID>-status.md
  - [ ] Intent, Current State, What Changed, What Was Verified, What Remains Unproven captured
  - [ ] References plan file for full detail
- [ ] Decision captured (if item resolves a question):
  - [ ] ADR file in work/decisions/YYYY/MM/DD/<ID>-decision.md (if significant)
  - [ ] Link added to plan file's "Recommended Fix" section
```

**Why:** Moves from procedural guidance to structured checklist. Agents can verify completeness. A librarian can check "is this close complete?" by looking for checklist marks.

---

### 2.2 — Evidence Discipline Rule (Single Source)

**Rule:**

- All verification claims in plan files MUST list supporting evidence.
- Evidence goes ONLY in `work/evidence/YYYY/MM/DD/` folder or in the plan file if <50 words.
- Large evidence (screenshots, command output, logs) → separate file in evidence/.
- Small evidence (one-line confirmations, status codes) → inline in plan file.
- Status reports REFERENCE plan file evidence; they do not duplicate it.

**Format for plan files:**

```markdown
### Verification

- [x] New instruction file exists: **File:** `.github/instructions/hawp-backlog-alignment.instructions.md` **Evidence:** reviewed file structure in root and core overlays
- [x] New prompt file exists: **File:** `.github/prompts/hawp-backlog-alignment.prompt.md` **Evidence:** see work/evidence/2026/05/02/TASK-007-file-audit.md
- [ ] Build passes — NOT YET VERIFIED (requires live environment; see status report)
```

**Why:** Makes evidence explicit and traceable. A librarian can collect all claims, find their evidence, and flag missing evidence.

---

### 2.3 — Simplified Backlog + Status Report Integration

**Rule:**

- BACKLOG.md is ONLY the active index. No Done rows stay there long.
- Recently Closed is capped at last 10 items OR last 14 days (pick one, apply consistently).
- When Recently Closed is at cap: create a `COMPACT-<date>.md` note in `work/notes/` that records which items were archived and when. Reference this from BACKLOG.md as history.
- Status reports are companions to plan files for context transfer, not replacements.

**When to write a status report:**

- **Always write if:** item reveals a pattern, decision, or lesson.
- **Always write if:** something was unproven or requires live environment.
- **Optional if:** item is trivial (e.g., typo fix, one-line config).

**Format for status report:**

```markdown
# Status Report: TASK-010

**Related Plan:** work/closed/2026/05/03/TASK-010.md

## Intent

Verify legacy .work folder cleanup was complete.

## Current State

Root .work/ folder was removed; contents migrated to .hawp/work/.

## What Changed

[Brief summary; see plan file Outcome section for full detail]

## What Was Directly Verified

- Confirmed all dated subfolders moved to work/closed/YYYY/MM/DD/
- Confirmed BACKLOG.md Active Work is now in .hawp/work/

## What Remains Unproven

- Whether any external references to root .work/ exist outside this repo
```

**Why:** Status reports become context-transfer artifacts, not primary evidence. Backlog stays lean. Librarian can read plan file for truth, status report for context.

---

## Part 3: Backlog Alignment Rules (Librarian-Safe)

### 3.1 — Compact Backlog Rules

1. **BACKLOG.md contains only:**
   - Purpose
   - Status key
   - Active Work table (working items)
   - Blocked / Parked table (non-active but tracked)
   - Recently Closed table (last 10 items OR last 14 days)
   - Archive index (links to closed/, decisions/, evidence/ folders by date)

2. **BACKLOG.md must never contain:**
   - Every completed item ever
   - Copied plan-file content
   - Full implementation notes (link to plan file instead)
   - Historical changelog or update log

3. **Recently Closed cap rule:**
   - When adding a new row would exceed the cap, create a compaction note and remove oldest rows.
   - Compaction notes live in `work/notes/YYYY/MM/DD/COMPACT-<date>.md` and simply record which items were archived.
   - Link from BACKLOG.md archive section to compaction notes so history is discoverable.

### 3.2 — Active Work Rows

Each active row:

- ID (e.g., TASK-010, BUG-003)
- Type (task | bug | improvement | decision)
- Title
- Status (inbox | analyzing | plan-ready | approved | in-progress | parked | blocked)
- Owner (human | agent | unassigned)
- Plan File (link to work/active/<ID>.md)
- Updated (date of last update)

---

## Part 4: Active → Closed Archival Rules

### 4.1 — What Must Happen at Close

**Prerequisite:** Close checklist (see Part 2.1) must be 100% complete.

**Steps (in order):**

1. **Update plan file** (`work/active/<ID>.md`):
   - Add Outcome section with what was implemented
   - Complete Verification section with all checks and direct evidence
   - Add final Status checklist (copy from close checklist, all marked ✓)

2. **Move plan file:**
   - From `work/active/<ID>.md` to `work/closed/YYYY/MM/DD/<ID>.md`
   - Never delete; always preserve in closed/

3. **Link evidence (if exists):**
   - Create `work/evidence/YYYY/MM/DD/<ID>-*.md` files for large evidence
   - Reference files by name in plan file Verification section

4. **Write status report (if applies):**
   - Save to `work/status/YYYY/MM/DD/<ID>-status.md`
   - Reference plan file, don't duplicate analysis

5. **Update BACKLOG.md:**
   - Remove row from Active Work or Blocked/Parked
   - Add short row to Recently Closed (ID, Type, Title, Closed date, Link to plan)
   - If Recently Closed at cap, trigger compaction

6. **Record decision (if applicable):**
   - If item resolves a design question, create `work/decisions/YYYY/MM/DD/<ID>-decision.md`
   - Link from plan file Recommended Fix section

### 4.2 — Archive Index (Librarian Discovery)

Keep a `work/ARCHIVE.md` file that lists:

```markdown
# Archive Index

## By Date

- 2026-05-03: TASK-010 (cleanup)
- 2026-05-02: TASK-008 (compaction), TASK-007 (feature), TASK-006 (refactor)

## By Type

### Bugs

- BUG-005: ...
- BUG-004: ...

### Tasks

- TASK-010: ...
- TASK-008: ...

### Improvements

- TASK-007: ...
- TASK-006: ...

### Decisions

- ADR-001 (2026-05-02): ...
```

**Why:** Gives a librarian a single entry point to discover what's archived and where.

---

## Part 5: Prerequisites Before Librarian/Indexing Work Starts

### Must Be True

1. ✅ **Kit is stable and documented.**
   - Current state: Kit is well-documented (spec.md, authoring-patterns.md, usage guides).
   - No changes needed; kit is the source of protocol truth.

2. ⚠️ **Work files follow a consistent schema.**
   - Current state: Plan files follow template, but closed plan files vary slightly (no checksum of completeness).
   - **Action needed:** Enforce close checklist so every closed item has same structure.

3. ⚠️ **BACKLOG.md is compact and trustworthy.**
   - Current state: Recently Closed is capped at 10; good. But no archive index or compaction history.
   - **Action needed:** Add ARCHIVE.md index and compaction note trail.

4. ⚠️ **Evidence is traceable.**
   - Current state: Most evidence lives in plan files; very few separate evidence/ files exist.
   - **Action needed:** Enforce evidence discipline so each claim can be verified.

5. ⚠️ **Decisions are discoverable from work items.**
   - Current state: Decisions exist in separate tree; not referenced from backlog or plan files.
   - **Action needed:** Add decision-to-work-item linking in plan files and BACKLOG.md.

6. ⚠️ **Overlapping-file detection is checked at plan-write time.**
   - Current state: Manual table in plan files; no validation.
   - **Action needed:** Create a quick reference script that checks active plan files for overlaps before auto-approving low-risk items.

7. ✅ **Parked items have clear scope.**
   - Current state: Parked folder exists; current backlog has no parked items.
   - **No change needed yet, but:** Add a re-triage rule (e.g., "parked items auto-moved to blocked after 30 days for re-evaluation").

### Implementation Priority

**Before building librarian:**

1. Document and enforce close checklist in intake-workflow.md
2. Add ARCHIVE.md index to work/
3. Implement evidence discipline rule in intake-plan.md template
4. Create a simple validation script that checks:
   - All rows in Recently Closed have plan files in closed/ folders
   - All claims in plan files have supporting evidence or explicit "unproven" tag
   - Overlapping files in plan files don't collide across active items
5. Add decision-linking rule to plan template

---

## Part 6: What Must Be True for Librarian to Safely Index

### Librarian Safety Invariants

When the librarian reads a plan file, it must be able to safely assume:

1. **Closed items are complete.**
   - Every plan file in `work/closed/` has passed close checklist.
   - Outcome, Verification, and Status sections are populated.
   - No plan file in closed/ is missing evidence links.

2. **Claims are traceable.**
   - Every "verified" claim has evidence.
   - Every "unproven" claim is explicitly marked.
   - Evidence files are named predictably (e.g., `<ID>-<claim-type>.md`).

3. **Backlog is consistent with reality.**
   - Every active item in BACKLOG.md has a file in `work/active/`.
   - Every closed item in BACKLOG.md has a file in `work/closed/YYYY/MM/DD/`.
   - No plan file exists in active/ without a backlog row.
   - No plan file is duplicated (one per ID, ever).

4. **Decisions are linked.**
   - If a plan file references a decision, the decision file exists and is linkable.
   - Conversely, if a decision file creates work, that work has a backlog row.

5. **Overlaps are documented.**
   - Every active plan file lists overlapping files in Work Coordination section.
   - Librarian can detect potential collisions by comparing overlap tables.

6. **History is preserved.**
   - No item is deleted; only archived.
   - Archived items are discoverable via ARCHIVE.md and dated folders.
   - Compaction notes record what was archived and when.

---

## Part 7: Roadmap Summary (Not Librarian Design, Just Enablement)

### Phase 1: Enforce Close Completeness (1–2 days)

- [ ] Update intake-workflow.md with close checklist
- [ ] Add close checklist section to intake-plan.md template
- [ ] Apply close checklist retroactively to last 5 closed items (in closed/2026/05/)
- [ ] Document evidence naming convention

### Phase 2: Archive Index & History (1 day)

- [ ] Create work/ARCHIVE.md with date and type indices
- [ ] Create work/notes/2026/05/10/COMPACT-2026-05-02.md (retroactive record of last compaction)
- [ ] Update BACKLOG.md to link to ARCHIVE.md

### Phase 3: Decision Linking (1 day)

- [ ] Update intake-plan.md template to include "Decision" row in plan files
- [ ] Retroactively link TASK-006 to its ADR (once ADR file is verified)
- [ ] Add work/decisions/2026/05/02/ links to BACKLOG.md

### Phase 4: Evidence Discipline (1–2 days)

- [ ] Update intake-plan.md to show evidence tagging format
- [ ] Review last 5 closed items and create evidence files where missing
- [ ] Document which claims are proven vs. unproven for each

### Phase 5: Validation Script (optional, 1 day)

- [ ] Write simple shell or Python script that:
  - Checks active/ items have BACKLOG.md rows
  - Checks closed/ items have BACKLOG.md rows (in archive or recently-closed)
  - Flags missing evidence files
  - Detects ID collisions or missing links

---

## Conclusion

**Current State:** Well-intentioned, mostly working, but lacks atomic close and evidence discipline at scale.

**Risk if unchanged:**

- Librarian indexing will find incomplete, inconsistent, or orphaned items.
- Will require custom reconciliation logic to handle edge cases.
- Will make it unsafe to trust searchable results without manual verification.

**Fix complexity:** Low. Mostly documentation, checklists, and naming conventions. No schema changes needed.

**Timeline to librarian-ready:** ~1 week of incremental enforcement and backward-compatibility fixes.

---

**Audit completed:** 2026-05-10
