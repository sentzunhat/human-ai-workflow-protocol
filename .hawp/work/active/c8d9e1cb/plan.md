# Interactive MCP intake refinement with structured results

**UUID:** `c8d9e1cb-2189-45d4-bd4c-7cbc51710822`
**Type:** feature
**Reported:** 2026-09-13
**Status:** `plan-ready`
**Target Version:** `v0.0.24` subset, remainder `v0.0.25`

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
