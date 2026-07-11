## Task: TASK-074 — Provider packs + split distribution guides

**Backlog ID:** TASK-074
**Type:** task
**Reported:** 2026-06-04
**Risk Level:** medium
**Status:** done
**Closed:** 2026-06-04

---

### Outcome

All agent provider phases complete: GitHub, Cursor, and Continue packs with per-provider distribution guides and composed install/update scripts.

### Phases

- [x] Phase 1 — GitHub move to `core/providers/.github/`
- [x] Phase 2 — Split guides under `distribution/generated/<provider>/`
- [x] Phase 2b — Provider-scoped overlays (`PROVIDER` per guide)
- [x] Phase 3 — Cursor pack and guides
- [x] Phase 3b — Composed scripts (core + provider + footer)
- [x] Phase 4 — Continue (TASK-073)

### Verification

```bash
npm --prefix librarian run distribution:validate   # 12/12 guides
```

Evidence: `.hawp/work/evidence/2026/06/04/provider-doc-verification.md`

Related closed plan: TASK-072 (architecture). Continue implementation: TASK-073.

## Close Checklist

**Before marking this task done in BACKLOG.md, verify ALL:**

- [x] Outcome section filled (what was actually implemented)
- [x] Verification section filled (all checks listed, each with direct evidence or "unproven" tag)
- [x] All evidence files referenced exist in `../evidence/2026/06/04/` or are noted as inline
- [x] Plan file will be moved to `../closed/2026/06/04/TASK-074.md`
- [x] BACKLOG.md row moved from Active to Recently Closed (or marked done)
- [x] Status report written (optional: only if non-trivial or if something remains unproven)
- [x] Decision file created if applicable (only if this task resolves a design question)
- [x] Staged-path proof captured before commit:
  - [x] `git diff --name-status`
  - [x] `git diff --check`
  - [x] `git diff --cached --name-status`
  - [x] `git diff --cached --check`
  - [x] `git status --short`
- [x] If basename-only paths appeared during path-sensitive work, mark report unsafe and restart after correction (fresh path-locked pass after two failures)

**Status:**

- [x] Plan written
- [x] Approved / auto-approved (low risk)
- [x] Implemented
- [x] Verified
- [x] Closed
