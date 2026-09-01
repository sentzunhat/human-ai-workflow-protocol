# Raw Output — No HAWP

Prompt given to the agent (verbatim, no framing added):

> Repository: human-ai-workflow-protocol
>
> Is HAWP ready for a v1.0 release? What do we still need to do?
>
> Return your complete answer as your final response.

Output captured as-is below.

---

## Short answer

**Not yet — but it's close on substance.** The content and tooling are in good shape; what's missing is release mechanics (versioning, tagging, main-branch sync) plus a handful of known defects and in-flight work that should land first.

## What's already in good shape

- **Backlog is clean.** No active, blocked, or parked items in `.hawp/work/BACKLOG.md`; the last big efforts (provider rollout TASK-072–075, audit remediation TASK-076, scripts alignment TASK-077) are all closed.
- **Everything validates on the current tree.** I ran the repo's own gates: TypeScript typecheck passes, all 38 librarian unit tests pass, and `providers:validate`, `distribution:validate`, and `workflow:validate` all pass (the latter with one tolerated legacy-file warning from April).
- **Distribution works.** Three provider packs (GitHub/Copilot, Cursor, Continue) with generated install/update guides that are in sync with their sources, plus two CI workflows guarding drift and quality.
- **Evidence of value exists.** The benchmark folder now has a completed comparison run (2026-06-11) with a written-up result.

## What still needs to happen before a v1.0

**1. Land the in-flight work.** There are 34 modified files plus untracked additions sitting uncommitted: a kit restructure moving `spec.md` and `authoring-patterns.md` into `references/` (~841 lines deleted across mirrors and generated guides) and the new benchmark run. The validators pass against this state, but it needs to be committed — and notably there's no backlog row tracking it, which goes against the repo's own rule that work started outside the loop still gets a row.

**2. Decide what "v1.0" even means here.** The protocol spec is *deliberately* locked at v0.1 — `references/spec.md` is titled "HAWP v0.1 Specification" and `start-here.md` says "The shape is locked to v0.1." There are currently **no git tags, no CHANGELOG, and `librarian/package.json` is at 0.0.0**. You need an explicit decision: does v1.0 mean promoting the shape spec itself to 1.0 (declaring it stable), or versioning the repo/distribution separately from the protocol? Either is fine, but right now releasing "HAWP 1.0" with a spec that says v0.1 would be contradictory.

**3. Sync `dev` to `main`.** `dev` is **68 commits ahead of `origin/main`**, and your README points "Stable" users at main. A release means merging dev → main and confirming both CI workflows are green there.

**4. Triage the defects your own benchmark found.** The 2026-06-11 run surfaced real issues in `librarian/`, and none of them have backlog rows yet. I verified two are still live:

- `npm test` is broken on Node 20 (the `--test` glob pattern needs Node 21+) while `engines` claims `>=20` — I reproduced this locally just now.
- A hardcoded `"hawp backlog upgrade v1.1.0 (TASK-028 dry-run detection)"` string still sits in `backlog-upgrade/cli.ts` line 235 — confusing user-facing versioning for a repo that has never tagged a release.

The run also flagged cross-domain imports violating the documented folder-boundary rules (in two directions), `--apply --export-plan` silently ignoring the export flag, and the distribution path-leak check silently skipping missing files. At minimum, the engines mismatch and the phantom version string should be fixed before anything gets called 1.0.

**5. Smaller polish items.**

- The README roadmap notes UUID-based work item IDs (TASK-013) as still pending — fine to defer past 1.0, but say so explicitly.
- The benchmark README's installation note ("copy `.github/` and `.hawp/` into `core/`") reads inverted/stale relative to how installation actually flows.
- Only one benchmark run exists; a second run on a different task type would strengthen the value claim, though that's optional.

## Suggested release checklist, in order

1. Commit and close out the `references/` restructure and benchmark run (with a backlog row).
2. File and fix the benchmark-found defects, at minimum the Node engines mismatch and the `v1.1.0` string.
3. Decide and document the versioning story (protocol spec version vs. repo release version), add a CHANGELOG, and set up tagging.
4. Merge dev → main, confirm CI green, tag the release.
