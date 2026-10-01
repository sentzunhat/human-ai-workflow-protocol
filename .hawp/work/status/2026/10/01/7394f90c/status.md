---
uuid: 7394f90c
title: "Parallel coordination and .hawp/.space continuation checkpoint"
type: status
date: 2026-10-01
---

# 2026-10-01 — Parallel coordination and `.hawp/.space` direction agreed

## Intent

Preserve the continuation state from the design discussion that turned a
downstream multi-agent workflow lesson into a public HAWP improvement plan.

This checkpoint is intentionally public-safe. It records the reusable HAWP
architecture lesson without naming or exposing the private downstream
organization, repository, infrastructure, tickets, credentials, or personal
context that motivated the analysis.

## Before

HAWP already had:

- a small task-shaping protocol and investigation → plan → review → implement →
  verify → close workflow;
- durable coordination state under `.hawp/work/`;
- owner/file-overlap/parallel-risk fields in work plans;
- guardrails that default conflicting same-file work to hold;
- Workflow Loop roles for executor/reviewer/approver;
- provider overlays generated from shared behavior sources.

What it did **not** have was a canonical live workspace convention for parallel
Git worktrees or a clear distinction between a coordinator checkout and worker
execution lanes.

The repository also had an older plural `.hawp/.spaces/` concept used for
embedded/frozen agent-worktree mirror snapshots. A prior status record explicitly
treated those mirrors as reference snapshots rather than live docs.

## During

The conversation:

1. inspected a real downstream HAWP installation that had evolved a coordinator
   checkout plus isolated per-branch worktrees;
2. confirmed that the pattern reduced checkout collisions but could still drift
   when implementation leaked into the coordinator checkout;
3. compared that behavior with upstream HAWP's existing parallel-work guardrails;
4. concluded that HAWP should absorb the reusable coordination lesson without
   copying private project-specific branch/ticket conventions;
5. renamed the generic "manager" concept to **Coordinator**, making it a role
   rather than a mandatory branch;
6. agreed that worker implementation should use local isolated worktrees;
7. agreed on the singular canonical root `.hawp/.space/`;
8. agreed that branch paths mirror filesystem paths directly;
9. agreed that generic categories such as `feature/`, `bug/`, `fix/`,
   `infra/`, `docs/`, `test/`, `refactor/`, `chore/`, and `spike/`
   may be recommended while tracker prefixes remain downstream policy;
10. kept the HAWP core Shape unchanged and rejected turning this into a runtime
    orchestration framework.

Repository inspection then surfaced the legacy plural `.hawp/.spaces/` history,
so explicit reconciliation of plural snapshot semantics vs singular live-worktree
semantics was added to the work item.

## After

The durable design direction is now:

```text
.hawp/work/
    committed, shared, durable workflow/project state

.hawp/.space/
    local, gitignored, disposable worker worktrees
```

Branch/worktree mapping:

```text
feature/add-context-search
→ .hawp/.space/feature/add-context-search/

fix/worktree-collision
→ .hawp/.space/fix/worktree-collision/
```

The Coordinator:

```text
shape
→ inspect active work
→ resolve dependencies/overlap
→ dispatch worker lanes
→ review/integrate
→ update durable HAWP state
```

Workers implement from their assigned lane. A coordinator checkout should not
silently accumulate implementation that belongs to an active worker lane.

No implementation has been completed yet. This checkpoint and the linked active
plan represent an approved design direction / plan-ready state, not a claim that
the feature already exists.

## Strategic Impact

This direction extends HAWP from "parallel work can be coordinated" to a concrete,
portable execution convention while preserving its existing scope:

- the HAWP Shape stays small;
- Git remains the source of branch/worktree mechanics;
- `.hawp/work/` remains durable truth;
- `.hawp/.space/` becomes disposable local execution;
- provider agents can receive equivalent parallel-lane behavior;
- coordinator drift becomes an explicit safety concern;
- repository-specific ticket naming remains outside HAWP core.

The main architectural dependency is the relationship between legacy
`.hawp/.spaces/` frozen mirrors and the proposed live `.hawp/.space/` lane.
That must be decided deliberately during implementation.

## Continuation State

**Current state:** plan-ready.

**Owning work item:** `7394f90c` — Parallel coordination with isolated
`.hawp/.space/` worktree lanes.

**Latest completed milestone:** public-safe knowledge transfer reconstructed,
existing HAWP coordination conventions inspected, singular `.space` direction
agreed, and legacy plural `.spaces` conflict discovered.

**Next milestone:** Phase 1 investigation and design reconciliation.

**Next actions:**

1. inspect all current references to `.hawp/.spaces/`;
2. decide snapshot/plural compatibility vs deprecation strategy;
3. inspect install/update/provider/distribution paths for local-state preservation;
4. finalize the smallest guidance-first implementation surface;
5. implement the core guide, guardrail update, plan template/metadata, provider
   behavior, and ignore/preservation rules;
6. run validation and dogfood two independent lanes.

**Blockers:** none known, but implementation should not begin until active
provider/distribution path overlap is rechecked.

**Unresolved questions:**

- exact legacy plural `.spaces` transition;
- whether singular `.space` ignore behavior is installer-managed or documented;
- whether coordinator-drift protection stays guidance-only or gains a small
  validator using existing tooling.

## Decisions Not to Revisit Without New Evidence

- singular `.hawp/.space/` is the target live execution root;
- branch path mirrors `.space` path;
- `.space` is local/disposable and `.hawp/work` is durable/shared;
- Coordinator is a role, not a mandatory branch;
- tracker prefixes are downstream policy;
- core Shape does not grow;
- no orchestration runtime in this work item.

## Resume From

Open `.hawp/work/active/7394f90c/plan.md` and begin **Phase 1 —
Investigation**, starting with the legacy `.hawp/.spaces/` references and
install/update preservation boundaries.
