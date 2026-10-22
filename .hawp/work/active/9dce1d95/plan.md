# Implement reusable MCP prompts for HAWP workflows

**UUID:** `9dce1d95-c017-4741-859f-ec7e9763df86`
**Type:** improvement
**Reported:** 2026-10-16
**Status:** `in-progress`
**Target Version:** `v0.1.0` follow-up lane; implementation explicitly resumed by maintainer on 2026-10-16

---

## Original Input (verbatim)

> After the v0.1.0 release, investigate and plan reusable MCP prompt templates for HAWP workflows. Review intake shaping, work review, status report, release checkpoint, and MCP connection check candidates; define prompt names/arguments, keep dynamic repository data separate, preserve user approval points, and verify client compatibility and tests. Prompts must not create, update, delete, configure, commit, publish, or release as implicit side effects. Do not implement prompts before separate review and implementation approval.

## Resume Decision (2026-10-16)

The maintainer explicitly authorized continuing from the current HAWP capability
catalog and implementing curated resources/prompts, followed by later provider
and distribution propagation. This supersedes the prior parked/no-implementation
gate for this item. Implement only prompts that add a useful client-selectable
workflow entry point and preserve explicit user decision points.

## Investigation

### Directly verified

- Before the current edits, the source stdio MCP server advertised only `tools` and dispatched `tools/list` and `tools/call`.
- The implementation now advertises prompts and handles `prompts/list` and `prompts/get` for the single `hawp_intake_request` template.
- The prompt accepts a required request and optional context, treats both as JSON-encoded user data, and makes no tool calls or writes.
- Focused Go tests now cover list metadata, required/optional argument declarations, safe delimiter handling, the no-write boundary, and invalid/unknown arguments (`server/prompts_test.go`).
- HAWP already documents multi-step MCP workflows in `.hawp/kit/examples/mcp-intake-to-work-doc.md` and provider instructions under `.hawp/kit/usage/mcp/`.
- The active catalog `2eea565c` proposes intake, work review, status report, release checkpoint, and connection-check prompts.

### Inference / open questions

- Client-side prompt discovery and argument rendering remain unverified in HAWP's supported clients.
- Static prompt text must not substitute for current repository state. Dynamic context should be read through separately scoped resources or tools.

### Scope and risk

- Scope: implement and verify a small prompt capability over local stdio. No prompt may create, update, delete, configure, commit, publish, or release as an implicit side effect.
- Risk: low to medium; protocol capability is small, but confusing or overbroad prompts could encourage unintended writes.
- Review gate: maintainer authorized implementation on 2026-10-16; keep the initial prompt set small and explicit.

## Plan

1. Keep prompt list/get response schemas compatible with the retained `2024-11-05` protocol.
2. Start with `hawp_intake_request`, since it provides a distinct selectable entry point for the common intake workflow.
3. Preserve required/optional argument metadata, safe interpolation, manual clarification, and explicit create/approval boundaries.
4. [x] Add focused tests for list/get, missing/unknown arguments, and untrusted argument interpolation.
5. Verify prompt selection and argument experience in intended clients before documenting provider support.

Initial implementation: one intake prompt. Other candidates remain deferred until the first prompt's client discovery value is verified.

## Backlog / Work State

- Status: `in-progress` as of 2026-10-16.
- Plan path: `active/9dce1d95/plan.md`.

## Verification

- `go test ./internal/platform/mcp/server` passes, including the new prompt protocol tests.
- `./.hawp/bin/hawp kit validate`, `./.hawp/bin/hawp work validate`, and `git diff --check` pass; work validation reports 0 issues and 1 existing verification-clarity warning.
- Prompt discovery and argument UI remain unverified in intended clients.
