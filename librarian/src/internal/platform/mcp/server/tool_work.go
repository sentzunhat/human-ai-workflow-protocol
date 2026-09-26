package mcp

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	appcheck "github.com/sentzunhat/hawp/librarian/src/internal/application/check"
	appsearch "github.com/sentzunhat/hawp/librarian/src/internal/application/search"
	appuuid "github.com/sentzunhat/hawp/librarian/src/internal/application/uuidgen"
	appdoc "github.com/sentzunhat/hawp/librarian/src/internal/application/work/doc"
	appwork "github.com/sentzunhat/hawp/librarian/src/internal/application/work/intake"
	domainwork "github.com/sentzunhat/hawp/librarian/src/internal/domain/work"
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

func toolWorkDoc(args json.RawMessage, repoRoot string) rpcResponse {
	var a struct {
		Type       string `json:"type"`
		Title      string `json:"title"`
		WorkItemID string `json:"work_item_id"`
	}
	if err := json.Unmarshal(args, &a); err != nil {
		return toolErr("invalid args: " + err.Error())
	}
	if a.Type == "" || a.Title == "" {
		return toolErr("type and title are required")
	}

	workDir := filepath.Join(repoRoot, ".hawp", "work")
	result, err := appdoc.CreateWorkDoc(a.Type, a.Title, workDir, a.WorkItemID)
	if err != nil {
		return toolErr("work doc failed: " + err.Error())
	}

	return text(fmt.Sprintf(
		"Created %s document %s\nFile: %s\n\nOpen the file and fill in the content.",
		result.DocType, result.UUID, result.FilePath,
	))
}

func toolWorkReshape(args json.RawMessage, repoRoot string) rpcResponse {
	_ = repoRoot // stateless: no repo access needed for shaping
	var a struct {
		Input   string `json:"input"`
		Context string `json:"context"`
		Backend string `json:"backend"`
		Model   string `json:"model"`
		URL     string `json:"url"`
	}
	if err := json.Unmarshal(args, &a); err != nil {
		return toolErr("invalid args: " + err.Error())
	}
	if strings.TrimSpace(a.Input) == "" {
		return toolErr("input is required")
	}
	if a.Backend == "" {
		a.Backend = "ollama"
	}
	if a.Backend != "ollama" && a.Backend != "onnx" {
		return toolErr(`backend must be "ollama" or "onnx"`)
	}

	shaper, cleanup, err := newMCPReshapeShaper(a.Backend, a.URL, a.Model)
	if err != nil {
		return toolErr("shaper unavailable: " + err.Error())
	}
	defer cleanup()

	draft, err := appwork.DraftIntake(context.Background(), appwork.DraftRequest{
		Input:   a.Input,
		Context: a.Context,
	}, shaper)
	if err != nil {
		return toolErr("reshape failed: " + err.Error())
	}

	return jsonResult(WorkReshapeResponse{
		Mission:     draft.Mission,
		Constraints: draft.Constraints,
		Output:      draft.Output,
		Checkpoint:  draft.Checkpoint,
	})
}

func toolWorkIntake(args json.RawMessage, repoRoot string) rpcResponse {
	var a struct {
		Input     string `json:"input"`
		Backend   string `json:"backend"`
		Model     string `json:"model"`
		URL       string `json:"url"`
		Limit     *int   `json:"limit"`
		MaxTokens *int   `json:"max_tokens"`
	}
	if err := json.Unmarshal(args, &a); err != nil {
		return toolErr("invalid args: " + err.Error())
	}
	if strings.TrimSpace(a.Input) == "" {
		return toolErr("input is required")
	}
	if a.Backend == "" {
		a.Backend = "ollama"
	}
	if a.Backend != "ollama" && a.Backend != "onnx" {
		return toolErr(`backend must be "ollama" or "onnx"`)
	}
	limit := 10
	if a.Limit != nil {
		limit = *a.Limit
	}
	if limit <= 0 {
		return toolErr("limit must be positive")
	} else if limit > 500 {
		return toolErr("limit must be 500 or less")
	}
	maxTokens := 2000
	if a.MaxTokens != nil {
		maxTokens = *a.MaxTokens
	}
	if maxTokens <= 0 {
		return toolErr("max_tokens must be positive")
	} else if maxTokens > 32000 {
		return toolErr("max_tokens must be 32000 or less")
	}

	response, err := runWorkIntake(context.Background(), repoRoot, a.Input, limit, maxTokens, func() (appwork.RequestShaper, func(), error) {
		return newMCPReshapeShaper(a.Backend, a.URL, a.Model)
	})
	if err != nil {
		return toolErr("intake failed: " + err.Error())
	}
	return jsonResult(response)
}

func runWorkIntake(ctx context.Context, repoRoot, input string, limit, maxTokens int, shaperFactory func() (appwork.RequestShaper, func(), error)) (WorkIntakeResponse, error) {
	contextBlock, err := searchContextResponse(input, limit, maxTokens, repoRoot)
	if err != nil {
		var notFound appsearch.IndexNotFoundError
		if errors.As(err, &notFound) {
			return WorkIntakeResponse{
				State: "blocked_missing_index",
				Retrieval: WorkIntakeRetrieval{
					Query:  input,
					Budget: maxTokens,
				},
				Questions: []string{"Run `hawp search index` and `hawp search embed`, then retry intake."},
				Warnings:  []string{notFound.Error()},
			}, nil
		}
		return WorkIntakeResponse{}, fmt.Errorf("search context: %w", err)
	}

	if contextBlock.ChunksUsed == 0 {
		return WorkIntakeResponse{
			State: "needs_user_input",
			Retrieval: WorkIntakeRetrieval{
				Query:  contextBlock.Query,
				Budget: maxTokens,
			},
			Questions: []string{"No indexed HAWP context matched this request. Ask the user for the missing background or index the relevant documents, then retry intake."},
			Warnings:  []string{"No search results were found; no draft was produced."},
		}, nil
	}

	shaper, cleanup, err := shaperFactory()
	if err != nil {
		return WorkIntakeResponse{
			State: "blocked_reshape_failed",
			Retrieval: WorkIntakeRetrieval{
				Query:         contextBlock.Query,
				ChunksUsed:    contextBlock.ChunksUsed,
				ChunksDropped: contextBlock.ChunksDropped,
				Budget:        contextBlock.Budget,
			},
			Questions: []string{"Review the retrieved context and retry intake with a more focused request or a different local model."},
			Warnings:  []string{"shaper unavailable: " + err.Error()},
			TokenAccounting: TokenAccounting{
				ContextTokens: contextBlock.TokenCount,
			},
		}, nil
	}
	defer cleanup()

	draft, err := appwork.DraftIntake(ctx, appwork.DraftRequest{
		Input:   input,
		Context: contextBlock.Content,
	}, shaper)
	if err != nil {
		return WorkIntakeResponse{
			State: "blocked_reshape_failed",
			Retrieval: WorkIntakeRetrieval{
				Query:         contextBlock.Query,
				ChunksUsed:    contextBlock.ChunksUsed,
				ChunksDropped: contextBlock.ChunksDropped,
				Budget:        contextBlock.Budget,
			},
			Questions: []string{"Review the retrieved context and retry intake with a more focused request or a different local model."},
			Warnings:  []string{"reshape failed: " + err.Error()},
			TokenAccounting: TokenAccounting{
				ContextTokens: contextBlock.TokenCount,
			},
		}, nil
	}

	response := WorkIntakeResponse{
		State: "ready_for_work_new",
		Draft: &WorkIntakeDraft{
			Input:       draft.Input,
			Context:     draft.Context,
			Mission:     draft.Mission,
			Constraints: draft.Constraints,
			Output:      draft.Output,
			Checkpoint:  draft.Checkpoint,
		},
		Retrieval: WorkIntakeRetrieval{
			Query:         contextBlock.Query,
			ChunksUsed:    contextBlock.ChunksUsed,
			ChunksDropped: contextBlock.ChunksDropped,
			Budget:        contextBlock.Budget,
		},
		TokenAccounting: TokenAccounting{
			ContextTokens: contextBlock.TokenCount,
			ShapedTokens:  estimateDraftTokens(draft),
			SavingsPct:    savingsPct(contextBlock.TokenCount, estimateDraftTokens(draft)),
		},
	}
	return response, nil
}

func estimateDraftTokens(d domainwork.Draft) int {
	chars := len(d.Mission) + len(d.Constraints) + len(d.Output) + len(d.Checkpoint)
	if chars == 0 {
		return 0
	}
	return (chars + 3) / 4
}

func savingsPct(contextTokens, shapedTokens int) int {
	if contextTokens <= 0 {
		return 0
	}
	saved := contextTokens - shapedTokens
	if saved <= 0 {
		return 0
	}
	return (saved * 100) / contextTokens
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
