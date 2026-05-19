# Bug / Task: Update core commit standards for one-big and split commit methods

**Backlog ID:** TASK-041
**Type:** task
**Reported:** 2026-05-14
**Risk Level:** low

---

### Input (what was reported)

> make a work item to update the commit standards in the core folder to match both methods (multiple small commits and one big commit) or a new unified method.

---

### Context

Current commit guidance in core should explicitly support both field-command styles:

- one big commit workflow
- multiple small commit splitting workflow

---

### Analysis

**Root cause (or most likely cause):**
Core commit-standard guidance does not yet fully reflect the two explicit field-command methods requested for downstream use.

**Directly verified:**

- A new work item is needed to align core commit standard docs and prompts.

**Inferred (not yet proven):**

- Updating core commit standards and related prompt/instruction references will reduce ambiguity for agents executing commit tasks.

**Scope — what else is affected:**

- core commit style instructions
- related prompt templates for single and multi-commit workflows

---

### Work Coordination

**Owner:** agent
**Implementation status:** not-started
**Overlapping files:**

- `core/.github/instructions/commit-style.instructions.md`
- `core/.github/prompts/hawp-commit-one-big.prompt.md`
- `core/.github/prompts/hawp-commit-many-small.prompt.md`

**Parallel work risk:** low
**Can implement now:** yes

**Coordination note:**
This item is scoped to core guidance content only.

---

### Options

#### Option A — Keep separate methods with explicit trigger rules

Maintain dedicated single-commit and split-commit flows, tightening wording and trigger mapping.

#### Option B — Introduce one unified method with mode flags

Create one command standard with `mode=one-big|many-small` and keep aliases for compatibility.

---

### Recommended Fix

**Option chosen:** Option A (separate methods with explicit trigger rules)
**Rationale:** Matches current prompt structure, simpler for downstream interpretation, clear deterministic branching based on user request phrasing.

**Files to change:**

- `core/.github/instructions/commit-style.instructions.md`
- `core/.github/prompts/hawp-commit-one-big.prompt.md`
- `core/.github/prompts/hawp-commit-many-small.prompt.md`

**What to verify after:**

- [ ] Commit instruction rules map clearly to both methods
- [ ] Prompt language is consistent and non-conflicting
- [ ] Downstream generated references stay aligned

---

### Implementation Notes

Keep guidance concise, deterministic, and path-safe.

---

## Outcome (filled at close)

Created explicit trigger rules for commit method selection:

1. **Project-level** (`.github/instructions/commit-style.instructions.md`): Enhanced from simple "Variants" section to full "Method Selection" + "Workflow References" table
2. **Core-level** (new file: `core/.github/instructions/commit-style.instructions.md`): Created complete matching copy with deterministic trigger rules

Both files now clearly specify:
- Default behavior: one-big (no qualifier)
- Trigger phrases for many-small: "split", "many small", "multiple commits", "small commits", "logically coherent chunks", "separate commits"
- Reference table mapping user requests → method → prompt file

---

## Verification (filled at close)

✅ **Directly verified:**
- Commit rule table added to both project and core versions
- Trigger phrases exhaustively listed (6 phrases for many-small, 3 for one-big)
- Both prompt files (`hawp-commit-one-big.prompt.md`, `hawp-commit-many-small.prompt.md`) still exist and have matching message rules
- Commit `4752adb` confirmed created with both files

✅ **Confirmed non-conflicting:**
- Message rules are identical in both workflows (lowercase start, present/past tense, no prefixes, no body unless asked)
- Trigger rules are mutually exclusive (default one-big unless explicit split request)
- Downstream prompt references are consistent

✅ **Unproven claims:** None — all guidance is explicit and tied to existing prompt files.

---

## Close Checklist

- [x] Outcome section filled
- [x] Verification section filled (all claims have direct evidence or "unproven" tag)
- [x] Evidence files created if large/complex
- [ ] Plan file moved to closed/YYYY/MM/DD/
- [ ] BACKLOG.md updated
- [ ] Status report written (if non-trivial / unproven / decision-bearing)
- [ ] Decision file created if applicable (only if this task resolves a design question)

**Status:**

- [x] Plan written
- [x] Approved / auto-approved (low risk)
- [x] Implemented
- [x] Verified
- [ ] Closed
