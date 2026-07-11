## Task: TASK-073 — Continue provider pack and guides

**Backlog ID:** TASK-073
**Type:** task
**Reported:** 2026-06-04
**Risk Level:** low
**Status:** done
**Closed:** 2026-06-04

---

### Outcome

Shipped Continue provider: `core/providers/.continue/`, distribution sources, four generated guides, and composition registration (12 guides total). Verified against [Continue Rules docs](https://docs.continue.dev/customize/deep-dives/rules).

### Delivered

- [x] Four local rules (`hawp-01-core.md` … `hawp-04-docs-alignment.md`)
- [x] `distribution/sources/providers/continue/**` and `distribution/generated/continue/**`
- [x] `ACTIVE_PROVIDERS` includes `continue`
- [x] Doc verification in `core/providers/DOC-VERIFICATION.md`
- [x] Smoke test: continue-only install (4 rules, no cross-provider paths)

### Verification

```bash
npm --prefix librarian run distribution:validate   # 12/12
```

Evidence: `.hawp/work/evidence/2026/06/04/provider-doc-verification.md`

## Close Checklist

**Before marking this task done in BACKLOG.md, verify ALL:**

- [x] Outcome section filled (what was actually implemented)
- [x] Verification section filled (all checks listed, each with direct evidence or "unproven" tag)
- [x] All evidence files referenced exist in `../evidence/2026/06/04/` or are noted as inline
- [x] Plan file will be moved to `../closed/2026/06/04/TASK-073.md`
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
