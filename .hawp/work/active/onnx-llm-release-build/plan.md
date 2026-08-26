# onnx-llm-release-build — Ship ONNX LLM in the release binary

**Type:** feature
**Status:** plan-ready
**Branch:** `feature/v0.0.13`
**Opened:** 2026-08-25

## Problem

The ONNX Text2Text reshaper (SmolLM2-360M-Instruct via hugot ORT backend) was
verified working end-to-end on 2026-07-27 (macOS arm64, ~1.1s/reshape, ChatML
prompt format). It has never shipped in a release binary because:

1. The release build uses `CGO_ENABLED=0` with no `-tags ORT`.
2. The ORT backend needs three native libraries (`libonnxruntime`,
   `libonnxruntime-genai`, `libtokenizers`) that must be bundled per platform.
3. Intel Mac (macOS/amd64) has no official Microsoft prebuilt for either
   onnxruntime or onnxruntime-genai — needs special handling.

Everything from Go source through hugot is already working. This is a
CI/packaging task.

## What "working" means (already verified)

- `hawp search "query" --context` with `llm.backend: "onnx"` and
  `llm.model: "homen3/SmolLM2-360M-Instruct-ort-genai-int4-cpu"` calls
  `ContextReshaper.Reshape()` → `ONNXLLMClient.Reshape()` → hugot ORT session
  → decoded ChatML output in ~1.1s on arm64.
- `ReshapeBatch()` loops single-prompt calls (batch EOS bug fixed 2026-07-27).
- Prompt format: `<|im_start|>system\n...\n<|im_end|>\n<|im_start|>user\n...<|im_end|>\n<|im_start|>assistant\n`

## Required platforms

| Platform | onnxruntime | onnxruntime-genai | libtokenizers | Status |
|----------|-------------|-------------------|---------------|--------|
| darwin/arm64 | ✅ official | ✅ official | ✅ static | Ready |
| linux/arm64 | ✅ official | ✅ official | ✅ static | Ready |
| linux/amd64 | ✅ official | ✅ official | ✅ static | Ready |
| windows/amd64 | ✅ official | ✅ official | ✅ static | Ready |
| windows/arm64 | ✅ official | ✅ official | ✅ static | Ready |
| darwin/amd64 | ❌ no official build | ❌ no official build | ✅ static | Skip or from-source |

Intel Mac: ship `hawp` without `-tags ORT` on darwin/amd64 (ONNX LLM
unavailable; Ollama LLM still works). Print a clear error when `llm.backend:
"onnx"` is selected on an unsupported platform.

## Implementation tasks

### 1. CI: update `release-librarian-go.yml`

For each ORT-supported platform, the build step becomes:

```bash
# Download native libs into a temp dir
# e.g. for darwin/arm64:
curl -L https://github.com/microsoft/onnxruntime/releases/download/v1.x.y/\
  onnxruntime-osx-arm64-1.x.y.tgz | tar xz -C $NATIVE_DIR
curl -L https://github.com/microsoft/onnxruntime-genai/releases/download/v0.x.y/\
  onnxruntime-genai-osx-arm64-0.x.y.tgz | tar xz -C $NATIVE_DIR
# libtokenizers.darwin-arm64.tar.gz from daulet/tokenizers releases

CGO_ENABLED=1 \
CGO_CFLAGS="-I$NATIVE_DIR/include" \
CGO_LDFLAGS="-L$NATIVE_DIR/lib -Wl,-rpath,@executable_path/../lib \
             -lonnxruntime -lonnxruntime-genai -ltokenizers -ldl -lm -lstdc++" \
go build -tags ORT -trimpath -ldflags="$LDFLAGS" -o hawp ./cmd/hawp
```

The release artifact for each ORT platform becomes a tarball containing:
- `hawp` binary
- `lib/libonnxruntime.dylib` (or `.so`, `.dll`)
- `lib/libonnxruntime-genai.dylib`
- `README.txt` — install instructions (copy lib/ next to binary or to
  `/usr/local/lib`)

Darwin/amd64 keeps the current CGO_ENABLED=0 single-binary build, no ORT.

### 2. `hawp mcp` / `hawp search` — unsupported platform error

When `llm.backend: "onnx"` is requested and the binary was built without
`-tags ORT`, return a clear error:
```
hawp: ONNX LLM is not available on this platform (darwin/amd64).
Use llm.backend: "ollama" instead, or build from source with -tags ORT.
```
This error already fires via hugot's own guard; just make sure it surfaces
cleanly to the user rather than as a Go panic.

### 3. Prompt quality pass

Run 10 real reshaping tasks (mix of kit docs and work items) through
SmolLM2-360M-Instruct and score output for:
- Completeness (key facts retained)
- Clarity (readable by a human, not just an LLM)
- Length (within --max-tokens budget)

Document results in `.hawp/work/evidence/2026/08/25/onnx-llm-quality/`.

### 4. Version + CHANGELOG

Bump to `0.0.13`, add CHANGELOG entry documenting the ORT build flag,
supported platforms, and model.

## Acceptance criteria

- [ ] `hawp search "query" --context` with `llm.backend: "onnx"` works on
  darwin/arm64, linux/arm64, linux/amd64, windows/amd64 release binaries
- [ ] darwin/amd64 binary prints clear unsupported error (no panic)
- [ ] Release tarballs include the native lib directory alongside the binary
- [ ] `go test ./...` passes with `-tags ORT` on arm64
- [ ] Quality pass: 8/10 reshaping tasks score ≥ 7/10 on completeness + clarity
- [ ] CHANGELOG updated

## Dependencies

- v0.0.12 (usage-log) should land first — quality pass benefits from token
  count visibility to verify the reshape actually reduces context size.
- onnxruntime ≥ 1.19, onnxruntime-genai ≥ 0.4, libtokenizers from daulet/tokenizers

## Relationship to v0.1.0

Shipping the ONNX LLM in the release binary completes the last Tier 1 Core
item in the v0.1.0 vision doc. After this + usage-log token evidence, v0.1.0
is ready to tag.
