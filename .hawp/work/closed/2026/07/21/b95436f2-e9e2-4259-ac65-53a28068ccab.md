# draft-and-approve release workflow: workflow_dispatch cuts the tag, builds, lands as draft; hawp update only sees published releases

**Backlog ID (Legacy):** — (UUID-native item)
**UUID:** `b95436f2-e9e2-4259-ac65-53a28068ccab`
**Type:** feature
**Reported:** 2026-07-21
**Risk Level:** low

---

### Input (what was reported)

> Want a workflow that cuts the tag, builds the binary, then drafts the
> release for approval — `hawp update` should only pick up the release
> after it's approved/published, checking against the current version.
>
> Follow-up same day: approval doesn't need to be mandatory — auto-publish
> is fine as the default; drafting should be an opt-in per-run choice, not
> a forced step every release.

---

### Context

Today's `release-librarian-go.yml` (from `cdcf9f78`) only triggers on a
manually-pushed tag and publishes immediately — no tagging automation,
no review gate before a release becomes live. This item adds both.

---

### Analysis

**Root cause (or most likely cause):**
_Manually running `git tag` + `git push` locally is exactly the step
that should be automatable from the GitHub UI, and publishing
immediately on tag push means there's no chance to review the built
binaries/changelog before `hawp update` can pull them._

**Directly verified (empirical test, 2026-07-21):**
_Created a real draft release via `gh api ... -F draft=true` on this
repo, then queried the unauthenticated `/repos/.../releases` list
endpoint (the same call `hawp update`'s `githubrelease.Client` makes) —
the draft did not appear; only the two already-published `v0.0.1`/
`v0.0.2` releases did. Confirms GitHub hides drafts from unauthenticated
requests, which is the entire approval-gate mechanism this item relies
on. Test release deleted after confirming._

**Inferred (not yet proven):**
_A `workflow_dispatch` input can drive both the git tag creation and the
subsequent build/draft steps without duplicating logic across two
trigger paths._

**Scope — what else is affected:**
_`.github/workflows/release-librarian-go.yml` only; no `hawp update`
client changes needed — the existing unauthenticated `Client.Latest`
already can't see drafts, which is exactly the desired behavior._

---

### Recommended Fix

- Add a `workflow_dispatch` trigger with a required `version` input
  (e.g. `0.1.0`); resolve the tag as `librarian-go-v<version>`.
- When dispatched manually, have the workflow create and push the
  annotated tag itself (`git tag` + `git push`) before building — so
  cutting a release no longer requires a local `git tag`/`git push`.
- Keep the existing tag-push trigger working unchanged (still supported
  for anyone who prefers tagging locally).
- Require a non-empty `CHANGELOG.md` section for the version — fail the
  job instead of silently falling back to a placeholder body, so a
  release can't go out without release notes.
- Add an optional `draft` boolean input (default `false`). Default
  behavior is auto-publish immediately, matching the "approval doesn't
  need to be mandatory" follow-up; setting `draft: true` on a dispatch
  run requires a maintainer to review the built binaries/changelog on
  GitHub and click "Publish release" before it becomes visible to
  `hawp update`. Tag-push triggers (no inputs available) always
  auto-publish.
- Document the new flow in `librarian/go/README.md`.

**What to verify after:**

- [x] `workflow_dispatch` with a version input creates the tag, builds,
      and publishes (or drafts, when requested)
      (Evidence: real dispatch run 29868596357 for `v0.0.3` — created
      and pushed the tag, built 6 binaries, published a non-draft
      prerelease with the correct changelog body)
- [x] The draft is invisible to `hawp update` (empirically confirmed for
      drafts in general via a real test release on this repo)
      (Evidence: dispatch run 29868859182 with `draft=true` for `v0.0.4`
      — release confirmed `isDraft: true`; a real `v0.0.3` binary ran
      `hawp update --check` against it and reported "Already up to
      date", correctly unable to see the draft. Draft + tag deleted
      after confirming.)
- [x] A version with no matching `CHANGELOG.md` section fails the workflow
      rather than publishing a placeholder-body release
      (Evidence: dispatch run 29869090753 for bogus version `99.99.99`
      — the changelog-extraction step failed with the expected error
      and the job stopped before the publish step ever ran;
      `gh release list` confirms no release was created. Stray tag
      deleted after confirming.)

---

## Outcome (filled at close)

Closed 2026-07-21. `.github/workflows/release-librarian-go.yml` now
supports cutting a release entirely from the GitHub Actions UI:
`workflow_dispatch` with a `version` input creates and pushes the tag,
builds all six platform binaries, and publishes. Publishing is
auto-approved by default (per the same-day follow-up that mandatory
approval wasn't wanted); an optional `draft` input makes a given release
require manual review and "Publish release" before `hawp update` can see
it — no client-side code needed for this gate, since GitHub already
hides drafts from the unauthenticated requests `hawp update` makes. The
workflow also now hard-fails without a matching `CHANGELOG.md` section,
so notes can never be skipped.

All three behaviors were proven with real dispatch runs against this
repo, not just design review: auto-publish (`v0.0.3`, now the real
`hawp update` target with genuine new functionality — folder-context
`index build`), the draft gate (`v0.0.4`, confirmed invisible to a real
binary then cleaned up), and the changelog guard (bogus `99.99.99`
version, confirmed the job fails before any release is created). A
minor doc-comment inaccuracy in `githubrelease.Client.Latest` (claimed
drafts were included in the list results) was corrected to match this
verified behavior.

## Verification (filled at close)

- Evidence: dispatch run 29868596357 (`v0.0.3`, auto-publish) — tag
  `librarian-go-v0.0.3` created by the workflow, release published
  non-draft, prerelease, correct changelog body; real `hawp update`
  from a `v0.0.2` binary found and applied it successfully.
- Evidence: dispatch run 29868859182 (`v0.0.4`, `draft=true`) — release
  confirmed draft via `gh release view --json isDraft`; a real `v0.0.3`
  binary's `hawp update --check` reported up to date, unable to see it.
- Evidence: dispatch run 29869090753 (`99.99.99`, no changelog section)
  — job failed at the changelog-extraction step with the expected
  error; `gh release list` shows no release was created for it.
- Evidence: `make check` (vet + tests + build) passes after the
  `githubrelease` doc-comment correction.

## Close Checklist

- [x] Outcome section filled
- [x] Verification section filled
- [x] Plan file saved under `closed/2026/07/21/b95436f2-e9e2-4259-ac65-53a28068ccab.md`
- [x] BACKLOG.md updated
