## Task: Build private-content review rubric and red-flag lexicon

**Backlog ID:** TASK-057
**Type:** governance
**Reported:** 2026-05-15
**Risk Level:** medium
**Status:** done

---

### Input (what was reported)

> please be very careful and do multiple plan revie work items

---

### Context

TASK-054 requires a repeatable way to classify each shared standard as public-safe, repo-specific, or private/proprietary.

---

### Analysis

**Root cause:**

Classification criteria are currently implicit and can vary by reviewer.

**Directly verified:**

- TASK-054 requires per-entry classification and split follow-ups.
- No formal rubric or red-flag lexicon exists yet in `.hawp/work/active/`.

**Inferred (not yet proven):**

- A rubric with explicit red flags will improve review consistency and reduce accidental leaks.

**Scope — what else is affected:**

- `.hawp/work/active/TASK-054.md`
- standards absorption commits under `core/.hawp/kit/standards/**`

---

### Work Coordination

**Owner:** agent
**Implementation status:** done
**Overlapping files:**

- `.hawp/work/active/TASK-054.md`
- optional reference note under `.hawp/work/evidence/`

**Parallel work risk:** medium
**Can implement now:** yes (approved by continuation request)

---

### Recommended Fix

1. Define a review rubric with three labels:
   - `public-safe`
   - `repo-specific`
   - `private/proprietary`
2. Define red-flag terms and patterns that force `private/proprietary` unless explicitly cleared.
3. Add a mandatory reviewer note format:
   - evidence (direct quote/snippet)
   - classification
   - rationale
   - follow-up destination
4. Add a quick checklist for final merge gating.

**What to verify after:**

- [x] Rubric exists and is referenced by TASK-054.
- [x] Red-flag lexicon exists and is usable for scans.
- [x] At least one sample classification is documented.

---

## Outcome

- Added `.hawp/kit/guidance/shared-standards-review-rubric.md` with:
   - three-label classifier (`public-safe`, `repo-specific`, `private/proprietary`)
   - hard red-flag lexicon and pattern-based red flags
   - required reviewer note format and merge gate checklist
   - sample classification entry
- Updated `.hawp/work/active/TASK-054.md` to require this rubric for all shared-standards intake entries.

## Verification

Direct evidence:

- Repo-root proof captured pre-edit:
   - `pwd` => `<repo-root-abs>`
   - `git rev-parse --show-toplevel` => `<repo-root-abs>`
   - `git rev-parse --show-prefix` => `(empty)`
- Rubric artifact created at `.hawp/kit/guidance/shared-standards-review-rubric.md`.
- TASK-054 references rubric usage in Step 1 gate.

Unproven:

- End-to-end classification consistency across a full shared-standards intake run remains to be proven in TASK-058.

## Close Checklist

- [x] Outcome section filled
- [x] Verification section filled (all claims have direct evidence or "unproven" tag)
- [x] Evidence files created if large/complex (not required)
- [x] Plan file moved to closed/YYYY/MM/DD/
- [x] BACKLOG.md updated
- [x] Status report written (not required)
- [x] Decision file created (not required)
