# Reshape Token-Savings Benchmark

Backend: **onnx / /Users/beltrd/.hawp/models/llm/homen3_SmolLM2-360M-Instruct-ort-genai-int4-cpu** | Token estimate: `(len(text)+3)/4` | Date: 2026-09-11

| # | Request (truncated) | Raw tokens | Shaped tokens | Saved | % saved | Note |
|---|---------------------|------------|---------------|-------|---------|------|
| 1 | add a --dry-run flag to work normalize that sho... | 23 | 37 | -14 | -61% | |
| 2 | the hawp search benchmark is slow, optimize it ... | 15 | — | — | — | error: shape intake: no JSON object found in LLM response |
| 3 | document the mcp server tools with examples for... | 24 | — | — | — | error: shape intake: no JSON object found in LLM response |
| 4 | migrate the TypeScript scripts to Go and deprec... | 21 | — | — | — | error: shape intake: no JSON object found in LLM response |
| 5 | add support for selecting multiple providers in... | 19 | — | — | — | error: shape intake: no JSON object found in LLM response |
| 6 | the kit validate command should check for broke... | 23 | — | — | — | error: shape intake: no JSON object found in LLM response |
| 7 | implement context deduplication to avoid sendin... | 39 | — | — | — | error: shape intake: no JSON object found in LLM response |
| 8 | improve error messages when index.sqlite does n... | 18 | — | — | — | error: shape intake: no JSON object found in LLM response |
| 9 | add a hawp providers list command that shows in... | 23 | — | — | — | error: shape intake: no JSON object found in LLM response |
| 10 | the work normalize folder migration reports fal... | 42 | — | — | — | error: shape intake: no JSON object found in LLM response |
| — | **TOTAL (1/10 succeeded)** | **247** | **37** | **+210** | **+85%** | |

_Raw tokens = `(len(request)+len(context)+3)/4` on the verbatim user input._
_Shaped tokens = `(len(mission)+len(constraints)+len(output)+len(checkpoint)+3)/4` on the DraftIntake output._
_Negative savings (expansion) is expected for short requests: structured intake adds labeled fields._
_The value of reshaping is precision and downstream filtering, not raw token compression._
_v0.1.0 gate: avg shaped tokens < avg raw tokens (any net savings across the 10-query suite)._
