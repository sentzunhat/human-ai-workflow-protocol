## Task: Define private/opinionated standards exclusion policy

**Backlog ID:** TASK-056
**Type:** governance
**Reported:** 2026-05-15
**Risk Level:** low
**Status:** done

---

### Input (what was reported)

> and the splitting of the private and opinionated should be private for the work project or no screct sauce reciepy or arch of tekit or micltan please be very careful and do multiple plan revie work items

---

### Context

The standards absorption stream needs an explicit policy boundary so only public-safe standards are copied into `core/.hawp/kit/standards/**`. Private or proprietary material must remain in work-scoped planning lanes.

---

### Analysis

**Root cause:**

No single policy file currently states a hard exclusion rule for private/opinionated content during standards absorption.

**Directly verified:**

- `core/.hawp/kit/standards/README.md` defines where public-safe standards belong.
- TASK-054 is the absorption coordination item.

**Inferred (not yet proven):**

- A written exclusion policy will reduce accidental leakage during future standards sync passes.

**Scope — what else is affected:**

- `core/.hawp/kit/standards/README.md`
- `.hawp/work/active/TASK-054.md`
- future standards review tasks

---

### Work Coordination

**Owner:** agent
**Implementation status:** in-progress
**Overlapping files:**

- `core/.hawp/kit/standards/README.md`
- `.hawp/work/active/TASK-054.md`

**Parallel work risk:** low
**Can implement now:** yes, after review

---

### Recommended Fix

1. Add a short exclusion policy section under `core/.hawp/kit/standards/README.md`.
2. Include explicit do-not-absorb examples: secret sauce recipes, internal architecture playbooks, and internal domain-specific patterns (for example Tekit/Micltan/Zacatl internals).
3. Link the exclusion policy from TASK-054.
4. Validate links and run `./.hawp/bin/hawp backlog validate`.

**What to verify after:**

- [x] Exclusion policy exists in standards index.
- [x] TASK-054 references the policy.
- [x] Validation passes with no FAIL checks.

---

## Outcome

- Added `Exclusion Policy (Private and Proprietary Content)` to `core/.hawp/kit/standards/README.md` with hard do-not-absorb rules for private/proprietary material.
- Added an explicit policy-link requirement to `.hawp/work/active/TASK-054.md` so each absorption pass checks the canonical exclusion policy.
- Kept scope narrow to governance and documentation policy only.

## Verification

Direct evidence captured:

- Repo-root proof commands run before edits:
	- `pwd` => `<repo-root-abs>`
	- `git rev-parse --show-toplevel` => `<repo-root-abs>`
	- `git rev-parse --show-prefix` => `(empty)`
- Link validation:
	- `BROKEN_LINKS=0` for `core/.hawp/kit/standards/README.md`, `.hawp/work/active/TASK-054.md`, `.hawp/work/active/TASK-056.md`, `.hawp/work/BACKLOG.md`
- Workflow validation:
	- `./.hawp/bin/hawp backlog validate` => `Result: PASS` and `Both kit/work validators completed without FAIL checks.`

Unproven:

- Future reviewer behavior consistency is not yet proven until TASK-057 rubric adoption and TASK-058 triage execution are completed.

## Close Checklist

- [x] Outcome section filled
- [x] Verification section filled (all claims have direct evidence or "unproven" tag)
- [x] Evidence files created if large/complex (not required for this small policy update)
- [x] Plan file moved to closed/YYYY/MM/DD/
- [x] BACKLOG.md updated
- [x] Status report written (not required; low-risk and fully verified)
- [x] Decision file created (not required)
