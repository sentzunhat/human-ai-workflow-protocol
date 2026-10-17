# Discoverability baseline and comparison plan

Captured: 2026-09-24

This is the pre-website launch baseline. The values below were transcribed from
user-provided GitHub Insights screenshots captured on 2026-09-24. They are
direct observations from the displayed GitHub UI, not API-readback values.

## GitHub baseline

| Metric | Baseline |
| --- | ---: |
| Stars | 0 |
| Forks | 0 |
| Watchers | 0 |
| Repository description | Unset |
| Topics | 0 |
| Homepage | Unset |
| GitHub Pages | Disabled at capture |

## Insights traffic baseline

The GitHub Traffic view showed the following totals for the displayed 14-day
window, dated 2026-09-10 through 2026-09-23:

| Metric | Pre-website value |
| --- | ---: |
| Total views | 328 |
| Unique visitors | 4 |
| Git clones | 971 |
| Unique cloners | 280 |

The screenshots showed the daily graph shapes, but not the exact value for
every day. Referring sites and popular paths were not displayed in the supplied
evidence and remain unrecorded.

## Repository activity baseline

The Pulse view covered 2026-08-24 through 2026-09-24:

| Metric | Pre-website value |
| --- | ---: |
| Active pull requests | 1 |
| Merged pull requests | 27 |
| Active issues | 0 |
| Closed issues | 0 |
| New issues | 0 |
| Authors with commits | 4 |
| Commits pushed to `main` | 36 |
| Commits pushed to all branches | 81 |
| Files changed on `main` | 374 |
| Additions on `main` | 14,888 |
| Deletions on `main` | 12,864 |
| Published releases | 19 |
| Release publishers | 1 |

## GitHub Actions baseline

The Actions metrics views showed the current-month period beginning
2026-09-01. The usage view said the data was current to approximately 29
minutes before the screenshot; the performance view said approximately 50
minutes before the screenshot.

### Usage

| Metric | Current-month value |
| --- | ---: |
| Total minutes | 1,046 |
| Total job runs | 392 |

| Workflow | Minutes | Runs |
| --- | ---: | ---: |
| `copilot-pull-request-reviewer` | 571 | 77 |
| `quality.yml` | 307 | 148 |
| `sync-distribution-generated.yml` | 148 | 147 |
| `website.yml` | 14 | 14 |
| `tag-on-merge.yml` | 2 | 2 |
| `test-auto-update.yml` on macOS | 2 | 2 |
| `test-auto-update.yml` on Linux | 2 | 2 |

### Performance

| Metric | Current-month value |
| --- | ---: |
| Average job run time | 2m 6s |
| Average queue time | 4s |
| Job failure rate | 10% |
| Failed-job minutes | 79 |

| Workflow | Failure rate | Average run time | Runs |
| --- | ---: | ---: | ---: |
| `copilot-pull-request-reviewer` | 0% | 6m 56s | 76 |
| `quality.yml` | 21% | 1m 31s | 141 |
| `website.yml` | 0% | 26s | 14 |
| `sync-distribution-generated.yml` | 2% | 25s | 140 |
| `test-auto-update.yml` | 0% | 9s | 2 |
| `tag-on-merge.yml` | 0% | 7s | 2 |

## Community standards baseline

The Community Standards view showed these completed items:

- README
- License
- Repository administrators accept content reports

The following items were shown as incomplete or needing action:

- Repository description
- Code of conduct
- Contributing guidelines
- Security policy
- Issue templates
- Pull request template

The screenshots did not show the repository topics or homepage fields as
configured, so those remain the values recorded in the GitHub baseline table
above.

## Adjacent-project scan

A useful adjacent reference is AGENTS.md: the strongest part of its public positioning is not visual complexity, but a direct explanation, an immediate example, and a dedicated website. HAWP should keep that clarity while remaining distinct:

- repository instruction files primarily tell an agent **how to behave**;
- HAWP shapes and carries **the work itself** through intent, context, mission, constraints, output, evidence, and handoff.

## Compare at +7 and +14 days

Record:

- GitHub views and unique visitors
- clones and unique cloners
- referring sites
- popular paths
- stars, forks, watchers
- Search Console impressions, clicks, indexed pages, and average position for hawp.online
- branded query visibility for HAWP, Human-AI Workflow Protocol, and HAWP AI agents

Keep the launch window clean: publish the website, then set the repository
description, topics, and homepage together so the before/after window has a
clear starting point. Record a fresh snapshot at +7 days and +14 days using
the same GitHub views and periods where possible.
