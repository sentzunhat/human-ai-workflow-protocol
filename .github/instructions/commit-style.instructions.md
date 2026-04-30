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

## Variants

- If the user says "split", "many small", "multiple commits", or similar,
  follow [.github/prompts/hawp-commit-many-small.prompt.md](../prompts/hawp-commit-many-small.prompt.md) instead.
