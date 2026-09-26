# Apply MCP configuration preflight before public wrapper writes

## Outcome

The PR-review finding was implemented and reconciled with the current branch.

## Verification

The focused regression coverage and repository-wide Go, HAWP, distribution,
formatting, and diff-hygiene checks passed before close.

## Close Checklist

- [x] Outcome recorded.
- [x] Verification evidence recorded or referenced.
- [x] Backlog row removed from active coordination.
- [x] Plan archived under the close date.

**UUID:** `73bd501c-086e-4fc9-982b-01e43f89f42a`
**Type:** fix
**Reported:** 2026-09-23
**Source plan:** `37703223`

## Context

The public provider-config wrapper is used by `hawp init --provider`, but currently validates only provider names. The repository-aware prerequisite and selected-configuration preflight exists only in `Configure`.

## Fix Work

Move the common all-provider preflight behind the public wrapper: expand provider names, require a non-empty selection, reject redirected/non-regular prerequisites, inspect every requested provider configuration before writes, then call the existing private writer. Keep the private writer as the already-preflighted low-level operation and route `Configure` through the public wrapper.

## Verification

Prove direct-wrapper rejection for symlinked prerequisites and a later invalid provider configuration without provider or `.gitignore` writes. Prove regular prerequisite configurations still write successfully. Then run focused configuration and init tests, full Go tests, vet, build, HAWP validation, distribution validation, formatting, and diff-hygiene checks.

## Outcome

Done. `WriteProviderConfigs` expands provider names, validates all
prerequisites and selected destinations, and only then calls the low-level
writer; `Configure` delegates through that same boundary. Focused MCP/init
tests, full Go tests, vet, build, HAWP checks, and diff-hygiene verification
passed during the remediation sequence.
