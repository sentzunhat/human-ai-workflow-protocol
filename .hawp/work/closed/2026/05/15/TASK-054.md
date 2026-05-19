## Task: Review shared standards and absorb public-safe guidance

**Backlog ID:** TASK-054
**Type:** task
**Reported:** 2026-05-15
**Risk Level:** low
**Status:** done

---

### Input (what was reported)

> i was thinkin for the folder to be added to the core hawp folder where it can be in standards and having those files in there with their opinionated decisions and have the mongo standards guide that we had before i will paste the shared_standards file and make another action item to plan the review of the shared stanads folder where it can review the public and private standards and take the best public or do a diff of what we have and what the stands have and then absorpt the standards that are public and keep it simple and if and pprivate stanadrds or zacatl standards that need to be added add them or make smaller action work items to plan how to integrate those pieces

---

### Context

The repository now has an opinionated database standards home under `core/.hawp/kit/standards/database/`, but the broader shared standards source still needs a bounded review. The next planning step is to compare the shared public standards against the repo's current standards surface and decide what to absorb, what to keep private, and what to split into smaller follow-up work items.

Hard boundary for this item: private or proprietary standards must stay out of `core/.hawp/kit/standards/**` and be handled as separate work-project artifacts.

### Analysis

**Root cause (or most likely cause):**
The repo has accumulated local standards and public kit standards independently, so a review is needed before more folders or guidance files are added.

**Directly verified:**

- `core/.hawp/kit/standards/` exists and already holds public absorbed standards.
- `core/.hawp/kit/standards/database/` now exists as the canonical database standards home.
- No `shared_standards` source file has been pasted into the workspace yet.

**Inferred (not yet proven):**

- A diff-first review of public vs private standards will make it easier to absorb only the portable parts and keep repository-specific opinions small.
- Any private or Zacatl-specific gap should be split into isolated follow-up work items instead of being merged into one large standards sweep.

**Privacy boundary (must enforce):**

- Do not absorb any secret sauce recipes, internal architecture playbooks, or proprietary implementation patterns into `core/.hawp/kit/standards/**`.
- Treat references to Tekit, Micltan, Zacatl, or other internal domains as private by default unless explicitly cleared for public-safe export.
- If uncertain, classify as private and open a separate review work item.
- Apply the canonical exclusion policy in `core/.hawp/kit/standards/README.md` before each absorption pass.

**Scope — what else is affected:**

- `core/.hawp/kit/standards/**`
- `shared_standards/**` once provided
- `.hawp/kit/guidance/da-schema-planning.md`
- related standards indexes and follow-up task files

---

### Work Coordination

**Owner:** agent
**Implementation status:** done
**Overlapping files:**

- `core/.hawp/kit/standards/**`
- `shared_standards/**`
- `.hawp/kit/guidance/`

**Parallel work risk:** medium
**Can implement now:** planning only, with review gate before absorption edits

**Coordination note:**
Keep this review bounded to public-safe absorption and gap classification. If the shared standards set is large, split the private or opinionated follow-ups into smaller tasks instead of widening this item.

Decomposed review items created:

- `TASK-056` — exclusion policy for private/opinionated standards
- `TASK-057` — review rubric and red-flag lexicon
- `TASK-058` — triage run and split follow-up implementation plan

---

### Recommended Fix

#### Step 1: Apply policy and rubric gates first

- [x] Confirm `TASK-056` policy wording is approved
- [x] Confirm `TASK-057` rubric and red-flag lexicon are approved
- [x] Use `.hawp/kit/guidance/shared-standards-review-rubric.md` as the mandatory classifier for all shared-standards intake entries

#### Step 2: Paste the shared_standards file below this line

---

#### Step 3: Review and classify each entry

- [x] Inventory all entries in the shared_standards file
- [x] Mark each as:
  - [x] public-safe (can be absorbed directly)
  - [x] repo-specific (needs adaptation)
  - [x] private/zacatl (requires separate follow-up)

#### Step 3: Absorb and update

- [x] Absorb public-safe MongoDB schema design guidance into `core/.hawp/kit/standards/database/mongodb-schema-design.md` (now complete)
- [x] Add reference in `core/.hawp/kit/standards/database/README.md` to absorbed MongoDB guidance
- [x] Continue review for any other public-safe standards

**Note:** The public MongoDB schema standards have been reviewed and absorbed into the core standards tree. No major conflicts or repo-specific overrides were required. The absorbed file is now referenced in the standards README. Remaining work: continue review for any other public-safe standards in the shared set.

- [x] Propose edits or new files for repo-specific items
- [x] Create new work items for any private/zacatl standards
- [x] Update planning guidance and references as needed

#### Step 5: Privacy safety check before merge

- [x] Run red-flag lexicon scan against candidate absorbed text
- [x] Confirm no private/internal architecture details are present
- [x] Confirm all uncertain items were split to follow-up tasks, not merged

---

**Next steps:**

1. Review/approve `TASK-056`, `TASK-057`, and `TASK-058` plan files.
2. Paste the shared_standards input for triage.
3. I will classify and diff each entry, then propose absorption or follow-up actions.
4. Only public-safe standards will be merged into the core tree; private/opinionated items stay in separate work-project tasks.

### Progress Update (2026-05-15, triage execution)

- TASK-058 completed with full per-file triage over `shared_standards/**`.
- Evidence artifact: `.hawp/work/evidence/2026/05/15/TASK-058-shared-standards-triage.md`.
- Split follow-up items created:
  - `TASK-059` for `shared_standards/private/**`
  - `TASK-060` for `shared_standards/project-specific/**`
  - `TASK-061` for `shared_standards/public/standards/zacatl/**` boundary review
  - `TASK-062` for repo-specific adaptation candidates

### Progress Update (2026-05-15, Zacatl boundary review)

- TASK-061 completed with content-level review of all `shared_standards/public/standards/zacatl/*.md` files.
- Evidence artifact: `.hawp/work/evidence/2026/05/15/TASK-061-zacatl-boundary-review.md`.
- Boundary decisions:
  - Zacatl `README.md` remains `split-private` (routed to `TASK-060`).
  - Six principle-level docs routed to `TASK-062` as `adapt` candidates (no direct absorb).

### Progress Update (2026-05-15, overlap adaptation lane)

- TASK-062 completed docs/template overlap diff review with no-content-change adaptation decision.
- Evidence artifact: `.hawp/work/evidence/2026/05/15/TASK-062-overlap-adaptation-review.md`.
- Split follow-up created: `TASK-063` for generalized extraction of six Zacatl principle docs.

### Progress Update (2026-05-15, generalized extraction lane)

- TASK-063 completed neutral extraction of six principle-level standards into `core/.hawp/kit/standards/service-design/**`.
- Evidence artifact: `.hawp/work/evidence/2026/05/15/TASK-063-generalized-extraction-review.md`.
- The extracted content avoids framework/domain labels and preserves boundary-safe exclusions.

---

### What to verify after

- [x] Shared standards content has been inventoried.
- [x] Public-safe items are identified.
- [x] Repo-specific/private items are split into follow-up tasks where needed.
- [x] The core standards tree reflects the approved absorbed guidance.
- [x] Planning guidance points at the final standards location.
- [x] No private/proprietary architecture details were absorbed into `core/.hawp/kit/standards/**`.

---

## Outcome

- Completed shared-standards umbrella review and absorption flow with explicit private-boundary enforcement.
- Completed all decomposed execution lanes (`TASK-056` through `TASK-063`) and preserved split-private handling for sensitive paths.
- Landed public-safe standards outcomes in `core/.hawp/kit/standards/**`, including database and generalized service-design standards.

## Verification

Direct evidence:

- Triage and boundary artifacts:
  - `.hawp/work/evidence/2026/05/15/TASK-058-shared-standards-triage.md`
  - `.hawp/work/evidence/2026/05/15/TASK-061-zacatl-boundary-review.md`
  - `.hawp/work/evidence/2026/05/15/TASK-062-overlap-adaptation-review.md`
  - `.hawp/work/evidence/2026/05/15/TASK-063-generalized-extraction-review.md`
- Standards index and destination tree updated under `core/.hawp/kit/standards/README.md`.
- Workflow validation command reports PASS with no FAIL checks:
  - `./.hawp/bin/hawp backlog validate`

Unproven:

- Long-term downstream reuse quality remains to be evaluated in future standards-review cycles.

## Close Checklist

- [x] Outcome section filled
- [x] Verification section filled (all claims have direct evidence or "unproven" tag)
- [x] Evidence files created if large/complex
- [x] Plan file moved to closed/YYYY/MM/DD/
- [x] BACKLOG.md updated
- [x] Status report written (not required)
- [x] Decision file created (not required)

---

## Outcome (filled at close)

Shared standards review and bounded absorption is complete.

Completed outcomes:

- Full triage executed across `shared_standards/**` with evidence captured in TASK-058 artifact.
- Public-safe standards absorbed into core standards tree.
- Repo-specific and private/proprietary lanes decomposed into dedicated follow-up tasks and completed through TASK-059 to TASK-063.
- Service-design generalized standards extracted into `core/.hawp/kit/standards/service-design/**` with red-flag domain terms removed.
- Merge-gate rubric hardened to enforce private and project-specific lane exclusions.

## Verification (filled at close)

Directly verified:

- Triage evidence exists: `.hawp/work/evidence/2026/05/15/TASK-058-shared-standards-triage.md`.
- Generalized extraction evidence exists: `.hawp/work/evidence/2026/05/15/TASK-063-generalized-extraction-review.md`.
- Red-flag scan over `core/.hawp/kit/standards/service-design/**` for `tekit|mictlan|micltan|zacatl` returns no matches.
- Backlog validation returns PASS.

Confidence: high. The shared-standards review objective is complete with policy boundaries preserved.

## Close Checklist

- [x] Plan written
- [x] Review/classification completed
- [x] Public-safe absorption completed
- [x] Private/repo-specific decomposition completed
- [x] Privacy/boundary checks completed
- [x] Verification completed
- [x] Closed
