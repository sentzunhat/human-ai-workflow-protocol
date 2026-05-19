# Git Workflow & Versioning Guide

Commit message format, branch naming, versioning strategy, and release procedure for all projects.

## Table of Contents

1. [Commit Message Format](#commit-message-format)
2. [Branch Naming](#branch-naming)
3. [Semantic Versioning](#semantic-versioning)
4. [Pre-PR Validation](#pre-pr-validation)
5. [PR Guidelines](#pr-guidelines)

---

## Commit Message Format

Follow **[Conventional Commits](https://www.conventionalcommits.org/)** format for clear, parseable commit history. This enables automated changelog generation and semantic versioning decisions.

### Format

```
<type>(<scope>): <subject>

<body>

<footer>
```

### Type Classification

| Type       | Meaning                                  | Release Impact  |
| ---------- | ---------------------------------------- | --------------- |
| `feat`     | New feature                              | Minor bump      |
| `fix`      | Bug fix                                  | Patch bump      |
| `refactor` | Code reorganization (no behavior change) | Patch bump      |
| `perf`     | Performance improvement                  | Patch bump      |
| `docs`     | Documentation only                       | No version bump |
| `test`     | Test additions/modifications             | No version bump |
| `chore`    | Build, deps, tooling                     | No version bump |
| `ci`       | CI/CD configuration                      | No version bump |
| `style`    | Formatting (no logic change)             | No version bump |

### Scope (Optional)

Scope references the area affected. Examples:

- `service`, `handler`, `repository`, `adapter`
- `error`, `validation`, `logging`
- `config`, `di-container`, `types`
- `tests`, `docs`

```
feat(error): add ForbiddenError class
fix(service): handle null database connections
docs(api): update endpoint documentation
refactor(di-container): simplify auto-registration logic
```

### Subject Line

- Start with **lowercase** (after type/scope)
- Use **imperative mood** ("add" not "added" or "adds")
- No trailing period (.)
- Max 50 characters

```
# ✅ Good
feat(error): add ForbiddenError class
fix(service): handle null database connections
docs(api): update endpoint documentation

# ❌ Avoid
feat(error): Added ForbiddenError class  # Wrong case, past tense
fix(service): Fixes null database connections  # Wrong mood
docs(api): updated endpoint documentation  # Wrong case
```

### Body (Optional but Recommended for Substantive Changes)

Include body for non-trivial commits:

```
feat(error): add ForbiddenError class

Introduces new ForbiddenError for 403 Forbidden responses.
Used when a user is authenticated but lacks permissions for the resource.

- Extends CustomError with code 403
- Supports metadata for additional context
- Aligns with existing error hierarchy
```

**Guidelines:**

- Explain **what** changed and **why**
- Use imperative mood ("refactor" not "refactored")
- Wrap at 72 characters
- Separate from subject with blank line

### Footer (Optional)

Use for breaking changes or issue references:

```
feat(service)!: make Service constructor async

BREAKING CHANGE: Service constructor is now async and requires await.
Before: const service = new Service(config);
After: const service = await new Service(config);

Fixes #42
Relates to #88
```

**Footer Keywords:**

- `BREAKING CHANGE:` — Signals breaking change (requires major version bump)
- `Fixes #<number>` — Closes issue
- `Relates to #<number>` — References related issue
- `Reviewed-by:` — Code reviewer(s)

### Full Example

```
feat(handler): add support for DELETE requests in Express

Implements DeleteRouteHandler abstract class and concrete implementations
for Express framework, matching existing Fastify handler patterns.

- Extends AbstractRouteHandler with signature typed for DELETE
- Supports optional request body and dynamic URL parameters
- Includes type-safe response handling

BREAKING CHANGE: Handler interface now requires explicit HTTP method type.

Fixes #42
Reviewed-by: @alice @bob
```

---

## Branch Naming

**Pattern:** `<type>/<scope>/<description>`

Use the same type prefix as commits for consistency. Keep branch names lowercase, kebab-case, and descriptive.

```
# ✅ Good
feat/error/forbidden-error-class
fix/service/handle-null-connections
docs/api/endpoint-documentation
refactor/di-container/auto-registration
chore/deps/upgrade-typescript

# ❌ Avoid
feature/my-new-thing         # Type too generic, scope missing
bugfix/error                 # "bugfix" instead of "fix"
WIP-service-changes          # All caps, vague
123-add-handler              # Issue number prefix only
```

### Branch Type Options

- `feat/` — Feature branch
- `fix/` — Bug fix branch
- `refactor/` — Code reorganization
- `docs/` — Documentation only
- `test/` — Test additions
- `chore/` — Dependencies, tooling
- `ci/` — CI/CD updates

### Scope & Description

- **Scope:** Area affected (same as commit scope)
- **Description:** Lowercase, kebab-case, concise (max 50 chars total)

```
feat/di-container/resolve-token-type
fix/repository/handle-null-models
docs/readme/add-examples
```

### Cleanup

Delete branches after merging:

```bash
# Locally
git branch -d feat/service/async-constructor

# On remote
git push origin --delete feat/service/async-constructor
```

---

## Semantic Versioning

Follow **[SemVer](https://semver.org/)**: `MAJOR.MINOR.PATCH`.

| Change                                       | Bump    |
| -------------------------------------------- | ------- |
| Breaking change (removes/renames public API) | MAJOR   |
| New feature (backwards-compatible)           | MINOR   |
| Bug fix, refactor, perf                      | PATCH   |
| Docs, tests, tooling only                    | No bump |

Start new projects at `0.1.0`. Promote to `1.0.0` when the public API is stable.

---

## Pre-PR Validation

Before opening any PR, pass:

```bash
npm run type:check && npm run lint && npm run test
```

Or the single alias if configured: `npm run validate`.

- Zero TypeScript errors.
- Zero lint errors (warnings are acceptable, review intent).
- All tests pass.
- Build succeeds for library/service projects.

---

## PR Guidelines

- **Scope**: one logical change per PR. Split unrelated changes.
- **Title**: plain sentence matching the final commit message style.
- **Description**: explain _why_, link to issue or backlog item, list any migrations or config changes.
- **Self-review**: review your own diff before requesting review from others.
- **No force-pushing** to shared branches after review has started.
- **Squash merge** for feature branches to keep `main` history clean (project-level decision — document it).
