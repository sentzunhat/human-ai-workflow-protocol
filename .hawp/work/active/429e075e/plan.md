# Remote MCP transport and authorization readiness audit

**UUID:** `429e075e-397d-4ed9-845e-fe349f1559c9`
**Type:** improvement
**Reported:** 2026-09-13
**Status:** `plan-ready`
**Target Version:** `v0.0.25`

---

## Input

User shared the YouTube video `https://youtu.be/BqRhBq-_kgE?si=Kt0vrAV5XYL4269g`
and asked to use its MCP architecture lesson as learning for planning HAWP
improvements.

## Source Notes

- Confirmed video metadata via YouTube oEmbed: title `MCP Just Got a Whole Lot Better`, author `Neon Postgres`.
- Transcript retrieval was blocked by YouTube during intake. Treat this plan as
  architecture-inspired and research-first.
- HAWP currently documents and implements local stdio MCP server setup, with
  explicit `--repo-root` and provider-specific local config preservation.

## Lesson

As MCP matures beyond local stdio demos, HAWP should be clear about which
deployment modes it supports, which modes it intentionally rejects, and what
security gates are required before remote access exists.

## Release Scope

This is a `v0.0.25` audit lane item. It should not add a remote listener in
`v0.0.24`; the current release should preserve HAWP's local-first MCP posture.

Version-scope note: see
`.hawp/work/status/2026/09/13/6059dd4b/status.md`.

## Plan

### Slice A - Readiness audit

- [ ] Compare HAWP's current stdio server and provider config writers against current official MCP transport and authorization guidance.
- [ ] Identify which assumptions are local-only: executable path, repo-root trust,
  filesystem permissions, local search index, and user approval model.
- [ ] Document which remote modes are explicitly unsupported today.

### Slice B - Security boundary proposal

- [ ] Define minimum requirements before any remote MCP HAWP mode:
  authentication, scoped authorization, repository binding, audit logging,
  secret/path redaction, and explicit operator approval for writes.
- [ ] Decide whether remote support belongs in HAWP core, a gateway adapter, or
  a separate experimental package.
- [ ] Preserve the default local-first posture unless a later decision approves remote work.

### Slice C - Provider documentation update

- [ ] Add a short compatibility matrix to `.hawp/kit/usage/mcp/verification.md`.
- [ ] Update troubleshooting language so users distinguish local stdio readiness
  from remote MCP capability.
- [ ] Add explicit non-goals for HTTP/remote transport if implementation is deferred.

## Risk

Medium-high. Remote MCP touches trust boundaries and should remain planning-only
until the security model is explicit.

## Verification

- Audit references current official MCP documentation and exact HAWP source/docs paths.
- No remote listener is added in the audit slice.
- Documentation distinguishes configured, connected, authorized, and executed states.
- `go run ./cmd/hawp work validate` passes from `librarian/src`.
