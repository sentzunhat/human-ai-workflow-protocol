# f8685e56: Split the website and Librarian changes into separate main-targeted PRs — File Tracking

**Purpose:** Keep the two PRs file-disjoint, omit generated website output, and
preserve unrelated work on the existing source branch.

**Work Item:** `.hawp/work/active/f8685e56/plan.md`

**Last Updated:** 2026-10-22

## Boundary and Branches

The user confirmed `feature/website-improvements` and
`feature/mcp-providers`, both based on `origin/main`. The existing
`feature/website` branch ref remains the source/reference and stays at
`e2237b480c01e3707d618685e4f8e6f913d5a263`; its pre-existing dirty worktree
is preserved and contains the updated split/client-gate records.
Both local branches and their worktrees now exist at
`178fde6ebcb57f3e8100710c72ae71f06e62fb6d`, the recorded `origin/main` tip.
The selected website and provider/MCP paths have been populated into their
respective worktrees. The original reference branch remains unchanged.

Do not move files between canonical repository directories. Split by branch
ownership only. No path may appear in both PRs.

## Preparation Verification (2026-10-22)

- The website worktree contains exactly 39 paths from this manifest: 38
  website source/configuration paths and the Pages workflow.
- The code worktree contains 76 paths from this manifest: the 68 original
  split paths plus the 8 approved search-index Option A paths, including the
  minimal shared-backlog delta and selected work records.
- The two path sets have no overlap. Both branch heads match the recorded
  `origin/main` SHA.
- The website type check and build passed. `.svelte-kit/` and `build/` were
  absent before the build and are ignored generated output afterward. `npm ci`
  reported a pending `fsevents` install-script warning; no script was approved.
- The code branch passed `go test ./...`, provider sync, distribution sync,
  and `hawp check`. The original workspace MCP check still reports its two
  existing missing-plan links and one warning; see the split plan.
- At the initial branch-preparation checkpoint, the worktrees contained only
  uncommitted changes. No staging, commit, push, or PR creation had been
  performed at that point.

## Owned Files

### Website PR — branch `feature/website-improvements`

The website PR owns these 38 non-generated website paths and the Pages
workflow:

```txt
.github/workflows/website.yml
website/.gitignore
website/.nvmrc
website/README.md
website/docs/discoverability-baseline.md
website/docs/post-merge-discoverability-actions.md
website/package-lock.json
website/package.json
website/src/app.css
website/src/app.html
website/src/lib/components/Agents.svelte
website/src/lib/components/AnalyticsConsent.svelte
website/src/lib/components/Benchmarks.svelte
website/src/lib/components/Faq.svelte
website/src/lib/components/FinalCta.svelte
website/src/lib/components/Footer.svelte
website/src/lib/components/Header.svelte
website/src/lib/components/Hero.svelte
website/src/lib/components/Problem.svelte
website/src/lib/components/Protocol.svelte
website/src/lib/components/ReleaseDownloads.svelte
website/src/lib/components/ReleaseVersion.svelte
website/src/lib/components/ThemeToggle.svelte
website/src/lib/components/Workflow.svelte
website/src/lib/content/site.ts
website/src/lib/theme.ts
website/src/routes/+layout.svelte
website/src/routes/+layout.ts
website/src/routes/+page.svelte
website/static/CNAME
website/static/favicon.svg
website/static/llms.txt
website/static/robots.txt
website/static/sitemap.xml
website/static/social-card.png
website/static/social-card.svg
website/svelte.config.js
website/tsconfig.json
website/vite.config.ts
```

`.github/workflows/website.yml` must include
`feature/website-improvements` in the push filter and update the preview
condition from `refs/heads/feature/website` to
`refs/heads/feature/website-improvements`.
Do not make this workflow change on the Librarian branch.

### Librarian/MCP PR — branch `feature/mcp-providers`

These 56 committed non-website paths belong to the code/provider PR:

```txt
.github/agents/hawp.agent.md
.github/copilot-instructions.md
.github/skills/commit-many-small/SKILL.md
.github/skills/commit-one-big/SKILL.md
.github/skills/conservative-docs-drift-safe-cleanup/SKILL.md
.github/skills/docs-alignment-deterministic/SKILL.md
.github/skills/docs-alignment-simplicity-pass/SKILL.md
.github/skills/hawp-backlog-alignment/SKILL.md
.github/skills/hawp-change-review-and-reference-sync/SKILL.md
.github/skills/intake/SKILL.md
.github/skills/intent-first-handoff/SKILL.md
.github/skills/status-report/SKILL.md
.gitignore
.hawp/config/search.json
.hawp/kit/usage/mcp/README.md
.hawp/kit/usage/mcp/github-copilot.md
.hawp/kit/usage/mcp/lessons.md
.hawp/kit/usage/mcp/verification.md
core/providers/.claude/rules/hawp-core.md
core/providers/.codex/AGENTS.md.seed
core/providers/.continue/rules/hawp-01-core.md
core/providers/.cursor/rules/hawp-core.mdc
core/providers/.github/copilot-instructions.md
core/providers/doc-verification.md
core/providers/shared/behaviors/hawp-core.md
distribution/generated/claude/install/development.md
distribution/generated/claude/install/main.md
distribution/generated/claude/update/development.md
distribution/generated/claude/update/main.md
distribution/generated/codex/install/development.md
distribution/generated/codex/install/main.md
distribution/generated/codex/update/development.md
distribution/generated/codex/update/main.md
distribution/generated/continue/install/development.md
distribution/generated/continue/install/main.md
distribution/generated/continue/update/development.md
distribution/generated/continue/update/main.md
distribution/generated/cursor/install/development.md
distribution/generated/cursor/install/main.md
distribution/generated/cursor/update/development.md
distribution/generated/cursor/update/main.md
distribution/generated/github/install/development.md
distribution/generated/github/install/main.md
distribution/generated/github/update/development.md
distribution/generated/github/update/main.md
distribution/sources/shared/install.md
distribution/sources/shared/update.md
librarian/src/internal/platform/mcp/configure/config.go
librarian/src/internal/platform/mcp/configure/config_agent.go
librarian/src/internal/platform/mcp/configure/configure.go
librarian/src/internal/platform/mcp/configure/configure_test.go
librarian/src/internal/platform/mcp/server/prompts.go
librarian/src/internal/platform/mcp/server/prompts_test.go
librarian/src/internal/platform/mcp/server/resources.go
librarian/src/internal/platform/mcp/server/resources_test.go
librarian/src/internal/platform/mcp/server/server.go
```

The following current worktree changes also belong to the code PR:

```txt
.github/agents/docs-alignment.agent.md
.github/agents/master-builder.agent.md
core/.hawp/kit/usage/mcp/README.md
core/.hawp/kit/usage/mcp/github-copilot.md
core/providers/.github/prompts/hawp-change-review-and-reference-sync.prompt.md
.hawp/work/active/09ee9986/plan.md
.hawp/work/active/b152c56c/plan.md
.hawp/work/active/f8685e56/plan.md
.hawp/work/active/f8685e56/files.md
.hawp/work/evidence/2026/10/22/09ee9986/evidence.md
.hawp/work/status/2026/10/21/b152c56c/status.md
librarian/src/internal/platform/mcp/configure/config.go
librarian/src/internal/platform/mcp/configure/config_agent.go
librarian/src/internal/platform/mcp/configure/configure_test.go
.hawp/kit/usage/search.md
librarian/src/internal/application/index/ingest-service.go
librarian/src/internal/application/index/ingest-service_test.go
librarian/src/internal/infrastructure/repositories/index/index.go
librarian/src/internal/infrastructure/repositories/index/fts_sync_test.go
librarian/src/internal/platform/cli/index/ingest/command.go
librarian/src/internal/platform/cli/index/ingest/corpus_test.go
librarian/src/internal/platform/mcp/server/tools_e2e_test.go
```

Other modified paths in the current worktree are already listed above and
remain with the same code-PR owner.

`.hawp/work/BACKLOG.md` has one owner: the Librarian PR. Reconstruct its
minimal delta from `origin/main` for records actually included in that PR;
do not carry over the full mixed 47-path work-record delta. Do not duplicate
the backlog file on the website branch.

Provider documentation/agent changes remain subject to live client evidence.
The user's confirmation that the fresh Copilot chat already selected the HAWP
agent is recorded; do not ask them to select it again unless the session
changes. The selected agent's actual Tools list, client approval behavior,
and the VS Code migration-preview effect remain unverified. Do not describe
repository tool discovery as proof of those client behaviors. Keep
`.hawp/work/active/cb096b6a/plan.md` and
`.hawp/work/parked/70d8112f/plan.md` on the preserved source branch until their
separate migration-preview and installer-audit gates are addressed.

The tool-discovery query for client UI/session inspection surfaced HAWP MCP
tools only; no VS Code UI control was available in this conversation. If that
remains true, request the live Tools list, approval behavior, and preview
text/screenshot. Do not ask the user to select the HAWP agent again unless
their session changes.

### Search-index Option A addition (2026-10-22)

Option A is implemented on `feature/mcp-providers`. The eight paths are now
included in the code-PR worktree manifest above; the original website path
set is unchanged, so the two PRs remain file-disjoint. Focused tests and the
full Go suite passed. Provider/distribution and HAWP code-branch validation
passed after work-record sync. The original-checkout MCP validation still
reports the previously recorded missing-plan issues. The persistent local
index remains untouched, and PR creation is still pending.

`distribution/generated/` is derived output: edit
`distribution/sources/shared/install.md` and
`distribution/sources/shared/update.md`, then run the repository's
distribution sync. Provider overlays are derived from
`core/providers/shared/behaviors/hawp-core.md`; synchronize them from that
source rather than hand-editing generated copies.

## Read-Only Context Files

```txt
.hawp/kit/references/work-item-file-tracking.md
.hawp/kit/templates/work-item-files.md
.github/workflows/website.yml
website/.gitignore
website/package.json
```

## Do-Not-Touch / Keep on Existing Source Branch

These paths are outside the two product PRs for now. Preserve them on the
existing `feature/website` branch; do not delete, move, or silently add them
to either PR:

```txt
benchmark/runs/2026-06-11-bounded-repo-review/no-hawp/output.md
scripts/artwork/generate-hawp-featured-card.mjs
.hawp/work/active/098c8a5c/plan.md
.hawp/work/active/098c8a5c/hawp-featured-card.svg
```

The following tracked generated output is also excluded from the website PR.
The count of 86 paths is confirmed by the current branch comparison:

```txt
website/.svelte-kit/**
website/build/**
```

Hold all other `.hawp/work/` changes from the mixed branch until their owning
work items and backlog links are reviewed individually. This includes
unrelated closed history, parked website work, and records not selected above.

## Locked / Reserved Files

- `.github/workflows/website.yml` — website PR only.
- `website/` source paths listed above — website PR only.
- `.hawp/work/BACKLOG.md` — code PR only; minimal selected delta.
- `distribution/sources/shared/install.md` and
  `distribution/sources/shared/update.md` — code PR source of truth; generated
  distribution copies stay synchronized from these sources.
- `core/providers/shared/behaviors/hawp-core.md` — code PR source of truth for
  provider overlays.

## Pre-commit Status

The manifest and plan guided creation of the two local worktrees and
population of their selected paths. The approved search-index Option A is
implemented on the code branch and adds the eight paths listed above. The
user has authorized committing and pushing these two new branches. The
existing `feature/website` reference checkout remains out of scope; its
unrelated dirty files are preserved. No files were moved within that
checkout.

## Verification Notes

Planning checks already run:

- `git diff --check` passed.
- `hawp/hawp_work_validate` passed kit checks (3 checks, 0 issues), but work
  validation reported 2 issues and 1 warning. The link check found that
  `.hawp/work/BACKLOG.md` points to missing
  `.hawp/work/active/2eea565c/plan.md` and
  `.hawp/work/active/429e075e/plan.md`. These existing record/link issues are
  outside the split scope.
- `hawp/hawp_search` returned no indexed results for the file-ownership query;
  repository files and Git diffs supplied the direct evidence.
- The four focused search-index Go packages passed, and the full
  `go test ./...` suite passed on `feature/mcp-providers`. The persistent
  local index was not rebuilt; a fresh live MCP query against that index
  remains pending.
- After the work-record sync, `go run ./cmd/hawp distribution sync` updated
  0/40 files and `go run ./cmd/hawp check` passed all 3 validations with
  0 issues/warnings across 147 Markdown files. The separate read-only
  `hawp/hawp_work_validate` call against the original checkout still reported
  the two known missing active-plan links and one warning.

Before creating PRs:

```bash
git -C <website-worktree> status --short --branch
git -C <website-worktree> diff --name-status origin/main...HEAD
git -C <website-worktree> diff --check
git -C <code-worktree> status --short --branch
git -C <code-worktree> diff --name-status origin/main...HEAD
git -C <code-worktree> diff --check
```

Verify the website output is regenerated from source with the generated
directories absent, run the affected Go tests, synchronize providers and
distribution, then run HAWP work and link validation. Do not stage or commit
until the exact path manifest is checked against both branch diffs.
