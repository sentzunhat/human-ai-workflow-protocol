# Embedding port

This domain package owns the `Embedder` interface, `EmbeddingResult` value, and
`DefaultModel` constant. It contains no backend implementations or factories.

Paths below are repository-relative and remain valid when this port is moved by
the source-layout tool to `librarian/src/internal/domain/providers/embeddings/`.

## Implementations and wiring

- Backend adapters: `librarian/src/internal/infrastructure/models/{onnx,ollama,none}/`.
- Factory selection: `librarian/src/internal/infrastructure/models/factory.go`.
- Service composition: `librarian/src/internal/bootstrap/models.go`.
- Consumers: application search, index embedding, and context services receive
  the port through injected factories.

To add a backend, implement this interface in infrastructure, extend the model
factory there, and test it beside the adapter. Keep concrete construction out of
domain and application production code. Model availability and native-library
requirements belong to the implementation, not this interface package.
