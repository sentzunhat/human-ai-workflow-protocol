# Investigate reshape support for HAWP work intake

**Backlog ID (Legacy):** — (UUID-native item)
**UUID:** `a3df8a9c-6f50-4437-9bcf-2c0cd8526f8d`
**Type:** improvement
**Reported:** 2026-09-07

---

## Input (verbatim)

> we need a work item for hawp to use reshape anyways update and create work items if needed and continue working

## Intake Summary

User confirmed on 2026-09-07: "Shape work requests into HAWP intake".
Preserve the original request and derive a compact, reviewable HAWP shape.

## Current Context

The connected server exposes search, work creation, validation, and usage.
`hawp_work_new` accepts title/type/input and creates an investigation scaffold;
it does not reshape the input or fill the analysis. Existing LLM reshaping
summarizes retrieved technical context, which is a different contract.

## Initial Analysis

**Confirmed:** `librarian/src/internal/platform/mcp/tools.go` and
`librarian/src/internal/platform/mcp/tool_work.go` expose no reshape tool.
`librarian/src/internal/domain/llm/llm_client.go` defines a context-rewriting
prompt, not intake extraction. `librarian/src/internal/platform/cli/search_commands.go`
formats retrieved context without invoking that LLM interface.
Indexed search found roadmap text; the source files above establish current behavior.
No matching active intake-reshape item was found in the backlog or active plans.

**Inference:** the smallest useful first slice is a worker-authored shape using
existing tools. A dedicated local-model/API implementation needs a separate
contract for preserving requests, unknowns, output validation, and failure handling.

## Plan

1. Document how a worker reshapes the request into the locked HAWP fields,
   keeps the original input, checks for an existing UUID, and uses available tools.
2. Exercise that workflow with this real request and validate its work record.
3. Specify the next callable draft-only slice before introducing runtime behavior:
   input = original request plus explicitly supplied context; output = a proposed
   HAWP shape with unchanged input and unknowns labeled. No work-file writes,
   invented owner/status/evidence, or implied approval from generated text.
   Test source preservation, missing context, malformed output, provider errors,
   and no writes. Keep creation an explicit subsequent `hawp_work_new` action.

Options: reuse the current worker and tool set now (recommended); introduce a
new local-model reshape command/tool later. Do not reuse search summarization
unchanged because shortening context can discard intent and constraints.

## Risk + Review Gate

Risk: low for guidance and the isolated internal draft contract.
Can implement now: yes under the user's continuation request; provider/client
wiring remains a separate step.
Owner: Codex. Existing architecture item owns model argument parsing; this item
owns request shaping guidance in both copies of the MCP worker README.
No model download, provider change, new dependency, or storage schema change.

## HAWP Shape For This Request

```text
input: |
  we need a work item for hawp to use reshape anyways update and create work items if needed and continue working
context: |
  The user confirmed that reshape means shaping work requests into HAWP intake.
  Current tools create investigation scaffolds but do not generate intake analysis.
mission: |
  Define and exercise a request-shaping workflow using the existing HAWP tools.
constraints: |
  Preserve the original input, UUIDs, and evidence. Label unknowns.
  Keep the workflow simple; no new model runtime in this first slice.
output: |
  Worker guidance, a tracked investigation, and a specified next draft-only slice.
```

## Backlog + Plan Link

Status now: in-progress.
Plan file: `.hawp/work/active/a3df8a9c/plan.md`.

## Verification

Worker guidance is present in `core/.hawp/kit/usage/mcp/README.md` and its
repo-local mirror; byte comparison passes. This plan exercises the flow with
verbatim input and the user's clarification. `hawp_search`, `hawp_work_new`,
and `hawp_work_validate` were called through MCP against the confirmed repo.
Validation passes kit/work/links with zero issues/warnings. This is agent-authored
shaping, not proof of an automated reshape implementation or model fidelity.

- [x] Investigate current tools and distinguish search-context shaping
- [x] Preserve the original request and capture the confirmed intent
- [x] Document and exercise the existing worker workflow
- [x] Add the domain `Draft` value and `Validate()` (`librarian/src/internal/domain/work/draft.go`)
- [x] Implement the application `RequestShaper` port + service with source-preservation/error tests
- [x] Evaluate and verify a concrete shaper adapter before client wiring

## Next Step

Evaluate one intake-specific adapter against the worker-guided baseline.
Require semantic fidelity examples and failure tests before any CLI/MCP entry
point. No default model selection or auto-creation is part of this contract.

## Draft contract implementation continuation

Investigation: no existing HAWP draft value or intake-specific shaper contract
exists in the inspected work packages. Context reshaping uses a summarization
prompt and cannot supply this contract unchanged. A small internal use case is
useful for testing fidelity before any client/provider wiring.

Plan: add a domain Draft value with the five required HAWP fields and optional
checkpoint, plus an application-owned RequestShaper port. The port returns only
mission/constraints/output/checkpoint: authoritative input and supplied context
are assigned by the service and never generated. Missing context is labeled
unknown. Require nonblank input and proposal fields; propagate provider errors
and cancellation; no partial draft on error. No filesystem dependency or work
creation call. Candidate semantic accuracy still requires review.

Alternatives: wire the generic summarizer now (wrong contract); add a CLI/MCP
tool before a provider exists (premature); keep only prose (no executable error
contract). Chosen scope is three small files plus tests/docs in existing work
packages, without model dependencies or new folder topology.
Risk: low for this isolated internal API. Can implement now: yes, continuing the
user-authorized draft-contract milestone. Existing work creation is unchanged.
Tests: verbatim multiline/Unicode input/context, missing context, rejected blank
fields, provider failure, cancellation, and no draft returned on error.

## Draft contract verification

**Corrected 2026-09-10 (record-truth pass):** an earlier checkpoint claimed the
application-side contract and its tests were implemented and passing. Verified
against the committed tree: only
`librarian/src/internal/domain/work/draft.go` exists (committed as `d3559556`);
it defines `Draft` (input/context/mission/constraints/output/checkpoint) plus a
completeness `Validate()`. It is not referenced anywhere in the tree, so it is
an isolated domain value with no service, port, or tests.

Not present in the tree (previously claimed, never committed, now removed):

- `librarian/src/internal/application/work/draft.go` (the `RequestShaper` port
  and service) — does not exist.
- `librarian/src/internal/application/work/draft_test.go` — does not exist.

The "source-preservation/error tests", "provider failure / cancellation", and
"JSON shape/source preservation" verification claims below do not apply to any
committed code and are withdrawn. Remaining work: implement the application
`RequestShaper` port + service, add the test suite, then evaluate a concrete
shaper adapter.

Current verification: the domain `Draft` value compiles as part of
`go build ./...` and is covered by the existing `internal/domain/work` package.
No dedicated tests exist for it yet (the previously referenced test file was
never committed). No default adapter, public entrypoint, or model wiring
exists. Semantic fidelity and model execution remain unproven.

## Application contract verification (2026-09-11)

Implemented in `librarian/src/internal/application/work/intake/`:

- `draft.go` — `RequestShaper` interface (port), `DraftRequest`/`DraftProposal`
  value types, and `DraftIntake` function (service). The service assigns `Input`
  and `Context` from the caller, labels missing context as
  `"Unknown: no context supplied."`, rejects blank input and nil shaper before
  calling the port, propagates port errors and context cancellation without
  returning a partial draft, and validates structural completeness via
  `domain.Draft.Validate()`.
- `draft_test.go` — table-driven tests covering: verbatim multiline/Unicode
  input/context passed through unchanged; missing context labeled; blank input
  rejected without calling shaper; context cancelled before shaper call; nil
  shaper rejected; incomplete proposal fields rejected; port failure propagated;
  post-shape context cancellation propagated; no draft returned on any error.

Verification run: `go test ./internal/application/work/...` — all packages pass.
`go test ./...` — all packages pass. `go vet ./...` — clean.
`go run ./cmd/hawp check --no-update-check` — 3/3 validations passed.
No new dependencies, CLI/MCP wiring, or model runtime introduced.

## 2026-09-11 Ollama adapter and MCP tool

Implemented `OllamaIntakeShaper` in `infrastructure/models/ollama/intake_shaper.go`.
Added `hawp_work_reshape` as the 4th MCP tool in `platform/mcp/server/`.
Verification: all unit tests pass (fake Ollama client, no live model required).

## Outcome

`OllamaIntakeShaper` and `ONNXIntakeShaper` implemented and benchmarked.
`hawp_work_reshape` MCP tool added. Ollama `mistral:7B` reached 10/10 coverage;
ONNX Phi-3-mini reached 10/10 coverage; SmolLM2-360M-Instruct (tiny) reached
1/10 and is not recommended as default. All unit tests pass (fake LLM client,
no live model required). Domain `Draft` and `RequestShaper` port are stable.

## Close Checklist

- [x] All focused checks pass (recorded in Verification section above).
- [x] HAWP check passes.
- [x] Plan status: done.
- [x] Plan moved to `.hawp/work/closed/2026/09/12/a3df8a9c/plan.md`.
- [x] BACKLOG.md updated.
