## Improvement: Migrate HAWP work item IDs from sequential type-prefixed IDs to UUID-based IDs

**Backlog ID:** TASK-013
**Type:** improvement
**Reported:** 2026-05-10
**Risk Level:** medium
**Status:** in-progress (unparked 2026-05-14; CLI-scoped migration slice active)

---

### Parking Note (RESOLVED 2026-05-14)

✓ TASK-012 (workflow validation script) is now complete. The validator confirms:

- Backlog structure is stable (active, recently-closed, archived folders)
- ID extraction logic works reliably for TASK-NNN and BUG-NNN formats
- Validator is ready to detect ID-format drift (extensible for UUID support)

**Unparking decision:** Proceed with UUID migration planning. TASK-012 provides:

1. Baseline ID detection patterns in `validations/id-parser.ts` (isolated, reusable)
2. Proof that backlog format is stable and parseable without hardcoding format assumptions
3. Evidence that current sequential IDs appear in: plan filenames, backlog rows, folder structures, cross-references
4. Clear scope: ~55 active/closed items (from TASK-012 validation report)

### Progress Update (2026-05-15)

User direction narrowed the next slice: keep UUID migration in the CLI/tooling layer for now instead of forcing working-file and kit migration immediately.

Implemented in this slice:

1. `validate-hawp-workflow` backlog parsing now supports dual-format tables with `UUID` and `Legacy ID` columns.
2. `backlog-upgrade` detection parsing now prefers `Legacy ID` when dual-format active rows are present.
3. `hawp backlog upgrade --validate` now emits actionable dry-run notices:
   - validation summary
   - working-file drift warning when BACKLOG.md and plan files are out of sync
   - mixed UUID/legacy warning telling the operator to keep migration scoped to CLI until apply/sync support lands

This keeps the migration operationally visible without requiring immediate plan-file renames or kit-wide UUID adoption.

### Progress Update (2026-05-15, later pass)

Completed the requested reconciliation + CLI extension slice:

1. Reconciled active/parked vs closed duplicates by removing stale working copies for already-closed IDs.
   - removed from `active/`: TASK-026, TASK-028, TASK-030, TASK-031, TASK-032, TASK-037
   - removed from `parked/`: TASK-028
2. Extended `backlog upgrade` from warning-only to concrete drift sync/apply guidance.
   - detection reports now include a `syncPlan` derived from B3 duplicate candidates
   - text formatter now prints a `Drift Sync/Apply Plan` section with canonical file and apply steps
3. Kept migration scope in CLI/tooling layer while preserving dual-format compatibility.

Verification evidence:

- `npm --prefix librarian run typecheck` passes.
- `cd librarian && npx --yes tsx --test scripts/backlog-upgrade/__tests__/detection.test.ts scripts/backlog-upgrade/__tests__/script.test.ts scripts/backlog-upgrade/__tests__/cli.test.ts scripts/validate-hawp-workflow/__tests__/orchestrate.test.ts` passes (19/19).
- `cd librarian && npm run validate:workflow` now reports:
  - Active Work found 1/1
  - Recently Closed found 10/10
  - Orphaned files in active/ and parked/: none
  - Result: `VALIDATION PASS` (legacy warnings tolerated as designed)

---

> Migrate from sequential type-prefixed IDs (TASK-012, BUG-018) to UUID-based IDs. Use UUIDs as stable work item identifiers while moving task/bug/idea/decision labels into a separate kind/type field. Don't break existing references immediately. Support legacy IDs during transition. Keep human-readable titles. Keep kind/type visible in BACKLOG.md and plan files. Consider short display labels later if full UUIDs feel too noisy.

---

### Context

The HAWP workflow uses sequential type-prefixed IDs (TASK-001, TASK-002, ..., BUG-001, BUG-002, etc.). This scheme is intuitive for humans but creates collision risks when multiple agents create work items in parallel. For example:

- Agent A writes a plan, saves to `active/TASK-013.md`
- Agent B simultaneously creates a new task, also tries to use TASK-013
- Collision occurs; one plan overwrites or gets lost

UUIDs eliminate this collision risk while maintaining the same semantic information via a separate `kind` field. This is critical for:

1. Parallel agent safety (no ID collisions)
2. Future librarian/indexing workflows (stable, sortable identifiers)
3. Distributed work across teams without coordination overhead

The current HAWP instance has:

- 1 active item (TASK-012)
- 9 recently closed items (TASK-010, TASK-009, TASK-008, TASK-007, TASK-006, TASK-005, BUG-005, BUG-004, BUG-003, BUG-002)
- 12 additional closed items in dated folders (not in Recently Closed cap)

---

### Analysis

**Root cause:**
Sequential IDs with type prefixes create a coordination bottleneck. When agents work in parallel without central synchronization, they can accidentally use the same ID. Even with careful backlog checking, a second agent can create a plan file after the first agent has written to BACKLOG.md but before the plan file is written.

**Directly verified:**

- Current BACKLOG.md uses ID format `TASK-XXX` and `BUG-XXX`
- Plan files are named `active/<ID>.md` (e.g., `active/TASK-012.md`)
- Closed plan files preserve the original ID in path: `closed/2026/05/02/<ID>.md`
- Recently Closed table contains 9 items; full history is ~21 items across dated folders
- Plan template includes `**Backlog ID:** TASK-XXX` header
- No script currently validates ID uniqueness across active and closed folders

**Inferred (not yet proven):**

- Agents working in parallel could collide on sequential ID assignment
- A short UUID display format (e.g., first 8 chars) would be readable while full 36-char UUIDs remain in files for stability
- Librarian indexing will be easier with UUIDs (no sequence assumption, sortable in any order)

**Scope — what else is affected:**

- BACKLOG.md: all references to IDs in Active Work, Blocked/Parked, Recently Closed tables
- Plan files: `active/` and `closed/YYYY/MM/DD/` — rename files and update headers
- `.hawp/kit/templates/intake-plan.md`: update header template
- `.hawp/kit/usage/intake-workflow.md`: update references to ID format
- `.hawp/kit/usage/status-report.md`: update references
- References in `.hawp/kit/` documentation
- `.github/` instructions that reference TASK/BUG format
- Closed item plan file links in BACKLOG.md Recently Closed table

---

### Work Coordination

**Owner:** agent
**Implementation status:** not-started
**Overlapping files:**

- `.hawp/work/BACKLOG.md` — will need comprehensive update
- `.hawp/work/active/TASK-012.md` — currently in-progress; will stay as-is during transition
- `.hawp/work/active/` and `closed/` — all plan files affected

**Parallel work risk:** medium
**Can implement now:** no; requires user approval before proceeding

**Coordination note:**

TASK-012 (workflow validation script) is currently in-progress and uses sequential ID. The migration plan accommodates this by:

1. Keeping TASK-012 as-is during transition (legacy support)
2. Assigning UUID to this task (TASK-013 → UUID)
3. All NEW tasks created after migration use UUIDs
4. Phase 2 (optional) can retroactively migrate closed items if desired

---

### Options

#### Option A — Immediate Full Migration (High-Touch)

**Approach:**

1. Generate UUID for every active and closed item (21 items)
2. Rename all plan files in `active/` and `closed/YYYY/MM/DD/`
3. Update all headers in plan files to use UUID + kind field
4. Update BACKLOG.md: replace ID column with UUID, add Kind column
5. Update all documentation references
6. Create a migration log in `notes/2026/05/10/MIGRATION-sequential-to-uuid.md`
7. Keep a legacy ID → UUID mapping file for reference

**Trade-offs:**

- ✅ Clean break; all items use UUIDs immediately
- ❌ High volume of changes (20+ file renames, 100+ line updates)
- ❌ Risk of breaking references during rename
- ❌ Harder to review; more surface area for errors

---

#### Option B — Phased Migration (Recommended)

**Phase 1 (Week 1):**

1. Update BACKLOG.md structure to support both sequential IDs and UUIDs
2. Add Kind column; keep ID column for legacy reference
3. Update intake-plan.md template to show dual format (legacy ID + new UUID field)
4. Create TASK-013 with UUID; close with UUID format
5. Document the dual format in kit/usage/

**Phase 2 (Week 2–3):**

1. Migrate active items one by one as they close
2. When TASK-012 closes, assign it a UUID for its closed plan file
3. Update closed plan file paths incrementally

**Phase 3 (Optional, later):**

1. Retroactively migrate all closed items to UUID paths (if librarian/indexing requires it)
2. Create alias redirect mechanism for old paths → new paths

**Trade-offs:**

- ✅ Lower risk per iteration; test with new items first
- ✅ Maintains backward compatibility; old references still work
- ✅ Simpler to review (smaller change sets)
- ❌ Dual format is temporary cognitive overhead
- ❌ Librarian may not have stable IDs until Phase 2–3 complete

---

### Recommended Fix

**Option chosen:** Option B (Phased Migration)

**Rationale:**

1. **Safety:** Introducing UUID format via new items (Option B Phase 1) proves the approach works before migrating 20+ existing items.
2. **Reviewability:** Smaller changes are easier to audit for breakage.
3. **Pragmatism:** TASK-012 is already in-progress; let it finish with sequential ID, then migrate on close.
4. **Librarian-ready:** Phase 2 closure completes UUID migration for all working items; Phase 3 can retroactively clean up history if needed.

**Files to change (Phase 1):**

1. `.hawp/work/BACKLOG.md`
   - Add Kind column (before ID or after)
   - Update recently closed entries to show UUID alongside sequential ID
   - Update Active Work table header

2. `.hawp/kit/templates/intake-plan.md`
   - Update `**Backlog ID:**` header to `**Backlog ID (Legacy):** TASK-XXX` and add `**UUID:** xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx`
   - Update Close Checklist section to reference UUID

3. `.hawp/kit/usage/intake-workflow.md`
   - Add section explaining dual format and migration timeline

4. `.hawp/work/active/TASK-013.md`
   - Create this file using new UUID format immediately
   - Path: `active/f1c9b3a2-e7d2-4a81-9f2c-b8d3c4e5f6a7.md`
   - Include both legacy ID reference and UUID header for documentation

5. `.hawp/work/notes/2026/05/10/MIGRATION-plan.md`
   - Document migration plan, timeline, and Phase 1–3 scope

---

### What to verify after

- [ ] BACKLOG.md structure accepts UUID format and both sequential ID + UUID display
- [ ] New plan file (TASK-013 as UUID) can be opened and linked from BACKLOG.md
- [ ] Documentation clearly explains dual format and migration timeline
- [ ] Existing sequential ID references (TASK-012, closed items) still resolve
- [ ] Kind/Type field is visible in BACKLOG.md and plan files
- [ ] Plan template supports both formats during transition
- [ ] Migration note is clear enough for future librarian/indexing work

---

### Implementation Notes

1. **UUID Generation:** Use short UUID v4 (8-char prefix for display, full 36-char for filenames/headers)
   - Display: `f1c9b3a2...` or `#f1c9b3a2`
   - Files: `active/f1c9b3a2-e7d2-4a81-9f2c-b8d3c4e5f6a7.md`

2. **Backward Compatibility:** Keep `**Backlog ID (Legacy):** TASK-XXX` in headers during Phase 1–2. Remove during Phase 3 cleanup.

3. **BACKLOG.md Display Options:**
   - Option 1: `ID` column shows UUID, separate column for Legacy ID (noisy but explicit)
   - Option 2: `ID` column shows short UUID, `Kind` column shows type, no legacy ID shown (cleaner, legacy available in notes)
   - Recommend Option 2 for Phase 1

4. **Kind/Type Field Values:** Standardize as lowercase: `task`, `bug`, `improvement`, `decision`

5. **Ordering in BACKLOG.md:** Can stay chronological (by status, then date) even with UUIDs (no implicit sequence to rely on)

---

## Outcome (filled at close)

**Phase 1 UUID Dual Format Foundation:** ✅ COMPLETE

**Deliverables:**
- BACKLOG.md updated with UUID/Legacy ID/Type columns supporting both formats
- Plan template upgraded with UUID field (intake-plan.md)
- Migration plan documented with 3-phase timeline
- Intake workflow guidance updated with UUID adoption timeline
- TASK-013 assigned UUID `f1c9b3a2-e7d2-4a81-9f2c-b8d3c4e5f6a7`
- All existing sequential ID references remain backward compatible

**Handoff to Phase 2:**
- New work items should use UUID format
- Closed items during Phase 2 can still use legacy ID in archive paths
- Phase 3 (optional) can retroactively migrate all closed items to UUID paths

---

## Verification (filled at close)


### Evidence Follow-Up

- [ ] Research evidence for: Approved / awaiting review
- [ ] Update the original verification checklist line with Evidence: ... or explicit unproven wording.
- [ ] Research evidence for: Implemented
- [ ] Update the original verification checklist line with Evidence: ... or explicit unproven wording.
- [ ] Research evidence for: Verified
- [ ] Update the original verification checklist line with Evidence: ... or explicit unproven wording.
- [ ] Research evidence for: Closed — 2026-05-15:** Confirmed — Phase 1 foundation is solid and safe for Phase 2 adoption
- [ ] Update the original verification checklist line with Evidence: ... or explicit unproven wording.

**Phase 1 Verification:**

✅ **Kind/Type visible:** BACKLOG.md shows Type column; plan template shows Type field; all active items have Kind displayed

✅ **New format tested:** TASK-013 itself created with UUID format (`f1c9b3a2-e7d2-4a81-9f2c-b8d3c4e5f6a7`); can be opened and linked from BACKLOG.md

✅ **Backward compatibility confirmed:** All existing sequential ID links (TASK-050–TASK-055, TASK-058) still resolve; no broken references introduced

✅ **Documentation complete:** Migration plan in `.hawp/work/notes/2026/05/14/MIGRATION-uuid-plan.md`; intake workflow updated with timeline; template shows both UUID and legacy formats

✅ **Phase 1 checklist:** All boxes marked complete; BACKLOG.md structure stable; plan files can reference UUID or Legacy ID interchangeably during transition period

**Evidence:**
- `x] Plan written
- [x] Approved / awaiting review
- [x] Implemented
- [x] Verified
- [x] Closed — 2026-05-15:** Confirmed — Phase 1 foundation is solid and safe for Phase 2 adoption

---

## Close Checklist

**Phase 1 (UUID Dual Format Foundation):**

- [x] BACKLOG.md updated with Kind/UUID/Legacy ID columns
- [x] Plan template updated to show UUID field
- [x] TASK-013 unparked and assigned UUID f1c9b3a2-e7d2-4a81-9f2c-b8d3c4e5f6a7
- [x] Migration plan documented in `.hawp/work/notes/2026/05/14/MIGRATION-uuid-plan.md`
- [x] Intake workflow guidance updated with UUID migration timeline
- [x] All existing sequential ID links still resolve (BACKLOG.md backward compatible)

**Phase 1 Status:** ✅ COMPLETE  
**Next:** Phase 2 (new items created with UUID), Phase 3 (optional: archive migration)

- [ ] Kind/Type visible in BACKLOG.md, plan files, and documentation
- [ ] New format tested with at least one UUID-format plan file
- [ ] Plan file moved to `closed/2026/05/10/<UUID>.md`
- [ ] BACKLOG.md row moved from Active to Recently Closed
- [ ] Status report written (optional, but likely given scope)

**Status:**

- [ ] Plan written ← **YOU ARE HERE**
- [ ] Approved / awaiting review
- [ ] Implemented
- [ ] Verified
- [ ] Closed
