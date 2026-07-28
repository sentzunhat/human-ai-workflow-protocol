# llm

Text generation: the `LLMClient` port plus one adapter per backend (Ollama,
ONNX, "none"). Pure domain code, same shape as `../embeddings`.

## Exports

`LLMClient` (port), `NewLLMClient` / `NewLLMClientWithURL` (factory),
`OllamaLLMClient`, `ONNXLLMClient`, `NullLLMClient` (adapters),
`ReshapingPrompt`, `SupportedModels`, `DefaultModel`.

## Backends

| Backend | Adapter | Status |
|---|---|---|
| `ollama` | `OllamaLLMClient` | Working, default |
| `onnx` | `ONNXLLMClient` | Verified working (SmolLM2-360M via hugot's CGO ORT backend) but requires a `-tags ORT` build with 3 native libraries — not part of the default build. See `librarian/docs/v0.1.0_VISION.md`'s "ONNX Text2Text Model" section for the full setup and reproduction steps. |
| `none` | `NullLLMClient` | Zero-cost passthrough (returns input unchanged) |

## Quick use

```go
client, err := llm.NewLLMClient("ollama", "mistral")
if err != nil {
    log.Fatal(err)
}
defer client.Close()

reshaped, err := client.Reshape(ctx, content, maxTokens)
```

## Adding a new backend

Same steps as `../embeddings/README.md` — one adapter file, one case in each
`New*` factory function, one entry in `context.BackendCategory`.
