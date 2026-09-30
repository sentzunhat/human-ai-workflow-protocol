# HAWP Architecture Context

Current review: 2026-09-07. Use `.hawp/work/BACKLOG.md` and the linked active
plans for lifecycle and priorities; this note does not override them.

- `librarian/src/internal/domain/index/chunk.go` owns Chunk and DocumentMetadata.
  Chunk retains nullable FolderContext and source line ranges. SQLite names are
  compatibility aliases; importing domain and adapter packages alone is not a bug.
- The existing architecture continuation is `.hawp/work/active/47c793d6/plan.md`.
  Parser corrections and database round-trip tests are verified locally; broader
  provider composition and persisted work metadata remain open.
- `.hawp/work/closed/2026/09/07/0cb0f9b0/plan.md` records local migration
  completion. Native paths are aligned and the unused core launcher is retired;
  installed legacy compatibility files are preserved.
  Do not infer a v0.1.0 deferral or blanket completion from old checkpoints.
- `.hawp/work/active/a3df8a9c/plan.md` owns request-to-intake shaping.
  Worker guidance exists; no automatic reshape tool is implemented.

Run maintainer checks from `librarian/src` and follow `.hawp/kit/start-here.md`.
Historical source: `.hawp/work/active/47c793d6/architecture-state-prior.md`.
