# Bug / Task: expose backlog alignment as agent-facing prompt and instruction

**Backlog ID:** TASK-007
**Type:** improvement
**Reported:** 2026-05-02
**Risk Level:** low

---

### Input (what was reported)

> Expose compact backlog alignment as agent-facing instruction/prompt in root and core overlays, wire existing prompts/instructions to it, and keep core/BACKLOG_ALIGNMENT.md as canonical policy.

---

### Context

The canonical policy exists in `core/BACKLOG_ALIGNMENT.md`, but reusable prompt/instruction files for agent runs are missing. Existing intake/status guidance has partial compact-backlog hints but no dedicated reusable backlog-alignment artifacts.

---

### Analysis

**Root cause (or most likely cause):**
Policy exists only as canonical documentation, not as a dedicated always-on instruction plus reusable task prompt in the overlay packs most agents load.

**Directly verified:**

- `core/BACKLOG_ALIGNMENT.md` exists and defines compact active-index rules.
- `.github/instructions/` and `.github/prompts/` do not yet include `hawp-backlog-alignment.*` files.
- `core/.github/instructions/` and `core/.github/prompts/` also do not yet include those files.
- Existing intake/copilot files mention compact backlog behavior but do not reference a dedicated backlog-alignment instruction/prompt artifact.

**Inferred (not yet proven):**
Without dedicated overlay artifacts, agent behavior may regress to appending all completed work in long `Done` sections.

**Scope — what else is affected:**

- Root overlay files under `.github/`.
- Core overlay source files under `core/.github/`.
- Installation/update docs listing shipped prompts/instructions.

---

### Work Coordination

**Owner:** agent
**Implementation status:** done
**Overlapping files:**

- `.github/copilot-instructions.md`
- `.github/instructions/intake.instructions.md`
- `.github/prompts/intake.prompt.md`
- `.github/prompts/hawp-status-report.prompt.md`
- `core/.github/copilot-instructions.md`
- `core/.github/instructions/intake.instructions.md`
- `core/.github/prompts/intake.prompt.md`
- `core/.github/prompts/hawp-status-report.prompt.md`
- `core/install.md`
- `core/update.md`

**Parallel work risk:** low
**Can implement now:** yes

**Coordination note:**
Current active work in this repo does not list overlapping items on these files.

---

### Options

#### Option A — Add dedicated instruction and prompt in both root/core overlays

Create the two new files in both overlays and wire references from existing prompts/instructions/docs.

#### Option B — Only update existing intake/copilot files

Avoid new files, but this weakens reuse and discoverability.

---

### Recommended Fix

**Option chosen:** A
**Rationale:**
This matches the requested layering and keeps canonical policy centralized in `core/BACKLOG_ALIGNMENT.md` while exposing concise operational guidance to agents.

**Files to change:**

- `.github/instructions/hawp-backlog-alignment.instructions.md`
- `.github/prompts/hawp-backlog-alignment.prompt.md`
- `core/.github/instructions/hawp-backlog-alignment.instructions.md`
- `core/.github/prompts/hawp-backlog-alignment.prompt.md`
- Existing root/core copilot/intake/status prompt files for references
- `core/install.md` and `core/update.md` for shipped-files references

**What to verify after:**

- [ ] New instruction exists in root and core overlays.
- [ ] New prompt exists in root and core overlays with review/compact modes.
- [ ] Existing intake/status/copilot surfaces reference the new artifacts.
- [ ] Install/update docs mention and ship the new files.
- [ ] Canonical policy remains `core/BACKLOG_ALIGNMENT.md`.

---

### Implementation Notes

Do not compact the live `.hawp/work/BACKLOG.md` in this task. Keep historical records intact.

### Outcome

Added dedicated backlog-alignment instruction and prompt files in both root and core overlays, wired references from intake/status/copilot guidance, and updated install/update documentation so these files are shipped and validated during refresh.

### Verification

- Confirmed new files exist:
	- `.github/instructions/hawp-backlog-alignment.instructions.md`
	- `.github/prompts/hawp-backlog-alignment.prompt.md`
	- `core/.github/instructions/hawp-backlog-alignment.instructions.md`
	- `core/.github/prompts/hawp-backlog-alignment.prompt.md`
- Confirmed root/core intake and status prompts now reference backlog alignment artifacts.
- Confirmed `core/install.md` and `core/update.md` include the new overlay files in shipped/checklist expectations.

---

### Status

- [x] Plan written
- [x] Approved / auto-approved (low risk)
- [x] Implemented
- [x] Verified
- [x] Closed
