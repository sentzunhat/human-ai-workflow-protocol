//go:build ORT

package benchmark

import (
	"fmt"

	appintake "github.com/sentzunhat/hawp/librarian/src/internal/application/work/intake"
	ollamamodel "github.com/sentzunhat/hawp/librarian/src/internal/infrastructure/models/ollama"
	onnxmodel "github.com/sentzunhat/hawp/librarian/src/internal/infrastructure/models/onnx"
)

// newReshapeShaper returns a RequestShaper for the given backend, plus a
// cleanup function to call when done. Supports "onnx" (requires -tags ORT)
// and "ollama".
func newReshapeShaper(backend, url, model string) (appintake.RequestShaper, func(), error) {
	switch backend {
	case "onnx":
		if model == "" {
			model = onnxmodel.DefaultLLMModel
		}
		client, err := onnxmodel.NewONNXLLMClient(model)
		if err != nil {
			return nil, nil, fmt.Errorf("ONNX LLM: %w", err)
		}
		shaper := onnxmodel.NewONNXIntakeShaper(client, 512)
		return shaper, func() { client.Close() }, nil //nolint:errcheck
	default: // "ollama"
		if url == "" {
			url = "http://localhost:11434"
		}
		client, err := ollamamodel.NewOllamaLLMClient(url, model)
		if err != nil {
			return nil, nil, fmt.Errorf("Ollama LLM: %w", err)
		}
		shaper := ollamamodel.NewOllamaIntakeShaper(client, 256)
		return shaper, func() { client.Close() }, nil //nolint:errcheck
	}
}
