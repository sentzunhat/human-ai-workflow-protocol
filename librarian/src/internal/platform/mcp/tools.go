package mcp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strings"

	appcheck "github.com/sentzunhat/hawp/librarian/src/internal/application/check"
	appsearch "github.com/sentzunhat/hawp/librarian/src/internal/application/search"
	appuuid "github.com/sentzunhat/hawp/librarian/src/internal/application/uuidgen"
	appwork "github.com/sentzunhat/hawp/librarian/src/internal/application/work"
)

type toolContent struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type toolResult struct {
	Content []toolContent `json:"content"`
	IsError bool          `json:"isError,omitempty"`
}

func text(s string) rpcResponse {
	return rpcResponse{Result: toolResult{Content: []toolContent{{Type: "text", Text: s}}}}
}

func toolErr(s string) rpcResponse {
	return rpcResponse{Result: toolResult{
		Content: []toolContent{{Type: "text", Text: s}},
		IsError: true,
	}}
}

func toolDefs() []map[string]any {
	return []map[string]any{
		{
			"name":        "hawp_search",
			"description": "Search indexed HAWP kit and work documents. Returns ranked chunks with source paths and relevance scores. Run `hawp search index` then `hawp search embed` first to build the index.",
			"inputSchema": map[string]any{
				"type": "object",
				"properties": map[string]any{
					"query": map[string]any{
						"type":        "string",
						"description": "Search query",
					},
					"limit": map[string]any{
						"type":        "integer",
						"description": "Max results to return (default 5)",
					},
				},
				"required": []string{"query"},
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

func toolSearch(args json.RawMessage, repoRoot string) rpcResponse {
	var a struct {
		Query string `json:"query"`
		Limit int    `json:"limit"`
	}
	if err := json.Unmarshal(args, &a); err != nil {
		return toolErr("invalid args: " + err.Error())
	}
	if a.Query == "" {
		return toolErr("query is required")
	}
	if a.Limit <= 0 {
		a.Limit = 5
	}

	results, err := appsearch.Query(repoRoot, a.Query, a.Limit)
	if err != nil {
		return toolErr("search failed: " + err.Error())
	}
	if len(results) == 0 {
		return text("No results found for: " + a.Query)
	}

	var sb strings.Builder
	for i, r := range results {
		fmt.Fprintf(&sb, "## [%d] %s\n", i+1, r.Source)
		fmt.Fprintf(&sb, "relevance: %.3f\n\n", r.Relevance)
		fmt.Fprintf(&sb, "%s\n\n---\n\n", strings.TrimSpace(r.Content))
	}
	return text(strings.TrimRight(sb.String(), "\n-"))
}

func toolWorkNew(args json.RawMessage, repoRoot string) rpcResponse {
	var a struct {
		Title string `json:"title"`
		Type  string `json:"type"`
		Input string `json:"input"`
	}
	if err := json.Unmarshal(args, &a); err != nil {
		return toolErr("invalid args: " + err.Error())
	}
	if a.Title == "" {
		return toolErr("title is required")
	}
	if a.Type == "" {
		a.Type = "task"
	}

	workDir := filepath.Join(repoRoot, ".hawp", "work")
	result, err := appwork.NewItem(workDir, a.Type, a.Title, a.Input)
	if err != nil {
		return toolErr("work new failed: " + err.Error())
	}

	return text(fmt.Sprintf(
		"Created work item %s (%s)\nPlan: %s\nBacklog: %s\n\nNext: investigate and fill in the plan (see .hawp/kit/usage/intake-workflow.md).",
		appuuid.Short(result.UUID), result.Type,
		result.PlanFilePath, result.BacklogPath,
	))
}

func toolWorkValidate(repoRoot string) rpcResponse {
	var out, errOut bytes.Buffer
	code := appcheck.Run(&out, &errOut, repoRoot)

	combined := out.String()
	if errOut.Len() > 0 {
		combined += "\nstderr:\n" + errOut.String()
	}

	if code != 0 {
		return rpcResponse{Result: toolResult{
			Content: []toolContent{{Type: "text", Text: combined}},
			IsError: true,
		}}
	}
	return text(combined)
}
