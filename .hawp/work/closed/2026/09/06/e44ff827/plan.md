# Copilot MCP + Terminal Fallback Protocol

**UUID:** `e44ff827-0b71-4ef8-9288-74c2be1bd372`
**Type:** improvement
**Reported:** 2026-09-06
**Risk Level:** low

---

## Input (what was reported)

> Add MCP + terminal fallback protocol to `.github/copilot-instructions.md` so agents don't enter infinite retry loops when MCP tools are unavailable in the current conversation.

---

## Context

The agent repeatedly tries to invoke `hawp_work_validate` through MCP when it's explicitly told "MCP only" — but the conversation doesn't expose MCP tool calls for invocation. This created a 5-turn deadlock (sessions `7fe3ae29`, `2ee14b51`) where the agent couldn't proceed, and you had to reconfigure and retry each time.

**Directly verified:**

- Session `7fe3ae29` — 37 turns across multiple back-and-forths; MCP server started but "MCP TOOL UNAVAILABLE" in conversation context.
- Session `2ee14b51` — same pattern: agent confirmed MCP unavailable, walked you through starting it, then still couldn't invoke it.
- Terminal output confirmed HAWP had 4 tools exposed (`2026-09-06 12:29:15.143 [info] Discovered 4 tools`) but this session never accessed them.

**Inferred (not yet proven):**

- The same pattern likely occurs with other MCP tool invocations across repos where the agent is told "MCP only" or when the server starts but isn't auto-exposed in the conversation.
- Other HAWP CLI tools (`hawp_work_new`, `hawp_work_validate`, `hawp_check`) are similarly unreachable via MCP from Copilot conversations even when the server is running.

**Scope — what else is affected:**

- `.github/copilot-instructions.md` (add fallback rules)
- All repos using HAWP MCP with GitHub Copilot (local-print-farm, mictlan, infra-as-code, zacatl, human-ai-workflow-protocol)

---

## Analysis

**Root cause (or most likely cause):**
GitHub Copilot conversations don't automatically expose MCP tool calls for invocation. Even when the HAWP server is running and exposing 4 tools in VS Code's MCP list, the Copilot agent has no mechanism to call them — it only sees terminal output or file content. The agent doesn't know this limitation and keeps retrying indefinitely.

**Directly verified:**

- Agent reported "MCP TOOL UNAVAILABLE" after confirming server logs showed 4 tools discovered.
- No terminal fallback was defined, so the agent stopped instead of trying `go run ./cmd/hawp`.

**Scope — what else is affected:**

- Only `.github/copilot-instructions.md` needs changes (add explicit fallback protocol).

---

## Work Coordination

**Owner:** unassigned
**Implementation status:** done
**Parallel work risk:** low
**Can implement now:** yes
**Coordination note:** None — this is a single-file edit to instructions. Only one active item touches `.github/copilot-instructions.md`.

---

## Verification

Implemented 2026-09-06: Added "MCP Tool Invocation Protocol" section to `.github/copilot-instructions.md` with explicit fallback rules and one-attempt cap. Verified `go run ./cmd/hawp check --no-update-check` passes all 3 validations (kit, work, links). Plan closed via commit `92cb6ac`.

---

## Next Step

1. Read current `.github/copilot-instructions.md`
2. Add a dedicated section: "MCP Tool Invocation Protocol" with these rules:
   - **Default:** Try MCP first; report the exact tool name used.
   - **Fallback:** If MCP is unavailable (tool not found, server error, or conversation doesn't expose tools), fall back to terminal commands (`go run ./cmd/hawp <command>`) with explicit disclosure ("MCP tool unavailable, using CLI fallback").
   - **Never retry indefinitely** — one fallback attempt per requested operation.
3. Update the `Do not` section if needed to prevent MCP-only deadlocks.
4. Run HAWP validation to confirm no doc drift: `go run ./cmd/hawp check`.

---

## Outcome

- `.github/copilot-instructions.md` updated with "MCP Tool Invocation Protocol" section (rules for default→fallback, one-attempt cap, explicit disclosure)
- `go run ./cmd/hawp check --no-update-check` passes all 3 validations (kit validate, work validate, links check)
- No files modified beyond instruction file

---

## Close Checklist

- [x] Plan written and reviewed
- [x] Implementation completed
- [x] Verification evidence recorded
- [x] Changes committed (`92cb6ac` on `feature/v0.0.24-work-folder-normalization`)
