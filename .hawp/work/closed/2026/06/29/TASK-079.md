## Task: Workflow Loop Orchestration — Planning

**Backlog ID:** TASK-079
**UUID:** `a7f3c2e1-9b4d-4a6f-8e2c-1d5f9a3b7c4e`
**Type:** improvement
**Reported:** 2026-06-25
**Risk Level:** medium
**Status:** in-progress (Phase 0 delivered)

---

### Goal

Design a **Workflow Loop** for long-running human–AI work in HAWP-enabled repositories — first as **instruction and artifact guidance**, with optional CLI orchestration deferred until tooling is explicitly allowed.

Large tasks (multi-file refactors, feature epics, audit remediation) often exceed a single agent session. Today HAWP shapes each unit of work and tracks it in `.hawp/work/`, but **iteration, handoff between roles, and retry-with-reflection are manual**. The Workflow Loop closes that gap through:

- A repeatable **Execution Loop**: continue → execute → review → approve/reject → reflect → retry or close.
- **Interchangeable participants** (human or agent) assigned via prompts, not runtime adapters.
- **Reusing** existing HAWP kit, backlog workflow, status reports, and the optional `checkpoint` field — not replacing them.

HAWP core remains a shaping protocol. Phase 0 is docs-only; the CLI engine in `librarian/` is a **future optional layer** gated on stakeholder approval when CLI/bash orchestration is in scope.

---

## Phase 0: Instruction-Only Loop (Current Version)

**Status:** delivered for this library version  
**Gate:** Use this phase until CLI tooling and bash orchestration are explicitly allowed.

### What shipped

| Deliverable | Path |
| ----------- | ---- |
| Usage guide (agents + humans) | `.hawp/kit/usage/workflow-loop.md` |
| Iteration handoff template | `.hawp/kit/templates/workflow-loop-handoff.md` |
| Cross-links | `.hawp/kit/start-here.md`, `.hawp/kit/usage/intake-workflow.md` |

### How it works (no runtime)

1. **Continue** — Load BACKLOG row, active plan, prior status handoffs and Iteration Log.
2. **Execute** — Implement one scoped slice per iteration (intake Step 5 discipline).
3. **Review** — Human or agent (review-only session) checks against plan constraints.
4. **Transition** — Approver decides: success, retry, park, or escalate (intake risk gates apply).
5. **Reflect** — On retry, append structured reflection to plan + handoff (**Next try**).
6. **Retry or close** — New session for next iteration, or intake Step 7 close on success.

### State model (artifact-only)

- **BACKLOG.md** — coordination status
- **Active plan** — Workflow Loop section + Iteration Log table
- **Status handoffs** — `.hawp/work/status/YYYY/MM/DD/<ID>-iter-<NNN>.md`
- **Evidence** — existing evidence paths
- **`checkpoint` field** — optional one-line pause anchor in prompts (per spec v0.1)

No `loop-runs/` tree, no JSON manifests, no subprocess adapters in Phase 0.

### Phase 0 success criteria

- [x] Instruction guide documents full loop with copy-paste prompts
- [x] Handoff template supports review, transition, and reflection
- [x] Intake workflow and start-here cross-linked
- [x] CLI/librarian implementation explicitly deferred (see Future Phases below)

### Phase 0 close (when ready)

- [ ] Stakeholder review of `workflow-loop.md`
- [x] Instruction loop trial on TASK-079 (3 gated + 5 autonomous iterations; see [autonomous trial summary](../../../../status/2026/06/25/TASK-079-autonomous-loop-trial.md))
- [ ] Optional: close TASK-079 as planning + Phase 0 complete, or keep open for Future Phase 1 gate

### Phase 0 close decision

| Option | Result | Impact |
| ------ | ------ | ------ |
| Close TASK-079 now | Mark Phase 0 complete and move the open question out of the active queue | Clears the active backlog and makes TASK-080 the next compounding step |
| Keep TASK-079 open | Hold a future gate for CLI/tooling follow-up | Preserves an open anchor, but keeps a resolved docs layer in active status |

**Recommended:** close TASK-079 after stakeholder review of `workflow-loop.md`, then use TASK-080 for scale planning and any future Phase 1 tooling gate.

---

## Future Phases: CLI / Librarian Orchestration (Out of Scope Now)

**Gate:** Do not start until stakeholders approve CLI tooling, bash scripts, and librarian runtime orchestration for this repo/version.

The architecture below is a **vision document** for optional tooling layered on Phase 0 artifacts. It must not redefine HAWP core or make the loop mandatory.

### Original CLI-oriented goal (deferred)

Design and plan a CLI-based orchestration layer that **automates** the same Execution Loop using participant adapters and persisted loop-run state.

- Running a repeatable **Execution Loop** via `./.hawp/bin/hawp loop …`
- Orchestrating **Participants** (human, coding-agent CLIs, GitHub PR checks, QA) through thin adapters
- **Reusing** Phase 0 kit docs as the policy source — not replacing them

This is **optional tooling** layered on HAWP artifacts. HAWP core remains a shaping protocol; the loop engine would live in `librarian/` and `.hawp/bin/hawp`, consistent with existing backlog CLI commands.

---

### Naming (HAWP-native)

| Term | Meaning |
| ---- | ------- |
| **Workflow Loop** | The outer continuous cycle for one work item until done or parked |
| **Execution Loop** | One pass through the stage pipeline (may repeat on retry) |
| **Stage** | A step in the pipeline: Task, Planner, Executor, Reviewer, Verifier |
| **Participant** | An actor that can run a stage (human, Cursor CLI, Copilot CLI, etc.) |
| **Participant Adapter** | Code that invokes a participant and normalizes its output |
| **Transition** | Policy-driven move to the next stage, Success, or Retry-with-Reflection |
| **Loop Run** | A single orchestrated session with persisted state under `.hawp/work/` |
| **Reflection** | Structured summary appended on retry (what failed, what to change) |

Do **not** use "Ralph" in user-facing docs, CLI, or file names.

---

### Current State Inventory

What exists today that this plan should build on (directly observed in repo):

#### Protocol and workflow (manual, agent-in-IDE)

| Asset | Path | Reuse |
| ----- | ---- | ----- |
| HAWP v0.1 shape (5 fields + optional checkpoint) | `core/.hawp/kit/references/spec.md` | Prompt compilation input; mission/constraints/output drive stage prompts |
| Conceptual pipeline stages (draft, not runtime) | `core/.hawp/kit/references/spec.md` (Pipeline Draft) | Map to Planner/Executor stages; clarify is pre-loop |
| Work intake loop (7 steps, investigation-first) | `core/.hawp/kit/usage/intake-workflow.md` | Policy source for review gates, risk levels, verify/close |
| Intake plan template | `core/.hawp/kit/templates/intake-plan.md` | Planner stage output schema |
| Status report / handoff | `core/.hawp/kit/usage/status-report.md`, `core/.hawp/kit/templates/status-report.md` | Reflection and cross-iteration handoff |
| **Workflow Loop (Phase 0)** | `.hawp/kit/usage/workflow-loop.md`, `.hawp/kit/templates/workflow-loop-handoff.md` | Instruction-based multi-iteration pattern (current) |
| Work item file tracking | `core/.hawp/kit/references/work-item-file-tracking.md` | Parallel lane rules, path discipline |
| Parallel work guardrails | `core/.hawp/kit/standards/patterns/parallel-work-guardrails.md` | Loop must respect lane boundaries |

#### Work coordination (filesystem state)

| Asset | Path | Reuse |
| ----- | ---- | ----- |
| Active backlog index | `.hawp/work/BACKLOG.md` | Task selection source |
| Active / parked / closed plan files | `.hawp/work/active/`, `parked/`, `closed/YYYY/MM/DD/` | Per-item truth; loop attaches to existing plan IDs |
| Evidence and status archives | `.hawp/work/evidence/`, `.hawp/work/status/` | Verifier output, iteration audit trail |
| Decisions | `.hawp/work/decisions/` | Architecture choices from loop design |

#### Provider architecture (multi-agent context)

| Asset | Path | Reuse |
| ----- | ---- | ----- |
| Provider manifest and install routing | `core/providers/manifest.yaml` | Know which overlays exist per provider |
| Shared behaviors (materialized) | `core/providers/shared/behaviors/*.md` | Canonical HAWP integration text for all participants |
| GitHub/Copilot pack | `core/providers/.github/` (instructions, prompts) | Copilot CLI prompt library; `intake.prompt.md` is closest to full-loop behavior |
| Cursor pack | `core/providers/.cursor/` (rules, `AGENTS.md.seed`) | Cursor CLI context injection |
| Continue pack | `core/providers/.continue/rules/` | Additional agent CLI target |
| Materialization pipeline | `librarian/scripts/providers/materialize/` | Pattern for generating participant-specific prompt fragments |
| Distribution guides | `distribution/generated/{github,cursor,continue}/` | Install/update; not loop runtime |

#### Tooling and CLI foundation

| Asset | Path | Reuse |
| ----- | ---- | ----- |
| HAWP CLI wrapper | `.hawp/bin/hawp` | Extend with `loop` subcommands |
| Librarian domain pattern | `librarian/scripts/README.md` (`index.ts` / `cli.ts` / `script.ts`) | New `workflow-loop/` domain |
| Backlog parser | `librarian/scripts/validate-hawp-workflow/orchestrate.ts` (`parseBacklog`) | Task selection, status reads |
| Workflow validation | `librarian/scripts/validate-hawp-workflow/` | Post-iteration integrity checks |
| Backlog upgrade | `librarian/scripts/backlog-upgrade/` | JSON-first internal models, formatters — pattern for loop state |
| Strategic orchestration notes | `.hawp/work/notes/2026/05/11/hawp-backlog-upgrade-command-design.md` (Part 10–11.5) | JSON-first objects, idempotency, validator authority |

#### Explicit boundaries (must respect)

| Statement | Source |
| --------- | ------ |
| "HAWP is not a runtime, validator, orchestrator, or memory system" | `README.md`, `core/.hawp/kit/references/spec.md`, provider behaviors |
| "Runtime orchestration design" is a v0.1 non-goal for the **protocol** | `core/.hawp/kit/references/spec.md` (Explicit Non-Goals) |
| Future: "Multi-agent orchestration" after validator/upgrade baseline | `hawp-backlog-upgrade-command-design.md` |

**Inference:** A Workflow Loop engine is acceptable as **optional librarian tooling** that reads/writes `.hawp/work/` artifacts and invokes external CLIs. It must not redefine HAWP core or downstream kit as a mandatory runtime.

---

### Gap Analysis

| Gap | Impact | Notes |
| --- | ------ | ----- |
| ~~No multi-iteration instruction pattern~~ | ~~Manual ad-hoc handoffs~~ | **Addressed in Phase 0** — `workflow-loop.md` |
| No loop runner CLI | Cannot automate multi-iteration work | `.hawp/bin/hawp` only has `backlog upgrade` / `backlog validate` — **deferred** |
| No stage pipeline implementation | Task→Planner→Executor→Reviewer→Verifier is conceptual only | Intake workflow covers similar steps but inside IDE sessions |
| No participant adapter layer | Cannot swap Cursor/Copilot/Claude/Codex/Gemini uniformly | Provider packs target IDE context, not subprocess CLIs |
| No loop run state schema | Iterations lack machine-readable continuity | **Partially addressed (Phase 0)** — human-readable Iteration Log in plan + status handoffs; JSON schema deferred to Phase 1 |
| No prompt compilation for CLIs | Each CLI has different invocation flags and I/O | Need template + shape → prompt builder |
| No evaluation / transition engine | Success vs retry is human judgment today | Reviewer/Verifier stages need structured rubrics |
| No reflection artifact on retry | Failed iterations lose structured "what to try next" | **Addressed (Phase 0)** — handoff template + Reflect section + plan Iteration Log |
| GitHub PR / QA not wired | External verification is manual | gh CLI exists; no HAWP wrapper |
| UUID backlog IDs incomplete | Parallel loop runs may collide on TASK-XXX | TASK-013 roadmap; parser already has UUID hooks |

---

### Proposed Architecture

#### Layering

```text
┌─────────────────────────────────────────────────────────┐
│  External participants (CLIs, human, gh, test runners) │
└───────────────────────────┬─────────────────────────────┘
                            │ Participant Adapters
┌───────────────────────────▼─────────────────────────────┐
│  Workflow Loop Engine (librarian/scripts/workflow-loop/) │
│  — stage runner, transitions, prompt compiler, evaluator │
└───────────────────────────┬─────────────────────────────┘
                            │ reads/writes
┌───────────────────────────▼─────────────────────────────┐
│  HAWP work artifacts (.hawp/work/, BACKLOG, kit)       │
└─────────────────────────────────────────────────────────┘
```

#### Stages (one Execution Loop pass)

| Stage | Purpose | Typical participant | Output artifact |
| ----- | ------- | ------------------- | --------------- |
| **Task** | Select backlog item; load plan; confirm scope | Human or loop CLI | Loop run manifest linked to plan ID |
| **Planner** | Shape/refine mission, constraints, output for this iteration | Planner CLI or human | Updated plan section or iteration brief |
| **Executor** | Implement changes per plan | Executor CLI (Cursor, Copilot, etc.) | Working tree diff + executor summary |
| **Reviewer** | Code/work review against constraints | Reviewer CLI or human | Review findings (pass / issues) |
| **Verifier** | Run checks (tests, workflow validate, custom) | Verifier (npm test, `hawp backlog validate`, QA) | Verification record with evidence links |
| **Transition** | Decide Success, Retry-with-Reflection, Park, or Escalate | Policy engine + optional human gate | Transition record; update BACKLOG status |

#### Participants (initial target set)

| Participant ID | Role | Integration sketch |
| -------------- | ---- | ------------------ |
| `human` | Any stage; approval gates | stdin/stdout prompts or pause-for-input |
| `cursor-cli` | Planner, Executor, Reviewer | Subprocess: Cursor agent CLI with repo context |
| `copilot-cli` | Planner, Executor, Reviewer | Subprocess: GitHub Copilot CLI |
| `claude-cli` | Planner, Executor, Reviewer | Anthropic CLI adapter |
| `codex-cli` | Executor | OpenAI Codex CLI adapter |
| `gemini-cli` | Planner, Reviewer | Google Gemini CLI adapter |
| `github-pr` | Reviewer, Verifier | `gh pr create`, check runs, review comments |
| `qa` | Verifier | Test/lint scripts, optional external QA harness |

Participants are **pluggable**. Default loop config maps stages → participant IDs; policies can override per risk level (reuse intake-workflow low/medium/high).

#### Policies and Transitions

| Policy | Source | Behavior |
| ------ | ------ | -------- |
| **Review gate** | `intake-workflow.md` Step 4 | Low: auto-advance; medium/high: human Transition approval |
| **Max iterations** | Loop config (default TBD) | After N retries → park or escalate |
| **Retry-with-Reflection** | New | Append reflection block to loop run state; re-enter Planner or Executor |
| **Path discipline** | `intake-workflow.md` | Reject participant output that fails path-lock gate |
| **Parallel lanes** | `parallel-work-guardrails.md` | Loop scoped to plan's overlapping files |
| **Validator authority** | Backlog upgrade design | `hawp backlog validate` must pass before Success transition |
| **Clean working tree** | `backlog-upgrade/script.ts` pattern | Optional guard before Executor stage |

Transition outcomes:

- **Success** → update plan Outcome/Verification, close per intake workflow
- **Retry** → write reflection, increment iteration, return to configured stage
- **Park** → move plan to `parked/`, update BACKLOG
- **Escalate** → human gate, status report to `.hawp/work/status/`

#### Loop run state (proposed)

Store under `.hawp/work/loop-runs/YYYY/MM/DD/<plan-id>-<run-id>/`:

```text
manifest.json       # plan ID, stages, participants, config, started_at
iterations/
  001/
    task.json
    planner.prompt.md / planner.output.md
    executor.prompt.md / executor.output.md
    reviewer.output.md
    verifier.output.json
    transition.json   # outcome, reflection (if retry)
```

JSON-first internal model (consistent with backlog-upgrade design); human-readable markdown mirrors for audit.

---

### Provider Integration Strategy

1. **Do not duplicate provider content.** Loop prompts compose from:
   - HAWP shape fields (mission, constraints, output)
   - Existing provider prompts where applicable (e.g. `core/providers/.github/prompts/intake.prompt.md` patterns)
   - Stage-specific templates in `librarian/scripts/workflow-loop/templates/` (new)

2. **Participant Adapter interface** (TypeScript, in `workflow-loop/adapters/`):

   ```ts
   type ParticipantAdapter = {
     id: string;
     invoke(input: StageInput): Promise<StageOutput>;
     supports(stage: StageName): boolean;
   };
   ```

   Each adapter wraps CLI discovery, argv construction, timeout, stdout/stderr capture, and exit-code interpretation.

3. **Materialization optional.** If loop prompt fragments should appear in provider packs, register emits in `librarian/scripts/providers/materialize/composition.ts` — same pipeline as shared behaviors.

4. **Provider manifest extension (future).** `core/providers/manifest.yaml` could list CLI binary hints and env vars per provider; not required for Phase 1.

---

### CLI Orchestration Pattern

One Execution Loop iteration:

```text
1. SELECT   hawp loop run --task TASK-079
            → parse BACKLOG.md, load active plan, create loop run dir

2. COMPILE   For current stage:
            → build prompt from shape + plan + prior reflections + kit refs

3. INVOKE    participant adapter runs external CLI (or human pause)

4. CAPTURE   stdout/stderr, exit code, artifacts (diff, logs)

5. EVALUATE  stage rubric (Reviewer/Verifier) → structured pass/fail/issues

6. TRANSITION
            success  → next stage or complete loop
            retry    → append reflection, goto configured stage
            park     → update BACKLOG + plan location

7. REPEAT    until Success, max iterations, or human abort
```

Proposed CLI surface (sketch):

```bash
./.hawp/bin/hawp loop run --task TASK-079 [--config loop.yaml]
./.hawp/bin/hawp loop step --run <run-id>   # single stage (debug)
./.hawp/bin/hawp loop status --run <run-id>
./.hawp/bin/hawp loop resume --run <run-id>
```

---

### Phased Implementation Plan

| Phase | Deliverable | Depends on | Status |
| ----- | ----------- | ---------- | ------ |
| **0 — Instruction-only loop** | `workflow-loop.md`, handoff template, kit cross-links | — | **Delivered** |
| **1 — Loop state model** | `workflow-loop/` domain skeleton; manifest + iteration JSON schema; read-only `loop status` | CLI tooling gate | deferred |
| **2 — Manual step runner** | `loop step` executes one stage with `human` participant; writes artifacts | Phase 1 | deferred |
| **3 — Prompt compiler** | Shape + plan → stage prompts; path-lock gate on outputs | Phase 2 | deferred |
| **4 — Pilot participant** | One CLI adapter (recommend Cursor CLI or Copilot CLI — decision needed) | Phase 3 | deferred |
| **5 — Full stage pipeline** | Task→Planner→Executor→Reviewer→Verifier wired with transitions | Phase 4 | deferred |
| **6 — Policy layer** | Review gates, max iterations, retry-with-reflection | Phase 5 | deferred |
| **7 — Additional participants** | Claude, Codex, Gemini, GitHub PR, QA adapters | Phase 6 | deferred |
| **8 — Docs sync** | Ensure CLI usage examples align with Phase 0 guide | Phase 7 | deferred |

**CLI gate:** Phases 1–8 require explicit approval that CLI/bash/librarian orchestration is in scope for this library version. Phase 0 is sufficient for Ralph-loop-like continuity via prompts and artifacts alone.

Each CLI phase (when allowed) ends with unit tests (`node:test`) and a short evidence note under `.hawp/work/evidence/`.

---

### Non-Goals

- **Ralph TUI** or any branded terminal UI for loop monitoring
- **Editor-specific orchestration** beyond existing provider packs (no new IDE plugins)
- **Redefining HAWP core** as a runtime, orchestrator, or memory system
- **Mandatory loop usage** — repos keep working with manual intake workflow only
- **Autonomous unbounded agents** — no skip of review gates or validator authority
- **Central cloud service** — loop runs locally against repo filesystem
- **Replacing** `intake-workflow.md` — the loop automates/extends it, not supersedes it

---

### Open Questions / Decisions Needed

| # | Question | Options | Recommendation |
| - | -------- | ------- | -------------- |
| 1 | Pilot CLI participant? | Cursor CLI vs Copilot CLI vs both | Start with one; Cursor if available in maintainer env |
| 2 | Loop state location? | `.hawp/work/loop-runs/` vs inside plan file frontmatter | Separate `loop-runs/` tree (keeps plans portable) |
| 3 | Default max iterations? | 3 / 5 / configurable | Configurable; default 5 |
| 4 | Human gate default for medium/high risk? | Pause every iteration vs Reviewer stage only | Pause at Transition for medium/high (matches intake) |
| 5 | UUID for loop runs vs plan ID? | UUID run IDs | UUID run IDs; link to legacy plan ID until TASK-013 lands |
| 6 | Install loop tooling downstream? | Ship in kit vs source-repo-only initially | Source-repo-only first (matches current `hawp` CLI scope) |
| 7 | Verifier stage default checks? | `npm test` + `hawp backlog validate` | Both when librarian present |
| 8 | Relationship to spec non-goals? | Document as tooling exception | ADR in `.hawp/work/decisions/` clarifying protocol vs tooling boundary |

---

### Success Criteria

Planning phase (this item):

- [x] Plan file created with inventory, gaps, architecture, phases, non-goals
- [x] BACKLOG.md updated with active row
- [x] Phase 0 instruction guide and handoff template delivered
- [ ] Stakeholder review and status → `done` (or keep open for Future Phase 1 gate)

Implementation phase (follow-on work items — **blocked on CLI gate**):

- [ ] `./.hawp/bin/hawp loop run` completes one full Execution Loop on a test task with at least one CLI participant
- [ ] Loop run artifacts are auditable (prompts, outputs, transitions on disk)
- [ ] Retry-with-Reflection produces a structured artifact consumed by the next iteration
- [ ] Review gate blocks auto-Success for medium/high risk without human approval
- [ ] `hawp backlog validate` passes after Successful close
- [ ] No duplication of provider behavior text (composition from shared sources)
- [ ] Documented in `.hawp/kit/usage/workflow-loop.md` with copy-paste examples — **Phase 0 done**

---

### Recommended Fix (after approval)

**Phase 0 (current):** Use `.hawp/kit/usage/workflow-loop.md` for multi-iteration work. No CLI required.

**When CLI tooling is allowed:**

1. Record decision ADR: protocol-vs-tooling boundary for Workflow Loop CLI.
2. Open **TASK-080** (implementation): Future Phase 1 — loop state model + `workflow-loop/` domain skeleton.
3. Spike pilot participant CLI (1–2 hours) to validate invoke/capture before Future Phase 4.

---

### Work Coordination

**Owner:** unassigned
**Implementation status:** Phase 0 delivered; CLI phases not started
**Overlapping files:**

- `.hawp/work/BACKLOG.md`
- `.hawp/work/active/TASK-079.md`
- `.hawp/kit/usage/workflow-loop.md`
- `.hawp/kit/templates/workflow-loop-handoff.md`
- Future (CLI gate): `librarian/scripts/workflow-loop/**`, `.hawp/bin/hawp`

**Parallel work risk:** low
**Can implement now:** Phase 0 yes (docs); CLI phases no — await CLI tooling gate

---

### Workflow Loop

**Loop status:** active (autonomous trial — Phase 0+)
**Loop mode:** autonomous
**Iteration budget:** 8
**Current iteration:** 8
**Executor:** agent
**Reviewer:** agent (review-only hat within session)
**Approver:** agent-auto (low-risk doc pass; `auto-approve: true` for trial)
**Auto-approve:** true

**Trial scope:** Iterations 1–3 = gated manual trial (2026-06-25). Iterations 4–8 = autonomous self-driving trial (2026-06-26) — five passes without human "loop again" between each.

---

### Iteration Log

| Iter | Date | Outcome | Handoff |
| ---- | ---- | ------- | ------- |
| 1 | 2026-06-25 | retry | [status](../../../../status/2026/06/25/TASK-079-iter-001.md) |
| 2 | 2026-06-25 | retry | [status](../../../../status/2026/06/25/TASK-079-iter-002.md) |
| 3 | 2026-06-25 | approve (gated trial) | [status](../../../../status/2026/06/25/TASK-079-iter-003.md) |
| 4 | 2026-06-26 | auto-advance | [status](../../../../status/2026/06/25/TASK-079-iter-004.md) |
| 5 | 2026-06-26 | auto-advance | [status](../../../../status/2026/06/25/TASK-079-iter-005.md) |
| 6 | 2026-06-26 | auto-advance | [status](../../../../status/2026/06/25/TASK-079-iter-006.md) |
| 7 | 2026-06-26 | auto-advance | [status](../../../../status/2026/06/25/TASK-079-iter-007.md) |
| 8 | 2026-06-26 | success (autonomous trial) | [status](../../../../status/2026/06/25/TASK-079-iter-008.md) |

---

### Outcome

Phase 0 (instruction-only Workflow Loop) delivered and dogfooded. The loop guidance lives in `.hawp/kit/usage/workflow-loop.md`; the optional CLI engine remains deferred pending explicit approval. Eight loop iterations were run against this task itself: 3 gated (2026-06-25) and 5 autonomous (2026-06-26), ending in success. Recommendation adopted: autonomous mode for low-risk doc work, gated default for medium/high-risk changes.

### Verification

- Iteration Log above links all 8 handoff status reports under `status/2026/06/25/`.
- Autonomous trial summary: [status/2026/06/25/TASK-079-autonomous-loop-trial.md](../../../../status/2026/06/25/TASK-079-autonomous-loop-trial.md).
- Kit guidance file exists and is linked from `.hawp/kit/start-here.md`.

### Close Checklist

- [x] Outcome section filled
- [x] Verification section filled
- [x] Plan file moved to `closed/YYYY/MM/DD/`
- [x] BACKLOG.md updated

---

## Loop Trial Report (2026-06-25 — gated)

Prior summary referenced in BACKLOG notes. Gated trial: 3 iterations with explicit approve/retry between passes.

## Autonomous Loop Trial (2026-06-26)

Full summary: [status/2026/06/25/TASK-079-autonomous-loop-trial.md](../../../../status/2026/06/25/TASK-079-autonomous-loop-trial.md)

**Iterations:** 5 autonomous (4–8) | **Outcome:** Phase 0+ self-driving loop documented and dogfooded | **Recommendation:** adopt autonomous mode for low-risk doc work; keep gated default for medium/high production changes.
