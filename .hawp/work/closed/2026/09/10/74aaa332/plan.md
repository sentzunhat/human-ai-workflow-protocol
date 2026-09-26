# Strengthen HAWP slice harnesses and UUID-scoped artifacts

**UUID:** `74aaa332-c0f4-493c-a5e3-1e0dc1bde5bc`
**Type:** improvement
**Reported:** 2026-09-10

## Input

Make HAWP more useful on every compoundable implementation slice. Ensure
status, evidence, and supporting artifacts can be found by work-item UUID,
without colliding with the separate filesystem-boundary implementation.

## Investigation

Active and closed work items already use UUID folders. Status and evidence
currently use date-only flat folders, which makes filesystem lookup by work
item ambiguous. Existing guidance also says per-item evidence folders are not
required, so the desired UUID-scoped convention is not yet consistently
expressed.

## Mission

Define and later implement a small, repeatable HAWP slice harness that records
ownership, plan/checkpoint state, focused verification, full verification,
HAWP validation, attribution, and continuation instructions. Define a safe
UUID-scoped artifact convention without moving files owned by the active
filesystem-boundary agent.

## Constraints

- Do not edit or move the other agent's filesystem repository work.
- Preserve legacy flat artifacts and history during migration.
- Do not invent authors, model identities, co-author trailers, or evidence.
- Keep HAWP shaping guidance separate from runtime behavior.
- One harness slice first; no broad artifact migration in this item.

## Proposed Harness Contract

Every non-trivial slice records:

1. exact work UUID and owned paths;
2. intent, constraints, and non-goals;
3. focused checks before broad checks;
4. `go test`, vet/build, and HAWP validation results as applicable;
5. verified facts versus unproven claims;
6. attribution evidence and continuation instructions;
7. artifact paths under `status/YYYY/MM/DD/{uuid}/` or `evidence/YYYY/MM/DD/{uuid}/`.

## Acceptance

- A new slice can be resumed from its UUID without searching date-only files.
- HAWP validation accepts the new plan and checkpoint shape.
- The harness does not require a new field in the locked HAWP v0.1 schema.
- Existing legacy status/evidence files remain readable and are not rewritten.
- The filesystem-boundary agent can continue without merge conflicts.

## Expanded Scope (2026-09-10)

User requested that the harness cover all tools, AI models, and providers —
not just UUID-scoped artifacts. Expanded scope:

1. **Slice harness template** — universal pattern for any implementation slice:
   owned paths, intent/constraints, focused checks, full verification, HAWP
   validation, attribution, continuation instructions, artifact paths under
   `status/YYYY/MM/DD/{uuid}/` or `evidence/YYYY/MM/DD/{uuid}/`.

2. **Provider/model harness patterns** — standard test patterns for:
   - Embedding backends (Ollama, ONNX) — interface compliance, `Embed`/`EmbedBatch`/`Dimension`/`Close`
   - LLM clients (Ollama) — `Reshape`/`ReshapeBatch`/`Backend`/`Model`/`Close`
   - New providers added in future: what tests a new adapter must pass before wiring

3. **Tool harness patterns** — standard test patterns for:
   - CLI command handlers — typed parser rejection, handler golden-path, no-side-effect before validation
   - MCP tools — typed parser, JSON schema compliance, error propagation

4. **Standards documents** — written to `.hawp/kit/standards/`:
   - `slice-harness.md` — the universal slice harness template
   - `provider-harness.md` — what every provider adapter must verify
   - `tool-harness.md` — what every CLI/MCP tool parser must verify

5. **Guidelines doc** — written to `.hawp/kit/usage/`:
   - `harness-guide.md` — how to apply the harness on a new slice

## Status

`done`; all four standards/usage documents written and HAWP check passes.

## Completion checkpoint (2026-09-10)

Documents created:

- `.hawp/kit/standards/slice-harness.md` — universal slice harness template
- `.hawp/kit/standards/provider-harness.md` — embedding and LLM backend verifications
- `.hawp/kit/standards/tool-harness.md` — CLI and MCP tool parser contract
- `.hawp/kit/usage/harness-guide.md` — how to apply the harness on a new slice

HAWP check: passed (`go run ./cmd/hawp check --no-update-check` from `librarian/src`).
No Go source modified; `go test ./...` not required.

## Attribution

The architecture direction is user-directed. Implementation and verification
are Codex/model-assisted. Add co-author trailers only when real contributor
identity details are supplied.

## Outcome

Four slice-harness documents are in place and HAWP check passes. `slice-harness.md`,
`provider-harness.md`, and `tool-harness.md` define the universal per-slice record
(owned paths, focused-then-full checks, HAWP validation, attribution, continuation)
and the standard verification patterns for provider adapters and CLI/MCP tool
parsers; `harness-guide.md` explains how to apply the harness on a new slice. No Go
source, schema field, or legacy artifact was modified.

## Verification

`go run ./cmd/hawp check --no-update-check` passes (kit, work, and links).
All four documents exist under `.hawp/kit/standards/` and `.hawp/kit/usage/`.
No Go source was modified and no new locked-schema field was added, so no
`go test ./...` run was required; legacy flat status and evidence files remain
readable and unrewritten.

## Close Checklist

- [x] Four standards/usage documents written and linked from the harness guide.
- [x] HAWP check passes with the new plan and checkpoint shape.
- [x] Filesystem-boundary agent work left untouched; no merge conflicts.
- [x] Plan moved to `closed/2026/09/10/74aaa332/`.
