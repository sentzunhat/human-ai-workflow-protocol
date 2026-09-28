# 2026-07-28 — Agent review, planning, and HAWP adoption instructions established

- **Event date:** 2026-07-28
- **Archive date:** 2026-09-28
- **Canonical project:** Human-AI Workflow Protocol (HAWP)
- **Related projects:** `mode40/mast_infrastructure`, Sunshine Medical Clinic repository
- **Evidence basis:** archived conversation context, prior HAWP continuity records, and current `sentzunhat/human-ai-workflow-protocol` repository structure

## What changed

### Before

- Sage could treat a green CI result or a reassuring PR description as sufficient evidence for approval, even when a Terraform plan contained destructive or unrelated changes.
- Weekly ticket-hour reconstruction lacked a strict evidence model, a fixed daily total, and safeguards against double-counting a ticket and its PR.
- Infrastructure portfolio investigation and prioritization were combined into one large activity, making it easy to mix incomplete discovery with sequencing decisions.
- HAWP installation guidance existed as separate provider-specific instructions and did not yet form one coherent, repository-safe workflow for all supported providers.
- Repository branch roles were initially ambiguous between HAWP's source branch and each target repository's base/PR branch.

### During

Reusable digital-agent instructions were created and refined for the following workflows:

1. **Terraform PR review for Sage**
   - Treats the plan as evidence to inventory, not trusted instructions.
   - Requires plan freshness, target-environment verification, scope mapping, secret redaction, and explicit review outcomes.
   - Blocks approval for unexplained destroys/replacements, unrelated material actions, stale or incomplete plans, or sensitive-data exposure.
   - Prohibits approval or PR-state changes without explicit authorization.

2. **MES-2806 follow-up for Sage**
   - Requires discovery of the real Terraform plan/apply input path instead of assuming `terraform.examplevars` configures a deployed environment.
   - Keeps the work in Review until plan, apply, SNS confirmation, and mailbox delivery are verified.
   - Requires current-head and approval-SHA checks and avoids duplicate review-request comments.

3. **Weekly ticket-hour reconstruction for Jarvis**
   - Reconstructs Tuesday through Friday at exactly eight hours per day and 32 hours total.
   - Uses ticket-first reporting, associates related PRs without double-counting, and separates evidence from estimated duration.
   - Excludes activity outside the reporting window and does not submit or mutate a timesheet.

4. **Infrastructure epic investigation and sequencing for Sage**
   - Split into two ordered runs:
     1. read-only discovery and complete epic/ticket inventory;
     2. dependency-safe sequencing, compoundability ranking, blockers, and parallel lanes.
   - Preserves exact YouTrack Stage, Kanban State, resolution, GitHub delivery state, apply state, and operational-verification state.
   - Requires reconciliation of stale or conflicting environment, account, domain, CI/OIDC, and migration assumptions.
   - Keeps hard prerequisites separate from compoundability so high-leverage work cannot jump ahead of a blocker.

5. **HAWP plus five-provider adoption**
   - Consolidated installation guidance for Claude Code, Codex, Continue, Cursor, and GitHub/Copilot.
   - Requires preserving existing project instructions and reconciling the shared `AGENTS.md` surface used by Codex and Cursor.
   - Establishes `sentzunhat/human-ai-workflow-protocol@main` as the HAWP source.
   - Establishes `mode40/mast_infrastructure@development` as the mast infrastructure base and PR target; installation work must occur on a feature branch created from `origin/development`, with no direct push to the target branch.
   - Establishes the Sunshine Medical Clinic repository's `main` branch as its base and PR target.
   - Adds clinic-specific protections: patient information, PHI, credentials, and production data must not be copied into HAWP records, prompts, commits, or PR evidence.

### After

- The agent workflows are evidence-gated, read-only by default for investigation, explicit about authorization boundaries, and clear about the difference between merged, applied, and operationally verified work.
- Infrastructure discovery and prioritization are now separate phases with a verified inventory as the handoff contract.
- HAWP adoption has a single conceptual all-provider workflow with repository-specific branch policies and privacy constraints.
- **The instruction sets were created; installation into either target repository was not verified in this conversation.**

## Strategic impact

- HAWP is being used not only as a coding workflow but as a governance layer for agent review, evidence collection, prioritization, and safe repository adoption.
- The prompts reduce false-positive approval risk for Terraform changes and make destructive infrastructure actions visible before authorization.
- Separating discovery from sequencing makes infrastructure portfolio planning reproducible and allows later agents to challenge evidence without rebuilding the inventory from scratch.
- The all-provider installation pattern creates a reusable adoption path, while target-specific branch and privacy rules keep repository truth authoritative.
- The next strategic step is execution validation: run the installation instructions in each target repository, inspect the resulting diffs, and verify every provider surface without overwriting project-specific guidance.

## Continuation state

### Current state

- Source protocol: `sentzunhat/human-ai-workflow-protocol@main`.
- Provider set: Claude Code, Codex, Continue, Cursor, GitHub/Copilot.
- Mast infrastructure target: base and PR target `development`; work on a feature branch from `origin/development`.
- Sunshine Medical Clinic target: base and PR target `main`; work on a feature branch from `origin/main`.
- Terraform review, MES-2806 follow-up, weekly-hours reconstruction, and infrastructure portfolio analysis each have a defined agent workflow.
- Actual installation and repository-specific verification remain unresolved.

### Next milestone

Validate HAWP all-provider installation in one target repository end to end before treating the pattern as proven for both repositories.

### Next actions

1. Confirm the exact Sunshine Medical Clinic repository identity before any repository mutation.
2. Re-read current HAWP `main` installation/provider documentation; do not rely on the July prompt if `main` has changed.
3. Inspect target repository status, branch protections, existing agent instruction files, and local conventions.
4. Create a feature branch from the correct target base.
5. Generate an installation diff for HAWP and all five providers.
6. Reconcile `AGENTS.md` rather than overwriting target-specific Codex/Cursor instructions.
7. Run available HAWP validation plus repository-native tests and checks.
8. For the clinic repository, scan the diff for PHI, patient data, credentials, secrets, and production values before committing.
9. Open a draft PR to the correct target branch and report evidence; do not merge without explicit authorization.

### Blockers and unresolved questions

- The exact Sunshine Medical Clinic repository name was not preserved in the archived conversation context.
- It is unverified whether either target repository already contains a partial or newer HAWP installation.
- HAWP `main` may have evolved since the July instructions; current installer/provider behavior must be re-verified.
- The original reusable instruction files were conversation artifacts; their durable repository location was not established by this conversation.

### Decisions not to revisit without new evidence

- HAWP source is `main`.
- Mast infrastructure base/PR target is `development`, not HAWP's branch and not `mast_infrastructure/main`.
- Sunshine Medical Clinic base/PR target is `main`.
- Investigation workflows remain read-only unless the user separately authorizes mutation.
- A green check or merged PR is not evidence that Terraform was safely applied or operationally verified.

## Memory delta

- **Updated:** the existing HAWP workflow timeline with the July 28 agent-instruction and adoption-design milestone.
- **Merged:** Terraform evidence-gated review rules, MES-2806 deployment verification, weekly-hours reconstruction, infrastructure epic inventory/sequencing, and all-provider installation guidance.
- **Added:** canonical provider list, repository-specific branch roles, `AGENTS.md` collision handling, and clinic PHI/privacy constraints.
- **Not duplicated:** earlier HAWP implementation/review history, CI gates, symlink findings, release-candidate state, and unrelated infrastructure/GPU/financial context.

## Resume summary

**Resume from:** verify current HAWP `main` and inspect the chosen target repository before applying the consolidated five-provider installation workflow.

**Next objective:** produce a clean, validated draft PR for one target repository without overwriting project instructions or exposing sensitive data.
