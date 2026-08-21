# c1d2e401 — Fix CLI capability splits: extract run.go into per-group cmd_*.go files

**Type:** refactor
**Status:** plan-ready
**Updated:** 2026-08-21

## Goal

Split `librarian/go/internal/platform/cli/run.go` (~1226 lines) into per-capability-group files so each group is independently readable and testable. Leave `run.go` as a thin router (`Run` function + `ExitError` type only).

## Input

- Audit findings from `c1d2e400`: [plan](c1d2e400-audit-cli-capabilities.md)

## Target layout

| File | Contents |
| ---- | -------- |
| `run.go` | `Run` router + `ExitError` only |
| `cmd_links.go` | `runLinksCheck`, `runLinksClean` |
| `cmd_kit.go` | `runKitValidate`, `runKitNormalize` |
| `cmd_work.go` | `runWorkValidate`, `runWorkNormalize`, `runWorkNew` |
| `cmd_update.go` | update functions + `parseProviderFlags`, `doKitSync`, `sortedProviderNames` |
| `cmd_util.go` | `runUUID`, `runCheck`, `runInit`, `mustGetwd`, `helpText` |
| `cmd_index.go` | `runIndexBuild`, `runModelPull`, `runEmbed`, `modelsRoot` |
| `cmd_search.go` | all search functions + helpers |

## Split order

Follow the 8-step order in `c1d2e400`'s `## Handoff To c1d2e401` section. Verify `go build ./... && go vet ./... && go test ./internal/platform/cli/...` after each step.

## Known seams to defer

- `benchmark.go` direct sqlite access — document as follow-up, do not address here
- `download.NewHTTPFetcher()` / `githubrelease.NewClient()` constructor duplication — note as tech debt, defer to a subsequent refactor

## Constraints

- No behavior changes — this is a mechanical file split only
- Each step leaves the build green
- Do not merge with any feature or bug work
