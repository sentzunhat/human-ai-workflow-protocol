package mcp

import "encoding/json"

func toolDefs() []map[string]any {
	return []map[string]any{
		searchToolDef(),
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
	switch p.Name {
	case "hawp_search":
		return toolSearch(p.Arguments, repoRoot)
	case "hawp_work_new":
		return toolWorkNew(p.Arguments, repoRoot)
	case "hawp_work_validate":
		return toolWorkValidate(repoRoot)
	default:
		return errResp(-32602, "unknown tool: "+p.Name)
	}
}
