# HAWP Self-Update Review Summary

Date: 2026-05-02
Reviewer: GitHub Copilot (review agent)

---

## Summary

The HAWP repository is in a clean, well-structured state after the recent alignment work. All seven review checks pass or pass with minor observations. The install and update flows correctly preserve project-owned `.hawp/work/` and never copy the source repo's operating state into downstream projects. The dogfooding boundary is intact: `.hawp/work/` holds real project history, `core/.hawp/work/` holds only scaffold files, and the kit lives exclusively in `core/.hawp/kit/`. Backlog alignment policy is correctly placed at `core/.hawp/kit/references/backlog-alignment.md`, the old uppercase path no longer exists, and active agent instructions reference the correct locations. The live backlog is compact (empty Active Work, 10 capped Recently Closed rows, closed plans properly archived by date). The only minor follow-up is that `.hawp/work/STATUS.md` shows "Last Closed: None yet" — stale after TASK-008 closed today — and that `.hawp/work/status/README.md` exists but there is no matching `status/YYYY/MM/DD/` update to STATUS.md from TASK-008 closure. These are housekeeping notes, not structural problems.

---

## Final decision

**Safe to test update**

The repo is clean (git working tree empty), all structural checks pass, and no regressions are present. An update test can proceed.

---

## Checks

### 1. Self-update safety

Status: **pass**
Evidence:

- `core/install.md` line 32: "Existing target `.hawp/work/` files are project-owned and must not be overwritten."
- `core/install.md` line 21: "Never overwrite existing files under `.hawp/work/`"
- `core/install.md` (trust boundary notes): "Install never copies the source repository's repo-root `.hawp/work/` operating state into downstream repositories."
- `core/update.md` header: update touches only `.hawp/LICENSE`, `.hawp/kit/**`, `.github/instructions/*.instructions.md`, `.github/prompts/*.prompt.md`. `.hawp/work/**` is explicitly excluded.
- Update boundary note: "Existing `.hawp/work/**` files are never overwritten."
- Scaffold seeding: `copy_file_no_clobber` helper used throughout — only creates files when absent.
- Kit refresh uses `cp -Rf` (overwrite kit) but `.hawp/work/` is not in the kit refresh path.

### 2. Repo dogfooding boundary

Status: **pass**
Evidence:

- `.hawp/work/` contains real project history: 14 closed plan files under `closed/YYYY/MM/DD/`, a populated BACKLOG with 10 recently closed items, evidence files, and status summaries.
- No root `.work/` directory exists (root contains only `benchmark/`, `core/`, `LICENSE`, `README.md` — old layout fully migrated).
- `core/.hawp/work/BACKLOG.md` contains only scaffold content ("_No active work._", "_No recently closed items._").
- `core/.hawp/work/STATUS.md` contains only scaffold content ("_No active work._").
- All `core/.hawp/work/` subdirs contain only `README.md` placeholders — no real history.
- Root `.hawp/kit/` mirrors `core/.hawp/kit/` — correct; this is the installed kit layer for the repo's own dogfooding use.
- `.github/copilot-instructions.md` instructs agents: "Treat `core/.hawp/work/` as downstream scaffold source only."

### 3. Backlog alignment location

Status: **pass**
Evidence:

- `test -f core/.hawp/kit/references/backlog-alignment.md` → file exists, confirmed.
- `test ! -f core/BACKLOG_ALIGNMENT.md` → old uppercase path absent, confirmed.
- Active docs reference correct path:
  - `.github/instructions/hawp-backlog-alignment.instructions.md` → `core/.hawp/kit/references/backlog-alignment.md`
  - `.github/prompts/hawp-backlog-alignment.prompt.md` → `core/.hawp/kit/references/backlog-alignment.md`
  - `core/.github/instructions/hawp-backlog-alignment.instructions.md` → `.hawp/kit/references/backlog-alignment.md` (installed path — correct for downstream)
  - `core/.github/prompts/hawp-backlog-alignment.prompt.md` → `.hawp/kit/references/backlog-alignment.md` (installed path — correct for downstream)
- Stale references to `core/BACKLOG_ALIGNMENT.md` appear only in archived evidence/status files under `.hawp/work/closed/`, `.hawp/work/status/`, `.hawp/work/evidence/` — intentionally preserved historical records, not active guidance.
- Filename uses lowercase kebab-case (`backlog-alignment.md`). ✓

### 4. Prompt/instruction alignment

Status: **pass**
Evidence:

- `.github/copilot-instructions.md`: "Do not append completed work endlessly to BACKLOG.md. Move closed work to `.hawp/work/closed/YYYY/MM/DD/` and keep BACKLOG.md compact."
- `.github/copilot-instructions.md`: "When BACKLOG.md has many Done rows, create a work item titled 'Compact BACKLOG.md and archive closed work.' before adding more Done entries."
- `.github/copilot-instructions.md`: "Treat `core/.hawp/work/` as downstream scaffold source only." (preserves project-owned work)
- `.github/instructions/hawp-backlog-alignment.instructions.md`: rules include cap Recently Closed to 5-10 items or 14-30 days, archive closed work to `.hawp/work/closed/YYYY/MM/DD/`, do not append every completed item.
- `.github/prompts/hawp-backlog-alignment.prompt.md`: "Do not overwrite project-owned files."
- `core/.github/` scaffold versions of the same files are consistent with source-repo versions, using installed paths.
- No install/update preservation guarantee surfaced directly in `.github/copilot-instructions.md`, but `core/install.md` and `core/update.md` carry it at the implementation level.

### 5. Scaffold cleanliness

Status: **pass**
Evidence:

- `core/.hawp/work/BACKLOG.md`: Active Work = "_No active work._", Recently Closed = "_No recently closed items._" — pure starter template.
- `core/.hawp/work/STATUS.md`: "No active work.", "None yet." — pure starter template.
- All subdirs in `core/.hawp/work/` contain only `README.md` files:
  - `active/README.md`, `closed/README.md`, `decisions/README.md`, `evidence/README.md`, `notes/README.md`, `parked/README.md`, `status/README.md`
- No Tekit records, no repo-local history, no dated plan files exist in `core/.hawp/work/`.
- `core/.hawp/kit/` contains only generic reusable files: templates, patterns, examples, usage guides, spec, references.

### 6. Backlog state

Status: **pass**
Evidence:

- Active Work table is empty — no stale or duplicate open items.
- Recently Closed shows exactly 10 items (TASK-008 through BUG-001) — within the 10-item cap.
- All closed items have plan links pointing to `closed/YYYY/MM/DD/` paths.
- Closed plans confirmed under `.hawp/work/closed/2026/04/26/`, `04/30/`, `05/01/`, `05/02/` — date-based archive model working.
- Blocked/Parked table is empty.
- Minor observation: `.hawp/work/STATUS.md` reads "Last Closed: None yet." — stale after TASK-008 was closed 2026-05-02. Not a structural problem; listed as follow-up.

### 7. Diff risk

Status: **pass**
Evidence:

- `git status --short` → empty (clean working tree).
- `git diff --stat` → no output (no unstaged changes).
- `git diff --name-status` → no output (no unstaged changes).
- Most recent commit `1a07bcb`: "added STATUS.md to surface current work state at a glance" — added `.hawp/work/STATUS.md`.
- `git diff HEAD~1 --name-status` → only `A .hawp/work/STATUS.md` — single-file addition, no deletions or regressions.
- No unexpected deletion of history, no unrelated feature work, no package layout regression, no stale path references in active guidance.

---

## Changes made during review

None.

---

## Follow-up items

1. **Update `.hawp/work/STATUS.md`** — currently reads "Last Closed: None yet." Should reflect TASK-008 as last closed item (2026-05-02). Tiny housekeeping; safe to do in the same session or the next work item.

2. **Consider `status/YYYY/MM/DD/` README sync** — `status/` subdirectory exists with only a `README.md`. Future status reports should be filed at `.hawp/work/status/YYYY/MM/DD/` as per policy (already happening for evidence and manager reviews; this is a reminder to continue).

3. **No compaction needed** — the current backlog is already compact. Active Work is empty, Recently Closed is capped at 10. No separate compaction task required.

---

## Verification commands run

| Command                                                                                         | Result                                                                  |
| ----------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------- |
| `git status --short`                                                                            | Empty — clean working tree                                              |
| `git diff --stat`                                                                               | Empty — no unstaged changes                                             |
| `git diff --name-status`                                                                        | Empty — no unstaged changes                                             |
| `git log --oneline -10`                                                                         | 10 recent commits visible; latest = STATUS.md addition                  |
| `git diff HEAD~1 --name-status`                                                                 | `A .hawp/work/STATUS.md` only                                           |
| `find core/.hawp/kit/references -maxdepth 2 -type f`                                            | `backlog-alignment.md` (1 file, correct)                                |
| `find core/.hawp/work -maxdepth 4 -type f`                                                      | Only README.md scaffold files + BACKLOG.md + STATUS.md (clean)          |
| `find .hawp/work -maxdepth 4 -type f`                                                           | Real project work: BACKLOG, STATUS, closed plans, README scaffolds      |
| `find .hawp/work/closed -type f`                                                                | 14 closed plan files under dated subdirs — archiving working            |
| `rg "core/BACKLOG_ALIGNMENT.md\|BACKLOG_ALIGNMENT.md" README.md core .github .hawp`             | Matches only in archived historical records — zero in active docs       |
| `rg "preserve.*\.hawp/work\|project-owned" core/install.md core/update.md .github core/.github` | Multiple hits in install.md, update.md — preservation boundary explicit |
| `rg "compact backlog\|recently closed\|backlog-alignment.md" README.md core .github .hawp`      | Hits in all correct active locations                                    |
| `test -f core/.hawp/kit/references/backlog-alignment.md`                                        | exit 0 — EXISTS                                                         |
| `test ! -f core/BACKLOG_ALIGNMENT.md`                                                           | exit 0 — ABSENT (confirmed by prior agents + search)                    |

---

## Paste-back summary

> HAWP self-update review complete (2026-05-02). All seven checks passed.
>
> **What passed:**
>
> - Install/update correctly preserves project-owned `.hawp/work/` and never copies the source repo's operating state downstream. The `copy_file_no_clobber` helper and explicit trust-boundary notes confirm this in both `core/install.md` and `core/update.md`.
> - Dogfooding boundary intact: `.hawp/work/` holds real project work (14 closed plans, populated BACKLOG); `core/.hawp/work/` is clean scaffold-only (empty BACKLOG, empty STATUS, README-only subdirs).
> - Old uppercase `core/BACKLOG_ALIGNMENT.md` path is gone; canonical policy lives at `core/.hawp/kit/references/backlog-alignment.md`. Active agent instructions all reference the correct path. Stale references appear only in archived historical records — expected and correct.
> - Agent-facing instructions (`.github/copilot-instructions.md`, `hawp-backlog-alignment.instructions.md`, `hawp-backlog-alignment.prompt.md`) all teach compact-backlog behavior: no endless Done rows, cap Recently Closed, archive under `closed/YYYY/MM/DD/`, preserve project-owned work.
> - Live backlog is already compact: empty Active Work table, 10 capped Recently Closed items, all closed plans properly archived by date.
> - Git working tree is clean with no stale diffs.
>
> **What remains:**
>
> - Minor: `.hawp/work/STATUS.md` reads "Last Closed: None yet." — stale after TASK-008 closed today. Housekeeping follow-up only.
> - No compaction task needed — backlog is already in good shape.
>
> **Verdict:** Safe to test the self-update flow. No structural problems, no regressions, no stale active items, no path drift in active guidance.
