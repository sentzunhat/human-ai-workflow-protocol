# Split the website and Librarian changes into separate main-targeted PRs

**Backlog ID (Legacy):** — (UUID-native item)
**UUID:** `f8685e56-9bfc-4272-b102-5501bcffae8e`
**Type:** task
**Reported:** 2026-10-22

---

## Input (verbatim)

> yes i do approve, and also approve the branches and pr splitting and maybe add or update work items to make sure we know what to do next and also start with the splitting and then the option a thank you

## Intake Summary

The requested split is for the current HAWP website and Librarian/MCP changes.
The existing GitHub PR is closed and targets the deleted `development` base;
the current `feature/website` branch contains one mixed commit. No open PR
exists. The split needs two new PRs targeting `main`, but preserving the exact
`feature/website` name would require changing its current remote branch history.

## Current Context

- Current branch: `feature/website`; local and `origin/feature/website` point
  to `e2237b480c01e3707d618685e4f8e6f913d5a263` (`wip`).
- Local `main` and `origin/main` point to `178fde6ebcb57f3e8100710c72ae71f06e62fb6d`.
  The current website branch is one commit ahead of that `main`.
- The current committed comparison from `main` contains 230 changed paths:
  124 under `website/`, 52 under `.hawp/`, 22 under `distribution/`, 13 under
  `.github/`, 9 under `librarian/`, 7 under `core/`, and one each under
  `.gitignore`, `benchmark/`, and `scripts/`.
- Of the 124 `website/` paths, 86 are generated files under
  `website/.svelte-kit/` and `website/build/`; the remaining 38 are source,
  content, package, static, and configuration files.
- `website/.gitignore` excludes both generated directories. The site workflow
  installs dependencies, type-checks, builds, and uploads `website/build/`;
  it does not need the checked-in generated output as source.
- The 47 changed `.hawp/work/` paths include active, parked, closed, and status
  records from several unrelated work items. They must not be swept wholesale
  into either product PR.
- The 56 non-website source/provider/runtime paths in the committed comparison
  are enumerated in the file-ownership manifest. Current unstaged provider,
  agent, documentation, configuration-generator, and test edits are assigned
  to the Librarian PR there.
- GitHub reports no open PRs. PR #45, titled “Add the HAWP website and
  MCP/provider workflows,” is closed. Its base is `development`; its recorded
  head at closure was `b4922c69535afb5e981ad00c34e614e17dbb8061`, which differs
  from the current remote `feature/website` tip. Do not reopen or edit that
  closed PR as if it represented the current branch.
- The current worktree has 16 modified tracked files and untracked HAWP work
  records. These are provider/agent/MCP and search-work changes, not website
  source changes. Nothing is staged.
- `.github/workflows/website.yml` currently runs its push preview only for
  `feature/website`. A new website branch name needs a matching update to its
  push-branch filter and preview condition.
- No active or parked work item for splitting this PR was found. The new
  record is `f8685e56`. `hawp_work_intake` returned `needs_user_input` because
  the local index had no matching context; the branch and PR details above
  were confirmed directly from Git and GitHub.
- The user confirmed two new branches: `feature/website-improvements` and
  `feature/mcp-providers`. Both will start from `origin/main`; the current
  `feature/website` is the source/reference and will remain unchanged.
- `hawp/hawp_search` returned no results for the file-ownership search query.
  Repository files and Git path lists were used as the source of truth.
- The user previously confirmed that a fresh Copilot chat has the generated
  HAWP agent selected. Do not ask them to select it again unless the session
  changes. The live Tools list, client approval behavior, and VS Code MCP
  migration preview remain unverified.
- A tool-discovery query for VS Code UI/session inspection surfaced HAWP MCP
  tools only; no client UI inspection tool was available in this conversation.
  This is not live-session evidence.

## Initial Analysis

**Directly verified:**

- The previous PR #45 is closed against `development`, and no PR is currently
  open.
- The current `feature/website` remote tip differs from the frozen PR #45 head
  and contains both website files and HAWP/MCP/provider files.
- There are 38 non-generated website paths, one Pages workflow, and 86
  generated website paths in the current branch comparison.
- The website workflow builds the site from source and currently scopes its
  push preview to `feature/website`.
- The root `.gitignore` delta adds local Claude and Continue MCP configuration
  paths, so it belongs with the code/provider PR, not the website PR.
- The current worktree changes are unstaged and must be preserved as code-side
  changes.
- The current `.hawp/work/` branch delta spans 47 paths and includes unrelated
  artwork, parked website plans, historical records, and MCP-related plans.
  The shared backlog therefore needs a minimal, selected delta rather than
  wholesale copying.

**Inferred (not yet proven):**

- The 86 tracked generated website files can be omitted from a source-based
  website PR because the workflow regenerates `website/build/` before
  deployment. Verify the filtered branch build before opening the PR.
- The confirmed `feature/website-improvements` branch preserves the existing
  remote `feature/website` ref, but the workflow's preview filter must be
  updated as part of the website PR.
- `.hawp/work/BACKLOG.md` should have one owner: the Librarian PR. Rebuild its
  delta from `main` using only the work records included with that PR; do not
  copy the full mixed branch delta or duplicate the backlog edit in the
  website PR.

**Likely scope:**

- Prepare independent branches from `main`: one website/Pages branch and one
  Librarian/MCP/provider branch. Use the exact, non-overlapping ownership
  manifest in `files.md`; do not move files between canonical directories.
- Include current uncommitted provider/agent/documentation changes on the code
  branch. Keep the shared backlog exclusively on that branch and reconstruct
  its small delta from `main`.
- Exclude the 86 generated files under `website/.svelte-kit/` and
  `website/build/`; retain the 38 website source/configuration/content paths.
- Update `.github/workflows/website.yml` to trigger previews for
  `feature/website-improvements`. This workflow path belongs only to the
  website PR.
- Keep unrelated benchmark output, the private LinkedIn artwork task, and
  unselected HAWP history on the preserved source branch for separate review;
  do not delete, stage, or silently fold them into either PR.
- Create new PRs targeting `main` with separate titles and change summaries.
- Leave version and release files unchanged; no version bump was requested in
  this split task.

## Risk + Review Gate

**Risk:** medium
**Gate:** The user approved splitting and confirmed the two new branch names.
Use `origin/main` as each branch base and `feature/website` only as the source
reference. Do not rewrite or force-push the existing `feature/website` branch.

## Repository-Root Proof

Captured before the file-ownership manifest edit. Machine-local prefixes are
redacted; repository-relative status paths are unchanged.

```text
$ pwd
<repo-root-abs>
$ git rev-parse --show-toplevel
<repo-root-abs>
$ git rev-parse --show-prefix
<empty output>
$ git status --short
 M .github/agents/docs-alignment.agent.md
 M .github/agents/hawp.agent.md
 M .github/agents/master-builder.agent.md
 M .github/copilot-instructions.md
 M .github/skills/hawp-change-review-and-reference-sync/SKILL.md
 M .hawp/kit/usage/mcp/README.md
 M .hawp/kit/usage/mcp/github-copilot.md
 M .hawp/work/BACKLOG.md
 M .hawp/work/active/b152c56c/plan.md
 M .hawp/work/status/2026/10/21/b152c56c/status.md
 M core/.hawp/kit/usage/mcp/README.md
 M core/.hawp/kit/usage/mcp/github-copilot.md
 M core/providers/.github/prompts/hawp-change-review-and-reference-sync.prompt.md
 M librarian/src/internal/platform/mcp/configure/config.go
 M librarian/src/internal/platform/mcp/configure/config_agent.go
 M librarian/src/internal/platform/mcp/configure/configure_test.go
?? .hawp/work/active/09ee9986/
?? .hawp/work/active/f8685e56/
?? .hawp/work/evidence/2026/10/22/
```

## Verification So Far

- `git diff --check` passed after the file-manifest edits.
- `hawp/hawp_work_validate` ran. Kit validation passed 3 checks with 0 issues.
  Work validation reported 2 issues and 1
  warning; link checking found `.hawp/work/BACKLOG.md` still points to missing
  `.hawp/work/active/2eea565c/plan.md` and
  `.hawp/work/active/429e075e/plan.md`.
- These missing-plan links were reported in earlier validation and were not
  introduced or repaired by this split-planning change. Keep the unrelated
  work-record delta out of the split until each record is reconciled.

## Independent Client-Verification Gate

Repository files and HAWP MCP tool-search results do not prove the live
selected-agent Tools list or VS Code UI behavior. Before making claims about
client availability or reconciling config:

- Inspect the selected agent's live Tools list against the expected seven
  names: `hawp/hawp_search`, `hawp/hawp_usage`,
  `hawp/hawp_work_intake`, `hawp/hawp_work_new`,
  `hawp/hawp_work_validate`, `hawp/hawp_work_doc`, and
  `hawp/hawp_work_reshape`.
- Inspect approval/confirmation behavior without invoking a write-capable
  tool. Do not infer approval from agent text or tool availability.
- Read the existing `hawp` entry's migration-preview effect (preserve, merge,
  or replace). Do not apply an ambiguous migration.
- If the VS Code UI remains inaccessible, request the live Tools list,
  approval behavior, and migration-preview text or screenshot from the user.
- Keep `b152c56c` open until its live tool-list and approval evidence is
  recorded. Keep `cb096b6a` blocked on migration-preview behavior.
- Treat `70d8112f` as a separate medium-risk provider-installer audit; present
  findings and a minimal plan before editing shared provider code or docs.

This planning turn invoked only read-only `hawp/hawp_search` and
`hawp/hawp_work_validate`; neither establishes those client behaviors.

## Backlog + Plan Link

**Status now:** in-progress (branch strategy confirmed; preparing isolated branches)
**Plan file:** work/active/f8685e56/plan.md

## Split Plan for Review

The exact file ownership, pending/current distinction, and do-not-touch list
are in [files.md](files.md). Both new branches now exist as worktrees based on
`origin/main`; the existing `feature/website` reference branch remains
unchanged. Selected files have been copied into the new worktrees without
staging or committing.

### Confirmed branch layout

- Leave the current remote `feature/website` and closed PR #45 unchanged.
- Create two clean branches from `origin/main`: `feature/website-improvements`
  and `feature/mcp-providers`.
- Use the current `feature/website` commit and the preserved working-tree
  changes as the reference source for the selected file manifest.
- Put the 38 website source/configuration/content paths and Pages workflow on
  `feature/website-improvements`; omit generated `.svelte-kit/` and `build/`
  output.
- Update the workflow's push filter and preview condition for
  `feature/website-improvements`.
- Put MCP/provider code, source-generated provider/distribution outputs, tests,
  current code-side working-tree edits, and only the selected HAWP records on
  `feature/mcp-providers`.
- Reconstruct the shared backlog delta from `main`; do not copy the current
  47-path mixed work-record diff.
- Create the new PRs targeting `main`, with separate summaries.

Both branches were created locally from `origin/main` at
`178fde6ebcb57f3e8100710c72ae71f06e62fb6d`. The website worktree contains
the selected website paths and the branch-specific workflow update. The
Librarian worktree contains the selected provider/MCP paths and the current
code-side working-tree changes, including the minimal selected backlog
delta. The approved search-index Option A is implemented there; the existing
`feature/website` branch remains the reference. At the initial branch-
preparation checkpoint, no changes had been staged, committed, pushed, or
opened as PRs; that state was superseded by the user's later commit/push
authorization and the verified publication below.

**Implementation coordination:** The approved search/index item `09ee9986`
is implemented and tested on `feature/mcp-providers`; the persistent local
index was not modified. The implementation is included in the open PR #46;
PR #47 contains the separate website/Pages changes.

### Branch preparation result (2026-10-22)

**Directly verified:**

- `feature/website-improvements` and `feature/mcp-providers` both point to
  `178fde6ebcb57f3e8100710c72ae71f06e62fb6d`, matching `origin/main`.
  The `feature/website` branch ref remains at
  `e2237b480c01e3707d618685e4f8e6f913d5a263`; its existing dirty worktree
  remains in place.
- The website worktree has exactly 39 manifest paths (38 website source and
  configuration paths plus `.github/workflows/website.yml`). At initial
  branch preparation, the code worktree had 68 manifest paths, including the
  selected shared backlog delta. The later search Option A addition brought
  the current code-branch manifest to 76 paths. The path sets remain
  disjoint.
- `website/.svelte-kit/` and `website/build/` were absent before the build.
  The build created them as ignored outputs; `.gitignore` excludes them and
  they are not in the 39-path worktree inventory.
- Website `npm run check` passed with 0 errors and 0 warnings; `npm run build`
  passed. `npm ci` was needed because dependencies were not installed. It
  reported a pending install-script warning for `fsevents`; no script was
  approved.
- On `feature/mcp-providers`, `go test ./...` passed. `hawp providers sync`
  updated 0/15 files and validated all materialized files; distribution sync
  updated 0/40 files and validated generated outputs. `hawp check` passed all
  three validations, with 0 work issues/warnings and valid links across 147
  Markdown files.
- The read-only `hawp/hawp_work_validate` call against the original workspace
  (not the code worktree) still reports 2 issues and 1 warning; its link check
  identifies the existing missing active targets
  `active/2eea565c/plan.md` and `active/429e075e/plan.md`. The code-worktree
  CLI check passed independently.
- At the branch-preparation checkpoint, no files were staged or committed;
  neither branch had been pushed and no PR had been created.

**Next:** PR #47 contains the website/Pages changes and PR #46 contains the
Librarian/MCP/provider changes; both target `main`. The approved `09ee9986`
Option A implementation is included in PR #46. Keep the unverified live
client checks and VS Code migration decision separate while the PR checks
and review proceed.

### Search-index Option A validation (2026-10-22)

- The four focused search-index Go packages and the full `go test ./...`
  suite passed on `feature/mcp-providers`.
- `go run ./cmd/hawp distribution sync` updated 0/40 files.
- `go run ./cmd/hawp check` passed all 3 validations with 0 issues/warnings
  across 147 Markdown files.
- The separate read-only `hawp/hawp_work_validate` call against the original
  checkout still reports 2 work issues, 1 warning, and the known links to
  `active/2eea565c/plan.md` and `active/429e075e/plan.md`.
- No local index rebuild was performed.

### Commit and push result (2026-10-22)

- `feature/website-improvements` was committed as
  `f27855b08f0c11d388a22f1700bc124f751a672c` and pushed to
  `origin/feature/website-improvements`.
- `feature/mcp-providers` was committed as
  `a35dd925a04c397e0a8ce363b3066e7a1edc085b` and pushed to
  `origin/feature/mcp-providers`.
- Both remote refs matched their local commits after push. The follow-up
  work-record commit `dfb9f6f13c816d8472f59c808ea4bd7730bf58c8` is also
  pushed on `feature/mcp-providers`.
- The preserved source checkout is on `feature/website` at local commit
  `6e6fbacc08830a1915a92d9259ab4a55562f7f79`, one commit ahead of
  `origin/feature/website` at
  `e2237b480c01e3707d618685e4f8e6f913d5a263`; its worktree is clean. It was
  not changed or pushed during this PR operation.
- PR #47, `Add the HAWP website and Pages deployment`, is open from
  `feature/website-improvements` to `main`:
  https://github.com/sentzunhat/human-ai-workflow-protocol/pull/47
- PR #46, `Expand HAWP MCP tooling and provider setup`, is open from
  `feature/mcp-providers` to `main`:
  https://github.com/sentzunhat/human-ai-workflow-protocol/pull/46
- At PR creation, GitHub reported both PRs mergeable. Website build checks
  passed and deploy was skipped as expected; the code PR's `quality` and
  `validate-generated` checks were pending.
- After updating these records, `cd librarian/src && go run ./cmd/hawp check`
  passed all 3 validations with 0 issues, 0 warnings, and 147 valid Markdown
  files.
- The MCP tool `hawp/hawp_work_validate` ran against the original checkout;
  its check reported 2 work issues and 1 warning, with broken backlog links
  to `active/2eea565c/plan.md` and `active/429e075e/plan.md`. The separate
  code-branch validation passed.

### Opened PR summaries

- Website PR #47: `Add the HAWP website and Pages deployment`.
- Librarian/MCP PR #46: `Expand HAWP MCP tooling and provider setup`.

Both descriptions categorize their branch changes and list the validation
commands already run on the corresponding worktree. The website description
distinguishes preview artifacts from deployment on `main`; the code
description includes the implemented search-index sync and regression tests.

### Verification sequence after branch selection

1. [x] Preserve `feature/website` and create isolated worktrees from
   `origin/main`; verify each branch starts at the recorded main SHA.
2. [x] Apply only the paths listed as owned in `files.md`; compare the exact
   worktree and committed path sets and confirm no overlap. The 39 website
   paths and 76 code paths match the manifest.
3. [x] Run the website checks/build after confirming generated output was
   absent; the resulting build output remains ignored.
4. [x] Run Go tests, provider/distribution sync, and HAWP validation on the
   code branch. Record the separate original-workspace link issues above.
5. [x] Prepare the isolated branch worktrees before implementing search
   Option A.
6. Finalize or defer provider setup claims based on the independent live
   client checks above; do not change ignored machine-local config files.

## Next Step

- [x] Investigate current branch, PR, diff, and existing work records
- [x] Write the split plan and identify the branch-history risk
- [x] Inventory file paths and record non-overlapping ownership
- [x] Create and link this work item
- [x] Confirm branch names and source reference with the user
- [x] Create both isolated worktrees from `origin/main`
- [x] Reconstruct the selected code-branch backlog and verify both manifests
- [x] Implement approved search-index Option A on the Librarian branch
- [x] Create and verify two PRs targeting `main`; PR #47 and PR #46 are open
  and mergeable. Website checks passed. PR #46's `quality` and
  `validate-generated` checks passed on head `b0f048f2` before this
  status-only work-record update.
