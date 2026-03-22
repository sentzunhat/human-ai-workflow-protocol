# Status Report: HAWP Workflow Foundation (Phase 1)

**Date:** 2026-05-10
**Intent:** Establish clean, simple, indexable task-closure patterns without designing the librarian

## Summary

Updated HAWP documentation and templates to enforce atomic task closure and evidence discipline. Task files are now the single source of truth; backlog is compact active index; evidence is traceable; status reports are optional companions.

---

## Files Changed

### 1. `.hawp/kit/templates/intake-plan.md`

**Changes:**

- Added `Outcome` section (filled at close) — documents what was actually implemented
- Added `Verification` section (filled at close) — evidence discipline with inline/linked evidence format
- Added `Close Checklist` section — 7-point mandatory close gate before marking done in BACKLOG.md
- Updated Status checklist to align with close phase

**Outcome:** Plan files now have explicit structure for close-time documentation and evidence linkage.

### 2. `.hawp/kit/usage/intake-workflow.md`

**Changes:**

- Enhanced Step 6 (Verify) with evidence discipline guidance
- Rewrote Step 7 (Close) as structured 6-step handoff with prerequisites
- Added "Close Checklist Quick Reference" for easy copy-to-plan usage
- Clarified evidence storage (inline <50 words vs. linked evidence/ files)
- Added decision-file guidance (optional, only if resolves design question)

**Outcome:** Workflow now makes close explicit and traceable. No more vague "save status report if matters."

### 3. `.hawp/kit/usage/status-report.md`

**Changes:**

- Reframed status reports as **optional companions**, not required for every task
- Added "When to Write" guidance (4 scenarios: non-trivial, unproven, pattern-bearing, or cross-session)
- Clarified status reports are companions to plan files, not replacements
- Moved structure down (it's reference, not mandatory)

**Outcome:** Status reports become context-transfer artifacts, not perceived requirements, reducing friction for trivial tasks.

---

## Final Simplified Task Lifecycle

```
INTAKE
  ↓
ANALYZE (fill Context, Analysis sections)
  ↓
PLAN (fill Options, Recommended Fix; write to work/active/<ID>.md)
  ↓
REVIEW GATE (low: auto-approve; medium: ask; high: hold)
  ↓
IMPLEMENT (make changes)
  ↓
VERIFY (test, confirm; note direct evidence vs. unproven)
  ↓
CLOSE (atomic handoff using close checklist):
  1. Fill Outcome + Verification sections in plan file
  2. Move plan from active/ to closed/YYYY/MM/DD/
  3. Link evidence files to work/evidence/YYYY/MM/DD/ (if large)
  4. Write status report (optional: only if non-trivial/unproven/pattern)
  5. Update BACKLOG.md (move to Recently Closed, cap at 10)
  6. Create decision file (optional: only if design decision)
```

---

## New Close Checklist

**Mandatory before marking done in BACKLOG.md:**

```
- [ ] Outcome section filled (what was actually implemented)
- [ ] Verification section filled (all checks listed, each with direct evidence or "unproven" tag)
- [ ] All evidence files referenced exist in work/evidence/YYYY/MM/DD/ or noted as inline
- [ ] Plan file will be moved to work/closed/YYYY/MM/DD/<ID>.md
- [ ] BACKLOG.md row moved from Active to Recently Closed (or marked done)
- [ ] Status report written (optional: only if non-trivial OR unproven OR decision-bearing)
- [ ] Decision file created (optional: only if resolves design question)
```

**Location:** Checklist is part of plan template and intake-workflow quick reference.

---

## What Is Now True

✅ **Task files are the source of truth**

- Plan file contains all context needed to understand what was done, what was proved, what remains open
- Status reports reference plan file; they don't duplicate it
- Backlog is only a navigation index, not a detail store

✅ **Evidence is traceable**

- Every verification claim must have evidence (inline or linked)
- Evidence files are predictably named and dated
- Unproven claims are marked explicitly

✅ **BACKLOG.md is compact**

- Only active items and recently-closed (capped at 10)
- Recently Closed rows link to plan files for detail
- Archive structure (closed/YYYY/MM/DD/) is clear and time-indexed

✅ **Close is atomic**

- 6-step checklist ensures nothing slips between plan file, backlog, evidence, and status report
- Prerequisite: plan file must have Outcome + Verification filled before backlog is updated
- Status reports are optional, not required

---

## What Is Not Done (Deferred for Later Phases)

❌ **ARCHIVE.md index** — Future task to make archived items discoverable by date/type
❌ **Validation script** — Future task to check backlog consistency, evidence completeness, collision detection
❌ **Retrospective close-checklist application** — Closed tasks in 2026/05/ are as-is; new tasks use new checklist
❌ **Decision tree linking** — Decision files exist but backlog doesn't yet link to them; minor future enhancement
❌ **Librarian/indexer** — Future work; workflow foundation is now ready for safe indexing

---

## Follow-Up Work Items

### Suggested next phase (not started):

**TASK-011: Create ARCHIVE.md index for closed work discoverability** (optional)

- Creates work/ARCHIVE.md with date and type indices
- Links Recently Closed table to archive sections
- Makes closed items searchable without opening BACKLOG.md

**TASK-012: Create validation script for backlog/evidence consistency** (optional)

- Checks active/ items have BACKLOG.md rows
- Checks closed/ items have BACKLOG.md or archive entries
- Flags missing evidence files
- Detects ID collisions

**TASK-013: Retrospectively apply close checklist to last 5 closed tasks** (optional)

- Review 2026-05-01 through 2026-05-03 closed items
- Verify they have Outcome, Verification, evidence noted
- Document any gaps found

**TASK-014: Link ADRs to work items in backlog** (optional)

- Existing decisions in decisions/ folder to be cross-referenced
- Update plan template to include decision-link row (if not already done)
- Create example linking for TASK-006's ADR

---

## Phase 1 Completion Checklist

- [x] Updated intake-plan.md template with Outcome, Verification, Close Checklist sections
- [x] Enhanced intake-workflow.md with evidence discipline and structured close steps
- [x] Updated status-report.md to clarify optional/companion role
- [x] Evidence guideline added (inline <50w vs. linked files)
- [x] Close checklist documented and referenceable
- [x] Documentation complete; no breaking changes to existing closed tasks
- [x] No retroactive rewrites of old tasks
- [x] No database/schema/CLI work started

---

## Ready for Librarian Foundation Work

The workflow is now clean enough for a librarian to:

1. Read a closed plan file and trust it has Outcome + Verification sections
2. Find evidence by following named links in work/evidence/YYYY/MM/DD/
3. Distinguish proven claims from unproven claims
4. Build a consistent index without custom reconciliation logic

**Next:** Librarian work can begin with confidence that closed items follow a predictable schema.
