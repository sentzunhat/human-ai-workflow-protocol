# Source layout preview

Generated from the current source snapshot. **Proposed, not applied.**

273 files inventoried; 0 proposed moves; 0 content updates (including moved files); 273 retained paths. No dead-file deletions.

This report records destinations, not proof that every layer is already pure. Run `--check` separately for candidate compilation and vet; run `--diff` to inspect exact code changes. All file paths below are relative to `librarian/src`.

## Proposed directory tree

Counts are files directly in each directory, not recursive totals. Empty folders are not created.

```text
librarian/src/
  (5 root files)
  cmd/ (0 files)
    hawp/ (1 files)
  internal/ (0 files)
    application/ (0 files)
      check/ (1 files)
      context/ (11 files)
        dedup/ (2 files)
      db/ (1 files)
      distribution/ (1 files)
      embed/ (1 files)
      index/ (8 files)
      kit/ (2 files)
      kitsync/ (1 files)
      links/ (2 files)
      providersync/ (1 files)
      provision/ (1 files)
      search/ (2 files)
      update/ (2 files)
      usage/ (2 files)
      uuidgen/ (1 files)
      work/ (1 files)
        intake/ (3 files)
        normalize/ (3 files)
        validation/ (1 files)
    bootstrap/ (1 files)
    domain/ (0 files)
      context/ (6 files)
      distribution/ (4 files)
      index/ (3 files)
      kit/ (6 files)
      kitsync/ (5 files)
      providers/ (0 files)
        embeddings/ (2 files)
        llm/ (2 files)
      providersync/ (2 files)
      provision/ (3 files)
      search/ (3 files)
      update/ (2 files)
      usage/ (1 files)
      work/ (25 files)
    infrastructure/ (0 files)
      archive/ (3 files)
      clients/ (0 files)
        download/ (2 files)
        githubrelease/ (1 files)
      filesystem/ (4 files)
      markdown/ (1 files)
      models/ (4 files)
        none/ (2 files)
        ollama/ (4 files)
        onnx/ (4 files)
      repo/ (2 files)
      repositories/ (0 files)
        context/ (1 files)
        distribution/ (2 files)
        index/ (6 files)
        kit/ (2 files)
        kitsync/ (2 files)
        provision/ (1 files)
        usage/ (3 files)
        work/ (2 files)
      selfreplace/ (2 files)
      tomlconfig/ (3 files)
    platform/ (0 files)
      cli/ (8 files)
        distribution/ (1 files)
        index/ (1 files)
          build/ (2 files)
          ingest/ (3 files)
        init/ (2 files)
        kit/ (2 files)
          normalize/ (3 files)
          validate/ (3 files)
        links/ (1 files)
          check/ (1 files)
          clean/ (3 files)
        mcp/ (2 files)
          configure/ (2 files)
        model/ (1 files)
          embed/ (3 files)
          location/ (1 files)
          pull/ (3 files)
          search-embed/ (2 files)
        providers/ (1 files)
        search/ (1 files)
          benchmark/ (2 files)
          query/ (3 files)
        update/ (2 files)
        usage/ (1 files)
          enable/ (3 files)
        work/ (1 files)
          new/ (4 files)
          normalize/ (3 files)
          validate/ (3 files)
      exitcode/ (1 files)
      mcp/ (0 files)
        configure/ (11 files)
        server/ (7 files)
  tests/ (1 files)
    application/ (0 files)
      db/ (1 files)
      links/ (2 files)
      provision/ (1 files)
      search/ (2 files)
      update/ (1 files)
      uuidgen/ (1 files)
      work/ (2 files)
    domain/ (0 files)
      search/ (1 files)
      update/ (1 files)
    infrastructure/ (0 files)
      clients/ (0 files)
        githubrelease/ (1 files)
      filesystem/ (1 files)
    platform/ (0 files)
      cli/ (2 files)
```

## Moves and ownership

| Current file | Proposed file | Reason | Content changes |
| --- | --- | --- | --- |

## Content updates without moves

These callers stay in their current folders; only their imports/references change.


## Unchanged files

Retained with their existing owner; inclusion here does not claim complete architectural purity. Tests and assets are not treated as dead files.

- `CHANGELOG.md`
- `Makefile`
- `README.md`
- `cmd/hawp/main.go`
- `go.mod`
- `go.sum`
- `internal/application/check/check.go`
- `internal/application/context/README.md`
- `internal/application/context/config.go`
- `internal/application/context/config_test.go`
- `internal/application/context/dedup/dedup.go`
- `internal/application/context/dedup/dedup_test.go`
- `internal/application/context/encryption.go`
- `internal/application/context/encryption_test.go`
- `internal/application/context/format.go`
- `internal/application/context/format_test.go`
- `internal/application/context/rag.go`
- `internal/application/context/rag_test.go`
- `internal/application/context/reshaper.go`
- `internal/application/context/reshaper_test.go`
- `internal/application/db/init-service.go`
- `internal/application/distribution/distribution.go`
- `internal/application/embed/embed.go`
- `internal/application/index/build-service.go`
- `internal/application/index/build-service_test.go`
- `internal/application/index/embed-service.go`
- `internal/application/index/embed-service_test.go`
- `internal/application/index/embed_integration_test.go`
- `internal/application/index/ingest-service.go`
- `internal/application/index/ingest-service_test.go`
- `internal/application/index/ingest_e2e_test.go`
- `internal/application/kit/normalize.go`
- `internal/application/kit/validate.go`
- `internal/application/kitsync/kitsync.go`
- `internal/application/links/README.md`
- `internal/application/links/check.go`
- `internal/application/providersync/providersync.go`
- `internal/application/provision/provision.go`
- `internal/application/search/config.go`
- `internal/application/search/service.go`
- `internal/application/update/notifier.go`
- `internal/application/update/update.go`
- `internal/application/usage/service.go`
- `internal/application/usage/service_test.go`
- `internal/application/uuidgen/uuid.go`
- `internal/application/work/README.md`
- `internal/application/work/intake/draft.go`
- `internal/application/work/intake/draft_test.go`
- `internal/application/work/intake/intake.go`
- `internal/application/work/normalize/duplicate_links.go`
- `internal/application/work/normalize/duplicate_links_test.go`
- `internal/application/work/normalize/normalize.go`
- `internal/application/work/validation/validate.go`
- `internal/bootstrap/models.go`
- `internal/domain/context/document.go`
- `internal/domain/context/kit.go`
- `internal/domain/context/kit_test.go`
- `internal/domain/context/source.go`
- `internal/domain/context/work.go`
- `internal/domain/context/work_test.go`
- `internal/domain/distribution/binary_download_mocks_test.go`
- `internal/domain/distribution/binary_download_test.go`
- `internal/domain/distribution/distribution.go`
- `internal/domain/distribution/distribution_test.go`
- `internal/domain/index/chunk.go`
- `internal/domain/index/chunk_test.go`
- `internal/domain/index/document-scope.go`
- `internal/domain/kit/links.go`
- `internal/domain/kit/normalize.go`
- `internal/domain/kit/normalize_test.go`
- `internal/domain/kit/source.go`
- `internal/domain/kit/validate.go`
- `internal/domain/kit/validate_test.go`
- `internal/domain/kitsync/apply.go`
- `internal/domain/kitsync/detect.go`
- `internal/domain/kitsync/filecopier.go`
- `internal/domain/kitsync/kitsync_test.go`
- `internal/domain/kitsync/manifest.go`
- `internal/domain/providers/embeddings/README.md`
- `internal/domain/providers/embeddings/embedder.go`
- `internal/domain/providers/llm/README.md`
- `internal/domain/providers/llm/llm_client.go`
- `internal/domain/providersync/materialize.go`
- `internal/domain/providersync/materialize_test.go`
- `internal/domain/provision/assets.go`
- `internal/domain/provision/manifest.go`
- `internal/domain/provision/manifest_test.go`
- `internal/domain/search/index.go`
- `internal/domain/search/result.go`
- `internal/domain/search/similarity.go`
- `internal/domain/update/asset.go`
- `internal/domain/update/version.go`
- `internal/domain/usage/types.go`
- `internal/domain/work/backlog.go`
- `internal/domain/work/backlog_test.go`
- `internal/domain/work/clarity.go`
- `internal/domain/work/completeness.go`
- `internal/domain/work/consistency.go`
- `internal/domain/work/constants.go`
- `internal/domain/work/deadlinks.go`
- `internal/domain/work/draft.go`
- `internal/domain/work/evidence.go`
- `internal/domain/work/idparse.go`
- `internal/domain/work/idparse_test.go`
- `internal/domain/work/intake.go`
- `internal/domain/work/intake_table.go`
- `internal/domain/work/intake_test.go`
- `internal/domain/work/links.go`
- `internal/domain/work/normalize_active_rows.go`
- `internal/domain/work/normalize_apply.go`
- `internal/domain/work/normalize_migrate.go`
- `internal/domain/work/normalize_report.go`
- `internal/domain/work/normalize_rules.go`
- `internal/domain/work/normalize_scan.go`
- `internal/domain/work/normalize_test.go`
- `internal/domain/work/types.go`
- `internal/domain/work/validations_test.go`
- `internal/domain/work/work_source.go`
- `internal/infrastructure/archive/extract.go`
- `internal/infrastructure/archive/extract_test.go`
- `internal/infrastructure/archive/targz.go`
- `internal/infrastructure/clients/download/download.go`
- `internal/infrastructure/clients/download/download_test.go`
- `internal/infrastructure/clients/githubrelease/githubrelease.go`
- `internal/infrastructure/filesystem/hawp_home.go`
- `internal/infrastructure/filesystem/hawp_project.go`
- `internal/infrastructure/filesystem/layout-service.go`
- `internal/infrastructure/filesystem/readme_generator.go`
- `internal/infrastructure/markdown/markdown.go`
- `internal/infrastructure/models/benchmark_test.go`
- `internal/infrastructure/models/factory.go`
- `internal/infrastructure/models/factory_test.go`
- `internal/infrastructure/models/integration_test.go`
- `internal/infrastructure/models/none/embedder.go`
- `internal/infrastructure/models/none/llm.go`
- `internal/infrastructure/models/ollama/embedder.go`
- `internal/infrastructure/models/ollama/embedder_test.go`
- `internal/infrastructure/models/ollama/llm.go`
- `internal/infrastructure/models/ollama/llm_test.go`
- `internal/infrastructure/models/onnx/embedder.go`
- `internal/infrastructure/models/onnx/embedder_test.go`
- `internal/infrastructure/models/onnx/llm.go`
- `internal/infrastructure/models/onnx/llm_test.go`
- `internal/infrastructure/repo/root.go`
- `internal/infrastructure/repo/worktree.go`
- `internal/infrastructure/repositories/context/reader.go`
- `internal/infrastructure/repositories/distribution/reader.go`
- `internal/infrastructure/repositories/distribution/reader_test.go`
- `internal/infrastructure/repositories/index/domain_models_test.go`
- `internal/infrastructure/repositories/index/embedding_metadata_test.go`
- `internal/infrastructure/repositories/index/fts_sync_test.go`
- `internal/infrastructure/repositories/index/index.go`
- `internal/infrastructure/repositories/index/index_test.go`
- `internal/infrastructure/repositories/index/transaction_test.go`
- `internal/infrastructure/repositories/kit/validate_reader.go`
- `internal/infrastructure/repositories/kit/validate_reader_test.go`
- `internal/infrastructure/repositories/kitsync/filecopy.go`
- `internal/infrastructure/repositories/kitsync/filecopy_test.go`
- `internal/infrastructure/repositories/provision/manifest_adapter.go`
- `internal/infrastructure/repositories/usage/config.go`
- `internal/infrastructure/repositories/usage/store.go`
- `internal/infrastructure/repositories/usage/store_test.go`
- `internal/infrastructure/repositories/work/backlog_reader.go`
- `internal/infrastructure/repositories/work/backlog_reader_test.go`
- `internal/infrastructure/selfreplace/selfreplace.go`
- `internal/infrastructure/selfreplace/selfreplace_test.go`
- `internal/infrastructure/tomlconfig/edits.go`
- `internal/infrastructure/tomlconfig/table.go`
- `internal/infrastructure/tomlconfig/table_test.go`
- `internal/platform/cli/commands.go`
- `internal/platform/cli/commands_test.go`
- `internal/platform/cli/distribution/commands.go`
- `internal/platform/cli/embed_test.go`
- `internal/platform/cli/help.go`
- `internal/platform/cli/index/build/args.go`
- `internal/platform/cli/index/build/command.go`
- `internal/platform/cli/index/index_commands.go`
- `internal/platform/cli/index/ingest/command.go`
- `internal/platform/cli/index/ingest/corpus.go`
- `internal/platform/cli/index/ingest/corpus_test.go`
- `internal/platform/cli/init/args_test.go`
- `internal/platform/cli/init/command.go`
- `internal/platform/cli/init_test.go`
- `internal/platform/cli/kit/commands.go`
- `internal/platform/cli/kit/mutation_boundary_test.go`
- `internal/platform/cli/kit/normalize/args.go`
- `internal/platform/cli/kit/normalize/args_test.go`
- `internal/platform/cli/kit/normalize/command.go`
- `internal/platform/cli/kit/validate/args.go`
- `internal/platform/cli/kit/validate/args_test.go`
- `internal/platform/cli/kit/validate/command.go`
- `internal/platform/cli/links/check/command.go`
- `internal/platform/cli/links/clean/args.go`
- `internal/platform/cli/links/clean/args_test.go`
- `internal/platform/cli/links/clean/command.go`
- `internal/platform/cli/links/commands.go`
- `internal/platform/cli/mcp/commands.go`
- `internal/platform/cli/mcp/commands_test.go`
- `internal/platform/cli/mcp/configure/command.go`
- `internal/platform/cli/mcp/configure/command_test.go`
- `internal/platform/cli/mcp_test.go`
- `internal/platform/cli/model/embed/args.go`
- `internal/platform/cli/model/embed/args_test.go`
- `internal/platform/cli/model/embed/command.go`
- `internal/platform/cli/model/location/location.go`
- `internal/platform/cli/model/model_commands.go`
- `internal/platform/cli/model/pull/args.go`
- `internal/platform/cli/model/pull/args_test.go`
- `internal/platform/cli/model/pull/command.go`
- `internal/platform/cli/model/search-embed/args.go`
- `internal/platform/cli/model/search-embed/command.go`
- `internal/platform/cli/providers/commands.go`
- `internal/platform/cli/registry.go`
- `internal/platform/cli/run.go`
- `internal/platform/cli/search/benchmark/args.go`
- `internal/platform/cli/search/benchmark/command.go`
- `internal/platform/cli/search/commands.go`
- `internal/platform/cli/search/query/args.go`
- `internal/platform/cli/search/query/args_test.go`
- `internal/platform/cli/search/query/command.go`
- `internal/platform/cli/update/args_test.go`
- `internal/platform/cli/update/commands.go`
- `internal/platform/cli/usage/commands.go`
- `internal/platform/cli/usage/enable/args.go`
- `internal/platform/cli/usage/enable/args_test.go`
- `internal/platform/cli/usage/enable/command.go`
- `internal/platform/cli/work/commands.go`
- `internal/platform/cli/work/new/args.go`
- `internal/platform/cli/work/new/args_test.go`
- `internal/platform/cli/work/new/command.go`
- `internal/platform/cli/work/new/mutation_boundary_test.go`
- `internal/platform/cli/work/normalize/args.go`
- `internal/platform/cli/work/normalize/args_test.go`
- `internal/platform/cli/work/normalize/command.go`
- `internal/platform/cli/work/validate/args.go`
- `internal/platform/cli/work/validate/args_test.go`
- `internal/platform/cli/work/validate/command.go`
- `internal/platform/exitcode/error.go`
- `internal/platform/mcp/configure/config.go`
- `internal/platform/mcp/configure/config_codex.go`
- `internal/platform/mcp/configure/config_codex_args.go`
- `internal/platform/mcp/configure/config_codex_test.go`
- `internal/platform/mcp/configure/config_json.go`
- `internal/platform/mcp/configure/config_json_test.go`
- `internal/platform/mcp/configure/config_providers.go`
- `internal/platform/mcp/configure/config_providers_test.go`
- `internal/platform/mcp/configure/config_test.go`
- `internal/platform/mcp/configure/configure.go`
- `internal/platform/mcp/configure/configure_test.go`
- `internal/platform/mcp/server/server.go`
- `internal/platform/mcp/server/tool_search.go`
- `internal/platform/mcp/server/tool_usage.go`
- `internal/platform/mcp/server/tool_work.go`
- `internal/platform/mcp/server/tools.go`
- `internal/platform/mcp/server/tools_e2e_test.go`
- `internal/platform/mcp/server/types.go`
- `tests/README.md`
- `tests/application/db/init_service_test.go`
- `tests/application/links/check_test.go`
- `tests/application/links/clean_test.go`
- `tests/application/provision/provision_test.go`
- `tests/application/search/repository_test.go`
- `tests/application/search/service_test.go`
- `tests/application/update/update_test.go`
- `tests/application/uuidgen/uuid_test.go`
- `tests/application/work/intake_test.go`
- `tests/application/work/normalize_test.go`
- `tests/domain/search/similarity_test.go`
- `tests/domain/update/version_test.go`
- `tests/infrastructure/clients/githubrelease/githubrelease_test.go`
- `tests/infrastructure/filesystem/hawp_project_test.go`
- `tests/platform/cli/run_test.go`
- `tests/platform/cli/search_test.go`

## Still requires semantic work

- domain/work normalization still has historical closed-record evidence review; folder normalization is clean
- domain/context, kit, kitsync, providersync and distribution still mix policy with I/O; imports alone cannot separate them
- application/context configuration and index storage wiring need consumer ports before further layer extraction
- MCP server tool handlers share private RPC contracts; keep them together until that contract is extracted
- No production file is deleted based on absent filename references; moved sources are removed only after verification
