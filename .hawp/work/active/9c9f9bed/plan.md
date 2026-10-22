# Add structured output to MCP validation responses

**Backlog ID (Legacy):** — (UUID-native item)
**UUID:** `9c9f9bed-b08c-4e33-9005-66e60ac48348`
**Type:** improvement
**Reported:** 2026-10-17
**Status:** `plan-ready`
**Target Version:** TBD; candidate: latest official Go SDK with MCP `2026-07-28` plus its legacy handshake path, pending client verification

---

## Input (verbatim)

> Add an additive machine-readable result for the existing read-only hawp_work_validate MCP tool while preserving current text output, names, inputs, and CLI behavior. Clarify existing tool descriptions for intent, prerequisites, limits, and effects. Keep validation logic in existing application services and response mapping at the MCP boundary. Define client compatibility and contract verification.

## Intake Summary

Additive machine-readable output for the existing read-only validation tool,
while retaining its human-readable text and identity. This is a response
contract change, not a new tool.

## Current Context

The HAWP server currently advertises MCP protocol `2024-11-05`. The official
`2024-11-05` tool-result shape includes `content` and `isError`; structured
content and output schemas appear in the `2025-06-18` tool specification.
Therefore a standards-based `structuredContent` response is not a
compatibility-safe field-only addition under the currently advertised revision.
The source code already has application-level validation results, and the MCP
boundary currently returns the composite human-readable check output.

The maintainer's direction is to use bleeding-edge MCP and SDK capabilities
while keeping Codex, Claude, and GitHub usable. The current server is a small
hand-written MCP implementation; `2026-07-28` changes lifecycle behavior as
well as tool results, so this is a protocol/SDK migration prerequisite rather
than a response-field-only change. The official Go SDK supports modern
`2026-07-28` and legacy initialize-based clients through one server API. Thus
using the latest SDK and offering the latest protocol does not require dropping
legacy client compatibility. Modern-only support is a separate, stricter
product choice and is not required to adopt the latest SDK.

## Initial Analysis

**Directly verified:**

- `server/server.go` advertises `2024-11-05` and `server/tool_work.go` returns
  `hawp_work_validate` as text via `application/check.Run`.
- HAWP application work validation returns a domain `Report`; CLI rendering is
  separate. MCP can map results at the platform edge without relocating
  validation rules.
- Official MCP spec `2024-11-05` shows tool results with `content` and
  `isError`; official `2025-06-18` adds `outputSchema` and `structuredContent`.
- Existing names and input arguments need no change for this response feature.
- The official `2026-07-28` Tools specification supports `outputSchema` and
  `structuredContent`; structured results should also be serialized into a
  text content block for backward-compatible result consumers.
- The official `2026-07-28` protocol removes the initialize handshake and
  carries protocol metadata on every request; modern servers implement
  `server/discover`.
- The official Go SDK documents support for both the `2026-07-28` stateless
  lifecycle and legacy handshake versions.
- The official Go SDK lifecycle guide documents configuring modern
  `2026-07-28` alongside legacy `2025-11-25`, while using the same application
  APIs; legacy requests continue through `initialize`, modern requests use
  `server/discover` and per-request protocol metadata.
- The latest stable Go SDK version returned by `go list -m ...@latest` during
  this investigation is `github.com/modelcontextprotocol/go-sdk v1.8.0`.
- A temporary Go SDK v1.8.0 stdio probe restricted to `2026-07-28` passed
  `server/discover`, tool listing, and an echo tool call returning both
  serialized text and `structuredContent`.
- Local Codex CLI `0.163.0-alpha.10.2` did not use the modern lifecycle against
  the probe: its captured wire traffic sent `initialize` with
  `protocolVersion: 2025-06-18`, then `tools/list`; it sent no
  `server/discover` request or per-request protocol metadata. The Go SDK's
  modern-only configuration answered the legacy initialize with
  `2025-11-25`, per the SDK's documented fallback behavior, so this probe did
  not prove a modern-only server would accept Codex.
- Codex emitted an in-progress `probe_echo` call, but no `tools/call` request
  reached the probe process during roughly 40 seconds. A second attempt with
  Codex approval auto-review enabled behaved the same way. The reason for the
  client-side dispatch stall remains unknown; treat the Codex check as failed
  for modern-only support and separately unresolved as a tool-dispatch issue.
- The official latest stable npm package was `@openai/codex 0.163.0`; it was
  installed only under `/tmp` for an isolated retest. That CLI completed the
  same `probe_echo` call and received both serialized text and
  `structuredContent`, so the earlier stall did not recur in this run. Its wire
  traffic still used `initialize` with `2025-06-18`, followed by
  `tools/list` and `tools/call`; it did not use `server/discover` or modern
  per-request protocol metadata. This proves legacy-compatible tool use, not
  modern protocol support.
- The app-bundled Codex CLI remains `0.163.0-alpha.10.2`. Its `codex update`
  command could not detect the installation method; the Codex app update check
  returned `busy`, which did not establish whether an app update is available.
- Local Claude Code `2.1.300` could not run the probe because its OAuth session
  had expired and could not be refreshed. No MCP compatibility result was
  obtained.
- Copilot CLI is not installed locally. GitHub Copilot's hosted client surface
  was not available for a direct smoke in this investigation.
- Public documentation found in this investigation does not establish that
  Codex, Claude, and GitHub's intended client surfaces all accept a
  `2026-07-28`-only server. GitHub's MCP debugging guide still demonstrates
  `initialize` with `2024-11-05`; generic MCP support documentation is not
  proof of latest-revision compatibility.

**Inferred (not yet proven):**

- Implementing `structuredContent` and offering `2026-07-28` requires moving
  the server lifecycle to the modern protocol and checking all intended
  clients; that is broader than a local response mapper. Latest SDK support
  can retain a legacy path for clients such as current Codex, so this need not
  be a modern-only migration.
- The official Go SDK is a candidate implementation base for that lifecycle;
  the latest stable SDK version should be checked and pinned when implementation
  begins. Replacing the hand-written MCP boundary is in scope only if it keeps
  the HAWP service boundary and tool contracts simple.
- Returning JSON-formatted text can preserve machine-readable data for clients
  that parse text, but does not provide the protocol's structured result field
  and is not equivalent to structured MCP output.

**Likely scope:**

- The dispatch stall from the app-bundled alpha did not recur with stable CLI
  `0.163.0`; its cause is not established. The stable CLI can invoke tools
  through the legacy lifecycle, but fails the modern-only criterion because it
  still negotiates through `initialize`.
- Authenticate Claude Code and repeat the smoke, and identify an available
  GitHub Copilot surface for a direct smoke. Do not claim compatibility from
  generic MCP support or server SDK documentation.
- Target the latest stable official Go SDK at implementation time and offer
  MCP `2026-07-28` alongside its legacy handshake path. Verify that modern
  clients use `2026-07-28` and that intended legacy clients, including Codex,
  remain usable through the SDK's compatibility path.
- Treat modern-only operation as a separate decision. Do not require every
  client to speak `2026-07-28` before adopting the latest SDK; only choose
  modern-only if all intended clients pass direct discovery and invocation
  verification and the maintainer explicitly selects that narrower policy.
- Add an additive structured result and output schema for
  `hawp_work_validate`, preserving readable text, its name and inputs, and the
  CLI output.
- Clarify existing MCP tool descriptions in the same boundary-only slice.
- Add protocol contract tests and a targeted client compatibility smoke.

## Risk + Review Gate

**Risk:** medium
**Gate:** not met for implementation. The Go SDK's modern-only probe passed.
Latest stable Codex CLI `0.163.0` can list and invoke through legacy
initialize, but did not negotiate MCP `2026-07-28`; the dispatch stall seen in
the app-bundled alpha did not recur, though its cause is unknown. Claude is
blocked by expired local authentication, and GitHub Copilot was not directly
available. These results do not block adopting the latest SDK with both
protocol paths, but direct client verification remains required before claiming
the migration works across the intended client set. No code change is implied
by this plan.

## Backlog + Plan Link

**Status now:** plan-ready
**Plan file:** work/active/9c9f9bed/plan.md

## Next Step

- [x] Investigation recorded above (required before planning)
- [x] Write or update the plan file
- [x] Move backlog status to plan-ready pending protocol-version review

## Proposed Contract

- **Tool:** existing `hawp_work_validate`; no new tool or input arguments.
- **Success:** preserve human-readable `content`; add object-valued
  `structuredContent` and matching `outputSchema` only when serving a protocol
  revision that supports them.
- **Data:** overall status; named validation checks; stable finding codes,
  severity, repository-relative path, and concise message. Keep absolute machine
  paths, user-home data, and secrets out of structured results.
- **Errors:** retain existing `isError` behavior and human-readable error text;
  do not mislabel validation FAIL as a transport failure if current semantics
  treat it as a tool result.
- **Architecture:** use existing application validation services; map to the
  MCP wire contract at the platform boundary. CLI text stays owned by CLI
  renderers. Avoid a cross-layer MCP dependency.
- **Compatibility:** preserve name and inputs; verify current client versions
  against the proposed negotiated protocol version. Text remains the fallback,
  but does not by itself prove structured-output support.

## Plan

1. [x] Verify current and proposed MCP tool-result contracts from official
   versioned protocol documentation.
2. [ ] Verify modern `2026-07-28` discovery and invocation in available modern
   clients, and verify tool listing/invocation through the legacy path in
   Codex, Claude, and GitHub's intended client surfaces.
3. [ ] Select and pin the latest stable official Go SDK available at
   implementation time; migrate the MCP boundary to offer `2026-07-28` plus
   the SDK-supported legacy handshake path.
4. [ ] Add the structured validation result and output schema to the existing
   `hawp_work_validate` tool; clarify existing tool descriptions without
   changing names, inputs, or CLI behavior.
5. [ ] Add application mapping and MCP boundary contract tests for pass, warn,
   fail, and operational error cases.
6. [ ] Smoke named versions of Codex, Claude, and GitHub clients for tool
   discovery, invocation, resource/prompt use where applicable, and structured
   result visibility. Record exact surfaces and versions.

## Sources

- [MCP 2024-11-05 tool specification](https://modelcontextprotocol.io/specification/2024-11-05/server/tools)
- [MCP 2025-06-18 tool specification](https://modelcontextprotocol.io/specification/2025-06-18/server/tools)
- [MCP 2026-07-28 tool specification](https://modelcontextprotocol.io/specification/2026-07-28/server/tools)
- [MCP Go SDK lifecycle documentation](https://go.sdk.modelcontextprotocol.io/protocol/)
- [GitHub Copilot MCP debugging guide](https://docs.github.com/en/copilot/how-tos/copilot-sdk/troubleshooting/mcp-debugging)
