# Design and prioritize improved HAWP MCP tools

**UUID:** `1673de65-63e0-4a8e-bade-cbdde4095912`
**Type:** improvement
**Reported:** 2026-10-17
**Status:** `closed`
**Target Version:** TBD after design review
**Evidence:** [MCP capability catalog](../../closed/2026/10/17/2eea565c/plan.md); [tool-first audit](../../closed/2026/10/16/c147f2ae/plan.md)

---

## Input (verbatim)

> Revisit the audit and the active plan to identify concrete improvements—what to add or change, why it helps, and its inputs, outputs, and safety boundaries—then create a new work item and close this one it will help as evidence and we dont blow up the scope of this

## Intake Summary

Design a small, prioritized set of concrete MCP tool improvements from the source-backed catalog and tool-first usage audit. Keep the catalog as closed evidence, separate tool design from implementation, and avoid expanding the unrelated resource, prompt, usage-metering, provider propagation, and remote transport work.

## Current Context

The local stdio server currently exposes seven tools, three fixed read-only resources, and one prompt. The completed catalog records inputs, outputs, behavior, filesystem scope, and workflows. The tool-first audit found that client allowlists and session exposure affect whether tools are usable, while the server has no existing-record update operation. Existing work creation scaffolds an inbox plan but does not complete investigation/plan content; agents then edit files directly. Search returns ranked chunks/context, while intake composes search and local model shaping.

## Initial Analysis

**Directly verified:**

- `hawp_work_new` creates a scaffold and backlog inbox row; it does not fill investigation or plan details.
- No MCP update or delete operation exists. Existing plan/backlog updates are client/agent file edits.
- `hawp_work_intake` combines indexed retrieval and local-model shaping without creating work. `hawp_search` and `hawp_work_reshape` remain independently useful primitives.
- `hawp_work_doc` creates a canonical secondary-document path; content is filled separately.
- Existing tool allowlists may filter tools at the client surface; a server catalog alone does not establish client availability.

**Inferred (to be assessed in this work):**

- The biggest workflow gap may be a typed, bounded way to update existing work records, but safe scope, conflict detection, preview/confirmation, and recovery need design before implementation.
- A guided create-from-intake operation might reduce scaffold follow-up, but risks hiding investigation and user approval. Prefer improving structured outputs or a narrowly typed create payload unless evidence shows a separate operation is useful.
- Tool descriptions and input schemas may need to make read/write effects, prerequisites, and workflow fit clearer; preserve stable names unless a compatibility case supports change.

## Architecture Direction (maintainer guidance, 2026-10-17)

Keep the design **simple to use and scalable to extend**. Use Zacatl as an
architecture reference for deliberate boundaries and composition, not as a
template to copy. HAWP must express those ideas in its own domain, Go codebase,
existing contracts, and lightweight workflow purpose.

Apply HAWP's existing standards:

- Define capability contracts before implementation details; state inputs,
  outputs, errors, ownership, and effects.
- Keep MCP handlers at the transport edge: parse/validate arguments and shape
  responses. Put work-record decisions and validation behavior in application
  or domain services; keep filesystem and protocol details in adapters.
- Compose dependencies explicitly so CLI and MCP can reuse the same behavior
  without duplicating business rules or pulling transport concerns inward.
- Prefer the smallest cohesive operation set. Add a tool only when it expresses
  a distinct user intent; retain primitives when workflows compose them.
- Avoid generic CRUD, arbitrary path access, wrapper layers, speculative
  extension points, and new abstractions without demonstrated reuse.
- Keep docs evidence-linked: distinguish what exists now from proposals and
  link tool contracts to source/tests. Do not copy Zacatl-specific structure,
  terminology, dependency injection choices, or framework conventions where
  HAWP's Go architecture differs.

Relevant HAWP guidance: `.hawp/kit/standards/service-design/` (service
boundaries, layered composition, handler responsibilities, dependency
composition, contract testing, evidence-linked docs),
`.hawp/kit/standards/guidelines/architecture.md`,
`.hawp/kit/standards/guidelines/documentation.md`,
`.hawp/kit/standards/guidelines/testing.md`, and
`.hawp/kit/standards/patterns/evidence-discipline.md`. The Go structure is
documented in `librarian/src/README.md`; follow its current application/domain/
infrastructure/platform boundaries rather than imposing another repository's
folder map.

### Design test

For each proposed capability, answer: (1) which user workflow is materially
better, (2) why an existing tool or composition cannot serve it, (3) which
layer owns the behavior, (4) whether CLI and MCP can reuse it, and (5) what is
the smallest contract that can prove value. If those answers require a generic
framework or duplicate an active item, defer the proposal.

## Risk + Review Gate

**Risk:** medium
**Gate:** review the proposed capability contract and compatibility boundaries before implementation planning.

## Mission

Produce a compact, evidence-backed MCP tool improvement proposal with prioritized add/change/keep decisions. For each proposal define the user problem, intended workflow, name or change, inputs, output contract, effects, repository/home scope, prerequisites, failure states, confirmation/authorization boundary, compatibility impact, and acceptance evidence.

## Scope

**In scope:** existing tool contract review; concrete improvements to create, typed update, search/intake, document creation, and validation; client discoverability implications; prioritization and release-sized slices.

**Out of scope:** implementation, resource/prompt expansion, usage metering/retention, provider/distribution propagation, remote/HTTP MCP, generic filesystem writes, delete operations, and changing the HAWP core shape.

## Candidate Design Questions

1. **Typed update of an existing work item (candidate add):** Which canonical fields and statuses can change? Accept a work-item ID plus explicit field patch and expected current value/revision; return changed fields, plan/backlog paths, and validation summary. Define preview/confirmation, stale-write conflict, atomicity/recovery, allowed active/parked/closed transitions, and audit trail. Never accept arbitrary paths or file contents. Decide whether this should be deferred if the confirmation/restore contract cannot stay small.
2. **Create from a completed intake (candidate change to `hawp_work_new`):** Should the tool accept structured investigation/plan fields rather than only title/type/input? Return created ID and paths plus validation result. Keep creation explicit, reject incomplete/unsupported fields, avoid implying the model's draft is verified, and preserve a user decision before writing. Compare this against a separate create-from-intake operation and the current scaffold flow.
3. **Search/intake ergonomics (candidate change):** Consider clearer intent-specific schema/outputs or a single recommended discovery contract. Keep standalone search and reshape if they remain composable. Inputs/outputs must distinguish index retrieval from local-model shaping, expose missing index/model and no-match states, bound result/token controls, and avoid claims that estimated usage is provider token usage.
4. **Work document lifecycle (candidate change):** Clarify create-only/idempotent behavior and address how a caller safely fills document content without introducing a generic path writer. If a typed append/fill operation is proposed, specify allowed document types, ownership, overwrite/conflict behavior, size bounds, and recovery.
5. **Validation output (candidate change):** Assess whether structured machine-readable findings with stable codes, paths, and severities would improve agent follow-through while retaining human-readable summary and read-only behavior.
6. **Tool descriptions and client discovery (candidate change):** Identify schema/description improvements that aid MCP clients in choosing the right capability and recognizing effects. Preserve current names and allowlists by default; document any compatibility consequence for a proposed rename or removal.

These are hypotheses to evaluate, not approved implementation decisions. Select no more than two high-value proposals for the first implementation slice; defer the rest with rationale.

## Source Review and Overlap Resolution (2026-10-17)

Reviewed `server/tools.go`, `tool_work.go`, `tool_search.go`, `types.go`,
`tool_usage.go`, `server.go`, `application/work/intake/intake.go`,
`application/work/validation/validate.go`, `application/work/doc/doc.go`, and
the active plans named above.

### HAWP architecture fit (source-backed)

- `librarian/src/README.md` documents application services and infrastructure
  adapters with CLI/MCP composition at the platform edge; it also records
  ongoing provider/composition separation rather than claiming a finished
  architecture.
- Work validation already returns a reusable domain `Report` through
  `application/work/validation.Validate`. The CLI renders it in
  `application/work/validation.Render`; the MCP handler currently calls the
  composite `application/check.Run` and captures its human-readable CLI text.
- Therefore the smallest HAWP-specific structured-output seam is response
  mapping at the MCP platform adapter, based on the existing work validation
  report. Keep domain validation rules transport-agnostic. Retain the existing
  composite CLI text contract; do not make CLI depend on MCP response types.
- MCP `rpcResponse` currently has no explicit `structuredContent` field and
  the handler response helper serializes JSON into text. Protocol compatibility
  and intended-client exposure of any structured result must be checked before
  selecting the wire representation. Keep text as the guaranteed fallback.
- HAWP standard references used: service boundaries, layered composition,
  handler responsibilities, dependency composition, contract testing,
  evidence-linked documentation, architecture, documentation, testing, and
  evidence discipline. Their reusable principles are contract-first design,
  edge orchestration, application/domain ownership of business rules,
  deterministic composition, observable contract tests, and evidence-linked
  docs. This plan does not import Zacatl's language/framework choices or
  structure.

### Candidate decisions

| Priority / decision | User problem and workflow | Proposed contract (design target) | Effects, boundary, compatibility |
| --- | --- | --- | --- |
| **Later — Design typed `hawp_work_update` only if direct-edit pain justifies a new write capability** | A caller can create a work scaffold, but updating an existing plan/backlog requires direct filesystem edits. The need is real, but safety and lifecycle breadth make it larger than the first slice. | If pursued: accept a work item ID and a typed patch against a deliberately small field set; return canonical ID, changed fields, paths, and validation result. A separate design must settle stale-write detection, preview/confirmation, valid status transitions, atomicity/recovery, and provenance before implementation. | Repository-scoped canonical record mutation only; no arbitrary paths, document bodies, or delete. Keep deferred until the contract is small and reusable from application services and CLI/MCP. |
| **P1 — Improve the existing `hawp_work_validate` response, without adding a tool** | Agents currently receive a human-oriented composite CLI report and must infer findings. Structured findings can support follow-through without another operation. | No new arguments. Preserve existing text and add structured overall status/checks/findings, each with stable code, severity, repo-relative path, concise message, and optional next step. Map from application validation results at the MCP boundary. | Read-only selected-repository validation; no mutation. Preserve current PASS/WARN/FAIL and CLI behavior. Deterministic output; no home paths/secrets. Keep text fallback and do not move validation rules into the handler. |
| **P1 — Clarify descriptions and response contracts for the existing tools** | Similar names and compound behavior can make it hard to choose search, reshape, intake, create, and validate; JSON is currently returned as text for several tools. | Keep current names and inputs. Clarify retrieval vs supplied-context shaping vs composed intake; distinguish estimates from provider token usage; state prerequisites, result states, limits, and read/write scope. Add structured fields only where supported without replacing text. | No new side effects or filesystem scope. Preserve input and allowlist compatibility. Keep the changes at the MCP boundary and reuse existing application behavior. Verify schema visibility in clients before claiming improved usability. |

### Deferred candidates and rationale

- **Create directly from completed intake:** defer. `hawp_work_new` currently
  deliberately scaffolds an investigation record; accepting model-shaped fields
  as if verified could skip investigation or imply approval. Revisit only with
  a clear user-confirmed handoff and no duplication of `c8d9e1cb`.
- **Document fill/append tool:** defer. `hawp_work_doc` intentionally creates
  canonical paths only. A fill operation needs document ownership, overwrite
  and conflict rules, bounded content and recovery; it is not required to solve
  the existing-record plan/backlog update gap.
- **Search/intake redesign:** keep current tools and their distinct primitives.
  No evidence justifies another wrapper; structured-state/manual-question work
  is already active under `c8d9e1cb`.
- **Usage controls/measurement:** defer to `e1dc28ac`; do not duplicate usage
  logging, retention, or token-accounting work here.
- **Resources, prompts, provider distribution, remote transport:** remain in
  `3cf52223`, `9dce1d95`, `3c137bad`, and `429e075e` respectively. This item
  does not reopen or expand those lanes.

### Recommended bounded first implementation slice

Proceed with one cohesive, read-only MCP contract improvement: structured
validation results plus a concise pass over existing tool descriptions. This
does not add a capability or change names/inputs. Keep typed update as deferred
design work until preview, confirmation, conflict, transition,
atomicity/recovery, and provenance are resolved. This keeps the first slice
small and leaves write-capable design for separate review.

Acceptance evidence for the first slice should include unchanged legacy text
and tool names, structured findings with stable codes and repository-relative
paths, deterministic ordering, no sensitive absolute paths, application tests
for the report mapping plus MCP boundary tests for pass/warn/fail output, and at
least one intended-client call confirming the structured result is available
or an explicit recorded client limitation. Documentation should state the
current implementation and fallback separately from the proposed contract.

## Plan

- [x] Read the closed catalog and `c147f2ae` audit as source evidence; verify relevant current handlers/schemas before finalizing recommendations.
- [x] Compare candidate improvements against observed workflows and failure points; identify duplication with active work (`c8d9e1cb`, `3cf52223`, `9dce1d95`, `3c137bad`, `e1dc28ac`, `b152c56c`).
- [x] Define prioritized tool contracts with names, inputs, outputs, effects, scope, prerequisites, errors, safety/confirmation, recovery, and compatibility impact.
- [x] Recommend a first bounded implementation slice (maximum two proposals) and explicit deferred list.
- [x] Identify acceptance checks and client verification needed; keep client availability evidence distinct from server tests.
- [x] Review the recommended structured validation/description slice and create a separate implementation plan for the protocol-version-gated work.

## Consolidation and Final Recommendation

### Combine these workflows, not their public tool identities

- `hawp_search` + `hawp_work_reshape` are composed by `hawp_work_intake` for
  repository-grounded intake. Keep all three: search serves retrieval alone;
  reshape serves caller-supplied context; intake is the explicit convenience
  composition. Removing the primitives would reduce reuse, while adding more
  wrapper variants would duplicate the same sequence.
- Keep `hawp_work_new` and `hawp_work_doc` separate. They create different
  lifecycle records at different canonical paths, and a status/evidence doc
  may be needed without creating a new work item.
- Keep `hawp_usage` separate. It reads user-home telemetry, has opt-in/privacy
  semantics, and is not part of repository work mutation or validation.
- Keep `hawp_work_validate` as the single MCP verification entry point. Do not
  add another `hawp_check` tool; improve its current composite result contract.
- Keep resources and prompts as separate MCP capability classes. They are not
  candidates for folding into the tool list, and their implementation/client
  proof remains owned by active work.

### What needs to be done

1. **Read-only validation response:** implementation design is in
   `9c9f9bed`. Its next gate is an explicit protocol revision decision because
   the currently advertised `2024-11-05` revision does not define
   `structuredContent`/`outputSchema`. Do not implement a non-standard field.
2. **Typed update:** defer; do not create this capability now. The current
   evidence shows direct file edits but not enough repeatable need to justify
   introducing write authorization, conflict handling, state transitions, and
   recovery machinery. Continue direct plan editing under HAWP guidance. Reopen
   only if repeated user workflows show direct edits are a material blocker and
   a small safe contract is possible.
3. **Provider prompts/distribution:** follow the maintainer's next-release
   direction in existing `3c137bad`; keep prompt/resource provider guidance
   and generated distribution propagation together there. Do not duplicate it
   or claim client support without live proof.
4. **Other active owners:** continue resource, prompt, elicitation, provider
   configuration, and usage work under their existing UUIDs.

### Explicitly do not combine

Do not merge create and secondary-document creation; do not make intake create
work implicitly; do not fold usage into search or validation; and do not create
a single catch-all work tool. Compose existing primitives in documented
workflows, while keeping public capabilities aligned to distinct user intent
and effect.

## Verification

- Each recommendation maps to direct source/audit evidence or is clearly labeled an inference.
- Add/change/keep decisions state why and identify relevant user workflow.
- Every proposed capability defines input/output and read/write/delete effects, repository/home scope, prerequisites, errors, and decision boundary.
- No proposal duplicates an active item or silently introduces arbitrary path access, implicit writes, deletion, or remote access.
- Compatibility impact and a staged acceptance strategy are explicit.

## Next Step

## Outcome

The investigation and consolidation design are complete. The validation
response contract is tracked in `9c9f9bed` and waits on the protocol-version
decision. Typed update is deferred as not justified by current evidence.
Provider prompts/distribution follow-up stays in `3c137bad` for the next
release candidate.

**Close date:** 2026-10-17

## Close Checklist

- [x] Compared proposed improvements against current source and active work.
- [x] Recorded which tools compose and which remain separate.
- [x] Deferred typed update based on its current safety/complexity cost and
  insufficient evidence of repeated user need.
- [x] Created the protocol-gated validation response plan `9c9f9bed`.
- [x] Kept provider prompt/distribution work in `3c137bad` and aligned it with
  the maintainer's next-release direction.
