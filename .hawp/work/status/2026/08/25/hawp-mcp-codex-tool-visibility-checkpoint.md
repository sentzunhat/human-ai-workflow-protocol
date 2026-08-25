# Status Report

## Intent

Determine whether HAWP MCP is failing because of repo setup or because Codex is not surfacing the tools into live sessions.

## Current State

As of Tuesday, August 25, 2026, `beltrd/local-print-farm` appears correctly configured for HAWP MCP at the repo level, and host-side Codex CLI checks now confirm the `hawp` server is registered with `hawp_search`, `hawp_work_new`, and `hawp_work_validate` enabled. Even so, fresh Codex tasks still do not expose those tools in the live deferred tool inventory. Within the inspected evidence, the failure point is Codex-side MCP client readiness/session loading rather than repo HAWP registration.

## What Was Inspected

- Fresh-session deferred tool probes for `hawp_search`, `hawp_work_new`, `hawp_work_validate`, and broad `hawp`
- Host-side `codex mcp list` and `codex mcp get hawp` in `beltrd/local-print-farm`
- Repo-scoped HAWP MCP setup in `beltrd/local-print-farm/.codex/config.toml`
- Codex desktop logs under `~/Library/Logs/com.openai.codex/2026/08/25/`
- Codex local SQLite log store at `~/.codex/logs_2.sqlite`
- A known-good deferred MCP comparison using `node_repl`

## What Changed

- No repo HAWP files were changed within this checkpoint.
- This report now captures the exact observed failure point inside Codex session MCP resolution, plus the most justified host-side next steps.

## What Was Directly Verified

- `beltrd/local-print-farm/.codex/config.toml` contains `[mcp_servers.hawp]` with `command = "./.hawp/bin/hawp"`, `args = ["mcp"]`, `enabled = true`, and `enabled_tools = ["hawp_search", "hawp_work_new", "hawp_work_validate"]`.
- Host-side `codex mcp list` shows `hawp` as `enabled`.
- Host-side `codex mcp get hawp` reports:
  - `enabled: true`
  - `enabled_tools: hawp_search, hawp_work_new, hawp_work_validate`
  - `transport: stdio`
  - `command: ./.hawp/bin/hawp`
- Fresh-session deferred tool probes returned:
  - `hawp_search` -> `Found 0 tools`
  - `hawp_work_new` -> `Found 0 tools`
  - `hawp_work_validate` -> `Found 0 tools`
  - `hawp` -> `Found 0 tools`
- Fresh-session deferred tool probe for `node_repl` returned visible tools, proving deferred discovery is working for at least one other MCP server in the same live environment.
- No HAWP MCP tool was visible in the live session inventory during those probes.
- Within the reported evidence, none of the missing HAWP tools were callable from the live session.
- Codex desktop log `~/Library/Logs/com.openai.codex/2026/08/25/codex-desktop-59822e4b-336e-49bf-b3ae-fd13fbd77fb1-70843-t0-i1-164858-0.log` records `mcpServerStatus/list` activity during the relevant session.
- Codex SQLite log records repeated session-tool-catalog events stating `omitting MCP server without an exact ready client server_name=hawp`.
- The same SQLite log store contains successful MCP client initialization records for `node_repl`, including `Service initialized as client ... server_name=node_repl`.

## Assessment

Within the inspected scope, the failure point is Codex session MCP resolution, not repo HAWP registration. The host sees the `hawp` server and its enabled tools, but the live session omits `hawp` before tool surfacing because Codex does not have an exact ready client for that server.

## What Remains Unproven

- Why `hawp` is not reaching an exact ready client state inside Codex desktop
- Whether the missing ready client is caused by trust/loading boundary behavior, stale desktop cache, or app-lifecycle timing
- Whether a full desktop app restart is required before repo-scoped MCP tools appear in new tasks
- Whether a second repo-scoped custom MCP server on the same host would reproduce the same ready-client failure

## Constraints

- Investigate Codex desktop/session/tool-loading behavior only.
- Do not continue into printer work.
- Do not change repo HAWP files unless a host-side test directly points back to them.
- Prefer direct host/session evidence over repo-side reconfiguration guesses.

## Help Wanted

- Explain why `hawp` never reaches an exact ready client state while `node_repl` does.
- Challenge the current Codex-side diagnosis if a later host-side ready-client trace for `hawp` appears.
- Propose the next 2-3 highest-yield debugging steps after a full desktop restart and fresh-session retry.

## Suggested Next Step

Use a dedicated HAWP digital agent to investigate Codex host/session behavior only, with this checklist:

1. Confirm the task is rooted at `/Users/beltrd/Desktop/projects/beltrd/local-print-farm`.
2. Confirm `.codex/config.toml` still contains `[mcp_servers.hawp]`.
3. Run `codex mcp list` and `codex mcp get hawp` from the host and capture exact output.
4. In a brand-new Codex task, probe deferred inventory for:
   - `hawp_search`
   - `hawp_work_new`
   - `hawp_work_validate`
   - `hawp`
5. Confirm in `~/.codex/logs_2.sqlite` whether `hawp` still logs `omitting MCP server without an exact ready client server_name=hawp`.
6. Perform a full Codex desktop restart, then repeat the fresh-task deferred probes and log capture.
7. If the post-restart session still omits `hawp`, compare against another repo-scoped custom MCP server on the same host.
8. Do not spend more time changing `.mcp.json`, `.codex/config.toml`, or HAWP binaries unless a host-side ready-client test directly points back to them.

## Attached Artifact

### Prompt For A HAWP Digital Agent

```text
You are investigating a Codex-side MCP tool-visibility problem in /Users/beltrd/Desktop/projects/beltrd/local-print-farm.

Scope:
- Investigate Codex desktop/session/tool-loading behavior only.
- Do not continue into printer work.
- Do not change repo HAWP files unless a host-side test proves they are implicated.

Known context:
- Repo .codex/config.toml already contains [mcp_servers.hawp].
- Prior host checks have shown `codex mcp get hawp` with enabled tools:
  - hawp_search
  - hawp_work_new
  - hawp_work_validate
- Fresh Codex sessions on August 25, 2026 still returned:
  - `hawp_search` -> Found 0 tools
  - `hawp_work_new` -> Found 0 tools
  - `hawp_work_validate` -> Found 0 tools
  - `hawp` -> Found 0 tools

Your goals:
1. Verify whether Codex desktop is loading the hawp MCP server into live session tool inventory.
2. Distinguish among:
   - repo MCP misconfiguration
   - trust/loading boundary behavior
   - deferred tool indexing/caching
   - desktop session initialization behavior
3. Gather direct evidence, not guesses.

Required checks:
1. Confirm working repo root and inspect `.codex/config.toml`.
2. Run host-side `codex mcp list` and `codex mcp get hawp`.
3. In a fresh task, probe deferred tool inventory for the exact three tool names and broad `hawp`.
4. If tools are still missing, inspect Codex desktop logs around task creation and MCP/session startup.
5. If possible, compare with another local MCP server or known-good project-scoped MCP setup on the same host.

Return:
- observed facts
- exact failure point
- whether the issue is repo-side or Codex-side
- the next 2-3 most justified debugging steps

Do not propose repo HAWP changes unless the evidence points back to the repo.
```
