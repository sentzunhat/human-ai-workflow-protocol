## ADR: Intake Template Refinement for Bug/Task Work

### Title

Add minimal optional intake templates so agents can stop at analysis and planning before implementation.

### Context

HAWP is intentionally lean and non-runtime. The current repository already includes:

- Core protocol shape and boundaries in README/SPEC/types.
- Authoring guidance in AUTHORING_PATTERNS.
- Optional templates for micro/standard shaping, status reports, and audits.
- A repo-local intake workflow with backlog states and review gate.
- A generic intake plan template used in usage/status work items.

Requested outcome: decide whether to add or refine intake templates for incoming bugs, tasks, improvements, and small project work, while preserving HAWP's file-first, markdown-first posture.

### Problem

Current intake support is operationally present but not ergonomically split for common agent behavior:

- Intake exists as workflow plus one general plan template.
- The plan template is broad and not explicitly bug-first.
- There is no dedicated optional work-intake template focused on intake normalization before deeper planning.

This can increase variance in how agents triage bug/task requests and can encourage premature implementation when the intake artifact is underspecified.

### Current State Found in Repo

Directly verified:

- Intake loop and risk gate exist in core/.hawp/usage/INTAKE_WORKFLOW.md.
- Backlog states already include inbox -> analyzing -> plan-ready -> approved -> in-progress -> done/blocked/wont-fix in core/.hawp/usage/BACKLOG.md.
- A generic plan template exists at core/.hawp/templates/intake-plan.md with evidence-vs-inference sections.
- Intake automation guidance exists in .github/instructions/intake.instructions.md and .github/prompts/intake.prompt.md.
- Optional templates currently focus on general shaping, status, and audit artifacts (not bug-specific intake split).

Inference (not directly proven in this repo):

- A dedicated bug-focused plan template will likely improve consistency for asset/path/runtime symptom bugs like the landing hero image example.
- Two concise templates should be enough for most bug/task intake without adding category sprawl.

### Options

#### Option A: One general template only

Add core/.hawp/templates/work-intake.md and keep core/.hawp/templates/intake-plan.md as-is.

Pros:

- Smallest file addition.
- Lowest maintenance overhead.

Cons:

- Bug analysis still relies on generic plan framing.
- Less explicit structure for symptom -> root cause -> verification workflows.

#### Option B: Two-template split (recommended)

Add:

- core/.hawp/templates/work-intake.md
- core/.hawp/templates/bug-plan.md

Pros:

- Still small and optional.
- Separates intake normalization from bug planning.
- Aligns with existing backlog/risk workflow and evidence discipline.
- Preserves existing intake-plan.md for broad/non-bug work.

Cons:

- Slight overlap risk if wording duplicates intake-plan.md.
- Requires minimal navigation updates in onboarding docs.

#### Option C: Larger template set

Add work-intake.md, bug-plan.md, task-plan.md, improvement-plan.md, implementation-review.md.

Pros:

- Maximum specificity per work type.

Cons:

- Higher maintenance burden.
- Higher drift risk across near-duplicate templates.
- Over-structures v0.1 usage and increases cognitive load.

### Recommended Option

Recommend Option B.

Rationale:

- It is the smallest useful addition that materially improves AI intake behavior.
- It stays copy-first and optional, without runtime or schema expansion.
- It captures the practical split between intake and bug planning seen in the existing usage workflow.

### Files Proposed to Add/Change

Proposed additions:

- core/.hawp/templates/work-intake.md
- core/.hawp/templates/bug-plan.md

Proposed minimal navigation updates:

- core/.hawp/START_HERE.md
- core/.hawp/README.md
- README.md (only if needed for top-level discoverability)

No protocol/schema changes proposed.

### Files Explicitly Not to Touch

- core/.hawp/types/shape.ts
- core/.hawp/SPEC.md (unless a tiny docs-only clarification is explicitly approved)
- Runtime/tooling surfaces (no CLI, no validators, no package scripts, no SQLite, no orchestrator artifacts)

### Risks

- Template duplication drift against core/.hawp/templates/intake-plan.md.
- Confusion if users treat templates as required instead of optional.
- Scope creep toward project-management behavior if backlog/process prose is over-expanded.

### Acceptance Criteria

- Two new optional templates exist and are copy-first, concise, and markdown-only.
- Templates explicitly separate directly verified facts from inference.
- Templates map cleanly to existing backlog status flow and review gate.
- No changes to core shape/type semantics.
- Navigation updates are minimal and do not duplicate large sections from AUTHORING_PATTERNS.

### Verification Checklist

- [x] Confirm both new templates are under core/.hawp/templates/.
- [x] Confirm wording states optional usage aids and non-runtime scope.
- [x] Confirm bug-plan template includes: ID, type, reported date, risk level, evidence/inference split, options, recommendation, files to change, verification checklist.
- [x] Confirm START_HERE/README references resolve and stay concise.
- [x] Confirm core/.hawp/types/shape.ts unchanged.
- [x] Confirm no new scripts, runtime commands, validators, or data stores were introduced.

### Implementation Status

Implemented on 2026-04-30 after verification.

Changed files:

- core/.hawp/templates/work-intake.md
- core/.hawp/templates/bug-plan.md
- core/.hawp/START_HERE.md
- core/.hawp/README.md

Confirmed unchanged by this implementation:

- core/.hawp/types/shape.ts
- core/.hawp/SPEC.md
- runtime/tooling/package-script surfaces

### Original Implementation Plan (Historical)

1. Draft core/.hawp/templates/work-intake.md as a short intake-normalization artifact:
   - Intake summary
   - Classification (bug/task/improvement)
   - Impact and scope hint
   - Initial evidence vs inference
   - Risk estimate and gate expectation
   - Backlog linkage fields
2. Draft core/.hawp/templates/bug-plan.md as a bug-focused planning artifact:
   - ID/type/reported/risk
   - Symptom, root cause hypothesis, verified facts, inferred items
   - Options and recommended fix
   - File touch list and verification checklist
   - Status progression markers
3. Apply minimal navigation links in core/.hawp/START_HERE.md and core/.hawp/README.md, and README.md only if needed.
4. Run a docs consistency pass to ensure optional/non-runtime language is preserved.
5. Stop at review gate and request approval before any further scope expansion.
