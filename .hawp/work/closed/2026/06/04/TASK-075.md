## Task: TASK-075 — Shared provider behaviors + materialization

**Backlog ID:** TASK-075
**Type:** improvement
**Reported:** 2026-06-04
**Risk Level:** low
**Status:** done
**Closed:** 2026-06-04

---

### Outcome

One canonical source per HAWP integration behavior; provider packs are materialized artifacts. Distribution guides unchanged in structure (12 provider × branch guides).

### Delivered

- `core/providers/shared/behaviors/` — four canonical behavior bodies
- `librarian/scripts/providers/materialize/` — build, validate, composition emit map
- `providers:sync` npm script; wired into `distribution:sync`
- CI drift gate extended for materialized overlays
- Docs: `shared/README.md`, root README, librarian README, `DOC-VERIFICATION.md`

### Verification

Evidence: `.hawp/work/evidence/2026/06/04/TASK-075-shared-behaviors-materialization.md`

```bash
npm --prefix librarian run distribution:sync
```

## Close Checklist

**Before marking this task done in BACKLOG.md, verify ALL:**

- [x] Outcome section filled (what was actually implemented)
- [x] Verification section filled (all checks listed, each with direct evidence or "unproven" tag)
- [x] All evidence files referenced exist in `../evidence/2026/06/04/` or are noted as inline
- [x] Plan file will be moved to `../closed/2026/06/04/TASK-075.md`
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
