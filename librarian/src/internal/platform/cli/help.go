package cli

func helpText() string {
	return `hawp

Go librarian CLI for HAWP — small native workflow intelligence tool.

USAGE
  hawp <command> [options]

COMMANDS
  uuid [--short]                       generate a work item UUID
  links check                          validate local markdown links (.hawp, docs, README.md)
  links clean [--apply]                relink (or, failing that, neutralize) broken links found by links check
  kit validate [--kit-path <path>]     validate .hawp/kit/ structure
  kit normalize [--apply]              normalize .hawp/kit/ names and links (dry-run default)
  work validate [--work-root <path>]   backlog/plan/evidence integrity checks
  work new "<title>" [--type ...]      scaffold intake: UUID, plan file, inbox backlog row
  work status --title "<label>" [--work-item <id>]    create status/YYYY/MM/DD/{id}/status.md
  work evidence --title "<label>" [--work-item <id>]  create evidence/YYYY/MM/DD/{id}/evidence.md
  work decision --title "<label>" [--work-item <id>]  create decisions/YYYY/MM/DD/{id}/decision.md
  work note --title "<label>" [--work-item <id>]      create notes/YYYY/MM/DD/{id}/note.md
  work normalize [--apply --migrate-folders --validate]  normalize work record drift (dry-run default)
  providers materialize                materialize shared provider behaviors into provider packs
  providers validate                   validate generated provider-pack files
  providers sync                       materialize + validate provider-pack files
  distribution build                   build generated install/update guides
  distribution validate                validate generated install/update guides
  distribution sync                    providers sync + build + validate generated guides
  check                                combined kit + work + links validation
  init [--provider <name>|all]           provision ~/.hawp, sync kit, write MCP configs (claude|cursor|codex|continue|github|all)
  mcp [--repo-root <path>]               start stdio MCP server for the selected repository
  mcp configure --provider <name>       configure MCP only; no downloads or kit sync
  version                               print the running hawp version
  update                                update binary + kit + all providers (--no-providers for kit-only)
  update --check                        check whether an update is available without installing (exit 1 = update ready)
  update --disable-auto                 disable the 21-min countdown auto-install (notices still print)
  update --enable-auto                  re-enable the auto-install countdown
  update latest                         update binary only
  update sync [--provider <name>|all]   sync kit (+ providers if specified)
  update verify                         check whether an update is available (exit 1 = update ready)
  commands [--json]                    list every command; --json is the agent-facing discovery output
  backlog validate                     alias for check
  backlog upgrade                      alias for work normalize
  db init                              plan the ~/.hawp home layout (scaffold)
  index build [--scope all|work|kit] [--export <path>]  enrich kit/work docs with folder context
  search index                                     ingest configured paths into SQLite (reads .hawp/config/search.json)
  search embed --backend onnx|ollama [--model <name>]  embed all chunks with vectors
  search <query> [--limit <n>] [--semantic] [--context] [--format markdown|json] [--max-tokens <n>] [--hybrid-ratio <f>]
                                                   lexical + vector hybrid search; --semantic for pure-vector mode;
                                                   --context for LLM-ready context block;
                                                   --hybrid-ratio tunes the lexical/semantic blend (default 0.3)
  model pull <hf-org/repo> [--onnx-file <path>]   download any Hugging Face ONNX model into ~/.hawp/models
  embed <text>... [--model <hf-org/repo>]         embed text via a local model (default: all-MiniLM-L6-v2)
  usage                                           show MCP call log totals (opt-in; run hawp usage enable first)
  usage log                                       tail 20 most recent logged calls
  usage report [--export <path>]                  full Markdown report: totals, per-tool breakdown, recent queries
  usage enable [--log-bodies]                     enable call logging (--log-bodies also stores raw input/output)
  usage disable                                   stop recording new calls
  usage clear                                     delete all stored log entries (irreversible)

WORK NORMALIZE OPTIONS
  --dry-run | --apply    detection only (default) | normalize closed records
  --validate             run workflow validation summary afterwards
  --format text|json     report format (dry-run)
  --output <path>        write report to file
  --export-plan <path>   write plan JSON (dry-run)
  --export-research-queue <path>  write research queue JSON
  --force-dirty          skip the apply-mode dirty-tree guard

SEARCH --CONTEXT OPTIONS
  --context              output LLM-ready context block (deduped + formatted)
  --semantic             pure-vector search (requires embed step; skips FTS5)
  --format markdown|json output format (default markdown)
  --max-tokens <n>       token budget for context block (default 2000)
  --verbose | -v         print token accounting summary to stderr (chunks, ~tokens, saved via dedup)
  --hybrid-ratio <f>     lexical fraction for hybrid blend [0.0, 1.0] (default 0.3)`
}
