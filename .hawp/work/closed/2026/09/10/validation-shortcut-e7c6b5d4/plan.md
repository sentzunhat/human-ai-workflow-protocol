# Move HAWP validation commands into agent instructions for Copilot

**UUID:** `validation-shortcut-e7c6b5d4`
**Type:** improvement
**Reported:** 2026-09-10
**Status:** `done`

---

## Input (what was reported)

Validation commands currently live in `CLAUDE.md` as a convenience section for Claude Code.
Copilot and other agents do not have equivalent access, causing redundant lookups of the same
command set during architecture validation.

## Context

Key HAWP validation commands:

```bash
go test ./...
go run ./cmd/hawp distribution sync
go run ./cmd/hawp check
```

These are referenced across multiple instruction files but centralized only in CLAUDE.md.

## Analysis

**Root cause:** Commands defined in Claude-specific docs; no shared reference point for Copilot.
Agents need to cross-reference `.hawp/kit/start-here.md`, `CLAUDE.md`, and `.github/copilot-instructions.md`
to find the validation command set.

**Directly verified:**

- `CLAUDE.md` has a "Commands" section with full validation set.
- `.github/copilot-instructions.md` references HAWP guidance but does not include the commands.
- `AGENTS.md` does not list validation commands.

**Likely fix:**

1. Create `.github/skills/validation-shortcuts.md` (or add to existing instructions) with the
   canonical validation command set.
2. Reference this skill from `.github/copilot-instructions.md`, `AGENTS.md`, and any other
   agent-facing docs.
3. Remove duplication — CLAUDE.md remains source-of-truth, all others reference it.

**Scope — what else is affected:**

- `.github/copilot-instructions.md` — add validation commands or skill reference
- `AGENTS.md` — same
- Consider creating `.github/skills/` directory if needed

## Work Coordination

**Owner:** Codex
**Implementation status:** done (2026-09-10)
**Parallel work risk:** low — does not overlap with any in-progress plan.
**Can implement now:** yes

---

## Completion Record (2026-09-10)

The canonical validation command block was folded into
`.github/copilot-instructions.md` as part of the consolidated Copilot
agent-instruction update. The same block has now been propagated to
`core/providers/.github/copilot-instructions.md` so GitHub provider updates can
carry the guidance to downstream repositories.

## Outcome

Copilot-facing HAWP instructions now include the canonical validation command
and grouped validation commands, with `CLAUDE.md` named as the command source of
truth.

## Verification

Directly checked `.github/copilot-instructions.md` and
`core/providers/.github/copilot-instructions.md` for the validation command
block.

## Close Checklist

- [x] Validation commands added to Copilot instructions.
- [x] Provider source pack aligned.
- [x] Plan moved to `closed/2026/09/10/validation-shortcut-e7c6b5d4/`.
