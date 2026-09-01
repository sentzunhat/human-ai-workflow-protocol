# Release Playbook

HAWP releases are cut from `main` with plain semantic-version tags such as
`0.0.23`. Do not use a `v` prefix.

The current release workflow is `.github/workflows/release.yml`. It builds the
standard `hawp` binaries for all six platforms, optional ORT tarballs where
supported, `hawp-kit-bundle.tar.gz`, and one `checksums.txt` file.

## Prepare

1. Update `librarian/src/internal/domain/update/version.go`.
2. Add a matching section to `librarian/src/CHANGELOG.md`.
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

Commit the version, changelog, generated distribution files, materialized
provider overlays, and checked-in `.hawp/bin/hawp` wrapper changes that belong
to the release lane.

## Publish

Preferred path: merge the prepared release branch to `main`. The
`tag-on-merge.yml` workflow reads `version.go` and dispatches `release.yml`
unless that plain tag already exists.

Manual path:

```bash
VERSION=<version>
git tag -a "$VERSION" -m "release $VERSION"
git push origin main "$VERSION"
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

For Claude/Codex provider and binary staging, the current branch observed in
this repository is:

```text
claude/claude-codex-binary-update-n00dve
```

When testing a slash-named branch, use the visible install/update command block
from a generated guide and set `REF` to that branch name after review. The
script archive extraction supports slash-named refs.
