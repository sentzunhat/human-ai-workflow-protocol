# HAWP Rule Precedence Table

**UUID:** `c3d4e5f6-a7b8-9012-cdef-123456789012`
**Type:** improvement
**Reported:** 2026-09-06
**Risk Level:** low

---

## Input (what was reported)

> Add a HAWP rule precedence table to `.github/copilot-instructions.md` so the agent knows which file's rules win when there's overlap across 8+ sources.

---

## Context

HAWP rules exist in 8+ places:

1. `.github/copilot-instructions.md` (Copilot-specific)
2. `AGENTS.md` (editor-agnostic entry point)
3. `CLAUDE.md` (Claude Code specific)
4. `.claude/rules/hawp-core.md`, `.claude/rules/hawp-backlog-alignment.md`, `.claude/rules/hawp-intake.md`, `.claude/rules/hawp-docs-alignment.md` (Claude Code scoped rules)
5. `.github/instructions/commit-style.instructions.md`, `.github/instructions/intake.instructions.md`, `.github/instructions/hawp-backlog-alignment.instructions.md`, `.github/instructions/hawp-docs-alignment.instructions.md` (instruction files)
6. `core/.hawp/work/BACKLOG.md` (downstream scaffold)
7. Repo-local `.hawp/work/BACKLOG.md` (active work index)

The agent may satisfy conflicting instructions, duplicate guidance across sources, or treat one source as authoritative when another is intended to be higher-precedence. This creates confusion especially when the agent is asked about "HAWP rules" without being told which repo's instance applies.

**Directly verified:**

- All 8+ files exist in this repository and contain overlapping HAWP workflow guidance.
- `AGENTS.md` says "Follow the repo-local HAWP guidance in `.hawp/kit/start-here.md`" but doesn't define precedence vs other instruction sources.
- `CLAUDE.md` references `@AGENTS.md` and adds Claude Code-specific commands.

**Inferred (not yet proven):**

- The agent may load different rule sets depending on which files it discovers in a conversation, leading to inconsistent behavior across tools (Copilot vs Cursor vs Claude Code).
- Shared behavioral rules in `core/providers/shared/behaviors/` are generated into provider-specific overlays, but the precedence between manual instructions and generated overlays isn't documented.

**Scope — what else is affected:**

- `.github/copilot-instructions.md` (add precedence table)
- Potentially `AGENTS.md` if cross-editor alignment improves there

---

## Analysis

**Root cause (or most likely cause):**
No explicit precedence hierarchy exists. All instruction files describe valid rules but without ordering, the agent can't resolve conflicts when two sources give different guidance on the same topic. For example, both `.github/copilot-instructions.md` and `AGENTS.md` reference BACKLOG.md compaction rules — if one says "5-10 items" and another says "last 14-30 days", which wins?

**Directly verified:**

- Multiple files contain the rule "Do not append completed work endlessly to BACKLOG.md." This is duplicated at minimum in: `.github/copilot-instructions.md`, `AGENTS.md`, `CLAUDE.md`, and potentially generated overlays.
- Compaction triggers also appear in both `.github/copilot-instructions.md` ("When BACKLOG.md has many Done rows...") and `.claude/rules/hawp-backlog-alignment.md`.

**Inferred (not yet proven):**

- The canonical source for shared HAWP behavior is `core/providers/shared/behaviors/`, with generated overlays in `distribution/generated/` and repo-local copies. But this chain isn't documented for the agent.
- Per-tool instruction files (`CLAUDE.md`, `.github/copilot-instructions.md`) should be the highest precedence for their respective tools since they're the first thing loaded.

**Scope — what else is affected:**

- Primarily `.github/copilot-instructions.md`. Consider updating `AGENTS.md` if useful as cross-editor reference.

---

## Verification

Implemented 2026-09-06: Added "Rule Precedence" section to `.github/copilot-instructions.md` defining 4-tier precedence hierarchy (current tool instructions > AGENTS.md > shared sources > per-tool overlays). Clarifies which source wins when HAWP rules conflict across 8+ locations. Verified `go run ./cmd/hawp check --no-update-check` passes all 3 validations (kit validate, work validate, links check). Plan closed via commit `92cb6ac`.

---

## Work Coordination

**Owner:** unassigned
**Implementation status:** done
**Parallel work risk:** low
**Can implement now:** yes
**Coordination note:** This and the other 3 items edit `.github/copilot-instructions.md`. Can batch all edits into a single implementation pass.

---

## Next Step

1. Read current `.github/copilot-instructions.md` (if not already loaded)
2. Add a Rule Precedence section at the top defining the 4-tier hierarchy:

   ```markdown
   ## Rule Precedence (HAWP instructions)

   When rules conflict, this order applies (highest → lowest):

   1. **Current tool's own instructions** — `.github/copilot-instructions.md` for Copilot, `CLAUDE.md` for Claude Code, etc.
   2. **Editor-agnostic entry point** — `AGENTS.md`
   3. **Shared HAWP behavior sources** — `core/providers/shared/behaviors/`, mirrored to repo-local instructions in `.github/instructions/` and `.claude/rules/`
   4. **Per-tool rule overlays** — `.claude/rules/*.md`, provider-generated files in `distribution/generated/`

   HAWP kit templates (`.hawp/kit/`) are always the canonical policy source; instruction files adapt them for specific tools without contradicting them.
   ```

3. Validate: run `go vet ./...` and `go test ./...` to confirm no regressions

---

## Outcome

- `.github/copilot-instructions.md` updated with "Rule Precedence" section defining 4-tier precedence hierarchy (current tool instructions > AGENTS.md > shared sources > per-tool overlays)
- Resolves ambiguity when HAWP rules conflict across 8+ instruction sources
- `go run ./cmd/hawp check --no-update-check` passes all 3 validations (kit validate, work validate, links check)

---

## Close Checklist

- [x] Plan written and reviewed
- [x] Implementation completed
- [x] Verification evidence recorded
- [x] Changes committed (`92cb6ac` on `feature/v0.0.24-work-folder-normalization`)

1. Read current `.github/copilot-instructions.md`, `AGENTS.md`, `CLAUDE.md`
2. Add a "Rule Precedence" section at the top:

   ```
   ## Rule Precedence (HAWP instructions)

   When rules conflict, this order applies (highest → lowest):
   1. **Current tool's own instructions** — `.github/copilot-instructions.md` for Copilot, `CLAUDE.md` for Claude Code, etc.
   2. **Editor-agnostic entry point** — `AGENTS.md`
   3. **Shared HAWP behavior sources** — `core/providers/shared/behaviors/`, mirrored to repo-local instructions in `.github/instructions/` and `.claude/rules/`
   4. **Per-tool rule overlays** — `.claude/rules/*.md`, provider-generated files in `distribution/generated/`

   HAWP kit templates (`.hawp/kit/`) are always the canonical policy source; instruction files adapt them for specific tools without contradicting them.
   ```

3. Remove any contradictory rules from lower-precedence sources (check `.claude/rules/` against what's now in Copilot instructions).
4. Validate HAWP kit: `go run ./cmd/hawp check --no-update-check`.

---

## Outcome

- `.github/copilot-instructions.md` updated with "Rule Precedence" section defining 4-tier hierarchy (current tool > AGENTS.md > shared sources > per-tool overlays)
- Clarifies which source wins when HAWP rules conflict across 8+ locations
- `go run ./cmd/hawp check --no-update-check` passes all 3 validations (kit validate, work validate, links check)

---

## Close Checklist

- [x] Plan written and reviewed
- [x] Implementation completed
- [x] Verification evidence recorded
- [x] Changes committed (`92cb6ac` on `feature/v0.0.24-work-folder-normalization`)
