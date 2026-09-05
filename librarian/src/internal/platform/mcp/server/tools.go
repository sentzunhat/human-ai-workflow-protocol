package mcp

import (
	"encoding/json"
	"os"

	"github.com/sentzunhat/hawp/librarian/src/internal/infrastructure/filesystem"
	usageinfra "github.com/sentzunhat/hawp/librarian/src/internal/infrastructure/repositories/usage"
)

func toolDefs() []map[string]any {
	return []map[string]any{
		searchToolDef(),
		usageToolDef(),
		{
			"name":        "hawp_work_intake",
			"description": "Run HAWP search and request reshape in one compound intake call. Returns structured draft fields, retrieval status, warnings/questions, and token accounting. Use this before hawp_work_new so agents do not skip context retrieval. Requires a local search index and Ollama running locally; shape failures return a structured blocked state.",
			"inputSchema": map[string]any{
				"type": "object",
				"properties": map[string]any{
					"input": map[string]any{
						"type":        "string",
						"description": "Verbatim user request to retrieve context for and shape",
					},
					"limit": map[string]any{
						"type":        "integer",
						"description": "Max search results to retrieve before context shaping (default 10)",
					},
					"max_tokens": map[string]any{
						"type":        "integer",
						"description": "Token budget for retrieved context block (default 2000)",
					},
					"model": map[string]any{
						"type":        "string",
						"description": "Ollama model name (default: mistral)",
					},
					"url": map[string]any{
						"type":        "string",
						"description": "Ollama server URL (default: http://localhost:11434)",
					},
				},
				"required": []string{"input"},
			},
		},
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
		{
			"name":        "hawp_work_doc",
			"description": "Create a secondary work document (status report, evidence file, decision record, or note) at the canonical UUID-subfolder path: {type}/YYYY/MM/DD/{id}/{type}.md. Pass work_item_id to use the work item's own ID as the folder name so all artifacts for that item share the same path prefix. Returns the created file path so the caller can append content without knowing path conventions.",
			"inputSchema": map[string]any{
				"type": "object",
				"properties": map[string]any{
					"type": map[string]any{
						"type":        "string",
						"description": "Document type: status|evidence|decision|note",
					},
					"title": map[string]any{
						"type":        "string",
						"description": "Short label for the document (recorded in front-matter)",
					},
					"work_item_id": map[string]any{
						"type":        "string",
						"description": "Work item ID to link this document to (8-char short or full UUID). When provided, the document folder uses this ID instead of a fresh UUID, grouping all artifacts for that work item together.",
					},
				},
				"required": []string{"type", "title"},
			},
		},
		{
			"name":        "hawp_work_reshape",
			"description": "Shape a raw user request into HAWP intake fields (mission, constraints, output, checkpoint) using a local Ollama model. Returns a proposed shape for review — does NOT create a work item. Requires Ollama running locally.",
			"inputSchema": map[string]any{
				"type": "object",
				"properties": map[string]any{
					"input": map[string]any{
						"type":        "string",
						"description": "Verbatim user request to shape",
					},
					"context": map[string]any{
						"type":        "string",
						"description": "Optional background context the user supplied",
					},
					"model": map[string]any{
						"type":        "string",
						"description": "Ollama model name (default: mistral)",
					},
					"url": map[string]any{
						"type":        "string",
						"description": "Ollama server URL (default: http://localhost:11434)",
					},
				},
				"required": []string{"input"},
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
	case "hawp_work_intake":
		resp = toolWorkIntake(p.Arguments, repoRoot)
	case "hawp_work_new":
		resp = toolWorkNew(p.Arguments, repoRoot)
	case "hawp_work_validate":
		resp = toolWorkValidate(repoRoot)
	case "hawp_work_doc":
		resp = toolWorkDoc(p.Arguments, repoRoot)
	case "hawp_work_reshape":
		resp = toolWorkReshape(p.Arguments, repoRoot)
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
	cfg := usageinfra.LoadConfig(h.UsageConfigFile)
	if !cfg.Enabled {
		return
	}
	store, err := usageinfra.Open(h.UsageDB)
	if err != nil {
		return
	}
	defer store.Close()

	outJSON, _ := json.Marshal(resp.Result)
	_ = store.Write(tool, inputArgs, outJSON, cfg.LogBodies)
}
