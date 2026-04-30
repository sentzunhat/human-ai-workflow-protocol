---
name: commit-many-small
description: Split uncommitted changes into multiple small, logically coherent commits using the repo commit style
---

Review the current uncommitted changes and split them into small, logically coherent chunks of work. Each chunk should represent a single, focused change with a clear purpose.

Create one separate commit per chunk.

For each commit, write a single concise sentence that:

- starts with a lowercase letter
- uses a present or past-tense verb
- does not include prefixes such as `feat:`, `fix:`, `chore:`, etc.
- avoids implementation details unless essential
- clearly describes what changed

Examples:

- "implemented a feature flag for potato peeling drill speed control."
- "refactored the image upload process to improve performance."

Do not include multi-line bodies unless the user explicitly asks for one.
