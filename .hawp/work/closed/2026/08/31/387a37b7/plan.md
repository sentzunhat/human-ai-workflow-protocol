# v0.0.23 refactor search pipeline and layout unification

**Backlog ID (Legacy):** — (UUID-native item)
**UUID:** `387a37b7-e8e4-4c6b-9972-3c9712bc9774`
**Type:** improvement
**Reported:** 2026-08-29

---

## Input (verbatim)

> Do the code audit refactoring needs please on a patch version branch and multiple agents working on the work and benchmark don't merge

## Intake Summary

Implement the most compoundable architecture cleanup from the current v0.1.0
audit on a non-merge patch branch: unify the active `.hawp` runtime layout
contract and preserve richer search provenance through the typed formatting path.

## Current Context

- Current branch for this lane: `feature/v0.0.23-refactor-v010`
- Repo-root proof captured before edits:
  - `pwd` -> `<repo-root-abs>`
  - `git rev-parse --show-toplevel` -> `<repo-root-abs>`
  - `git rev-parse --show-prefix` -> empty output (repo root)
  - `git status --short` -> clean before this work item was opened
- Audit findings already confirmed in code:
  - active runtime code reads `.hawp/db/index.sqlite` and `.hawp/config/context.json`
  - filesystem helper + generated README text still describe `.hawp/.data/db`
    and `.hawp/.data/config`
  - typed search results define line/provenance fields that are not preserved
    through `FormatAsMarkdown` and reference deduplication
- Parallel lanes opened:
  - benchmark/evidence lane: `.hawp/work/active/c804eec0/plan.md`
  - follow-up audit lane: `.hawp/work/active/5957aaf4/plan.md`

## Initial Analysis

**Directly verified:**

- `.hawp/work/BACKLOG.md` now tracks three separate v0.0.23 lanes
- `librarian/src/internal/application/context/config.go` loads project config from
  `.hawp/config/context.json`
- `librarian/src/internal/application/search/service.go` and
  `librarian/src/internal/platform/cli/run.go` use `.hawp/db/index.sqlite`
- `librarian/src/internal/infrastructure/filesystem/hawp_project.go` and
  `librarian/src/internal/infrastructure/filesystem/readme_generator.go` still
  point at `.hawp/.data/...`
- `librarian/src/internal/domain/search/result.go` includes `ChunkID`,
  `LexicalRank`, `SemanticScore`, `LineStart`, and `LineEnd`, but
  `librarian/src/internal/application/context/format.go` does not preserve the
  line fields

**Inferred (not yet proven):**

- Aligning the helper/runtime layout first should reduce future drift when the
  reshape path is exposed in CLI
- Preserving line provenance now should make later reshape/reference output safer
  without forcing the full reshape-path change into this patch

**Likely scope:**

- Update the filesystem helper and generated README text to match the runtime
  `.hawp/{db,config}` contract
- Preserve line/ranking provenance through typed search results and context
  reference formatting
- Add focused tests for the layout contract and context/reference provenance
- Keep reshape-path exposure, benchmark artifacts, and broader CLI splitting in
  separate lanes

## Risk + Review Gate

**Risk:** medium (shared search/context contracts)
**Gate:** user approved implementation on a patch branch; do not merge

## Backlog + Plan Link

**Status now:** done
**Plan file:** work/closed/2026/08/31/387a37b7/plan.md

## Next Step

- [x] Investigation recorded above
- [x] Plan direction captured in-place for this bounded patch lane
- [x] Implement scoped refactor + focused tests
- [x] Verify and summarize without merging
- [ ] Continue sparse-result shaping cleanup until token benchmark negatives are removed or explicitly justified
- [ ] Move to the next shared-search consolidation slice after the shaping pass

## 2026-08-29 Compact Render Follow-Up

**Directly verified:**

- `librarian/src/internal/application/context/format.go` now keeps sparse
  markdown output body-only (no title header at 5 results or fewer) and emits
  compact inline provenance as `Ref: <source[:line[-line]]>`
- richer title/relevance metadata is still preserved on typed
  `FormattedResult` / `DocumentReference` shapes and JSON output
- focused verification passed:
  - `cd librarian/src && go test ./internal/application/context ./internal/application/search`
- source-run token benchmark still shows three sparse negatives:
  - intake ordering: `1707 -> 1753` (`-46`, `-3%`)
  - provider distribution: `1796 -> 1870` (`-74`, `-4%`)
  - core shape fields: `591 -> 613` (`-22`, `-4%`)

**Inference:**

- the cleanup improved the formatting boundary and kept the shaping change
  minimal, but fully removing the remaining sparse negatives likely needs a
  deliberate sparse passthrough policy rather than more wording compression

## 2026-08-31 Release Checkpoint

**Directly verified:**

- the source-backed rerun of `hawp search benchmark --token` now reports
  aggregate savings of `4838` tokens (`21%`)
- prior sparse-query negatives were reduced to `0` in two cases and a single
  negligible `-1` token case in the third
- `go run ./cmd/hawp check` passes on the live source path after the
  verification-clarity cleanup

## Outcome

Closed 2026-08-31.

The search/layout unification lane is complete for `v0.0.23`. Runtime layout
drift was resolved, sparse shaping is now effectively passthrough for small
contexts, typed provenance stays intact, and the benchmark rerun no longer
shows materially negative sparse-query behavior.

## Verification

- [x] Shared search/runtime cleanup shipped in the `0.0.23` changelog. Evidence:
      [CHANGELOG.md](/Users/beltrd/Desktop/projects/sentzunhat/human-ai-workflow-protocol/librarian/src/CHANGELOG.md)
      `shared search execution path` and `sparse context shaping to passthrough`
- [x] Branch-local benchmark proof is recorded. Evidence:
      [.hawp/work/evidence/2026/08/31/c804eec0-v0-0-23-benchmark-rerun.md](/Users/beltrd/Desktop/projects/sentzunhat/human-ai-workflow-protocol/.hawp/work/evidence/2026/08/31/c804eec0-v0-0-23-benchmark-rerun.md)
      shows `21%` aggregate savings and only a `-1` sparse outlier
- [x] Source-backed validation is green. Evidence:
      [.hawp/work/evidence/2026/08/31/c804eec0-v0-0-23-benchmark-rerun.md](/Users/beltrd/Desktop/projects/sentzunhat/human-ai-workflow-protocol/.hawp/work/evidence/2026/08/31/c804eec0-v0-0-23-benchmark-rerun.md)
      records `hawp check` and `work validate` pass results

## Close Checklist

- [x] Layout/runtime drift resolved
- [x] Sparse shaping checkpoint justified by rerun evidence
- [x] Follow-up implementation moved out of the release lane
