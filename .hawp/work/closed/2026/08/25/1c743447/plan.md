# 1c743447 — Configurable --hybrid-ratio flag

**Type:** feature  
**Status:** done  
**Closed:** 2026-08-25

## Outcome

`HybridRank` signature extended with `lexicalWeight float32` (0 → default 0.3/0.7). `--hybrid-ratio <f>` CLI flag added with [0.0, 1.0] validation; out-of-range exits 1 with clear error. All call sites updated (mcp/tool_search.go, benchmark.go, Query() pass 0 for default). Merged cleanly with `--verbose` flag from token-reduction branch (manual conflict resolution in run.go). Smoke tested: 0.5 and 0.7 return different orderings; 1.5 rejected with exit 1.

## Verification

- [x] `hawp search "backlog alignment" --hybrid-ratio 0.5` returns results without error
- [x] `hawp search "backlog alignment" --hybrid-ratio 0.7` returns different ordering
- [x] `hawp search "test" --hybrid-ratio 1.5` exits 1 with clear error message
- [x] All Go tests pass
- [x] Build clean; binary at `.hawp/bin/hawp`
- [x] Benchmark results in `librarian/docs/benchmarks-v008.md`

## Close Checklist

- [x] Implementation complete
- [x] Smoke tests verified
- [x] Benchmark ratio sensitivity section in benchmarks-v008.md
- [x] BACKLOG updated; plan moved to closed

## What was done

- `HybridRank` signature: added `lexicalWeight float32`; 0 → default 0.3/0.7
- `--hybrid-ratio <f>` CLI flag: parses + validates [0.0, 1.0]; out-of-range exits 1
- All call sites (mcp/tool_search.go, benchmark.go, Query()) pass 0 for default behavior
- Merged with `--verbose` flag from token-reduction branch (manual conflict resolution)
- Build clean; binary updated; smoke-tested 0.5 and 0.7 produce different orderings
