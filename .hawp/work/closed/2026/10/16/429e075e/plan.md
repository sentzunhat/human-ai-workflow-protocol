# Remote MCP transport and authorization readiness audit

**UUID:** `429e075e-397d-4ed9-845e-fe349f1559c9`
**Type:** improvement
**Reported:** 2026-09-13
**Status:** `closed`
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

- [x] Compare HAWP's current stdio server and provider config writers against current official MCP transport and authorization guidance.
- [x] Identify which assumptions are local-only: executable path, repo-root trust,
  filesystem permissions, local search index, and user approval model.
- [x] Document which remote modes are explicitly unsupported today.

#### Audit findings (2026-10-16)

**Directly observed in this checkout:** `librarian/src/internal/platform/mcp/server/server.go`
implements newline-delimited JSON-RPC over stdin/stdout and returns
`2024-11-05` from `initialize`. It dispatches only `initialize`, `tools/list`,
and `tools/call`; no HTTP listener, Streamable HTTP framing, HTTP auth layer,
resources, or prompts are implemented. The current tool set includes both
read-only operations and work-record creation (`hawp_work_new`), so a remote
deployment would require a principal/permission model and filesystem scope; a
single server-wide bearer token alone is not a complete authorization design.

The provider config writers produce stdio launch entries. The JSON writer
refuses an existing remote `url` or non-stdio `type` rather than converting it;
the Codex TOML writer similarly requires manual review for remote/non-stdio
entries. These safeguards are evidence that remote configuration is not
currently supported by HAWP's setup path.

**Official protocol guidance checked:** the current `2026-07-28` MCP spec names
stdio and Streamable HTTP as standard transports. MCP HTTP authorization is
optional at protocol level, but HTTP implementations that support it should
follow MCP OAuth requirements; stdio implementations should obtain credentials
from their environment instead. The `2025-11-25` Streamable HTTP guidance
requires Origin validation against DNS rebinding and recommends localhost-only
binding for local servers plus authentication. The current authorization spec
requires protected-resource metadata, audience-bound token validation, least
privilege scopes, and per-request bearer authorization. MCP security guidance
also calls out SSRF in OAuth metadata discovery, confused-deputy/token
passthrough risks, local-server privilege, and explicit user consent for local
server setup.

**Supported now:** local, client-launched stdio for a selected repository,
using the process user's filesystem permissions and the configured repository
root. The client/server process boundary is the access boundary; this is not
multi-tenant isolation.

**Unsupported now:** Streamable HTTP (local or hosted), remote MCP endpoints,
OAuth authorization/discovery, bearer-token authentication, per-user or
per-tool authorization, tenant/repository identity binding, network exposure,
and remote write/delete operations. No remote listener or HTTP+SSE legacy
transport exists. Do not describe HAWP as remotely deployable based on the
MCP protocol's general capabilities.

**Security gates before any future remote proposal:**

1. Choose the deployment boundary and supported protocol revision; use
   Streamable HTTP for standard remote transport and explicitly reject legacy
   HTTP+SSE unless a compatibility requirement justifies it.
2. Require HTTPS in production, validate Origin on every request, and define
   safe local binding behavior; do not expose a listener by default.
3. Implement OAuth resource-server discovery and validate token signature,
   issuer, expiry, audience, and required scopes on every request. Never accept
   or relay tokens issued for downstream APIs.
4. Bind each principal to an explicitly authorized repository root and
   filesystem boundary. Define read/create/update/delete permissions separately;
   require explicit human approval for destructive or externally consequential
   actions.
5. Threat-model metadata-discovery SSRF, confused deputy, cross-tenant data
   access, credential/path leakage, request size/rate limits, audit logging,
   revocation, and secret storage. Establish tests for denial cases before
   implementation.
6. Keep stdio as the default local deployment. Do not add a remote mode until
   a separate implementation plan is reviewed and approved.

### Slice B - Security boundary proposal

- [x] Define minimum requirements before any remote MCP HAWP mode:
  authentication, scoped authorization, repository binding, audit logging,
  secret/path redaction, and explicit operator approval for writes.
- [x] Decide whether remote support belongs in HAWP core, a gateway adapter, or
  a separate experimental package: recommendation is a separately scoped
  gateway/adapter or experimental package, not the current stdio server, until
  a concrete deployment use case justifies the operational and security cost.
- [x] Preserve the default local-first posture unless a later decision approves remote work.

### Slice C - Provider documentation update

- [x] Add a short compatibility matrix to `.hawp/kit/usage/mcp/verification.md`.
- [x] Update troubleshooting language so users distinguish local stdio readiness
  from remote MCP capability.
- [x] Add explicit non-goals for HTTP/remote transport if implementation is deferred.

## Risk

Medium-high. Remote MCP touches trust boundaries and should remain planning-only
until the security model is explicit.

## Verification

- [x] Audit references current official MCP documentation and exact HAWP source/docs paths.
- [x] No remote listener is added in the audit slice.
- [x] Documentation distinguishes local stdio support from unsupported remote modes.
- [x] `go run ./cmd/hawp work validate` passes from `librarian/src` (5 checks, 0 issues, 0 warnings).

## Audit conclusion

This item is a planning/documentation audit and does not add remote transport,
authorization code, or runtime behavior. The recommended decision is to keep
HAWP local-first and stdio-only; treat remote MCP as unsupported until a
separately reviewed use case and implementation plan pass the security gates
above. The capabilities catalog (`2eea565c`) is the more compoundable next
planning lane because it can define every capability and its read/write/delete
boundary before adding surfaces. That catalog should inventory mutations as
well as resources/prompts; it must not imply that write or delete tools are
safe to expose remotely.

**Decision (2026-10-16):** Keep HAWP on local stdio while it serves the local,
standalone use case. Consider Streamable HTTP only when a concrete need arises
for a remote listener, shared HAWP service, or multiple clients connecting to
one HAWP server. Do not build HAWP's own OAuth authorization server for the
initial HTTP deployment; use an existing identity provider, with HAWP acting as
the protected resource server and validating token identity, audience, scopes,
and repository access.

**Future dependency:** Any future HTTP MCP implementation plan must depend on
this completed audit and carry forward its security requirements: standard
Streamable HTTP, external OAuth authorization service, authenticated principal
to repository binding, least-privilege operation permissions, Origin/HTTPS and
SSRF protections, and denial/isolation verification before deployment. This
record is the prerequisite; no HTTP implementation work is authorized here.

## Outcome and Closeout

Outcome: HAWP remains local-stdio-first. Any future HTTP MCP implementation is
gated on this audit's decisions and security requirements. The authorization
server remains an external identity-provider responsibility; HAWP would act as
the protected resource server.

## Close Checklist

- [x] Outcome and security decision recorded.
- [x] Future HTTP MCP work dependency recorded.
- [x] Compatibility and unsupported-mode guidance added to the MCP verification guide.
- [x] Backlog status synchronized and plan moved to the dated closed archive.
- [x] HAWP work validation passed (5 checks, 0 issues, 0 warnings).

**Close date:** 2026-10-16

## Modular implementation roadmap (proposal only)

The layers below are a proposed sequence, not approved implementation work.
Keep the existing local stdio path as the default and preserve the application
use cases as the transport-independent core.

### Stage 0 - Decide whether remote hosting is needed

- Name the concrete deployment, users, repository ownership model, and threat
  boundary. If the goal is only to let a user run HAWP in their own trusted
  agent workspace, retain stdio; do not add OAuth or an HTTP service.
- If multiple users or hosted agents must access HAWP, select a supported
  external identity provider / authorization server before implementation.
- Define who provisions users and repository access, and how access is removed.
  Without those operating decisions, implementation should not proceed.

### Stage 1 - Inventory operations and define authorization policy

- Use `2eea565c` to catalog every tool/use case by data access and effect:
  read/search, local model invocation, create, update, delete, and validation.
- Define principal, repository, and permission types in an application-facing
  policy contract. Resolve the repository from server-side authorized identity
  and configuration; do not accept an arbitrary root path supplied by a remote
  caller.
- Start with least-privilege read-only access. Keep creation and any future
  update/delete actions separately scoped; require explicit confirmation for
  destructive actions. No delete MCP tool currently exists.
- Add policy checks at the use-case/tool boundary, not only at the HTTP route,
  so a second transport cannot bypass them. Keep stdio's local trust policy
  explicit and compatible with current behavior.

### Stage 2 - Separate protocol/transport adapters

- Current `internal/platform/mcp/server` owns a handwritten JSON-RPC loop,
  dispatch, tool metadata, and tool handlers. Extract stable application
  capability handlers from transport framing and request decoding before
  introducing another transport.
- Evaluate migration to the official `github.com/modelcontextprotocol/go-sdk`
  transport/server implementation. The current server advertises protocol
  `2024-11-05`; remote implementation should target a selected current MCP
  revision and preserve tested client compatibility. Avoid maintaining a
  second handwritten protocol implementation.
- Keep stdio and a future HTTP adapter independently startable. Do not make
  either adapter choose a different repository or authorization policy by
  accident.

### Stage 3 - Add HTTP resource-server authentication

- HAWP should normally **not create its own OAuth authorization server**. The
  MCP client is the OAuth client; an external identity provider / authorization
  service issues tokens; HAWP is the OAuth-protected MCP resource server that
  validates tokens and enforces access. This matches the MCP authorization
  roles and avoids building account, consent, token issuance, refresh, key
  rotation, and recovery operations into HAWP.
- Implement Protected Resource Metadata discovery and bearer-token validation
  for the selected issuer using a maintained library or trusted provider
  integration. Verify issuer, signature, expiry, audience/resource, and scopes
  on every request; reject tokens meant for downstream services. Do not invent
  a homegrown OAuth server or accept unsigned identity headers from a proxy.
- Keep secrets and signing keys with the identity provider. HAWP should not
  store user passwords or issue refresh tokens. The MCP client performs its
  supported authorization flow and stores/refreshes client-side credentials.
- If deployment uses a gateway to terminate OAuth, specify and verify the
  authenticated-principal handoff (for example, a mutually authenticated
  gateway-to-service channel plus signed/verified identity context). A
  caller-controlled `X-User`/`X-Scopes` header is not authentication.

### Stage 4 - Enforce repository and operation boundaries

- Bind verified identity to an explicit repository allowlist or server-side
  repository mapping. Constrain all file access to that root and reject
  symlink/path escapes; do not use client-supplied paths as authority.
- Map scopes/policies to operations (initially read-only; then separately
  approved create/update/delete capabilities). Enforce policy in handlers or
  application services and return consistent unauthenticated/forbidden errors.
- Make human confirmation a distinct gate for irreversible or destructive
  operations. OAuth scope grants alone must not silently approve a particular
  deletion or overwrite.

### Stage 5 - Harden and prove the deployment

- Require HTTPS outside local development; validate `Origin`; define CORS,
  request/body limits, timeouts, rate limits, and safe shutdown.
- Threat-model OAuth discovery SSRF and redirects, DNS rebinding, token theft,
  confused deputy, cross-tenant access, path/secret leakage, and audit-log
  privacy. Restrict outbound discovery and redact authorization headers and
  sensitive paths from logs.
- Test auth denial, wrong issuer/audience, expired/revoked token, missing scope,
  repository mismatch, path traversal, symlink escape, cross-user isolation,
  Origin rejection, and concurrent requests before enabling any hosted mode.
- Add deployment and operator guidance only after an actual provider/client
  combination has been exercised end to end. Keep remote mode opt-in and off by
  default.

### OAuth role decision

**Recommendation:** HAWP does not operate an OAuth authorization server for
the first remote deployment. Select an existing OIDC/OAuth provider and
configure it as the authorization server. HAWP implements the protected
resource side (discovery metadata, token validation, audience and scope checks,
principal/repository mapping). MCP clients implement the OAuth client side.
If no suitable external provider exists, prefer evaluating a managed/self-hosted
identity provider as a separate dependency; only consider HAWP-owned token
issuance if there is a clear product requirement and a separately reviewed
identity/security operations plan.

The available evidence supports this separation: current HAWP is an MCP server
that directly invokes local application use cases; no MCP client authorization
flow, identity store, OAuth endpoints, or account lifecycle exists in the
repository. The official MCP authorization spec describes the protected MCP
server as a resource server and the MCP client as OAuth client; it defines
authorization for HTTP transports and advises stdio servers to source
credentials from the environment. The official Go SDK documents server-side
bearer-token middleware and Protected Resource Metadata helpers, but those
remain implementation options to evaluate after Stage 0, not a decision to add
an SDK dependency now. References:

- [MCP 2026-07-28 transport specification](https://modelcontextprotocol.io/specification/2026-07-28/basic/transports)
- [MCP 2026-07-28 authorization specification](https://modelcontextprotocol.io/specification/2026-07-28/basic/authorization)
- [MCP security best practices](https://modelcontextprotocol.io/docs/2026-07-28/tutorials/security/security_best_practices)
- [Official Go SDK authorization and transport guidance](https://github.com/modelcontextprotocol/go-sdk/blob/main/docs/protocol.md)
