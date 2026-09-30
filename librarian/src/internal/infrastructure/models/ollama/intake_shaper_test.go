package ollama

import (
	"context"
	"errors"
	"strings"
	"testing"

	appintake "github.com/sentzunhat/hawp/librarian/src/internal/application/work/intake"
)

// fakeOllamaLLMClient satisfies llmReshaper without a live Ollama server.
type fakeOllamaLLMClient struct {
	reshapeFunc func(prompt string) (string, error)
}

func (f *fakeOllamaLLMClient) Reshape(_ context.Context, prompt string, _ int) (string, error) {
	return f.reshapeFunc(prompt)
}

func shaperWith(fn func(string) (string, error)) *OllamaIntakeShaper {
	return &OllamaIntakeShaper{client: &fakeOllamaLLMClient{reshapeFunc: fn}, maxTokens: 512}
}

const validJSON = `{"mission":"Fix the login bug.","constraints":"Do not touch unrelated code.","output":"A passing test suite and patched handler."}`

func TestShape_HappyPath(t *testing.T) {
	shaper := shaperWith(func(_ string) (string, error) { return validJSON, nil })
	got, err := shaper.Shape(context.Background(), appintake.DraftRequest{Input: "Fix the login bug."})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Mission != "Fix the login bug." {
		t.Errorf("mission: got %q", got.Mission)
	}
	if got.Constraints != "Do not touch unrelated code." {
		t.Errorf("constraints: got %q", got.Constraints)
	}
	if got.Output != "A passing test suite and patched handler." {
		t.Errorf("output: got %q", got.Output)
	}
	if got.Checkpoint != "" {
		t.Errorf("checkpoint should be blank, got %q", got.Checkpoint)
	}
}

func TestShape_JSONEmbeddedInProse(t *testing.T) {
	prose := "Sure, here is the JSON:\n" + validJSON + "\nHope that helps!"
	shaper := shaperWith(func(_ string) (string, error) { return prose, nil })
	got, err := shaper.Shape(context.Background(), appintake.DraftRequest{Input: "Fix it."})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Mission == "" {
		t.Error("mission should not be blank")
	}
}

func TestShape_ContextAbsentOmitsBlock(t *testing.T) {
	var capturedPrompt string
	shaper := shaperWith(func(p string) (string, error) {
		capturedPrompt = p
		return validJSON, nil
	})
	_, err := shaper.Shape(context.Background(), appintake.DraftRequest{Input: "Do something.", Context: ""})
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(capturedPrompt, "Supplied context:") {
		t.Errorf("prompt should not contain context block when Context is blank, got:\n%s", capturedPrompt)
	}
}

func TestShape_ContextPresentIncludesBlock(t *testing.T) {
	var capturedPrompt string
	shaper := shaperWith(func(p string) (string, error) {
		capturedPrompt = p
		return validJSON, nil
	})
	_, err := shaper.Shape(context.Background(), appintake.DraftRequest{Input: "Do something.", Context: "Background info here."})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(capturedPrompt, "Supplied context:") {
		t.Errorf("prompt should contain context block, got:\n%s", capturedPrompt)
	}
	if !strings.Contains(capturedPrompt, "Background info here.") {
		t.Errorf("prompt should contain the context text, got:\n%s", capturedPrompt)
	}
}

func TestShape_BlankRequiredField_ReturnsError(t *testing.T) {
	shaper := shaperWith(func(_ string) (string, error) {
		return `{"mission":"","constraints":"none","output":"something"}`, nil
	})
	_, err := shaper.Shape(context.Background(), appintake.DraftRequest{Input: "Do it."})
	if err == nil {
		t.Fatal("expected error for blank mission")
	}
}

func TestShape_NoJSONInResponse_ReturnsError(t *testing.T) {
	shaper := shaperWith(func(_ string) (string, error) {
		return "I cannot help with that request.", nil
	})
	_, err := shaper.Shape(context.Background(), appintake.DraftRequest{Input: "Do it."})
	if err == nil {
		t.Fatal("expected error when no JSON in response")
	}
}

func TestShape_LLMError_Propagates(t *testing.T) {
	want := errors.New("network failure")
	shaper := shaperWith(func(_ string) (string, error) { return "", want })
	_, err := shaper.Shape(context.Background(), appintake.DraftRequest{Input: "Do it."})
	if !errors.Is(err, want) {
		t.Fatalf("expected wrapped network failure, got: %v", err)
	}
}

func TestShape_OptionalCheckpointAbsent(t *testing.T) {
	shaper := shaperWith(func(_ string) (string, error) {
		return `{"mission":"Ship it.","constraints":"none identified","output":"Binary released."}`, nil
	})
	got, err := shaper.Shape(context.Background(), appintake.DraftRequest{Input: "Ship it."})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Checkpoint != "" {
		t.Errorf("checkpoint should be blank when absent from JSON, got %q", got.Checkpoint)
	}
}

func TestShape_OptionalCheckpointObjectBecomesText(t *testing.T) {
	shaper := shaperWith(func(_ string) (string, error) {
		return `{"mission":"Ship it.","constraints":"none identified","output":"Binary released.","checkpoint":{"review":"before create"}}`, nil
	})
	got, err := shaper.Shape(context.Background(), appintake.DraftRequest{Input: "Ship it."})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(got.Checkpoint, `"review"`) {
		t.Errorf("checkpoint should preserve object text, got %q", got.Checkpoint)
	}
}
