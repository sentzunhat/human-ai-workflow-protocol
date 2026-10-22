# Audit HAWP tool-first usage, availability, and fallback behavior

**Backlog ID (Legacy):** — (UUID-native item)
**UUID:** `c147f2ae-28e8-4e8b-abca-2f4e4a61838b`
**Type:** improvement
**Reported:** 2026-10-15
**Status:** `closed`
**Risk Level:** low for audit; medium for any configuration or provider changes
**Priority:** immediate; unblocks resumption of the parked website lane

---

## Input (verbatim)

> Audit HAWP tool-first usage, availability, and fallback behavior

## Intake Summary

The repository has a HAWP-first guide and prior closed work covering MCP search defaults and Copilot fallback, but this chat did not expose callable HAWP MCP tools. The local Codex MCP configuration is enabled; actual invocation capability and configuration drift need to be distinguished.

## Current Context

The source repository is configured for stdio MCP at `.codex/config.toml`; `codex mcp get hawp` reports the server enabled. The config initially pointed to a missing `.hawp/bin/hawp`. The user then authorized and supplied the Codex provider update command, which fetched `sentzunhat/human-ai-workflow-protocol@main`, updated binary `0.0.24`, verified SHA-256, and refreshed the kit. The installed binary now launches; a direct JSON-RPC initialize plus `tools/list` returned all seven tools. The current `enabled_tools` allowlist still has four entries, so the Codex client will filter out `hawp_work_intake`, `hawp_work_doc`, and `hawp_work_reshape`. This conversation itself still exposes no HAWP tools. The read-only HAWP usage report covers 2026-08-27 through 2026-10-15 and records 148 calls: 96 `hawp_work_validate`, 26 `hawp_search`, 13 `hawp_work_new`, 12 `hawp_work_intake`, and 1 `hawp_usage`. Prior closed records include `c98518bb` (HAWP-first session guide) and `e44ff827` (Copilot MCP fallback protocol).

**Continuation findings (2026-10-16):** HAWP tools were exposed in the fresh task. `hawp_work_validate` confirmed repository root `/Users/beltrd/Desktop/projects/sentzunhat/human-ai-workflow-protocol` and passed; read-only searches initially returned no results. With user approval, `.hawp/bin/hawp search index` rebuilt the ignored local index (597 documents, 5,001 chunks, 118 work-item metadata records). Follow-up read-only searches returned the c147f2ae plan and relevant MCP setup/verification guidance. The repo-local Codex config points to this checkout's `.hawp/bin/hawp` and allowlists all seven tools.

## Initial Analysis

**Directly verified:**

- Before the user-authorized update, the Codex MCP entry pointed to a missing `.hawp/bin/hawp`; the update restored it to version `0.0.24` and verified its SHA-256.
- A direct read-only JSON-RPC probe through the installed binary now succeeds and returns all seven tools.
- This task exposes callable HAWP tools; `hawp_work_validate` confirmed the target repository root and passed with zero issues and warnings.
- The binary's `mcp configure --provider all` updated `.codex/config.toml` to allow all seven tools and wrote/updated the Claude, Cursor, and VS Code/Copilot MCP configuration files.
- `.gitignore` now excludes `.codex/config.toml`, `.mcp.json`, `.cursor/mcp.json`, `.vscode/mcp.json`, `.claude/settings.local.json`, and `.continue/config.yaml`; `git check-ignore` confirmed each path.
- Ollama 0.31.1's local API responded and listed eight installed models; `ollama ps` showed none loaded at inspection time.
- Two read-only `hawp_search` calls returned no results, so a useful repository search is not yet verified.
- The usage log shows historical MCP use, including `hawp_work_intake`, but does not establish that this checkout's current configured command was launchable.
- `hawp_work_validate` passed before this workflow update with zero issues, warnings, or broken links.

**Inferred (not yet proven):**

- Confirmed and repaired primary cause: the project-local binary path was missing, so the configured MCP server could not launch.
- The fresh trusted Codex session check is still required to confirm the updated project configuration is loaded by the client.
- The empty search results may mean the index is absent/stale or the query did not match; this has not been diagnosed.

**Likely scope:**

- Audit config, provider guidance, tool inventory, and recent usage by client/session.
- Reconcile configured versus expected tools and identify whether the MCP server is callable when it is merely registered as enabled.
- Propose a clear tool-first rule with an explicit, bounded CLI fallback when tools are not exposed.
- Keep future HAWP runtime/config/provider changes out of scope unless separately authorized; this continuation synchronized the requested local client configs.

## Risk + Review Gate

**Risk:** low for read-only audit; medium for any configuration/provider changes
**Gate:** review findings before configuration or guidance changes

## Backlog + Plan Link

**Status now:** closed
**Plan file:** work/closed/2026/10/16/c147f2ae/plan.md

## Next Step

- [x] Investigation recorded above (required before planning)
- [x] High-level audit plan written
- [x] Move backlog status accordingly

## Mission

Determine why HAWP MCP tools are not being used as the default in sessions where the repository is configured for HAWP, separate global usage from current-session availability, and recommend the smallest fix before implementation.

## High-Level Audit Plan

1. Compare repo and user-level MCP configuration, allowlists, provider instructions, and the tools actually exposed in representative sessions.
2. Trace server startup and handshake evidence, distinguishing “registered/enabled” from “callable in this conversation.”
3. Compare configured tools against HAWP’s supported tool set and usage records, including the `hawp_work_intake` discrepancy.
4. Review the existing HAWP-first guide and prior fallback work; identify instruction or client gaps without duplicating completed work.
5. The user-authorized provider update restored and checksum-verified the missing executable; direct server discovery succeeds.
6. [x] Synchronize repo-local MCP configs with `mcp configure --provider all`; confirm all seven Codex tools are allowed and every identified repo-local client config is ignored.
7. [x] In a fresh trusted Codex task, confirm all seven tools are exposed, validate repository identity, and verify useful read-only search after the approved index rebuild.

## Outcome and Closeout

Closed 2026-10-16 after all verification gates passed. Fresh-task MCP exposure included all seven tools; repository-scoped validation passed (3 checks, 0 issues, 0 warnings; links valid in 148 Markdown files). Read-only search returned `.hawp/work/active/c147f2ae/plan.md` for the distinctive audit phrase and `.hawp/kit/usage/mcp/install-agent.md` plus `.hawp/kit/usage/mcp/verification.md` for Codex MCP guidance. The local Codex config points to this repository's `.hawp/bin/hawp` and allows all seven tools. The approved index rebuild reported 597 documents, 5,001 chunks, and 118 work-item metadata records.

**Close date:** 2026-10-16

## Verification

- Fresh task exposed all seven HAWP MCP tools; `.codex/config.toml` points to the repository binary and allowlists all seven.
- `hawp_work_validate` passed for `/Users/beltrd/Desktop/projects/sentzunhat/human-ai-workflow-protocol`: 3 checks, 0 issues, 0 warnings; 148 Markdown files had valid local links.
- Read-only `hawp_search` returned this plan for the audit phrase and relevant MCP guidance in `.hawp/kit/usage/mcp/install-agent.md` and `.hawp/kit/usage/mcp/verification.md`.
- Approved index rebuild completed with 597 documents, 5,001 chunks, and 118 work-item metadata records.

## Close Checklist

- [x] Fresh-task exposure, repository-scoped validation, useful search, and Codex config evidence recorded above.
- [x] Existing plan archived and backlog active row removed; parked item statuses preserved.

## Resume Coordination

The seven items parked during this intake resume after this audit plan is reviewed and the audit is ready for execution. Website sequencing remains pentagon (`96c45719`) before cable motion (`1607037f`); analytics (`b4300599`) remains in the same resume-priority group.
