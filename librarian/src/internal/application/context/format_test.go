package context

import (
	"strings"
	"testing"

	"github.com/sentzunhat/hawp/librarian/src/internal/domain/search"
)

func TestFormatAsMarkdown(t *testing.T) {
	results := []search.Result{
		{
			Content:   "Vector embeddings are dense representations",
			Source:    "README.md",
			Title:     "Embeddings Intro",
			Relevance: 0.95,
			LineStart: 10,
			LineEnd:   18,
		},
		{
			Content:   "ONNX Runtime provides efficient inference",
			Source:    "guide.md",
			Title:     "ONNX Integration",
			Relevance: 0.87,
			LineStart: 30,
			LineEnd:   35,
		},
	}

	block := FormatAsMarkdown(results, "vector embedding", 500)

	if block.Title == "" {
		t.Error("Title should not be empty")
	}

	if block.ResultCount != 2 {
		t.Errorf("ResultCount = %d, want 2", block.ResultCount)
	}

	if block.Query != "vector embedding" {
		t.Errorf("Query = %q, want %q", block.Query, "vector embedding")
	}

	if len(block.Results) != 2 {
		t.Errorf("Results length = %d, want 2", len(block.Results))
	}

	// Check results are sorted by relevance (descending)
	if block.Results[0].Relevance < block.Results[1].Relevance {
		t.Error("Results not sorted by relevance descending")
	}
	if block.Results[0].LineStart == 0 || block.Results[0].LineEnd == 0 {
		t.Error("Line ranges should be preserved in formatted results")
	}
}

func TestFormatAsMarkdownRespectsBudget(t *testing.T) {
	results := []search.Result{
		{
			Content:   strings.Repeat("word ", 100), // ~100 words, ~400 tokens
			Source:    "a.md",
			Relevance: 0.9,
		},
		{
			Content:   strings.Repeat("word ", 100), // ~100 words, ~400 tokens
			Source:    "b.md",
			Relevance: 0.8,
		},
	}

	// With budget of 300 tokens, only first result should fit (partially truncated)
	block := FormatAsMarkdown(results, "test", 300)

	if block.TokenCount > 300 {
		t.Errorf("TokenCount %d exceeds budget 300", block.TokenCount)
	}
}

func TestFormatAsMarkdownEmpty(t *testing.T) {
	block := FormatAsMarkdown([]search.Result{}, "query", 1000)

	if block.ResultCount != 0 {
		t.Errorf("Empty results should have ResultCount 0, got %d", block.ResultCount)
	}

	if len(block.Results) != 0 {
		t.Errorf("Empty results should have length 0, got %d", len(block.Results))
	}
}

func TestEstimateTokens(t *testing.T) {
	tests := []struct {
		text string
		want int
	}{
		{"hello world", 3},           // 11 chars ÷ 4 = 2.75 → 3
		{"a", 1},                     // 1 char ÷ 4 = 0.25 → 1
		{"", 0},                      // 0 chars
		{"this is a longer text", 6}, // 21 chars ÷ 4 = 5.25 → 6
	}

	for _, tt := range tests {
		got := estimateTokens(tt.text)
		if got != tt.want {
			t.Errorf("estimateTokens(%q) = %d, want %d", tt.text, got, tt.want)
		}
	}
}

func TestTruncateToTokens(t *testing.T) {
	tests := []struct {
		text      string
		maxTokens int
		check     func(string) bool
	}{
		{
			text:      "short text",
			maxTokens: 100,
			check: func(s string) bool {
				return s == "short text" // Should not truncate
			},
		},
		{
			text:      strings.Repeat("word ", 100),
			maxTokens: 50,
			check: func(s string) bool {
				tokens := estimateTokens(s)
				return tokens <= 55 && strings.HasSuffix(s, "...")
			},
		},
	}

	for i, tt := range tests {
		got := truncateToTokens(tt.text, tt.maxTokens)
		if !tt.check(got) {
			t.Errorf("Test %d: truncateToTokens result unexpected", i)
		}
	}
}

func TestContextBlockString(t *testing.T) {
	block := ContextBlock{
		Title:       "Test Results",
		ResultCount: 1,
		TokenCount:  100,
		Results: []FormattedResult{
			{
				Rank:      1,
				Relevance: 0.95,
				Source:    "test.md",
				Title:     "Test Title",
				Content:   "Test content",
				Tokens:    5,
				LineStart: 7,
				LineEnd:   9,
			},
		},
	}

	str := block.String()

	if !strings.Contains(str, "Test Results") {
		t.Error("String output should contain title")
	}

	if !strings.Contains(str, "Test Title") {
		t.Error("String output should contain result title")
	}

	if !strings.Contains(str, "Test content") {
		t.Error("String output should contain result content")
	}

	if !strings.Contains(str, "95%") {
		t.Error("String output should contain relevance percentage")
	}
	if !strings.Contains(str, "test.md:7-9") {
		t.Error("String output should contain line-ranged source path")
	}
}

func TestContextBlockReferences(t *testing.T) {
	// Test case: 3 results from 2 unique sources → References should have 2 entries (deduplicated)
	results := []search.Result{
		{
			Content:   "First result from deployment docs",
			Source:    "deployment.md",
			Title:     "Deployment Guide",
			Relevance: 0.95,
			LineStart: 4,
			LineEnd:   12,
		},
		{
			Content:   "Architecture overview",
			Source:    "architecture.md",
			Title:     "System Architecture",
			Relevance: 0.82,
			LineStart: 20,
			LineEnd:   25,
		},
		{
			Content:   "More deployment info",
			Source:    "deployment.md",
			Title:     "Deployment Setup",
			Relevance: 0.78,
			LineStart: 30,
			LineEnd:   34,
		},
	}

	block := FormatAsMarkdown(results, "deployment", 1000)

	// Should have 3 results
	if block.ResultCount != 3 {
		t.Errorf("ResultCount = %d, want 3", block.ResultCount)
	}

	// Should have 2 references (deduplicated)
	if len(block.References) != 2 {
		t.Errorf("References length = %d, want 2", len(block.References))
	}

	// Check references are sorted by relevance descending
	if len(block.References) >= 2 {
		if block.References[0].Relevance < block.References[1].Relevance {
			t.Error("References should be sorted by relevance descending")
		}
	}

	// Verify reference sources
	sourceMap := make(map[string]bool)
	for _, ref := range block.References {
		sourceMap[ref.Source] = true
	}
	if !sourceMap["deployment.md"] {
		t.Error("References should contain deployment.md")
	}
	if !sourceMap["architecture.md"] {
		t.Error("References should contain architecture.md")
	}

	// deployment.md should have the highest relevance from its results (0.95 > 0.78)
	for _, ref := range block.References {
		if ref.Source == "deployment.md" && ref.Relevance != 0.95 {
			t.Errorf("deployment.md relevance should be 0.95 (highest), got %f", ref.Relevance)
		}
		if ref.Source == "deployment.md" && (ref.LineStart != 4 || ref.LineEnd != 12) {
			t.Errorf("deployment.md line range = %d-%d, want 4-12", ref.LineStart, ref.LineEnd)
		}
	}
}

func TestContextBlockStringInterleavesReferences(t *testing.T) {
	// Each result's **Reference:** line must appear immediately above its own
	// content — not collected into one list at the end — so a reader sees
	// which source a chunk came from right where it's used.
	block := ContextBlock{
		Title:       "Test Results",
		ResultCount: 2,
		TokenCount:  150,
		Results: []FormattedResult{
			{
				Rank:      1,
				Relevance: 0.95,
				Source:    "deployment.md",
				Title:     "Deployment Guide",
				Content:   "Deployment info",
				Tokens:    50,
				LineStart: 8,
				LineEnd:   14,
			},
			{
				Rank:      2,
				Relevance: 0.82,
				Source:    "architecture.md",
				Title:     "Architecture",
				Content:   "Architecture info",
				Tokens:    50,
				LineStart: 21,
				LineEnd:   24,
			},
		},
	}

	str := block.String()

	if !strings.Contains(str, "**Reference:** deployment.md:8-14") {
		t.Error("String output should contain an inline Reference line for deployment.md")
	}
	if !strings.Contains(str, "**Reference:** architecture.md:21-24") {
		t.Error("String output should contain an inline Reference line for architecture.md")
	}

	// The deployment.md reference must precede its own content, and precede
	// the architecture.md reference (interleaved order, not batched at end).
	refDeployIdx := strings.Index(str, "**Reference:** deployment.md:8-14")
	contentDeployIdx := strings.Index(str, "Deployment info")
	refArchIdx := strings.Index(str, "**Reference:** architecture.md:21-24")
	contentArchIdx := strings.Index(str, "Architecture info")

	if refDeployIdx < 0 || contentDeployIdx < 0 || refArchIdx < 0 || contentArchIdx < 0 {
		t.Fatalf("expected all reference/content markers present, got: %s", str)
	}
	if !(refDeployIdx < contentDeployIdx && contentDeployIdx < refArchIdx && refArchIdx < contentArchIdx) {
		t.Errorf("expected order ref1 < content1 < ref2 < content2, got positions %d,%d,%d,%d",
			refDeployIdx, contentDeployIdx, refArchIdx, contentArchIdx)
	}
}
