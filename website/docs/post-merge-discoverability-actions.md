# Post-merge discoverability actions

Captured: 2026-09-24

This document records the repository and domain actions that are intentionally
deferred until the website pull request is merged. They are administrative or
external-state changes, so they cannot be proven by the website build alone.

## Why these actions are deferred

The website branch can ship the public HTML, metadata, crawl instructions, and
static assets. It cannot change GitHub repository settings, enable Pages for
`main`, configure DNS, or submit the site to search engines. Those actions
should be performed as one post-merge launch step so the discoverability
baseline has a clear boundary.

## GitHub repository settings to apply after merge

Use the following proposed values. Confirm the final wording against the
merged README and release state before saving settings.

### Description

> HAWP is an open-source, Markdown-first workflow protocol and Go CLI for shaping AI agent work, preserving context, reducing drift, and making handoffs portable across tools.

This description is deliberately specific about the category, format, and
user benefit. It does not claim that HAWP is an agent runtime, hosted platform,
or proprietary memory system.

### Homepage

`https://hawp.online/`

### Suggested topics

Review these against GitHub's topic rules and the final product vocabulary:

- `ai-agents`
- `ai-workflow`
- `context-engineering`
- `developer-tools`
- `human-ai-collaboration`
- `markdown`
- `mcp`
- `workflow-automation`

Topics are a proposed set, not a claim that they have already been applied.

## GitHub Pages and domain actions

After the website workflow is merged into `main`:

1. Confirm the Pages deployment succeeds and serves the built artifact.
2. Set the Pages custom domain to `hawp.online`.
3. Configure the apex and `www` DNS records at the domain provider according
   to the GitHub Pages values shown for the repository.
4. Wait for GitHub's domain validation, then enable **Enforce HTTPS**.
5. Verify `https://hawp.online/`, `/robots.txt`, `/sitemap.xml`, and
   `/llms.txt` over HTTPS.

Do not record DNS, Pages, or HTTPS as complete based only on a local build.

## Community-standard items still requiring GitHub UI or API changes

The pre-website Community Standards screenshot showed these items as
incomplete:

- repository description
- code of conduct
- contributing guidelines
- security policy
- issue templates
- pull request template

The repository already contains project guidance and security-related material
in its source tree, but the GitHub Community Standards panel must be checked
after merge to confirm which files and settings it recognizes. The panel state
is external evidence and is not inferred from local files.

## Search and AI discovery actions

After the canonical HTTPS URL is live:

- submit `https://hawp.online/sitemap.xml` to Google Search Console and Bing
  Webmaster Tools, if those accounts are available;
- inspect URL indexing and canonical selection for the homepage;
- verify that the GitHub repository homepage links to the same canonical URL;
- record branded queries for `HAWP`, `Human-AI Workflow Protocol`, and
  `HAWP AI agents` before and after launch;
- compare GitHub views, unique visitors, clones, unique cloners, referring
  sites, and popular paths at +7 and +14 days using the baseline in
  [discoverability-baseline.md](./discoverability-baseline.md).

Search ranking, indexing, referrals, and AI answer inclusion remain unproven
until these external checks are performed. Structured data and `llms.txt` can
make the site's meaning easier to retrieve, but cannot guarantee indexing or
ranking.

## Evidence to attach after launch

Capture the following as direct evidence rather than changing this document
from assumption to completion by memory:

- merged commit or release identifier;
- successful GitHub Pages deployment URL;
- HTTPS responses for the canonical URL and discovery files;
- final GitHub description, homepage, topics, and Community Standards state;
- Search Console/Bing submission or indexing readback, where available;
- the +7-day and +14-day comparison tables.

