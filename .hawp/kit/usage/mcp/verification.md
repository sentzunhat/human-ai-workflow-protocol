# Verify And Extend Provider Setup

## Current Support

| Provider | HAWP MCP setup behavior | Connection evidence |
| --- | --- | --- |
| Codex | Writes project TOML | Live validation passed in this repository |
| Claude Code | Writes project JSON | User-reported MCP PASS on 2026-09-06; not independently observed |
| Cursor | Writes project JSON | Config tests pass; client test still required |
| Continue | Prints manual configuration advice | Client test still required |
| GitHub Copilot | Prints manual setup advice | User subsequently reported actual MCP PASS; not independently observed |

This is a maintenance snapshot, not a guarantee for other machines or versions.
The user accepts Cursor/Continue live testing as deferred for this release.
Copilot's successful reported tool was `mcp_hawp2_hawp_work_validate`; the
client prefix differs without changing the underlying HAWP tool identity.
The five names identify client integrations, not model or infrastructure backends.
`all` retains its four-provider meaning: Claude, Cursor, Continue, Codex.
GitHub advice is explicitly selected. Unknown names fail in the configuration
writer before its writes; this does not make the entire init command atomic or
prevent earlier provisioning/network operations.

## Connection Check

1. Confirm the intended repository and compatible binary version.
2. Inspect the client's server status and approval requirements.
3. In that client, ask: "Call hawp_work_validate through MCP. Report the
   repository path and PASS/WARN/FAIL. Do not create or modify work items."
4. Compare the returned root with the workspace. Stop on mismatch or tool error.
5. Ask for a small `hawp_search` query after checking index availability.
   A missing index is a search setup issue, not necessarily a transport failure.
6. Test creation only when authorized, using a real new task or a disposable
   fixture. Inspect the UUID plan and backlog, then validate again.

Do not count CLI fallback as MCP proof. A server listing proves configuration
visibility, while a successful tool call proves more of the actual connection.
Keep each client's result separate. See the [setup guides](README.md) for
Codex, Claude Code, and VS Code Copilot commands and official references.

## Keep Extension Simple

Keep the current five for v0.0.24. Improve configuration preservation and prove
connections before adding another provider. No plugin framework is needed.

For a new client:

- Confirm its MCP transport, config format, execution host, and approval model
  from official documentation.
- Reuse the native HAWP server and explicit repository-root argument.
- Add a small client-specific configuration adapter and an explicit supported
  provider name; reuse JSON merging only when the schema actually matches.
- Decide deliberately whether it belongs in `all`; test ordering, duplicate
  selection, and unknown-name rejection.
- Test preservation of unrelated settings, malformed-input refusal, paths with
  spaces, idempotency, platform filenames, and a client-side tool call.
- Add its setup guide. Provider instruction overlays and distribution packs
  are separate contracts and need their own review when included.

Current configuration dispatch is a small explicit switch, not a dynamic plugin
system. The MCP protocol provides reuse across clients; configuration differences
still require adapters. No throughput or concurrent-work-item safety claim is
established by adding more client names.
