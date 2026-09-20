# Reshape Token-Savings Benchmark

Backend: **onnx / microsoft/Phi-3-mini-4k-instruct-ort-genai-int4-cpu** | Token estimate: `(len(text)+3)/4` | Date: 2026-09-11

| # | Request (truncated) | Raw tokens | Shaped tokens | Saved | % saved | Note |
|---|---------------------|------------|---------------|-------|---------|------|
| 1 | add a --dry-run flag to work normalize that sho... | 23 | 96 | -73 | -317% | |
| 2 | the hawp search benchmark is slow, optimize it ... | 15 | 52 | -37 | -247% | |
| 3 | document the mcp server tools with examples for... | 24 | 81 | -57 | -238% | |
| 4 | migrate the TypeScript scripts to Go and deprec... | 21 | 65 | -44 | -210% | |
| 5 | add support for selecting multiple providers in... | 19 | 89 | -70 | -368% | |
| 6 | the kit validate command should check for broke... | 23 | 51 | -28 | -122% | |
| 7 | implement context deduplication to avoid sendin... | 39 | 58 | -19 | -49% | |
| 8 | improve error messages when index.sqlite does n... | 18 | 66 | -48 | -267% | |
| 9 | add a hawp providers list command that shows in... | 23 | 72 | -49 | -213% | |
| 10 | the work normalize folder migration reports fal... | 42 | 80 | -38 | -90% | |
| — | **TOTAL (10/10 succeeded)** | **247** | **710** | **-463** | **-187%** | |

_Raw tokens = `(len(request)+len(context)+3)/4` on the verbatim user input._
_Shaped tokens = `(len(mission)+len(constraints)+len(output)+len(checkpoint)+3)/4` on the DraftIntake output._
_Negative savings (expansion) is expected for short requests: structured intake adds labeled fields._
_The value of reshaping is precision and downstream filtering, not raw token compression._
_v0.1.0 reshape coverage gate: 10/10 queries return valid structured output.
**Met (10/10).** The short-request token-compression gate is **not met** here
(247 raw tokens versus 710 shaped tokens); downstream savings is measured
separately in `benchmark/runs/2026-09-12-downstream-token-savings-ollama.md`._
