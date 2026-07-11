# TASK-071 Public Sync and Docs Classification

Date: 2026-06-04

## Scope

- Source standards: `/Users/beltrd/Desktop/projects/sentzunhat/docs/standards/public/**`
- Source docs: `/Users/beltrd/Desktop/projects/sentzunhat/docs/docs/**`
- Destination: `core/.hawp/kit/standards/public/**`

## Direct Evidence

### Public standards full mirror

Command result:

- `source_count=32`
- `dest_count=32`

Result: all files from source `standards/public/**` were mirrored into `core/.hawp/kit/standards/public/**`.

### Destination sample paths

- `core/.hawp/kit/standards/public/context/README.md`
- `core/.hawp/kit/standards/public/exports/hawp-absorbable/manifest.json`
- `core/.hawp/kit/standards/public/exports/machine-readable/standards-boundary-classification.json`
- `core/.hawp/kit/standards/public/guidelines/testing.md`
- `core/.hawp/kit/standards/public/standards/database/mongodb-schema-design.md`
- `core/.hawp/kit/standards/public/standards/docs/hawp-install-update-safety.md`
- `core/.hawp/kit/standards/public/standards/nodejs/project-structure.md`
- `core/.hawp/kit/standards/public/standards/zacatl/service-boundaries.md`
- `core/.hawp/kit/standards/public/templates/ADR.template.md`

### Docs folder classification pass

- `docs_total_count=119`
- `docs_public_path_count=0`

Interpretation:

- `docs/docs/**` has no explicit public-lane path marker.
- By boundary policy, files without explicit public classification remain non-absorbed in this pass.
- Public-safe intake was therefore completed from `standards/public/**` only.

## Decision

- Implemented exhaustive import of all explicitly public standards files.
- Did not directly mirror `docs/docs/**` due missing explicit public boundary in file paths.
- Private/workflow adaptation remains tracked separately under `TASK-070`.
