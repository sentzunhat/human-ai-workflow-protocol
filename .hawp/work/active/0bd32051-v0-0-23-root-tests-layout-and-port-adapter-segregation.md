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

## Initial Analysis

**Directly verified:**

- The user explicitly wants test files moved into a replicated root-level
  `tests/` folder beside `src/`
- That change touches folder topology and verification layout, so it should be
  handled as its own work item rather than piggybacked onto the current
  formatter pass

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
**Gate:** plan/investigate first; do not merge

## Backlog + Plan Link

**Status now:** analyzing
**Plan file:** work/active/0bd32051-v0-0-23-root-tests-layout-and-port-adapter-segregation.md

## Next Step

- [x] Work item created from the requested next compoundable slice
- [ ] Inspect current Go package/test constraints for a replicated `tests/` tree
- [ ] Propose the smallest safe migration plan and boundary cleanup set
- [ ] Implement only after the migration path is explicit
