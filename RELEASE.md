# Release Playbook

HAWP releases are cut from `main` with plain semantic-version tags such as
`0.0.24`. Do not use a `v` prefix.

The current release workflow is `.github/workflows/release.yml`. It builds the
standard `hawp` binaries for all six platforms, optional ORT tarballs where
supported, `hawp-kit-bundle.tar.gz`, and one `checksums.txt` file.

## Prepare

1. Update `librarian/src/internal/domain/update/version.go`.
2. Add a matching section to `librarian/CHANGELOG.md`.
3. Run the local release checks from `librarian/src`.

```bash
cd librarian/src
go test ./...
go run ./cmd/hawp providers sync
go run ./cmd/hawp distribution sync
go run ./cmd/hawp kit validate
go run ./cmd/hawp work validate
go run ./cmd/hawp check
make dist VERSION=<version>
```

Commit the version, changelog, generated distribution files, and materialized
provider overlays that belong to the release lane. `.hawp/bin/` is local build
and install output; no binary is committed to the source repository. For a
source checkout, `make install VERSION=<version>` builds the host executable at
`.hawp/bin/hawp`. The generated install and update scripts select the host
platform's executable from the latest published GitHub release and verify it
against `checksums.txt` before replacing the local executable. Until `0.0.24`
is published, those scripts install the latest published version instead.

## Publish

Preferred path: merge the prepared release branch to `main`. The
`tag-on-merge.yml` workflow reads `version.go` and dispatches `release.yml`
unless that plain tag already exists.
Do not create the `0.0.24` tag before merging to `main`.
The release workflow builds the standard platform binaries, prepares the kit
bundle and checksums, and extracts release notes before it creates the tag.
If publication fails after tagging, rerun that release workflow run: it accepts
an existing tag only when it points to the same commit. A later `main` push
with the same version does not dispatch another release.

Manual path:

```bash
VERSION=<version>
gh workflow run release.yml --field version="$VERSION" --field draft=false
```

Or use GitHub Actions manual dispatch:

1. Open **Actions**.
2. Run **Release**.
3. Enter the release version being published.
4. Leave `draft` off for immediate publication, or enable it for manual review.

## Verify

After the workflow finishes, confirm the release page has:

- `hawp-darwin-amd64`
- `hawp-darwin-arm64`
- `hawp-linux-amd64`
- `hawp-linux-arm64`
- `hawp-windows-amd64.exe`
- `hawp-windows-arm64.exe`
- `hawp-kit-bundle.tar.gz`
- `checksums.txt`

Optional ORT tarballs may also be present for supported platforms.

Check one downloaded binary:

```bash
VERSION=<version>
curl -L -o hawp-darwin-arm64 "https://github.com/sentzunhat/human-ai-workflow-protocol/releases/download/${VERSION}/hawp-darwin-arm64"
curl -L -o checksums.txt "https://github.com/sentzunhat/human-ai-workflow-protocol/releases/download/${VERSION}/checksums.txt"
grep ' hawp-darwin-arm64$' checksums.txt | shasum -a 256 -c -
chmod +x hawp-darwin-arm64
./hawp-darwin-arm64 version
```

Expected version output: the value of `VERSION`.

## Downstream Provider Updates

Downstream repositories should get a normal branch and PR. For a stable update,
use the provider update guide from `main`.

For Claude/Codex provider and binary staging, use a temporary branch when you
need to test a slash-named ref. Run the visible install/update command block
from a generated guide, review it first, then set `REF` to the branch name you
are testing. The script archive extraction supports slash-named refs.
