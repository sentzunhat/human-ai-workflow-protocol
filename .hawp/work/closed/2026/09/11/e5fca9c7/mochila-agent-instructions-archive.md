# Mochila Archive Viewer HAWP Cleanup Instructions

Use these instructions in `/path/to/projects/mochila-archive-viewer`.

## Mission

Test the HAWP `v0.0.24` work-normalization path against Mochila's existing
`.hawp/work` state, produce a reviewable patch, and do not commit.

## Constraints

- Do not commit, push, merge, tag, or publish.
- Preserve unrelated dirty changes already present in the repo.
- Work on a new branch before applying changes.
- Keep scope limited to `.hawp/work/**` unless validation points to a real
  broken reference outside that tree.
- Do not rewrite historical closed records beyond the mechanical scaffolding
  needed for validation.
- Separate directly observed command output from inference.

## Starting Evidence

Read-only validation from the HAWP `v0.0.24` branch found:

- Active rows `048`, `050`, `051`, `052`, and `053` still appear active while
  their plans already live under `.hawp/work/closed/2026/09/01/`.
- The corresponding closed plans are missing modern `Outcome`, `Verification`,
  and `Close Checklist` sections.
- Flat active plans remain for numeric work items and can be migrated into
  folder-per-item layout.
- The repo had unrelated dirty changes at audit time; do not overwrite them.

## Commands

From the HAWP source repo, build or use the `0.0.24` binary:

```bash
cd /path/to/projects/human-ai-workflow-protocol
./.hawp/bin/hawp version
```

Expected version:

```text
0.0.24
```

In Mochila, create a local branch:

```bash
cd /path/to/projects/mochila-archive-viewer
git status --short --branch
git switch -c chore/hawp-work-normalization-2026-09-05
```

Dry-run first:

```bash
/path/to/projects/human-ai-workflow-protocol/.hawp/bin/hawp work normalize --dry-run --validate --hawp-root .hawp
/path/to/projects/human-ai-workflow-protocol/.hawp/bin/hawp work normalize --dry-run --migrate-folders --validate --hawp-root .hawp
```

If the dry-runs match the expected mechanical cleanup, apply in this order:

```bash
/path/to/projects/human-ai-workflow-protocol/.hawp/bin/hawp work normalize --apply --force-dirty --validate --hawp-root .hawp
/path/to/projects/human-ai-workflow-protocol/.hawp/bin/hawp work normalize --apply --migrate-folders --force-dirty --validate --hawp-root .hawp
/path/to/projects/human-ai-workflow-protocol/.hawp/bin/hawp work validate --hawp-root .hawp
```

Use `--force-dirty` only because this target repo is already dirty. Do not use
it to hide unrelated changes; inspect `git diff -- .hawp/work` afterward.

## Expected Result

- `BACKLOG.md` no longer lists completed rows `048`, `050`, `051`, `052`, and
  `053` under Active Work.
- The closed records for those IDs have `Outcome`, `Verification`, and
  `Close Checklist` scaffolding.
- Remaining active flat files such as `active/035.md` are moved to
  folder-per-item paths such as `active/035/plan.md`.
- Final `hawp work validate --hawp-root .hawp` passes.
- No commit is made.

## Report Back

Return:

- branch name
- exact commands run
- final `git status --short --branch`
- final validation summary
- list of changed `.hawp/work` files
- any skipped or manually reviewed items
