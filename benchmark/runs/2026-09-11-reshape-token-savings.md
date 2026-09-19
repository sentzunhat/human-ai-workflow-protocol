# Reshape Token-Savings Benchmark

Backend: **onnx / SmolLM2-360M-Instruct** | Token estimate: `(len(text)+3)/4` | Date: 2026-09-11

## Result: Live run failed — 0/10 queries produced JSON output

**Root cause (diagnosed):** SmolLM2-360M-Instruct does not follow a 770-char
intake-extraction instruction embedded in the ChatML user turn when the system
turn contains the separate `ReshapingPrompt`. With the full intake prompt in
the user turn, the model ignores it and replies with a generic assistant opener
("I'm ready to assist you…"). With a shorter/simpler prompt (direct question,
no nested role instruction), the model does generate JSON-like output but using
the wrong schema (JSON Schema format, not the `mission/constraints/output`
fields we request).

**Pathway:** `ONNXIntakeShaper.Shape` → `buildONNXIntakePrompt` (770 chars) →
`ONNXLLMClient.Reshape` → `reshapeOne` wraps in ChatML with `ReshapingPrompt`
in system turn → model ignores user-turn instruction at this size.

**Implication:** The ONNX reshape path requires either (a) a larger/stronger
model, or (b) restructuring the ChatML call so the intake instruction is the
system prompt and the user request is the only user turn content. The current
`ONNXLLMClient.Reshape` API puts a fixed reshaping system prompt and treats
all user content as the user turn — that design works for context reshaping but
not for structured intake extraction.

**Gate status: NOT MET** for ONNX reshape. Blocked on model capability or API
restructuring. Slice B remains open. See plan `288d543c`.

---

| # | Request (truncated) | Raw tokens | Shaped tokens | Saved | % saved | Note |
|---|---------------------|------------|---------------|-------|---------|------|
| 1 | add a --dry-run flag to work normalize that sho... | 23 | — | — | — | error: no JSON object found in LLM response |
| 2 | the hawp search benchmark is slow, optimize it ... | 15 | — | — | — | error: no JSON object found in LLM response |
| 3 | document the mcp server tools with examples for... | 24 | — | — | — | error: no JSON object found in LLM response |
| 4 | migrate the TypeScript scripts to Go and deprec... | 21 | — | — | — | error: no JSON object found in LLM response |
| 5 | add support for selecting multiple providers in... | 19 | — | — | — | error: no JSON object found in LLM response |
| 6 | the kit validate command should check for broke... | 23 | — | — | — | error: no JSON object found in LLM response |
| 7 | implement context deduplication to avoid sendin... | 39 | — | — | — | error: no JSON object found in LLM response |
| 8 | improve error messages when index.sqlite does n... | 18 | — | — | — | error: no JSON object found in LLM response |
| 9 | add a hawp providers list command that shows in... | 23 | — | — | — | error: no JSON object found in LLM response |
| 10 | the work normalize folder migration reports fal... | 42 | — | — | — | error: no JSON object found in LLM response |
| — | **TOTAL (0/10 succeeded)** | **247** | **0** | — | — | |

_Raw tokens = `(len(request)+len(context)+3)/4` on the verbatim user input._
_Shaped tokens = `(len(mission)+len(constraints)+len(output)+len(checkpoint)+3)/4` on the DraftIntake output._
_Note: negative savings (expansion) is expected for short requests even on success — structured intake adds labeled fields._
_The value of reshaping is precision and downstream filtering, not raw token compression._
_v0.1.0 gate: avg shaped tokens < avg raw tokens (any net savings across the 10-query suite). Gate NOT MET._
