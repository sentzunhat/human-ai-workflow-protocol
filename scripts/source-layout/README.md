# Source layout migration

Repository-specific migration tooling for the whole `librarian/src` tree.
It preserves `src/internal`, the existing layers, and cohesive package ownership.
Default mode is preview; applying requires an explicit, current reviewed plan.

## Organization

```text
scripts/source-layout/
  run.sh                     # launcher; no duplicated migration logic
  go.mod                     # independent, standard-library-only Go tool
  cmd/source-layout/         # executable entrypoint
  internal/
    mapping/                 # destinations, revision, and mapping-only tests
    migration/               # inventory, Go rewrites, preview, apply, safety
  review/
    README.md                # snapshot status
    preview.md               # previous generated proposal, marked stale
    plan.json                # previous fingerprints, rejected by new revision
```

The top-level setup script was redundant: use
`hawp mcp configure --provider codex --repo-root .` for MCP configuration.
The old work-only and source-layout shell entrypoints are removed; use `run.sh`.

## Checkpoint before running

The current cleanup adds the context `dedup/` destination and changes
the mapping revision. The [saved preview](review/preview.md) is the earlier
39-move snapshot, **not a regenerated or validated result for this revision**.
The [saved plan](review/plan.json) cannot be applied with the revised tool.

The requested ordering is: commit the cleanup, generate/check a new preview,
review its diff, then separately decide whether to apply. Do not treat the earlier
passing candidate compilation as proof of this deeper mapping.

## Preview commands (after the checkpoint)

From the repository root:

```sh
bash scripts/source-layout/run.sh --preview --check --write-plan review/plan.json --report review/preview.md
bash scripts/source-layout/run.sh --preview --plan review/plan.json --diff
```

The first command inventories all tracked/non-ignored source files and explicitly
writes only the requested review artifacts. It compiles/vets a temporary candidate,
including `integration,benchmark` test compilation, without running live models.
The second shows exact code/import changes against that same snapshot.
Source files are not moved or deleted in either preview command.
Go may populate its ordinary module/build caches.

The launcher resolves the repository root through Git and refuses to run if the
script is not inside the expected checkout. Relative flag paths resolve inside
`scripts/source-layout/`; absolute paths also work. Review artifacts cannot be
written into `librarian/src`. Reports describe a proposal, not completed
migration or proof that every layer is pure.

## Proposed ownership

Paths below are relative to `librarian/src`.

| Owner | Destination |
| --- | --- |
| Embedding and LLM ports | `internal/domain/providers/{embeddings,llm}/` |
| Work create/draft, validation, normalization | `internal/application/work/{intake,validation,normalize}/` |
| Context deduplication and adjacent tests | `internal/application/context/dedup/` |
| SQLite search index | `internal/infrastructure/repositories/index/` |
| Download and release clients | `internal/infrastructure/clients/{download,githubrelease}/` |
| MCP configuration and JSON-RPC server | `internal/platform/mcp/{configure,server}/` |
| Live model integration/benchmark tests | `internal/infrastructure/models/` |

Context's larger package has an independent deduplication owner. Configuration
stays cohesive because `encryption.go` defines methods on `ContextConfig`.
Rendering/reshaping remain together because they share private reference helpers.
MCP configuration has private cross-file dependencies and mixed-provider tests:
no additional provider folders are introduced until those contracts are extracted.
Do not nest small packages merely to reduce a file count.

The CLI hierarchy, command entrypoint, model adapters, and `internal/bootstrap`
remain in place. Bootstrap wires adapters outside domain; it is not another
business-policy layer. Configuration I/O still needs a later port/adapter
extraction; grouping its files does not make it pure.

See the [architecture map](../../.hawp/work/active/47c793d6/package-boundaries.md)
and [work-domain separation record](../../.hawp/work/closed/2026/09/10/742aa60b/plan.md).

## Safety and apply gate

- Syntax-aware rewrites update imports and split aliases while respecting local
  shadowing. Exported cross-file references are qualified; private cross-package
  dependencies and unresolved symbols are refused.
- Saved plans must match the mapping revision, every source path, hash, mode, and
  generated output. Changed/partial source or older mappings require fresh review.
- Occupied destinations, symlink paths, and ignored Go files are rejected.
  No file is considered dead merely because its name is not imported.
- Assets, tests, and module files are retained. Only successfully moved originals
  and their now-empty old folders are removed during explicit apply.

After the new preview is reviewed and source writers are stopped:

```sh
bash scripts/source-layout/run.sh --plan review/plan.json --apply
```

Apply compiles/vets the candidate, rechecks the source snapshot, backs up changed
originals, writes destinations exclusively, and checks output hashes. Ordinary
errors attempt rollback; termination may leave a partial migration with the
printed backup retained for manual recovery. This is not a crash-atomic or
concurrent-writer transaction. Rerunning the same applied plan is a no-op.

Run full source tests and HAWP checks after an actual migration. Native ORT builds,
live models, non-Go path references, and full migrated-tree behavior need their
own validation; compilation alone does not prove them.

## Development checks

Before the checkpoint, these check only syntax/types, flags, and destination
policy; they do not prepare a preview, call apply, or regenerate snapshots:

```sh
go -C scripts/source-layout test -run '^(TestDestination|TestCommandFlags|TestObsoleteMappingRejected)$' ./...
go -C scripts/source-layout vet ./...
bash -n scripts/source-layout/run.sh
```

The complete regression suite also creates temporary preview/apply fixtures:

```sh
go -C scripts/source-layout test ./...
```

Run it after the requested checkpoint. It covers rewrites, tags/assets, stale
plans, path safety, a fixture apply, post-write rollback, and already-applied state.
