# Token-Savings Benchmark

Shape budget: **2000 tokens** | Token estimate: `(len(text)+3)/4`

| # | Query (intent) | Results | Raw tokens | Shaped tokens | Saved | % saved |
|---|----------------|---------|------------|---------------|-------|---------|
| 1 | Work tracking policy | 10 | 3064 | 1991 | 1073 | 35% |
| 2 | Context transfer between sessions | 10 | 3031 | 1992 | 1039 | 34% |
| 3 | Evidence standards for findings | 10 | 2372 | 1989 | 383 | 16% |
| 4 | Intake process and investigation ordering | 5 | 1707 | 1707 | 0 | 0% |
| 5 | Provider distribution and materialization | 5 | 1796 | 1797 | -1 | -0% |
| 6 | MCP server configuration for AI agents | 8 | 2043 | 1992 | 51 | 2% |
| 7 | Core HAWP protocol shape fields | 2 | 591 | 591 | 0 | 0% |
| 8 | Plan file structure and fields | 10 | 3292 | 1990 | 1302 | 40% |
| 9 | Kit maintenance and validation commands | 10 | 2581 | 1991 | 590 | 23% |
| 10 | Binary update and install flow | 10 | 3044 | 1989 | 1055 | 35% |
| — | **TOTAL** | — | **23521** | **18029** | **5492** | **23%** |

_Context shaping applies deduplication + token-budget truncation._
_Raw tokens = sum of `len(chunk text)/4` across all ranked results._
_Shaped tokens = `ContextBlock.TokenCount` after `FormatAsMarkdown`._
_Negative savings = sparse result set already under budget; shaper adds Markdown formatting overhead._
