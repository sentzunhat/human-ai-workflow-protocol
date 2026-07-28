# adopt hugot as a transformers.js-style Go wrapper for hawp (pull, embed, generate)

**Backlog ID (Legacy):** — (UUID-native item)
**UUID:** `748609a8-ac1a-4e27-a0dc-7a67d3b9fd14`
**Type:** feature
**Reported:** 2026-07-21
**Risk Level:** medium

---

### Input (what was reported)

> Want something like a transformers.js wrapper on the CLI librarian so
> it's easier to use for other things — generating text with small
> models or other Hugging Face models/modes.
>
> Follow-up: if a good open-source project already exists, use that
> instead of hand-rolling it.

---

### Context

This generalizes `fbf12a93` Slice 2 (which was scoped narrowly to
embedding the pinned MiniLM model for vector search) into a broader
local-model wrapper: pull arbitrary Hugging Face ONNX models and run
multiple pipeline types (embeddings, classification, and eventually
text generation) from the CLI — the actual "transformers.js for Go"
experience the input asked for.

---

### Analysis

**Root cause (or most likely cause):**
_Hand-rolling ONNX inference + tokenization in Go (cgo bindings,
tokenizer implementation, per-architecture decode loops) is a large,
error-prone undertaking that a maintained open-source project already
solves well._

**Directly verified:**
_[hugot](https://github.com/knights-analytics/hugot) (Apache 2.0) is
exactly this: HF-compatible ONNX pipelines in Go, with a HF Hub
downloader built in, and a pluggable backend architecture. Read its
README and source (`hugot.go`, `hugot_go.go`, `downloader.go`)
directly. Critically, it has a **pure-Go backend** (`NewGoSession`,
backed by GoMLX's `simplego`) that needs **no cgo and no external ONNX
Runtime binary** — confirmed empirically on this repo:
`CGO_ENABLED=0 go build`/`go test` of a package importing
`hugot.NewGoSession` succeeds, and `CGO_ENABLED=0 go build` for
darwin/amd64, linux/arm64, and windows/amd64 all succeed with the real
package linked in (not just an unlinked no-op). This means adopting
hugot's Go backend for embeddings does **not** break the existing
single-Linux-runner cross-compile model in `make dist`._

**Inferred (not yet proven):**
_hugot's text-generation pipeline is currently ORT-backend only (needs
`-tags ORT` + cgo + the ONNX Runtime shared library) — this is a real
fork: adding text generation later means either a real per-OS build
matrix for that build variant, or shipping it as a separate opt-in
binary/build tag rather than baking it into the default cross-compiled
`hawp` binary. Not resolved here; tracked as follow-up scope below._

**Scope — what else is affected:**
_`librarian/go/go.mod` (first external dependency — previously zero),
new `internal/application/embed/` package, `hawp embed`/`hawp model
pull` CLI commands, `make dist` (must keep working with `CGO_ENABLED=0`
for the default build), `librarian/go/README.md`._

---

### Recommended Fix

- Add `github.com/knights-analytics/hugot` as a dependency; use its
  pure-Go (`NewGoSession`) backend by default — no cgo, no ONNX Runtime
  binary dependency for this feature.
- `internal/application/embed`: wrap `hugot.DownloadModel` (model
  pulling) and the feature-extraction pipeline (embedding) behind a
  small API; default to the already-provisioned `Xenova/all-MiniLM-L6-v2`
  model so `hawp embed` works without an extra pull step.
- CLI: `hawp model pull <hf-repo>` (generic HF model pull, any pipeline
  type) and `hawp embed <text>` (embedding via the default or a pulled
  model).
- Defer text generation and other cgo/ORT-only pipelines to a follow-up
  item once the build-matrix question is decided — do not silently ship
  a feature that only works on whichever platform happens to build it.
- Verify with a real, non-mocked run: pull the real model over the
  network, embed real sentences, sanity-check that semantically similar
  sentences score higher cosine similarity than dissimilar ones.

**What to verify after:**

- [x] `CGO_ENABLED=0 go build`/`test` pass for the new package and cross-compile
      cleanly for all six release platforms with hugot actually linked in
      (Evidence: `CGO_ENABLED=0 go build`/`go test` of the real
      `internal/application/embed` package pass; `CGO_ENABLED=0 go build`
      for darwin/amd64, linux/arm64, windows/amd64 succeed with the
      package actually linked; `make dist` produces all six real
      binaries — `bin/dist/hawp-{darwin,linux,windows}-{amd64,arm64}`)
- [x] Real model pull + real embedding run produces plausible similarity
      ordering (similar sentences score higher than dissimilar ones)
      (Evidence: real run — "cat/feline" cosine 0.557 vs "cat/revenue"
      cosine 0.017, 384-dim vectors matching MiniLM's architecture)
- [x] `hawp embed` and `hawp model pull` work end to end on this machine
      (Evidence: real CLI run — `hawp embed "The quick brown fox" "A
      fast auburn fox"` printed two distinct 384-dim vectors from the
      compiled binary, not just the internal package)
- [x] `make dist` still produces all six binaries after adding the dependency
      (Evidence: `make dist VERSION=v0.0.5` succeeded for all six
      platforms; binary size grew from ~6.6 MB to ~16-18 MB per
      platform, disclosed honestly in the README rather than hidden)

---

## Outcome (filled at close)

Closed 2026-07-21. Adopted
[hugot](https://github.com/knights-analytics/hugot) (Apache 2.0) rather
than hand-rolling ONNX bindings/tokenization, per the explicit follow-up
to prefer an existing open-source project. `internal/application/embed`
wraps hugot's model downloader (`PullModel`/`PullDefaultModel`) and
feature-extraction pipeline (`Embed`) behind a small API; `hawp model
pull <hf-org/repo>` and `hawp embed <text>...` are the new CLI commands,
both using hugot's **pure-Go backend** (`NewGoSession`, GoMLX's
`simplego`) — confirmed empirically to need no cgo and no external ONNX
Runtime binary, and to cross-compile cleanly for all six `make dist`
platforms with the dependency actually linked in (not just an unlinked
no-op check). `hawp embed` defaults to the pinned
`Xenova/all-MiniLM-L6-v2` model already known from `e98de8c4`'s
provisioning work, pulling it on first use.

**Real, not mocked, verification**: pulled the actual model over the
network, ran real inference on three sentences, and confirmed the
semantic behavior is correct (similar sentences score far higher cosine
similarity than dissimilar ones) — both via a temporary internal probe
and via the real compiled `hawp embed` CLI command.

**Deliberately deferred, not silently shipped**: hugot's text-generation
pipeline requires its ORT backend (cgo + libonnxruntime), which would
break the current single-Linux-runner cross-compile model the pure-Go
path preserves. This is registered as `Available: false` in the command
registry (`hawp generate`, planned) rather than wired in half-working —
a real per-OS build matrix decision is needed first, tracked as this
item's own follow-up, not bundled into `fbf12a93`.

**Disclosed trade-off**: binary size grew roughly 2.5–3x (≈6.6 MB →
≈16–18 MB per platform) from hugot's GoMLX/onnx-gomlx dependency tree,
even though only the pure-Go feature-extraction path is used. Documented
plainly in `librarian/go/README.md` rather than left for someone to
discover later.

## Verification (filled at close)

- Evidence: `CGO_ENABLED=0 go build`/`go vet`/`go test` all pass for
  `internal/application/embed` and the full repo (`make check`).
- Evidence: real cross-compiles (darwin/amd64, linux/arm64,
  windows/amd64) of the actual linked package succeed under
  `CGO_ENABLED=0`; `make dist VERSION=v0.0.5` produces all six real
  release binaries.
- Evidence: real network run — model pull of `Xenova/all-MiniLM-L6-v2`,
  real inference on 3 sentences, correct similarity ordering (0.557 vs
  0.017 cosine).
- Evidence: real compiled `hawp embed` CLI run on this machine printed
  two correct 384-dim vectors.
- Evidence: `TestRegistryStaysInSyncWithHelpText`,
  `TestRunModelPullRequiresRepoArg`, `TestRunEmbedRequiresTextArg` all
  pass; full validation sweep (Go tests, TS validator, both link
  checkers, `hawp check`) green.

## Close Checklist

- [x] Outcome section filled
- [x] Verification section filled
- [x] Plan file saved under `closed/2026/07/21/748609a8-ac1a-4e27-a0dc-7a67d3b9fd14.md`
- [x] BACKLOG.md updated
