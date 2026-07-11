# Go Librarian Model Strategy

Related work:

- `.hawp/work/active/1a3b32a4-ab37-4c86-ade7-71e2eb42b440.md`
- `.hawp/work/active/54e68af7-4622-4383-8482-cc4d4e1e21ee.md`
- `.hawp/work/active/8dedc4e2-69c5-42ca-aaad-93b62d7fb899.md`

## Decision Direction

Use a small Go librarian binary for the HAWP CLI/library direction. Do not embed
models in the binary. Download, verify, cache, and load models after install.

Primary product job:

- ingest HAWP work records and kit guidance
- make them easy to search lexically and semantically
- summarize the most relevant chunks locally
- hand stronger, smaller, more targeted context to larger AI models

This makes HAWP not just a folder convention, but a local intelligence layer
that improves AI-agent usage without requiring cloud retrieval infrastructure.

Default runtime path:

- CPU-only ONNX Runtime first.
- Go orchestrates repo/chunk parallelism with worker pools.
- ONNX Runtime handles per-session intra-op and inter-op CPU threading.
- GPU/WebGPU stays out of scope until CPU behavior is proven useful.

## Direct Evidence

- `onnxruntime_go` exposes `SessionOptions.SetIntraOpNumThreads` and
  `SessionOptions.SetInterOpNumThreads`.
- ONNX Runtime documentation describes intra-op threads for parallel execution
  inside graph operators and inter-op threads for parallel execution across
  graph nodes.
- `onnx-models/all-MiniLM-L6-v2-onnx` is an ONNX embedding model with 384
  dimensions, max sequence length 256, and documented file size around 0.08 GB.
- `sentence-transformers/all-MiniLM-L6-v2` maps sentences and paragraphs to a
  384-dimensional vector space for semantic search and clustering.

## Inference

- A 20 MB-ish Go librarian binary is plausible only if model files and native
  ONNX Runtime libraries are external/cache-managed.
- The right first-class librarian story is not "local chatbot." It is
  "structure HAWP data, search it well, and compress the right context for the
  next agent/model step."
- Embeddings are the right first ML capability for HAWP because vector search
  benefits from compact sentence models and does not require autoregressive
  generation loops.
- Text2text/summarization via ONNX is possible but should not be the first local
  CPU feature. It has more tokenizer/decoder/runtime complexity than embeddings.
- Default summaries should start as extractive summaries over ranked chunks;
  generative summaries can be optional later through a model provider.

## Model Cache Shape

Recommended default:

```text
~/.hawp/
  models/
    registry.json
    onnx/
      all-minilm-l6-v2/
        model.onnx
        tokenizer.json
        config.json
        manifest.json
  cache/
    downloads/
```

Every model should have a manifest:

```json
{
  "id": "all-minilm-l6-v2",
  "provider": "huggingface",
  "task": "embedding",
  "runtime": "onnxruntime",
  "executionProvider": "cpu",
  "embeddingDimensions": 384,
  "maxSequenceLength": 256,
  "files": [
    {
      "path": "model.onnx",
      "sha256": "<filled by downloader>"
    }
  ]
}
```

## CLI Shape

Initial model commands:

```bash
hawp models list
hawp models install all-minilm-l6-v2
hawp models status
hawp models remove all-minilm-l6-v2
hawp models path
```

Initial database/index commands:

```bash
hawp db init
hawp db migrate
hawp db status
hawp index build --scope work
hawp index build --scope kit
hawp index build --scope all
hawp index rebuild --scope all
```

Initial index/search commands:

```bash
hawp index build --mode lexical
hawp index build --mode embeddings --model all-minilm-l6-v2
hawp search "provider install contract"
hawp search "provider install contract" --semantic
hawp summarize .hawp/work --mode extractive
```

Prompt/context commands:

```bash
hawp context work 8dedc4e2
hawp context search "vector search summarization"
hawp context build --query "provider install contract"
hawp prompt handoff --id 8dedc4e2
```

The important distinction:

- `search` is for humans and tools asking questions.
- `context` is for building compact, ranked evidence packets for a larger model.
- `prompt handoff` is for producing model-ready instruction/context text from
  the best matching chunks.

## Go Runtime Shape

Suggested packages:

```text
cmd/hawp/
internal/db/
internal/db/migrations/
internal/ingest/work/
internal/ingest/kit/
internal/chunks/
internal/models/registry/
internal/models/download/
internal/models/cache/
internal/inference/onnx/
internal/inference/embedding/
internal/index/lexical/
internal/index/vector/
internal/search/
internal/summarize/
internal/context/
internal/prompt/
```

Suggested data flow:

```text
work/ + kit/
  -> ingest
  -> normalize metadata
  -> chunk content
  -> store docs/chunks in db
  -> optional embeddings
  -> lexical + vector search
  -> ranked chunk selection
  -> local extractive summary / handoff prompt
```

Threading/concurrency policy:

- Use a bounded Go worker pool for file parsing, chunking, and batch assembly.
- Use one or a small number of ONNX Runtime sessions per model, not one session
  per goroutine by default.
- Expose settings:
  - `models.onnx.intraOpThreads`
  - `models.onnx.interOpThreads`
  - `models.embedding.batchSize`
  - `index.workers`
- Default ONNX thread values can start at `0` so ONNX Runtime chooses, then tune
  after measuring.
- Avoid oversubscription: if Go workers are high, cap ONNX threads; if ONNX
  threads are high, keep embedding workers low.

## Settings Shape

Example user config:

```toml
[models]
directory = "~/.hawp/models"

[models.onnx]
execution_provider = "cpu"
intra_op_threads = 0
inter_op_threads = 0

[models.embedding]
default_model = "all-minilm-l6-v2"
batch_size = 16

[index]
workers = 0

[db]
path = "~/.hawp/index/librarian.db"

[context]
default_top_k = 12
max_chunk_chars = 1200
```

## Database Shape

Likely initial choice:

- SQLite for metadata and chunk storage
- SQLite FTS5 for lexical search
- vector storage in SQLite or adjacent files, depending the first practical Go
  library path

Initial records:

- `documents`
  - source type: `work` | `kit`
  - repo-relative path
  - title
  - modified timestamp
  - content hash
- `chunks`
  - document id
  - chunk index
  - text
  - token/char estimates
- `embeddings`
  - chunk id
  - model id
  - vector payload or pointer

`hawp db init` should:

- create the database
- apply migrations
- create lexical indexes
- prepare model/cache directories

`hawp index build --scope all` should:

- scan `.hawp/work/` and `.hawp/kit/`
- normalize and chunk content
- upsert documents/chunks
- optionally generate embeddings when a model is installed

## Feature Order

1. Go CLI + `db init` + lexical index/search over `work/` and `kit/`.
2. Model registry, download, checksum, cache, and settings.
3. ONNX Runtime CPU embedding proof with `all-MiniLM-L6-v2`.
4. Vector index/search over HAWP work/docs.
5. `context` and `prompt handoff` commands that rank and compress chunks for a
   larger model.
6. Extractive summarization from top-ranked chunks.
7. Optional text2text provider only after embeddings are stable.

## Open Questions

- Whether to ship ONNX Runtime native libraries with the Go binary, download
  them into the same model/runtime cache, or require users to install them.
- Whether the first vector store is SQLite vector extension, embedded Go index,
  or a simple file-backed ANN structure.
- Which tokenizer implementation gives the cleanest Go path for the selected
  ONNX embedding model.
- Whether text2text belongs in the Go librarian at all, or should remain an
  optional provider/sidecar.
- Whether `db init` should auto-install a default embedding model or keep model
  install explicit.
