# Manager Review — Backlog alignment kit reference

## Summary

Ready to commit.

## Decision reviewed

Backlog alignment policy moved to:
`core/.hawp/kit/references/backlog-alignment.md`

Old path `core/BACKLOG_ALIGNMENT.md` removed.
Downstream installed path: `.hawp/kit/references/backlog-alignment.md`

## Worker reports reviewed

- `.hawp/work/evidence/2026/05/02/worker-backlog-alignment-kit-reference-check.md` — first worker pass (move completed by this worker)
- `.hawp/work/evidence/2026/05/02/worker-backlog-alignment-kit-reference-check-2.md` — second worker verification (confirmed existence, no changes needed)
- `.hawp/work/evidence/2026/05/02/worker-backlog-alignment-kit-reference-check-3.md` — third worker verification (confirmed clean state, no changes needed)

## Findings

### Canonical file location and naming

Status: Pass

Evidence:

- `test -f core/.hawp/kit/references/backlog-alignment.md` → exit 0 (EXISTS)
- File is 149 lines, well-formed, lowercase kebab-case filename.
- Content confirms compact active-index model, Recently Closed cap rules, archive-by-date, agent rules, and acceptance checklist.
- `find core/.hawp/kit/references -maxdepth 2 -type f` shows exactly one file: `backlog-alignment.md`. No duplicates.

### Old path cleanup

Status: Pass

Evidence:

- `test ! -f core/BACKLOG_ALIGNMENT.md` → exit 0 (ABSENT)
- `rg "core/BACKLOG_ALIGNMENT.md|BACKLOG_ALIGNMENT.md" README.md core/install.md core/update.md .github/copilot-instructions.md` → no matches in active docs.
- Stale references to `core/BACKLOG_ALIGNMENT.md` appear only in historical/archived records under `.hawp/work/closed/`, `.hawp/work/status/`, and `.hawp/work/evidence/`. These are intentionally preserved and not updated (correct behavior).

### Reference correctness

Status: Pass

Evidence:

- Source-repo overlay `.github/instructions/hawp-backlog-alignment.instructions.md` references `core/.hawp/kit/references/backlog-alignment.md` — correct source path.
- Source-repo overlay `.github/prompts/hawp-backlog-alignment.prompt.md` references `core/.hawp/kit/references/backlog-alignment.md` — correct source path.
- Downstream scaffold overlay `core/.github/instructions/hawp-backlog-alignment.instructions.md` references `.hawp/kit/references/backlog-alignment.md` — correct installed path.
- Downstream scaffold overlay `core/.github/prompts/hawp-backlog-alignment.prompt.md` references `.hawp/kit/references/backlog-alignment.md` — correct installed path.
- No confusion between source path and installed path detected. The distinction is maintained correctly across the two overlay layers.

### Prompt/instruction alignment

Status: Pass

Evidence:

- Both `.github/instructions/hawp-backlog-alignment.instructions.md` and `core/.github/instructions/hawp-backlog-alignment.instructions.md` explicitly state:
  - `BACKLOG.md` is an active coordination index, not permanent history.
  - `Do not append every completed item forever to BACKLOG.md.`
  - `Keep Recently Closed capped to the last 5-10 items or the last 14-30 days.`
  - `Preserve history in archive files; never delete records just to shorten the backlog.`
  - `If BACKLOG.md already has many Done rows, create or recommend a work item titled "Compact BACKLOG.md and archive closed work." before adding more Done rows.`
- Both `.github/prompts/hawp-backlog-alignment.prompt.md` and `core/.github/prompts/hawp-backlog-alignment.prompt.md` cover the same compact-index checks.
- `.github/copilot-instructions.md` and `core/.github/copilot-instructions.md` reference the new backlog-alignment overlay files.
- All four overlay files are untracked new files, consistent with TASK-007 output.

### Install/update safety

Status: Pass

Evidence:

- `core/install.md` step 4 covers refreshing `core/.hawp/kit/**` — the new `references/` subfolder will be included automatically.
- `core/install.md` explicitly states existing `.hawp/work/` files are project-owned and must not be overwritten.
- `core/update.md` preserves the same boundary: `core/.hawp/work/` is scaffold source only; project-owned `.hawp/work/` is never overwritten.
- `README.md` states "Install/update seed target `.hawp/work/` only from `core/.hawp/work/` scaffold files."
- No stale path references to the old `core/BACKLOG_ALIGNMENT.md` found in install or update scripts.

### Work/scaffold boundary

Status: Pass

Evidence:

- `.hawp/work/` — this repo's project work state. Contains active BACKLOG.md, closed/, evidence/, status/, decisions/. Not included in downstream install.
- `core/.hawp/work/` — scaffold source only. Contains READMEs and starter BACKLOG.md template. No live project records.
- `core/.hawp/kit/` — reusable protocol guidance. Now includes `references/backlog-alignment.md`. No project-local records found.
- No live work records from `.hawp/work/` were copied into `core/.hawp/work/` or `core/.hawp/kit/`.
- `find core/.hawp/work` confirms only scaffold placeholder files (READMEs + BACKLOG.md template).

### Diff risk

Status: Pass

Evidence:

- `git diff --stat` shows 32 files changed with net -2252 lines. The large deletion is entirely from `.work/**` (19 files deleted — migration of legacy `.work/` to `.hawp/work/` completed). This is expected.
- Modified files are all previously reviewed: `.github/**`, `core/.github/**`, `core/.hawp/work/`, `core/install.md`, `core/update.md`, `README.md`, `.gitignore`.
- Untracked new files: `.hawp/` (this repo's work state), `core/.hawp/kit/references/` (canonical policy), `.github/instructions/hawp-backlog-alignment.instructions.md`, `.github/prompts/hawp-backlog-alignment.prompt.md`, `core/.github/instructions/hawp-backlog-alignment.instructions.md`, `core/.github/prompts/hawp-backlog-alignment.prompt.md`, `core/.hawp/work/status/`.
- No duplicate policy files. Only one copy of the alignment policy: `core/.hawp/kit/references/backlog-alignment.md`.
- No uppercase/snake_case naming drift in active files.
- No broken links detected in active docs.
- No unrelated scope added.

## Tiny corrections made

None. State was clean upon review.

## Follow-up items

- **TASK-008** (already in backlog as `inbox`): Compact live `.hawp/work/BACKLOG.md` — replace long `## Done` section with capped Recently Closed + archive links per compact-index policy. This is tracked separately and does not block the current commit.

## Final decision

**Ready to commit.**

All required checks pass. The canonical file is at `core/.hawp/kit/references/backlog-alignment.md` with lowercase kebab-case naming. The old uppercase-named root file is absent. Source-repo and downstream-installed path references are correctly distinct. No work/scaffold boundary regression. No unrelated scope. Three worker evidence artifacts confirm consistent state. TASK-008 for live backlog compaction is registered and does not block this commit.

## Verification commands run

| Command                                                                                                   | Result                                                                                                                                              |
| --------------------------------------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------- |
| `git status --short`                                                                                      | `.work/**` deleted (migration); `.hawp/` and `core/.hawp/kit/references/` untracked (new); expected modifications to core docs and .github overlays |
| `git diff --stat`                                                                                         | 32 files, +131 / -2252; large deletion is `.work/**` migration, not data loss                                                                       |
| `git diff --name-status`                                                                                  | All D entries are `.work/**`; M entries are expected modified files                                                                                 |
| `test -f core/.hawp/kit/references/backlog-alignment.md`                                                  | exit 0 — EXISTS                                                                                                                                     |
| `test ! -f core/BACKLOG_ALIGNMENT.md`                                                                     | exit 0 — ABSENT                                                                                                                                     |
| `rg "core/BACKLOG_ALIGNMENT.md\|BACKLOG_ALIGNMENT.md\|backlog_alignment" README.md core .github .hawp`    | Matches only in archived historical records (expected); zero matches in active docs                                                                 |
| `rg "backlog-alignment.md\|compact backlog\|active index\|Recently Closed" README.md core .github .hawp`  | Matches in all four overlay files at correct paths; matches in scaffold BACKLOG.md and README.md                                                    |
| `rg "preserve.*\.hawp/work\|project-owned" README.md core/install.md core/update.md .github core/.github` | Matches confirmed in install.md and README.md at correct boundary statements                                                                        |
| `find core/.hawp/kit/references -maxdepth 2 -type f \| sort`                                              | `core/.hawp/kit/references/backlog-alignment.md` (one file, correct)                                                                                |
| `find core/.hawp/work -maxdepth 4 -type f \| sort`                                                        | Scaffold READMEs + BACKLOG.md only; no live project records                                                                                         |
