---
applyTo: "**"
description: Apply the repo commit style whenever the user asks to commit changes
---

# Commit Style — Ambient Rules

Apply these rules whenever the user asks to commit, e.g.:

- "commit the changes"
- "commit this"
- "please commit"
- "stage and commit"
- any similar phrasing requesting a Git commit

## Default Behavior

Unless the user explicitly asks for multiple commits or a multi-line body,
follow the single-commit flow defined in
[.github/prompts/hawp-commit-one-big.prompt.md](../prompts/hawp-commit-one-big.prompt.md).

## Commit Message Rules

The commit message sentence must:

- start with a lowercase letter
- use a present- or past-tense verb
- omit Conventional-Commit prefixes (`feat:`, `fix:`, `chore:`, etc.)
- avoid implementation details unless essential
- clearly describe what changed

Do not add a multi-line body unless the user explicitly asks for one.

## Method Selection

### Trigger Rules

**Use one-big (default):**

- User says "commit the changes" (no qualifier)
- User says "commit this"
- No explicit split/multiple mention

**Use many-small:**

- User says "split", "many small", or "multiple commits"
- User says "small commits"
- User says "logically coherent chunks"
- User says "separate commits"

## Workflow References

| User Request                       | Method     | Prompt File                      |
| ---------------------------------- | ---------- | -------------------------------- |
| "commit the changes"               | one-big    | hawp-commit-one-big.prompt.md    |
| "split into small commits"         | many-small | hawp-commit-many-small.prompt.md |
| "commit this" / "stage and commit" | one-big    | hawp-commit-one-big.prompt.md    |
| "multiple commits" / "break it up" | many-small | hawp-commit-many-small.prompt.md |

## Implementation Notes

Both workflows use the same commit message rules. The difference is in _how many_ commits are created and _how_ changes are grouped.

- **one-big:** All uncommitted changes → one atomic commit
- **many-small:** All uncommitted changes → split into focused logical chunks → one commit per chunk

For ambiguous requests, default to one-big. Clarify with the user if needed.
