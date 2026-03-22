# build folder-context enrichment layer for kit/work documents (fbf12a93 Slice 1)

**Backlog ID (Legacy):** — (UUID-native item)
**UUID:** `f93bee55-1ff5-4ebd-a135-c95c40fc6a68`
**Type:** feature
**Reported:** 2026-07-21
**Risk Level:** low

---

### Input (what was reported)

> Add more details to the vector search data: give context to the data
> in each folder (kit and work) before searching from there. Scoped as
> Slice 1 of `fbf12a93` — no ONNX dependency, land first.

---

### Context

`fbf12a93`'s 2026-07-21 scope refinement calls for a folder-context
layer before any embedding work: kit/ folder role + work/ record
metadata attached to each document so retrieved chunks carry "this is a
closed feature from 2026-07-20" context, not just raw prose. This item
implements that layer standalone, replacing the `hawp index build`
scaffold (which today only echoes the requested scope — "poc scaffold
only").

---

### Analysis

**Root cause (or most likely cause):**
_Embedding raw file chunks in isolation loses the surrounding context
that makes a kit/work document meaningful._

**Directly verified:**
_Real kit folders: `examples/, instructions/, patterns/, references/,
reviews/, standards/, templates/, types/, usage/` — only `standards/`
has a README.md today. Real work folders: `active/, closed/, parked/,
decisions/, evidence/, notes/, status/`. `domain/work` already has
backlog parsing (`ParseBacklog`) and ID extraction
(`ExtractIDFromFilename`) from the Phase 1 port — reusable here instead
of re-implementing._

**Inferred (not yet proven):**
_Folder role can be classified purely from path segments (kit's
top-level subfolder name; work's active/closed/parked/other), and each
work record's metadata (type, status, closed date, ID) can be resolved
by cross-referencing BACKLOG.md via the existing parser._

**Scope — what else is affected:**
_`internal/domain/index/` (new context/enrichment types),
`internal/application/index/build-service.go` (replace the stub),
`hawp index build` CLI route (add `--export <path>` for JSON output)._

---

### Recommended Fix

- Classify kit documents by folder role (top-level subfolder name, or
  `root` for `.hawp/kit/*.md` directly) and prepend the folder's
  README.md content when one exists.
- Classify work documents by their backlog section (active/closed/
  parked/other) and attach type, status/reason, and the closed date
  (from the `closed/YYYY/MM/DD/` path) resolved via the existing
  `domain/work` backlog parser and ID extractor — falling back to
  filename-derived metadata for records with no matching backlog row
  (decisions/, evidence/, notes/, status/).
- Produce a `ContextPrefix` string per document (e.g. `[kit/usage]` or
  `[work/closed] TASK-086 (feature, closed 2026-07-03)`) ready to
  prepend before chunking — this is the interface Slice 2's embedding
  step will consume.
- `hawp index build [--scope all|work|kit] [--export <path>]`: walk the
  real corpus, enrich, print a summary (counts per role/type), and
  optionally write the full enriched document list as JSON.

**What to verify after:**

- [x] Folder-context layer attaches correct metadata for a sample of
      kit/ and work/ documents (manual spot-check + unit tests)
      (Evidence: `internal/domain/context/kit_test.go` and `work_test.go`
      — role classification, folder-README summary propagation,
      backlog-metadata resolution, closed-date extraction, ID resolution
      for legacy/UUID/date-prefixed filenames)
- [x] Real `hawp index build --scope kit` and `--scope work` run against
      this repo produce accurate role/metadata counts
      (Evidence: real run 2026-07-21 — kit: 100 documents across 9 roles
      matching the actual `.hawp/kit/` subfolder names exactly; work:
      225 documents across 7 roles, 14/225 resolved backlog metadata —
      the rest are archived closed/ records no longer listed in
      BACKLOG.md's capped Recently Closed section, exactly as the
      repo's own backlog-compaction policy intends, not a bug)
- [x] `--export` output is valid JSON with the expected fields
      (Evidence: real `--export` run — 325 = 100 kit + 225 work
      documents, parsed successfully with `json.load`, sample kit doc
      shows `contextPrefix` correctly includes the kit README's
      first descriptive line; sample work doc shows resolved
      type/status/id)
- [x] Documents with no matching backlog row (decisions/, evidence/,
      notes/, status/) still get sensible fallback metadata, not errors
      (Evidence: real run processed all 5 decisions/, 45 evidence/, 17
      notes/, 39 status/ documents without error, each with Role set and
      empty Type/Status as designed)

---

## Outcome (filled at close)

Closed 2026-07-21. Replaces the `hawp index build` scaffold ("poc
scaffold only") with a real folder-context enrichment layer:
`internal/domain/context/` classifies every kit document by its
top-level subfolder (or `root`) and threads in the folder's README
summary when one exists; every work document by its
active/closed/parked/decisions/evidence/notes/status folder, with
type/status/closed-date/ID resolved by reusing the existing
`domain/work` backlog parser and ID extractor from the Phase 1 port —
no new parsing logic duplicated. Each document gets a `ContextPrefix`
(e.g. `[kit/standards] Rules to follow in real work.` or `[work/closed]
TASK-086 (improvement, closed 2026-07-03)`) that is the exact interface
Slice 2's chunking/embedding step will consume.

`hawp index build [--scope all|work|kit] [--export <path>]` prints a
role/type-count summary (never raw content) and can export the full
enriched corpus as JSON. `fbf12a93` is now scoped purely to Slice 2
(ONNX embedding), unblocked to start directly from this enriched corpus
rather than raw files.

## Verification (filled at close)

- Evidence: `go test ./...` — new `internal/domain/context` package (2
  test files) and rewritten `internal/application/index` tests all pass;
  full suite green.
- Evidence: `make check` (vet + tests + build) passes.
- Evidence: real `hawp index build` runs against this repo for
  `--scope kit`, `--scope work`, `--scope all`, and `--export` (counts
  and JSON validity transcribed above).
- Evidence: `hawp check`, `work:validate`, and `check:markdown-links`
  all pass after the change.

## Close Checklist

- [x] Outcome section filled
- [x] Verification section filled
- [x] Plan file saved under `closed/2026/07/21/f93bee55-1ff5-4ebd-a135-c95c40fc6a68.md`
- [x] BACKLOG.md updated
