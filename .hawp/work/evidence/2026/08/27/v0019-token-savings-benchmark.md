# Token-Savings Benchmark

Shape budget: **2000 tokens** | Token estimate: `(len(text)+3)/4`

| # | Query (intent) | Results | Raw tokens | Shaped tokens | Saved | % saved |
|---|----------------|---------|------------|---------------|-------|---------|
| 1 | Work tracking policy | 10 | 2707 | 1994 | 713 | 26% |
| 2 | Context transfer between sessions | 10 | 2856 | 1987 | 869 | 30% |
| 3 | Evidence standards for findings | 10 | 2642 | 1994 | 648 | 25% |
| 4 | Intake process and investigation ordering | 5 | 1707 | 1896 | -189 | -11% |
| 5 | Provider distribution and materialization | 5 | 1796 | 1981 | -185 | -10% |
| 6 | MCP server configuration for AI agents | 5 | 1266 | 1453 | -187 | -15% |
| 7 | Core HAWP protocol shape fields | 2 | 591 | 707 | -116 | -20% |
| 8 | Plan file structure and fields | 10 | 3220 | 1987 | 1233 | 38% |
| 9 | Kit maintenance and validation commands | 10 | 2457 | 1995 | 462 | 19% |
| 10 | Binary update and install flow | 10 | 2847 | 1992 | 855 | 30% |
| — | **TOTAL** | — | **22089** | **17986** | **4103** | **19%** |

_Context shaping applies deduplication + token-budget truncation._
_Raw tokens = sum of `len(chunk text)/4` across all ranked results._
_Shaped tokens = `ContextBlock.TokenCount` after `FormatAsMarkdown`._
