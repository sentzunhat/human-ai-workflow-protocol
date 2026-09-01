# Evidence: Codex MCP project-trust gate

**Date:** 2026-08-25  
**Repo under test:** `local-print-farm`  
**Codex binary:** `codex` on PATH

---

## Verified facts

### Positive case (trusted project + CLI)

- `local-print-farm` has `.codex/config.toml` pointing at `hawp mcp`
- The project is listed as trusted in `~/.codex/config.toml` (line 164)
- `codex mcp list` shows `hawp` enabled with `hawp_search`, `hawp_work_new`, `hawp_work_validate`
- Source: `.codex/config.toml` in that repo (not `.mcp.json`)

### Negative case (untrusted throwaway repo)

- `codex-mcp-probe/.codex/config.toml` had same shape, pointed at OpenAI Docs MCP server
- `codex mcp list` omitted it; `codex mcp get openaiDeveloperDocs` returned "No MCP server named ... found."
- Same config, different project trust status → server not loaded

### Desktop behavior

- App was making `mcpServerStatus/list` requests on 2026-08-25
- No `hawp` startup failure found in inspected log slice
- Fresh-task probe found 0 tools for both exact HAWP tool names and broad `hawp`
- Conclusion: desktop session does not hot-reload project MCP config changes

---

## Loading boundary (verified)

| Condition | CLI `mcp list` shows hawp | Desktop task sees tools |
|-----------|--------------------------|------------------------|
| Trusted project + `.codex/config.toml` present + fresh session | ✅ | Unknown (not confirmed) |
| Trusted project + config present + existing desktop session | ✅ | ❌ (probe found 0 tools) |
| Untrusted project + config present | ❌ | ❌ |

---

## What this means for HAWP

1. **`.codex/config.toml` is the correct path** — confirmed by the positive case
2. **Project trust is a required second gate** — file presence alone is not enough
3. **Desktop UI does not hot-reload** — users must start a fresh task/session after config or trust changes
4. **`.mcp.json` is non-authoritative for Codex** — Codex ignores it; Claude Code reads it

---

## Impact on v0.0.9

- `hawp init --provider codex` already writes `.codex/config.toml` (fixed in v0.0.9)
- CLI output now prints the trust note and verification command
- `search.md` updated with trust requirement and per-agent config table
- The desktop exposure gap (trusted project + config but 0 tools in desktop task) is a
  Codex behavior question, not a HAWP config issue — no further HAWP change needed

---

## Root cause found — 2026-08-25 follow-up

The "0 tools in desktop task" failure was **not** a trust gate and not a Codex hot-reload
issue. It was a **relative path resolution bug** in the generated `.codex/config.toml`:

| Field | Old value | Effect |
|-------|-----------|--------|
| `command` | `"./.hawp/bin/hawp"` | Resolved relative to Codex app dir → binary not found |
| `cwd` | `"."` | Also resolves to Codex app dir, not project root |

Evidence: Codex SQLite log entry `omitting MCP server without an exact ready client server_name=hawp`.
Contrast with `node_repl` (global config, absolute path) → `Service initialized as client`.

Manual proof: running `"./.hawp/bin/hawp"` from `/tmp` → `no such file or directory`.
Running from the actual project root → binary found and MCP handshake succeeds.

### Fix applied

- `local-print-farm/.codex/config.toml` updated to use absolute `command` and `cwd` (hand-patched)
- `hawp init --provider codex` (Go source `mcp/config.go`) updated to write absolute paths using repo root
- `writeCodexTOML` signature changed to `writeCodexTOML(path, repoRoot string)` — passes absolute paths
- All tests updated and passing; binary rebuilt as v0.0.10
- `local-print-farm/.hawp/bin/hawp-bin` updated to v0.0.10

### Verification needed

Start a fresh Codex task in `local-print-farm` and probe `hawp_search` via deferred tool
inventory — should now surface rather than returning 0 tools.
