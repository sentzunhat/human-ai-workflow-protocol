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
