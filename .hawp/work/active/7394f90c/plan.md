# Parallel coordination with isolated `.hawp/.space/` worktree lanes

**UUID:** `7394f90c-8498-4abf-91a7-74bd1070162a`
**Type:** improvement
**Reported:** 2026-10-01
**Status:** `plan-ready`
**Risk Level:** medium
**Target Version:** unscheduled

---

## Input

Add the parallel-work lessons learned from a real downstream HAWP deployment back
into HAWP as a reusable, public-safe work item.

The agreed target is a local execution-space convention where every worker branch
maps to an isolated worktree beneath:

```text
.hawp/.space/<branch>
```

Example:

```text
repo/
├── .hawp/
│   ├── kit/
│   ├── work/
│   └── .space/                 # local parallel worktrees, gitignored
│       ├── feature/
│       │   └── add-ecr-publishing/
│       ├── fix/
│       │   └── postgres-restore/
│       ├── infra/
│       │   └── staging-oidc/
│       └── docs/
│           └── migration-guide/
│
├── src/
└── ...
```

The public work item intentionally omits private downstream organization, project,
ticket, infrastructure, and personal context. The lesson is the reusable workflow
pattern, not the private source project.

---

## Context

HAWP already has the core coordination pieces needed for parallel work:

- `.hawp/work/BACKLOG.md` as the active coordination index;
- one work item / one plan ownership;
- owner, implementation status, overlapping files, risk, and implementation gate
  fields in active plans;
- exact repo-relative path discipline;
- parallel-work guardrails that default same-file overlap to hold;
- Workflow Loop roles such as executor, reviewer, and approver;
- status/evidence/decision artifacts under `.hawp/work/`.

What is missing is a canonical local execution-space convention that keeps
parallel worker implementation isolated from the coordinator checkout.

A historical repository detail must also be reconciled: HAWP previously used
the plural path `.hawp/.spaces/` for embedded agent-worktree mirror snapshots.
Those mirrors were explicitly treated as frozen reference material and excluded
from live link validation. This work proposes the singular `.hawp/.space/` for
**live, disposable, gitignored execution lanes**. The implementation must define
the relationship between the old plural snapshot concept and the new singular
live-worktree concept rather than silently conflating them.

---

## Mission

Define and implement an optional **HAWP Parallel Coordination** pattern that:

1. introduces `.hawp/.space/` as the canonical local execution-space root;
2. maps one Git branch to one isolated local worktree below that root;
3. mirrors branch paths directly in the filesystem;
4. defines a generic **Coordinator** role and explicit Worker lanes;
5. keeps worker implementation out of the coordinator checkout when a lane owns it;
6. integrates with existing HAWP backlog, plan, status, evidence, overlap, and
   Workflow Loop conventions;
7. propagates equivalent guidance to supported provider overlays;
8. keeps `.hawp/.space/` disposable, local, and gitignored;
9. preserves existing HAWP repositories that do not adopt parallel coordination;
10. leaves the HAWP core Shape unchanged.

---

## Decisions Already Made

These are accepted design constraints for this item and should not be reopened
without new evidence:

### 1. Canonical live execution root

Use:

```text
.hawp/.space/
```

Singular `.space`, nested inside `.hawp`.

### 2. Branch-to-space mapping

The worktree path mirrors the branch name exactly beneath `.hawp/.space/`:

```text
feature/add-context-search
→ .hawp/.space/feature/add-context-search/

fix/worktree-collision
→ .hawp/.space/fix/worktree-collision/

infra/staging-oidc
→ .hawp/.space/infra/staging-oidc/

docs/parallel-coordination
→ .hawp/.space/docs/parallel-coordination/
```

### 3. Generic branch categories

HAWP should recommend a small extensible vocabulary such as:

```text
feature/<short-summary>
bug/<short-summary>
fix/<short-summary>
infra/<short-summary>
refactor/<short-summary>
docs/<short-summary>
test/<short-summary>
chore/<short-summary>
spike/<short-summary>
```

Repository-specific ticket prefixes remain downstream policy.

### 4. Coordinator is a role, not a mandatory branch

A Coordinator may be a human, agent, session, original checkout, or optional
coordination branch. HAWP must not require a permanent "manager branch."

### 5. Durable truth vs disposable execution

```text
.hawp/work/   = durable, committed, shared project/workflow state
.hawp/.space/ = local, gitignored, disposable execution worktrees
```

### 6. No protocol-shape growth

Do not add branch, worktree, lane, manager, coordinator, or workspace fields to
the HAWP core Shape.

### 7. No orchestration runtime

This remains guidance, templates, provider behavior, and narrowly scoped safety
support. Do not build a scheduler, daemon, lock service, merge engine, or agent
orchestrator as part of this item.

---

## Analysis

### Directly verified

- HAWP's current plan template already includes `Work Coordination`,
  overlapping files, parallel-work risk, and `Can implement now`.
- Current parallel-work guidance already defaults overlapping edits to a hold
  unless explicitly approved.
- Current status-report convention is
  `.hawp/work/status/YYYY/MM/DD/<work-id>/status.md`.
- Current `.gitignore` ignores legacy `.hawp/.spaces/`, not the newly
  proposed singular `.hawp/.space/`.
- A 2026-09-10 status record describes legacy `.hawp/.spaces/` as embedded
  frozen agent-worktree mirrors/reference snapshots rather than live docs.
- At checkpoint intake on 2026-10-01, `main` and `development` pointed to
  the same repository state immediately after the v0.0.24 merge.

### Inferred / to verify during implementation

- Whether legacy `.hawp/.spaces/` should remain documented as a deprecated
  snapshot-only concept, be retired in a separate hygiene item, or be migrated.
- Whether install/update scripts should add the `.hawp/.space/` ignore rule
  automatically to downstream repositories or only document the requirement.
- Whether existing validation should merely ignore `.hawp/.space/` or also
  detect unsafe tracked contents.
- Whether coordinator-drift protection should remain guidance-only in this item
  or gain a lightweight validation check using existing HAWP tooling.

---

## Work Coordination

**Owner:** unassigned
**Implementation status:** not-started

**Expected overlapping / touched paths:**

- `core/.hawp/kit/standards/patterns/parallel-work-guardrails.md`
- `core/.hawp/kit/start-here.md`
- `README.md`
- `.gitignore`
- `core/.hawp/kit/usage/parallel-coordination.md` — planned, not present yet
- `core/.hawp/kit/templates/parallel-work-plan.md` — planned, not present yet
- `core/providers/shared/behaviors/hawp-parallel.md` — planned, not present yet
- install/update/distribution source paths discovered during Phase 1
- tests/validation paths discovered during Phase 1

**Parallel work risk:** medium
**Can implement now:** only after approval / explicit assignment

**Coordination note:**

Before implementation, re-read `.hawp/work/BACKLOG.md` and verify no active
item owns the same provider/distribution/install paths. Keep this item separate
from unrelated release or MCP work.

**Path discipline:** use exact repo-relative paths. Planned files that do not yet
exist must be explicitly marked planned-not-present until created.

---

## Options

### Option A — Guidance-first parallel coordination

Add the usage guide, template, expanded guardrails, provider behavior, ignore
boundary, and preservation tests. Reuse existing HAWP coordination state and Git
itself for worktree operations.

**Pros:** aligns with HAWP's protocol-first scope; small conceptual surface;
portable; no new runtime.

**Cons:** branch/worktree creation remains manual or repository-specific.

### Option B — Add a first-class HAWP worktree CLI/runtime

Add commands that create, resume, dispatch, lock, and clean up worktrees.

**Pros:** stronger automation.

**Cons:** expands HAWP toward orchestration/runtime ownership, increases safety
surface, and is not required to capture the proven workflow lesson.

### Recommended direction

Choose **Option A** for this work item.

If later evidence shows repetitive worktree operations cause real failures,
shape a separate follow-up item for minimal tooling rather than expanding this
scope preemptively.

---

## Proposed Public Model

```text
                       Coordinator
                            │
                    reads .hawp/work/
                            │
          ┌─────────────────┼─────────────────┐
          │                 │                 │
          ▼                 ▼                 ▼
      Worker A          Worker B          Worker C
     feature/foo        fix/bar           infra/baz
          │                 │                 │
          ▼                 ▼                 ▼
 .hawp/.space/      .hawp/.space/      .hawp/.space/
   feature/foo          fix/bar           infra/baz
          │                 │                 │
          └────────────┬────┴────────────┬────┘
                       │
                normal Git review
                       │
                       ▼
                    base branch
```

The repository's normal Git review/merge policy remains authoritative.

---

## Parallel Coordination Plan Block

Add or recommend an optional plan section similar to:

```markdown
### Parallel Coordination

**Mode:** parallel
**Coordinator:** human | agent | unassigned

**Base branch:** <repo policy>
**Base revision:** <sha, optional>

**Worker owner:** human | agent | unassigned
**Branch:** infra/staging-oidc
**Space:** .hawp/.space/infra/staging-oidc

**Dependencies:**
- none

**Expected touched paths:**
- src/environments/staging/
- src/modules/oidc/

**Parallel work risk:** low | medium | high
**Can implement now:** yes | no | only after approval

**Coordination note:**
<short note>
```

This is operational plan metadata, not a new protocol Shape.

---

## Dispatch Gate

Before creating or resuming a lane, the Coordinator should establish:

1. owning HAWP work item;
2. repository base/integration branch;
3. current/fresh base state;
4. whether the branch already exists locally;
5. whether the branch already exists remotely;
6. whether a worktree already owns the branch;
7. expected touched paths;
8. overlap with active HAWP work;
9. dependency state;
10. whether implementation may proceed.

Same-file overlap defaults to hold unless sequencing or explicit approval
resolves ownership.

---

## Worker Rules

A Worker lane should:

1. operate from its declared `.hawp/.space/<branch>` checkout;
2. read the owning HAWP plan first;
3. modify only assigned scope;
4. respect dependencies and file-overlap boundaries;
5. leave unrelated work untouched;
6. verify changes using repository-specific checks;
7. separate direct evidence from inference;
8. hand results back through `.hawp/work/` artifacts;
9. follow normal repository PR/merge policy;
10. remove the disposable worktree only after safe integration/closure.

One work item has one implementation owner at a time. If one item needs multiple
independent executors, split it into separate HAWP work items with explicit
dependencies and file boundaries.

---

## Coordinator Drift Protection

### Problem

A long-lived coordination checkout can gradually accumulate implementation
changes that belong to worker lanes. That creates ambiguous ownership and makes
parallel work unsafe.

### Required behavior

When a worker lane owns a change, the Coordinator should not silently implement
that change in the coordinator checkout.

If implementation appears there:

```text
Coordinator detects implementation work
               ↓
Does an active lane own it?
        ┌──────┴──────┐
       yes            no
        │              │
 resume lane      shape/assign work
        │              │
        └──── dispatch worker
```

A Coordinator may intentionally become the Executor, but ownership must be made
explicit first.

---

## Relationship to Workflow Loop

Parallel Coordination and Workflow Loop are complementary:

- **Parallel Coordination** answers where/under whose lane work executes.
- **Workflow Loop** answers how one lane iterates through
  execute → review → reflect → retry/close.

A single `.hawp/.space/<branch>` lane may run a Workflow Loop without changing
the parallel coordination model.

---

## Proposed Implementation Phases

### Phase 1 — Investigation

Inspect current:

- parallel-work guardrails;
- work-item file tracking;
- intake workflow;
- Workflow Loop;
- provider shared behaviors/materialization;
- install/update preservation rules;
- generated distribution structure;
- `.gitignore`;
- link validation and legacy `.hawp/.spaces/` handling.

Explicitly decide how legacy plural `.spaces` relates to singular live
`.space`.

### Phase 2 — Core guidance

Create planned
`core/.hawp/kit/usage/parallel-coordination.md`.

Expand
`core/.hawp/kit/standards/patterns/parallel-work-guardrails.md`.

Update
`core/.hawp/kit/start-here.md`.

Add a reusable parallel-work plan template if it improves authoring without
duplicating the intake-plan template.

### Phase 3 — Provider parity

Add shared provider behavior (planned
`core/providers/shared/behaviors/hawp-parallel.md` or the current canonical
equivalent).

Regenerate provider overlays using the repository's normal sync process.

### Phase 4 — Local-state safety

Add `.hawp/.space/` to the appropriate ignore/preservation model.

Ensure install/update/distribution behavior never copies, deletes, normalizes,
or overwrites live local worktrees.

### Phase 5 — Tests and validation

Add evidence for:

- ignore behavior;
- update/install preservation;
- provider materialization consistency;
- multiple sibling worktrees coexisting;
- cleanup preserving `.hawp/work/`;
- no tracked `.hawp/.space/` contents in normal operation.

Run the repository's current quality and workflow validation commands.

### Phase 6 — Dogfood

Exercise at least two independent low-risk lanes under the new convention.

Example:

```text
.hawp/.space/docs/parallel-coordination-example/
.hawp/.space/test/space-preservation/
```

Verify Coordinator cleanliness, independent branch ownership, and preserved
HAWP durable state.

---

## Acceptance Criteria

### Documentation

- [ ] Parallel Coordination is documented as optional HAWP usage guidance.
- [ ] `.hawp/.space/` is defined as ephemeral local execution state.
- [ ] `.hawp/work/` remains durable shared project/workflow state.
- [ ] Coordinator and Worker responsibilities are explicit.
- [ ] Branch → `.space` mapping is deterministic.
- [ ] Generic branch categories are documented and extensible.
- [ ] Tracker-specific naming remains downstream policy.
- [ ] Coordinator drift protection is documented.
- [ ] Legacy `.hawp/.spaces/` semantics are explicitly reconciled.

### Repository / distribution safety

- [ ] `.hawp/.space/` is gitignored or otherwise reliably excluded from commits.
- [ ] No live `.hawp/.space/` contents are distributed.
- [ ] Install/update logic does not overwrite or delete `.hawp/.space/`.
- [ ] Existing HAWP repositories remain compatible without adopting the pattern.

### Providers

- [ ] Claude guidance has equivalent semantics.
- [ ] Codex guidance has equivalent semantics.
- [ ] Cursor guidance has equivalent semantics.
- [ ] Continue guidance has equivalent semantics.
- [ ] GitHub/Copilot guidance has equivalent semantics.
- [ ] Generated provider outputs are synchronized from canonical shared sources.

### Verification

- [ ] Existing HAWP validation passes.
- [ ] Provider/distribution sync passes.
- [ ] Relevant install/update preservation tests pass.
- [ ] A test or controlled dogfood run proves multiple sibling worktrees coexist.
- [ ] Cleanup preserves durable `.hawp/work/` state.
- [ ] Core HAWP Shape remains unchanged.
- [ ] No runtime orchestrator or mandatory coordinator branch is introduced.

---

## Out of Scope

- automatic branch deletion;
- automatic force reset;
- automatic merging;
- background agent scheduling;
- runtime lock services;
- mandatory GitHub PR semantics;
- mandatory `main`, `development`, or other base branch names;
- project-specific ticket naming;
- private downstream infrastructure details;
- automatic migration/deletion of historical `.hawp/.spaces/` mirrors unless
  separately approved.

---

## Expected Output

Deliver:

1. reviewed Parallel Coordination design;
2. singular `.hawp/.space/` convention;
3. legacy `.hawp/.spaces/` reconciliation decision;
4. expanded parallel-work guardrails;
5. Coordinator/Worker guidance;
6. reusable plan metadata/template;
7. provider parity;
8. local-state ignore/preservation behavior;
9. tests/validation;
10. dogfood evidence;
11. final status/handoff explaining what changed and what remains optional.

---

## Verification (current planning state)

- [x] Existing parallel-work guardrails inspected.
- [x] Existing intake plan coordination block inspected.
- [x] Existing status artifact convention inspected.
- [x] Current `.gitignore` inspected.
- [x] Legacy `.hawp/.spaces/` history identified and added as a reconciliation requirement.
- [x] Public checkpoint/work item stripped of private downstream project and personal context.
- [ ] Implementation not started.
- [ ] Provider changes not started.
- [ ] Install/update preservation not yet verified for singular `.hawp/.space/`.
- [ ] Dogfood not started.

## Close Checklist

- [ ] Outcome section filled after implementation.
- [ ] Verification claims backed by direct evidence.
- [ ] Evidence artifacts linked if needed.
- [ ] Final status report updated.
- [ ] Plan moved to closed archive on completion.
- [ ] BACKLOG entry moved to Recently Closed on completion.
