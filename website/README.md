# HAWP website

Static marketing and documentation entry point for `hawp.online`.

## Stack

- SvelteKit with static adapter
- Svelte 5 + TypeScript (strict)
- Tailwind CSS 4 through the Vite plugin
- Full prerendering with a tiny hydrated theme control: the production site is static HTML/CSS/JS and requires no application server

## Commands

```bash
npm install
npm run check
npm run build
npm run dev
```

## Deployment

`.github/workflows/website.yml` builds `website/` and deploys `website/build` to GitHub Pages after changes reach `main`. The custom domain is `hawp.online`.

After the workflow is merged and Pages is enabled, set the repository Pages custom domain to `hawp.online` and point the domain's apex DNS records at GitHub Pages. Enable **Enforce HTTPS** after GitHub validates the domain.

## SEO and AI discovery

The landing page includes:

- static prerendered HTML
- canonical URL and descriptive title/description
- Open Graph/Twitter metadata
- `SoftwareSourceCode` JSON-LD
- `robots.txt`
- `sitemap.xml`
- `llms.txt` with a concise, human-readable project summary and canonical links
- semantic headings and crawlable links
- JSON-LD for the website, software source, page, and FAQ content

Repository settings that cannot be changed from this branch are recorded in
[`docs/post-merge-discoverability-actions.md`](./docs/post-merge-discoverability-actions.md).

The first release intentionally stays one page. Future `/docs`, `/mcp`, and
`/getting-started` routes can be added without changing the deployment model.

## Themes

The site supports **Light**, **Dark**, and **System** preferences. The choice
is stored locally, System follows the OS color scheme live, and the initial
theme is applied in `app.html` before hydration to avoid a flash of the wrong
palette.

## Analytics baseline

Captured 2026-09-23 before the website/repository metadata release:

- GitHub stars: 0
- forks: 0
- watchers: 0
- repository description: unset
- repository topics: none
- repository homepage: unset
- GitHub Pages: disabled
- exact-name web searches currently surface LinkedIn posts before the GitHub repository in the sampled results

GitHub's authenticated **Insights → Traffic** views/clones/referrers endpoint is not exposed by the available GitHub connector, so capture those four traffic panels manually on release day for a true before/after comparison.

## Version alignment

The website branch is based on the HAWP `0.0.24` release lane. Version-specific capability and benchmark copy should be checked against the root `README.md`, `librarian/CHANGELOG.md`, and repository benchmark evidence before release.
