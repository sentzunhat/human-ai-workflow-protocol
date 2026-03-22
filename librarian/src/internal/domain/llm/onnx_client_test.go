package llm

import (
	"strings"
	"testing"
)

func TestUnsupportedModel(t *testing.T) {
	_, err := NewONNXLLMClient("nonexistent-model")
	if err == nil {
		t.Error("unsupported model should return error")
	}
}

// TestNewONNXLLMClientFailsFastWithoutDownloading locks in that
// NewONNXLLMClient creates the ORT session (and, without a `-tags ORT`
// build, fails there) BEFORE downloading anything — so a default `go test`
// run of this package never attempts a ~1GB+ network download. hugot's
// disabled-ORT-backend stub returns its error instantly (see
// hugot_ort_disabled.go: "to enable ORT, run `go build -tags ORT`"), which
// is what this test runs against (this package's tests are not built with
// `-tags ORT`). A `-tags ORT` build with the native libraries actually
// present would instead proceed to download — that path is exercised
// manually, not by this default test run (see v0.1.0_VISION.md).
func TestNewONNXLLMClientFailsFastWithoutDownloading(t *testing.T) {
	_, err := NewONNXLLMClient(DefaultModel)
	if err == nil {
		t.Fatal("NewONNXLLMClient should fail without a -tags ORT build")
	}
	if !strings.Contains(err.Error(), "ORT") {
		t.Errorf("error = %v, want it to mention ORT (hugot's own disabled-backend message)", err)
	}
}

func TestLLMClientInterface(t *testing.T) {
	// Test interface is properly defined
	var _ LLMClient = (*ONNXLLMClient)(nil)
}

func TestNewLLMClientFactory(t *testing.T) {
	// Test that factory function works
	_, err := NewLLMClient("onnx", "nonexistent-model")
	if err == nil {
		t.Error("unsupported model should return error")
	}

	// Test unsupported backend
	_, err = NewLLMClient("unsupported-backend", "some-model")
	if err == nil {
		t.Error("unsupported backend should return error")
	}
}

// TestSupportedModels verifies SmolLM2-360M-Instruct is registered with the
// fields NewONNXLLMClient needs to download and run it. Chosen over
// FLAN-T5-small after benchmarking (2026-07-26) — see BENCHMARKS_v003.md.
func TestSupportedModels(t *testing.T) {
	info, ok := SupportedModels["SmolLM2-360M-Instruct"]
	if !ok {
		t.Fatal(`SupportedModels should contain "SmolLM2-360M-Instruct"`)
	}
	if info.HFRepo == "" {
		t.Error("SmolLM2-360M-Instruct should have a non-empty HFRepo")
	}
	if info.ONNXFile == "" {
		t.Error("SmolLM2-360M-Instruct should have a non-empty ONNXFile")
	}
}

func TestDefaultModelIsSupported(t *testing.T) {
	if DefaultModel == "" {
		t.Fatal("DefaultModel should not be empty — an empty default means NewONNXLLMClient(\"\") always fails")
	}
	if _, ok := SupportedModels[DefaultModel]; !ok {
		t.Errorf("DefaultModel %q must be a key in SupportedModels", DefaultModel)
	}
}

func TestGetModelPath(t *testing.T) {
	home := "/home/user"
	model := "SmolLM2-360M-Instruct"

	path := GetModelPath(home, model)

	if path != "/home/user/.hawp/models/llm/SmolLM2-360M-Instruct" {
		t.Errorf("model path incorrect: %s", path)
	}
}

// TestReshapingPromptFormat verifies the reshaping prompt template is correct.
func TestReshapingPromptFormat(t *testing.T) {
	if !strings.Contains(ReshapingPrompt, "{maxTokens}") {
		t.Error("ReshapingPrompt should contain {maxTokens} placeholder")
	}

	if !strings.Contains(ReshapingPrompt, "clarity") {
		t.Error("ReshapingPrompt should mention clarity")
	}

	if !strings.Contains(ReshapingPrompt, "redundancy") {
		t.Error("ReshapingPrompt should mention redundancy")
	}
}

// TestGetModelPathLLM tests the model path resolution for LLM models.
func TestGetModelPathLLM(t *testing.T) {
	home := "/home/user"
	model := "mistral-7b"

	path := GetModelPath(home, model)

	if path != "/home/user/.hawp/models/llm/mistral-7b" {
		t.Errorf("model path incorrect: %s", path)
	}
}
