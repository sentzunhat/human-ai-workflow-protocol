# context

Application-layer orchestration of the search → embed → reshape pipeline
(`hawp search --context --llm-reshape`). Wires the `embeddings`/`llm` domain
ports together; owns config loading, none of the actual embedding/generation
logic.

## Exports

`ContextBlock`, `FormattedResult`, `DocumentReference` (data shapes) ·
`FormatAsMarkdown` (dedupe + format search results) ·
`ContextReshaper` / `NewContextReshaper` / `ReshapingConfig` (embed + LLM
orchestration) · `RAGPipeline` / `DefaultRAGPipeline` / `NewDefaultRAGPipeline`
(the CLI-facing entry point — wraps `ContextReshaper`, adds reference
tracking) · `ContextConfig` / `LoadContextConfig` / `BackendCategory`
(config: file → env → CLI flag priority).

## Quick use

```go
config, err := context.LoadContextConfig(hawpHome, projectRoot)
pipeline, err := context.NewDefaultRAGPipeline(context.ReshapingConfig{
    EmbeddingsBackend: config.Embeddings.Backend,
    EmbeddingsModel:   config.Embeddings.Model,
    LLMBackend:        config.LLM.Backend,
    LLMModel:          config.LLM.Model,
})
defer pipeline.Close()

block := context.FormatAsMarkdown(searchResults, query, maxTokens)
output, err := pipeline.Reshape(ctx, block, maxTokens)
```

## Design notes

- **References render inline**, not batched — each result's compact `Ref:`
  line sits immediately above its own content (`format.go`'s
  `formatResultsInline`), not collected into a footer list.
- **Structured metadata stays typed**, while markdown stays lean — title,
  relevance, and matched excerpt still live on `FormattedResult` /
  `DocumentReference` and JSON output, but sparse markdown rendering keeps only
  the source+line provenance needed at read time.
- **The `"none"` backend is a deliberate mode**, not an error fallback —
  `ContextReshaper.Reshape` skips the embed/LLM stages entirely rather than
  routing a no-op call through them, so it's genuinely zero network.
- **`LLMClient` was deliberately not renamed to `Reshaper`** during the
  v0.1.0 push — see `librarian/docs/v0.1.0-vision.md`'s "Unified Interfaces"
  section for why.

See `librarian/docs/v0.1.0-vision.md` for the full architecture writeup.
