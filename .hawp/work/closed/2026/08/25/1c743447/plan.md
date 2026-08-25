# 1c743447 — Configurable --hybrid-ratio flag

**Type:** feature  
**Status:** done  
**Closed:** 2026-08-25

## What was done

- `HybridRank` signature: added `lexicalWeight float32`; 0 → default 0.3/0.7
- `--hybrid-ratio <f>` CLI flag: parses + validates [0.0, 1.0]; out-of-range exits 1
- All call sites (mcp/tool_search.go, benchmark.go, Query()) pass 0 for default behavior
- Merged with `--verbose` flag from token-reduction branch (manual conflict resolution)
- Build clean; binary updated; smoke-tested 0.5 and 0.7 produce different orderings
