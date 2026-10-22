# 2026-10-19 — HAWP MCP dual-era protocol paths and client verification

## Status Report

### Intent

Preserve the decision and evidence needed to continue the MCP protocol/SDK
review for active work item `9c9f9bed-b08c-4e33-9005-66e60ac48348`. The
maintainer clarified that HAWP must support both the bleeding-edge protocol
path and the legacy path used by existing clients.

### Current State

The active plan proposes migrating the MCP boundary to the latest official Go
SDK, exposing MCP `2026-07-28` and retaining the SDK's legacy initialize-based
path. Both paths are required outcomes and must be verified separately. This
is a planning decision; no server implementation change has been made in this
checkpoint.

The local checkout is on `feature/website`, at `db9b16b3065e0d4696440e556eb24520d51fda4f`
before this archival commit. The branch matched `origin/feature/website` at
inspection time. Other local work is present and excluded from this checkpoint
commit.

### What Was Inspected

- Existing October 15–17 HAWP status checkpoints and MCP work records,
  including the closed catalog `2eea565c`, closed tool design `1673de65`, and
  active protocol/structured-validation plan `9c9f9bed`.
- `.hawp/kit/usage/status-report.md` and repository archival conventions.
- Official Go SDK lifecycle documentation and the MCP `2026-07-28` tools
  specification.
- The local Codex CLI probe evidence recorded in the active plan.
- GitHub repository search and the local Git remote/branch state.

### What Changed

- Reframed the active plan's target from “latest protocol only if every
  client passes” to “latest official SDK with both modern and legacy protocol
  paths.” A modern-only server is now a separate choice and is not required
  for HAWP to use the newest SDK.
- Updated the verification gate to require direct evidence for both paths:
  modern `server/discover` plus invocation, and legacy `initialize` plus
  invocation. Record the client surface and version for each observation.
- Clarified the provider-install boundary: provider settings normally choose
  the executable and its arguments/environment; the client and server
  negotiate the MCP wire version. A client or server binary version is not
  itself proof of the negotiated protocol. Do not pass a protocol version as
  a generic server argument unless that particular implementation documents
  it.
- Created this continuation checkpoint. The repository validator was run
  after the plan change.

### What Was Directly Verified

- The latest stable official Go SDK documentation describes modern
  `2026-07-28` discovery and legacy handshake operation through one server API.
- The isolated stable Codex CLI `0.163.0` probe listed and invoked a tool
  through `initialize`; the response included `structuredContent`. Captured
  traffic did not include `server/discover` or modern per-request protocol
  metadata. This verifies the legacy path only.
- The earlier app-bundled Codex CLI `0.163.0-alpha.10.2` probe did not return
  a result within roughly 40 seconds and was stopped. This remains
  inconclusive; the stable CLI's successful legacy invocation does not explain
  the earlier stall.
- The HAWP validation command completed with `VALIDATION PASS`, zero issues,
  and one verification-clarity warning: 307 of 382 checks were evidence-backed,
  71 explicitly unproven, and 4 ambiguous.
- GitHub repository search identified `sentzunhat/human-ai-workflow-protocol`;
  local `origin` points to that repository. `feature/website` matched the
  remote branch at the inspected SHA before archival.

### What Remains Unproven

- The HAWP server has not yet been migrated to the Go SDK or verified to serve
  both paths.
- Codex has not demonstrated modern `2026-07-28` discovery/invocation.
- Claude Code could not be tested because its local OAuth session had expired;
  GitHub Copilot's intended client surface was unavailable for a direct smoke.
- The intended clients' modern and legacy behavior, resource/prompt support,
  and structured-result visibility still need client-specific evidence.
- The Codex app's update status was not established: the app update check
  returned `busy`, and the bundled CLI could not update itself because its
  installation method was not detected.

### Constraints

- Keep the public MCP tool names and privacy boundaries. The structured
  validation improvement remains read-only and additive; this decision adds no
  tool names or side effects.
- Preserve the existing service boundary and use SDK lifecycle support rather
  than hand-writing parallel protocol implementations.
- Keep provider/prompt distribution, usage measurement, resources/prompts, and
  unrelated local implementation work under their existing items.
- Do not claim modern compatibility from SDK documentation, server tests, or a
  legacy invocation. Do not claim provider-reported token counts from local
  estimates.
- Exclude machine-local paths, credentials, and unrelated repository changes
  from this public checkpoint.

### Help Wanted

Review the dual-path policy in `work/active/9c9f9bed/plan.md`, especially the
minimum client matrix that should pass before implementation is considered
ready.

### Suggested Next Step

Continue `9c9f9bed`: identify the modern and legacy protocol configurations
supported by the pinned latest Go SDK, then run separate instrumented client
smokes for `server/discover`/modern invocation and `initialize`/legacy
invocation. Resume with Codex first, then obtain usable Claude and GitHub
Copilot surfaces. Keep protocol negotiation evidence distinct from provider
installation settings and client binary version.

## Conversation Checkpoint Summary

### Strategic Impact

HAWP can adopt bleeding-edge MCP and SDK support without requiring every
provider client to upgrade in lockstep. The intended architecture serves
modern clients through `2026-07-28` and keeps a compatible legacy route for
clients still using initialize. Maintaining both paths is now an explicit
requirement; modern-only is not the default migration target.

The release remains small: the structured validation response and existing
tool-description clarifications stay within `9c9f9bed`. Provider distribution
improvements remain planned for a later release lane.

### Continuation State

- **Latest completed milestone:** dual-path protocol choice recorded in the
  active plan; HAWP plan validation passed.
- **Active work:** `9c9f9bed-b08c-4e33-9005-66e60ac48348`,
  `.hawp/work/active/9c9f9bed/plan.md`.
- **Next milestone:** prove both modern and legacy MCP paths with named client
  versions and surfaces, then implement the SDK migration and additive
  validation contract within the same bounded work item.
- **Open client gates:** Codex modern-path proof; Claude authentication and
  smoke; GitHub Copilot surface and smoke.
- **Existing related work:** catalog `2eea565c` and tool-improvement design
  `1673de65` are closed evidence. Provider/prompt distribution remains separate
  from this implementation lane.
- **Working-tree boundary:** many pre-existing MCP/provider/docs/work-item
  changes are present locally. Only this checkpoint and the directly related
  active plan are included in the archival commit.

### Memory Delta

The existing HAWP MCP memory group already records the dual-era Go SDK
direction, Codex's legacy-only probe result, provider-setting distinction, and
active IDs. This checkpoint consolidates the October 19 clarification that
both protocol paths are required, and records that the plan and validation
were updated. No `MEMORY.md` content is duplicated. A compact follow-up note
will be added under the authorized ad-hoc memory extension path.

### GitHub Archival

- **Related repository found:** Yes
- **Organization:** `sentzunhat`
- **Repository:** `sentzunhat/human-ai-workflow-protocol`
- **Branch:** `feature/website`
- **Checkpoint file:** `.hawp/work/status/2026/10/19/3626745c/status.md`
- **Initial archival commit:** `06d3a579` (checkpoint and active plan)
- **Push status:** successful; `origin/feature/website` matched this commit
  after push. The final metadata follow-up is recorded in repository history.

### Resume Summary

**Resume from:** active work item `9c9f9bed`, with dual modern/legacy protocol
support as a required migration outcome.

**Next objective:** verify both wire paths with intended clients, then deliver
the bounded Go SDK migration and structured read-only validation response.
