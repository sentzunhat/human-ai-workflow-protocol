// Package bootstrap owns concrete service wiring for command entrypoints.
package bootstrap

import (
	appcontext "github.com/sentzunhat/hawp/librarian/src/internal/application/context"
	"github.com/sentzunhat/hawp/librarian/src/internal/application/index"
	appsearch "github.com/sentzunhat/hawp/librarian/src/internal/application/search"
	inframodels "github.com/sentzunhat/hawp/librarian/src/internal/infrastructure/models"
)

func NewEmbedService(dbPath string) *index.EmbedService {
	return index.NewEmbedServiceWithFactory(dbPath, inframodels.NewEmbedder)
}

func NewSearchService() appsearch.Service {
	return appsearch.NewServiceWithEmbedder(nil, nil, inframodels.NewEmbedder)
}

func NewContextReshaper(config appcontext.ReshapingConfig) (*appcontext.ContextReshaper, error) {
	return appcontext.NewContextReshaper(config, inframodels.NewEmbedderWithURL, inframodels.NewLLMClientWithURL)
}

func NewRAGPipeline(config appcontext.ReshapingConfig, repoRoot string) (*appcontext.DefaultRAGPipeline, error) {
	return appcontext.NewDefaultRAGPipeline(config, repoRoot, inframodels.NewEmbedderWithURL, inframodels.NewLLMClientWithURL)
}
