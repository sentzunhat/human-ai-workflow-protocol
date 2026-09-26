# Structured checkpointing to break session compaction loops

**UUID:** `session-checkpoint-c4a8d7e1`
**Type:** improvement
**Reported:** 2026-09-10
**Status:** `done`

---

## Input (what was reported)

Sessions reaching ~160 turns require multiple `/compact` commands because context is
lost between compactions. Each compact resets the agent's working memory without a
structured handoff, causing re-investigation of previously solved problems.

## Context

HAWP architecture cleanup sessions involve many sequential edits (25+ files). The
agent loses track of implementation progress after each compact — repeating file reads,
re-tracing git status, re-identifying what was already committed.

## Analysis

**Root cause:** No structured checkpoint mechanism before compaction. Compaction resets
context without capturing: current plan state, completed slices, pending files, and
verification status.

**Directly verified:**

- Sessions using `/chronicle` require repeated re-investigation because compacted context
  does not persist implementation progress.
- HAWP already has `.hawp/work/status/` for checkpoint summaries — this is underutilized.

**Likely fix:**

1. Before any `/compact`, write a structured checkpoint to `.hawp/work/status/<timestamp>/checkpoint.md`
   capturing: current plan IDs, completed commits (SHA), pending files, verification status.
2. After compaction, agent reads the latest checkpoint and resumes without re-investigation.
3. Add this as a discipline rule in `.github/copilot-instructions.md`.

**Scope — what else is affected:**

- `.github/copilot-instructions.md` — add checkpoint-before-compact rule
- `CLAUDE.md` — same rule
- Agent mode docs (`AGENTS.md`, `.claude/rules/`)

## Work Coordination

**Owner:** Codex
**Implementation status:** done (2026-09-10)
**Can implement now:** yes

---

## Completion Record (2026-09-10)

The checkpoint-before-compaction rule was folded into
`.github/copilot-instructions.md` as part of the consolidated Copilot
agent-instruction update. The same rule has now been propagated to
`core/providers/.github/copilot-instructions.md` so GitHub provider updates can
carry the guidance to downstream repositories.

## Outcome

Copilot-facing HAWP instructions now direct agents to write a structured
checkpoint under `.hawp/work/status/YYYY/MM/DD/checkpoint.md` before compaction
or context reset, then read it afterward before resuming.

## Verification

Directly checked `.github/copilot-instructions.md` and
`core/providers/.github/copilot-instructions.md` for the session checkpoint rule.

## Close Checklist

- [x] Rule added to Copilot instructions.
- [x] Provider source pack aligned.
- [x] Plan moved to `closed/2026/09/10/session-checkpoint-c4a8d7e1/`.
