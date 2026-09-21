# MCP capability catalog for HAWP resources and prompts

**UUID:** `2eea565c-f56d-4139-b4eb-06eebfc2b456`
**Type:** improvement
**Reported:** 2026-09-13
**Status:** `plan-ready`
**Target Version:** `v0.0.25`

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

### Slice A - Capability inventory

- [ ] Inventory HAWP MCP surfaces and classify each as tool, resource, prompt, or future extension.
- [ ] Identify high-value read-only resources, starting with kit start-here,
  backlog status, active-plan summaries, and provider setup guidance.
- [ ] Identify prompt candidates, starting with intake shaping, status report,
  release checkpoint, and provider connection verification.

### Slice B - Thin resource/prompt design

- [ ] Propose stable URI names for read-only HAWP resources.
- [ ] Propose prompt names and arguments that avoid copying large kit files into
  every agent prompt.
- [ ] Keep tools for actions and mutations; do not turn resources or prompts into
  hidden write paths.

### Slice C - Implementation candidate

- [ ] Add the smallest read-only MCP resource list/read path after the design is reviewed.
- [ ] Add one prompt candidate only if the MCP SDK/server path supports it cleanly.
- [ ] Update `.hawp/kit/usage/mcp/README.md` and provider guides with the resulting capability map.

## Risk

Medium. This changes how clients discover HAWP, but the first slice is
read-only planning and should not alter existing tools.

## Verification

- Capability inventory cites exact source files.
- New resources, if implemented later, are read-only and scoped to the selected repository.
- MCP tool list remains backward compatible.
- `go test ./...`, `go vet ./...`, and `go run ./cmd/hawp work validate` pass from `librarian/src`.
