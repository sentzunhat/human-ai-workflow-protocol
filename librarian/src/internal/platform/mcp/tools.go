package mcp

import (
	"encoding/json"
	"os"

	"github.com/sentzunhat/hawp/librarian/src/internal/domain/usage"
	"github.com/sentzunhat/hawp/librarian/src/internal/infrastructure/filesystem"
)

func toolDefs() []map[string]any {
	return []map[string]any{
		searchToolDef(),
		usageToolDef(),
		{
			"name":        "hawp_work_new",
			"description": "Create a new HAWP work item. Generates a UUID, writes active/{uuid}/plan.md from the intake template, and adds an inbox row to BACKLOG.md. Returns the UUID and plan file path.",
			"inputSchema": map[string]any{
				"type": "object",
				"properties": map[string]any{
					"title": map[string]any{
						"type":        "string",
						"description": "Work item title",
					},
					"type": map[string]any{
						"type":        "string",
						"description": "Item type: task|bug|feature|fix|test|improvement|infrastructure|release|decision (default: task)",
					},
					"input": map[string]any{
						"type":        "string",
						"description": "Verbatim original request, recorded in the plan's Input section (defaults to title)",
					},
				},
				"required": []string{"title"},
			},
		},
		{
			"name":        "hawp_work_validate",
			"description": "Validate HAWP kit structure, work item integrity, and local markdown links. Returns PASS/WARN/FAIL with a list of issues.",
			"inputSchema": map[string]any{
				"type":       "object",
				"properties": map[string]any{},
			},
		},
	}
}

func callTool(params json.RawMessage, repoRoot string) rpcResponse {
	var p struct {
		Name      string          `json:"name"`
		Arguments json.RawMessage `json:"arguments"`
	}
	if err := json.Unmarshal(params, &p); err != nil {
		return errResp(-32602, "invalid params: "+err.Error())
	}

	var resp rpcResponse
	switch p.Name {
	case "hawp_search":
		resp = toolSearch(p.Arguments, repoRoot)
	case "hawp_work_new":
		resp = toolWorkNew(p.Arguments, repoRoot)
	case "hawp_work_validate":
		resp = toolWorkValidate(repoRoot)
	case "hawp_usage":
		resp = toolUsage(p.Arguments)
	default:
		return errResp(-32602, "unknown tool: "+p.Name)
	}

	// Log synchronously: sqlite writes complete in <1ms and the goroutine
	// approach caused entries to be lost when the process exited before the
	// goroutine ran (observed with hawp_work_validate in short-lived sessions).
	logCall(p.Name, p.Arguments, resp)

	return resp
}

func logCall(tool string, inputArgs json.RawMessage, resp rpcResponse) {
	home, err := os.UserHomeDir()
	if err != nil {
		return
	}
	h := filesystem.ResolveHawpHome(home)
	cfg := usage.LoadConfig(h.UsageConfigFile)
	if !cfg.Enabled {
		return
	}
	store, err := usage.Open(h.UsageDB)
	if err != nil {
		return
	}
	defer store.Close()

	outJSON, _ := json.Marshal(resp.Result)
	_ = store.Write(tool, inputArgs, outJSON, cfg.LogBodies)
}
