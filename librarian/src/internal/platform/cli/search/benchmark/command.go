package benchmark

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	appcontext "github.com/sentzunhat/hawp/librarian/src/internal/application/context"
	appsearch "github.com/sentzunhat/hawp/librarian/src/internal/application/search"
	appintake "github.com/sentzunhat/hawp/librarian/src/internal/application/work/intake"
	domainsearch "github.com/sentzunhat/hawp/librarian/src/internal/domain/search"
	"github.com/sentzunhat/hawp/librarian/src/internal/infrastructure/filesystem"
	inframodels "github.com/sentzunhat/hawp/librarian/src/internal/infrastructure/models"
	"github.com/sentzunhat/hawp/librarian/src/internal/infrastructure/repo"
	sqlite "github.com/sentzunhat/hawp/librarian/src/internal/infrastructure/repositories/index"
)

func Run(args []string, cwd string) error {
	opts, err := parseArgs(args)
	if err != nil {
		return err
	}

	root, err := repo.FindBacklogRepoRoot(cwd)
	if err != nil {
		fmt.Println("Not in a HAWP repo; no index to benchmark.")
		return nil
	}

	// Reshape-only benchmark does not require an indexed database.
	if opts.reshapeToken {
		return runReshapeBenchmark(opts.reshapeBackend, opts.reshapeURL, opts.reshapeModel, opts.exportPath)
	}

	dbPath, err := filesystem.ResolveSafeSearchIndexPath(root)
	if err != nil {
		return err
	}
	db, err := sqlite.Open(dbPath)
	if err != nil {
		fmt.Printf("Index not found at %s. Run `hawp search index` first.\n", dbPath)
		return nil
	}
	defer db.Close()

	if opts.downstreamToken {
		return runDownstreamBenchmark(opts.downstreamBackend, opts.downstreamURL, opts.downstreamModel, db, opts.exportPath)
	}
	if opts.tokenMode {
		return runTokenBenchmark(db, opts.exportPath)
	}
	return runBenchmark(db)
}

// ─── Search-pattern benchmark ────────────────────────────────────────────────

type benchmarkQuery struct {
	Query            string
	Intent           string
	RelevantKeywords []string
}

type benchmarkResult struct {
	Query            string
	Pattern          string
	LatencyMS        float64
	ResultCount      int
	TopResultQuality string
}

var benchmarkQueries = []benchmarkQuery{
	{
		Query:            "backlog alignment rules",
		Intent:           "Work tracking policy",
		RelevantKeywords: []string{"backlog", "alignment", "active", "closed"},
	},
	{
		Query:            "status report handoff",
		Intent:           "Context transfer between sessions",
		RelevantKeywords: []string{"status", "report", "handoff", "session"},
	},
	{
		Query:            "evidence discipline patterns",
		Intent:           "Evidence standards for findings",
		RelevantKeywords: []string{"evidence", "discipline", "finding", "inference"},
	},
	{
		Query:            "intake workflow investigation first",
		Intent:           "Intake process and investigation ordering",
		RelevantKeywords: []string{"intake", "investigation", "plan", "workflow"},
	},
	{
		Query:            "provider overlay sync",
		Intent:           "Provider distribution and materialization",
		RelevantKeywords: []string{"provider", "overlay", "sync", "distribution"},
	},
	{
		Query:            "hawp mcp stdio server tools",
		Intent:           "MCP server configuration for AI agents",
		RelevantKeywords: []string{"mcp", "server", "tool", "stdio"},
	},
	{
		Query:            "HAWP shape template mission constraints output",
		Intent:           "Core HAWP protocol shape fields",
		RelevantKeywords: []string{"mission", "constraints", "output", "shape"},
	},
	{
		Query:            "work item plan file format",
		Intent:           "Plan file structure and fields",
		RelevantKeywords: []string{"plan", "work", "item", "status"},
	},
	{
		Query:            "kit normalize validate",
		Intent:           "Kit maintenance and validation commands",
		RelevantKeywords: []string{"kit", "normalize", "validate", "naming"},
	},
	{
		Query:            "hawp update binary install",
		Intent:           "Binary update and install flow",
		RelevantKeywords: []string{"update", "binary", "install", "release"},
	},
}

func runBenchmark(db *sqlite.IndexDB) error {
	fmt.Println("\n╔════════════════════════════════════════════════════════════════╗")
	fmt.Println("║        SEARCH BENCHMARK: LEXICAL / SEMANTIC / HYBRID          ║")
	fmt.Println("╚════════════════════════════════════════════════════════════════╝")

	hasVectors, _ := db.HasVectors()
	availablePatterns := []string{"lexical"}
	if hasVectors {
		availablePatterns = append(availablePatterns, "semantic", "hybrid")
	}

	fmt.Printf("Available patterns: %s\n", strings.Join(availablePatterns, ", "))
	if !hasVectors {
		fmt.Println("Note: Vectors not found. Run 'hawp search embed' first to enable semantic and hybrid search.")
	}
	fmt.Printf("Running %d queries across %d patterns...\n\n", len(benchmarkQueries), len(availablePatterns))

	var results []benchmarkResult
	for _, query := range benchmarkQueries {
		for _, pattern := range availablePatterns {
			result := benchmarkOneQuery(query, pattern, db)
			results = append(results, result)
			fmt.Printf("✓ %s (%s): %dms, %d results [%s]\n",
				query.Intent, pattern, int(result.LatencyMS), result.ResultCount, result.TopResultQuality)
		}
	}

	printBenchmarkSummary(results)
	return nil
}

func benchmarkOneQuery(query benchmarkQuery, pattern string, db *sqlite.IndexDB) benchmarkResult {
	result := benchmarkResult{Query: query.Query, Pattern: pattern}
	start := time.Now()

	switch pattern {
	case "lexical":
		rows, _ := db.QueryChunksLexical(query.Query, 10)
		result.ResultCount = len(rows)
		result.TopResultQuality = assessQuality(query.RelevantKeywords, rows)

	case "semantic":
		rows := appsearch.SemanticSearchWithEmbedder(query.Query, db, 10, inframodels.NewEmbedder)
		result.ResultCount = len(rows)
		result.TopResultQuality = assessQuality(query.RelevantKeywords, rows)

	case "hybrid":
		rows, _ := db.QueryChunksLexical(query.Query, 30)
		if len(rows) > 0 {
			rows = appsearch.HybridRankWithEmbedder(rows, query.Query, db, 10, 0, inframodels.NewEmbedder)
		}
		result.ResultCount = len(rows)
		result.TopResultQuality = assessQuality(query.RelevantKeywords, rows)
	}

	result.LatencyMS = float64(time.Since(start).Milliseconds())
	return result
}

func assessQuality(keywords []string, results []map[string]interface{}) string {
	if len(results) == 0 {
		return "low"
	}
	top := strings.ToLower(getStr(results[0], "text"))
	for _, kw := range keywords {
		if strings.Contains(top, strings.ToLower(kw)) {
			return "high"
		}
	}
	return "low"
}

func printBenchmarkSummary(results []benchmarkResult) {
	fmt.Println("\n╔════════════════════════════════════════════════════════════════╗")
	fmt.Println("║                      RESULTS SUMMARY                          ║")
	fmt.Println("╚════════════════════════════════════════════════════════════════╝")

	patternStats := make(map[string][]float64)
	patternQuality := make(map[string]map[string]int)

	for _, r := range results {
		patternStats[r.Pattern] = append(patternStats[r.Pattern], r.LatencyMS)
		if patternQuality[r.Pattern] == nil {
			patternQuality[r.Pattern] = make(map[string]int)
		}
		patternQuality[r.Pattern][r.TopResultQuality]++
	}

	fmt.Println("Pattern      | Avg Latency | Min/Max    | Queries | Quality (High/Medium/Low)")
	fmt.Println("-------------|-------------|-----------|---------|-------------------------")

	patterns := []string{"lexical", "semantic", "hybrid"}
	for _, pattern := range patterns {
		latencies, ok := patternStats[pattern]
		if !ok {
			continue
		}

		avgLatency := averageFloat(latencies)
		minLatency := minFloat(latencies)
		maxLatency := maxFloat(latencies)

		quality := patternQuality[pattern]
		high := quality["high"]
		med := quality["medium"]
		low := quality["low"]

		fmt.Printf("%-12s | %8.1fms  | %5.1f/%5.1fms | %7d | %2d / %2d / %2d\n",
			pattern, avgLatency, minLatency, maxLatency,
			len(latencies), high, med, low)
	}

	fmt.Println("\n" + "════════════════════════════════════════════════════════════════")
	fmt.Println("RELATIVE PERFORMANCE (vs Lexical):")
	fmt.Println("════════════════════════════════════════════════════════════════")

	if lexicalLatencies, ok := patternStats["lexical"]; ok {
		lexicalAvg := averageFloat(lexicalLatencies)

		for _, pattern := range []string{"semantic", "hybrid"} {
			if latencies, ok := patternStats[pattern]; ok {
				patternAvg := averageFloat(latencies)
				if lexicalAvg < 1 {
					fmt.Printf("%s:  %.1fms avg (lexical sub-millisecond; relative multiple not meaningful)\n",
						pattern, patternAvg)
				} else {
					multiple := patternAvg / lexicalAvg
					fmt.Printf("%s:  %.1fx slower than lexical (%.1fms vs %.1fms)\n",
						pattern, multiple, patternAvg, lexicalAvg)
				}
			}
		}
	}

	fmt.Println("\n" + "════════════════════════════════════════════════════════════════")
	winner := ""
	winnerHigh := -1
	winnerLatency := 1e9
	for _, pattern := range patterns {
		if _, ok := patternStats[pattern]; !ok {
			continue
		}
		high := patternQuality[pattern]["high"]
		avg := averageFloat(patternStats[pattern])
		if high > winnerHigh || (high == winnerHigh && avg < winnerLatency) {
			winner = pattern
			winnerHigh = high
			winnerLatency = avg
		}
	}
	if winner != "" {
		total := len(benchmarkQueries)
		fmt.Printf("BENCHMARK WINNER: %s\n", strings.ToUpper(winner))
		fmt.Println("════════════════════════════════════════════════════════════════")
		fmt.Printf("✓ %.1fms avg latency\n", winnerLatency)
		fmt.Printf("✓ %d/%d queries returned keyword-matched results\n", winnerHigh, total)
	}
	fmt.Println()
}

func averageFloat(vals []float64) float64 {
	if len(vals) == 0 {
		return 0
	}
	sum := 0.0
	for _, v := range vals {
		sum += v
	}
	return sum / float64(len(vals))
}

func minFloat(vals []float64) float64 {
	if len(vals) == 0 {
		return 0
	}
	min := vals[0]
	for _, v := range vals {
		if v < min {
			min = v
		}
	}
	return min
}

func maxFloat(vals []float64) float64 {
	if len(vals) == 0 {
		return 0
	}
	max := vals[0]
	for _, v := range vals {
		if v > max {
			max = v
		}
	}
	return max
}

// ─── Token-Savings Benchmark ─────────────────────────────────────────────────

type tokenBenchResult struct {
	Query        string
	Intent       string
	ResultCount  int
	RawTokens    int
	ShapedTokens int
}

func runTokenBenchmark(db *sqlite.IndexDB, exportPath string) error {
	fmt.Println("\n╔════════════════════════════════════════════════════════════════╗")
	fmt.Println("║       TOKEN-SAVINGS BENCHMARK: RAW vs SHAPED CONTEXT          ║")
	fmt.Println("╚════════════════════════════════════════════════════════════════╝")
	const budget = 2000
	fmt.Printf("Shape budget: %d tokens  |  Queries: %d\n\n", budget, len(benchmarkQueries))

	hasVectors, _ := db.HasVectors()
	results := make([]tokenBenchResult, 0, len(benchmarkQueries))

	for _, q := range benchmarkQueries {
		rows, err := db.QueryChunksLexical(q.Query, 30)
		if err != nil || len(rows) == 0 {
			results = append(results, tokenBenchResult{Query: q.Query, Intent: q.Intent})
			fmt.Printf("  %-45s  (no results)\n", q.Intent)
			continue
		}
		if hasVectors {
			rows = appsearch.HybridRankWithEmbedder(rows, q.Query, db, 10, 0, inframodels.NewEmbedder)
		} else if len(rows) > 10 {
			rows = rows[:10]
		}

		domainResults := make([]domainsearch.Result, len(rows))
		rawTokens := 0
		for i, r := range rows {
			text := getStr(r, "text")
			var score float32
			if v, ok := r["_hybrid_score"].(float64); ok {
				score = float32(v)
			}
			domainResults[i] = domainsearch.Result{
				Source:    getStr(r, "path"),
				Title:     getStr(r, "folder_role"),
				Content:   text,
				Relevance: score,
			}
			rawTokens += (len(text) + 3) / 4
		}

		block := appcontext.FormatAsMarkdown(domainResults, q.Query, budget)
		saved := rawTokens - block.TokenCount
		pct := 0.0
		if rawTokens > 0 {
			pct = float64(saved) / float64(rawTokens) * 100
		}
		results = append(results, tokenBenchResult{
			Query:        q.Query,
			Intent:       q.Intent,
			ResultCount:  len(rows),
			RawTokens:    rawTokens,
			ShapedTokens: block.TokenCount,
		})
		fmt.Printf("  %-45s  raw=%4d  shaped=%4d  saved=%4d (%3.0f%%)\n",
			q.Intent, rawTokens, block.TokenCount, saved, pct)
	}

	report := formatTokenReport(results, budget)
	fmt.Print("\n" + report)

	if exportPath != "" {
		if err := os.WriteFile(exportPath, []byte(report), 0o644); err != nil {
			return fmt.Errorf("export: %w", err)
		}
		fmt.Printf("Evidence written to: %s\n", exportPath)
	}
	return nil
}

func formatTokenReport(results []tokenBenchResult, budget int) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "# Token-Savings Benchmark\n\n")
	fmt.Fprintf(&sb, "Shape budget: **%d tokens** | Token estimate: `(len(text)+3)/4`\n\n", budget)
	fmt.Fprintln(&sb, "| # | Query (intent) | Results | Raw tokens | Shaped tokens | Saved | % saved |")
	fmt.Fprintln(&sb, "|---|----------------|---------|------------|---------------|-------|---------|")

	totalRaw, totalShaped := 0, 0
	for i, r := range results {
		saved := r.RawTokens - r.ShapedTokens
		pct := 0.0
		if r.RawTokens > 0 {
			pct = float64(saved) / float64(r.RawTokens) * 100
		}
		fmt.Fprintf(&sb, "| %d | %s | %d | %d | %d | %d | %.0f%% |\n",
			i+1, r.Intent, r.ResultCount, r.RawTokens, r.ShapedTokens, saved, pct)
		totalRaw += r.RawTokens
		totalShaped += r.ShapedTokens
	}

	totalSaved := totalRaw - totalShaped
	totalPct := 0.0
	if totalRaw > 0 {
		totalPct = float64(totalSaved) / float64(totalRaw) * 100
	}
	fmt.Fprintf(&sb, "| — | **TOTAL** | — | **%d** | **%d** | **%d** | **%.0f%%** |\n\n",
		totalRaw, totalShaped, totalSaved, totalPct)

	fmt.Fprintln(&sb, "_Context shaping applies deduplication + token-budget truncation._")
	fmt.Fprintln(&sb, "_Raw tokens = sum of `len(chunk text)/4` across all ranked results._")
	fmt.Fprintln(&sb, "_Shaped tokens = `ContextBlock.TokenCount` after `FormatAsMarkdown`._")
	fmt.Fprintln(&sb, "_Negative savings = sparse result set already under budget; shaper adds Markdown formatting overhead._")
	return sb.String()
}

func getStr(m map[string]interface{}, key string) string {
	if v, ok := m[key]; ok && v != nil {
		if s, ok := v.(*string); ok && s != nil {
			return *s
		}
		return fmt.Sprintf("%v", v)
	}
	return ""
}

// ─── Reshape Token-Savings Benchmark ─────────────────────────────────────────

type reshapeBenchResult struct {
	Request      string
	RawTokens    int
	ShapedTokens int
	ErrNote      string
}

var reshapeTestCases = []struct{ Request, Context string }{
	{"add a --dry-run flag to work normalize that shows what would change without making changes", ""},
	{"the hawp search benchmark is slow, optimize it to run faster", ""},
	{"document the mcp server tools with examples for each tool and the expected input/output format", ""},
	{"migrate the TypeScript scripts to Go and deprecate the npm-based workflow entirely", ""},
	{"add support for selecting multiple providers in a single hawp init command", ""},
	{"the kit validate command should check for broken links inside kit files not just work files", ""},
	{"implement context deduplication to avoid sending the same chunks twice in a single search response", "already implemented Jaccard dedup in the context pipeline"},
	{"improve error messages when index.sqlite does not exist or is corrupted", ""},
	{"add a hawp providers list command that shows installed providers and their current versions", ""},
	{"the work normalize folder migration reports false positives for legacy slug IDs in parked rows", "fixed in the recent WIP commit but not yet validated against a fresh clone"},
}

func runReshapeBenchmark(backend, url, model, exportPath string) error {
	fmt.Println("\n╔════════════════════════════════════════════════════════════════╗")
	fmt.Println("║      RESHAPE TOKEN-SAVINGS BENCHMARK: RAW vs SHAPED INTAKE    ║")
	fmt.Println("╚════════════════════════════════════════════════════════════════╝")
	fmt.Printf("Backend: %s", backend)
	if model != "" {
		fmt.Printf(" | Model: %s", model)
	}
	fmt.Printf("  |  Queries: %d\n\n", len(reshapeTestCases))

	shaper, cleanup, err := newReshapeShaper(backend, url, model)
	if err != nil {
		return fmt.Errorf("initialize reshape shaper: %w", err)
	}
	defer cleanup()

	ctx := context.Background()
	results := make([]reshapeBenchResult, 0, len(reshapeTestCases))

	for _, tc := range reshapeTestCases {
		rawTokens := (len(tc.Request) + len(tc.Context) + 3) / 4

		draft, shapeErr := appintake.DraftIntake(ctx, appintake.DraftRequest{
			Input:   tc.Request,
			Context: tc.Context,
		}, shaper)

		var shapedTokens int
		var errNote string
		if shapeErr != nil {
			errNote = shapeErr.Error()
			shapedTokens = 0
		} else {
			shapedTokens = (len(draft.Mission) + len(draft.Constraints) + len(draft.Output) + len(draft.Checkpoint) + 3) / 4
		}

		saved := rawTokens - shapedTokens
		pct := 0.0
		if rawTokens > 0 && errNote == "" {
			pct = float64(saved) / float64(rawTokens) * 100
		}

		short := tc.Request
		if len(short) > 45 {
			short = short[:42] + "..."
		}
		if errNote != "" {
			fmt.Printf("  %-45s  ERROR: %s\n", short, errNote)
		} else {
			fmt.Printf("  %-45s  raw=%4d  shaped=%4d  saved=%4d (%+.0f%%)\n",
				short, rawTokens, shapedTokens, saved, pct)
		}

		results = append(results, reshapeBenchResult{
			Request:      tc.Request,
			RawTokens:    rawTokens,
			ShapedTokens: shapedTokens,
			ErrNote:      errNote,
		})
	}

	report := formatReshapeReport(results, backend, model)
	fmt.Print("\n" + report)

	if exportPath != "" {
		if err := os.MkdirAll(filepath.Dir(exportPath), 0o755); err != nil {
			return fmt.Errorf("create export directory: %w", err)
		}
		if err := os.WriteFile(exportPath, []byte(report), 0o644); err != nil {
			return fmt.Errorf("export: %w", err)
		}
		fmt.Printf("Evidence written to: %s\n", exportPath)
	}
	return nil
}

func formatReshapeReport(results []reshapeBenchResult, backend, model string) string {
	var sb strings.Builder
	modelNote := backend
	if model != "" {
		modelNote = backend + " / " + model
	}
	fmt.Fprintf(&sb, "# Reshape Token-Savings Benchmark\n\n")
	fmt.Fprintf(&sb, "Backend: **%s** | Token estimate: `(len(text)+3)/4` | Date: 2026-09-11\n\n", modelNote)
	fmt.Fprintln(&sb, "| # | Request (truncated) | Raw tokens | Shaped tokens | Saved | % saved | Note |")
	fmt.Fprintln(&sb, "|---|---------------------|------------|---------------|-------|---------|------|")

	totalRaw, totalShaped, errorCount := 0, 0, 0
	for i, r := range results {
		short := r.Request
		if len(short) > 50 {
			short = short[:47] + "..."
		}
		note := ""
		if r.ErrNote != "" {
			note = "error: " + r.ErrNote
			errorCount++
			fmt.Fprintf(&sb, "| %d | %s | %d | — | — | — | %s |\n", i+1, short, r.RawTokens, note)
			totalRaw += r.RawTokens
			continue
		}
		saved := r.RawTokens - r.ShapedTokens
		pct := 0.0
		if r.RawTokens > 0 {
			pct = float64(saved) / float64(r.RawTokens) * 100
		}
		fmt.Fprintf(&sb, "| %d | %s | %d | %d | %+d | %+.0f%% | |\n",
			i+1, short, r.RawTokens, r.ShapedTokens, saved, pct)
		totalRaw += r.RawTokens
		totalShaped += r.ShapedTokens
	}

	successCount := len(results) - errorCount
	totalSaved := totalRaw - totalShaped
	totalPct := 0.0
	if totalRaw > 0 && successCount > 0 {
		totalPct = float64(totalSaved) / float64(totalRaw) * 100
	}
	fmt.Fprintf(&sb, "| — | **TOTAL (%d/%d succeeded)** | **%d** | **%d** | **%+d** | **%+.0f%%** | |\n\n",
		successCount, len(results), totalRaw, totalShaped, totalSaved, totalPct)

	fmt.Fprintln(&sb, "_Raw tokens = `(len(request)+len(context)+3)/4` on the verbatim user input._")
	fmt.Fprintln(&sb, "_Shaped tokens = `(len(mission)+len(constraints)+len(output)+len(checkpoint)+3)/4` on the DraftIntake output._")
	fmt.Fprintln(&sb, "_Negative savings (expansion) is expected for short requests: structured intake adds labeled fields._")
	fmt.Fprintln(&sb, "_The value of reshaping is precision and downstream filtering, not raw token compression._")
	fmt.Fprintln(&sb, "_v0.1.0 gate: avg shaped tokens < avg raw tokens (any net savings across the 10-query suite)._")
	return sb.String()
}

// ─── Downstream Savings Benchmark ────────────────────────────────────────────
// Measures: (user request + retrieved search context) vs (shaped HAWP draft fields).
// This models the real downstream LLM call: without reshape, the caller passes
// request+context; with reshape, only the compact draft fields are needed.

type downstreamBenchResult struct {
	Intent       string
	RawTokens    int // request + formatted search context
	ShapedTokens int // mission + constraints + output + checkpoint
	ErrNote      string
}

func runDownstreamBenchmark(backend, url, model string, db *sqlite.IndexDB, exportPath string) error {
	fmt.Println("\n╔════════════════════════════════════════════════════════════════╗")
	fmt.Println("║   DOWNSTREAM SAVINGS BENCHMARK: REQUEST+CONTEXT vs DRAFT      ║")
	fmt.Println("╚════════════════════════════════════════════════════════════════╝")
	fmt.Printf("Backend: %s", backend)
	if model != "" {
		fmt.Printf(" | Model: %s", model)
	}
	fmt.Printf("  |  Queries: %d\n\n", len(benchmarkQueries))

	shaper, cleanup, err := newReshapeShaper(backend, url, model)
	if err != nil {
		return fmt.Errorf("initialize downstream shaper: %w", err)
	}
	defer cleanup()

	const budget = 2000
	hasVectors, _ := db.HasVectors()
	ctx := context.Background()
	results := make([]downstreamBenchResult, 0, len(benchmarkQueries))

	for _, q := range benchmarkQueries {
		// Retrieve and shape the search context (same path as token benchmark).
		rows, err := db.QueryChunksLexical(q.Query, 30)
		if err != nil || len(rows) == 0 {
			results = append(results, downstreamBenchResult{Intent: q.Intent, ErrNote: "no search results"})
			fmt.Printf("  %-45s  (no search results)\n", q.Intent)
			continue
		}
		if hasVectors {
			rows = appsearch.HybridRankWithEmbedder(rows, q.Query, db, 10, 0, inframodels.NewEmbedder)
		} else if len(rows) > 10 {
			rows = rows[:10]
		}

		domainResults := make([]domainsearch.Result, len(rows))
		for i, r := range rows {
			text := getStr(r, "text")
			var score float32
			if v, ok := r["_hybrid_score"].(float64); ok {
				score = float32(v)
			}
			domainResults[i] = domainsearch.Result{
				Source:    getStr(r, "path"),
				Title:     getStr(r, "folder_role"),
				Content:   text,
				Relevance: score,
			}
		}
		block := appcontext.FormatAsMarkdown(domainResults, q.Query, budget)
		contextText := block.String()

		// Raw: what a downstream LLM would receive without reshape.
		rawTokens := (len(q.Query) + len(contextText) + 3) / 4

		// Shaped: run DraftIntake with the real retrieved context as input context.
		draft, shapeErr := appintake.DraftIntake(ctx, appintake.DraftRequest{
			Input:   q.Query,
			Context: contextText,
		}, shaper)

		var shapedTokens int
		var errNote string
		if shapeErr != nil {
			errNote = shapeErr.Error()
		} else {
			shapedTokens = (len(draft.Mission) + len(draft.Constraints) + len(draft.Output) + len(draft.Checkpoint) + 3) / 4
		}

		saved := rawTokens - shapedTokens
		pct := 0.0
		if rawTokens > 0 && errNote == "" {
			pct = float64(saved) / float64(rawTokens) * 100
		}

		if errNote != "" {
			fmt.Printf("  %-45s  raw=%4d  ERROR: %s\n", q.Intent, rawTokens, errNote)
		} else {
			fmt.Printf("  %-45s  raw=%4d  shaped=%4d  saved=%4d (%+.0f%%)\n",
				q.Intent, rawTokens, shapedTokens, saved, pct)
		}

		results = append(results, downstreamBenchResult{
			Intent:       q.Intent,
			RawTokens:    rawTokens,
			ShapedTokens: shapedTokens,
			ErrNote:      errNote,
		})
	}

	report := formatDownstreamReport(results, backend, model)
	fmt.Print("\n" + report)

	if exportPath != "" {
		if err := os.MkdirAll(filepath.Dir(exportPath), 0o755); err != nil {
			return fmt.Errorf("create export directory: %w", err)
		}
		if err := os.WriteFile(exportPath, []byte(report), 0o644); err != nil {
			return fmt.Errorf("export: %w", err)
		}
		fmt.Printf("Evidence written to: %s\n", exportPath)
	}
	return nil
}

func formatDownstreamReport(results []downstreamBenchResult, backend, model string) string {
	var sb strings.Builder
	modelNote := backend
	if model != "" {
		modelNote = backend + " / " + model
	}
	fmt.Fprintf(&sb, "# Downstream Savings Benchmark\n\n")
	fmt.Fprintf(&sb, "Backend: **%s** | Token estimate: `(len(text)+3)/4` | Date: 2026-09-12\n\n", modelNote)
	fmt.Fprintln(&sb, "| # | Query (intent) | Raw (req+ctx) tokens | Shaped tokens | Saved | % saved |")
	fmt.Fprintln(&sb, "|---|----------------|----------------------|---------------|-------|---------|")

	totalRaw, totalShaped, errorCount := 0, 0, 0
	for i, r := range results {
		if r.ErrNote != "" {
			errorCount++
			fmt.Fprintf(&sb, "| %d | %s | %d | — | — | — |\n", i+1, r.Intent, r.RawTokens)
			totalRaw += r.RawTokens
			continue
		}
		saved := r.RawTokens - r.ShapedTokens
		pct := 0.0
		if r.RawTokens > 0 {
			pct = float64(saved) / float64(r.RawTokens) * 100
		}
		fmt.Fprintf(&sb, "| %d | %s | %d | %d | %+d | %+.0f%% |\n",
			i+1, r.Intent, r.RawTokens, r.ShapedTokens, saved, pct)
		totalRaw += r.RawTokens
		totalShaped += r.ShapedTokens
	}

	successCount := len(results) - errorCount
	totalSaved := totalRaw - totalShaped
	totalPct := 0.0
	if totalRaw > 0 && successCount > 0 {
		totalPct = float64(totalSaved) / float64(totalRaw) * 100
	}
	fmt.Fprintf(&sb, "| — | **TOTAL (%d/%d succeeded)** | **%d** | **%d** | **%+d** | **%+.0f%%** |\n\n",
		successCount, len(results), totalRaw, totalShaped, totalSaved, totalPct)

	fmt.Fprintln(&sb, "_Raw tokens = `(len(query) + len(formatted_search_context) + 3) / 4`._")
	fmt.Fprintln(&sb, "_This models the downstream LLM call without reshape: request text + retrieved context._")
	fmt.Fprintln(&sb, "_Shaped tokens = `(len(mission)+len(constraints)+len(output)+len(checkpoint)+3)/4` on the DraftIntake output._")
	fmt.Fprintln(&sb, "_Positive savings = shaped draft is more compact than request + context (downstream token reduction)._")
	fmt.Fprintln(&sb, "_v0.1.0 gate: ≥20% average savings across succeeded queries._")
	return sb.String()
}
