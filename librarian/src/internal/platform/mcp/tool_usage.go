package mcp

import (
	"encoding/json"

	"github.com/sentzunhat/hawp/librarian/src/internal/domain/usage"
	"github.com/sentzunhat/hawp/librarian/src/internal/infrastructure/filesystem"
	"os"
)

func usageToolDef() map[string]any {
	return map[string]any{
		"name":        "hawp_usage",
		"description": "Return HAWP MCP call log statistics. Without arguments returns totals (calls, tokens in/out). Pass report:true for a full Markdown breakdown by tool and recent query list. Usage logging must be enabled first with `hawp usage enable`.",
		"inputSchema": map[string]any{
			"type": "object",
			"properties": map[string]any{
				"report": map[string]any{
					"type":        "boolean",
					"description": "When true, return a full Markdown usage report instead of a one-line totals summary",
				},
			},
		},
	}
}

func toolUsage(args json.RawMessage) rpcResponse {
	var a struct {
		Report bool `json:"report"`
	}
	if len(args) > 0 {
		if err := json.Unmarshal(args, &a); err != nil {
			return toolErr("invalid args: " + err.Error())
		}
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return toolErr("cannot determine home directory: " + err.Error())
	}
	h := filesystem.ResolveHawpHome(home)

	cfg := usage.LoadConfig(h.UsageConfigFile)
	if !cfg.Enabled {
		return text("Usage logging is disabled. Run `hawp usage enable` to start recording MCP calls.")
	}

	store, err := usage.Open(h.UsageDB)
	if err != nil {
		return toolErr("cannot open usage DB: " + err.Error())
	}
	defer store.Close()

	if a.Report {
		rep, err := store.GetReport()
		if err != nil {
			return toolErr("GetReport: " + err.Error())
		}
		return text(usage.FormatReport(rep))
	}

	totals, err := store.GetTotals()
	if err != nil {
		return toolErr("GetTotals: " + err.Error())
	}
	return text(usage.FormatTotals(totals))
}
