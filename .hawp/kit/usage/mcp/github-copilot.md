# GitHub Copilot In VS Code

Follow the [shared setup and worker instructions](README.md) first.
This guide covers a local VS Code workspace, not GitHub-hosted coding agents,
Copilot CLI, or remote containers. Those require host-specific setup.

## Setup

`hawp init --provider github` syncs the GitHub instruction overlay and writes
the repo-local MCP configuration plus a HAWP custom agent. Use
`hawp mcp configure --provider github` when you want MCP configuration without
provisioning or kit sync. It writes `.vscode/mcp.json` and generates
`.github/agents/hawp.agent.md` with the seven HAWP tools in its `tools` list.
The HAWP agent scopes tool availability; it does not auto-approve calls. If an
unmanaged `.github/agents/hawp.agent.md` already exists, setup stops and
preserves it for manual review.

Merge this entry into `.vscode/mcp.json`; retain other servers and settings:

```json
{
  "servers": {
    "hawp": {
      "type": "stdio",
      "command": "${workspaceFolder}/.hawp/bin/hawp",
      "args": ["mcp", "--repo-root", "${workspaceFolder}"]
    }
  }
}
```

Use `hawp.exe` on Windows. In a multi-root workspace, configure the intended
folder explicitly. A remote extension host needs a compatible binary and paths
on that host, not your local computer.

## Verify

Use this strict connection-test prompt:

```text
Discover the MCP tools available in this conversation. If hawp_work_validate
is exposed (possibly with a client prefix), invoke it through MCP and report
the exact tool name, repository path, and validation result.
Do not use terminal commands, scripts, tasks, or shell fallback for this test.
If the tool is unavailable, report MCP TOOL UNAVAILABLE and stop.
Do not create or modify work items or configurations.
```

The user-supplied 2026-09-06 report ran a terminal command, so it establishes
CLI validation only. `work validate` checks work integrity, while the MCP tool
also checks kit and links. The CLI supports `--work-root` or `--hawp-root`,
not `--repo-root`; older builds silently ignored that unsupported flag.

Run VS Code's `MCP: List Servers`, select HAWP, and review/start the server.
In a Copilot conversation, select the HAWP agent to expose the HAWP tool set.
Call `hawp_work_validate` and confirm the repository before authorizing writes.
Copilot's normal server trust and tool-call confirmation remain separate and
must be handled by the workspace owner.

VS Code uses `servers`, unlike Claude's `mcpServers`. Its Agent Host does
not read `.vscode/mcp.json` directly; VS Code forwards supported configurations.
See the [official configuration reference](https://code.visualstudio.com/docs/agents/reference/mcp-configuration)
and [custom-agent tool configuration](https://code.visualstudio.com/docs/agent-customization/custom-agents).
This guide describes setup; see [connection verification](verification.md)
for current client test evidence.
