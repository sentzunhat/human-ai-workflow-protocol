---
uuid: f3dbf906-1327-4bed-9f42-8beaeb02143f
title: "v0.0.24 PR #40 archive extraction hardening and provider-installation boundary checkpoint"
type: status
date: 2026-09-21
---
# Status Report

#### Intent

Close the remaining Copilot-reported Markdown offset defect and the queued
v0.0.24 filesystem-boundary findings while preserving the existing PR #40
scope and keeping provider installation boundaries explicit.

#### Current State

The working tree contains the implementation and regression tests for all
findings inspected in this checkpoint. Existing commits `e48ad0c` and
`7d0b7aa` remain unchanged. No commit, push, merge, or new Copilot review was
requested or performed.

#### What Was Inspected

- `librarian/src/internal/domain/work/markdown/links.go` and all direct
  `BlankFences` callers.
- Search corpus path resolution and Markdown walkers.
- Kitsync domain `FileCopier` and repository filesystem adapter.
- MCP provider `.gitignore` writes and `links clean --apply` collection/write
  paths.
- Install/update distribution source templates and regenerated provider
  outputs.
- Existing HAWP active plans and the current dirty worktree.

#### What Changed

- `BlankFences` now masks bytes rather than runes, preserving byte offsets for
  non-ASCII fenced content.
- Search indexing rejects absolute/traversal paths and skips symlinked or
  non-regular Markdown entries.
- Kitsync checks destination ancestors before directory creation, temp-file
  creation, and rename operations.
- MCP `.gitignore` mutation rejects symlinked ancestors and non-regular files.
- Link cleanup collects regular Markdown files only and rechecks the write
  target before applying changes.
- Distribution source scripts reject `.hawp` symlinks, nested `.hawp`
  symlinks, and `.`/`..` backlog destination components; all 40 generated
  provider artifacts were regenerated and validated.

#### What Was Directly Verified

- Focused tests for Markdown offsets, indexing boundaries, kitsync symlinked
  destinations, and `.gitignore` symlinks passed.
- `go test ./...` passed.
- `go vet ./...` passed.
- `git diff --check` passed.
- `go run ./cmd/hawp distribution sync` regenerated and validated outputs.
- Generated shell syntax checks passed for representative outputs.
- `go run ./cmd/hawp check` passed.
- `go run ./cmd/hawp work validate` passed with one pre-existing verification
  clarity warning and no issues.

#### What Remains Unproven

- Fresh GitHub checks, Copilot review, and human approval for PR #40 remain
  external and were not refreshed.
- The shell guards were syntax-checked and generated-output validated, but no
  destructive install/update was run against a malicious fixture.
- The portable symlink checks still have the normal check/use race of path-
  based filesystem APIs; descriptor-relative no-follow operations were not
  introduced.
- Physical provider availability and live local-model readiness remain outside
  this checkpoint.

#### Constraints

Unrelated dirty HAWP work items and existing backlog changes were preserved.
The work stayed within the current PR #40 hardening lane; no planning-only
security work was moved into implementation.

#### Help Wanted

Request a fresh Copilot review and inspect the resulting diff before approval,
with particular attention to generated shell behavior and the remaining
path-based filesystem race boundary.

#### Suggested Next Step

Review the working-tree diff, then commit and push the compoundable fixes only
after explicit authorization. Request fresh CI/Copilot assessment before
calling PR #40 approved or merged.
