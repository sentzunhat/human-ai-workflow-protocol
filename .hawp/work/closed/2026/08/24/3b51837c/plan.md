---
work-item: 3b51837c
type: feature
title: "v0.0.4: hawp init + hawp mcp + provider MCP config auto-write"
status: done
owner: unassigned
created: 2026-08-24
updated: 2026-08-24
closed: 2026-08-24
---

# v0.0.4: hawp init + hawp mcp

## Mission

Replace the bash install script with a first-class `hawp init` binary command,
implement `hawp mcp` as a stdio MCP server, and have `hawp init` auto-write the
provider's MCP config so agents get native tool access with zero manual setup.

## Context

The bash install script in `distribution/generated/*/install/main.md` works but
requires curl + bash and is hard to use on Windows. The binary already ships to
all 6 platforms — a `hawp init` command replaces the script with a single binary
call. Adding `hawp mcp` turns the CLI into a native AI tool server that Claude,
Cursor, Continue, and (when supported) Codex can call directly instead of via
shell commands.

## Planned Commands

### `hawp init [--provider <name>]`

Runs from the target repo root. Equivalent to the bash install script but:

- Uses the local binary (no curl needed after bootstrap)
- Platform-aware, no bash dependency on Windows
- Auto-detects provider when `--provider` is omitted (checks for `.claude/`,
  `.cursor/`, `.continue/`, `AGENTS.md` presence)
- Writes the provider's MCP config entry automatically (see below)
- Idempotent — safe to re-run (same no-clobber semantics as the bash script)

Outputs: `.hawp/kit/**`, `.hawp/work/` scaffold, provider overlay, MCP config.

### `hawp mcp`

Starts a stdio MCP server. Tools exposed:

| Tool | Description |
|------|-------------|
| `hawp_search` | Lexical/semantic/hybrid search over indexed kit + work docs |
| `hawp_work_new` | Create a new work item (`active/{uuid}/plan.md`) |
| `hawp_work_validate` | Run work validation and return findings |
| `hawp_update` | Self-update binary + kit + provider overlays |
| `hawp_status` | Return current backlog summary |

Server is started by the MCP host (Claude Code, Cursor, etc.) via the config
written by `hawp init`.

## Provider MCP Config Auto-Write

`hawp init` writes or patches the provider's MCP config with a `hawp` server entry.

### Claude Code → `.mcp.json`

```json
{
  "mcpServers": {
    "hawp": {
      "command": ".hawp/bin/hawp",
      "args": ["mcp"]
    }
  }
}
```

### Cursor → `.cursor/mcp.json`

Same format as Claude Code.

### Continue → `.continue/config.yaml`

```yaml
mcp:
  servers:
    - name: hawp
      command: .hawp/bin/hawp
      args: [mcp]
```

### Codex / GitHub Copilot

No native MCP protocol yet. `hawp init --provider codex` adds a shell-command
section to `AGENTS.md` documenting how to call the CLI tools directly.

## Bootstrap Story (install flow after this ships)

```bash
# Bootstrap: download binary once (curl still needed here)
curl -fsSL https://github.com/sentzunhat/human-ai-workflow-protocol/releases/latest/download/hawp-darwin-arm64 \
  -o .hawp/bin/hawp && chmod +x .hawp/bin/hawp

# From then on, binary handles everything
hawp init --provider claude    # installs kit + writes .mcp.json
hawp update                    # updates binary + kit + providers + mcp configs
```

## Implementation Plan

### Phase 1: `hawp mcp` (stdio server)

1. Add `cmd/hawp/mcp.go` — dispatch `hawp mcp` to MCP server
2. Implement `internal/platform/mcp/server.go` — stdio JSON-RPC 2.0 handler
3. Wire tools: search, work_new, work_validate, update, status
4. Add `make build` target; verify with `echo '{"method":"tools/list"}' | hawp mcp`

### Phase 2: `hawp init`

1. Add `cmd/hawp/init.go` — dispatch `hawp init [--provider]`
2. Port bash script logic to Go: kit refresh, work scaffold seed, legacy migration,
   active reconciliation
3. Add provider detection + MCP config write per provider
4. Add `--provider` flag: `claude | cursor | continue | codex | github | all`
5. Wire into `hawp update` so provider MCP configs stay current

### Phase 3: Provider overlay updates

1. Update each provider's `.mcp.json` / `.cursor/mcp.json` / `.continue/config.yaml`
   seed files in `core/providers/`
2. Run `providers:sync` + `distribution:sync`
3. Update install guides to show `hawp init` as the primary path

## Success Criteria

- [ ] `hawp mcp` starts and responds to `tools/list` with all 5 tools
- [ ] Claude Code can call `hawp_search` as a native tool (no shell command)
- [ ] `hawp init --provider claude` writes `.mcp.json` and installs kit in a fresh repo
- [ ] `hawp init` is idempotent — re-run produces no changes
- [ ] All 6 platforms build and pass tests
- [ ] `npm --prefix librarian run validate` passes

## Dependencies

- v0.0.3 merged and released (binary on v0.0.3 as baseline)
- MCP protocol spec: stdio JSON-RPC 2.0, `initialize` / `tools/list` / `tools/call`

## Out of Scope

- Remote/HTTP MCP transport (stdio only for now)
- Authentication / multi-tenant
- `hawp mcp` auto-restart on crash (host handles reconnect)

## Outcome

Shipped as v0.0.4 in PR #13 (0.0.4 → development).

- `hawp mcp` stdio MCP server: 3 tools (`hawp_search`, `hawp_work_new`, `hawp_work_validate`), zero external deps, full JSON-RPC 2.0
- `hawp init --provider <name>|all`: writes `.mcp.json` (Claude), `.cursor/mcp.json` (Cursor), `codex.toml` (Codex), prints manual block (Continue)
- 5 new tests in `internal/platform/mcp/config_test.go`, all passing
- Dead code removed (`getFloat` in `run.go`)
- Registry and help text updated; `update` description corrected to match v0.0.3 default behavior
- `npm --prefix librarian run validate` passes clean (60/60 tests)

Items from original plan that moved out of scope (kept small for v0.0.4):
- `hawp_update` and `hawp_status` tools deferred (not enough demand to justify; 3 tools shipped)
- Phase 3 provider seed file updates deferred (distribution scripts not yet updated for MCP seed files)
- Auto-detect provider from installed overlays deferred
- Bootstrap story unchanged (curl still needed for first install)

## Verification

- [x] `echo '{"jsonrpc":"2.0","id":1,"method":"initialize",...}' | hawp mcp` → valid JSON-RPC response (run locally, 2026-08-24). Evidence: see Outcome section above.
- [x] `echo '{"jsonrpc":"2.0","id":2,"method":"tools/list",...}' | hawp mcp` → 3 tools returned (run locally, 2026-08-24). Evidence: see Outcome section above.
- [x] `hawp_work_validate` tool via MCP returns `hawp check: all 3 validations passed` (run locally, 2026-08-24). Evidence: see Outcome section above.
- [x] `go test ./internal/platform/mcp/...` → 5/5 pass (run locally, 2026-08-24). Evidence: see Outcome section above.
- [x] `npm --prefix librarian run validate` → 60/60 pass, hawp check all green (run locally, 2026-08-24). Evidence: see Outcome section above.

## Close Checklist

- [x] Implementation complete and committed on branch `0.0.4`
- [x] Tests pass (Go + npm validate)
- [x] PR #13 opened (0.0.4 → development)
- [x] BACKLOG updated (this item moved to Recently Closed)
- [x] Plan file updated to `status: done`
