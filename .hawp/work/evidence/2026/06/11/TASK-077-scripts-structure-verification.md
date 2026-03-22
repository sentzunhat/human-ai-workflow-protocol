# Evidence: TASK-077 — Scripts structure guidelines and boundary alignment

Date: 2026-06-11
Verified by: agent (Cursor session), commands run locally from repo root.

## Direct evidence (commands and observed output)

- `npm --prefix librarian run typecheck` — exit 0, no diagnostics.
- `npm --prefix librarian test` — 38 tests, 38 pass, 0 fail (new `parseArgs` unit test for `backlog-validate/cli.ts` included).
- `npm --prefix librarian run validate:backlog` — `Result: PASS` (distribution + workflow validators).
- `./.hawp/bin/hawp backlog upgrade --dry-run --validate` — 0 auto-fixable, 0 blocked, `VALIDATION PASS`.

## Changes verified

- `backlog-validate/` split from a single monolithic `index.ts` into the documented boundary pattern: `index.ts` (executable boundary), `cli.ts` (arg parsing + help), `script.ts` (combined validator logic). Behavior unchanged; help test still passes.
- `validate-hawp-workflow/` gained the same split: new `script.ts` exports `runWorkflowValidation(argv): number`; `index.ts` reduced to the executable boundary setting `process.exitCode` (replaced scattered `process.exit` calls).
- New `librarian/scripts/README.md` documents the domain map, the `index/cli/script` boundary rules, `lib/` constraints, conventions (node: imports, error surfacing, path containment, fail-closed guards), and the procedure for adding a new script/CLI subcommand.
- `librarian/README.md` links the new guidelines.
- Backlog reconciled: no active or parked items; Recently Closed capped at 10 rows; STATUS.md points at the latest closed item.

## Inference / not directly proven

- CI runs on GitHub-hosted runners (workflows unchanged today; still unexecuted remotely from this session).
