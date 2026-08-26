# Hawp-First Session Workflow

Use `hawp_search` (MCP) as the default context strategy for kit and work files. Direct reads
load whole files and surface irrelevant content; `hawp_search` returns ranked chunks sized to
your token budget.

---

## Why hawp-first

- `hawp_search` ranks chunks by relevance and trims to `--max-tokens N`; a direct read of a
  large kit file can cost 3-10x more tokens with most of that content unused.
- MCP exposes the same index to any agent without requiring file-system access.
- Ranked results separate signal from noise; direct reads require you to skim the noise yourself.

---

## Session-start pattern

Run once per session if kit or work files changed since the last index:

```bash
hawp search index
hawp search embed --backend ollama   # or --backend onnx for offline
```

Skip if the index is current and no kit/work files changed this session.

---

## During-work pattern

Before reading any `.hawp/kit/` or `.hawp/work/` file, search first:

```
hawp_search  query="<topic>"  context=true  max_tokens=2000
```

Use the returned chunks to answer the question. Fall through to a direct read only if the
chunks are insufficient (see below).

Adjust `max_tokens` to your remaining budget. For quick lookups, 500-1000 is usually enough.
For broad context gathering, use 2000-4000.

---

## When to fall through to direct reads

Direct reads are appropriate when:

- The target is **implementation code** (Go source, TypeScript scripts, CI configs) — not kit
  or work docs.
- You already know the **exact file and section** and need the full text verbatim.
- The search index is stale or unavailable and re-indexing is not practical mid-session.
- A single small file (under ~100 lines) where reading costs less than formulating a query.

When falling through, prefer `Read` with `offset` and `limit` to read only the relevant section.

---

## Re-index trigger

Re-run `hawp search index && hawp search embed --backend ollama` any time:

- Kit files were edited this session.
- Work files (backlog, active plans, evidence) were added or modified.
- A new status report or evidence file was committed.

The index is fast (seconds for typical repo sizes); err on the side of re-indexing.

---

## Token budget and session continuity

When approaching the context window limit mid-session:

1. **Batch all remaining parallel work** — spawn agents or write commits for everything in flight before the window closes. Agents run in background and survive the context window reset.
2. **Set a reminder** — use `/loop` or note the session state explicitly so the next session can pick up from a known checkpoint.
3. **Commit early, commit often** — work on the manager branch is the handoff artifact. The next session reads `git log` and BACKLOG to reconstruct state.
4. **Save lessons to memory** — anything non-obvious about the current task goes into a memory file *before* the window fills.

The session transcript is the gold standard, but BACKLOG.md + memory files are the fastest cold-start for a new context.

---

## Parallel agent worktrees

When spawning agents with `isolation: "worktree"`, the worktree is auto-removed only
if the agent makes **no changes**. When agents commit work, the worktree persists and
must be cleaned up manually after the branch is merged:

```bash
git worktree remove .claude/worktrees/agent-<id> --force
```

Pattern for managing parallel sub-branches:

1. Create sub-branches from the manager branch (`feature/vX.Y.Z`).
2. Spawn agents with `isolation: "worktree"` — one per sub-branch.
3. On agent completion: squash-merge the sub-branch into the manager branch, then remove the worktree.
4. Do not merge sub-branches into `development` or `main` — only into the manager branch.

List live worktrees at any time: `git worktree list`

For teams that also separate HAWP coordination from a product integration branch,
see [manager-branch.md](manager-branch.md) for the optional manager-branch pattern.

After squash-merging all sub-branches, delete the stale local and remote tracking branches:

```bash
# Delete stale local branches (not main, not the current feature branch)
git branch -d feature/v008-some-subtask   # or -D if not fully merged
# Delete corresponding remote branches
git push origin --delete feature/v008-some-subtask
```

---

## Lessons learned

### hawp init: provision failure blocked provider config writes (v0.0.8)

**What broke:** `hawp init --provider codex` exited 1 before writing `codex.toml` whenever
any asset download failed — including BGE model files whose SHA-256 checksums had been
placeholder values since v0.0.2 and could never verify correctly.

**Root cause:** two independent bugs compounding:
1. `runInit` returned `ExitError{Code:1}` immediately on any provision step failure, before
   kit sync or `WriteProviderConfigs` ran. Provision and provider config are independent;
   the early exit was wrong.
2. BGE-base-en-v1.5 SHA-256s were fake hex strings (`// v0.0.2: verify from HuggingFace`)
   that had never been replaced with real values. Every init attempt failed on BGE.

**Fix:** decouple provision failure from provider config write — both steps always run;
`ExitError{Code:1}` deferred to the end. Remove BGE from `ModelAssets` (moved to
`BGEModelAssets` with empty SHA256 and TODO) so unverified assets can't silently block users.

**Pattern to follow:** when a command has multiple independent phases (asset download,
config write, kit sync), each phase must run to completion regardless of prior phase
failures. Carry failures forward; report them at the end. Never gate unrelated work on
optional steps.

**Verification:** if `hawp init --provider <name>` exits 1 but the config file
(`codex.toml`, `.mcp.json`, etc.) IS written, the command succeeded at its core job.
Asset download failures are a separate concern and require separate remediation.

---

## Co-authoring conventions

When an AI agent authors or co-authors a commit, include a `Co-Authored-By` trailer so
attribution is visible in git history and on GitHub. Use the format:

```
Co-Authored-By: <Model Name> (<Company>) <noreply@company.com>
```

Common agents:

| Agent | Trailer |
|-------|---------|
| Claude (Anthropic) | `Co-Authored-By: Claude Sonnet 4.6 (Anthropic) <noreply@anthropic.com>` |
| GPT-4o (OpenAI) | `Co-Authored-By: GPT-4o (OpenAI) <noreply@openai.com>` |
| o3 / Codex (OpenAI) | `Co-Authored-By: Codex (OpenAI) <noreply@openai.com>` |
| Gemini (Google) | `Co-Authored-By: Gemini 2.5 Pro (Google) <noreply@google.com>` |
| Cursor (Anysphere) | `Co-Authored-By: Cursor (Anysphere) <noreply@cursor.sh>` |
| Continue | `Co-Authored-By: Continue (Continue) <noreply@continue.dev>` |

Update the model name (e.g. `Sonnet 4.6` → `Opus 5`) when the version changes; the company
and email stay stable per vendor.

For human + AI pair sessions, include both authors:

```
Co-Authored-By: Diego Beltran <beltrd@gmail.com>
Co-Authored-By: Claude Sonnet 4.6 (Anthropic) <noreply@anthropic.com>
```

---

## Quick reference

| Situation | Action |
|-----------|--------|
| Session start, files changed | `hawp search index && hawp search embed --backend ollama` |
| Need context on a topic | `hawp_search` via MCP, `max_tokens` to budget |
| Broad context pull | `hawp_search context=true max_tokens=3000` |
| Implementation code | Direct `Read` — search index does not cover non-kit files |
| Exact known section | Direct `Read` with `offset`/`limit` |
| Index stale mid-session | Re-run index + embed, then search |
