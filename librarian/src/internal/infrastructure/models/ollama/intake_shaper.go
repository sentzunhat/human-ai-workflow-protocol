package ollama

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	appintake "github.com/sentzunhat/hawp/librarian/src/internal/application/work/intake"
)

// llmReshaper is the subset of OllamaLLMClient used by OllamaIntakeShaper,
// extracted as an interface so tests can inject a fake without a live server.
type llmReshaper interface {
	Reshape(ctx context.Context, prompt string, maxTokens int) (string, error)
}

// OllamaIntakeShaper implements application/work/intake.RequestShaper using an
// Ollama LLM. It sends a HAWP-specific structured-extraction prompt and parses
// the JSON response into a DraftProposal.
type OllamaIntakeShaper struct {
	client    llmReshaper
	maxTokens int
}

// NewOllamaIntakeShaper wraps an existing OllamaLLMClient.
// If maxTokens <= 0, 512 is used.
func NewOllamaIntakeShaper(client *OllamaLLMClient, maxTokens int) *OllamaIntakeShaper {
	if maxTokens <= 0 {
		maxTokens = 512
	}
	return &OllamaIntakeShaper{client: client, maxTokens: maxTokens}
}

// Shape satisfies appintake.RequestShaper.
func (s *OllamaIntakeShaper) Shape(ctx context.Context, req appintake.DraftRequest) (appintake.DraftProposal, error) {
	prompt := buildIntakePrompt(req)

	raw, err := s.client.Reshape(ctx, prompt, s.maxTokens)
	if err != nil {
		return appintake.DraftProposal{}, fmt.Errorf("ollama reshape: %w", err)
	}

	proposal, err := parseProposal(raw)
	if err != nil {
		return appintake.DraftProposal{}, err
	}
	return proposal, nil
}

// buildIntakePrompt assembles the HAWP intake extraction prompt.
func buildIntakePrompt(req appintake.DraftRequest) string {
	var b strings.Builder
	b.WriteString("You are a HAWP work intake assistant. Given a user request, extract the following fields and respond in JSON only.\n\n")
	b.WriteString("Fields:\n")
	b.WriteString("- mission: What needs to be accomplished (the goal, 1-3 sentences)\n")
	b.WriteString("- constraints: What must not be done, hard limits, invariants (1-3 sentences; use \"none identified\" if truly absent)\n")
	b.WriteString("- output: What the deliverable looks like when complete (1-2 sentences)\n")
	b.WriteString("- checkpoint: An optional progress marker or interim artifact (omit or leave empty if not applicable)\n\n")
	b.WriteString("Rules:\n")
	b.WriteString("- Preserve the user's intent precisely — do not invent scope.\n")
	b.WriteString("- Label gaps as \"unknown\" rather than guessing.\n")
	b.WriteString("- Return valid JSON only, no surrounding text.\n\n")
	b.WriteString("User request:\n")
	b.WriteString(req.Input)
	b.WriteString("\n")
	if strings.TrimSpace(req.Context) != "" {
		b.WriteString("Supplied context:\n")
		b.WriteString(req.Context)
		b.WriteString("\n")
	}
	b.WriteString("JSON response:\n")
	return b.String()
}

// parseProposal extracts the JSON block from raw LLM output and unmarshals it.
func parseProposal(raw string) (appintake.DraftProposal, error) {
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
