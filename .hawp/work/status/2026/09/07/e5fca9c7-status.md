# Status Report — Code Progress And Next Work

## Intent

Summarize the recent implementation/cleanup sequence, assess the code within
inspected scope, and order the next compoundable work. Parent:
[release/coordination plan](../../../../active/e5fca9c7/plan.md).

## Current State

The branch is `feature/v0.0.24-work-folder-normalization`; HEAD is `78c058a`.
Recent implementation, tests, guidance, and work-record changes remain
uncommitted. Local checks are passing; this is not a published release or a
blanket readiness certification. Three active items remain: release coordination,
architecture/correctness, and intake reshaping. Migration `0cb0f9b0` is closed
within its documented local compatibility scope.

## What Was Inspected

Current git state and backlog; active plans; migration closeout; CLI parsers,
corpus walker and ingestion; search service/port; downloader tests; source/kit
guidance; the latest local test result and connected MCP validation.

Key source references:

- `librarian/src/internal/domain/index/chunk.go`
- `librarian/src/internal/infrastructure/sqlite/index.go`
- `librarian/src/internal/platform/cli/index_corpus.go`
- `librarian/src/internal/application/index/ingest-service.go`
- `librarian/src/internal/domain/search/index.go`
- `librarian/src/internal/application/search/service.go`
- `librarian/src/internal/domain/distribution/binary_download_test.go`

## What Changed

| Area | Completed locally | Why it matters |
| --- | --- | --- |
| Index models | Reviewed and verified the pre-existing Chunk/DocumentMetadata consolidation; domain ownership with SQLite aliases, nullable context, and line ranges | Reduces duplicate definitions without a schema change |
| Embedding CLI | Fixed flag values being treated as text and discarded parser errors; eight parser cases | Invalid input is caught before model setup |
| Model pull | Restored documented trailing options while preserving leading options; 16 cases | Documentation examples now match parser behavior |
| Mutation boundaries | Reject conflicting kit modes and explicit empty/whitespace path overrides | Prevents accidental apply mode or current-repo fallback |
| Mutation proof | Four handler tests snapshot temporary files and directories around rejection | Tests observable preservation, beyond parser return values |
| Native migration | Native install/config paths aligned; unused core source launcher retired; installed legacy files preserved | Removes the wrapper dependency for new installations/configs |
| Download proof | 40 fresh/legacy scenarios plus four repeated successes; checks failures, installed bytes, executable permissions, and staging cleanup | Makes compatibility/failure behavior repeatable |
| Work/guidance cleanup | Preserved old intake/history, reconciled migration lifecycle, corrected stale architecture and Go import guidance, repaired references | Reduces misleading instructions for subsequent work |
| Intake reshaping | Created a dedicated item and documented/exercised worker-guided shaping with original input preserved | Establishes a usable workflow; automated reshape is not implemented |

## How The Code Is Coming Along

Assessment: clearer and more testable, with specific correctness improvements;
architecture separation is still partial.

Confirmed strengths: CLI parsing can be tested without storage/models; shared
index types preserve storage semantics; search has an injectable repository
contract; provider configuration has preservation/refusal tests. Recent passes
added no new dependencies or general parser framework.

Confirmed gaps: the search service's default constructor still wires filesystem
and SQLite adapters; the Index contract returns loosely typed maps; provider
implementations remain alongside domain ports. CLI corpus collection still walks
files directly. These are bounded future seams, not reasons for a broad rewrite.

A concrete correctness gap takes precedence: `walkWorkFiles` skips BACKLOG.md and
does not populate WorkUUID. Ingestion persists work metadata only when WorkUUID
is supplied. Exact backlog status and persisted UUID/status parity therefore
remain unresolved in this path. The separate index-build enrichment path needs
comparison before changing the shared behavior.

## What Was Directly Verified

- Latest full Go suite and vet passed during the launcher-retirement continuation;
  no implementation changes have occurred since that run. Not rerun just to write
  this report.
- Distribution validation passed in that continuation.
- Local release-bundle assembly produced 195 kit/provider members without the
  retired launcher; host copy metadata was disabled and temporary files removed.
- Fresh connected `hawp_work_validate` passes kit/work/links with zero issues and
  warnings. Fresh `git diff --check` and root/core kit comparison also pass.

## What Remains Unproven

No fresh live release downloads, all-platform execution, model inference, or
new provider-client sessions. Historical provider/build evidence stays historical.
Downloader fixtures are not whole-installer transaction proof. No performance,
token-saving, concurrency, or broad security claims follow from these changes.
The installed native binary has not been refreshed for the latest source fixes.

## Constraints

This reporting pass changes documentation/work coordination only. Preserve the
uncommitted implementation and unrelated prior work. No commits, pushes, release
publication, or downstream changes are performed here. Reuse existing UUIDs;
these milestones do not need duplicate work items.

## Help Wanted

No immediate clarification needed. Before adding a callable reshaper, assess its
value over the current worker-guided flow and define what fidelity tests can prove.

## Suggested Next Step

| Order | Existing item | Bounded deliverable | Done when |
| --- | --- | --- | --- |
| 0 | e5fca9c7 | Consolidate the verified local checkpoint into focused review groups: models, CLI safety, migration, guidance/work records | Untracked tests/history are included, source/binary differences explicit, and groups are reviewable before any commit |
| 1 | a3df8a9c | Define and test a request-to-intake draft contract | Original input unchanged; only HAWP fields; unknowns explicit; malformed/failing output handled; no work-file writes; decide whether runtime integration is warranted |
| 2 | 47c793d6 | Persist correct work UUID and exact status | Compare corpus paths, resolve authoritative metadata without guessing, and round-trip active/parked/closed/unknown fixtures through SQLite |
| 3 | 47c793d6 | Extract the next useful boundary | Move corpus/storage policy behind an application-owned contract only where it removes duplication or enables meaningful failure tests |
| 4 | e5fca9c7 | Release-readiness review | Choose patch scope, refresh the binary from verified source, validate generated assets/builds/smokes, and distinguish local checks from release/provider evidence |

Do not make the patch release wait for every architecture or reshape idea.
Release scope should be an explicit decision. Unknown-folder routing, historical
benchmark backfill, and cloud/runtime expansion stay outside this next sequence.
