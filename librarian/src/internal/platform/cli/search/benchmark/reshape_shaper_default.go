//go:build !ORT

package benchmark

import (
	"fmt"

	appintake "github.com/sentzunhat/hawp/librarian/src/internal/application/work/intake"
	ollamamodel "github.com/sentzunhat/hawp/librarian/src/internal/infrastructure/models/ollama"
)

// newReshapeShaper returns a RequestShaper for the given backend, plus a
// cleanup function to call when done. Without -tags ORT, only "ollama" is
// supported; requesting "onnx" returns an immediate error.
func newReshapeShaper(backend, url, model string) (appintake.RequestShaper, func(), error) {
	if backend == "onnx" {
		return nil, nil, fmt.Errorf("ONNX reshape requires a -tags ORT build (rebuild with CGO_LDFLAGS pointing at the native libs)")
	}
	// Default: Ollama
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
