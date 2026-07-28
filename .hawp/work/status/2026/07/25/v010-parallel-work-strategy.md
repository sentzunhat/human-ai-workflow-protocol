# v0.1.0 Parallel Work Strategy — 2026-07-25

## Overview

v0.0.3 is ship-ready. Five parallel work tracks are defined, scoped, and ready to start immediately after v0.0.3 ships (2026-07-30).

---

## 5 Parallel Tracks

### Track 1: OpenAI Embeddings (Phase 3.2c) — 6 hours
- **Lead:** Assign to strongest developer (most impactful)
- **Blockers:** None (independent)
- **Dependencies:** OpenAI SDK
- **Files:** `openai_embedder.go`, tests, docs
- **Detail plan:** `.hawp/work/active/v010-phase3-2c-openai-embeddings.md`

**Start date:** 2026-07-30
**Est. complete:** 2026-07-31 (end of day)

### Track 2: OpenAI LLM (Phase 3.3c) — 8 hours
- **Lead:** Strong developer + token counting knowledge
- **Blockers:** None (independent)
- **Dependencies:** OpenAI SDK, tiktoken
- **Files:** `openai_client.go`, `token_counter.go`, tests, docs
- **Detail plan:** `.hawp/work/active/v010-phase3-3c-openai-llm.md`

**Start date:** 2026-07-30
**Est. complete:** 2026-08-01

**Note:** Can learn from Track 1 (OpenAI patterns). Shares `token_counter.go` utility with Track 4.

### Track 3: Anthropic Embeddings Stub (Phase 3.2d) — 3 hours
- **Lead:** Junior developer (low complexity)
- **Blockers:** None (independent)
- **Dependencies:** None (stub only)
- **Files:** `anthropic_embedder.go`, tests, docs
- **Detail plan:** `.hawp/work/active/v010-phase3-2d-anthropic-embeddings.md`

**Start date:** 2026-07-31 (can start after 3.2c to learn from pattern)
**Est. complete:** 2026-08-01

### Track 4: Anthropic LLM (Phase 3.3d) — 8 hours
- **Lead:** Developer who did Track 2 (OpenAI), or pair with Token Counting dev
- **Blockers:** None (independent)
- **Dependencies:** Anthropic SDK, token counting (from Track 2)
- **Files:** `anthropic_client.go`, tests, docs
- **Detail plan:** `.hawp/work/active/v010-phase3-3d-anthropic-llm.md`

**Start date:** 2026-08-01 (waits for token_counter from Track 2)
**Est. complete:** 2026-08-03

### Track 5: Cost Tracking + Rate Limiting (Cross-cutting) — 10 hours
- **Lead:** Developer with systems experience (most complex)
- **Blockers:** Waits for cloud backend PRs (need to wrap them)
- **Dependencies:** Completed Track 1, Track 2, Track 4
- **Files:** `cost_tracker.go`, `rate_limiter.go`, integrations, docs
- **Detail plan:** `.hawp/work/active/v010-cost-tracking-rate-limiting.md`

**Start date:** 2026-08-01 (after 3.2c ready to integrate)
**Est. complete:** 2026-08-04

---

## Recommended Developer Assignment (5 devs)

| Track | Dev | Start | Complete |
|---|---|---|---|
| 1. OpenAI Embed | Alice (strongest) | 07-30 | 07-31 |
| 2. OpenAI LLM | Bob | 07-30 | 08-01 |
| 3. Anthropic Embed Stub | Charlie (junior) | 07-31 | 08-01 |
| 4. Anthropic LLM | Bob (continues from track 2) | 08-01 | 08-03 |
| 5. Cost + Rate Limit | Diana (systems) | 08-01 | 08-04 |

**Timeline with 5 devs:** Full completion by 2026-08-04 (6 days, parallel work)

---

## Recommended Timeline (1-2 devs)

If only 1-2 developers available, run sequentially:

**Week 1:**
1. Phase 3.2c (OpenAI Embed) — 6 hours
2. Phase 3.3c (OpenAI LLM) — 8 hours
→ Total: ~14 hours = ~2 days

**Week 2:**
3. Phase 3.2d (Anthropic Stub) — 3 hours
4. Phase 3.3d (Anthropic LLM) — 8 hours
→ Total: ~11 hours = ~1.5 days

**Week 3:**
5. Cost Tracking + Rate Limiting — 10 hours
→ Total: ~10 hours = ~1.5 days

**Total sequential:** ~3 weeks (2026-07-30 to 2026-08-20)

---

## Dependencies & Ordering

```
Start (2026-07-30)
    ├─ Phase 3.2c (OpenAI Embed) ──┐
    │                               ├─ Cost Tracking (integrates both)
    ├─ Phase 3.3c (OpenAI LLM) ────┤
    │   (with token_counter.go)     │
    │                               │
    ├─ Phase 3.2d (Anthropic Stub)  │
    │                               │
    └─ Phase 3.3d (Anthropic LLM) ──┘
       (uses token_counter from 3.3c)
```

**Critical path:**
1. OpenAI Embed (3.2c) — no dependencies
2. OpenAI LLM (3.3c) with token counting — no dependencies
3. Cost tracking integrates 3.2c + 3.3c

**All other phases can overlap without blocking.**

---

## Quality Gates

Before marking each track complete:

| Phase | Gate 1 | Gate 2 | Gate 3 |
|---|---|---|---|
| 3.2c (OAI Embed) | Unit tests pass | Integration tests pass | Docs written |
| 3.3c (OAI LLM) | Unit tests pass | Integration tests pass | Token counting accurate |
| 3.2d (Anthro Stub) | Compiles | Interface defined | Docs written |
| 3.3d (Anthro LLM) | Unit tests pass | Integration tests pass | Cost comparable |
| Cost + Rate Limit | Unit tests pass | Tracked in wrapped clients | Config works |

---

## Success Metrics

**By 2026-08-04 (parallel):**
- ✅ All 4 cloud backend phases complete
- ✅ Cost tracking + rate limiting integrated
- ✅ 40+ new tests passing
- ✅ Full documentation updated
- ✅ Config system extended

**Ship v0.1.0: 2026-08-07**

---

## Communication Plan

### Daily Standup Items
- **Monday 07-30:** All 5 tracks start (async standup)
- **Wed 08-01:** Cost tracking dev starts integrating
- **Fri 08-03:** Anthropic LLM complete, all phases ready for cost integration
- **Mon 08-05:** All phases integrated, start final testing

### PR Review Rotation
- Track 1 → Alice reviews Track 2, 3
- Track 2 → Bob reviews Track 4, 5 (cost tracker parts)
- Track 3 → Charlie reviews Track 1, 4
- Track 4 → Bob reviews Track 5
- Track 5 → Diana reviews Track 1, 2 (integration points)

---

## Rollback Plan

Each track is independent; can rollback without affecting others:

- Rollback 3.2c: Revert `openai_embedder.go`, remove from factory
- Rollback 3.3c: Revert `openai_client.go`, remove from factory
- Rollback 3.2d: Remove stub (no production code affected)
- Rollback 3.3d: Revert `anthropic_client.go`, remove from factory
- Rollback Cost: Unwrap backends (revert integration)

---

## Pre-Ship Checklist (2026-08-05)

- [ ] All 5 tracks PR'd and reviewed
- [ ] All 40+ tests passing on main
- [ ] No flaky tests
- [ ] Docs complete (each backend + cost guide)
- [ ] Changelog updated (4 cloud backends + cost tracking)
- [ ] Version bumped to 0.1.0
- [ ] Performance benchmarked (compare to Ollama/ONNX)

---

## Post-Ship (2026-08-07+)

### Immediate
- v0.1.0 released and tested
- Auto-update propagates to v0.0.3 users
- Monitor for API integration issues

### Short-term (Week 2)
- Analyze usage patterns (which backends are users choosing?)
- Collect cost data (are budgets working?)
- Address any bugs found in cloud backends

### Medium-term (v0.2.0, Sept)
- Batch API integration (cost savings)
- GPU acceleration for local backends
- Agentic loops + multi-turn conversations

---

## Work Item Links

| Phase | Link |
|---|---|
| 3.2c | `.hawp/work/active/v010-phase3-2c-openai-embeddings.md` |
| 3.3c | `.hawp/work/active/v010-phase3-3c-openai-llm.md` |
| 3.2d | `.hawp/work/active/v010-phase3-2d-anthropic-embeddings.md` |
| 3.3d | `.hawp/work/active/v010-phase3-3d-anthropic-llm.md` |
| Cost | `.hawp/work/active/v010-cost-tracking-rate-limiting.md` |

---

**Status: READY TO SHIP v0.0.3 & START v0.1.0 (2026-07-30)**
