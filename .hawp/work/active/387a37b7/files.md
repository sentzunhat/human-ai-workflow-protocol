# File Tracking — 387a37b7

## Owned Files

```txt
.hawp/work/BACKLOG.md
.hawp/work/active/387a37b7/plan.md
.hawp/work/active/387a37b7/files.md
librarian/src/internal/application/context/format.go
librarian/src/internal/application/context/format_test.go
librarian/src/internal/application/search/service.go
librarian/src/internal/application/search/service_test.go
librarian/src/internal/domain/search/result.go
librarian/src/internal/infrastructure/filesystem/hawp_project.go
librarian/src/internal/infrastructure/filesystem/hawp_project_test.go
librarian/src/internal/infrastructure/filesystem/readme_generator.go
```

## Read-Only Context Files

```txt
.hawp/kit/references/work-item-file-tracking.md
.hawp/kit/usage/intake-workflow.md
librarian/src/internal/application/context/config.go
librarian/src/internal/application/context/rag.go
librarian/src/internal/infrastructure/sqlite/index.go
librarian/src/internal/platform/cli/run.go
```

## Do-Not-Touch Files

```txt
.hawp/work/active/c804eec0/plan.md
.hawp/work/active/5957aaf4/plan.md
```

## Locked / Reserved Files

```txt
librarian/src/internal/application/context/format.go
librarian/src/internal/application/context/format_test.go
librarian/src/internal/application/search/service.go
librarian/src/internal/application/search/service_test.go
librarian/src/internal/domain/search/result.go
librarian/src/internal/infrastructure/filesystem/hawp_project.go
librarian/src/internal/infrastructure/filesystem/hawp_project_test.go
librarian/src/internal/infrastructure/filesystem/readme_generator.go
```

## Changed Files

```txt
.hawp/work/BACKLOG.md
.hawp/work/active/387a37b7/plan.md
.hawp/work/active/387a37b7/files.md
librarian/src/internal/application/context/format.go
librarian/src/internal/application/context/format_test.go
librarian/src/internal/application/search/service.go
librarian/src/internal/application/search/service_test.go
librarian/src/internal/domain/search/result.go
librarian/src/internal/infrastructure/filesystem/hawp_project.go
librarian/src/internal/infrastructure/filesystem/hawp_project_test.go
librarian/src/internal/infrastructure/filesystem/readme_generator.go
```

## Verification Notes

```txt
git diff --name-status
git diff --check
go test ./internal/application/search ./internal/application/context ./internal/infrastructure/filesystem
```
