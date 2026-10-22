# Implement repository-scoped read-only MCP resources

**UUID:** `3cf52223-b5bc-4e75-af79-a7a8763f704c`
**Type:** improvement
**Reported:** 2026-10-16
**Status:** `in-progress`
**Target Version:** `v0.1.0` follow-up lane; implementation explicitly resumed by maintainer on 2026-10-16

---

## Original Input (verbatim)

> After the v0.1.0 release, investigate and plan a small read-only MCP resources capability for local HAWP. Review candidate resources such as the operating guide, backlog summary, selected active work plan, and MCP setup guidance. Define stable URI/content contracts, repository scope and path safety, handling of missing or stale files, client compatibility, and tests. Do not add resources before separate review and implementation approval.

## Resume Decision (2026-10-16)

The maintainer explicitly authorized continuing from the current HAWP capability
catalog and implementing curated resources/prompts, followed by later provider
and distribution propagation. This supersedes the prior parked/no-implementation
gate for this item. Keep the resource set minimal, read-only, and bounded to
fixed repository paths; no arbitrary path reads or HTTP support.

## Investigation

### Directly verified

- Before the current edits, the source stdio MCP server advertised only `tools` and dispatched `tools/list` and `tools/call`.
- The implementation now registers a fixed resource catalog for the start guide, backlog, and MCP setup guide, all resolved beneath the selected `repoRoot`.
- Resource reads reject symlinked ancestors and hard-linked files, require a regular file, and cap content at 1 MiB.
- The new resource protocol methods are `resources/list` and `resources/read`; no resource templates or arbitrary path reads are exposed.
- Focused Go tests now cover listing, all allowlisted reads, invalid/unknown/missing URIs, symlink and hard-link rejection, and the size limit (`server/resources_test.go`).

### Inference / open questions

- Client-side resource discovery and rendering remain unverified in HAWP's supported clients.
- File-read race behavior between path validation and `os.ReadFile` has not been tested or hardened with descriptor-relative opens.
- Content freshness is read at request time; consumers must treat repository content as untrusted data, not server instructions.

### Scope and risk

- Scope: implement and verify a local stdio, read-only resource capability. No HTTP transport, writes, resource templates, arbitrary path reads, user-home data, or arbitrary active-plan IDs.
- Risk: medium because this adds a new protocol surface and filesystem read path.
- Review gate: maintainer authorized implementation on 2026-10-16; keep the resource allowlist and repository boundary explicit.

## Plan

1. Keep the fixed URI/MIME catalog and error contract stable.
2. Review all allowlisted paths for symlink, hard-link, type, size, and concurrent replacement behavior.
3. [x] Add focused protocol/resource tests for listing, allowed reads, invalid URIs, missing files, symlink/hard-link rejection, and size limits.
4. Run client discovery/read smokes on intended supported clients; document per-client support.
5. Propagate provider/distribution guidance through the parked `3c137bad` item after client verification.

Initial implementation: minimal fixed resources for the operating guide, backlog, and MCP guide. Selected-plan resources remain deferred pending a typed ID contract and client value evidence.

## Backlog / Work State

- Status: `in-progress` as of 2026-10-16.
- Plan path: `active/3cf52223/plan.md`.

## Verification

- `go test ./internal/platform/mcp/server` passes, including the new resource protocol tests.
- `./.hawp/bin/hawp kit validate`, `./.hawp/bin/hawp work validate`, and `git diff --check` pass; work validation reports 0 issues and 1 existing verification-clarity warning.
- Client discovery/read behavior remains unverified. Resource file replacement between path checks and read remains a review item.
