# Bug / Task: Harden intake workflow against machine-local path leakage

**Backlog ID:** TASK-035
**Type:** task
**Reported:** 2026-05-13
**Risk Level:** low

---

### Input (what was reported)

> update the core files if needed and then i can run the update-dev command

---

### Context

The user asked to remove machine-local absolute paths and requested core workflow updates so future update runs carry the same privacy guardrails.

---

### Analysis

**Root cause (or most likely cause):**
Current path-discipline guidance requires repo-root proof command outputs, but does not explicitly enforce redaction when writing those outputs into plans/evidence. This can leak machine-local prefixes.

**Directly verified:**

- Workflow guidance in `.hawp/kit/usage/intake-workflow.md`, `.hawp/kit/templates/intake-plan.md`, `.hawp/kit/start-here.md`, and `.github/instructions/intake.instructions.md` requires repo-root proof commands.
- Core source equivalents exist under `core/.hawp/kit/**` and `core/.github/**`, which feed distribution/update workflows.
- One remaining generic user-home absolute-path mention exists in `.hawp/work/closed/2026/05/13/TASK-034.md`.

**Inferred (not yet proven):**
Adding explicit redaction-safe rules in both repo-local and core source guidance should prevent recurrence across direct usage and future `update-dev` propagation.

**Scope — what else is affected:**

- `.hawp/work/BACKLOG.md`
- `.hawp/work/active/TASK-035.md`
- `.hawp/work/closed/2026/05/13/TASK-034.md`
- `.hawp/kit/start-here.md`
- `.hawp/kit/usage/intake-workflow.md`
- `.hawp/kit/templates/intake-plan.md`
- `.github/instructions/intake.instructions.md`
- `.github/prompts/intake.prompt.md`
- `core/.hawp/kit/start-here.md`
- `core/.hawp/kit/usage/intake-workflow.md`
- `core/.hawp/kit/templates/intake-plan.md`
- `core/.github/instructions/intake.instructions.md`
- `core/.github/prompts/intake.prompt.md`
- `core/distribution/sources/shared/safety.md`
- `core/distribution/sources/shared/update.md`

---

### Work Coordination

**Owner:** agent
**Implementation status:** done
**Overlapping files:**

- `.hawp/work/BACKLOG.md`
- `.hawp/work/active/TASK-035.md`
- `.hawp/work/closed/2026/05/13/TASK-034.md`
- `.hawp/kit/start-here.md`
- `.hawp/kit/usage/intake-workflow.md`
- `.hawp/kit/templates/intake-plan.md`
- `.github/instructions/intake.instructions.md`
- `.github/prompts/intake.prompt.md`
- `core/.hawp/kit/start-here.md`
- `core/.hawp/kit/usage/intake-workflow.md`
- `core/.hawp/kit/templates/intake-plan.md`
- `core/.github/instructions/intake.instructions.md`
- `core/.github/prompts/intake.prompt.md`
- `core/distribution/sources/shared/safety.md`
- `core/distribution/sources/shared/update.md`

**Parallel work risk:** low
**Can implement now:** yes

**Coordination note:**
Changes are constrained to privacy and workflow guidance plus this work item. Unrelated lanes remain untouched.

**Path discipline:**

- Use exact repo-relative paths for file references.
- Do not persist machine-local absolute paths in work artifacts.
- If repo-root proof commands output absolute paths, redact only the machine-local prefix in saved artifacts (for example `<repo-root-abs>`), while keeping command identity and relative-path evidence intact.

**Repo-root proof (required before path-sensitive edits):**

```bash
pwd
git rev-parse --show-toplevel
git rev-parse --show-prefix
git status --short
```

Output logged for this task (machine-local prefix redacted):

- `pwd` -> `<repo-root-abs>`
- `git rev-parse --show-toplevel` -> `<repo-root-abs>`
- `git rev-parse --show-prefix` -> `` (repo root)
- `git status --short` -> `M .hawp/work/BACKLOG.md`

---

### Options

#### Option A — Local-only hardening

Update only repo-local guidance files and sanitize current leaks. Lower touch count but does not propagate via `update-dev` to downstream consumers.

#### Option B — Local + core source hardening

Update local guidance and corresponding `core/` source files, then sanitize remaining historical leak references.

---

### Recommended Fix

**Option chosen:** B
**Rationale:**
Satisfies immediate privacy requirements and ensures future update flows inherit the same guardrails.

**Files to change:**

- See scoped list above.

**What to verify after:**

- [ ] No machine-local absolute user-home path prefixes remain in `.hawp/work/**`
- [ ] Intake/workflow docs contain explicit redaction-safe guidance
- [ ] Core source docs/prompts include the same guardrails
- [ ] Only TASK-035 scoped files changed

---

### Implementation Notes

Keep semantics stable; make targeted wording updates only.

---

## Outcome (filled at close)

- Added redaction-safe path handling rules to repo-local intake guidance in `.hawp/kit/start-here.md`, `.hawp/kit/usage/intake-workflow.md`, `.hawp/kit/templates/intake-plan.md`, `.github/instructions/intake.instructions.md`, and `.github/prompts/intake.prompt.md`.
- Mirrored the same guardrails to core source files used by update distribution: `core/.hawp/kit/start-here.md`, `core/.hawp/kit/usage/intake-workflow.md`, `core/.hawp/kit/templates/intake-plan.md`, `core/.github/instructions/intake.instructions.md`, and `core/.github/prompts/intake.prompt.md`.
- Added privacy-safe guidance to core shared distribution docs in `core/distribution/sources/shared/safety.md` and `core/distribution/sources/shared/update.md`.
- Sanitized the remaining generic user-home absolute-path mention in `.hawp/work/closed/2026/05/13/TASK-034.md`.
- Created and executed TASK-035 as a new intake item in `.hawp/work/BACKLOG.md`.

---

## Verification (filled at close)

- [x] Claim 1: No exact personal path prefix remains in repository. **Evidence:** `rg -n --hidden --glob '!.git' '<redacted-user-prefix>' .` returned no matches when executed with the user-provided absolute prefix.
- [x] Claim 2: No literal macOS user-home markers remain in scoped workflow/work/core artifacts. **Evidence:** scoped user-home marker scan returned no matches.
- [x] Claim 3: Intake and core update-source guidance now includes explicit redaction-safe logging policy. **Evidence:** Updated files listed in Outcome include placeholder-based redaction requirements.
- [x] Claim 4: Path-sensitive proof captured with scoped changes only. **Evidence:** `git diff --name-status`, `git diff --check`, `git diff --cached --name-status`, `git diff --cached --check`, and `git status --short` were run; only TASK-035 scoped files are changed and no diff-check errors were reported.

---

## Close Checklist

- [x] Outcome section filled
- [x] Verification section filled
- [x] Staged-path proof captured before commit
- [x] Only approved TASK-035 paths are staged
- [ ] Commit created with focused message (not requested in this task)

**Status:**

- [x] Plan written
- [x] Approved / auto-approved (low risk)
- [x] Implemented
- [x] Verified
- [x] Closed
