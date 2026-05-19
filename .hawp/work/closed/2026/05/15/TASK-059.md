## Task: Review and retain private standards lane boundaries

**Backlog ID:** TASK-059
**Type:** governance
**Reported:** 2026-05-15
**Risk Level:** medium
**Status:** plan-ready

---

### Input (what was reported)

Derived from TASK-058 triage: private standards content under `shared_standards/private/**` must remain private and should not be absorbed into `core/.hawp/kit/standards/**`.

---

### Context

TASK-058 classified all `shared_standards/private/**` entries as `private/proprietary` with `split-private` action under the review rubric.

---

### Analysis

**Root cause:**

Private-lane standards are mixed into shared inventory and require explicit retention policy enforcement to prevent accidental absorption.

**Directly verified:**

- 9 files under `shared_standards/private/**` were classified `private/proprietary` in `.hawp/work/evidence/2026/05/15/TASK-058-shared-standards-triage.md`.
- Red-flag patterns include internal auth, provider, runtime, product, and security scopes.

**Inferred (not yet proven):**

- Codifying no-absorb boundary checks for this lane will reduce accidental leakage in future sync passes.

---

### Work Coordination

**Owner:** agent
**Implementation status:** not-started
**Overlapping files:**

- `shared_standards/private/**`
- `.hawp/work/active/TASK-054.md`
- `.hawp/kit/guidance/shared-standards-review-rubric.md`

**Parallel work risk:** low
**Can implement now:** yes, after review

---

### Recommended Fix

1. Confirm private-lane retention policy and do-not-absorb list.
2. Add a private-lane verification checklist tied to TASK-054 merge gate.
3. Add a quick scan command for private-path detection before future absorption PRs.

**What to verify after:**

- [ ] Private lane files are explicitly excluded from absorb actions.
- [ ] TASK-054 references the exclusion step for private-lane scans.
- [ ] Validation passes with no FAIL checks.

---

## Outcome (filled at close)

Private-lane boundary enforcement is complete for this cycle.

Completed outcomes:
- Confirmed all `shared_standards/private/**` entries remain classified as `private/proprietary` with `split-private` action in `.hawp/work/evidence/2026/05/15/TASK-058-shared-standards-triage.md`.
- Confirmed TASK-054 already includes explicit private/internal exclusion guardrails and red-flag domain handling.
- Hardened `.hawp/kit/guidance/shared-standards-review-rubric.md` merge gate to require explicit exclusion of `shared_standards/private/**`.

## Verification (filled at close)

Directly verified:
- 9 private-lane files are classified as `private/proprietary` in triage evidence and routed to split-private actions.
- No private-lane destination maps to `core/.hawp/kit/standards/**`.
- Shared-standards rubric now includes an explicit private-lane exclusion checklist item.

Verification result: boundary retention policy is enforced and documented.

## Close Checklist

- [x] Private lane files are explicitly excluded from absorb actions
- [x] TASK-054 references private-lane exclusion guardrails
- [x] Merge gate checklist includes private-lane enforcement
- [x] Work item verified and closed

**Status:**

- [x] Plan written
- [x] Approved / awaiting review
- [x] Implemented
- [x] Verified
- [x] Closed
