---
uuid: bae2c316
title: "Optional work-item tracking providers (Vikunja first)"
type: status
date: 2026-10-22
event-date: 2026-10-20
---

# 2026-10-20 — Optional Work-Item Tracking Providers

**Event date:** 2026-10-20
**Archive date:** 2026-10-22
**State:** Architecture direction agreed in discussion; design and implementation not yet approved or completed.

## Intent

Preserve the evolution of an optional visual work-tracking integration for HAWP so that a future contributor can continue without conflating a third-party task board with HAWP's repository-local work records. This is a public-safe architectural checkpoint, not a shipped-feature announcement or executable work item.

## Before

- HAWP's canonical shape carries human/agent intent and constraints; it is not a task-board product or workflow runtime.
- Project-local `.hawp/work/` holds the backlog, plans, evidence, decisions, and handoffs. These records remain authoritative.
- An initial integration idea used a self-hosted Vikunja Docker service as a human-facing Kanban board and a small Go adapter to map HAWP work-item IDs to Vikunja tasks using its REST API.
- The initial proposal treated Vikunja as an optional visualization service, not as the source of truth. No deployment or adapter implementation was reported as completed.

## During — 2026-10-20

- The integration direction broadened from a Vikunja-specific adapter to **optional, provider-agnostic work-item tracking adapters**.
- The same HAWP project should be able to opt into an external tracker; different projects may select different providers independently.
- Vikunja remains the proposed first concrete integration. GitHub Issues and other compatible trackers are possible future adapters, not commitments or implemented integrations.
- The phrase **work-item provider** here means an external tracking destination. It is distinct from HAWP's existing coding-agent/provider guidance and MCP client configurations.

## After — agreed constraints

1. **Canonical ownership:** HAWP work identities and repository-local `.hawp/work/` artifacts remain authoritative for intent, plans, decisions, evidence, verification, and handoffs.
2. **Optionality:** HAWP must work normally with no tracking provider configured, including offline use and provider outages.
3. **Per-project choice:** Provider selection and credentials belong to opt-in, project-scoped integration configuration, not to the canonical HAWP protocol fields.
4. **Stable identity:** Map a stable HAWP work-item identifier to a provider-native task identifier without replacing or regenerating the HAWP ID when a provider changes.
5. **Safe synchronization:** Start by exploring idempotent HAWP-to-provider publishing. Changes coming from a provider must not silently rewrite authoritative repository files; conflict and approval policy are required before any write-back.
6. **Small boundary:** Keep provider-specific clients and schema translation behind narrow integration interfaces; do not couple HAWP's protocol or local work-domain rules to Vikunja.
7. **Provider independence:** A visual Kanban board is a coordination view of HAWP work, not a new mandatory source of truth.

## Strategic impact

The decision changes the design target from **HAWP + one task board** to **HAWP + optional interchangeable external work tracking**. This improves portability across projects and avoids locking work history to a hosted tracker. The smallest useful validation remains one self-hosted Vikunja instance and one adapter, with no requirement to build a general plugin marketplace in the first slice.

## What was inspected and directly verified during archive

- `README.md` describes HAWP as a repo-local protocol with optional tooling and `.hawp/work/`.
- `AGENTS.md`, `.hawp/kit/usage/status-report.md`, and `.hawp/work/status/README.md` establish the repository's date/ID-based checkpoint convention.
- `.hawp/work/BACKLOG.md` defines work states including `inbox`, `analyzing`, `plan-ready`, `approved`, `in-progress`, `blocked`, and `done`.
- `librarian/src/` contains the Go CLI/MCP structure and existing work-domain separation.
- A targeted search of the default branch found no Vikunja adapter or dedicated work-item tracking provider implementation within the inspected scope.
- This archive adds documentation only; it does not deploy Vikunja, create a task-board sync, alter canonical schema, or change application code.

## Unresolved design questions

- Exact provider port/interface and capabilities; where mapping state is stored and how it survives clone, migration, and provider replacement.
- Read/write permissions, credential storage, tenant/project isolation, and handling of private work descriptions.
- Mapping between HAWP lifecycle states and a provider's columns, labels, or status fields.
- Idempotency and duplicate prevention; conflict detection, reconciliation, rate limits, offline recovery, and failure visibility.
- Whether reverse synchronization is appropriate at all and, if so, what explicit validation/approval gates authorize mutations to `.hawp/work/`.
- Which CLI/MCP commands expose the integration, if any, without expanding the canonical HAWP shape.
- Resources and maintenance impact of the self-hosted service; no performance or availability testing has been performed.

## Continuation state

**Current state:** A provider-agnostic optional integration is an agreed design direction only. HAWP remains usable without it. No implementation milestone has been verified.

**Next logical milestone:** A small, reviewable proposal for the provider interface and first Vikunja adapter, using existing HAWP work/plan conventions.

**Resume sequence:**

1. Inspect current `.hawp/work/BACKLOG.md` and relevant active plans; update an existing matching work item if present rather than duplicating it.
2. Investigate current Go work-domain/application/adapter seams; propose the smallest optional interface and explicit per-project configuration.
3. Define a one-way idempotent HAWP-to-Vikunja task projection with stable ID mapping, no mandatory credentials, and safe no-provider behavior.
4. Specify binary acceptance checks for an offline/unconfigured HAWP session, repeat sync without duplicates, failed provider access, distinct project configurations, and unchanged authoritative HAWP files.
5. Present executable work items and an implementation plan for explicit approval before coding. Preserve the existing REVIEW → VERIFY → PRIORITIZE → CREATE WORK ITEMS → PRESENT PLAN → APPROVAL gate.

**Blockers:** Interface and mapping ownership are undecided; no verified Docker deployment, credentials, or end-to-end tests exist. These are design/deployment dependencies, not failures in current HAWP functionality.

**Cross-project effect:** Multiple repositories may independently opt into the same adapter contract. This checkpoint is canonical to HAWP; downstream projects should reference it rather than copy its text.

**Resume from:** Provider-contract design and existing-work-item discovery; no code changes yet.
**Next objective:** Approved, minimal Vikunja-first adapter plan with acceptance criteria.

## Privacy and scope

This public artifact intentionally omits personal details, internal company context, private infrastructure topology, credentials, and unrelated product names. It does not introduce or change a formal HAWP protocol field.
