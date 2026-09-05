# MCP Setup Lessons

Recorded 2026-09-06 from HAWP source, configuration tests, and local MCP calls.

- Configuration written, server connected, tool exposed, and tool executed are
  separate checks. Record each honestly; validate the selected repository first.
- Native executable paths plus explicit `--repo-root` avoid dependence on the
  client's working directory. Absolute local paths must not leak into shared
  repository artifacts.
- `init` provisions assets and syncs kit/provider files. It is not a lightweight
  MCP-only command. A later asset failure can coexist with an already-written
  provider config; examine both the error and resulting files.
- Codex and Claude config writers exist. GitHub currently prints manual MCP
  setup advice only. Provider instruction overlays are not MCP connections.
- Preserve custom environment, timeout, policy, and unrelated server settings.
  JSON and explicit Codex TOML table migrations have preservation tests. Codex
  migration uses parsed source ranges rather than replacing the whole table;
  unsupported layouts and embedded launch-value comments require manual review.
- Do not reuse one provider's JSON wrapper key indiscriminately: Claude uses
  `mcpServers`; VS Code workspace MCP configuration uses `servers`.
- Trust and approval remain user decisions. Do not bypass them to make a
  connection test pass.
- Local in-place executable replacement produced exit 137 during maintenance.
  Temporary-file replacement restored execution. The Makefile and CI now share
  that install path; the precise OS-level cause was not established.
- Search can be stale or incomplete; use the backlog and source files to verify.
  The current tools create and validate work, but do not update existing plans.
- Single-executable retirement and routing unmatched records into an unknown
  folder remain planned. Preserve copies, evidence, identity, and sidecars.

## Verification Boundary

Local configuration tests and actual HAWP MCP calls in Codex are evidence.
They do not prove every provider version, platform, user policy, or hosted agent
environment. Official provider documentation is linked in each setup guide.
No live Claude Code or Copilot connection was exercised for these guides.
