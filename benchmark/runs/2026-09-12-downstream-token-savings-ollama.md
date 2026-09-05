# Downstream Savings Benchmark

Backend: **ollama / mistral** | Token estimate: `(len(text)+3)/4` | Date: 2026-09-12

| # | Query (intent) | Raw (req+ctx) tokens | Shaped tokens | Saved | % saved |
|---|----------------|----------------------|---------------|-------|---------|
| 1 | Work tracking policy | 1996 | 115 | +1881 | +94% |
| 2 | Context transfer between sessions | 1997 | 86 | +1911 | +96% |
| 3 | Evidence standards for findings | 1996 | 211 | +1785 | +89% |
| 4 | Intake process and investigation ordering | 1716 | 30 | +1686 | +98% |
| 5 | Provider distribution and materialization | 1802 | 84 | +1718 | +95% |
| 6 | MCP server configuration for AI agents | 1998 | — | — | — |
| 7 | Core HAWP protocol shape fields | 603 | 61 | +542 | +90% |
| 8 | Plan file structure and fields | 1996 | 73 | +1923 | +96% |
| 9 | Kit maintenance and validation commands | 1997 | 84 | +1913 | +96% |
| 10 | Binary update and install flow | 1995 | 87 | +1908 | +96% |
| — | **TOTAL (9/10 succeeded)** | **18096** | **831** | **+17265** | **+95%** |

_Raw tokens = `(len(query) + len(formatted_search_context) + 3) / 4`._
_This models the downstream LLM call without reshape: request text + retrieved context._
_Shaped tokens = `(len(mission)+len(constraints)+len(output)+len(checkpoint)+3)/4` on the DraftIntake output._
_Positive savings = shaped draft is more compact than request + context (downstream token reduction)._
_v0.1.0 gate: ≥20% average savings across succeeded queries._
