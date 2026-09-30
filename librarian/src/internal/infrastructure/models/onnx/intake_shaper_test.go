//go:build ORT

package onnx

import (
	"context"
	"errors"
	"strings"
	"testing"

	appintake "github.com/sentzunhat/hawp/librarian/src/internal/application/work/intake"
)

// fakeReshaper satisfies llmReshaper without a live ONNX runtime.
type fakeReshaper struct {
	reshapeFunc func(userContent string) (string, error)
}

func (f *fakeReshaper) IntakeReshape(_ context.Context, _, userContent string, _ int) (string, error) {
	return f.reshapeFunc(userContent)
}

func onnxShaperWith(fn func(string) (string, error)) *ONNXIntakeShaper {
	return &ONNXIntakeShaper{client: &fakeReshaper{reshapeFunc: fn}, maxTokens: 256}
}

const onnxValidJSON = `{"mission":"Fix the login bug.","constraints":"Do not touch unrelated code.","output":"A passing test suite and patched handler."}`

func TestONNXShape_HappyPath(t *testing.T) {
	shaper := onnxShaperWith(func(_ string) (string, error) { return onnxValidJSON, nil })
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

func TestONNXShape_JSONEmbeddedInProse(t *testing.T) {
	prose := "Sure, here is the JSON:\n" + onnxValidJSON + "\nHope that helps!"
	shaper := onnxShaperWith(func(_ string) (string, error) { return prose, nil })
	got, err := shaper.Shape(context.Background(), appintake.DraftRequest{Input: "Fix it."})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Mission == "" {
		t.Error("mission should not be blank")
	}
}

func TestONNXShape_ContextAbsentOmitsBlock(t *testing.T) {
	var capturedPrompt string
	shaper := onnxShaperWith(func(p string) (string, error) {
		capturedPrompt = p
		return onnxValidJSON, nil
	})
	_, err := shaper.Shape(context.Background(), appintake.DraftRequest{Input: "Do something.", Context: ""})
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(capturedPrompt, "Supplied context:") {
		t.Errorf("prompt should not contain context block when Context is blank, got:\n%s", capturedPrompt)
	}
}

func TestONNXShape_ContextPresentIncludesBlock(t *testing.T) {
	var capturedPrompt string
	shaper := onnxShaperWith(func(p string) (string, error) {
		capturedPrompt = p
		return onnxValidJSON, nil
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

func TestONNXShape_BlankRequiredField_ReturnsError(t *testing.T) {
	shaper := onnxShaperWith(func(_ string) (string, error) {
		return `{"mission":"","constraints":"none","output":"something"}`, nil
	})
	_, err := shaper.Shape(context.Background(), appintake.DraftRequest{Input: "Do it."})
	if err == nil {
		t.Fatal("expected error for blank mission")
	}
}

func TestONNXShape_NoJSONInResponse_ReturnsError(t *testing.T) {
	shaper := onnxShaperWith(func(_ string) (string, error) {
		return "I cannot help with that request.", nil
	})
	_, err := shaper.Shape(context.Background(), appintake.DraftRequest{Input: "Do it."})
	if err == nil {
		t.Fatal("expected error when no JSON in response")
	}
}

func TestONNXShape_LLMError_Propagates(t *testing.T) {
	want := errors.New("inference failure")
	shaper := onnxShaperWith(func(_ string) (string, error) { return "", want })
	_, err := shaper.Shape(context.Background(), appintake.DraftRequest{Input: "Do it."})
	if !errors.Is(err, want) {
		t.Fatalf("expected wrapped inference failure, got: %v", err)
	}
}

func TestONNXShape_OptionalCheckpointAbsent(t *testing.T) {
	shaper := onnxShaperWith(func(_ string) (string, error) {
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

func TestONNXShape_OptionalCheckpointPresent(t *testing.T) {
	shaper := onnxShaperWith(func(_ string) (string, error) {
		return `{"mission":"Deploy the service.","constraints":"no downtime","output":"Service running in prod.","checkpoint":"Staging deploy verified."}`, nil
	})
	got, err := shaper.Shape(context.Background(), appintake.DraftRequest{Input: "Deploy it."})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Checkpoint != "Staging deploy verified." {
		t.Errorf("checkpoint: got %q", got.Checkpoint)
	}
}
