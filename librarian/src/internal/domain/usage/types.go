package usage

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// Config is the user preference for usage logging.
type Config struct {
	Enabled   bool `json:"enabled"`
	LogBodies bool `json:"log_bodies"`
}

// Entry is a recorded tool call.
type Entry struct {
	ID         int64
	TS         time.Time
	Tool       string
	QueryHash  string
	QueryText  *string
	TokensIn   int
	TokensOut  int
	InputBody  *string
	OutputBody *string
}

type Totals struct {
	Calls     int
	TokensIn  int
	TokensOut int
}

// Store is the usage-log port. Infrastructure owns concrete persistence.
type Store interface {
	Write(tool string, inputJSON, outputJSON []byte, logBodies bool) error
	Recent(n int) ([]Entry, error)
	GetTotals() (Totals, error)
	GetReport() (Report, error)
	Clear() error
	Close()
}

type ToolStat struct {
	Tool      string
	Calls     int
	TokensIn  int
	TokensOut int
}

type Report struct {
	Calls      int
	TokensIn   int
	TokensOut  int
	ByTool     []ToolStat
	TopEntries []Entry
	Since      *time.Time
	Until      *time.Time
}

// FormatReport renders a Report as human-readable Markdown.
func FormatReport(rep Report) string {
	if rep.Calls == 0 {
		return "No calls recorded. Run `hawp usage enable` then make some MCP tool calls.\n"
	}
	var sb strings.Builder
	fmt.Fprintln(&sb, "# HAWP Usage Report")
	fmt.Fprintln(&sb)
	if rep.Since != nil && rep.Until != nil {
		fmt.Fprintf(&sb, "Period: %s → %s\n\n", rep.Since.Format("2006-01-02 15:04 UTC"), rep.Until.Format("2006-01-02 15:04 UTC"))
	}
	fmt.Fprintln(&sb, "## Totals")
	fmt.Fprintln(&sb)
	fmt.Fprintf(&sb, "| Metric | Value |\n|--------|-------|\n| Calls | %d |\n| Tokens in (est.) | ~%d |\n| Tokens out (est.) | ~%d |\n", rep.Calls, rep.TokensIn, rep.TokensOut)
	saved := rep.TokensIn - rep.TokensOut
	if saved > 0 {
		fmt.Fprintf(&sb, "| Tokens saved (est.) | ~%d (~%.0f%%) |\n", saved, float64(saved)/float64(rep.TokensIn)*100)
	}
	fmt.Fprintln(&sb)
	if len(rep.ByTool) > 0 {
		fmt.Fprintln(&sb, "## By Tool")
		fmt.Fprintln(&sb)
		fmt.Fprintln(&sb, "| Tool | Calls | Tokens In | Tokens Out | Saved |\n|------|-------|-----------|------------|-------|")
		for _, ts := range rep.ByTool {
			s := ts.TokensIn - ts.TokensOut
			pct := 0.0
			if ts.TokensIn > 0 {
				pct = float64(s) / float64(ts.TokensIn) * 100
			}
			fmt.Fprintf(&sb, "| %s | %d | ~%d | ~%d | ~%d (~%.0f%%) |\n", ts.Tool, ts.Calls, ts.TokensIn, ts.TokensOut, s, pct)
		}
		fmt.Fprintln(&sb)
	}
	if len(rep.TopEntries) > 0 {
		fmt.Fprintln(&sb, "## Recent Queries")
		fmt.Fprintln(&sb)
		fmt.Fprintln(&sb, "| # | Time | Tool | Tokens In | Tokens Out | Query |\n|---|------|------|-----------|------------|-------|")
		for i, e := range rep.TopEntries {
			fmt.Fprintf(&sb, "| %d | %s | %s | %d | %d | %s |\n", i+1, e.TS.Format("01-02 15:04"), e.Tool, e.TokensIn, e.TokensOut, EntrySummary(e))
		}
		fmt.Fprintln(&sb)
	}
	fmt.Fprintln(&sb, "_Token estimates: chars/4 (MCP request/response JSON byte length)_")
	return sb.String()
}

func QuerySummary(inputJSON []byte, queryHash string) string {
	var m map[string]any
	if err := json.Unmarshal(inputJSON, &m); err != nil {
		return queryHash
	}
	for _, key := range []string{"query", "title"} {
		if value, ok := m[key].(string); ok && value != "" {
			if len(value) > 60 {
				return value[:57] + "..."
			}
			return value
		}
	}
	return queryHash
}

func EntrySummary(e Entry) string {
	if e.QueryText != nil && *e.QueryText != "" {
		s := *e.QueryText
		if len(s) > 60 {
			return s[:57] + "..."
		}
		return s
	}
	if e.InputBody != nil {
		return QuerySummary([]byte(*e.InputBody), e.QueryHash)
	}
	return e.QueryHash
}

func FormatTotals(t Totals) string {
	if t.Calls == 0 {
		return "No calls recorded. Run `hawp usage enable` then make some MCP tool calls."
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "Calls:       %d\nTokens in:   ~%d\nTokens out:  ~%d\n", t.Calls, t.TokensIn, t.TokensOut)
	saved := t.TokensIn - t.TokensOut
	if saved > 0 {
		fmt.Fprintf(&sb, "Saved:       ~%d (~%.0f%% reduction via context shaping)\n", saved, float64(saved)/float64(t.TokensIn)*100)
	}
	return sb.String()
}
