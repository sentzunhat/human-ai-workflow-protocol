# 2026-09-29 — Reusable agent instruction set checkpoints

Status: archived public-safe status checkpoint
Conversation/event date: 2026-09-29
Archive date: 2026-09-29
Repository: `sentzunhat/human-ai-workflow-protocol`
Path: `.hawp/work/status/2026/09/29/checkpoint-summary.md`

This same-day status file is intentionally updated rather than duplicated. It preserves compact public-safe checkpoints for reusable HAWP-adjacent agent instruction sets created on 2026-09-29.

## Checkpoint 1 — Reusable PR review agent instruction set

### What changed

#### Before

- Prior reusable agent guidance existed for reviewing pull requests against tickets, including ticket alignment, code audit, quality signals, and review outcomes.
- Earlier variants were larger and sometimes tied to specific external work contexts.
- The durable repository location for the generalized, public-safe checkpoint was not established.

#### During

- A shorter reusable instruction set was created with replaceable placeholders:
  - `{{PR_REFERENCE}}`
  - `{{TICKET_REFERENCE}}`
- The instruction set asks the agent to:
  - read the ticket first,
  - inspect the pull request diff,
  - compare implementation to acceptance criteria,
  - audit correctness, security, authorization, tests, migrations, architecture, and CI/quality signals,
  - leave only useful review comments,
  - and choose exactly one review outcome: `APPROVE`, `COMMENT`, or `REQUEST CHANGES`.
- A smaller fast-mode version was also produced for lightweight reuse.
- This checkpoint was made public-safe by preserving the reusable workflow pattern without including private customer, employer, ticket, or repository details from adjacent conversations.

#### After

- Current reusable pattern: a compact PR-review-agent prompt with two caller-replaced placeholders and a fixed review output structure.
- The prompt is suitable for general PR audit/review workflows and can be adapted into a repo-local skill, prompt file, or issue/PR review template.
- No code or runtime behavior changed in this repository.

### Strategic impact

- This strengthens HAWP-adjacent review workflows by preserving a concise agent instruction pattern for PR review, ticket validation, and review decision discipline.
- It supports future automation or human-agent collaboration without requiring the agent to infer the ticket or PR context from memory.
- The most important design decision is that the agent must separate evidence from assumptions and must not approve when blocking correctness, security, authorization, data, migration, architecture, or quality-gate concerns remain unresolved.

### Continuation state

- Current state: reusable PR-review-agent instruction text exists in conversation history and is summarized here as a checkpoint.
- Next milestone: decide whether to promote the instruction set into a repo-local reusable artifact, such as a prompt, skill, or template.
- Next actions:
  1. Choose the durable artifact format.
  2. Add the compact instruction set in that format.
  3. Test it against one real PR and ticket pair.
  4. Refine wording only where the real review exposes ambiguity.
- Blockers: no target artifact path has been selected yet.
- Unresolved questions:
  - Should this become a HAWP kit artifact, a `.github` prompt/skill, or project-specific review guidance?
  - Should the fast-mode version be kept as a separate artifact or folded into the main prompt?
- Dependencies: access to the relevant PR and ticket system when executing an actual review.

### Memory delta

- Updated durable context should record only the compact state: a reusable PR-review-agent instruction set now exists with `{{PR_REFERENCE}}` and `{{TICKET_REFERENCE}}` placeholders and approve/comment/request-changes outcomes.
- Do not duplicate the full prompt in memory.
- Keep this checkpoint as the detailed archive.

### GitHub archival

- Related repository found: Yes
- Organization: `sentzunhat`
- Repository: `sentzunhat/human-ai-workflow-protocol`
- Branch: `main`
- Checkpoint file: `.hawp/work/status/2026/09/29/checkpoint-summary.md`
- Fallback archive: No

### Resume summary

**Resume from:** The compact reusable PR-review-agent instruction set exists in conversation form and has been checkpointed here.

**Next objective:** Choose a durable prompt/skill/template location and test the instruction set against a real PR plus ticket reference.

## Checkpoint 2 — Repository audit agent instruction set

### What changed

#### Before

- HAWP needed a reusable digital-agent instruction set for reviewing and auditing repository files and listing actionable improvements.
- The desired agent behavior was read-only repository inspection, not code modification.
- No repository files were audited and no tests, linters, or build commands were run as part of the instruction-set drafting.
- It was unresolved whether the instruction set should later become a canonical HAWP prompt, checklist, template, or implementation workflow.

#### During

- A full read-only repository-audit agent instruction set was drafted.
- A compact quick-run version was also drafted.
- The full instruction set directed the agent to orient around project purpose, language, framework, package manager, entry points, commands, and folder structure.
- The audit scope covered source code, domain logic, API/routes/controllers, infrastructure/adapters, persistence, tests, docs, configuration, scripts/tooling, CI/CD, security/auth/secrets, observability, errors, and developer experience.
- The audit also included HAWP-specific review of `.hawp`, `docs`, `standards`, or workflow/protocol folders for instruction clarity, repeatability, decision records, standards, templates, review checklists, public/private documentation boundaries, and drift between docs and code.
- Findings were required to be evidence-based, concrete, and explicit about uncertainty.
- Improvements were classified by priority, type, effort, and confidence:
  - Priority: `P0`, `P1`, `P2`, `P3`
  - Type: bug risk, security, architecture, testing, documentation, developer experience, performance, maintainability, product behavior, CI/CD, observability
  - Effort: small, medium, large
  - Confidence: high, medium, low
- The required report structure included executive summary, project map, strengths, key risks, improvement backlog, phased work plan, recommended first five tasks, and open questions.

#### After

- A reusable full repository-audit prompt and a shorter quick-run prompt exist in conversation history.
- No application/source code changes were made.
- This repository was not itself audited in this conversation; the archived artifact only preserves the audit-agent instruction set and continuation state.
- The next decision is whether to promote the drafted instruction set into canonical HAWP documentation or reusable kit material.

### Completed work

- Drafted a full repository-audit digital-agent instruction set.
- Drafted a shorter quick-run repository-audit prompt.
- Established that the audit agent should be read-only by default, evidence-based, practical, and explicit about uncertainty.
- Archived the checkpoint by updating the existing same-date status file instead of creating a duplicate same-day checkpoint.

### Decisions

- Repository audit agents should not modify or delete files unless explicitly asked.
- Repository audit findings should cite concrete repository evidence and avoid inventing missing files or commands.
- Agents should not claim tests, builds, linters, or other commands pass unless they actually ran those commands.
- Improvements should be classified by priority, type, effort, and confidence.
- HAWP audit coverage should include workflow/protocol documentation, standards, decision records, public/private documentation boundaries, and drift between docs and code.

### Unresolved questions

- Where should the repository-audit prompt live as a canonical HAWP asset?
- Should the instruction set be maintained as a prompt, checklist, template, or all three?
- Should future audit-agent runs optionally open issues or propose patches after human approval?
- Should the quick-run version be a separate artifact or part of the same canonical prompt file?

### Strategic impact

- The repository-audit instruction set expands HAWP-adjacent workflows beyond PR review into full repository understanding, risk discovery, and improvement planning.
- It supports human/AI collaboration by requiring project orientation, evidence-based findings, uncertainty labeling, and an immediately actionable backlog.
- It reinforces a key HAWP operating principle: agents should preserve human agency by surfacing clear findings and next tasks rather than making broad, unapproved changes.
- The next strategic move is to decide whether this remains a conversation-level pattern or becomes maintained HAWP kit/template material.

### Continuation state

Current state:

- Full and compact repository-audit-agent prompts exist in conversation history.
- This public-safe checkpoint records the material behavior and continuation state without copying private surrounding context.
- No repo audit has been executed yet from this instruction set.
- No canonical HAWP artifact path has been chosen yet.

Next milestone:

- Promote the repository-audit-agent instruction set into a maintained HAWP prompt, checklist, template, or kit artifact if it is intended for reuse.

Next actions:

1. Inspect existing HAWP docs, kit, and provider template conventions.
2. Choose the durable artifact type and path for the repository-audit agent instructions.
3. Add the full instruction set as the maintained version.
4. Add the compact quick-run prompt as either a separate variant or a section in the same artifact.
5. Run it against one repository and refine only where the real audit exposes ambiguity.

Blockers:

- No explicit decision yet on canonical artifact type or path.
- No real repository audit has been run from the prompt, so wording has not been validated against an actual audit session.

Dependencies:

- Existing HAWP documentation and kit layout.
- Agreement on whether this belongs in HAWP core, kit templates, provider-specific prompts, or project-local guidance.

### Memory delta

- Updated durable context should record only the compact state: a reusable repository-audit-agent instruction set now exists and is read-only by default, evidence-based, repo-wide, HAWP-aware, and classifies improvements by priority, type, effort, and confidence.
- The full prompt should not be duplicated in memory.
- This checkpoint remains the detailed public-safe archive.

### GitHub archival

- Related repository found: Yes
- Organization: `sentzunhat`
- Repository: `sentzunhat/human-ai-workflow-protocol`
- Branch: `main`
- Checkpoint file: `.hawp/work/status/2026/09/29/checkpoint-summary.md`
- Fallback archive: No

### Resume summary

**Resume from:** The repository-audit-agent instruction set exists in conversation form and is checkpointed here, but has not yet been promoted into a canonical HAWP artifact.

**Next objective:** Choose a durable HAWP prompt/checklist/template location and convert the drafted repository-audit-agent instructions into a maintained reusable artifact.

## Checkpoint 3 — HAWP onboarding link drift and lightweight Markdown link validation

### What changed

#### Before

- HAWP already treated its small core shape and lean-scope guardrail as durable design constraints: improvements should prefer optional patterns, examples, evidence discipline, and workflow guidance rather than schema/runtime expansion.
- The Guardrail ADR lived in dated work-history storage, while onboarding/usage guidance depended on a path that could drift as work records moved.
- There was no lightweight repository validation dedicated to catching broken local Markdown links across the main HAWP documentation surfaces.
- Prior project history already contained documentation-link drift findings, so this was a recurring maintenance risk rather than a new category of problem.

#### During

- The repository `sentzunhat/human-ai-workflow-protocol` was identified as the project-specific public repository.
- The live Guardrail ADR content was inspected from the dated work-history record and its lean-scope decision was preserved rather than replaced with unrelated new policy.
- A stable onboarding reference was created at `core/.hawp/kit/usage/GUARDRAIL_ADR.md`.
- `core/.hawp/kit/usage/init.md` was changed to point to the stable sibling reference instead of sending onboarding readers into date-organized work history.
- A small Node script was added at `librarian/scripts/check-markdown-links.mjs` to scan local Markdown links under `.hawp`, `docs` when present, and root `README.md`.
- `librarian/package.json` was changed to expose `npm run check:markdown-links` and include it in the aggregate `npm run validate` flow.
- Four GitHub write operations returned commit SHAs during the conversation:
  - `a6065d74df431a8130dd52cab14e4975e85e9a13` — stable Guardrail ADR reference
  - `34031f1955181203d636ac5429d151b597825ca5` — onboarding Guardrail link update
  - `b28fc41cd340cd50134cd7a3b4684487ad16a0a6` — Markdown link checker
  - `137ff0c81a4f011ef816b82a4810b5f00eeeabca` — package-script integration
- At the end of that implementation turn, GitHub reported no CI statuses for the last returned SHA.
- During archival, the `dev` ref no longer resolved through the GitHub connector, and commit search did not rediscover those four commit messages/SHAs. Therefore their current reachability from a live branch is **unverified** and must not be treated as durable branch state.

#### After

- The intended fix is well-defined: onboarding should use a stable Guardrail ADR reference, and local Markdown link drift should be caught by a small repository check rather than a documentation platform.
- The implementation was accepted by GitHub write calls during the conversation, but its current presence on a live branch is unresolved because `dev` no longer resolves and the returned commits are not currently discoverable through commit search.
- The canonical same-day checkpoint remains on `main`; this entry records the implementation attempt and the verification gap without claiming that the live repository currently contains those changes.
- No private project names, personal information, or unrelated conversation context are included in this public checkpoint.

### Completed work

- Identified the stale/date-sensitive onboarding reference problem.
- Inspected the existing Guardrail ADR and preserved its lean-scope intent.
- Designed and wrote the stable ADR reference, onboarding-link update, lightweight Markdown link checker, and npm validation integration.
- Captured the exact commit SHAs returned by the GitHub write operations.
- Re-checked repository state during archival and detected that the earlier `dev` state can no longer be verified.

### Decisions

- Onboarding documentation should link to stable kit/reference paths, not date-based work-history paths.
- The dated ADR remains historical evidence; the stable onboarding reference is the durable navigation target.
- Link validation should stay lightweight: local Markdown links in `.hawp`, `docs`, and `README.md` are sufficient scope.
- Do not build a documentation platform for this concern.
- Do not claim CI/link validation passed unless it actually runs successfully.
- Do not assume the earlier `dev` writes remain live merely because the write API returned commit SHAs.

### Strategic impact

- This closes a recurring class of onboarding failure at the design level: stable public references are separated from dated work-history organization.
- Adding a small link check to normal validation would turn documentation drift from a release-time/manual discovery into an ordinary quality-gate failure.
- The approach remains consistent with HAWP's guardrail: solve the concrete reliability problem with a small optional/tooling aid rather than expanding the protocol schema or creating documentation infrastructure.
- The immediate priority is repository-state reconciliation, not additional documentation features.

### Continuation state

Current state:

- Canonical archive branch: `main`.
- The intended link-drift patch is fully specified in this checkpoint.
- Earlier GitHub writes targeted `dev` and returned four commit SHAs, but `dev` is currently unresolved through the connector and those commits are not currently rediscovered by commit search.
- CI success for the patch is not established.

Next milestone:

- Reconcile the intended patch with the repository's current live development branch, then run the Markdown link check and the relevant validation suite successfully.

Next actions:

1. Determine the current development branch/ref replacing or superseding `dev`, if any.
2. Inspect whether `core/.hawp/kit/usage/GUARDRAIL_ADR.md`, the `init.md` stable link, `librarian/scripts/check-markdown-links.mjs`, and the package scripts already exist on that live branch.
3. Reapply only missing pieces; do not duplicate files or overwrite newer equivalent work.
4. Run `npm run check:markdown-links` from `librarian`.
5. Run the repository's relevant aggregate validation/CI path and record evidence.
6. Only after successful verification, treat the link-drift fix as completed durable branch state.

Blockers:

- The previously targeted `dev` ref does not currently resolve through GitHub.
- The four returned implementation commits are not currently discoverable through repository commit search.
- No successful CI/link-check execution was captured in this conversation.

Unresolved questions:

- Was `dev` deleted, renamed, force-moved, or otherwise superseded after the write operations?
- Are the four returned commits reachable through another branch/ref even though commit search does not currently surface them?
- Does a newer branch already contain an equivalent stable Guardrail ADR/link-check implementation?

Dependencies:

- Current GitHub branch topology.
- Existing HAWP distribution/source synchronization rules if the stable kit file must also be materialized into repo-local/generated copies.
- Node/npm versions required by `librarian/package.json` for validation.

### Memory delta

- Merge into existing HAWP durable context rather than creating a parallel project memory.
- Durable addition: on 2026-09-29, HAWP onboarding link drift was addressed conceptually and via GitHub write operations by introducing a stable Guardrail ADR reference and a lightweight local Markdown-link checker integrated with librarian validation.
- Durable caution: the implementation commits returned by those writes are not currently verified as reachable because the `dev` ref stopped resolving during archival; future work must reconcile live branch state before considering the fix complete.
- Preserve the existing HAWP lean-scope guardrail and prior v0.0.24/review history; do not duplicate those timelines here.
- Do not preserve unrelated personal, business, infrastructure, finance, or private-project context in this open-source checkpoint.

### GitHub archival

- Related repository found: Yes
- Organization: `sentzunhat`
- Repository: `sentzunhat/human-ai-workflow-protocol`
- Archive branch: `main`
- Implementation branch requested earlier: `dev` (currently unresolved through GitHub during archival)
- Checkpoint file: `.hawp/work/status/2026/09/29/checkpoint-summary.md`
- Fallback archive: No

### Resume summary

**Resume from:** Reconcile the four intended onboarding/link-check changes against the repository's current live development branch; do not assume the earlier `dev` writes are still reachable.

**Next objective:** Establish a live, verified stable Guardrail ADR link plus lightweight Markdown-link validation, then run and record successful validation evidence.



## Checkpoint 3 — CLI-driven iterative participant loop direction

### What changed

#### Before

- HAWP's durable model remains a lightweight task-shaping and workflow method; repository guidance explicitly says it is not a runtime engine or orchestrator.
- HAWP already supports iterative human/AI workflow stages, durable handoffs, provider guidance, MCP/CLI surfaces, and review/verification checkpoints.
- CLI participant adapters for Codex, Claude, and GitHub were already parked as work item `bee15107` because provider packs were sufficient at that stage.
- No HAWP runtime loop for repeatedly invoking coding agents was implemented.

#### During

- The conversation explored the "Ralph loop" pattern: repeatedly select/shape work, invoke an AI coding participant, inspect/evaluate the result, update state, and continue until a completion or stop condition is reached.
- Ralph TUI was discussed as an existing interface/orchestration implementation, but the preferred integration direction for HAWP is CLI-first rather than depending on a TUI.
- Cursor and GitHub Copilot were identified as initial examples of coding participants that could be invoked through CLI-oriented adapters/prompts; other agents or a human could occupy equivalent workflow roles.
- A key abstraction emerged: HAWP should model the role/participant and transition semantics independently from any specific AI vendor or editor.
- The idea remains architectural exploration only. No loop engine, participant interface, retry policy, CLI adapter, or runtime orchestration was implemented in this conversation.

#### After

- Current direction: investigate a generic, CLI-driven iterative workflow capability that can hand a work unit between AI agents and humans while preserving HAWP's existing shaping, evidence, review, verification, and continuation semantics.
- The capability should not be described as implemented and should not silently change HAWP's current "not an orchestrator/runtime" contract.
- Any implementation should begin as an explicit proposal/work item that resolves whether orchestration belongs inside HAWP core, beside HAWP as a companion runner, or as an adapter layer.
- "Ralph loop" is useful prior-art terminology, but the durable HAWP concept should remain vendor-neutral and not depend on Ralph TUI.

### Completed work

- Identified the iterative CLI-agent loop as a potentially useful direction for larger tasks/projects.
- Established CLI-first integration as preferable to coupling the design to a TUI.
- Identified humans and AI agents as interchangeable participants at workflow boundaries, subject to explicit role/capability and approval rules.
- Connected the new idea to existing HAWP participant/provider, handoff, review, verification, and continuation concepts.

### Decisions

- Prefer CLI-oriented participant integration over editor/TUI-specific coupling.
- Keep participant identity vendor-neutral at the workflow level; Cursor, GitHub Copilot, Codex, Claude, or a human are adapters/participants rather than the workflow definition itself.
- Preserve explicit evaluation/review/verification between iterations instead of treating repetition as an uncontrolled infinite prompt loop.
- Do not claim HAWP is now an orchestrator. Runtime/orchestration ownership remains unresolved and requires an explicit architecture decision.

### Experiments / hypotheses

- Hypothesis: a small participant adapter contract plus a state/transition loop could let HAWP continue larger work across multiple agent invocations without losing durable workflow state.
- Hypothesis: existing HAWP handoffs and work records can provide the continuation packet between iterations, reducing dependence on a single model's conversation context.
- Hypothesis: the parked CLI-participant work item `bee15107` may be relevant prior work, but it should be re-evaluated rather than automatically reopened because the proposed loop has broader orchestration semantics.

### Unresolved questions

- Should the loop runner live inside HAWP, as a separate companion CLI, or as an external orchestrator consuming HAWP artifacts?
- What is the minimal participant contract: prompt in/result out, structured result, capabilities, status, evidence, cost/usage, cancellation, and/or resume token?
- Which transitions require human approval, and which may automatically continue?
- What are the termination controls: success criteria, max iterations, time/cost budget, repeated-failure detection, cancellation, and manual stop?
- How should retry/reflection differ from handing the task to a different participant?
- How should Cursor and GitHub Copilot CLI capabilities be detected and invoked without making HAWP depend on proprietary editor behavior?
- Should `bee15107` be superseded, extended, or remain parked?

### Strategic impact

- This direction could extend HAWP from shaping and durable continuation toward optional execution continuity for larger tasks while retaining human review gates.
- The architectural risk is scope creep: embedding orchestration directly into HAWP could contradict its established lightweight protocol boundary. The next milestone therefore must be an architecture decision/prototype boundary, not immediate broad implementation.
- Existing handoff, evidence, verification, MCP/CLI, provider, and work-record primitives should be reused rather than creating a parallel state model.

### Continuation state

Current state:

- HAWP v0.0.24 remains active in PR #41 on `feature/v0.0.24`; this checkpoint does not alter that release implementation.
- CLI-driven iterative execution is a proposed post/current-release design direction, not completed functionality.
- Existing parked item `bee15107` is the closest known prior work concerning CLI participant adapters.

Next milestone:

- Write a focused architecture/work-intake proposal for an optional CLI iteration runner and decide its boundary relative to HAWP core.

Next actions:

1. Inspect `.hawp/work/parked/bee15107/plan.md` and current CLI/provider contracts.
2. Define a minimal participant interface and iteration state machine on paper before implementation.
3. Define mandatory stop/budget/approval/verification semantics.
4. Decide whether the runner is HAWP core, a companion package/CLI, or an external consumer.
5. Prototype one adapter path with a single CLI coding agent before adding multiple providers.
6. Validate that handoffs/work records can resume an interrupted iteration without relying on hidden conversation state.

Blockers:

- Runtime ownership/boundary is unresolved.
- Exact current CLI contracts for Cursor and GitHub Copilot have not been verified as part of this checkpoint.
- No acceptance criteria or work item has yet been approved for implementation.

Dependencies:

- Existing HAWP work/handoff/evidence model.
- Current HAWP CLI and provider architecture.
- The parked CLI participant adapter investigation `bee15107`.
- Provider CLI capabilities and stable invocation contracts.

### Memory delta

- Add one durable change: on 2026-09-29, HAWP gained a proposed CLI-first iterative participant-loop direction for continuing larger tasks across AI agents or humans.
- Preserve that this is architectural exploration only, not implemented runtime behavior.
- Preserve the boundary decision: vendor-neutral participant semantics, explicit review/verification/stop controls, and an unresolved choice between HAWP core versus a companion/external runner.
- Do not duplicate existing v0.0.24, Librarian, handoff, provider, or review-protocol history.

### GitHub archival

- Related repository found: Yes
- Organization: `sentzunhat`
- Repository: `sentzunhat/human-ai-workflow-protocol`
- Source branch inspected: `feature/v0.0.24` at `9cbf94864f14370e0af998f892f917971590a137`
- Archival branch: `checkpoint/2026-09-29-cli-iteration-loop`
- Checkpoint file: `.hawp/work/status/2026/09/29/checkpoint-summary.md`
- Fallback archive: No
- Note: archival was isolated on a checkpoint branch so the open v0.0.24 PR was not modified.

### Resume summary

**Resume from:** Treat the CLI-driven participant loop as an unimplemented architecture proposal; start by reviewing `bee15107` and defining the minimal participant/state-machine boundary.

**Next objective:** Produce and review a focused architecture/work-intake proposal for a bounded CLI iteration runner with explicit human/agent handoff, verification, and termination semantics.

## Checkpoint 3 — GitHub-enabled master-builder agent setup request

### What changed

#### Before

- HAWP already had public reusable agent-instruction checkpoints for PR review and repository audit patterns.
- The repository's active work index reported no active work items, with open follow-up ideas parked rather than in progress.
- The repository supports provider-oriented install/update distribution, including GitHub Copilot, Codex, Cursor, Claude Code, and Continue.
- A browser view suggested a `dev` branch, but the GitHub API branch listing exposed `development` as the available development branch. The `dev` ref was not available through the API.
- No dedicated implementation branch had yet been created for a GitHub agent artifact in this session.

#### During

- GitHub connector access was verified for `sentzunhat/human-ai-workflow-protocol` with write/admin-capable repository permissions.
- Repository discovery checked accessible `beltrd` and `sentzunhat` repositories and confirmed `sentzunhat/human-ai-workflow-protocol` as the directly associated repository for this checkpoint.
- Existing HAWP context was reviewed, including the README, active backlog, branch list, and existing 2026-09-29 checkpoint summary.
- A safe archival branch was created from `development`: `checkpoint/2026-09-29-github-agent-access`.
- This existing same-day public-safe status file was updated instead of creating a parallel checkpoint.

#### After

- GitHub access is available for future HAWP repository work through the connector.
- The checkpoint branch records the current handoff state without changing application/runtime code.
- The requested GitHub agent/master-builder artifact has not yet been implemented.
- The next implementation should start from `development`, not an API ref named `dev`, unless a `dev` branch is later created or clarified.

### Completed work

- Verified GitHub access to the HAWP repository.
- Confirmed the directly related repository is `sentzunhat/human-ai-workflow-protocol`.
- Confirmed the current development branch available through the GitHub API is `development`.
- Created archival branch `checkpoint/2026-09-29-github-agent-access` from `development`.
- Updated the existing 2026-09-29 status checkpoint with this continuation state.

### Decisions

- Public checkpoint content must remain public-safe and must not include private personal background, private project context, or private organization details.
- The future GitHub agent should be grounded in repository conventions and HAWP patterns rather than inferred personal biography.
- No implementation work should be merged automatically from this checkpoint branch.
- The future implementation branch should be created from `development` unless the maintainer explicitly asks for another base branch.

### Unresolved questions

- What should the implementation branch be named?
- Should the GitHub agent live as a HAWP kit artifact, a `.github` prompt, a provider-specific rule set, or a generated distribution artifact?
- What exact tool boundaries should be granted to the GitHub agent beyond repository search/read/write, issue/PR work, CI inspection, and optional vision-capable review when images are part of an issue or PR?
- Should the agent be optimized for GitHub Copilot agent mode, Codex, Claude Code, or a provider-neutral HAWP template first?
- What public-safe wording should replace the request to base the agent on the maintainer's personal style?

### Strategic impact

- This moves the HAWP GitHub-agent effort from an access/setup question into a ready-to-start implementation state.
- The work should build on the existing reusable PR-review and repository-audit agent instruction sets instead of starting from scratch.
- The immediate priority is to create a canonical public-safe agent artifact that expresses master-builder behavior as repeatable engineering principles: architecture awareness, evidence-first reading, safe tool use, issue/PR discipline, testing discipline, and clear handoff reporting.

### Continuation state

Current state:

- Repository access is available.
- `development` is the verified development branch through the GitHub API.
- A checkpoint-only branch exists: `checkpoint/2026-09-29-github-agent-access`.
- No agent artifact has been added yet.
- No PR has been opened yet.

Next milestone:

- Create the implementation branch from `development` and add the first public-safe GitHub agent artifact.

Next actions:

1. Choose or confirm the implementation branch name.
2. Inspect provider/rules distribution paths before adding the artifact.
3. Decide whether the first artifact belongs under `.github`, `core/.hawp/kit`, provider rules, or distribution sources.
4. Draft the GitHub master-builder agent as a public-safe reusable instruction set.
5. Validate that only documentation/prompt files changed.
6. Open a PR only after human approval.

Blockers:

- Implementation branch name is not yet confirmed.
- Artifact path and provider target are not yet decided.
- The `dev` branch name from the request does not currently resolve through the GitHub API; `development` is the usable base branch.

Dependencies:

- Existing HAWP provider distribution layout.
- Existing PR-review and repository-audit agent instruction-set checkpoints.
- Maintainer approval before any PR, merge, or broader code changes.

### Memory delta

- Durable memory should record only the compact state: GitHub connector access to `sentzunhat/human-ai-workflow-protocol` is working; the usable development branch is `development`; an archival checkpoint branch was created; the GitHub master-builder agent artifact remains unimplemented and should start from `development` after branch/path confirmation.
- Do not duplicate private personal context or the full checkpoint in memory.
- This checkpoint remains the detailed public-safe archive.

### GitHub archival

- Related repository found: Yes
- Organization: `sentzunhat`
- Repository: `sentzunhat/human-ai-workflow-protocol`
- Branch: `checkpoint/2026-09-29-github-agent-access`
- Checkpoint file: `.hawp/work/status/2026/09/29/checkpoint-summary.md`
- Fallback archive: No

### Resume summary

**Resume from:** GitHub access is verified and the HAWP repository checkpoint is archived on `checkpoint/2026-09-29-github-agent-access`; no GitHub master-builder agent artifact has been implemented yet.

**Next objective:** Create the implementation branch from `development`, choose the canonical public-safe artifact path, and draft the first GitHub master-builder agent instruction set.

## Checkpoint 4 — v0.0.24 architecture verification and AI collaboration attribution

- Status: archived public-safe continuation checkpoint
- Conversation/event dates: 2026-09-18 through 2026-09-29
- Archive date: 2026-09-29
- Repository: `sentzunhat/human-ai-workflow-protocol`
- Source branch inspected: `feature/v0.0.24` at `4d39675cd0b73f765b9594c5a168dcf53307afd1`
- Pull request: #41, `feature/v0.0.24` → `development`

### What changed

#### Before

- The v0.0.24 work had previously been inspected on `feature/v0.0.24-work-folder-normalization`.
- The architecture diagram under discussion treated `commands.go` and `registry.go` as the CLI execution path, collapsed indexing into one service, reversed the context/search dependency, and referenced several stale infrastructure paths.
- The checked-in `.hawp/bin/hawp` provenance was initially ambiguous because the binary stored in the `0.0.23` tag differed from the published `0.0.23` release asset.
- HAWP's canonical commit convention normally used one lowercase plain sentence with no body or trailers. AI co-author attribution existed in history for individual tools/models, but there was no provider-family collaboration convention.

#### During

Architecture and binary verification established the following:

- `cmd/hawp/main.go` delegates into the CLI composition/router in `internal/platform/cli/run.go`; `commands.go` and `registry.go` primarily describe command discovery/help metadata rather than dispatch application work.
- Indexing is a staged flow:
  - corpus build/enrichment in `index/build-service.go`,
  - chunk/index ingestion in `index/ingest-service.go`,
  - and a separate embedding pass in `index/embed-service.go`.
- Bounded context/RAG retrieves through the shared search service; the dependency direction is context → search, not search → context.
- Important current paths at the time of the audit included:
  - SQLite index repository: `internal/infrastructure/repositories/index/index.go`
  - MCP server: `internal/platform/mcp/server/server.go`
  - GitHub release client: `internal/infrastructure/clients/githubrelease/githubrelease.go`
  - project layout: `internal/infrastructure/filesystem/hawp_project.go`
- Provisioning and kit synchronization are distinct responsibilities: provisioning handles runtime/model setup, while kit sync obtains and materializes the release kit/provider bundle.
- The corrected architecture diagram produced from this audit is a useful historical reference, but it was not made canonical repository documentation in this conversation and must be revalidated before publication because PR #41 subsequently changed a large portion of the branch.

Binary provenance work also exposed an important release lesson:

- The published `hawp-darwin-arm64` v0.0.23 release asset was 21,151,570 bytes and corresponded to the historical verified branch blob `59348c5ee9d5a5fa7380324d5b5c6f4d525b1459`.
- The binary checked into the `0.0.23` tag was a different blob, `e2fae97f520dc94a7d49e7bec31da302498bce8d`, with a different size.
- A tag-binary substitution was briefly committed as `57fdc971ecca5966e5f4f21e538c2a29a3ed5414`, then corrected after the mismatch was discovered.
- Commit `287c3e0f7ec8d814fe4bcfbcad125e74234ae87e` restored the historically verified v0.0.23 release binary.
- Durable lesson: do not treat a tag's repository copy of a binary as equivalent to the published release asset without verifying provenance.

AI collaboration attribution was then standardized at the conversation level:

- Prefer stable provider/family labels rather than listing every transient model variant:
  - OpenAI GPT
  - Anthropic Claude
  - GitHub Copilot
  - Qwen
- Model-specific names such as Luna, Astra, individual GPT variants, and Sonnet point releases need not become separate collaborator identities.
- The selected human-readable collaborator block is:

  - `OpenAI GPT`
  - `Anthropic Claude`
  - `GitHub Copilot`
  - `Qwen`

- The selected raw Git trailer form is:

  ```text
  Co-authored-by: OpenAI GPT <noreply@openai.com>
  Co-authored-by: Claude Sonnet 4.6 <noreply@anthropic.com>
  Co-authored-by: Copilot <175728472+Copilot@users.noreply.github.com>
  Co-authored-by: Qwen <UNVERIFIED-QWEN-EMAIL>
  ```

- These trailers are an explicit exception to HAWP's default no-trailer commit rule when AI attribution is intentionally requested.
- `Qwen <UNVERIFIED-QWEN-EMAIL>` remains a placeholder, not a verified GitHub identity.
- `OpenAI GPT <noreply@openai.com>` is the selected convention from the conversation, but its account-association behavior was not independently verified here.
- No repository commit-style rule was changed to make this automatic; the decision currently exists as an explicit collaboration convention for commits where it is intentionally used.

Repository state changed after the earlier branch audit:

- PR #40 merged `feature/v0.0.24-work-folder-normalization` into `feature/v0.0.24` on 2026-09-21.
- The active release work is now PR #41 from `feature/v0.0.24` into `development`.
- At archive time, PR #41 head is `4d39675cd0b73f765b9594c5a168dcf53307afd1`.
- The head's `Quality` and `Validate Distribution Generated` workflow runs are green.
- The live PR review state has 80 review threads with 2 unresolved Copilot findings. This supersedes an older PR-description snapshot that reported zero unresolved threads.
- Both unresolved findings concern explicit repository roots accepting symlinked ancestor components:
  - `librarian/src/internal/platform/cli/mcp/commands.go`: explicit `--repo-root` validation starts too low in the path and can trust a symlinked parent.
  - `librarian/src/internal/platform/cli/mcp/configure/command.go`: explicit configuration roots can likewise follow a symlinked parent before provider config writes.
- The latest review summary therefore still recommends changes; PR #41 is open and must not be treated as approved or complete.

#### After

- The authoritative v0.0.24 continuation point is PR #41 on `feature/v0.0.24`, not the older work-folder-normalization branch.
- The earlier architecture corrections remain useful, but any canonical diagram should be regenerated or revalidated against the current PR #41 tree before publication.
- The v0.0.23 binary provenance question is resolved for the earlier branch state: the historically verified published-release binary was restored after distinguishing it from the tag's checked-in copy.
- The AI collaboration naming decision is provider-family based: GPT, Claude, Copilot, and Qwen.
- Machine-readable Git trailers should use raw `Name <email>` syntax, not Markdown links.
- Qwen's canonical co-author email remains unresolved, and the selected generic GPT email has not been independently proven to map to a GitHub co-author identity.
- No application code, release tag, merge, or publication action is performed by this checkpoint.

### Strategic impact

- PR #41 review closure is the immediate release priority. Green CI is necessary but not sufficient while the two live filesystem-boundary findings remain unresolved.
- The two findings fit the broader v0.0.24 hardening theme: an explicit root should not become trusted until symlink ancestry is validated from an appropriate volume/root boundary.
- Architecture documentation should follow actual composition boundaries rather than file names that merely sound like dispatch or storage layers.
- Release verification should distinguish repository-tag artifacts from published release assets and preserve provenance evidence before replacing checked-in binaries.
- Provider-family AI attribution is more durable than model-version attribution and avoids churn as individual models change.
- Formalizing AI attribution in repository guidance is optional future work; the current decision does not require changing the default commit style for ordinary commits.

### Continuation state

Current state:

- PR #41 is open and mergeable, with green `Quality` and `Validate Distribution Generated` runs at `4d39675`.
- Two unresolved Copilot threads remain, both involving explicit-root symlink ancestry in MCP CLI/configuration paths.
- The architecture audit and binary provenance investigation are historical inputs, not a substitute for validating the current PR head.
- The four-family AI collaborator convention is decided for intentional attribution, but Qwen's trailer email remains unresolved.

Next milestone:

- Close the two live PR #41 filesystem-boundary findings and obtain a fresh review state with no unresolved blocking findings.

Next actions:

1. Apply the same volume/root-anchored symlink-ancestor validation already used by the work-new mutation boundary to explicit MCP repository/configuration roots.
2. Add or extend regression tests covering symlinked parent components for both affected paths.
3. Rerun the repository's required v0.0.24 validation, including `Quality`, generated-distribution validation, and the relevant Go test/check commands.
4. Request or wait for a fresh Copilot review and verify the live thread list rather than relying on the PR description's older zero-thread snapshot.
5. Only after review closure, decide whether PR #41 is ready to merge into `development`.
6. If the architecture diagram is to become canonical documentation, re-audit it against the then-current head before committing it.
7. If AI attribution is to become a repository-wide convention, verify a stable Qwen identity and decide whether the selected GPT email should remain informational or machine-readable.

Blockers:

- Two unresolved PR #41 symlink-ancestor review findings.
- No verified Qwen co-author email.
- The generic GPT trailer address has not been independently verified as a GitHub identity.

Important dependencies:

- Existing filesystem/mutation-boundary symlink guards and their regression tests.
- PR #41's current MCP CLI/configuration implementation.
- Current HAWP commit-style guidance.
- GitHub's co-author association rules for trailer email identities.

Decisions that should not be revisited without new evidence:

- Treat `run.go` as the CLI composition/dispatch point; do not model `commands.go` or `registry.go` as the application execution engine.
- Keep build/ingest/embed as distinct indexing responsibilities.
- Model bounded context as a consumer of search.
- Distinguish provisioning from kit synchronization.
- Distinguish published release assets from same-tag repository binaries when verifying provenance.
- Prefer GPT / Claude / Copilot / Qwen provider-family labels over enumerating every model version for collaboration attribution.
- Use raw Git trailer syntax, not Markdown mail links, for machine-readable co-authorship.

### Memory delta

- Reviewed the existing HAWP durable context, prior v0.0.24 architecture/binary checkpoint state, the same-day checkpoint timeline, and the live PR #41 state.
- Merged the older work-folder-normalization snapshot into the current timeline instead of preserving it as the active branch.
- Added the provider-family AI collaboration decision and the selected trailer format, while preserving the unverified status of Qwen's email and the unverified GitHub mapping of the generic GPT address.
- Added the live 2026-09-29 PR #41 continuation state: head `4d39675`, green required workflows, and two unresolved explicit-root symlink-ancestor findings.
- Preserved the binary-provenance correction and the lesson that tag-tree binaries and published release assets can differ.
- Intentionally did not duplicate the broader HAWP origin timeline, historical benchmark values, general Librarian history, or unrelated project/private context already archived elsewhere.
- This checkpoint is the detailed durable archive; future compact memory should retain only the current branch/PR continuation state, the architecture corrections that affect future diagrams, the release-provenance lesson, and the AI attribution decision.

### GitHub archival

- Related repository found: Yes
- Organization: `sentzunhat`
- Repository: `sentzunhat/human-ai-workflow-protocol`
- Source branch: `feature/v0.0.24`
- Source head: `4d39675cd0b73f765b9594c5a168dcf53307afd1`
- Archival branch: `checkpoint/2026-09-29-v0024-ai-attribution`
- Checkpoint file: `.hawp/work/status/2026/09/29/checkpoint-summary.md`
- Fallback archive: No
- Active release branch modified: No

### Historical resume summary (superseded by Checkpoint 5)

**Resume from:** PR #41 at `4d39675` with the two explicit-root symlink-ancestor review findings still unresolved in MCP CLI/configuration code.

**Next objective:** Fix and test those two boundary issues, obtain a fresh clean review state, then evaluate PR #41 for merge into `development`.

## Checkpoint 5 — PR #41 review and archive/download boundary follow-up

**Event date:** 2026-09-30

**Archive date:** 2026-09-30

### Previous state

- Checkpoint 4 captured PR #41 at `4d39675`, with two unresolved explicit-root symlink findings and the related implementation work still pending.
- The previous continuation note at `ea7d4ffc` was already stale; the live branch and review state were re-read before this update.

### Changes and verified current state

- The current PR tree includes fixes for bounded UTF-8 context truncation overflow, SQLite main/WAL/SHM/rollback-journal link protection, corrected successful-query token accounting, volume-anchored symlink-ancestor checks for explicit MCP roots, and archive/download destination revalidation before temporary-file creation.
- The latest implementation commit is `646529bd56cda2e6be3ccd9c279ca146b8e5de84` (`fix: revalidate archive and download destinations before temp creation`). The preceding checkpoint whitespace correction is `4061c411c4c0c0e5447289c81cf7e46d0bbe3990`.
- Local verification recorded for `646529bd`: Go 1.26.4 `make check`, `make install`, source-layout tests, generated-distribution sync with no drift, `hawp check`, `hawp work validate` with zero issues and one existing verification-clarity warning, formatting, and branch hygiene passed.
- Live GitHub readback: PR #41 is open against `development` at `646529bd`; `quality` and `validate-generated` pass; merge state is `CLEAN`; review decision is empty. The PR has 82 review threads and zero unresolved. The latest Copilot overview reports no findings and resolves the two archive/download threads.
- The 2026-09-29 Quality failure after `dd56b05` was a checkpoint whitespace/diff-hygiene issue; `4061c411` corrected it and the subsequent required checks passed.

### Decisions and limits

- Keep the PR open until human review/approval and an explicit merge decision. Green checks, zero unresolved threads, and a clean merge state do not constitute human approval.
- No merge or provider setup was performed. The latest filesystem checks reduce path substitution risk but do not claim elimination of every concurrent filesystem race.
- Checkpoint 4 remains the historical record for architecture and attribution context; this section supersedes only its PR continuation state.

### Resume state

- Repository: `sentzunhat/human-ai-workflow-protocol`
- Branch: `feature/v0.0.24`, tracking `origin/feature/v0.0.24`
- PR: [#41](https://github.com/sentzunhat/human-ai-workflow-protocol/pull/41), open against `development`
- Latest implementation head: `646529bd56cda2e6be3ccd9c279ca146b8e5de84`
- Next action: obtain human review/approval, decide whether to merge PR #41, then resume deferred provider setup only after the merge gate.

### Memory delta

- Supersedes the older PR #41 memory continuation with current head, checks, review-thread state, and the archive/download follow-up.
- Retains the existing detailed review history here; does not duplicate the PR transcript or unrelated architecture and attribution history.
- Memory update: `2026-09-30T15-27-31Z-hawp-pr41-review-continuation.md`.

### GitHub archival

- Related repository found: Yes
- Organization/repository: `sentzunhat/human-ai-workflow-protocol`
- Branch: `feature/v0.0.24`
- Checkpoint file: `.hawp/work/status/2026/09/29/checkpoint-summary.md`
- Source implementation head at archival: `646529bd56cda2e6be3ccd9c279ca146b8e5de84`
- Checkpoint archive commit: `887fb9c6290555e83b2f97c3122049c18e4e3917`
- Push status: checkpoint and archival metadata commits pushed to `origin/feature/v0.0.24`.

### Resume summary

**Resume from:** PR #41 at `646529bd`; checks pass, the latest Copilot review has no findings, and all review threads are resolved, but human approval is still absent.

**Next objective:** Obtain human review/approval and decide whether to merge PR #41 into `development`; only after that gate, resume provider setup.

## Checkpoint 6 — PR #41 archival and live-state reconciliation

**Event date:** 2026-09-30

### Reconciled state

- The September 30 timeline was archived at `887fb9c6290555e83b2f97c3122049c18e4e3917`; follow-up archive bookkeeping commits are `7e6828dc` and `fd703bae`. The branch and `origin/feature/v0.0.24` both currently point to `fd703bae367d49350e720dbf651b49ed5cf4da93`; the checkout was clean when checked.
- Live GitHub readback: PR #41 is open against `development` at `fd703bae`; `quality` and `validate-generated` pass; merge state is `CLEAN`; review decision is empty; 82 review threads exist and zero are unresolved.
- The latest Copilot overview is dated 2026-09-30 02:47:56 UTC and is attached to implementation commit `646529bd`. It reports no findings and says two archive/download findings were resolved. The subsequent commits `887fb9c`, `7e6828d`, and `fd703bae` only record/correct checkpoint archival details.
- No human approval or merge is recorded. The supplied page snapshot's “Still in progress?” label is not evidence of a new review request or approval; no new human review request was found in the live PR metadata.

### Resume state

- Resume from PR #41 at `fd703bae367d49350e720dbf651b49ed5cf4da93`; implementation review evidence remains tied to `646529bd`.
- Next action remains human review/approval and an explicit decision whether to merge PR #41. Resume deferred provider setup only after that gate.
