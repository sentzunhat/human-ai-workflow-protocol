//go:build ORT

package onnx

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	appintake "github.com/sentzunhat/hawp/librarian/src/internal/application/work/intake"
)

// llmReshaper is the subset of ONNXLLMClient used by ONNXIntakeShaper,
// extracted as an interface so tests can inject a fake without native libs.
type llmReshaper interface {
	IntakeReshape(ctx context.Context, systemPrompt, userContent string, maxTokens int) (string, error)
}

// ONNXIntakeShaper implements application/work/intake.RequestShaper using the
// ONNX LLM backend. It sends the same HAWP intake extraction prompt as
// OllamaIntakeShaper and parses the JSON response into a DraftProposal.
type ONNXIntakeShaper struct {
	client    llmReshaper
	maxTokens int
}

// NewONNXIntakeShaper wraps an existing ONNXLLMClient.
// If maxTokens <= 0, 256 is used.
func NewONNXIntakeShaper(client *ONNXLLMClient, maxTokens int) *ONNXIntakeShaper {
	if maxTokens <= 0 {
		maxTokens = 256
	}
	return &ONNXIntakeShaper{client: client, maxTokens: maxTokens}
}

// intakeSystemPrompt is placed in the ChatML system turn so the model actually
// follows it. Previously the full extraction prompt was in the user turn while
// ReshapingPrompt held the system turn — SmolLM2-360M ignored the user-turn
// instruction entirely at 770 chars.
const intakeSystemPrompt = `You extract structured work intake fields from a user request. Respond ONLY with valid JSON containing exactly these fields: "mission", "constraints", "output", and optionally "checkpoint". No prose, no explanation, no markdown fences.`

// Shape satisfies appintake.RequestShaper.
func (s *ONNXIntakeShaper) Shape(ctx context.Context, req appintake.DraftRequest) (appintake.DraftProposal, error) {
	userContent := buildONNXUserContent(req)

	raw, err := s.client.IntakeReshape(ctx, intakeSystemPrompt, userContent, s.maxTokens)
	if err != nil {
		return appintake.DraftProposal{}, fmt.Errorf("onnx reshape: %w", err)
	}

	proposal, err := parseONNXProposal(raw)
	if err != nil {
		return appintake.DraftProposal{}, err
	}
	return proposal, nil
}

// buildONNXUserContent builds the user turn for intake extraction.
// Kept short so the system instruction in the system turn dominates.
func buildONNXUserContent(req appintake.DraftRequest) string {
	var b strings.Builder
	b.WriteString("Extract HAWP intake fields from this request:\n\nRequest: ")
	b.WriteString(req.Input)
	b.WriteString("\n")
	if strings.TrimSpace(req.Context) != "" {
		b.WriteString("Supplied context:\n")
		b.WriteString(req.Context)
		b.WriteString("\n")
	}
	b.WriteString("\nReturn JSON with: mission (one sentence goal), constraints (limits/requirements), output (deliverable), checkpoint (optional: success criteria).")
	return b.String()
}

// parseONNXProposal extracts the JSON block from raw LLM output and unmarshals it.
// Identical to the Ollama shaper's parseProposal.
func parseONNXProposal(raw string) (appintake.DraftProposal, error) {
	start := strings.Index(raw, "{")
	if start < 0 {
		return appintake.DraftProposal{}, fmt.Errorf("no JSON object found in LLM response")
	}
	end := strings.LastIndex(raw, "}")
	if end < start {
		return appintake.DraftProposal{}, fmt.Errorf("malformed JSON in LLM response: no closing brace")
	}
	jsonStr := raw[start : end+1]

	var fields struct {
		Mission     json.RawMessage `json:"mission"`
		Constraints json.RawMessage `json:"constraints"`
		Output      json.RawMessage `json:"output"`
		Checkpoint  json.RawMessage `json:"checkpoint"`
	}
	if err := json.Unmarshal([]byte(jsonStr), &fields); err != nil {
		return appintake.DraftProposal{}, fmt.Errorf("unmarshal LLM JSON: %w", err)
	}

	mission := jsonFieldText(fields.Mission)
	constraints := jsonFieldText(fields.Constraints)
	output := jsonFieldText(fields.Output)
	checkpoint := jsonFieldText(fields.Checkpoint)

	if strings.TrimSpace(mission) == "" {
		return appintake.DraftProposal{}, fmt.Errorf("LLM returned blank mission")
	}
	if strings.TrimSpace(constraints) == "" {
		return appintake.DraftProposal{}, fmt.Errorf("LLM returned blank constraints")
	}
	if strings.TrimSpace(output) == "" {
		return appintake.DraftProposal{}, fmt.Errorf("LLM returned blank output")
	}

	return appintake.DraftProposal{
		Mission:     mission,
		Constraints: constraints,
		Output:      output,
		Checkpoint:  checkpoint,
	}, nil
}

func jsonFieldText(raw json.RawMessage) string {
	if len(raw) == 0 || string(raw) == "null" {
		return ""
	}
	var s string
	if err := json.Unmarshal(raw, &s); err == nil {
		return strings.TrimSpace(s)
	}
	return strings.TrimSpace(string(raw))
}
