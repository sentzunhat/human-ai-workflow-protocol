# Provider Harness — Embedding and LLM Backends

Every provider adapter (embedding backend or LLM client) must satisfy this
harness before being wired into the application. Apply the slice harness
(`slice-harness.md`) alongside this document.

## Interfaces

**Embedder** (`internal/domain/providers/embeddings`):

```go
type Embedder interface {
    Embed(ctx context.Context, text string) ([]float32, error)
    EmbedBatch(ctx context.Context, texts []string) ([][]float32, error)
    Dimension() int
    Backend() string
    Model() string
    Close() error
}
```

**LLMClient** (`internal/domain/providers/llm`):

```go
type LLMClient interface {
    Reshape(ctx context.Context, packedContext string, maxTokens int) (string, error)
    ReshapeBatch(ctx context.Context, contexts []string, maxTokens int) ([]string, error)
    Backend() string
    Model() string
    Close() error
}
```

## Required verifications (all must be recorded before wiring)

### 1. Interface compliance

The adapter type must satisfy the interface at compile time. Assert this
with a blank assignment in the adapter package:

```go
var _ embeddings.Embedder = (*MyEmbedder)(nil)
// or
var _ llm.LLMClient = (*MyLLMClient)(nil)
```

This check is compile-time and produces no runtime output. Record that it
compiles without error.

### 2. Dimension consistency

`Embed` and `EmbedBatch` must return vectors whose length equals
`Dimension()`. Verify with a table-driven test:

```go
func TestDimensionConsistency(t *testing.T) {
    e := newTestEmbedder(t)
    dim := e.Dimension()

    vec, err := e.Embed(context.Background(), "hello")
    if err != nil { t.Fatalf("Embed: %v", err) }
    if len(vec) != dim {
        t.Errorf("Embed returned %d dims, want %d", len(vec), dim)
    }

    vecs, err := e.EmbedBatch(context.Background(), []string{"a", "b"})
    if err != nil { t.Fatalf("EmbedBatch: %v", err) }
    for i, v := range vecs {
        if len(v) != dim {
            t.Errorf("EmbedBatch[%d] returned %d dims, want %d", i, len(v), dim)
        }
    }
}
```

### 3. Backend() and Model() non-empty

```go
if e.Backend() == "" { t.Error("Backend() must not be empty") }
if e.Model() == ""   { t.Error("Model() must not be empty") }
```

### 4. Close() idempotent

Close must not panic when called more than once:

```go
if err := e.Close(); err != nil { t.Logf("first Close: %v", err) }
if err := e.Close(); err != nil { t.Logf("second Close: %v", err) }
// neither call may panic
```

### 5. Backend-unavailable error path

When the backend is not reachable, the adapter must return a meaningful
error — not panic, not silent failure:

```go
// configure adapter to point at an invalid endpoint or missing file
e := newUnavailableEmbedder(t)
_, err := e.Embed(context.Background(), "test")
if err == nil {
    t.Error("expected error when backend unavailable, got nil")
}
```

## Test table structure

Minimum: 3 success cases, 4 error/rejection cases.

```go
func TestEmbed(t *testing.T) {
    cases := []struct {
        name    string
        text    string
        wantErr bool
    }{
        {name: "short text",   text: "hello",              wantErr: false},
        {name: "longer text",  text: "the quick brown fox", wantErr: false},
        {name: "empty text",   text: "",                   wantErr: false}, // backend-specific
        {name: "backend down", text: "x",                  wantErr: true},  // use unavailable backend
        // add backend-specific rejection cases
    }
    // ...
}
```

## Factory wiring

Factories live in `infrastructure/models/factory.go`. The domain layer
never constructs backends directly. When adding a new backend:

1. Add the adapter under `infrastructure/models/<name>/`.
2. Add a `case "<name>":` branch to `NewEmbedderWithURL` or
   `NewLLMClientWithURL` in `factory.go`.
3. Add the interface compliance assertion to the adapter package.
4. Run and record: `go build ./...`, `go vet ./...`, focused package test.

## Performance baseline (record actuals, not claims)

| Backend  | Latency (warm, single embed) | Note |
|----------|------------------------------|------|
| Lexical  | <1 ms                        | baseline for comparison |
| ONNX     | ~8 ms warm / ~476 ms cold    | all-MiniLM-L6-v2 |
| Ollama   | ~32 ms warm batch / ~110 ms cold | nomic-embed-text |

Record actual latency when testing a new backend. Do not assert specific
numbers without running a benchmark and recording the output.
