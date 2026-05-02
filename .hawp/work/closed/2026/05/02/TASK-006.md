# Bug / Task: dogfood source-repo work state under repo-root .hawp/work

**Backlog ID:** TASK-006
**Type:** improvement
**Reported:** 2026-05-02
**Risk Level:** medium

---

### Input (what was reported)

> Title: Dogfood HAWP in the HAWP repo using repo-root .hawp/work
>
> Mission: move this source repo operating state from root .work/ to repo-root .hawp/work/, preserve install/update trust boundary, create ADR first, and align README/install/update/instructions/prompts.

---

### Context

The repo currently uses `core/.hawp/` as the installable source package and has source-repo operating state in root `.work/`. This differs from downstream projects, which use repo-root `.hawp/work/`, and causes mental-model drift.

---

### Analysis

**Root cause (or most likely cause):**
Repo docs/instructions currently codify root `.work/` as source-repo operating state, while downstream guidance codifies `.hawp/work/`. The mismatch exists in README/update/copilot/intake guidance.

**Directly verified:**
- Root `.work/` exists and contains real backlog/active/closed/decisions/evidence/notes state.
- Root `.hawp/` is currently absent.
- `README.md` and `core/update.md` explicitly describe root `.work/` as source-repo operating state.
- `.github/copilot-instructions.md` and `.github/instructions/intake.instructions.md` still direct work tracking to `.work/`.

**Inferred (not yet proven):**
- Some historical links may require follow-up after folder migration if they assumed root `.work/`.

**Scope — what else is affected:**
- Source-repo docs and AI instruction/prompt guidance.
- Source-repo work tree location and migration notes.
- Backlog and active plan paths as they move from `.work/` to `.hawp/work/`.

---

### Work Coordination

**Owner:** agent
**Implementation status:** done
**Overlapping files:**

- `README.md`
- `core/install.md`
- `core/update.md`
- `.github/copilot-instructions.md`
- `.github/instructions/intake.instructions.md`
- `.github/prompts/intake.prompt.md`
- `.hawp/work/BACKLOG.md`

**Parallel work risk:** medium
**Can implement now:** yes

**Coordination note:**
TASK-005 is not currently listed as active in the observed backlog. This task stays scoped to repo-root work-location dogfooding and trust-boundary alignment, not backlog compaction policy.

---

### Options

#### Option A — Keep root .work/

Preserves current source-repo state location and avoids migration churn. Trade-off: keeps source/downstream layout mismatch and mental-model drift.

#### Option B — Move source repo state to repo-root .hawp/work/ while keeping core/.hawp as package source

Aligns source repo with downstream shape while preserving the trust boundary between source package (`core/.hawp/*`) and source-repo operating state (`.hawp/work/*`).

#### Option C — Move package source to repo-root .hawp/

Would collapse source/runtime boundaries and require architecture changes to install/update assumptions.

---

### Recommended Fix

**Option chosen:** B
**Rationale:**
Aligns dogfooding behavior with downstream usage while preserving install/update packaging boundaries and avoiding architecture changes.

**Files to change:**

- `.hawp/work/**` — create and migrate source-repo work state from root `.work/`
- `README.md` — update boundary explanation and source layout
- `core/install.md` — clarify source-repo root `.hawp/work/` never ships; install seeds only from `core/.hawp/work/`
- `core/update.md` — same trust-boundary clarification
- `.github/copilot-instructions.md` — redirect source-repo task tracking to `.hawp/work/`
- `.github/instructions/intake.instructions.md` and `.github/prompts/intake.prompt.md` — replace root `.work/` guidance
- `.hawp/work/decisions/2026/05/02/adr-root-hawp-work-for-source-repo.md` — decision record required before implementation
- `.hawp/work/status/2026/05/02/*` or `.hawp/work/evidence/2026/05/02/*` — migration note

**What to verify after:**

- [x] No stale root `.work/` guidance remains in targeted docs/instructions/prompts.
- [x] Root `.hawp/work/` contains preserved source-repo work records.
- [x] `core/.hawp/work/` remains clean scaffold content.
- [x] `core/install.md` and `core/update.md` explicitly preserve project-owned target `.hawp/work/`.

---

### Implementation Notes

Perform ADR first, then migration move, then documentation/instruction alignment, then stale-reference validation and link-follow-up summary.

---

### Status

- [x] Plan written
- [x] Approved / auto-approved (low risk)
- [x] Implemented
- [x] Verified
- [x] Closed
