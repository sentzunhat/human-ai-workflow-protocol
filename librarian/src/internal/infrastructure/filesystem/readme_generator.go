package filesystem

import (
	"fmt"
	"os"
	"path/filepath"
)

// GenerateREADMEs creates documentation for key HAWP folders.
// Safe to call multiple times — only creates if not present.
func GenerateREADMEs(home, projectRoot string) error {
	hawpHome := ResolveHawpHome(home)
	hawpProj := ResolveHawpProject(projectRoot)

	// Home folder READMEs
	homeReadmes := map[string]string{
		filepath.Join(hawpHome.Root, "README.md"):             homeRootReadme(),
		filepath.Join(hawpHome.Models, "README.md"):           modelsRootReadme(),
		filepath.Join(hawpHome.ModelsEmbedding, "README.md"):  embeddingModelsReadme(),
		filepath.Join(hawpHome.ModelsLLM, "README.md"):        llmModelsReadme(),
		filepath.Join(hawpHome.Config, "README.md"):           configReadme(),
	}

	for path, content := range homeReadmes {
		if err := createReadmeIfNotExists(path, content); err != nil {
			return err
		}
	}

	// Project folder READMEs
	projectReadmes := map[string]string{
		filepath.Join(hawpProj.Root, "README.md"): projectRootReadme(),
		filepath.Join(hawpProj.Work, "README.md"): workReadme(),
		filepath.Join(hawpProj.Kit, "README.md"):  kitReadme(),
	}

	for path, content := range projectReadmes {
		if err := createReadmeIfNotExists(path, content); err != nil {
			return err
		}
	}

	return nil
}

// createReadmeIfNotExists creates a README file only if it doesn't already exist.
func createReadmeIfNotExists(path, content string) error {
	// Check if file exists
	if _, err := os.Stat(path); err == nil {
		return nil // File already exists, don't overwrite
	}

	// Create directory if needed
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// Write file
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write %s: %w", path, err)
	}

	return nil
}

// ============================================================================
// README Templates
// ============================================================================

func homeRootReadme() string {
	return `# HAWP Home (~/.hawp/)

Shared configuration and models for HAWP across all projects on this machine.

## Folders

- models/ - Embedding and LLM models (downloaded once, shared)
- config/ - Global default context configuration
- runtime/ - ONNX Runtime library (auto-downloaded)
- cache/ - Download staging area

## First-Time Setup

No manual setup needed! Models are auto-downloaded on first use.

## Configuration

Create ~/.hawp/config/context.json to set global defaults.
Project-specific configs override these defaults.
`
}

func modelsRootReadme() string {
	return `# Models (~/.hawp/models/)

Embedding and LLM models, auto-downloaded on first use.

## Folders

- embedding/ - Embedding models (BGE, MiniLM, etc.)
  Used to convert text to vectors
  Identify key concepts in context packing
  ~768 or 384-dimensional vectors

- llm/ - Large Language Models
  Used to reshape and improve context
  TinyLlama (fast), Mistral (higher quality), etc.
  Run locally via ONNX or use API backends

## About Models

Models are only downloaded when actually used. Default is ONNX (local execution):
- No internet required after download
- No API costs
- Private (never sends data to external servers)

To use different models, update your config or set environment variables.
`
}

func embeddingModelsReadme() string {
	return `# Embedding Models (~/.hawp/models/embedding/)

Models that convert text to numerical vectors for semantic search.

## Supported Models

- bge-base-en-v1.5 (default) - 768-dim, best quality (95%+ MTEB)
- all-MiniLM-L6-v2 - 384-dim, lighter, good fallback (90% MTEB)

## Usage

Models are auto-downloaded and cached here. No setup needed.

To use a specific model, set HAWP_EMBEDDINGS_MODEL environment variable.

## For v0.0.3+

Additional backends will be available:
- Ollama - Local embedding server
- OpenAI - text-embedding-3-small/large
- Anthropic - Coming when API available
`
}

func llmModelsReadme() string {
	return `# LLM Models (~/.hawp/models/llm/)

Large Language Models that reshape and improve context for better prompt injection.

## Supported Models

- TinyLlama-1.1B (default) - Smallest, fastest, runs anywhere
- Mistral-7B-v0.1 (optional) - Better quality, needs more resources

## Usage

Models are auto-downloaded on first use. No setup needed.

To use a specific model, set HAWP_LLM_MODEL environment variable.

## For v0.0.3+

Additional backends available:
- Ollama - Run any Ollama model locally
- OpenAI - gpt-3.5-turbo or gpt-4-turbo
- Anthropic - claude-3-sonnet or claude-3-opus
`
}

func configReadme() string {
	return `# Global Config (~/.hawp/config/)

Default configuration shared across all HAWP projects.

## Setup

Create ~/.hawp/config/context.json with your defaults.

## Priority

Config priority (highest to lowest):
1. CLI flags
2. Project config (.hawp/.data/config/context.json)
3. Home config (~/.hawp/config/context.json)
4. Built-in defaults: ONNX + BGE

## Environment Variables

Override via env vars (HAWP_*):
- HAWP_EMBEDDINGS_BACKEND=openai
- HAWP_LLM_BACKEND=anthropic
- HAWP_OPENAI_API_KEY=sk-...
- HAWP_ANTHROPIC_API_KEY=sk-ant-...
`
}

func projectRootReadme() string {
	return `# Project HAWP (.hawp/)

Project-specific configuration, work tracking, and patterns.

## Folders

- work/ - Task tracking and planning (committed to git)
- kit/ - Patterns and standards (committed to git)
- .data/ - Auto-created runtime data (NOT in git)

## First-Time Setup

After cloning or installing, build the search index once:

  hawp search index

This creates .hawp/.data/db/index.sqlite containing:
- Lexical search index (FTS5)
- Vector embeddings (semantic search)
- Project-specific caches

Takes ~2-3 seconds. Only needed once after clone/install.

## .data/ Folder (Auto-Created)

NOT committed to git. Contains:
- .data/db/index.sqlite - Search index
- .data/config/context.json - Project-specific config
- .data/db/embeddings/ - Embedding cache (optional)

Regenerated on demand if the codebase changes.

## Using Context Packing

  hawp search "query" --context              # LLM-ready markdown
  hawp search "query" --context --format json # JSON output
`
}

func workReadme() string {
	return `# Work Tracking (.hawp/work/)

Task backlog, plans, evidence, and status reports.

See BACKLOG.md for the active work index.

Files are committed to git for team collaboration.
`
}

func kitReadme() string {
	return `# HAWP Kit (.hawp/kit/)

Canonical patterns, standards, and usage guides.

These are the source of truth for how HAWP is used in this project.

Files are committed to git.
`
}
