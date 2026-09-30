# Slice Harness — Universal Template

Every non-trivial implementation slice records a harness. This template
applies to any compoundable unit of work: a new feature, a refactor, a
standards document, a provider adapter, or a CLI command.

## Header block (fill in before writing code)

```
Work UUID:   <uuid>
Plan file:   .hawp/work/active/<uuid>/plan.md
Branch:      <branch-name>
Owned paths: <list every file this slice creates or modifies>
Session:     <YYYY-MM-DD>
```

## Intent, constraints, non-goals

**Intent** — one sentence: what this slice achieves.

**Constraints** — what the slice must not do (do not move files, do not
introduce new dependencies, do not edit generated output, etc.).

**Non-goals** — explicitly out of scope for this UUID.

## Verification order

Run focused checks first. Only escalate to broad checks when focused checks
pass. Do not assert results — record actual command output.

### 1. Focused package or file check

```bash
# Go source changes: run the affected package before the full suite
cd librarian/src && go test ./internal/platform/cli/init/... -v
```

Record actual output, not expected output.

### 2. Build and vet

```bash
cd librarian/src && go build ./...
cd librarian/src && go vet ./...
```

### 3. HAWP validation

```bash
cd librarian/src && go run ./cmd/hawp kit validate
cd librarian/src && go run ./cmd/hawp work validate
cd librarian/src && go run ./cmd/hawp check --no-update-check
```

### 4. Full test suite (only if Go source was modified)

```bash
cd librarian/src && go test ./...
```

## Evidence record

Paste actual output here. Separate verified facts from inference.

| Claim | Evidence type | Source |
|-------|---------------|--------|
| Tests pass | direct — `go test` output | paste below |
| Build clean | direct — `go build` output | paste below |
| HAWP check pass | direct — `hawp check` output | paste below |

```
<paste actual command output here>
```

## Verified facts vs. unproven claims

- **Verified:** list only what the evidence above directly shows.
- **Inferred:** note anything that follows from reasoning but is not
  directly proven by the recorded output.

## Attribution

Real contributors only. Do not invent identities or model names.

```
Author: <real name or omit>
```

If assisted by a model, note the task shape and that output was reviewed;
do not fabricate a co-author trailer.

## Continuation instructions

What the next session must know to resume this slice:

- Which step was last completed.
- Which commands were last run and what they returned.
- Remaining work and the first command to run on restart.
- Any deferred items with reason.

## Artifact paths

```
Status checkpoints: .hawp/work/status/YYYY/MM/DD/<uuid>/
Evidence files:     .hawp/work/evidence/YYYY/MM/DD/<uuid>/
Plan file:          .hawp/work/active/<uuid>/plan.md
Closed (when done): .hawp/work/closed/YYYY/MM/DD/<uuid>/plan.md
```

Use UUID-scoped subfolders so a slice can be found by work item without
searching date-only paths.

## Close checklist

- [ ] All focused checks pass (recorded above).
- [ ] HAWP check passes.
- [ ] Plan status updated to `done`.
- [ ] Plan moved to `.hawp/work/closed/YYYY/MM/DD/<uuid>/plan.md`.
- [ ] BACKLOG.md updated (Recently Closed capped at 5–10 items).
- [ ] No invented evidence, authors, or identities.
