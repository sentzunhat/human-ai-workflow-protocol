# defer CLI participant adapters for Codex, Claude, and GitHub

**Backlog ID (Legacy):** — (UUID-native item)
**UUID:** `bee15107-e4d1-41b3-b67e-04c87328fef3`
**Type:** improvement
**Reported:** 2026-07-06
**Risk Level:** medium
**Status:** parked

---

### Input (what was reported)

> Amazing the codex cli is not needed yet same as claude cli or github cli let's keep those in a parked worked item

---

### Context

The Codex provider pack is useful now because it installs repo-local `AGENTS.md`
instructions. Runtime participant adapters are a different layer: they would
invoke external CLIs as part of a future workflow-loop runner.

The user explicitly does not want `codex-cli`, `claude-cli`, or `github-cli` work
started yet.

---

### Parked Scope

Future work, when resumed:

- Define whether workflow-loop runtime orchestration is in scope.
- Decide participant adapter contract for external CLIs.
- Add adapters only after the loop runner, prompt compiler, output capture, and
  transition/evaluation model are ready enough to use them.
- Candidate participants:
  - `codex-cli`
  - `claude-cli`
  - `github-cli` or GitHub CLI-backed PR/check participant

---

### Reason Parked

Provider packs are enough for current Codex/Claude/GitHub usage. CLI participant
adapters would add runtime orchestration scope before the repo needs it.

---

### Resume Trigger

Resume this item only when the user asks for workflow-loop automation,
participant adapters, external CLI orchestration, or a concrete Codex/Claude/GitHub
CLI runner.

---

## Outcome

_Parked before implementation by user request._

## Verification

_No implementation performed. Parking verified by BACKLOG row pointing to this file._

