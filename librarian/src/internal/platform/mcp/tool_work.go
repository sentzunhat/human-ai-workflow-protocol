package mcp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"path/filepath"

	appcheck "github.com/sentzunhat/hawp/librarian/src/internal/application/check"
	appuuid "github.com/sentzunhat/hawp/librarian/src/internal/application/uuidgen"
	appwork "github.com/sentzunhat/hawp/librarian/src/internal/application/work"
)

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
