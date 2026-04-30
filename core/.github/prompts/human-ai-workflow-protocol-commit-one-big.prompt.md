---
name: commit-one-big
description: Stage all uncommitted changes and create a single commit using the repo commit style
---

Review all uncommitted changes in the repository and create a single Git commit.

Generate one concise commit message sentence that:

- starts with a lowercase letter
- uses a present or past-tense verb
- does not include prefixes such as `feat:`, `fix:`, `chore:`, etc.
- avoids implementation details unless essential
- clearly describes what changed

Then stage all changes and commit them with that message.

Do not include a multi-line body unless the user explicitly asks for one.
