# TASK-036 Status

## Summary

TASK-036 added automatic repo-side syncing for `core/distribution/generated/**`, a one-step local `distribution:sync` command, and execution-first wording in the shared install/update distribution docs.

## Verified

- `npm --prefix librarian run distribution:sync` rebuilt and validated generated outputs.
- Generated branch guides now include `Install Work Item Contract` or `Update Work Item Contract` near the top.
- Generated guide footers now tell maintainers to rely on automatic sync or local `distribution:sync` instead of hand-editing generated files.

## Not Yet Verified

- The GitHub Actions workflow in `.github/workflows/sync-distribution-generated.yml` still needs one real push on `main` or `dev` to prove end-to-end auto-commit behavior.

## Follow-up Note

If the first live workflow run exposes branch-permission or token-policy issues, adjust only `.github/workflows/sync-distribution-generated.yml`; the generated-guide content changes are already locally verified.