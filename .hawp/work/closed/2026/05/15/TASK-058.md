## Task: Run shared-standards triage and split follow-up implementation

**Backlog ID:** TASK-058
**Type:** review
**Reported:** 2026-05-15
**Risk Level:** medium
**Status:** done

---

### Input (what was reported)

> ...review the public and private standards and take the best public or do a diff of what we have and what the stands have and then absorpt the standards that are public and keep it simple

---

### Context

After TASK-056 and TASK-057 are approved, this item executes the actual triage run over shared standards input and emits split implementation tasks.

---

### Analysis

**Root cause:**

The shared standards review has not been run yet because source input and review gates were not fully prepared.

**Directly verified:**

- TASK-054 requests shared standards paste and classification.
- No completed triage artifact exists yet for current intake.

**Inferred (not yet proven):**

- Splitting triage execution into its own item will keep TASK-054 bounded and easier to audit.

**Scope — what else is affected:**

- `core/.hawp/kit/standards/**`
- `.hawp/work/active/TASK-054.md`
- new follow-up task files for private/opinionated items

---

### Work Coordination

**Owner:** agent
**Implementation status:** done
**Overlapping files:**

- `.hawp/work/active/TASK-054.md`
- `.hawp/work/BACKLOG.md`
- `core/.hawp/kit/standards/**` (future)

**Parallel work risk:** medium
**Can implement now:** yes (shared standards input exists in workspace)

---

### Recommended Fix

1. Take shared standards source input and normalize entries.
2. Apply TASK-057 rubric to each entry.
3. Generate a triage table with:
   - source path/section
   - classification
   - action (`absorb`, `adapt`, `split-private`)
4. For each `split-private` entry, create a separate backlog/plan item.
5. Implement only `absorb` items that are public-safe and approved.

**What to verify after:**

- [x] Triage table covers all input entries.
- [x] Every private/proprietary entry has a follow-up task.
- [x] Absorbed content excludes proprietary/internal architecture details.

---

## Outcome

- Produced full per-file shared standards triage artifact: `.hawp/work/evidence/2026/05/15/TASK-058-shared-standards-triage.md`.
- Coverage completed for all 55 files under `shared_standards/**`.
- Classification/action summary:
   - `absorb`: 10
   - `adapt`: 25
   - `split-private`: 20
- Created split follow-up tasks for non-absorb lanes:
   - `TASK-059` private standards lane retention
   - `TASK-060` project-specific boundary lane (Tekit/Mictlan/Zacatl)
   - `TASK-061` Zacatl public-lane red-flag boundary review
   - `TASK-062` repo-specific adaptation lane for docs/template overlap

## Verification

Direct evidence:

- Repo-root proof captured pre-edit:
   - `pwd` => `<repo-root-abs>`
   - `git rev-parse --show-toplevel` => `<repo-root-abs>`
   - `git rev-parse --show-prefix` => `(empty)`
- File inventory captured with `rg --files shared_standards | sort` and triaged in evidence file.
- Triage table includes all 55 entries with required fields: source path, classification, action, destination, rationale.

Unproven:

- Content-level redaction/extraction decisions for Zacatl docs remain pending TASK-061.
- Line-level adaptation outcomes for docs/template overlap remain pending TASK-062.

## Close Checklist

- [x] Outcome section filled
- [x] Verification section filled (all claims have direct evidence or "unproven" tag)
- [x] Evidence files created if large/complex
- [x] Plan file moved to closed/YYYY/MM/DD/
- [x] BACKLOG.md updated
- [x] Status report written (not required)
- [x] Decision file created (not required)
