### Source Layout Preview Checkpoint

Historical verification for the first mapping. The later scripts cleanup relocates
the tool and adds deeper context destinations; that newer revision has not run a
preview/apply. The previous snapshot is retained and labeled, not regenerated.

#### Checkpoint-First Scripts Cleanup

- `scripts/` now contains only `source-layout/`, with `cmd/source-layout`,
  `internal/mapping`, `internal/migration`, `review`, and one `run.sh` launcher.
- Removed the redundant MCP setup script; the README uses the existing
  `hawp mcp configure --provider codex --repo-root .` CLI. Removed the work-only
  wrapper and relocated the newer launcher under its actual owner.
- Added proposed `application/context/dedup` destinations, including adjacent
  tests. Context configuration remains cohesive because its type methods are
  split across files. Shared MCP and rendering helpers remain
  together. Mapping revision checks reject the old saved plan before any writes.
- Verified tool compilation, mapping/flag/revision-only tests, vet, shell syntax,
  current MCP/configuration/context tests, and diff whitespace. Full migration
  fixture tests and fresh preview/apply deliberately remain deferred until after
  the user-requested checkpoint commit.
- Commit scope includes the previously verified source/documentation work and the
  already-recorded retirement of `core/.hawp/bin/hawp`. No installed launcher or
  native executable is removed, refreshed, or published. Removed tracked scripts
  remain recoverable from Git history.

The sections below preserve the earlier verification record, not new-map proof.

#### Scripts Standards Update

Added [scripts/README.md](../../../../scripts/README.md) with the maintenance
standards, structural-tool pattern, and lessons learned from the source-layout
work. It documents why shell launchers stay thin, why tagged tests and mapping
revisions matter, and why semantic boundaries are not established by folder moves.

#### Intent

Make the whole-source migration inspectable and testable before applying any
layout changes. Preserve `librarian/src/internal` and its existing layer roots.

#### Current State

The preview inventories 252 files: 39 proposed moves, 42 content updates
(overlapping moved files), and 213 retained paths. No dead-file deletions are
proposed. The broad migration is **not applied** in the working repository.

See the [prior tree and file decisions](../../../../../../scripts/source-layout/review/preview.md),
[tool guide](../../../../../../scripts/source-layout/README.md), and
[architecture plan](../../../../active/47c793d6/plan.md).

#### What Was Inspected

Canonical layered-composition rules, the source inventory, Go declarations and
callers, MCP/work package boundaries, tagged integration tests, and stale workflow
links. The generated report distinguishes retained owners from unresolved policy/I/O.

#### What Changed

- Added an explicit preview, exact diff, Markdown report, fingerprinted plan,
  destination preflight, and regression coverage. The old work-only helper now
  forwards to `scripts/migrate-source-layout.sh`.
- Corrected stale adapter constructor imports and two missing `HybridRank`
  arguments in the existing tagged index integration test. This is a compilation
  repair, not a new semantic-ranking implementation.
- Updated port/work documentation and planned package locations; repaired stale
  links to the closed runtime-adapter item. Kept semantic work separation open.

#### What Was Directly Verified

- Tool unit tests and vet pass. Fixtures cover import splitting/shadowing,
  build tags/assets, unchanged previews, stale/partial plans, path collisions and
  symlinks, apply, post-write mismatch rollback, and already-applied detection.
- `bash scripts/migrate-source-layout.sh --preview --check` passes ordinary and
  `integration,benchmark` candidate compilation plus candidate vet.
- Saved-plan preview, exact diff, and compatibility launcher pass. Before/after
  source diffs match across preview execution. Mixed preview/apply flags are refused.
- Current source `go test ./...`, `go vet ./...`, and tagged index test compilation
  pass. Full default tests run on the current layout, not the migrated candidate.
- HAWP kit/work/links validation passes; `git diff --check` passes.

#### What Remains Unproven

The full real repository has not been applied/migrated. Live model execution,
native ORT builds, full runtime tests on the migrated tree, and interruption or
concurrent-writer recovery are not proved. Apply requires a reviewed snapshot and
exclusive source access; retained backups are not a crash-atomic transaction.

The existing `HybridRank` compatibility helper supplies no embedder factory;
MCP/benchmark callers need a separate semantic-ranking wiring review. Folder and
import changes do not repair that runtime behavior or make mixed work-domain I/O pure.

#### Constraints

No source file moves/deletions, commit, push, release, or binary refresh. Existing
`core/.hawp/bin/hawp` deletion is unrelated and preserved. Source edits in this
iteration are limited to the tagged test repair and documentation alignment.

#### Suggested Next Step

Review the generated destination tree and diff. Apply only after accepting that
mechanical map; then run the full suite on the migrated checkout before continuing
the separately tracked pure work-rule/filesystem extraction.
