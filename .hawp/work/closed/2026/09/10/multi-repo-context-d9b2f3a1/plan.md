# Multi-repo context skill for agent

**UUID:** `multi-repo-context-d9b2f3a1`
**Type:** improvement
**Reported:** 2026-09-10
**Status:** `done`

---

## Input (what was reported)

Create a `.github/skills/multi-repo-context.md` skill that helps the agent navigate and
understand multi-repo dependency graphs when working across the human-ai-workflow-protocol
codebase and its downstream consumers.

## Context

The HAWP kit is designed to be installed into other repositories (per `core/providers/`
manifest). When agents work on changes that affect installed files in a shared `.hawp/kit`,
they need awareness of the install/update contract and how changes propagate.

## Analysis

**Root cause:** No existing skill captures the multi-repo lifecycle: kit → install → update →
consumer repos. Agents treat each repo in isolation without considering downstream impact.

**Directly verified:**

- `core/providers/manifest.yaml` defines provider-specific installation paths.
- `distribution/generated/` and `distribution/sources/` contain install/update scripts.
- No `.github/skills/` directory exists yet.

**Likely fix:**

1. Create `.github/skills/multi-repo-context.md` capturing: HAWP kit lifecycle, provider
   installation paths, update propagation rules, common failure modes across repos.
2. Reference the skill from provider-specific instruction files (CLAUDE.md, AGENTS.md, etc.).

**Scope — what else is affected:**

- `.github/skills/multi-repo-context.md` — new file
- Provider instruction files in `core/providers/` — add skill reference
- Distribution docs — link back to the skill

## Work Coordination

**Owner:** Codex
**Implementation status:** done (2026-09-10)
**Can implement now:** yes

---

## Completion Record (2026-09-10)

- Created `.github/skills/multi-repo-context.md` — three-lane model (source → downstream → consumer), propagation rules, provider overlay map, cross-repo failure modes, and downstream-consumer guidance.
- Referenced the skill from `AGENTS.md` (shared entry point).
- Distribution docs intentionally left unchanged: the skill is a source-repo navigation aid and is not part of the install/update contract propagated to consumers.

## Outcome

A source-repo navigation skill now captures the multi-repo HAWP lifecycle
(source → downstream → consumer), provider installation paths, update propagation
rules, and common cross-repo failure modes. The skill is referenced from the shared
`AGENTS.md` entry point so agents working on installed `.hawp/kit` files see the
propagation contract. Distribution install/update docs were intentionally left
unchanged because the skill is not part of the consumer-facing contract.

## Verification

`.github/skills/multi-repo-context.md` exists and is referenced from `AGENTS.md`.
The skill content covers the three-lane model, propagation rules, provider
overlay map, and cross-repo failure modes. No install/update distribution output
or provider overlay was modified, so no runtime or build checks were required for
this documentation-only change.

## Close Checklist

- [x] Skill file created and referenced from the shared entry point.
- [x] Distribution docs left unchanged by design.
- [x] No new dependency, runtime, or install-contract change.
- [x] Plan moved to `closed/2026/09/10/multi-repo-context-d9b2f3a1/`.
- [x] BACKLOG.md updated (moved to Recently Closed).
