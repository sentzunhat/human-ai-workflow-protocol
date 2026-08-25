# v0.0.8 Session Close Snapshot — 2026-08-24 (goodnight)

**Manager branch:** `feature/v0.0.8` — 12 commits ahead of `main`. **Do NOT merge to `development` or `main`.**

---

## Commits on feature/v0.0.8 (not yet in main)

```
e37ff98  chore: close 0ca7cf49 in BACKLOG
c07fd65  chore: 0ca7cf49 — stale llm-reshape doc traces removed
787ae7b  chore: 8672216a — backlog audit + Recently Closed compacted 15→10
5a5d827  chore: CHANGELOG v0.0.8 full; close 4c88f451; token-budget lesson in workflow doc
b0991b8  feat: 4c88f451 — Jaccard dedup + dynamic chunk cap for --context (~30% token savings)
9a51a42  chore: add work items 1c743447, 0ca7cf49, 8672216a to active backlog
dec41c8  chore: worktree-cleanup lesson + close c98518bb
f75a93f  feat: c98518bb — hawp-first session workflow kit doc
12bcba8  chore: open active plans c98518bb and 4c88f451
9617e5d  docs: v0.0.8 cleanup — remove llm-reshape docs, fix CHANGELOG, add work items
b489a22  chore: remove --llm-reshape from CLI; engine key canonical
12fdb41  fix: v0.0.8 — engine as canonical JSON key for context config
```

---

## Agents still running (check on wake)

| Agent worktree | Branch | Work item | Status |
|---|---|---|---|
| `agent-a7b4cede444e43295` | `feature/v008-benchmark` | benchmarks-v008 (ONNX vs Ollama, context, hybrid ratios) | may have completed |
| `agent-abfc239b8df2f01ef` | `feature/v008-hybrid-ratio` | `1c743447` — `--hybrid-ratio` flag | may have completed |

**On wake-up, for each completed agent:**
1. `git fetch origin feature/v008-<name>`
2. `git merge --squash origin/feature/v008-<name>`
3. `git commit` + `git push origin feature/v0.0.8`
4. `git worktree remove .claude/worktrees/agent-<id> --force`
5. Close work item in BACKLOG (move plan to `closed/2026/08/24/<uuid>/`)

---

## Next compoundable work after agents land

1. **Rerun benchmarks** — once `--hybrid-ratio` flag is merged, run `hawp search benchmark` with `--hybrid-ratio 0.5` and `--hybrid-ratio 0.7` to generate real comparison data. Update `benchmarks-v008.md` or add a v008-hybrid section.
2. **Update `search.md`** — document `--hybrid-ratio` flag in the Full flag reference and search modes table.
3. **`hawp kit validate` + `npm --prefix librarian run validate`** — full validation suite before shipping v0.0.8.
4. **v0.0.8 ship prep** — once all above is clean, cut PRs: `feature/v0.0.8` → `development` → `main`. Release pipeline fires automatically.

---

## Active BACKLOG items remaining

| UUID | Title | Branch |
|---|---|---|
| `1c743447` | `--hybrid-ratio` flag | `feature/v008-hybrid-ratio` (agent) |

---

## Key v0.0.8 changes (for CHANGELOG review)

- `engine` key canonical (was `backend`) in context config — breaking for old `backend` key configs
- `--llm-reshape` removed from CLI; ContextReshaper/RAG pipeline still in Go codebase
- Real `ContentJaccardDedup` (~30% token savings); old cosine dedup was silent no-op
- Dynamic chunk cap for `--context`; `--verbose` prints token savings to stderr
- `--hybrid-ratio <f>` flag (if/when merged from agent)
- New kit doc: `.hawp/kit/usage/hawp-first-workflow.md`
- Benchmarks-v008: ONNX vs Ollama comparison (when benchmark agent lands)
- Code cleanup: 3 stale doc comments removed; all tests green

---

## Lessons learned this session

- Worktree cleanup is manual after agents commit (`git worktree remove --force`)
- Approaching token limit: batch all parallel work NOW, commit to manager branch, set wakeup
- Old `DeduplicateResults` was a silent no-op — embeddings were always `[]float32{}`
- `v010-3-3a` parked reason was wrong — FLAN-T5-small IS feasible per its plan file
- `feature/vX.Y.Z/sub-branch` naming fails when parent branch exists — use flat names (`feature/v008-*`)

All lessons documented in memory files and `.hawp/kit/usage/hawp-first-workflow.md`.
