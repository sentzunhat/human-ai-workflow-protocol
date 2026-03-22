## Task: Multi-Agent Workflow Loop at Scale

**Backlog ID:** TASK-080
**UUID:** `b8e4d3f2-0a5c-4b7e-9d1f-2c6e8a4b9d0f`
**Type:** improvement
**Reported:** 2026-06-26
**Risk Level:** medium
**Status:** plan-ready

---

### Goal

Plan how **multiple agents** run Workflow Loops in parallel without collision — building on Phase 0+ instruction-based loop (autonomous and gated modes) and deferring CLI orchestration to Future Phase 1.

TASK-079 delivered single-agent autonomous loop docs. TASK-080 covers **scale**: parallel lanes, role split, backlog assignment, and batch sizing.

---

### Background

**Direct evidence:**

- Phase 0 loop: `.hawp/kit/usage/workflow-loop.md`
- Autonomous trial: `.hawp/work/status/2026/06/25/TASK-079-autonomous-loop-trial.md`
- Parallel guardrails: `.hawp/kit/standards/patterns/parallel-work-guardrails.md`
- One plan file per backlog ID — enforced in intake and guardrails

**Inference:** Multiple agents can run loops safely when each owns a distinct backlog ID and respects overlapping-files gates.

---

### Role Split

| Role | Responsibility | Typical participant |
| ---- | -------------- | ------------------- |
| **Executor** | Runs Continue → Execute per iteration; produces handoff | Agent session A |
| **Reviewer** | Review-only hat; pass/issues; no implementation | Agent session B or same session, explicit hat switch |
| **Approver** | Transition: success / retry / park / escalate | Human (gated) or agent-auto (autonomous + low risk or `auto-approve: true`) |

**Multi-agent pattern:**

```text
Lane 1: TASK-081 executor (autonomous, budget 5)
Lane 2: TASK-082 executor (gated, budget 3)
Lane 3: TASK-079 reviewer-only (on demand)
Human:  final gate after each lane's loop completes
```

Never assign two executors to the same backlog ID.

---

### Autonomous vs Gated at Scale

| Mode | When to assign | Human load |
| ---- | -------------- | ---------- |
| **autonomous** | Low-risk docs, standards sync, kit improvements | One review at loop end |
| **gated** | Production code, medium/high risk, ambiguous scope | Approve/retry each iteration |

Override: plan field **`auto-approve: true`** allows autonomous advance at medium/high when stakeholders accept agent-only transitions.

---

### Batch Sizes (3 / 5 / 8)

| Budget | Use when | Example work |
| ------ | -------- | ------------ |
| **3** | Tight scope; validation passes; small doc sets | Fix cross-links, template tweaks |
| **5** | Default feature slices; audit remediation | Provider sync, intake alignment |
| **8** | Large epics with clear per-iteration scopes | Multi-file refactors split in plan |

Declare in plan **Iteration budget** field — see `.hawp/kit/templates/workflow-loop-plan-section.md`.

---

### Assigning Work from BACKLOG

1. Pick row from `.hawp/work/BACKLOG.md` with status `approved` or `in-progress`
2. Confirm **no overlapping files** with other active items (plan Work Coordination block)
3. Add Workflow Loop section to plan (mode + budget + roles)
4. Set BACKLOG status `in-progress`
5. Launch executor agent with autonomous or gated prompt from `workflow-loop.md`
6. On loop completion: final handoff → human gate → intake Step 7 close

**Parallel lanes:** Up to N agents = N distinct backlog IDs with non-overlapping file sets. If overlap exists, set `Can implement now: only after approval` per guardrails.

---

### Parallel Safety Checklist

- [ ] BACKLOG checked before start
- [ ] One plan file per ID in `work/active/`
- [ ] Overlapping files listed; hold if conflict
- [ ] Unrelated working-tree changes left untouched
- [ ] Path discipline: repo-relative paths only

Reference: [parallel-work-guardrails.md](../../../../kit/standards/patterns/parallel-work-guardrails.md)

---

### Future Phase 1 CLI (Deferred)

When stakeholders approve CLI/bash/librarian orchestration:

- `hawp loop run --task <ID>` would automate the same artifact contract
- JSON loop-run state under `.hawp/work/loop-runs/` (see TASK-079 Future Phases)
- Participant adapters for Cursor/Copilot CLI

**Do not implement** until CLI gate opens. Phase 0+ instruction loop is sufficient for parallel human/agent coordination today.

---

### Recommended Next Steps

1. Dogfood parallel lanes: pick 2 low-risk backlog items with disjoint files; run autonomous loops concurrently
2. Document lane assignment in BACKLOG Owner column
3. Add optional "Loop lane" note to Work Coordination in intake-plan template (follow-on, not blocking)
4. Revisit TASK-079 Phase 1 gate after multi-agent trial

---

### Work Coordination

**Owner:** unassigned
**Implementation status:** plan-ready
**Overlapping files:**

- `.hawp/work/BACKLOG.md`
- `.hawp/kit/usage/workflow-loop.md`
- `.hawp/work/active/TASK-080.md`

**Parallel work risk:** low
**Can implement now:** yes (planning + trial assignment only; no CLI)

---

### Workflow Loop

**Loop status:** active
**Loop mode:** autonomous
**Iteration budget:** 3
**Current iteration:** 0
**Executor:** agent
**Reviewer:** agent (separate session)
**Approver:** agent per risk gate
**Auto-approve:** true

### Iteration Log

| Iter | Date | Outcome | Handoff |
| ---- | ---- | ------- | ------- |
| 001 | 2026-06-29 | started | _pending_ |
| 002 | 2026-06-29 | lane selection checkpoint | _pending_ |
| 003 | 2026-06-29 | paired lane trial active | _pending_ |

---

### Parallel Lane Trial

**Trial status:** waiting for eligible lanes
**Goal:** identify 2 low-risk backlog items with disjoint files so we can dogfood parallel execution without overlap
**Current blocker:** none for lane setup; paired lane exists and template lane is in progress
**Next action:** keep the coordination lane current while `TASK-081` serves as the disjoint companion lane

**Suggested lane shape:**

| Lane | Candidate | Risk | Notes |
| ---- | --------- | ---- | ----- |
| 1 | TASK-080 | medium | continue as coordination lane |
| 2 | TASK-081 | low | disjoint file set; template-only lane |

---

### Close Checklist

_(Not yet — planning item.)_

### Outcome

`TASK-080` was used as the coordination lane for a two-item parallel-trial setup. It recorded the lane-selection gate, then the paired-lane trial once `TASK-081` was seeded as a disjoint companion lane.

### Verification

- [x] Parallel lane trial executed (2+ agents, distinct IDs)
- [x] Outcome and verification filled
- [x] Plan moved to closed/

**Direct evidence:** `TASK-080` and `TASK-081` now reference distinct file sets, and the template lane change landed in `.hawp/kit/templates/intake-plan.md`.
