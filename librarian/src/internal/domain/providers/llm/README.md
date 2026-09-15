# LLM port

This domain package owns the `LLMClient` interface and `ReshapingPrompt` constant.
It contains no backend implementations, model sessions, or factories.

Paths below are repository-relative and remain valid when this port is moved by
the source-layout tool to `librarian/src/internal/domain/providers/llm/`.

## Implementations and wiring

- Backend adapters: `librarian/src/internal/infrastructure/models/{ollama,onnx,none}/`.
- Factory selection: `librarian/src/internal/infrastructure/models/factory.go`.
- Context-service composition: `librarian/src/internal/bootstrap/models.go`.

To add a backend, implement `LLMClient` in infrastructure, extend the model
factory there, and add adapter tests. Application services receive injected
factories; domain code never imports infrastructure. Backend-specific defaults,
supported models, and native-library requirements belong to their adapters.
