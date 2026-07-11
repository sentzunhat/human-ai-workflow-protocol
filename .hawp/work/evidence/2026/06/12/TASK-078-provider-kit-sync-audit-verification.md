# TASK-078 — Provider/kit sync audit verification

Date: 2026-06-12
Repo root: `<repo-root-abs>` (redacted)

## Direct evidence

### Kit sync (core vs root)

```
$ diff -rq core/.hawp/kit .hawp/kit
(no output — identical)
```

Top level of both `core/.hawp/kit/` and `.hawp/kit/` contains only `README.md`, `start-here.md`, and category folders (`examples`, `instructions`, `patterns`, `references`, `reviews`, `standards`, `templates`, `types`, `usage`). `spec.md` and `authoring-patterns.md` live in `references/`.

### Provider packs vs root dot folders

```
$ diff -rq core/providers/.cursor/rules .cursor/rules        # identical
$ diff -rq core/providers/.continue/rules .continue/rules    # identical
$ diff -rq core/providers/.github/instructions .github/instructions   # identical (after code-style fix)
$ diff -rq core/providers/.github/prompts .github/prompts    # identical
$ diff -q core/providers/.github/copilot-instructions.md .github/copilot-instructions.md  # identical
$ diff -q core/providers/.cursor/AGENTS.md.seed AGENTS.md    # identical
```

### Librarian validators and tests (Node 22 per .nvmrc)

```
provider validation passed: 11 materialized file(s) are current
distribution validation passed: generated outputs are current
workflow validation: VALIDATION PASS (0 issues, 1 pre-existing legacy warning)
tsc --noEmit: exit 0
npm test: 38/38 pass
```

### Stale path references

Only remaining old-path mention of `.hawp/kit/spec.md` / `.hawp/kit/authoring-patterns.md` is in the closed archive `closed/2026/06/04/TASK-072.md` (historical record, preserved intentionally).

## Drift found and fixed

`.github/instructions/code-style.instructions.md` existed only at repo root (not in `core/providers/.github/instructions/`), had mangled YAML frontmatter (heading inserted before the frontmatter fence, escaped glob characters), and referenced a nonexistent path `shared_standards/public/guidelines/code-style.md`. It also claimed "No bundler — direct tsc compilation only," contradicting the kit standard baseline.

Fix: rewrote the file with valid frontmatter, pointed it at `.hawp/kit/standards/guidelines/code-style.md`, aligned the key rules with that standard, and added it to the core pack so root and pack match.

## Inference

The earlier `git log` shows `code-style.instructions.md` was introduced in commit `89bf269` and likely mangled by a markdown formatter pass afterward; the `shared_standards/` path appears to be a leftover from a pre-migration layout. Not directly proven.
