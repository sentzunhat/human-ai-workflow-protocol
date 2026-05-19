## Task: Handle project-specific standards boundaries (Tekit/Mictlan/Zacatl)

**Backlog ID:** TASK-060
**Type:** governance
**Reported:** 2026-05-15
**Risk Level:** medium
**Status:** plan-ready

---

### Input (what was reported)

Derived from TASK-058 triage: `shared_standards/project-specific/**` paths contain internal domain-specific boundary pointers and must remain outside core public-safe standards.

---

### Context

The user explicitly requested caution around private/opinionated material and internal domains. TASK-058 classified project-specific entries as `private/proprietary` by default with `split-private` action.

---

### Analysis

**Root cause:**

Project-specific domain pointers are present in the shared standards tree but are not portable standards assets.

**Directly verified:**

- 4 files under `shared_standards/project-specific/**` were triaged as `private/proprietary` in `.hawp/work/evidence/2026/05/15/TASK-058-shared-standards-triage.md`.

**Inferred (not yet proven):**

- Maintaining this lane as private-only avoids accidental publication of domain-specific architecture guidance.

---

### Work Coordination

**Owner:** agent
**Implementation status:** not-started
**Overlapping files:**

- `shared_standards/project-specific/**`
- `.hawp/work/active/TASK-054.md`

**Parallel work risk:** low
**Can implement now:** yes, after review

---

### Recommended Fix

1. Define explicit keep-private handling for project-specific boundary pointers.
2. Add a reviewer checklist item requiring domain-name red-flag confirmation.
3. Ensure no project-specific files are mapped to core standards destinations.

**What to verify after:**

- [ ] Project-specific lane remains excluded from absorb actions.
- [ ] Red-flag domains are documented and enforced.
- [ ] Validation passes with no FAIL checks.

---

## Outcome (filled at close)

Project-specific lane boundary handling is complete for this cycle.

Completed outcomes:
- Confirmed all `shared_standards/project-specific/**` entries remain classified as `private/proprietary` with `split-private` action in `.hawp/work/evidence/2026/05/15/TASK-058-shared-standards-triage.md`.
- Confirmed project-specific domains (Tekit/Mictlan/Zacatl) remain default-private in TASK-054 boundary guidance.
- Hardened `.hawp/kit/guidance/shared-standards-review-rubric.md` merge gate to require explicit exclusion of `shared_standards/project-specific/**`.

## Verification (filled at close)

Directly verified:
- 4 project-specific files are classified as `private/proprietary` in triage evidence and routed to split-private actions.
- No project-specific destination maps to `core/.hawp/kit/standards/**`.
- Shared-standards rubric now includes an explicit project-specific lane exclusion checklist item.

Verification result: project-specific boundary policy is enforced and documented.

## Close Checklist

- [x] Project-specific lane remains excluded from absorb actions
- [x] Red-flag domains are documented and treated as private by default
- [x] Merge gate checklist includes project-specific lane enforcement
- [x] Work item verified and closed

**Status:**

- [x] Plan written
- [x] Approved / awaiting review
- [x] Implemented
- [x] Verified
- [x] Closed
