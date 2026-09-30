# Restore coordination and harness guidance to the canonical kit

**UUID:** `a91cfdbe-5a15-4475-b588-ea82d3bec588`  
**Type:** improvement  
**Status:** done  
**Opened:** 2026-09-30  
**Closed:** 2026-09-30  
**Risk level:** low (documentation only)

## Input

The Codex provider update removed the local kit copies of the parallel-agent
worktree, manager-branch, HAWP-first workflow, and slice/provider/tool harness
guides. Check whether these concepts moved elsewhere and preserve the useful
guidance in the canonical core kit.

## Investigation

### Directly verified

- `core/.hawp/kit/start-here.md` now includes HAWP MCP context guidance: prefer
  `hawp_search`, check the repository scope, reuse existing plans, and validate
  workflow changes. `core/.hawp/kit/usage/search.md` covers CLI/MCP search.
- `core/.hawp/kit/references/work-item-file-tracking.md` and
  `core/.hawp/kit/patterns/parallel-work-guardrails.md` retain file ownership
  and collision-prevention guidance, but not an end-to-end worktree workflow.
- The manager-branch, parallel-agent-worktrees, HAWP-first, harness-guide,
  slice-harness, provider-harness, and tool-harness documents are absent from
  `core/.hawp/kit/`.
- The current MCP usage index links to the removed
  `examples/mcp-intake-to-work-doc.md`, so the provider refresh also leaves a
  broken in-kit link unless that current example is restored.
- The removed TypeScript scoped instruction overlaps the canonical Node.js
  style guide on `node:` imports, but its TypeScript module-resolution,
  repo-root script, and CLI-boundary rules are not all present there.
- Historical versions remain in the root kit at `HEAD` and work records
  `.hawp/work/closed/2026/08/25/manager-branch-kit-pattern/plan.md` and
  `.hawp/work/closed/2026/09/10/74aaa332/plan.md`; these are not current
  downstream kit sources.

### Assessment

Some HAWP-first context rules moved into `start-here.md`, and basic file
ownership concepts remain. Manager-branch coordination, safe worktree setup,
and the slice/provider/tool harness contracts no longer have current portable
guides in the source kit. The historical documents are useful design evidence,
but need current links, safer cleanup guidance, and current package/contracts
before they are restored. The old TypeScript rule partly duplicates the
canonical Node.js style guide; its useful module/script boundaries can be
preserved in a narrower file-scoped instruction without enforcing obsolete
project-specific npm script names.

## Mission

Restore portable, current documentation for HAWP-first context gathering,
parallel worktrees and manager branches, and focused slice/provider/tool
verification under `core/.hawp/kit/`. Keep the protocol schema unchanged and
make the refreshed root kit match the canonical source.

## Constraints

- No runtime, CLI, MCP schema, or provider behavior changes.
- Keep manager branches and parallel agents optional; avoid claims about a
  specific agent product's worktree defaults.
- Update package references and verification guidance to current HAWP layout.
- Preserve existing `.hawp/work/**` records; add only this task's backlog row
  and plan. Keep the pre-existing local helper-binary state out of documentation
  changes.
- Keep machine-local paths out of durable documents.

## Plan

1. Restore focused guides and standards under canonical `core/.hawp/kit/`.
2. Link them from `core/.hawp/kit/start-here.md` and the standards index.
3. Refresh root `.hawp/kit/` from the canonical core, preserving work records.
4. Verify cross-links, HAWP work validation, source/root kit parity, and a clean
   diff check for the documentation change.

## Acceptance

- The seven named concepts have current, linked documentation in `core/.hawp/kit/`.
- Existing kit links to the MCP intake example resolve again.
- The scoped TypeScript instruction retains useful rules and defers module
  syntax to each project's TypeScript configuration.
- Start-here explains which guide to use for context, coordination, and harness
  work without duplicating the full guide contents.
- The root kit mirrors `core/.hawp/kit/`; `.hawp/work/**` is unchanged.
- Documentation links and HAWP work validation pass.

## Implementation Notes

Restored the current guides in `core/.hawp/kit/` for HAWP-first context,
parallel worktrees, optional manager branches, slice harnesses, provider
harnesses, and CLI/MCP tool harnesses. Restored the MCP intake example required
by its existing index link and promoted the useful TypeScript scoped rules into
`instructions/typescript-conventions.md`, deferring module extensions to each
project's `tsconfig.json`.

Updated `start-here.md` and `standards/README.md` to surface the guides, then
refreshed root `.hawp/kit/` from `core/.hawp/kit/`. The context-search basics
were already partially moved into `start-here.md`; the worktree/manager and
provider/tool/slice harness details had no current core-kit replacement.

## Outcome

All named guidance concepts and the missing MCP example now have canonical,
linked files in `core/.hawp/kit/`. The core kit and root kit match. No HAWP
schema, runtime, CLI, or MCP behavior changed.

## Verification

- `go run ./cmd/hawp links check` from `librarian/src`: 144 Markdown files
  checked; local links valid.
- `go run ./cmd/hawp kit validate`: 3 checks passed, 0 issues.
- `go run ./cmd/hawp work validate`: validation passed with 0 issues; the
  repository-wide verification-clarity scan reported 25 ambiguous items.
- `git diff --check` passed; `diff -qr core/.hawp/kit .hawp/kit` and `cmp` for
  the LICENSE found no differences.

## Close Checklist

- [x] Current coordination and harness guides are in the canonical core kit.
- [x] Start-here and standards indexes link to the restored guidance.
- [x] Root kit matches the core source; `.hawp/work/**` existing records remain.
- [x] HAWP link, kit, and work validation completed with results recorded.
