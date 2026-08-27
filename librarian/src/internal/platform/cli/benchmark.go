// Package cli provides the benchmark command for search patterns.
package cli

import (
	"fmt"
	"os"
	"strings"
	"time"

	appcontext "github.com/sentzunhat/hawp/librarian/src/internal/application/context"
	appsearch "github.com/sentzunhat/hawp/librarian/src/internal/application/search"
	domainsearch "github.com/sentzunhat/hawp/librarian/src/internal/domain/search"
	"github.com/sentzunhat/hawp/librarian/src/internal/infrastructure/sqlite"
)

// BenchmarkQuery represents a test query with expected characteristics
type BenchmarkQuery struct {
	Query            string
	Intent           string // What semantic concept it's testing
	RelevantKeywords []string
}

// BenchmarkResult holds metrics for one query-pattern combination
type BenchmarkResult struct {
	Query            string
	Pattern          string // "lexical", "hybrid"
	LatencyMS        float64
	ResultCount      int
	TopResultQuality string // "high", "medium", "low"
}

// benchmarkQueries are 10 corpus-representative queries drawn from the actual
// .hawp/kit/ and .hawp/work/ document set. All queries are verified to return
// results against the current index.
var benchmarkQueries = []BenchmarkQuery{
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

// RunBenchmark executes benchmark tests across lexical, semantic, and hybrid patterns.
func RunBenchmark(db *sqlite.IndexDB) error {
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

	var results []BenchmarkResult
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

// benchmarkOneQuery runs a single query using the real search path for the pattern.
func benchmarkOneQuery(query BenchmarkQuery, pattern string, db *sqlite.IndexDB) BenchmarkResult {
	result := BenchmarkResult{Query: query.Query, Pattern: pattern}
	start := time.Now()

	switch pattern {
	case "lexical":
		rows, _ := db.QueryChunksLexical(query.Query, 10)
		result.ResultCount = len(rows)
		result.TopResultQuality = assessQuality(query.RelevantKeywords, rows)

	case "semantic":
		// Pure vector path: embed query → cosine rank all stored vectors → top-10.
		rows := appsearch.SemanticSearch(query.Query, db, 10)
		result.ResultCount = len(rows)
		result.TopResultQuality = assessQuality(query.RelevantKeywords, rows)

	case "hybrid":
		// Real hybrid path: lexical candidates → HybridRank (embeds query, cosine re-rank).
		rows, _ := db.QueryChunksLexical(query.Query, 30)
		if len(rows) > 0 {
			rows = appsearch.HybridRank(rows, query.Query, db, 10, 0)
		}
		result.ResultCount = len(rows)
		result.TopResultQuality = assessQuality(query.RelevantKeywords, rows)
	}

	result.LatencyMS = float64(time.Since(start).Milliseconds())
	return result
}

// assessQuality checks whether the top result's text contains any of the
// expected keywords. "high" = keyword found, "low" = no match or no results.
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

// printBenchmarkSummary prints the final comparison table
func printBenchmarkSummary(results []BenchmarkResult) {
	fmt.Println("\n╔════════════════════════════════════════════════════════════════╗")
	fmt.Println("║                      RESULTS SUMMARY                          ║")
	fmt.Println("╚════════════════════════════════════════════════════════════════╝")

	// Aggregate by pattern
	patternStats := make(map[string][]float64)
	patternQuality := make(map[string]map[string]int)

	for _, r := range results {
		patternStats[r.Pattern] = append(patternStats[r.Pattern], r.LatencyMS)
		if patternQuality[r.Pattern] == nil {
			patternQuality[r.Pattern] = make(map[string]int)
		}
		patternQuality[r.Pattern][r.TopResultQuality]++
	}

	// Print table header
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

	// Calculate relative performance
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

	// Data-driven winner: highest high-quality count; tie-break on latency.
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

// Helper functions for statistics
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

// tokenBenchResult records shaping savings for one query.
type tokenBenchResult struct {
	Query        string
	Intent       string
	ResultCount  int
	RawTokens    int
	ShapedTokens int
}

// RunTokenBenchmark measures how much context shaping reduces token count vs.
// returning all raw result text verbatim. For each benchmark query it:
//  1. Runs lexical (+ hybrid if vectors exist) search → sum raw content tokens
//  2. Runs FormatAsMarkdown with a 2000-token budget → reads shaped token count
//  3. Reports the difference as evidence of shaping effectiveness
//
// If exportPath is non-empty the Markdown report is also written to that file.
func RunTokenBenchmark(db *sqlite.IndexDB, exportPath string) error {
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
			rows = appsearch.HybridRank(rows, q.Query, db, 10, 0)
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

// formatTokenReport renders the token-savings results as a Markdown table
// suitable for committing as a work evidence artifact.
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
	return sb.String()
}
