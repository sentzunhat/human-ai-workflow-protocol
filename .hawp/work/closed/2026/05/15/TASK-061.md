## Task: Review Zacatl public-lane standards against private boundary rules

**Backlog ID:** TASK-061
**Type:** review
**Reported:** 2026-05-15
**Risk Level:** medium
**Status:** done

---

### Input (what was reported)

Derived from TASK-058 triage: `shared_standards/public/standards/zacatl/**` is currently in a public lane but triggers rubric red flags due to internal-domain reference.

---

### Context

TASK-058 classified all Zacatl standards entries as `private/proprietary` with `split-private` pending explicit clearance.

---

### Analysis

**Root cause:**

Path-level public placement conflicts with rubric red-flag treatment for internal-only domain references.

**Directly verified:**

- 7 files under `shared_standards/public/standards/zacatl/**` were triaged to `split-private` in `.hawp/work/evidence/2026/05/15/TASK-058-shared-standards-triage.md`.

**Inferred (not yet proven):**

- A targeted boundary decision is needed to determine whether these are truly public-safe framework docs or internal architecture guidance.

---

### Work Coordination

**Owner:** agent
**Implementation status:** done
**Overlapping files:**

- `shared_standards/public/standards/zacatl/**`
- `.hawp/work/active/TASK-054.md`
- `.hawp/kit/guidance/shared-standards-review-rubric.md`

**Parallel work risk:** medium
**Can implement now:** yes (approved by continuation request)

---

### Recommended Fix

1. Perform content-level review of each Zacatl doc against red-flag patterns.
2. Keep any internal/private guidance in split-private lane.
3. If any section is portable and public-safe, extract only generalized guidance as a separate candidate.

**What to verify after:**

- [x] Every Zacatl file has a boundary decision.
- [x] Internal/private architecture details stay out of core standards.
- [x] Any approved public-safe extraction is explicitly scoped and reviewed.

---

## Outcome

- Completed content-level review of all 7 files in `shared_standards/public/standards/zacatl/**`.
- Wrote evidence and per-file decisions in `.hawp/work/evidence/2026/05/15/TASK-061-zacatl-boundary-review.md`.
- Boundary decisions:
	- `README.md` => `private/proprietary` (`split-private`) routed to `TASK-060`
	- remaining 6 standards pages => `repo-specific` (`adapt`) routed to `TASK-062`
- No file from this lane was promoted as direct `public-safe` absorb in current form.

## Verification

Direct evidence:

- Repo-root proof captured pre-edit:
	- `pwd` => `<repo-root-abs>`
	- `git rev-parse --show-toplevel` => `<repo-root-abs>`
	- `git rev-parse --show-prefix` => `(empty)`
- Lexical red-flag scan executed on Zacatl files:
	- `rg -n "tekit|micltan|mictlan|zacatl|internal|private|proprietary|topology|runbook|endpoint|credential" shared_standards/public/standards/zacatl/*.md`
- Per-file decisions and evidence snippets are recorded in `.hawp/work/evidence/2026/05/15/TASK-061-zacatl-boundary-review.md`.

Unproven:

- Generalized extraction quality is pending implementation in `TASK-062`.

## Close Checklist

- [x] Outcome section filled
- [x] Verification section filled (all claims have direct evidence or "unproven" tag)
- [x] Evidence files created if large/complex
- [x] Plan file moved to closed/YYYY/MM/DD/
- [x] BACKLOG.md updated
- [x] Status report written (not required)
- [x] Decision file created (not required)
