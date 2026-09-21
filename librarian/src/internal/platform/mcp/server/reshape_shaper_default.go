//go:build !ORT

package mcp

import (
	"fmt"

	appintake "github.com/sentzunhat/hawp/librarian/src/internal/application/work/intake"
	ollamamodel "github.com/sentzunhat/hawp/librarian/src/internal/infrastructure/models/ollama"
)

// newMCPReshapeShaper returns a RequestShaper for the given backend plus a
// cleanup func. Without -tags ORT only "ollama" is supported; requesting
// "onnx" returns an error.
func newMCPReshapeShaper(backend, url, model string) (appintake.RequestShaper, func(), error) {
	if backend == "onnx" {
		return nil, nil, fmt.Errorf("ONNX reshape requires a -tags ORT build")
	}
	if url == "" {
		url = "http://localhost:11434"
	}
	client, err := ollamamodel.NewOllamaLLMClient(url, model)
	if err != nil {
		return nil, nil, fmt.Errorf("ollama unavailable: %w", err)
	}
	shaper := ollamamodel.NewOllamaIntakeShaper(client, 512)
	return shaper, func() { client.Close() }, nil //nolint:errcheck
}
