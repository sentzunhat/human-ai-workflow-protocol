# Remove Legacy Root .work Cleanup

## Summary

Legacy root `.work/` is already retired in this repository. Active guidance is aligned to repo-root `.hawp/work/`, and install/update boundaries remain aligned with project-owned `.hawp/work/**` preservation.

## Files checked

- `.work/**` (presence + file listing)
- `.hawp/work/**`
- `.hawp/work/BACKLOG.md`
- `.hawp/work/closed/**`
- `.hawp/work/decisions/**`
- `.hawp/work/evidence/**`
- `.hawp/work/status/**`
- `.hawp/work/notes/**`
- `README.md`
- `core/install.md`
- `core/update.md`
- `.github/copilot-instructions.md`
- `.github/instructions/**`
- `.github/prompts/**`
- `core/.github/**`
- `core/.hawp/kit/**`
- `core/.hawp/work/**`
- `.gitignore`

## Migration result

`removed safely`.

Observed state:

- Root `.work/` directory is absent.
- Repo-local work state exists under `.hawp/work/` with active backlog, closed archives, decisions, evidence, status, and notes trees.
- No additional unmigrated `.work/**` files were found because `.work/` no longer exists.

## Changes made

- Added `.hawp/work/status/2026/05/03/remove-legacy-root-work-cleanup.md` (this report).

## Work history preservation

Historical records remain preserved under `.hawp/work/closed/**`, `.hawp/work/decisions/**`, `.hawp/work/evidence/**`, `.hawp/work/status/**`, and `.hawp/work/notes/**`.

No archived historical files were rewritten during this cleanup pass.

## Stale references removed

No active stale `.work/` guidance references required edits in this pass.

Notes:

- Active docs/instructions/prompts already point to `.hawp/work/`.
- Remaining `.work/` mentions are in historical archived records and an explicit retirement note in `README.md`, which is intentional.

## Validation results

- `git status --short`
  - Result: only expected working-tree changes; root `.work/` not present.
- `git diff --stat`
  - Result: includes this cleanup report addition; no boundary-breaking structural changes.
- `find .work -maxdepth 5 -type f | sort 2>/dev/null || true`
  - Result: `.work` not found.
- `find .hawp/work -maxdepth 5 -type f | sort`
  - Result: expected populated `.hawp/work/**` structure present.
- `rg -n "\.work/|root \.work|\.work/BACKLOG\.md" README.md core .github .hawp .gitignore`
  - Result: matches only in intentional historical records and README retirement note.
- `rg -n "core/BACKLOG_ALIGNMENT.md|BACKLOG_ALIGNMENT.md|backlog_alignment|BacklogAlignment" README.md core .github .hawp`
  - Result: no active guidance matches; archived historical references only.
- `rg -n "single source of truth|Single source of truth|Done forever|permanent history" README.md core .github .hawp`
  - Result: active language supports compact-index model; archived records contain historical wording.
- `rg -n "preserve.*\.hawp/work|overwrite.*\.hawp/work|project-owned" README.md core/install.md core/update.md .github core/.github`
  - Result: install/update boundary language confirms preservation of project-owned `.hawp/work/**`.
- `test -f core/.hawp/kit/references/backlog-alignment.md`
  - Result: present.
- `test ! -f core/BACKLOG_ALIGNMENT.md`
  - Result: old path absent.

## Follow-up items

None.
