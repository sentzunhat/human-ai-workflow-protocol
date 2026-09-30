# Parallel Agent Worktrees

How to run multiple AI agents (Claude Code, Codex, Cursor, GitHub Copilot,
Continue, or any editor with Git worktree support) on independent slices of
the same feature branch without merge conflicts.

All agent workspace directories live under `.hawp/.spaces/` — the project's
workflow root. This keeps parallel lanes discoverable and provider-agnostic.

---

## The Problem With `isolation: "worktree"`

The Agent tool's built-in `isolation: "worktree"` creates a worktree from the
repo's **default branch** (typically `main`), not from the current feature
branch. If you are mid-way through a feature branch with several commits, the
agent works against stale source and the squash-merge produces conflicts.

**Symptom:** the second commit in the agent's `git log` is `main`'s latest
commit rather than the feature branch HEAD.

---

## The Manager-Branch Pattern

Create worktrees manually from the feature branch, then spawn agents without
isolation that work in those paths.

### Step 1 — create sub-branches from the feature branch

```bash
# Create sub-branches from the feature branch

git worktree add .hawp/.spaces/agent-<slug> -b agent/<slug>
```

Both sub-branches share the same base commit as the feature branch HEAD.

### Step 2 — spawn agents pointing at the worktree paths

In the agent prompt, set the working directory to the worktree path:

```
Working directory: <repo-root>/.hawp/.spaces/agent-<slug>
All Go commands run from: <repo-root>/.hawp/.spaces/agent-<slug>/librarian/src
```

Do NOT use `isolation: "worktree"` — the worktree already exists.

### Step 3 — verify base before merging

Before squash-merging any agent branch, check it shares the right base:

```bash
git log --oneline agent/<slug> | head -3
# Second line must match feature branch HEAD, not main's latest
```

If the second line is main's latest commit: re-apply the agent's diff manually
to the correct files rather than squash-merging the branch.

### Step 4 — squash-merge in dependency order, clean up

```bash
# Merge the dependency first
git merge --squash agent/<slug-1>
git commit -m "..."

# Then the dependent
git merge --squash agent/<slug-2>
git commit -m "..."

# Remove worktrees and branches — do this immediately after each merge
git worktree remove .hawp/.spaces/agent-<slug-1> --force
git worktree remove .hawp/.spaces/agent-<slug-2> --force
git branch -D agent/<slug-1> agent/<slug-2>
```

**Cleanup is mandatory after every squash-merge.** Leaving stale worktrees
and agent branches causes three problems:

1. `git worktree list` grows cluttered — hard to tell what's in-flight vs. done.
2. GUI clients (GitKraken, etc.) may not auto-refresh and will show stale branch
   topology until a manual refresh.
3. If the next parallel round creates a branch with the same slug, Git refuses
   to create it because the old branch still exists.

**Lesson learned (2026-09-10):** even when worktrees disappear from the GUI
before you remove them (because they complete fast and the GUI doesn't poll),
always run the `git worktree remove --force` + `git branch -D` cleanup
explicitly. Fast-completing agents still leave branches behind.

### Step 5 — verify clean state before the next round

```bash
git worktree list          # should show only the main repo
git branch -a | grep agent # should return nothing
git status                 # should be clean
```

Only create new worktrees once all prior agent branches are gone.

---

## File Ownership Rules

Assign non-overlapping files before spawning. Conflicts are always avoidable:

| Agent                      | Owns                                                                                                |
| -------------------------- | --------------------------------------------------------------------------------------------------- |
| Domain/infra boundary work | `internal/domain/work/`, `internal/infrastructure/repositories/work/`, `internal/application/work/` |
| CLI / platform work        | `internal/platform/cli/`, `internal/platform/mcp/`                                                  |

Never assign the same file to two agents. The `.hawp/work/active/<uuid>/` plan
folder for each item is also exclusive — one agent per plan file.

---

## Quick Reference

```bash
# Create two worktrees from current feature branch HEAD
git worktree add .hawp/.spaces/agent-A -b agent/A
git worktree add .hawp/.spaces/agent-B -b agent/B

# Verify both share the right base
git log --oneline agent/A | sed -n '2p'  # must equal feature HEAD
git log --oneline agent/B | sed -n '2p'  # must equal feature HEAD

# After agents complete — merge A first (if B depends on A), then B
git merge --squash agent/A && git commit -m "..."
git merge --squash agent/B && git commit -m "..."

# Cleanup
git worktree remove .hawp/.spaces/agent-A --force
git worktree remove .hawp/.spaces/agent-B --force
git branch -D agent/A agent/B
```
