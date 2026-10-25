---
uuid: b152c56c
title: "HAWP provider pilot and v0.0.26 measurement checkpoint"
type: status
date: 2026-10-25
---

# 2026-10-25 — HAWP provider pilot and release measurement checkpoint

## Status Report

### Intent

Continue the active provider setup item `b152c56c`, record the bounded MCP pilot evidence for the three priority clients, and clarify how v0.0.26 will measure HAWP-visible usage while keeping preference testing separate.

### Current State

`b152c56c` remains in progress for v0.0.25. The HAWP CLI configured repo-scoped entries for Codex, GitHub Copilot/VS Code (`github`), and Cursor. Configuration generation is not proof that each client's current session loaded the entry or approved a call. Only a synthetic Codex search has a recorded successful live call in this checkpoint; Copilot and Cursor live verification remain open.

The existing plan `1baa2747` now scopes v0.0.26 usage evidence to HAWP-observable input and output sizes, search item totals, and pipeline shaping counts. `b921093c` separately plans paired capability-choice and fallback measurements across the three clients. Neither study is complete and neither release is locked.

### What Was Inspected

- `.hawp/work/BACKLOG.md`
- `.hawp/work/active/b152c56c/plan.md`
- `.hawp/work/active/2eea565c/plan.md`
- `.hawp/work/active/1baa2747/plan.md`
- `.hawp/work/active/b921093c/plan.md`
- `.hawp/work/status/2026/10/21/b152c56c/status.md`
- The selected What to Ship project/release list
- GitHub repository listing for `beltrd` and `sentzunhat`

### What Changed

- The v0.0.26 usage plan now asks for HAWP-observable query/tool argument bytes and characters, fixed synthetic prompt size, search result/candidate totals, deduplicated candidates when available, chunks used/dropped, and post-shaping response sizes delivered to the client.
- The plans distinguish direct byte/character measurements from chars-per-four token estimates. The MCP server cannot see the provider's complete system/main prompt, provider model token counts, or billing unless a separate approved client integration exposes them. Real private conversation prompts are not to be captured for measurement.
- Preference and fallback evaluation remains a separate v0.0.26 work item (`b921093c`); capability listing alone is not adoption or preference evidence.
- The duplicate `b921093c` backlog row was removed. Local HAWP work validation and `git diff --check` passed.
- What to Ship was not changed because the connected project list exposed the Example and Mochilada projects, but no HAWP project.

### What Was Directly Verified

- HAWP's binary configure command completed for the repo-scoped Codex, GitHub Copilot/VS Code, and Cursor provider entries.
- A fixed synthetic Codex `hawp_search` query succeeded in raw and context-shaped modes. The query was 41 UTF-8 bytes; raw mode returned 2 results with 2,290 bytes of result content and 3,090 serialized bytes. Context-shaped mode used 2 chunks, dropped 0, returned 2,268 bytes, and reported 567 estimated tokens. These are HAWP payload observations, not provider token counts.
- The local usage configuration was enabled with request-body logging disabled. The current logger still stores a bounded query-text field; avoid private query text during future pilot runs.
- The repo-local validator passed all five groups with zero issues and warnings. `git diff --check` passed.
- The GitHub repository selected for this checkpoint is the dedicated `sentzunhat/human-ai-workflow-protocol` repository. The personal-account repository listing returned no repositories. Recent organization repositories were also inspected; none was a closer match to this HAWP checkpoint.
- The working branch was `feature/mcp-providers`, at `71b14bb65bd8a04f87ab450aa91baa7363bdf85c`, tracking `origin/feature/mcp-providers` with no reported ahead/behind commits. Multiple unrelated local plan/config changes were already present and are excluded from this archival checkpoint commit.

### What Remains Unproven

- Copilot's current session has not returned an inspectable result for the bounded synthetic HAWP call; effective tool visibility and approval behavior remain unverified.
- Cursor's recorded prior smoke reported that HAWP was absent from that session's MCP namespace. That older observation does not verify whether Cursor loaded the subsequently generated repo-scoped configuration; refresh/reconnect and repeat one bounded read-only call.
- The CLI-generated configuration being present does not prove a client loaded it. Record actual client version, server connection, visible HAWP tools, and successful read-only invocation for each provider.
- HAWP does not currently expose all candidate counts or byte sizes at every pre/post pipeline stage, and current usage logging does not count MCP resource/prompt protocol methods. The instrumentation plan must report missing stages as unobserved.
- Provider system/main prompt size, total provider prompt tokens, and model token consumption remain unavailable at the MCP boundary. Do not infer cost or preference from HAWP payload measurements.
- The active What to Ship connector did not list HAWP, so release-board synchronization remains outstanding until the HAWP project is available in that connected workspace.

### Constraints

- Pilot Codex, GitHub Copilot, and Cursor as tool-host clients. Keep model/intelligence backends distinct from tool-host clients.
- Approve only bounded read-only pilot calls, such as `hawp_search` and `hawp_work_validate`; do not use write-capable tools as approval tests.
- Use fixed synthetic prompts. Do not save real private conversation text; note that current logging may store a bounded query string even with request-body logging disabled.
- Do not claim provider token totals, billing savings, broad client preference, or release completion from these observations.
- Preserve the existing dirty working tree; commit and push only this status report.

### Help Wanted

No new user input is required for the planned read-only client smoke. If a client presents a permission request broader than the bounded read-only HAWP calls, stop and record the exact gate.

### Suggested Next Step

Resume `b152c56c`: refresh Copilot and Cursor separately, confirm each loads the generated HAWP server, record client/version and the visible tool list, then approve one fixed synthetic read-only `hawp_search` call per client. Update the plans with observed outcomes. Keep preference/fallback comparisons for `b921093c` and durable pipeline metrics/report work for `1baa2747` in v0.0.26. Revisit What to Ship only when the HAWP project is visible to the connected release-board tools.
