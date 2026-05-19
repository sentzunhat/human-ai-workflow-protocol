# librarian

HAWP maintenance and documentation generation tooling.

## Purpose

The `librarian` folder contains scripts that power HAWP packaging, validation, and workflow integrity checks. It is designed to be extensible for future CLI commands and automation integration.

## Script Organization

### Distribution Tools

Scripts for building and validating user-facing installation and update guides.

**Location:** `scripts/distribution/`

- **`build/index.ts`** — Composes distribution guide variants by combining source fragments with bash script extraction
  - Reads source fragments from `distribution/sources/`
  - Extracts bash blocks from authoritative scripts (`distribution/sources/install/script.md`, `distribution/sources/update/script.md`)
  - Applies branch-specific (`main` or `dev`) variable substitution
  - Outputs 4 guide variants: `install-main`, `install-dev`, `update-main`, `update-dev`

- **`validate/index.ts`** — Verifies generated outputs match expected composition
  - Compares actual files against source fragment lists
  - Reports missing or stale generated files
  - Exit code: 0 (success), 1 (validation failed)

### Workflow Validation Tools

Scripts for HAWP backlog and work item integrity verification.

**Location:** `scripts/validate-hawp-workflow/`

- **`index.ts`** — CLI entry point and orchestration
- **`cli.ts`** — Command-line argument parsing
- **`reporter.ts`** — Structured output formatting (plain text, JSON)
- **`orchestrate.ts`** — Validation pipeline coordination
- **`types.ts`** — Shared type definitions

**Validations** (`validations/`):

- `backlog-consistency.ts` — Ensures BACKLOG.md rows match actual plan files
- `id-parser.ts` — Extracts and validates work item IDs (extensible for UUID support)
- `closed-task-completeness.ts` — Checks closed plans include required sections (with legacy tolerance)
- `evidence-integrity.ts` — Validates evidence file references and completeness
- `verification-clarity.ts` — Ensures verification sections avoid ambiguous or unproven claims

## Common Tasks

### Regenerate distribution guides

When source fragments or authoritative scripts change:

```bash
cd librarian
npm run distribution:build
npm run distribution:validate
```

### Validate workflow records

Check backlog consistency and task completeness:

```bash
cd librarian
npm run workflow:validate
```

### Normalize workflow records

Apply closed-record normalization and run validation:

```bash
cd librarian
npm run workflow:normalize
```

### Type checking

Verify TypeScript code:

```bash
cd librarian
npm run typecheck
```

## Future Extensions

The librarian architecture supports:

- **CLI commands** (e.g., `hawp backlog upgrade`) for automated workflow improvements
- **CI/automation integration** via structured JSON output
- **UUID-based work item migration** (extensible id-parser)
- **Multi-root validation** for monorepo support
