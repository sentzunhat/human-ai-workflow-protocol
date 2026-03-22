# embeddings

Text-to-vector embedding: the `Embedder` port plus one adapter per backend
(ONNX, Ollama, "none"). Pure domain code — no CLI, config file, or HTTP
handler concerns live here, only the interface and its implementations.

## Exports

`Embedder` (port), `NewEmbedder` / `NewEmbedderWithURL` (factory),
`ONNXEmbedder`, `OllamaEmbedder`, `NullEmbedder` (adapters), `EmbeddingResult`.

## Backends

| Backend | Adapter | Network/model dependency |
|---|---|---|
| `onnx` | `ONNXEmbedder` | Local model file, no network after first download |
| `ollama` | `OllamaEmbedder` | Local/remote Ollama server over HTTP |
| `none` | `NullEmbedder` | Zero — returns empty vectors, for structured-output-without-a-model use |

## ONNX models

`bge-base-en-v1.5` (default, 768d), `all-MiniLM-L6-v2` (384d, smaller/faster),
`mdbr-leaf-ir` (384d — MongoDB's retrieval-tuned model; fine-tuned
specifically for retrieval quality rather than general-purpose similarity).
`mdbr-leaf-ir` is what "RAG retrieval" means in HAWP: a better embedding
model choice, not a separate IR subsystem — search and RAG share the same
`Embedder` (see `internal/application/search` and
`internal/application/context/rag.go`).

`bge-large-en-v1.5` was attempted and is deliberately **not** included —
see the code comment in `onnx_embedder.go` for why (an upstream download
deadlock, not a code fix).

## Quick use

```go
embedder, err := embeddings.NewEmbedder("onnx", "bge-base-en-v1.5")
if err != nil {
    log.Fatal(err)
}
defer embedder.Close()

vec, err := embedder.Embed(ctx, "some text")
```

## Adding a new backend

1. Implement `Embedder` in a new file, one adapter per file
   (`xyz_embedder.go`), matching the existing `onnx_embedder.go` /
   `ollama_embedder.go` shape.
2. Add a case to both `NewEmbedder` and `NewEmbedderWithURL` in `embedder.go`.
3. Add it to `context.BackendCategory` (`internal/application/context/config.go`)
   so it's classified as `offline`/`online`/`none`.
