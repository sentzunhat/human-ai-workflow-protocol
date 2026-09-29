# 2026-09-29 — Reusable PR review agent instruction set checkpoint

## What changed

### Before

- Prior reusable agent guidance existed for reviewing pull requests against tickets, including ticket alignment, code audit, quality signals, and review outcomes.
- Earlier variants were larger and sometimes tied to specific external work contexts.
- The durable repository location for the generalized, public-safe checkpoint was not established.

### During

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

### After

- Current reusable pattern: a compact PR-review-agent prompt with two caller-replaced placeholders and a fixed review output structure.
- The prompt is suitable for general PR audit/review workflows and can be adapted into a repo-local skill, prompt file, or issue/PR review template.
- No code or runtime behavior changed in this repository.

## Strategic impact

- This strengthens HAWP-adjacent review workflows by preserving a concise agent instruction pattern for PR review, ticket validation, and review decision discipline.
- It supports future automation or human-agent collaboration without requiring the agent to infer the ticket or PR context from memory.
- The most important design decision is that the agent must separate evidence from assumptions and must not approve when blocking correctness, security, authorization, data, migration, architecture, or quality-gate concerns remain unresolved.

## Continuation state

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

## Memory delta

- Updated durable context should record only the compact state: a reusable PR-review-agent instruction set now exists with `{{PR_REFERENCE}}` and `{{TICKET_REFERENCE}}` placeholders and approve/comment/request-changes outcomes.
- Do not duplicate the full prompt in memory.
- Keep this checkpoint as the detailed archive.

## GitHub archival

- Related repository found: Yes
- Organization: `sentzunhat`
- Repository: `sentzunhat/human-ai-workflow-protocol`
- Branch: `main`
- Checkpoint file: `.hawp/work/status/2026/09/29/checkpoint-summary.md`
- Fallback archive: No

## Resume summary

**Resume from:** The compact reusable PR-review-agent instruction set exists in conversation form and has been checkpointed here.

**Next objective:** Choose a durable prompt/skill/template location and test the instruction set against a real PR plus ticket reference.
