# Verify And Extend Provider Setup

## Current Support

HAWP's server currently supports **local stdio only**. It does not implement a
local HTTP endpoint, Streamable HTTP, OAuth, hosted/remote MCP, or per-user
authorization. A successful client-side tool call proves the configured local
stdio connection; it is not evidence of remote transport readiness.

| Deployment mode | HAWP status | Boundary |
| --- | --- | --- |
| Client-launched stdio | Supported | Runs with the local process user's permissions and selected repository root; not multi-tenant isolation. |
| Streamable HTTP on loopback | Unsupported | No HTTP listener or Origin validation is implemented. |
| Hosted/remote Streamable HTTP | Unsupported | No HTTP transport, OAuth resource-server flow, token validation, or repository/principal binding is implemented. |
| Legacy HTTP+SSE | Unsupported | No implementation; do not use as a remote fallback. |

MCP's authorization specification applies to HTTP-based transports when
authorization is implemented; it says stdio servers should obtain credentials
from the environment instead. Do not add credentials to MCP tool arguments or
pass downstream API tokens through HAWP. Before any future remote design, scope
read, create, update, and delete permissions separately, bind identities to
repository roots, require explicit approval for destructive operations, and
address Origin/DNS-rebinding, OAuth discovery SSRF, token audience validation,
least-privilege scopes, auditability, and secret/path redaction. See the
[MCP transport specification](https://modelcontextprotocol.io/specification/2026-07-28/basic/transports),
[MCP authorization specification](https://modelcontextprotocol.io/specification/2026-07-28/basic/authorization),
and [MCP security best practices](https://modelcontextprotocol.io/docs/2026-07-28/tutorials/security/security_best_practices).

| Provider | HAWP MCP setup behavior | Connection evidence |
| --- | --- | --- |
| Codex | Writes project TOML and allowlists all seven HAWP tools with `enabled_tools` | Codex MCP validation passed in this repository on 2026-10-16 |
| Claude Code | Writes project JSON | User-reported MCP PASS on 2026-09-06; not independently observed |
| Cursor | Writes project JSON | User-reported live `hawp_work_validate` and `hawp_search` tests passed in this repository |
| Continue | Reads `~/.continue/config.yaml` (`mcpServers`) | Live `hawp_work_validate` passed after adding the HAWP server; known link warnings remain |
| GitHub Copilot | Writes repo-local `.vscode/mcp.json` and generates `.github/agents/hawp.agent.md` with the seven HAWP tools | VS Code Copilot called `hawp-hawp_work_validate` through MCP on 2026-10-16 after session approval; kit validation passed (3 checks, 0 issues), work validation passed (0 issues, 1 warning), and links passed (150 Markdown files). The generated HAWP custom-agent tool scope is not yet independently verified |

This is a maintenance snapshot, not a guarantee for other machines or versions.
The user reported that Cursor passed the target conversation tests after
approving the project server. The successful calls were read-only and used the
provider's prefixed HAWP tool names. Continue's live call used the unprefixed
`hawp_work_validate` tool name exposed by its MCP bridge.
Copilot's successful reported tool was `mcp_hawp_hawp_work_validate`; the
client prefix differs without changing the underlying HAWP tool identity.
The live result reported the correct repository root, zero work issues, and
one warning; aggregate `hawp check` still failed its links check for two known
broken local Markdown links.
The five names identify client integrations, not model or infrastructure
backends. `all` selects Claude, Cursor, Continue, Codex, and GitHub Copilot.
Continue still prints a user-scoped configuration block because the HAWP
command does not rewrite personal settings. Unknown names fail in the
configuration writer before its writes; this does not make the entire init
command atomic or prevent earlier provisioning/network operations.

## Tool Availability and Approval

Tool availability controls which MCP tools an agent can select. Approval
controls whether a selected tool call may run. HAWP configures availability
allowlists only where the provider has a suitable generated surface:

| Provider | HAWP-managed availability | Approval behavior |
| --- | --- | --- |
| Codex | `.codex/config.toml` lists all seven names in `enabled_tools` | Codex's own tool-call approval controls remain in effect |
| VS Code GitHub Copilot | `.github/agents/hawp.agent.md` lists the seven tools as `hawp/<tool>`; select the HAWP agent | Copilot continues to ask for tool-call confirmation as configured by the user |
| Claude Code | Server config only; HAWP does not write an MCP availability allowlist | Do not confuse Claude permission/allow rules with tool visibility; permission rules can change approval behavior |
| Cursor | Server config only; HAWP does not write Cursor MCP admin allowlists | Cursor's tool allowlists can make selected tools run automatically, so HAWP leaves them to owner/admin policy |
| Continue | HAWP prints the user-scoped server snippet; no HAWP-managed tool list | Configure availability and approvals in the Continue client |

VS Code documents both MCP server configuration and a custom-agent `tools`
field for selecting tools. Its server trust and tool-call confirmation remain
separate controls. See [VS Code MCP server setup](https://code.visualstudio.com/docs/agent-customization/mcp-servers)
and [VS Code custom agents](https://code.visualstudio.com/docs/agent-customization/custom-agents).
Codex's allowlist is part of its MCP server configuration. See the
[Codex MCP guide](https://developers.openai.com/codex/mcp). Cursor documents
tool allowlists as permitting selected MCP calls to run automatically; HAWP
does not write them. See [Cursor MCP documentation](https://cursor.com/docs/mcp).

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

## Read-only search test

After `hawp_work_validate` succeeds, use this exact provider-conversation test:

```text
Now perform one additional read-only HAWP MCP test.

Call the HAWP search tool through MCP, preferably hawp_search or its
provider-prefixed equivalent, with:

query: "MCP provider configuration"
limit: 3
context: false

Do not use terminal commands or shell fallback.

Report:

1. The exact tool name used.
2. Whether MCP was used.
3. Number of results.
4. The source paths returned.
5. Whether any result contains an absolute home-directory path, credential,
   API key, token, or other personal information.

Do not modify files or work records.
```

Record the provider, exact tool name, result count, source paths, and leakage
assessment as direct evidence. Do not treat a server listing or CLI result as a
provider-side tool invocation.

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
