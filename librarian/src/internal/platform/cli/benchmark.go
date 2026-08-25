// Package cli provides the benchmark command for search patterns.
package cli

import (
	"fmt"
	"strings"
	"time"

	appsearch "github.com/sentzunhat/hawp/librarian/src/internal/application/search"
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
			rows = appsearch.HybridRank(rows, query.Query, db, 10)
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
