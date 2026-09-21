# Review groups — pre-reshape checkpoint

Captured before the draft-service continuation on 2026-09-07. No staging or
commits performed. Includes untracked tests and preserved work history; deleted
plan paths belong with their archived destinations, not standalone deletion.
Shared README/registry/kit changes should be reviewed with their relevant code
intent rather than assumed independent. No installed binary refresh is included.

## Index models

- `librarian/src/internal/application/index/ingest-service.go`
- `librarian/src/internal/domain/index/chunk.go`
- `librarian/src/internal/infrastructure/sqlite/domain_models_test.go`
- `librarian/src/internal/infrastructure/sqlite/index.go`

## CLI argument safety

- `librarian/src/internal/platform/cli/embed_args.go`
- `librarian/src/internal/platform/cli/embed_args_test.go`
- `librarian/src/internal/platform/cli/kit_normalize_args.go`
- `librarian/src/internal/platform/cli/kit_validate_args.go`
- `librarian/src/internal/platform/cli/model_pull_args.go`
- `librarian/src/internal/platform/cli/model_pull_args_test.go`
- `librarian/src/internal/platform/cli/mutation_boundary_args_test.go`
- `librarian/src/internal/platform/cli/registry.go`
- `librarian/src/internal/platform/cli/work_new_args.go`

## Migration and packaging

- `.hawp/work/active/0cb0f9b0/plan.md`
- `.hawp/work/closed/2026/09/07/0cb0f9b0/plan.md`
- `.hawp/work/closed/2026/09/07/0cb0f9b0/progress-archive.md`
- `core/.hawp/bin/hawp`
- `librarian/src/CHANGELOG.md`
- `librarian/src/internal/domain/distribution/binary_download_test.go`

## Guidance and work records

- `.claude/rules/architecture-state.md`
- `.github/copilot-instructions.md`
- `.hawp/kit/usage/mcp/README.md`
- `.hawp/work/BACKLOG.md`
- `.hawp/work/STATUS.md`
- `.hawp/work/active/47c793d6/architecture-state-prior.md`
- `.hawp/work/active/47c793d6/audit.md`
- `.hawp/work/active/47c793d6/launcher-audit.md`
- `.hawp/work/active/47c793d6/plan.md`
- `.hawp/work/active/a3df8a9c/plan.md`
- `.hawp/work/active/c42e04b7/plan.md`
- `.hawp/work/active/e5fca9c7/plan.md`
- `.hawp/work/closed/2026/09/06/c42e04b7/intake-original-archive.md`
- `.hawp/work/closed/2026/09/06/c42e04b7/plan.md`
- `.hawp/work/status/2026/09/06/0cb0f9b0-status.md`
- `.hawp/work/status/2026/09/06/47c793d6-status.md`
- `.hawp/work/status/2026/09/06/checkpoint-2026-09-06-architecture-audit.md`
- `.hawp/work/status/2026/09/07/0cb0f9b0-status.md`
- `.hawp/work/status/2026/09/07/47c793d6-status.md`
- `.hawp/work/status/2026/09/07/e5fca9c7-status.md`
- `core/.hawp/kit/usage/mcp/README.md`
- `librarian/src/README.md`

## Verification boundary

Latest full tests/vet and distribution checks passed; fresh MCP validation
passes. Code and test groups are covered together; this manifest is a review
map, not proof that each group builds independently as a separate commit.
