# Librarian CLI / Binary / Transformers.js WebGPU Research

Date: 2026-07-03

## Question

Should `librarian/` grow into a first-class installable CLI and potentially a
single-file binary, and should `transformers.js` + WebGPU be part of that plan?

## Directly verified

### Local repo state

- `.hawp/work/BACKLOG.md` has no pre-existing active work; the highest-leverage
  open repo-native items were `Verification clarity cleanup` and the optional
  MCP overlay note.
- `README.md` currently describes `librarian/` as a repo-local maintenance tree,
  not a separately installed CLI product.
- `.hawp/bin/hawp` is a repo-local wrapper over `librarian/package.json` scripts,
  not a published package/binary surface.

### Current external-source picture

- **Transformers.js docs** describe a Node tutorial and both ESM/CommonJS usage,
  which confirms Node-side execution is a supported lane for the library.
  Source: <https://huggingface.co/docs/transformers.js/tutorials/node>
- **Transformers.js docs** also show `device: 'webgpu'` as a first-class option
  and document a WebGPU guide.
  Source: <https://huggingface.co/docs/transformers.js/en/index>
- **Transformers.js v4 blog** says the new WebGPU runtime can run directly in
  Node, Bun, and Deno.
  Source: <https://huggingface.co/blog/transformersjs-v4>
- **Node SEA docs** say single executable applications are `Stability: 1.1 -
  Active development`, support `--build-sea`, and support a single embedded
  script using CommonJS or ESM.
  Source: <https://nodejs.org/api/single-executable-applications.html>
- **ONNX Runtime Node binding docs** list prebuilt Node EP support for CPU,
  DirectML (Windows), and CUDA (Linux), but not WebGPU.
  Source: <https://onnxruntime.ai/docs/get-started/with-javascript/node.html>
- **ONNX Runtime Web docs** explicitly show Node support only for single-threaded
  WASM in the Web package compatibility table, with WebGPU marked unsupported for
  Node there.
  Source: <https://onnxruntime.ai/docs/get-started/with-javascript/web.html>

## Inference

- The **CLI productization lane is viable now**: installable npm CLI + package
  surface + command contract cleanup are all straightforward repo-owned work.
- The **single-binary lane is plausible but not yet default-safe**: SEA is real,
  but still active-development and will need explicit handling for assets,
  native modules, and cross-platform release flow.
- The **`transformers.js` lane should be treated as a bounded spike, not the CLI
  foundation**:
  - Hugging Face’s own current story is optimistic for Node WebGPU.
  - Current ONNX Runtime official docs still present a more conservative Node EP
    matrix.
  - Therefore the repo should not assume "WebGPU in Node is solved everywhere"
    until a local proof matrix exists for this exact toolchain.

## Recommended staging

1. **CLI first**: make `librarian` a clean npm-installable CLI with stable
   command contracts and explicit subcommand ownership.
2. **Runtime spike second**: prove whether a real `transformers.js` Node command
   belongs inside librarian, and whether the supported baseline is WASM-only,
   optional WebGPU, or a different acceleration path.
3. **Binary packaging third**: only after the CLI surface and AI runtime
   dependency graph are explicit, because packaging native/large-model assets too
   early will create false certainty.

## Proposed work items

- `1c2db85b` — research librarian full CLI, binary packaging, and
  `transformers.js` WebGPU feasibility
- `8dedc4e2` — design installable librarian CLI surface and bin contract
- `1a3b32a4` — spike `transformers.js` Node and WebGPU viability for librarian
  subcommands
- `54e68af7` — spike single-binary packaging and model asset strategy for
  librarian
