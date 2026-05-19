# Code Style Instructions

---

applyTo: "**/\*.ts,**/_.tsx,\*\*/_.js,\*_/_.mjs"
description: Apply the shared code style standard when writing or reviewing TypeScript/JavaScript code in this repo

---

# Code Style

Full standard: `shared_standards/public/guidelines/code-style.md`

Key rules enforced in this repo:

- Prefer arrow functions for declarations and callbacks.
- Import groups in order: type-only → external → internal → sibling. No mixing.
- Extensionless local TypeScript imports (no `.ts` or `.js` suffix on relative paths).
- Kebab-case file and folder names; `README.md` and `CHANGELOG.md` are exceptions.
- `strict: true` always on; `noImplicitAny`, `exactOptionalPropertyTypes`, `strictNullChecks`.
- No bundler — direct `tsc` compilation only.
