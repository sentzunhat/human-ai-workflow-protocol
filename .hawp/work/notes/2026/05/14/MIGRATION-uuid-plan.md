# Migration Plan: Sequential IDs → UUID-Based IDs

**Date Started:** 2026-05-14  
**Dependency:** TASK-012 (workflow validation) — ✓ Complete  
**Tracking:** TASK-013 (UUID migration improvement)

---

## Executive Summary

Migrate HAWP work item IDs from sequential type-prefixed format (`TASK-001`, `BUG-018`) to UUID-based format (`f1c9b3a2-e7d2-4a81-9f2c-b8d3c4e5f6a7`) to enable parallel-safe work item creation. Phased approach: 3 phases over 2–4 weeks with backward compatibility during transition.

**Key motivations:**

- Parallel-safe ID generation (no collisions when multiple agents create tasks simultaneously)
- Future-proof for librarian indexing and distributed workflows
- Backward compatible during Phase 1–2 (legacy IDs remain readable)

---

## Phase 1: Dual Format Foundation (Week 1)

**Goal:** Prove UUID format works; support both sequential and UUID IDs in tooling.

### Changes

1. **BACKLOG.md Structure**
   - Add `UUID` column (short display format)
   - Add `Legacy ID` column (reference during migration)
   - Rename `ID` → columns split as above
   - Keep `Type`/`Kind` column for clarity
   - Format: `f1c9b3a2` (8-char short) in display; full `f1c9b3a2-e7d2-4a81-9f2c-b8d3c4e5f6a7` in files

2. **Plan Template** (`.hawp/kit/templates/intake-plan.md`)
   - Update header to show both formats:
     ```
     **UUID:** f1c9b3a2-e7d2-4a81-9f2c-b8d3c4e5f6a7
     **Legacy ID (Migration Period):** TASK-NNN
     **Type:** task | bug | improvement
     ```
   - Update Close Checklist to reference UUID

3. **Documentation** (`.hawp/kit/usage/intake-workflow.md`)
   - Add section: "During Migration: Sequential and UUID Support"
   - Explain both formats are valid
   - Timeline: Phase 1 (June), Phase 2 (late June), Phase 3 (optional, cleanup)

4. **TASK-013 Itself**
   - Create with UUID: `f1c9b3a2-e7d2-4a81-9f2c-b8d3c4e5f6a7`
   - File path: `active/TASK-013.md` (keep legacy path during Phase 1 for simplicity)
   - Header shows both UUID and Legacy ID
   - Serve as reference implementation

### Acceptance Criteria

- [x] BACKLOG.md shows both UUID and Legacy ID columns
- [x] New plan template supports UUID field
- [x] TASK-013 created with UUID header
- [x] Migration plan documented (this file)
- [x] Intake workflow guidance updated
- [x] All existing sequential ID links still resolve

---

## Phase 2: New Items Use UUID (Week 2–3)

**Goal:** All new work items created during Phase 2 use UUIDs. Closed items retain legacy IDs in archive until Phase 3.

### Changes

1. **Intake Workflow Update**
   - When creating new item, agent generates UUID (via `uuidgen` or online tool)
   - Assigns both UUID and legacy reference in header (for 1 cycle of cross-reference)

2. **BACKLOG.md**
   - Recently Closed section: still shows legacy IDs (TASK-024, BUG-010)
   - Will migrate to UUID format during Phase 3 cleanup

3. **File Paths**
   - New items: `active/{UUID}.md`
   - Closed items: `closed/YYYY/MM/DD/{legacy-id}.md` (unchanged during Phase 2)
   - Phase 3 can retroactively rename closed files to UUID paths

### Example: Closing an Item During Phase 2

```
**UUID:** g2a5d1b8-c3f4-4e9a-8d2f-e7c9a4b1f6d3
**Legacy ID:** TASK-025
**Type:** task

[... work content ...]

## Close Checklist

- [x] All items complete
- [x] File path: closed/2026/05/20/TASK-025.md (legacy during Phase 2)
- [x] Next phase will migrate to UUID path if needed
```

### Acceptance Criteria

- [ ] First new task created with UUID (not sequential ID)
- [ ] UUID format accepted by backlog parser
- [ ] Legacy ID still visible in closed archive for reference
- [ ] No breakage to existing sequential ID references

---

## Phase 3: Full UUID Archive (Optional, Late June)

**Goal:** Retroactively migrate all closed items to UUID paths. Optional if librarian/indexing doesn't require it.

### Changes

1. **File Renaming**
   - Rename `closed/YYYY/MM/DD/{legacy-id}.md` → `closed/YYYY/MM/DD/{full-uuid}.md`
   - Update BACKLOG.md Recently Closed links to new paths

2. **Cleanup**
   - Remove Legacy ID columns from BACKLOG.md (if desired)
   - Update documentation to reflect UUID-only format
   - Create archive of legacy ID → UUID mapping for historical reference

3. **Optional: Path Redirect**
   - Keep symlinks or redirect mechanism if external systems reference old paths
   - Document in `.hawp/kit/lib/legacy-path-mapping.md`

### Acceptance Criteria

- [ ] All closed items migrated to UUID paths
- [ ] BACKLOG.md Recently Closed section shows UUIDs only
- [ ] Legacy ID mapping document created
- [ ] No broken references (all old links updated or redirected)

---

## Reference: UUID Generation

### Quick Gen (macOS/Linux)

```bash
# Generate full UUID
uuidgen

# Output: F1C9B3A2-E7D2-4A81-9F2C-B8D3C4E5F6A7

# Short form (first 8 chars)
uuidgen | cut -c1-8
# Output: f1c9b3a2
```

### Online (if needed)

- https://www.uuidgenerator.net/ (v4 recommended)
- Copy first 8 chars for display short form

---

## Timeline

| Phase | Week                 | Milestone                               | Owner | Status      |
| ----- | -------------------- | --------------------------------------- | ----- | ----------- |
| 1     | Week of 2026-05-14   | Dual format support, TASK-013 with UUID | agent | In Progress |
| 2     | Week of 2026-05-20   | First new task created with UUID        | agent | Pending     |
| 3     | Late June (optional) | All closed items migrated to UUID paths | agent | Pending     |

---

## Risk & Mitigations

| Risk                                                   | Mitigation                                                                                  |
| ------------------------------------------------------ | ------------------------------------------------------------------------------------------- |
| Broken references during migration                     | Keep legacy ID references visible during Phase 1–2; test all links before Phase 3           |
| Confusion during dual format                           | Document clearly; add migration guidance to intake workflow; examples in README             |
| External tools/scripts hardcoded for sequential format | Audit `.github/` instructions, CI scripts, and librarian code before Phase 2                |
| UUID generation collisions (extremely rare)            | Use v4 UUIDs (cryptographic randomness); collision probability negligible for <10,000 items |

---

## Rollback Plan (If Needed)

If UUID format causes issues:

1. Revert BACKLOG.md to sequential ID format
2. Move UUID-created items back to sequential IDs (renumber appropriately)
3. Archive migration plan for future retry
4. Document lessons learned in `.hawp/kit/decisions/YYYY/MM/DD/UUID-migration-decision.md`

Current status: **No rollback needed; proceeding with Phase 1.**

---

## Notes for Future Phases

- **Phase 2 Checkpoint:** Review after 3–5 new items created with UUID; confirm no unexpected issues.
- **Phase 3 Decision:** Depends on librarian/indexing requirements. Defer if not needed for core workflow.
- **Documentation:** Keep this file updated as phases complete; move to `.hawp/kit/` or archived decisions folder when done.
