# TASK-063 Generalized Extraction Review

Generated: 2026-05-15

## Scope

Extracted neutral standards from six adaptation candidates:

- `shared_standards/public/standards/zacatl/service-boundaries.md`
- `shared_standards/public/standards/zacatl/handler-responsibilities.md`
- `shared_standards/public/standards/zacatl/dependency-registration.md`
- `shared_standards/public/standards/zacatl/layered-composition.md`
- `shared_standards/public/standards/zacatl/contract-testing.md`
- `shared_standards/public/standards/zacatl/evidence-linked-documentation.md`

Destination bundle:

- `core/.hawp/kit/standards/service-design/README.md`
- `core/.hawp/kit/standards/service-design/service-boundaries.md`
- `core/.hawp/kit/standards/service-design/handler-responsibilities.md`
- `core/.hawp/kit/standards/service-design/dependency-composition.md`
- `core/.hawp/kit/standards/service-design/layered-composition.md`
- `core/.hawp/kit/standards/service-design/contract-testing.md`
- `core/.hawp/kit/standards/service-design/evidence-linked-documentation.md`

## Extraction Rules Applied

- Removed framework/domain labels from headings and wording.
- Kept only principle-level guidance and explicit "does not include" boundaries.
- Excluded any internal-domain references and private/source-of-truth pointers.
- Preserved intent while avoiding framework-coupled implementation assumptions.

## Privacy Scan Targets

Must not appear in extracted files:

- tekit
- mictlan / micltan
- zacatl
- internal-only domain references tied to private architecture

## Result

- Extraction completed.
- Service-design standards are now indexed in `core/.hawp/kit/standards/README.md`.
- No direct private-domain references should remain in extracted text.
