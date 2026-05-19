# Git Workflow & Versioning Guide

Commit message format, branch naming, versioning strategy, and release procedure for all projects.

**Status:** Standard level - Required

---

## Commit Message Format

Follow **[Conventional Commits](https://www.conventionalcommits.org/)** format for clear, parseable commit history.

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

### Scope

Area affected (optional but recommended):

- `service`, `handler`, `repository`, `adapter`
- `error`, `validation`, `logging`
- `config`, `di-container`, `types`

```
feat(error): add ForbiddenError class
fix(service): handle null database connections
docs(api): update endpoint documentation
```

### Subject Line

- Start with **lowercase** (after type/scope)
- Use **imperative mood** ("add" not "added")
- No trailing period (.)
- Max 50 characters

---

## Branch Naming

**Pattern:** `<type>/<scope>/<description>`

```
# ✅ Good
feat/error/forbidden-error-class
fix/service/handle-null-connections
docs/api/endpoint-documentation

# ❌ Avoid
feature/my-new-thing      # Too generic
WIP-service-changes       # All caps
```

### Branch Type Options

- `feat/` — Feature branch
- `fix/` — Bug fix branch
- `refactor/` — Code reorganization
- `docs/` — Documentation only
- `test/` — Test additions
- `chore/` — Dependencies, tooling
- `ci/` — CI/CD updates

Delete branches after merging:

```bash
git branch -d feat/service/async-constructor
git push origin --delete feat/service/async-constructor
```

---

## Semantic Versioning

Follow **[SemVer](https://semver.org/)**: `MAJOR.MINOR.PATCH`

| Change                                | Bump  |
| ------------------------------------- | ----- |
| Breaking change (removes/renames API) | MAJOR |
| New feature (backwards-compatible)    | MINOR |
| Bug fix, refactor, perf               | PATCH |
| Docs, tests, tooling only             | None  |

Start new projects at `0.1.0`. Promote to `1.0.0` when the public API is stable.

---

## Pre-PR Validation

Before opening any PR, pass:

```bash
npm run type:check && npm run lint && npm run test
```

Or use the combined alias: `npm run validate`.

- Zero TypeScript errors
- Zero lint errors
- All tests pass
- Build succeeds for library/service projects

---

## PR Guidelines

- **Scope**: one logical change per PR
- **Title**: plain sentence matching commit message style
- **Description**: explain _why_, link to issue, list any migrations or config changes
- **Self-review**: review your own diff before requesting review
- **No force-pushing** to shared branches after review has started
- **Squash merge** for feature branches to keep history clean

**Standard level:** Required
