# MCP capability catalog for HAWP resources and prompts

**UUID:** `2eea565c-f56d-4139-b4eb-06eebfc2b456`
**Type:** improvement
**Reported:** 2026-09-13
**Status:** `closed`
**Target Version:** `v0.0.25`
**Successor:** [Design improved HAWP MCP tools](../../../active/1673de65-63e0-4a8e-bade-cbdde4095912/plan.md)

---

## Input

User shared the YouTube video `https://youtu.be/BqRhBq-_kgE?si=Kt0vrAV5XYL4269g`
and asked to use its MCP architecture lesson as learning for planning HAWP
improvements.

## Source Notes

- Confirmed video metadata via YouTube oEmbed: title `MCP Just Got a Whole Lot Better`, author `Neon Postgres`.
- Transcript retrieval was blocked by YouTube during intake. Treat this plan as
  derived from the confirmed video topic plus current MCP architecture themes,
  not as a transcript-backed summary.
- Current HAWP MCP guidance is mostly tool-first: `hawp_search`,
  `hawp_work_reshape`, `hawp_work_new`, `hawp_work_doc`, `hawp_work_validate`,
  and `hawp_usage`.

## Lesson

MCP is more useful when the server advertises the right kind of capability for
the job. HAWP currently exposes most behavior as tools, even when some material
is better modeled as discoverable, read-only context or reusable task templates.

## Release Scope

This is a `v0.0.25` lane item. Keep it out of the `v0.0.24` cutoff unless the
maintainer explicitly expands the release beyond the compound intake tool.

Version-scope note: see
`.hawp/work/status/2026/09/13/6059dd4b/status.md`.

## Plan

### Slice A - Capability and workflow inventory

- [x] Inventory every existing MCP operation and its behavior, inputs, outputs,
  side effects, filesystem scope, and suitable user workflow.
- [x] Complete the initial source-backed tool comparison below, including
  usage-log storage behavior and the intake/search/reshape overlap.
- [x] Classify capabilities as tools (discrete operations), resources
  (repository-scoped read-only context), prompts (reusable task templates), or
  documented workflows (multi-step recipes using one or more capabilities).
- [x] Map current workflows, including search/context retrieval, intake shaping,
  create-work, validation, and status/evidence documentation.
- [x] Review overlaps and possible consolidation, especially the composed
  `hawp_work_intake` path versus standalone `hawp_search` and
  `hawp_work_reshape`. Preserve standalone operations where they compose into
  other workflows; recommend merging only where it improves discoverability
  without hiding distinct behavior or adding surprise side effects.
- [x] Record mutation boundaries for create, update, and delete. Identify what
  exists today and what is absent; do not add update/delete tools as part of the
  inventory. Any future mutation proposal must define scope, confirmation,
  idempotency, and recovery expectations.
- [x] Identify initial high-value read-only resource candidates and their
  repository scope below.
- [x] Identify initial prompt candidates and their arguments below.
- [x] Map current MCP workflows and note which exist only as documentation.
- [ ] Identify high-value read-only resources, starting with kit start-here,
  backlog status, active-plan summaries, and provider setup guidance.
- [ ] Identify prompt candidates, starting with intake shaping, status report,
  release checkpoint, and provider connection verification.

#### Initial tool catalog (source review, 2026-10-16)

| Tool | Behavior and boundary | Side effects / prerequisites | Catalog recommendation |
| --- | --- | --- | --- |
| `hawp_search` | Searches the selected repository's HAWP index. Returns ranked chunks, or a deduplicated/token-capped context block when `context:true`. | Reads local SQLite index; requires an index. No repository mutation. | **Keep.** It is the general retrieval primitive; clarify when to use raw results versus context block. |
| `hawp_work_reshape` | Shapes supplied request/context into HAWP draft fields using local Ollama or ONNX backend. It does not retrieve repository context. | Invokes configured local model; no repository mutation. Backend availability required. | **Keep.** Distinct for user-supplied context and independent shaping. Clarify its stateless scope and distinguish from intake. |
| `hawp_work_intake` | Searches indexed repository context, shapes the request with retrieved context, and returns structured readiness, questions/blockers, warnings, and token accounting. Does not create a work item. | Reads local index and invokes local model; requires index and model. No repository mutation. | **Keep as the recommended compound workflow entry point.** It composes search plus reshape with explicit states, rather than duplicating their independent use cases. Consider clearer naming/description only if client discovery testing shows confusion. |
| `hawp_work_new` | Creates a new work plan and adds an inbox row to the selected repository's backlog. | Writes under `.hawp/work`; creates a new UUID; no update/delete operation. | **Keep, clearly label as a write.** Require callers to check for an existing matching item first; preserve explicit human decision to create. |
| `hawp_work_doc` | Creates a status, evidence, decision, or note document, optionally grouped under a work item ID. Repeated same type/work-item requests are documented as idempotent. | Writes a new file under canonical HAWP secondary-document paths; no existing content update or delete. | **Keep, clarify create-only semantics.** It is narrower than generic file write and preserves path conventions. Review idempotency and collision behavior before considering extensions. |
| `hawp_work_validate` | Runs repository HAWP structure, work-record, and link checks. | Reads repository state; may produce validation output, but no intended work-record mutation. | **Keep.** It provides a distinct verification capability and repository-scope confirmation. |
| `hawp_usage` | Returns local MCP call totals or a Markdown report with per-tool breakdown and recent query list. | Reads user-home HAWP usage config/database; available only when usage logging is enabled. Even with full-body logging disabled, the store keeps up to 256 characters from a `query` or `title`; `--log-bodies` additionally stores raw request and response bodies. | **Keep as local opt-in analytics, with an explicit privacy note.** Do not promote the recent-query report without explaining that queries/titles are retained. Review retention/clear guidance before broadening exposure. |

**Merge decision:** Do not merge `hawp_search`, `hawp_work_reshape`, and
`hawp_work_intake`. The first two are independent primitives; intake is an
opt-in composition with retrieval, model invocation, readiness states, and
accounting. Removing the primitives would make targeted retrieval or shaping
with caller-supplied context less direct. Improve descriptions and workflow
guidance instead; retain existing tool names for compatibility.

**Maintainer decision (2026-10-16):** Keep the existing MCP tool names and
group them in documentation/client guidance. Do not remove, rename, or replace
the current tools with grouped wrapper tools. The groups are descriptive:
Discover, Shape, Create, and Verify. This preserves provider allowlists and
existing client configurations while improving discoverability.

**Protocol implementation decision (2026-10-16):** Keep the handwritten local
stdio server for this bounded change. The retained MCP protocol version
`2024-11-05` already defines resources and prompts, so HAWP can add their
standard list/read and list/get methods without introducing a second server
stack or making an SDK migration part of this feature. The official Go SDK is
an implementation option, not proof of HAWP support; the new handlers still
require protocol tests and client verification.

**Mutation scope observed:** `hawp_work_new` creates plans/backlog entries and
`hawp_work_doc` creates secondary documents. No MCP update or delete tool is
present in the inspected server dispatch. Updating existing plans is currently
documented as direct file editing by the client/agent; deletion has no exposed
MCP operation. Any future update/delete proposal remains separate and must
specify exact allowed paths/fields, confirmation, idempotency, recovery, and
audit behavior. Usage storage was inspected: ordinary opt-in logging stores a
query/title summary capped at 256 characters; `--log-bodies` stores raw input
and output bodies too. Existing CLI help documents opt-in, disable, and clear
commands, but the MCP tool description only says logging must be enabled. Add a
privacy/retention note to MCP usage guidance before recommending report mode.
The remaining accurate metering and retention work is tracked in the resumed
`e1dc28ac` (legacy ID `usage-tracking`); do not duplicate it here. Intake elicitation/client
support is tracked separately in `c8d9e1cb`.

#### Current workflow map

| Workflow | Current capability sequence | Notes / decision boundary |
| --- | --- | --- |
| Repository context lookup | `hawp_search` (raw results or `context:true`) | Client verifies source files before treating hits as evidence. The index can be stale. |
| Request shaping only | `hawp_work_reshape` | Uses only caller-supplied input/context; no repository retrieval and no write. |
| Structured intake | `hawp_work_intake` -> inspect state/questions -> user clarification or retry as needed | Compound retrieval and shaping; no work item is created by this call. |
| New work intake | `hawp_work_intake` -> inspect backlog/plans and avoid duplicate -> user decides create -> `hawp_work_new` | Tool response does not itself authorize implementation. Existing matching work should be continued. |
| Work update | client/agent reviews plan -> directly edits existing plan/backlog -> `hawp_work_validate` | No MCP update tool exists; retain UUID and existing evidence. |
| Checkpoint/evidence note | `hawp_work_doc` -> write verified content into returned path -> `hawp_work_validate` | MCP tool creates the document path; caller supplies/appends its content outside this tool. |
| Repository validation | `hawp_work_validate` | Read/validation operation; does not require search index. |
| Usage inspection | optional `hawp_usage` totals/report when local logging is enabled | User-home local data; report includes recent query summaries. Body logging is separately opt-in. |
| Work deletion | none | No MCP capability; deleting work is intentionally outside this catalog's implementation scope. |

The local stdio MCP server advertises tools, resources, and prompts while
retaining protocol version `2024-11-05`. It handles `tools/list`, `tools/call`,
`resources/list`, `resources/read`, `prompts/list`, and `prompts/get`.
“Workflow” remains a documented sequence, not a separate live MCP capability.

#### Full capability and operation inventory (source review, 2026-10-17)

This inventory describes the checked-in server and CLI source. Provider setup
and client-side workflow docs are guidance; they do not add server operations.

| Class / operation | Inputs and output | Effects and scope | Workflow / boundary |
| --- | --- | --- | --- |
| Tool `hawp_search` | `query` required; optional `limit` (default 5, max 500), `context`, `max_tokens` (default 2000, max 32000). Returns ranked source chunks/line positions or a deduplicated token-capped context block. | Reads the selected repository's local search index; no repository writes. Requires an index; embedding-backed ranking is used when vectors exist. | Discover context; verify source files because indexed results may be stale. |
| Tool `hawp_usage` | Optional `report` boolean. Returns totals or a Markdown breakdown and recent query summaries. | Reads local user HAWP config/database outside the selected repository. Only available when opt-in logging is enabled. Logging can store query/title summaries; `--log-bodies` stores raw bodies. | Inspect local usage. Treat report output as potentially sensitive; retention is manual. |
| Tool `hawp_work_intake` | `input` required; optional backend/model/url/limit/max_tokens. Returns structured draft, retrieval state/warnings/questions, and token estimates. | Reads selected repo index and invokes configured local Ollama or ONNX model; no work-record write. | Compound retrieval + shaping. Missing index/results or shaping failures return blocked/needs-input states; does not authorize creation. |
| Tool `hawp_work_reshape` | `input` required; optional caller-provided `context`, backend/model/url. Returns shaped mission, constraints, output, checkpoint. | Invokes configured local Ollama or ONNX model; stateless with respect to repository, no write. | Shape caller-supplied context independently of repository search. |
| Tool `hawp_work_new` | `title` required; optional type/input. Returns generated short UUID, plan path, and backlog path. | Creates `.hawp/work/active/{uuid}/plan.md` and adds inbox row to selected repo `.hawp/work/BACKLOG.md`. | Explicit create step after duplicate/backlog review and user decision. No update/delete operation. |
| Tool `hawp_work_doc` | `type` and `title` required; optional `work_item_id`. Returns canonical secondary-document path. | Creates a status/evidence/decision/note path under selected repository `.hawp/work`; repeated same type/work-item calls return existing path without overwrite. Caller fills content separately. | Create checkpoint/evidence artifact; validate afterward. It is not an arbitrary file-write tool. |
| Tool `hawp_work_validate` | No arguments. Returns PASS/WARN/FAIL findings for kit structure, work records, and local Markdown links. | Reads selected repository; no intended record mutation and no search-index prerequisite. | Verify repository/work state before or after workflow changes. |
| Resource `hawp://repo/kit/start-here` | No arguments; `resources/read` returns Markdown contents. | Fixed `.hawp/kit/start-here.md` in selected repository; bounded regular-file read (1 MiB maximum), path/symlink/hard-link checks. | Read operating guidance. |
| Resource `hawp://repo/work/backlog` | No arguments; `resources/read` returns Markdown contents. | Fixed `.hawp/work/BACKLOG.md` in selected repository; same bounded safe-read behavior. | Read current work index; content is repository data, not server instruction. |
| Resource `hawp://repo/kit/usage/mcp` | No arguments; `resources/read` returns Markdown contents. | Fixed `.hawp/kit/usage/mcp/README.md` in selected repository; same bounded safe-read behavior. | Read MCP capability/setup/workflow guidance. |
| Prompt `hawp_intake_request` | Required `request`, optional `context`; `prompts/get` returns one user text message. | No repository reads or writes by the prompt handler; arguments are JSON-encoded and framed as user data. | Draft HAWP intake fields, preserve unknowns, ask focused questions; does not create or approve work. |
| Documented workflow: configure/connect | `hawp mcp configure --provider ...` (Codex, Claude, Cursor, Continue, GitHub; `all` selects supported set). | CLI configuration writes or prints provider config; no MCP server method. Provider-specific files/permissions differ; setup exit success is not live connection proof. | Configure, restart client, then perform client-side MCP smoke and verify repository root. |
| Documented workflow: intake to work artifact | `hawp_work_intake` -> inspect readiness/backlog -> user decision -> `hawp_work_new` -> implement -> `hawp_work_doc` -> `hawp_work_validate`. | Sequence crosses read/model operations, explicit work creation, and secondary-document creation. | Use `.hawp/kit/examples/mcp-intake-to-work-doc.md`; decision points stay visible, no all-in-one mutation. |
| Documented workflow: direct retrieval/shaping/update | `hawp_search` and/or `hawp_work_reshape`; update existing plan/backlog directly in client/agent, then validate. | Search and reshape are read/model operations. Existing-record edits are outside MCP tools and performed by the client/agent. | Keep standalone operations composable; no MCP update/delete capability exists. |
| Documented workflow: usage administration | CLI `hawp usage enable|disable|report|log|clear|totals`; MCP `hawp_usage` only reads totals/report. | CLI toggles local logging, reads reports, or clears local usage log; MCP tool does not toggle or clear. | Opt-in local telemetry; communicate summary/body retention before enabling or sharing reports. |

**Server dispatch boundary:** observed MCP methods are `initialize`,
`tools/list`, `tools/call`, `resources/list`, `resources/read`,
`prompts/list`, and `prompts/get`. Notifications without IDs receive no
response. No resource templates, subscriptions, roots negotiation, sampling,
elicitation, logging capability, HTTP transport, MCP update, or MCP delete
operation is advertised/handled in this server. This is a source observation,
not a claim that every provider surfaces every implemented capability.

**Inventory classification:** tools are discrete search/shape/create/verify
operations; resources expose three fixed repository-scoped read-only files;
the prompt is one static intake template; workflows are documented multi-step
recipes, not a fifth protocol capability. Provider configuration and usage
administration are CLI operations documented alongside MCP, not MCP methods.

#### Candidate read-only resources (proposal)

Prefer a small stable surface whose content is generated from the selected
repository at read time. Candidate URI scheme is `hawp://repo/...`; exact URI
syntax should be checked against target clients and the selected SDK before
implementation.

| Candidate | Example URI | Content and scope | Priority |
| --- | --- | --- | --- |
| Operating guide | `hawp://repo/kit/start-here` | `.hawp/kit/start-here.md`, selected repo only | High |
| Backlog summary | `hawp://repo/work/backlog` | Compact `.hawp/work/BACKLOG.md`; expose current active/parked/closed index as authored | High |
| Active work plan | `hawp://repo/work/active/{id}` | One explicitly selected active plan; do not enumerate or read arbitrary paths | High |
| MCP capability/setup guide | `hawp://repo/kit/usage/mcp` | MCP README and bounded provider/setup guidance | Medium |
| Status report template | `hawp://repo/kit/usage/status-report` | The canonical status-report guide | Low/medium; may duplicate operating context |

Resources must be read-only, repository-contained, reject traversal/symlink
escapes, avoid exposing unrelated home-level usage data or secrets, and report
missing/stale source state clearly. Backlog and plan contents remain untrusted
repository data, not server instructions. Prefer summaries only if they can be
produced deterministically and linked to their source; avoid maintaining a
second stale catalog.

#### Candidate prompts (proposal)

| Prompt candidate | Arguments | Suggested capability sequence | Boundary |
| --- | --- | --- | --- |
| `hawp_intake_request` | `request` (required), optional `context` | Read start guide/backlog (future resources) -> suggest `hawp_work_intake` or `hawp_work_reshape` -> ask user when uncertain | Shapes a request; never creates an item or implies approval. |
| `hawp_review_work` | `work_item_id` (required) | Read selected plan/backlog -> review evidence/status -> optionally suggest `hawp_work_validate` | Review instructions only; no writes. |
| `hawp_status_report` | `work_item_id` (required), `focus` (optional) | Read plan/context -> draft report using status guide -> suggest `hawp_work_doc` only after user asks to save | Draft only; saving remains an explicit tool call. |
| `hawp_release_checkpoint` | `version` (required), `scope` (optional) | Read release plans/evidence -> summarize gates and unknowns -> suggest validation | No tag, publish, commit, or release action. |
| `hawp_mcp_connection_check` | `provider` (optional) | Follow verification guide -> suggest `hawp_work_validate` then a bounded `hawp_search` smoke | Does not configure clients or claim a live connection without direct result. |

Do not implement a prompt if it merely repeats a concise workflow guide. Prompt
templates should be small, identify arguments, preserve user decision points,
and avoid embedding full kit documents. Dynamic repository state belongs in
resources or tool outputs, not copied static prompt text.

#### Client and SDK compatibility research (2026-10-16)

Official documentation confirms that the proposed capability types are
available in some target clients, with important client-surface differences:

| Target | Directly documented support | Boundary |
| --- | --- | --- |
| Cursor | Its MCP guide lists prompts as supported and describes MCP servers as tools/data sources. | The inspected page does not establish a resource reference workflow. Documentation-level evidence only; no HAWP smoke was run. |
| Claude Code | Anthropic documents resources via `@` references and server prompts as slash commands. | Documentation-level evidence only; no HAWP smoke was run. |
| GitHub Copilot in VS Code | GitHub documents MCP prompts via slash commands and resources via Add Context. | This applies to the VS Code surface; do not generalize to Copilot cloud agent or code review. No HAWP smoke was run. |
| Continue | Current config reference documents MCP context-provider/server configuration. | Inspected docs do not establish server-provided resource/prompt discovery UX; no HAWP smoke was run. |
| Codex local clients | Current docs enumerate stdio, Streamable HTTP, server instructions, and authentication features. | Inspected supported-feature list does not establish runtime MCP resource/prompt UX; no HAWP smoke was run. |

The official Go SDK supports registering server resources and prompts, but
HAWP does not depend on that SDK. HAWP's handwritten server advertises protocol
`2024-11-05` and now handles tools, fixed read-only resources, and one intake
prompt. Focused protocol tests pass. These checks establish HAWP server
behavior, not client rendering or HAWP-client compatibility; do not infer those
from the SDK or vendor feature lists.

Sources: [Cursor MCP protocol support](https://docs.cursor.com/context/model-context-protocol),
[Claude Code MCP](https://docs.anthropic.com/en/docs/claude-code/mcp),
[GitHub Copilot Chat MCP](https://docs.github.com/en/copilot/how-tos/copilot-in-your-ide/customize-copilot/extend-copilot-with-tools-and-context/extend-copilot-chat-with-mcp),
[Copilot cloud-agent boundary](https://docs.github.com/en/copilot/concepts/agents/cloud-agent/mcp-and-cloud-agent),
[Continue configuration reference](https://docs.continue.dev/reference),
[official Go SDK server guide](https://github.com/modelcontextprotocol/go-sdk/blob/main/docs/server.md).
These are documentation observations, not HAWP integration evidence. See
`core/providers/doc-verification.md` for the current vendor-doc capability
matrix and the explicit client smoke gap.

#### Capability grouping and naming recommendation

- Keep the seven existing tools as stable operations. Group them in user-facing
  documentation under **Discover** (`hawp_search`, `hawp_usage`), **Shape**
  (`hawp_work_reshape`, `hawp_work_intake`), **Create** (`hawp_work_new`,
  `hawp_work_doc`), and **Verify** (`hawp_work_validate`). These are catalog
  categories, not a new MCP namespace or wrapper tool.
- Do not combine create and document creation: they have distinct schemas,
  paths, and lifecycle meanings.
- Do not add a generic `hawp_work_update` or delete operation as part of this
  catalog. If a future need is approved, scope updates to typed HAWP fields and
  canonical records rather than arbitrary paths; define previews, approval,
  conflict handling, backup/recovery, and audit behavior first. Delete requires
  an explicit retention/recovery decision and separate approval.
- Keep multi-step flows as documentation/prompts where the client must make
  decisions or verify results. Avoid a monolithic “do everything” tool that
  hides search results, model uncertainty, user clarification, or write steps.

### Slice B - Capability map and thin resource/prompt/workflow design

- [x] Propose initial repository-scoped URI names for read-only HAWP resources.
- [x] Propose initial prompt names and arguments that never perform writes
  implicitly.
- [x] Map named end-to-end workflows, decision points, and current fallback
  behavior.
- [x] Produce an initial capability matrix/recommendation for existing tools,
  resources, prompts, and workflow guidance.
- [ ] Validate proposed resources and prompts against target client/SDK support
  and URI/content constraints before implementation selection.
- [x] Review the catalog with the maintainer before implementation selection;
  the maintainer authorized a minimal local stdio implementation on 2026-10-16.
- [x] Decide whether the proposed resource/prompt split is worth implementing:
  proceed with three fixed read-only resources and one intake prompt, with
  client verification still required before distribution/provider claims.
- [x] If planning publication, update the MCP setup guides with the reviewed
  capability map and privacy note for usage logging.
- [x] Decide whether grouping requires tool renames or removals: maintainer
  chose to retain names and group them in guidance; no wrapper tools or
  deprecation path are needed.
- [x] Keep tools for discrete operations and mutations; do not turn resources,
  prompts, or workflow descriptions into hidden write paths.
- [x] Treat future HTTP exposure as out of scope for implementation here; carry
  forward the security requirements in
  `.hawp/work/closed/2026/10/16/429e075e/plan.md`.

### Slice C - Implementation candidate

- [x] After review, select a small, independently useful implementation slice:
  fixed read-only repository resources and one intake prompt. Detailed tracking
  is in `3cf52223` and `9dce1d95`.
- [x] Start with read-only resources and a prompt on the existing stdio path;
  do not add arbitrary repository paths, resource templates, or HTTP transport.
- [ ] If consolidation is recommended, preserve backward compatibility or
  document an explicit migration/deprecation path before changing tool names.
- [x] Update `.hawp/kit/usage/mcp/README.md` with the resource URI contract and
  intake prompt behavior. Provider packs and generated distributions are
  deferred to `3c137bad` until client behavior is verified.

## Risk

Medium. This changes how clients discover HAWP, but the first slice is
read-only planning and should not alter existing tools.

## Verification

- Capability inventory cites exact source files and distinguishes observed
  behavior from design proposals.
- Catalog identifies the purpose and side effects of every current tool and
  evaluates merge/group options without presuming fewer tools are better.
- Proposed prompts, resources, and workflows have clear boundaries and no
  implicit write paths.
- Create/update/delete inventory distinguishes existing behavior from absent
  capabilities and proposed future work.
- New resources, if implemented later, are read-only and scoped to the selected repository.
- MCP tool list remains backward compatible.
- `go test ./...`, `go vet ./...`, and `go run ./cmd/hawp work validate` pass from `librarian/src`.


## Outcome and Closeout

The source-backed inventory and workflow classification are complete and retained here as evidence. The earlier framing focused on cataloging and backward compatibility; the maintainer clarified that the desired outcome is concrete proposals for better MCP tools. That design is split into successor work item `1673de65` so tool design and any later implementation remain independently scoped. This item closes the catalog/evidence phase only; it does not claim that the MCP tool set has been improved.

- **Status:** `closed` as of 2026-10-17.
- **Guidance update (2026-10-16):** `.hawp/kit/usage/mcp/README.md` now groups
  the existing tools as Discover, Shape, Create, and Verify without renaming
  or wrapping them. It explains intake versus standalone search/reshape and
  documents opt-in usage logging, the 256-character query/title summary,
  optional raw-body capture, report exposure, and manual-only retention.
- **Resume context:** `c147f2ae` tool-first audit is closed. The remote transport
  and authorization audit `429e075e` is also closed and is a prerequisite for
  any future HTTP MCP implementation, not for this local stdio catalog.
- **Continuation:** Successor item `1673de65` will use this inventory and the
  tool-first audit `c147f2ae` to propose and prioritize concrete tool changes.
  Client discovery/read/prompt verification remains separately tracked.
- **Implementation slices:** `3cf52223` adds fixed read-only resources;
  `9dce1d95` adds one intake prompt. Provider/distribution guidance is being
  propagated under active item `3c137bad`, with client use conditional on the
  capability appearing in that connected client.
- **Scope boundary:** This catalog phase added no MCP tool capabilities or HTTP
  transport. Proposed tool changes belong to successor item `1673de65` and need
  their own review and implementation scope.


**Close date:** 2026-10-17

## Close Checklist

- [x] Source-backed capability and workflow inventory recorded.
- [x] Follow-on tool-design work split into successor item `1673de65`.
- [x] Archived plan and removed active backlog row.
