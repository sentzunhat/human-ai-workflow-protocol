# Interactive MCP intake refinement with structured results

**UUID:** `c8d9e1cb-2189-45d4-bd4c-7cbc51710822`
**Type:** feature
**Reported:** 2026-09-13
**Status:** `in-progress`
**Target Version:** Structured contract shipped in `v0.0.24`; remaining elicitation/client review is post-`v0.1.0` planning

---

## Input

User shared the YouTube video `https://youtu.be/BqRhBq-_kgE?si=Kt0vrAV5XYL4269g`
and asked to use its MCP architecture lesson as learning for planning HAWP
improvements.

## Source Notes

- Confirmed video metadata via YouTube oEmbed: title `MCP Just Got a Whole Lot Better`, author `Neon Postgres`.
- Transcript retrieval was blocked by YouTube during intake. Treat this plan as
  MCP-architecture inspired, not as a transcript summary.
- Active plan `a8797f44` already proposes `hawp_work_intake`, combining search
  and reshape in one tool call.

## Lesson

The sexier MCP shape is not only "one more tool." It is a guided interaction:
structured results, explicit missing-information handling, and a way to keep
humans in the loop without making the agent guess.

## Plan

### Slice A - Structured result contract

- [x] Extend the planned `hawp_work_intake` result design with structured fields:
  `draft`, `retrieval`, `questions`, `warnings`, and `token_accounting`.
- [x] Add result states such as `ready_for_work_new`, `needs_user_input`, and
  `blocked_missing_index`.
- [x] Add `blocked_reshape_failed` for local model/parser failures after
  retrieval succeeds.
- [x] Make missing context visible instead of silently returning a confident-looking draft.

`v0.0.24` cutoff: include only the smallest structured-state contract needed to
make `hawp_work_intake` safe and clear. Defer richer elicitation, guidance, and
client-behavior work to `v0.0.25`.

`v0.0.24` subset landed in `a8797f44` on 2026-09-13. Keep this record open for
the remaining `v0.0.25` elicitation and guidance work.

### Slice B - Elicitation research

- [ ] Research current MCP elicitation/client support and decide whether HAWP can
  request focused user input during intake.
- [ ] If support is uneven, design a fallback where the tool returns concise
  questions for the agent to ask manually.
- [ ] Keep user-confirmation distinct from implementation approval.

Target version: `v0.0.25`.

### Slice C - Agent ergonomics

- [x] Add examples showing the complete path:
  `hawp_work_intake` -> optional user answer -> `hawp_work_new` -> `hawp_work_doc`.
- [x] Add tests for missing search index, empty/no-match context, structured
  success, and token accounting in the `v0.0.24` intake implementation.
- [ ] Add remaining `v0.0.25` tests for ambiguous requests and client guidance.
- [x] Update MCP worker guidance to prefer the structured status over prose parsing.

Target version: `v0.0.25`, except tests directly required by the `v0.0.24`
`hawp_work_intake` implementation.

## Risk

Medium. This depends on client support differences, so the first implementation
should work even when elicitation is unavailable.

## Verification

- Tests cover structured success, needs-input, and missing-index responses.
- Existing text responses remain understandable for clients that do not expose structured content.
- Agent guidance demonstrates no-guessing behavior for missing context.
- `go test ./...`, `go vet ./...`, and `go run ./cmd/hawp work validate` pass from `librarian/src`.


## Resume update (2026-10-16)

### Directly verified complete

- Structured intake states, retrieval metadata, warnings/questions, and token
  accounting are implemented in `librarian/src/internal/platform/mcp/server/tool_work.go`
  and `types.go`.
- The server returns concise manual questions in `needs_user_input`,
  `blocked_missing_index`, and `blocked_reshape_failed` states. This is the
  fallback when no interactive client elicitation path is available.
- The intake-to-work-document example and MCP worker guidance already describe
  the user-answer/manual fallback and separate work creation from approval.
- Repository tests cover ready, no-match, missing-index, reshape-failure, and
  invalid numeric option behavior in
  `librarian/src/internal/platform/mcp/server/tools_e2e_test.go`.

### Official protocol / SDK research

- The current MCP elicitation specification defines form and URL modes. Clients
  declare elicitation capability; users must be able to review/decline form
  requests. Form mode must not request secrets or credentials.
- The current official Go SDK documents server `ServerSession.Elicit`, client
  `ElicitationHandler`, and capability-dependent behavior. For protocol
  `2026-07-28`, elicitation uses the Multi Round-Trip Requests pattern; older
  protocol versions use the legacy server-initiated request path.
- Sources: [MCP elicitation specification](https://modelcontextprotocol.io/specification/2026-07-28/client/elicitation),
  [official Go SDK client feature guide](https://github.com/modelcontextprotocol/go-sdk/blob/main/docs/client.md),
  and [official Go SDK protocol guide](https://github.com/modelcontextprotocol/go-sdk/blob/main/docs/protocol.md).
- Current HAWP server is a handwritten stdio implementation advertising
  protocol `2024-11-05`; it does not currently implement elicitation. Do not
  retrofit interactive elicitation into the handwritten server without
  evaluating SDK adoption and client capability negotiation.

### Current client documentation review (2026-10-16)

- Cursor's current MCP protocol-support page lists elicitation as supported.
- GitHub documents MCP feature UX separately by host surface: Copilot Chat in
  VS Code supports server-provided prompts and resources; Copilot cloud agent
  and code review currently support MCP tools only, not resources or prompts.
- Claude Code documentation describes MCP resources and prompts, but the
  inspected source does not establish HAWP's elicitation experience in a
  connected Claude Code client.
- These are documentation observations, not successful HAWP-client calls.
  No end-to-end elicitation smoke was run. Keep manual questions as the
  supported fallback until the intended client and version are exercised.
- Sources: [Cursor MCP support](https://docs.cursor.com/context/model-context-protocol),
  [GitHub Copilot Chat MCP](https://docs.github.com/en/copilot/how-tos/copilot-in-your-ide/customize-copilot/extend-copilot-with-tools-and-context/extend-copilot-chat-with-mcp),
  [Copilot cloud-agent MCP boundary](https://docs.github.com/en/copilot/concepts/agents/cloud-agent/mcp-and-cloud-agent),
  and [Claude Code MCP](https://docs.anthropic.com/en/docs/claude-code/mcp).

### Remaining work

- [ ] Determine whether post-`v0.1.0` HAWP should add interactive elicitation
  or keep the current manual-question response as the supported path.
- [ ] Verify real support and user experience in the MCP clients HAWP intends to
  support; protocol/SDK capability does not prove provider-client support.
- [ ] If elicitation is adopted, select form-only versus other modes, define
  cancellation/decline/retry behavior, and retain the manual fallback.
- [ ] Add focused protocol/client integration tests for capability present,
  absent, decline, malformed response, and fallback states before implementation.
- [ ] Review whether current structured result/guide expectations need a small
  update; avoid duplicating resource or prompt work tracked in `2eea565c`,
  `3cf52223`, and `9dce1d95`.

**Decision boundary:** This item concerns server-requested user input during
intake. It is distinct from MCP prompt templates and resources, which are
tracked separately. The existing manual-question fallback remains the safe
behavior until client support and an implementation approach are verified.
