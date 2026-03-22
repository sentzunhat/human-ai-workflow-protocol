# Librarian Product Direction Alignment

Related work:

- `.hawp/work/active/1c2db85b-d57b-47a1-8c4b-b28dc55f2d5c.md`
- `.hawp/work/active/8dedc4e2-69c5-42ca-aaad-93b62d7fb899.md`
- `.hawp/work/active/1a3b32a4-ab37-4c86-ade7-71e2eb42b440.md`
- `.hawp/work/active/54e68af7-4622-4383-8482-cc4d4e1e21ee.md`

## Purpose

Record what the repository says today about `librarian/`, what our current
active work says about the future direction, and what should change next without
creating docs drift.

## Directly Verified Current State

- Root `README.md` still describes `librarian/` as repo-local distribution
  validation and backlog tooling.
- `README.md` explicitly says `librarian/` is a repo-local source tree for HAWP
  maintenance and should not be installed separately.
- `librarian/README.md` documents the existing Node 26 TypeScript maintenance
  tooling truthfully: distribution sync, provider materialization, work
  validation, kit validation, and the repo-local `./.hawp/bin/hawp` wrapper.
- `librarian/README.md` now also documents the experimental Node CLI/binary PoC
  under `scripts/hawp-cli-poc/`.
- `.github/workflows/librarian-quality.yml` still runs the Node-based quality
  path, including `npm test` and a stale command name `workflow:validate`.
- Active HAWP notes now describe a future Go librarian direction focused on
  `db init`, ingesting `.hawp/work/` and `.hawp/kit/`, lexical/vector search,
  local context building, and prompt handoff.

## Inference

- The public docs are not wrong today; they describe the current shipped
  TypeScript tooling.
- The Go librarian direction is still planning and strategy, not current
  repository behavior.
- Changing root or librarian public docs to describe the Go librarian as current
  would create premature drift.
- The best next compoundable item is to reconcile the parent librarian research
  lane around a staged transition:
  1. keep current Node maintenance tooling true and stable
  2. define the Go librarian product surface clearly
  3. decide what eventually moves, stays, or coexists

## What Should Change Next

- Update the parent librarian research item to reflect the Go librarian product
  direction and the reality that current docs remain correct for today.
- Keep public docs conservative until there is an implemented Go workspace,
  command surface, or migration plan.
- Treat the next highest-leverage work as product-definition and transition
  planning, not broad public README rewriting.

## Candidate Next Item

Refine `1c2db85b` into the parent decision lane for:

- current Node maintenance tooling boundaries
- future Go librarian boundaries
- DB/index/model lifecycle
- transition plan for `work/` + `kit/` intelligence features
- what stays in repo-local Node scripts vs what becomes the Go librarian product
