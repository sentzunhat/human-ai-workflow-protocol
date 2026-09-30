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

