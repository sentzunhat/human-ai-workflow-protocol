# 2026-10-15 — HAWP Codex MCP availability and provider update checkpoint

## Status Report

### Intent

Preserve the evidence and continuation state for work item `c147f2ae`, which audits why HAWP tools were not available by default in this repository's Codex sessions. The user authorized the Codex provider update and requested durable conversation, memory, and GitHub archival.

### Current State

The primary launch failure is repaired: the project-local HAWP executable now exists and responds. The Codex MCP configuration still allows only four of the seven tools the server advertises, and this conversation has not received callable HAWP tools. Work item `c147f2ae` remains plan-ready; fresh-session exposure and repository-scope validation are pending.

### What Was Inspected

- `.codex/config.toml` and the configured `hawp` MCP entry.
- HAWP provider update instructions and the user-supplied update command.
- Installed CLI version and direct MCP JSON-RPC initialize/tools-list response.
- `codex mcp get hawp` configuration output.
- HAWP usage summary and relevant prior closed work.
- Active and parked work plans, backlog, and HAWP status-report guidance.

### What Changed

- Executed the user-authorized Codex provider update from the repository root. It reported source `sentzunhat/human-ai-workflow-protocol@main`, provider `codex`, and remote archive mode.
- Restored `.hawp/bin/hawp` at release `0.0.24`; the updater reported SHA-256 verification. It refreshed managed kit/license files, seeded provider guidance only if missing, and preserved `.hawp/work/**`.
- Recorded the repaired executable, the still-restricted tool allowlist, and the fresh-session verification plan in active work item `c147f2ae`.
- Continued the existing `c147f2ae` investigation; no duplicate work item was created.

### What Was Directly Verified

- Before the update, `.codex/config.toml` pointed to a missing `.hawp/bin/hawp`. Afterward, `./.hawp/bin/hawp version` reported `0.0.24`.
- A direct read-only JSON-RPC initialize and `tools/list` through the installed executable succeeded and advertised seven tools.
- `codex mcp get hawp` reports the server enabled, but the explicit `enabled_tools` allowlist still contains only `hawp_search`, `hawp_usage`, `hawp_work_new`, and `hawp_work_validate`.
- `hawp_work_intake`, `hawp_work_doc`, and `hawp_work_reshape` are therefore still excluded by the configured allowlist. This conversation still exposes no callable HAWP tools.
- The usage report spans 2026-08-27 through 2026-10-15 and shows 148 historical calls: 96 `hawp_work_validate`, 26 `hawp_search`, 13 `hawp_work_new`, 12 `hawp_work_intake`, and 1 `hawp_usage`. These aggregate historical records do not establish current-session availability.
- Existing closed work includes `c98518bb` (HAWP-first session guide) and `e44ff827` (Copilot MCP fallback protocol).
- The HAWP validator passed after the workflow changes with active 1/1, closed 12/12, parked 17/17, and zero issues, warnings, or broken links.

### What Remains Unproven

- The provider update did not synchronize `.codex/config.toml`; whether to expand the allowlist to all seven advertised tools remains for review.
- No fresh trusted task/session has yet confirmed that Codex exposes the configured tools after the executable repair and any allowlist synchronization.
- A direct server probe proves the binary's MCP response, not client discovery, session exposure, or write scope.
- The usage totals are historical and aggregated; they do not identify whether a particular current session loaded this server.

### Constraints

- The provider update was explicitly authorized; no further provider/runtime/configuration changes are included in this checkpoint.
- The audit is separate from downstream consumer configuration and does not infer that other repositories' MCP setups are repaired.
- Parked plans and backlog edits are existing working state and are excluded from this checkpoint commit.
- This public repository checkpoint contains no machine-local absolute paths, personal account details, or private analytics URLs.

### Help Wanted

Review whether the explicit Codex allowlist should be synchronized to the seven tools the server advertises, and whether the proposed fresh-session checks are sufficient to close the audit.

### Suggested Next Step

Review the plan for `c147f2ae`. If approved for execution, synchronize the Codex allowlist through the supported setup path, start a fresh trusted task, verify the HAWP tools are exposed, confirm repository scope with `hawp_work_validate`, and run a read-only `hawp_search`. Keep the parked website lane paused until this audit is ready to resume.

## Conversation Checkpoint Summary

### Strategic Impact

The missing local executable explains the confirmed server startup failure and has been repaired. The allowlist mismatch is a separate confirmed configuration gap; actual tool visibility remains unverified until a fresh session. The update did not touch Zacatl or establish that any downstream consumer repository is fixed.

### Continuation State

- Branch: `feature/website`
- Checkpoint base before archival: `8e0afdf6fc1a8d962478468a7bd228c2b1f443f8`
- Active audit: `c147f2ae` (`plan-ready`)
- Parked work: `7394f90c`, `c8d9e1cb`, `2eea565c`, `429e075e`, `96c45719`, `1607037f`, `b4300599`
- The seven parked work items remain unchanged by this checkpoint commit.

### Memory Delta

The existing HAWP MCP memory note documents an earlier consumer-side missing-binary and fresh-chat symptom. This checkpoint adds the distinct source-repository provider update result, seven-tool server inventory versus four-tool allowlist, and the current fresh-session gate. A compact memory follow-up note is saved separately for integration with that MCP timeline; `MEMORY.md` is not edited directly.

### GitHub Archival

Repository: `sentzunhat/human-ai-workflow-protocol`

Branch: `feature/website`

Checkpoint: `.hawp/work/status/2026/10/15/c147f2ae/status.md`

### Resume Summary

Resume work item `c147f2ae` from its plan-ready state. Review the allowlist recommendation, then verify tool exposure from a fresh trusted task before claiming the tool-first availability issue is fully resolved.
