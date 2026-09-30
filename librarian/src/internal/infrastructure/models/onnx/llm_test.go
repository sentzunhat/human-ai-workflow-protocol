package onnx

import (
	"strings"
	"testing"

	domainllm "github.com/sentzunhat/hawp/librarian/src/internal/domain/providers/llm"
)

func TestONNXLLMUnsupportedModel(t *testing.T) {
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
	_, err := NewONNXLLMClient(DefaultLLMModel)
	if err == nil {
		t.Fatal("NewONNXLLMClient should fail without a -tags ORT build")
	}
	if !strings.Contains(err.Error(), "ORT") {
		t.Errorf("error = %v, want it to mention ORT (hugot's own disabled-backend message)", err)
	}
}

func TestONNXLLMClientInterface(t *testing.T) {
	var _ domainllm.LLMClient = (*ONNXLLMClient)(nil)
}

// TestSupportedLLMModels verifies SmolLM2-360M-Instruct is registered with the
// fields NewONNXLLMClient needs to download and run it.
func TestSupportedLLMModels(t *testing.T) {
	info, ok := SupportedLLMModels["SmolLM2-360M-Instruct"]
	if !ok {
		t.Fatal(`SupportedLLMModels should contain "SmolLM2-360M-Instruct"`)
	}
	if info.HFRepo == "" {
		t.Error("SmolLM2-360M-Instruct should have a non-empty HFRepo")
	}
	if info.ONNXFile == "" {
		t.Error("SmolLM2-360M-Instruct should have a non-empty ONNXFile")
	}
}

func TestDefaultLLMModelIsSupported(t *testing.T) {
	if DefaultLLMModel == "" {
		t.Fatal("DefaultLLMModel should not be empty — an empty default means NewONNXLLMClient(\"\") always fails")
	}
	if _, ok := SupportedLLMModels[DefaultLLMModel]; !ok {
		t.Errorf("DefaultLLMModel %q must be a key in SupportedLLMModels", DefaultLLMModel)
	}
}

func TestGetLLMModelPath(t *testing.T) {
	home := "/home/user"
	model := "SmolLM2-360M-Instruct"
	path := GetLLMModelPath(home, model)
	if path != "/home/user/.hawp/models/llm/SmolLM2-360M-Instruct" {
		t.Errorf("model path incorrect: %s", path)
	}
}

func TestGetLLMModelPathArbitrary(t *testing.T) {
	home := "/home/user"
	model := "mistral-7b"
	path := GetLLMModelPath(home, model)
	if path != "/home/user/.hawp/models/llm/mistral-7b" {
		t.Errorf("model path incorrect: %s", path)
	}
}

func TestReshapingPromptFormat(t *testing.T) {
	if !strings.Contains(domainllm.ReshapingPrompt, "{maxTokens}") {
		t.Error("ReshapingPrompt should contain {maxTokens} placeholder")
	}
	if !strings.Contains(domainllm.ReshapingPrompt, "clarity") {
		t.Error("ReshapingPrompt should mention clarity")
	}
	if !strings.Contains(domainllm.ReshapingPrompt, "redundancy") {
		t.Error("ReshapingPrompt should mention redundancy")
	}
}
