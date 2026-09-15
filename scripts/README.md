# Scripts guide

This directory contains repository maintenance tools. Scripts should be small,
explicit, reviewable, and safe to run from the repository root.

## Current layout

```text
scripts/
  README.md                 # shared standards and lessons
  source-layout/            # whole-source-tree migration tool
    run.sh                  # only shell launcher
    cmd/                     # executable entrypoints
    internal/mapping/        # reviewed destination policy
    internal/migration/      # inventory, rewrite, preview, safety, apply
    review/                  # generated or checkpointed review artifacts
```

There are no competing top-level migration wrappers. The source-layout tool is
the owner of source-tree migration behavior; its launcher is
`scripts/source-layout/run.sh`.

## Standards

- Run from an explicit repository root and resolve paths before reading or
  writing files. Do not depend on the caller's current directory.
- Default to read-only preview. Any mutation requires an explicit flag and a
  reviewed, current plan.
- Keep shell launchers thin. Put policy and transformations in tested code,
  preferably Go for this repository's Go source tree.
- Use standard-library-only tooling unless a dependency is clearly justified.
  Keep an independent `go.mod` when the tool is intentionally separate from
  the product module.
- Inventory the complete scoped tree. Git filename absence is not proof that a
  file is dead; retain assets, tests, module files, and support files unless a
  removal has direct evidence.
- Refuse unsafe paths, symlink traversal, occupied destinations, ignored source
  files, changed snapshots, duplicate destinations, and private cross-package
  dependencies.
- Make transformations syntax-aware. Update imports and identifiers through Go
  syntax trees; never use broad text replacement for code or silently rewrite
  comments and string literals.
- Keep tests beside the behavior they verify. Include fixture tests for preview,
  deterministic output, stale plans, rollback, already-applied state, and
  source preservation.
- Separate mechanical topology changes from semantic behavior changes. A new
  folder or import path does not prove that policy, filesystem I/O, or domain
  purity has been separated.
- Report exactly what was inspected, changed, verified, skipped, and left
  unproven. Preserve backups and recovery instructions for destructive steps.
- Do not refresh binaries, publish, push, merge, or modify external consumer
  configurations unless the task explicitly includes that action.

## Script pattern

Use this sequence for a structural tool:

1. Inspect the repository root, active plan, callers, tests, and dirty state.
2. Define a versioned mapping/policy with an explicit reason for each move.
3. Inventory every file in scope and record source, destination, mode, and
   before/after fingerprints.
4. Generate a preview report and exact diff without mutating the source tree.
5. Validate a temporary candidate with focused checks, then broader checks.
6. Apply only from the reviewed plan after rechecking the source snapshot.
7. Verify output hashes, retain a recovery backup, and document remaining risk.

## Lessons learned

- A preview-only shell map was insufficient for a whole-tree Go migration. The
  replacement needs an executable inventory, import rewrite, temporary compile
  check, and guarded apply path.
- Tagged tests are part of the source contract. Default compilation alone missed
  stale model constructor imports and missing `HybridRank` arguments.
- A folder move can split one package into several packages. Exported references
  can be qualified automatically; private cross-package references require a
  semantic extraction decision and must stop the tool.
- Destination safety must include ignored files and symlink ancestors, especially
  on case-insensitive filesystems. Preflight must happen before candidate apply.
- Saved plans need a mapping revision in addition to file hashes. A plan can be
  internally consistent yet stale after ownership rules change.
- One extra folder level is useful only when it represents a cohesive capability:
  context deduplication qualifies; configuration must stay cohesive with its
  type methods, while shared rendering/reshaping
  and MCP private contracts do not yet.
- Historical reports are evidence for their original mapping only. Move them into
  a clearly labeled review area and refuse to apply them after a mapping revision.
- Existing product CLIs should own configuration behavior. A second shell script
  for MCP setup duplicated logic and was removed in favor of `hawp mcp configure`.

## Checks

From the repository root:

```sh
go -C scripts/source-layout test ./...
go -C scripts/source-layout vet ./...
bash -n scripts/source-layout/run.sh
```

After committing tool/documentation changes, generate the source-layout preview:

```sh
bash scripts/source-layout/run.sh --preview --check \
  --write-plan review/plan.json --report review/preview.md
bash scripts/source-layout/run.sh --preview --plan review/plan.json --diff
```

Preview artifacts describe proposed changes. They are not proof that an apply was
performed or that live model/native-runtime behavior was exercised.
