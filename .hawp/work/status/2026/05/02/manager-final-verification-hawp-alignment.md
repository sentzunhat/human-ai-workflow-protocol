# Manager Final Verification — HAWP alignment work

## Summary

Final readiness decision:

- Needs follow-up work

## Worker verification reviewed

- Report path:
  - `.hawp/work/evidence/2026/05/02/worker-verification-hawp-alignment.md`
- Worker decision:
  - Not available (file not found at the required path)
- Manager assessment of worker report:
  - Worker verification quality cannot be confirmed because the expected report artifact is missing.
  - This is a process/evidence gap, not a protocol-architecture failure.

## Final checks

### 1. Three-folder boundary

Status:

- Pass
  Evidence:
- `README.md` describes:
  - repo-local work state at `.hawp/work/`
  - reusable kit source at `core/.hawp/kit/`
  - downstream scaffold source at `core/.hawp/work/`
- `core/.hawp/work/README.md` explicitly says this tree is scaffold source only.
- `core/.hawp/work/` contains generic starter files only (README + starter BACKLOG + folder README anchors), no source-repo live histories.

### 2. Install/update safety

Status:

- Pass
  Evidence:
- `core/install.md` trust-boundary text states downstream `.hawp/work/` is seeded from `core/.hawp/work/` scaffold only.
- `core/update.md` states source-repo `.hawp/work/` is never copied into downstream repositories.
- Both docs/scripts state no-overwrite behavior for existing target `.hawp/work/**` and scaffold seeding only when files are missing.
- `.github/copilot-instructions.md` is seeded only when missing (not overwritten).

### 3. Backlog compactness policy

Status:

- Partial
  Evidence:
- Compact backlog policy exists in `core/BACKLOG_ALIGNMENT.md`.
- Scaffold template `core/.hawp/work/BACKLOG.md` now reflects compact model (Recently Closed + Archive, capped guidance).
- Root and core intake instructions/prompts align with `.hawp/work/` path and source/scaffold boundary.
- Live repo backlog file `.hawp/work/BACKLOG.md` still includes a long `## Done` history section and lacks a capped Recently Closed section.

### 4. Source repo dogfooding

Status:

- Pass
  Evidence:
- Source-repo state exists under `.hawp/work/**` including backlog, active/closed history, decisions, notes, and status reports.
- `.work/**` content appears migrated to `.hawp/work/**`; working tree shows `.work/**` deletions while equivalent records exist under `.hawp/work/**`.
- ADR and migration notes under `.hawp/work/decisions/...` and `.hawp/work/notes/...` document this migration.

### 5. Scaffold cleanliness

Status:

- Pass
  Evidence:
- `core/.hawp/work/**` contains only starter/scaffold artifacts.
- No source-repo closed items, decisions, status history, or evidence artifacts are present in `core/.hawp/work/**`.

### 6. Instructions/prompts alignment

Status:

- Partial
  Evidence:
- `.github/instructions/intake.instructions.md` and `core/.github/instructions/intake.instructions.md` both direct project work to `.hawp/work/` and prohibit writing repo work into `core/.hawp/work/`.
- `.github/prompts/intake.prompt.md` and `core/.github/prompts/intake.prompt.md` reinforce same boundary.
- `.github/copilot-instructions.md` still says “Track all open and completed work in .hawp/work/BACKLOG.md,” which can be interpreted as allowing perpetual done accumulation; this should be tightened to compact-index wording.

### 7. Stale references and links

Status:

- Partial
  Evidence:
- Stale-reference search finds `root .work`/`.work/` mostly in historical artifacts:
  - `.hawp/work/closed/**`
  - `.hawp/work/notes/2026/05/01/migration-note.md`
  - `.hawp/work/active/TASK-006.md` (historical problem statement)
- Current active policy docs (`README.md`, `core/install.md`, `core/update.md`, intake instructions/prompts) are aligned to `.hawp/work`.
- “single source of truth” appears in general workflow language in intake workflow docs; not a direct permanent-history instruction by itself, but should be interpreted consistently with compact backlog policy.

### 8. Diff risk

Status:

- Partial (acceptable with follow-up)
  Evidence:
- `git diff --stat` shows changes are in expected scope: docs/instructions/install-update/backlog-alignment + `.work` to `.hawp/work` migration.
- No signs of CLI/DB/automation/protocol-runtime expansion.
- No package-layout merge between `core/.hawp/kit` and `core/.hawp/work`.
- Historical records do not appear deleted outright; they appear migrated from `.work/**` to `.hawp/work/**` (still present in repo under new location).
- Missing worker verification artifact is a release-readiness process gap.

## Tiny corrections made

- None.

## Required follow-up items

- Create/recover worker verification evidence artifact at `.hawp/work/evidence/2026/05/02/worker-verification-hawp-alignment.md` (or agreed equivalent path) so manager review has complete chain-of-evidence.
- Compact live `.hawp/work/BACKLOG.md` into active-index format (capped Recently Closed + archive links) in a dedicated work item.
- Tighten wording in `.github/copilot-instructions.md` from “open and completed work in BACKLOG.md” to compact-index guidance consistent with `core/BACKLOG_ALIGNMENT.md`.
- Optional: add a short release note documenting `.work -> .hawp/work` dogfooding migration and historical-reference expectations.

## Final decision

- Needs follow-up work

Brief reason:

- Architecture and trust boundaries are in good shape, but final verification is not fully complete because the worker evidence report is missing and compact-backlog policy is not yet reflected in the live source backlog file.

## Verification commands run

- `git status --short`
  - Result: expected in-scope modifications plus `.work/**` deletions and `.hawp/**` additions.
- `git diff --stat`
  - Result: changes concentrated in docs/instructions/install-update and `.work` migration artifacts.
- `git diff --name-status`
  - Result: no unrelated feature files; mostly policy/docs and migration paths.
- `git --no-pager diff -- README.md core/install.md core/update.md core/BACKLOG_ALIGNMENT.md core/.hawp/work .github core/.github .hawp/work .gitignore`
  - Result: alignment edits in expected files; no trust-boundary regression found.
- `rg -n --hidden -uu "\.work/|root \.work|single source of truth|Done forever|permanent history" README.md core .github .hawp`
  - Result: remaining `.work` references are primarily historical artifacts; active policy files are mostly clean.
- `rg -n --hidden -uu "core/\.hawp/work|\.hawp/work|preserve|overwrite|scaffold" README.md core .github .hawp`
  - Result: install/update and instructions consistently express source/scaffold/preservation boundary.
- `find .hawp/work -maxdepth 6 -type f | sort`
  - Result: repo-local work state present under `.hawp/work/**` with migrated history.
- `find core/.hawp/work -maxdepth 6 -type f | sort`
  - Result: scaffold tree remains generic starter content only.
