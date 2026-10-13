# Use a Node runtime supported by DigitalOcean website builds

**UUID:** `4b3f1581-906d-46b5-9f13-9e424b298ec4`
**Type:** bug
**Reported:** 2026-10-13
**Status:** `in-progress`
**Risk Level:** high
**Target Version:** website launch

## Investigation

The DigitalOcean App Platform build log rejects the website's `engines.node` range, `>=26.0.0 <27`, before dependency installation or the static build begins. The website CI workflow also pins Node `26.10.0`, so the declared/runtime version currently aligns with CI but not with the deployed builder's available runtime catalog.

DigitalOcean's current Node.js buildpack reference lists Node 24 through `24.15.0` and recommends active LTS or stable runtimes. The installed dependency lock includes Vite 8, whose declared engine range allows Node `>=22.12.0`; Node 24 is therefore within the package engine requirements.

**Confidence:** Confirmed for the supplied failure and supported builder versions; local Node 24 check/build pass. Remote deployment remains unverified.

## Mission

Make the website's declared Node runtime resolvable by DigitalOcean App Platform and keep GitHub Actions on a compatible version.

## Plan

1. Change `website/package.json` and the root package entry in `website/package-lock.json` to declare Node `24.x`.
2. Set `website/.nvmrc` and `.github/workflows/website.yml` to `24.15.0`, the latest Node 24 patch listed in the inspected DigitalOcean buildpack reference.
3. Run `npm run check` and `npm run build` with the installed Node 24 runtime; inspect lockfile and diff consistency.
4. Retry the DigitalOcean deployment after this change reaches the configured app source; record that result separately from local/CI proof.

## Options

- **Recommended:** align the package engine to Node 24 and CI to `24.15.0`. This uses a supported LTS major for the deploy builder and remains within Vite 8's engine range.
- Remove `engines.node` and let DigitalOcean choose its default Node 22 runtime. This reduces explicit runtime control and would leave CI on a different major.
- Keep Node 26 and add a custom Docker build. This adds deployment configuration and maintenance to work around a runtime mismatch.

## Work Coordination

**Owner:** Codex
**Implementation status:** approved by user; Node 24 will be used only for this static website build.
**Expected touched paths:**

- `website/package.json`
- `website/package-lock.json`
- `website/.nvmrc`
- `.github/workflows/website.yml`
- `.hawp/work/BACKLOG.md`
- `.hawp/work/active/4b3f1581/plan.md`

**Review gate:** explicit user approval received on 2026-10-13. User confirmed Zacatl is not built from this static-site package and is unaffected.

## Verification

- Observed the supplied DigitalOcean log: buildpack reports no matching Node version for `>=26.0.0 <27`.
- Inspected `website/package.json`, the corresponding root package entry in `website/package-lock.json`, and `.github/workflows/website.yml`.
- DigitalOcean's Node buildpack documentation lists supported Node 24 releases through `24.15.0`; Vite 8's lockfile engine range includes Node 24.
- User approved aligning this static site's build runtime to Node 24 on 2026-10-13.
- Updated `website/package.json` and the lockfile root engine to `24.x`, `website/.nvmrc` to `24.15.0`, and the website workflow setup-node version to `24.15.0`.
- With installed Node `24.17.0`, `npm run check` passed (`svelte-check`: zero errors and warnings) and `npm run build` passed; static output was generated.
- `git diff --check` passed after these changes.
- DNS currently delegates `hawp.online` to DigitalOcean nameservers, but authoritative and recursive lookups returned no apex A/AAAA or `www` CNAME answers. HTTPS could not resolve the hostname. Deployment and live-site verification remain blocked on DNS records and a reachable deploy.

## Close Checklist

- [x] Owner approved implementation.
- [x] Website check and build pass under the aligned runtime.
- [ ] DigitalOcean deployment succeeds with the supported Node 24 version.
- [ ] Outcome and verification recorded; plan archived; backlog updated.
