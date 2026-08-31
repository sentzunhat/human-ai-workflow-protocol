# v0.0.23 root tests layout and port adapter segregation

**Backlog ID (Legacy):** — (UUID-native item)
**UUID:** `0bd32051-fd2e-4bd5-a92f-65ef7c20ef29`
**Type:** refactor
**Reported:** 2026-08-29

---

## Input (verbatim)

> And also the test files should be in a replicated tests folder on the root beside the src folder please what about adding as a work item for this next compoundable work

## Intake Summary

Track the next architecture slice after the current shaping cleanup: move tests
to a root-level replicated `tests/` tree beside `src/`, and use that move to
clean up any adjacent port/adapter or ownership drift exposed by the layout
change.

## Current Context

- The current implementation lane `387a37b7` is still focused on shaping and
  shared search-service cleanup
- The follow-up audit lane `5957aaf4` exists to queue next simplification work
  without colliding with in-progress implementation
- `librarian/src/go.mod` makes `librarian/src/` the Go module root today
- Current tests are co-located under `librarian/src/internal/...`
- As of 2026-08-30, exported black-box tests can be moved under
  `librarian/src/tests/...`, but a sibling `librarian/tests/` tree beside `src/`
  would fall outside the module root and break normal `internal/...` import use

## Initial Analysis

**Directly verified:**

- The user explicitly wants test files moved into a replicated root-level
  `tests/` folder beside `src/`
- That change touches folder topology and verification layout, so it should be
  handled as its own work item rather than piggybacked onto the current
  formatter pass
- `librarian/src/` is the Go module root, so `librarian/src/tests/...` is the
  safe replicated target for black-box tests that still need `internal/...`
  imports
- On 2026-08-30, the first replicated tests slice was implemented:
  - moved exported-behavior tests from:
    - `internal/application/search/service_test.go`
    - `internal/application/work/intake_test.go`
    - `internal/application/work/normalize_test.go`
    - `internal/platform/cli/run_test.go`
  - into:
    - `tests/application/search/service_test.go`
    - `tests/application/work/intake_test.go`
    - `tests/application/work/normalize_test.go`
    - `tests/platform/cli/run_test.go`
- Verification passed on 2026-08-30:
  - `go test ./tests/... ./internal/application/search ./internal/application/work ./internal/platform/cli ./internal/domain/work`
  - `go test ./...`

**Inferred (not yet proven):**

- A mirrored `tests/` tree could make ownership boundaries clearer during the
  next port/adapter cleanup, especially if application, domain, and
  infrastructure tests are reviewed together
- The highest-risk part is not the file move itself but any import/package
  fallout from Go's current module/package structure

**Likely scope:**

- Audit the current Go test/package layout and constraints
- Decide the safest target structure for a root-level replicated `tests/` tree
- Identify the minimum port/adapter cleanup that naturally belongs with the move
- Implement only after the structure and verification strategy are explicit

## Risk + Review Gate

**Risk:** medium
**Gate:** implement incrementally; keep white-box tests co-located; do not merge

## Backlog + Plan Link

**Status now:** done
**Plan file:** work/closed/2026/08/31/0bd32051/plan.md

## Next Step

- [x] Work item created from the requested next compoundable slice
- [x] Inspect current Go package/test constraints for a replicated `tests/` tree
- [x] Propose the smallest safe migration plan and boundary cleanup set
- [x] Implement the first safe replicated `src/tests` slice for exported behavior
- [x] Continue moving exported black-box tests by capability where package internals are not required
- [x] Keep package-internal white-box tests beside code until a safer seam exists

## 2026-08-30 Migration Checkpoint

**Directly verified:**

- The replicated Go test tree now covers `14` black-box test files under
  `librarian/src/tests/...`
- The second migration pass added exported-behavior coverage for:
  - `tests/application/db/init_service_test.go`
  - `tests/application/links/check_test.go`
  - `tests/application/links/clean_test.go`
  - `tests/application/provision/provision_test.go`
  - `tests/application/update/update_test.go`
  - `tests/application/uuidgen/uuid_test.go`
  - `tests/domain/search/similarity_test.go`
  - `tests/domain/update/version_test.go`
  - `tests/infrastructure/filesystem/hawp_project_test.go`
  - `tests/infrastructure/githubrelease/githubrelease_test.go`
- The remaining `42` `_test.go` files under `internal/...` are concentrated in
  package-internal or integration-heavy areas such as:
  - `internal/application/context`
  - `internal/application/index`
  - `internal/domain/context`
  - `internal/domain/work`
  - `internal/infrastructure/sqlite`
  - `internal/platform/mcp`
- Verification passed after the second migration pass:
  - `go test ./tests/... ./internal/application/... ./internal/domain/... ./internal/infrastructure/... ./internal/platform/...`
  - `go test ./...`
  - `go run ./cmd/hawp work validate`

**Decision:**

- Treat `librarian/src/tests/...` as the home for exported black-box behavior
  tests
- Keep white-box and integration tests co-located until port seams or helper
  extraction make a move simpler than the current layout

## Outcome

Closed 2026-08-31.

The replicated test-topology lane is complete for `v0.0.23`. The repo now has
an explicit, verified boundary: exported black-box tests live under
`librarian/src/tests/...`, while package-internal or integration-heavy tests
stay beside the code until there is a cleaner seam.

## Verification

- [x] The replicated test tree now covers `14` black-box test files under
      `librarian/src/tests/...`. Evidence: [plan.md](/Users/beltrd/Desktop/projects/sentzunhat/human-ai-workflow-protocol/.hawp/work/closed/2026/08/31/0bd32051/plan.md)
      migration checkpoint lists the moved files
- [x] The remaining `42` tests intentionally stay co-located under
      `internal/...`. Evidence: [plan.md](/Users/beltrd/Desktop/projects/sentzunhat/human-ai-workflow-protocol/.hawp/work/closed/2026/08/31/0bd32051/plan.md)
      checkpoint records the retained package-internal and integration-heavy
      areas
- [x] The mixed topology passed verification after the migration pass. Evidence:
      [plan.md](/Users/beltrd/Desktop/projects/sentzunhat/human-ai-workflow-protocol/.hawp/work/closed/2026/08/31/0bd32051/plan.md)
      records `go test ./...` and `go run ./cmd/hawp work validate`

## Close Checklist

- [x] Black-box boundary exported into `src/tests`
- [x] White-box boundary explicitly retained in `internal/...`
- [x] Verification recorded
