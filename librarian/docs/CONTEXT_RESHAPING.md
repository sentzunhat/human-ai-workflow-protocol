# Context Reshaping with LLM Backend

The `--context --llm-reshape` flags enable semantic context improvement: search finds relevant results, embeddings identify key concepts, and an LLM restructures the content for clarity.

## Quick Start (Mac Max M1)

### Prerequisites
1. **Ollama** (required for LLM reshaping)
   ```bash
   # Install: https://ollama.ai/
   # Then download a model:
   ollama pull mistral
   
   # Start Ollama server (in a separate terminal):
   ollama serve
   ```

2. **HAWP** (v0.0.3+)
   ```bash
   hawp --version  # Should show v0.0.3 or higher
   ```

### Default Usage (ONNX + Ollama)

```bash
# Search + reshape in one command
hawp search "kubernetes deployment strategies" --context --llm-reshape

# Output includes:
#   1. Ranked search results
#   2. Key concepts extracted via ONNX embeddings
#   3. Reshaped markdown optimized for clarity
#   4. Ready to paste into Claude or your editor
```

**Default backends:**
- **Embeddings:** ONNX (bge-base-en-v1.5, 768-dim, downloaded on first use)
- **LLM:** Ollama (mistral, via localhost:11434)

**No configuration needed** — just run the command and embeddings are downloaded automatically on first use.

---

## Configuration

### Location
```
~/.hawp/config/context.json
```

### Full-Ollama Mode (Recommended for Speed)

```json
{
  "embeddings": {
    "backend": "ollama",
    "model": "nomic-embed-text"
  },
  "llm": {
    "backend": "ollama",
    "model": "mistral"
  }
}
```

**First run:**
```bash
ollama pull nomic-embed-text
ollama pull mistral
ollama serve  # Start in background
```

### ONNX + Ollama (Default, Recommended for Offline)

```json
{
  "embeddings": {
    "backend": "onnx",
    "model": "bge-base-en-v1.5"
  },
  "llm": {
    "backend": "ollama",
    "model": "mistral"
  }
}
```

**Benefits:** ONNX runs completely offline (no network), Ollama can use any model.

### Faster/Smaller Embeddings

```json
{
  "embeddings": {
    "backend": "onnx",
    "model": "all-MiniLM-L6-v2"  // 384-dim, smaller/faster than the default
  },
  "llm": {
    "backend": "ollama",
    "model": "neural-chat:latest"  // Faster, smaller
  }
}
```

Supported ONNX embedding models:
- `bge-base-en-v1.5` (768-dim, downloaded on first use, higher quality) ← default
- `all-MiniLM-L6-v2` (384-dim, smaller/faster, lower quality)
- `bge-large-en-v1.5` (1024-dim, largest download, highest quality, slowest)

Supported Ollama models:
- `mistral` (7B, recommended)
- `neural-chat` (7B, faster)
- `llama2` (7B/13B/70B)
- Any other GGUF model you've pulled

---

## Backend Comparison

| Aspect | ONNX Embeddings | Ollama Embeddings | Ollama LLM |
|--------|---|---|---|
| **Speed** | <50ms | 10-50ms | 500-2000ms |
| **Quality** | Good (384/768-dim vectors) | High (768/1024-dim) | High (7B+ models) |
| **Privacy** | ✅ 100% local | ✅ 100% local | ✅ 100% local |
| **Setup** | Auto-download (~50MB) | Manual: `ollama pull` | Manual: `ollama pull` |
| **Deps** | None (included) | Requires Ollama | Requires Ollama |
| **CPU Only** | ✅ Yes | ✅ Yes | ✅ Yes (slower) |
| **GPU Support** | Not yet | Via Ollama | Via Ollama |
| **Best For** | Fast, offline search | Quality + flexibility | Context improvement |

---

## Advanced: Environment Variables

Override configuration without editing `~/.hawp/config/context.json`:

```bash
# Use Ollama embeddings instead of ONNX
export HAWP_EMBEDDINGS_BACKEND=ollama
export HAWP_EMBEDDINGS_MODEL=nomic-embed-text

# Use a custom Ollama URL (default: localhost:11434)
export HAWP_OLLAMA_URL=http://192.168.1.100:11434

# Control LLM behavior
export HAWP_LLM_MAX_TOKENS=512
export HAWP_LLM_TEMPERATURE=0.7

# Then run:
hawp search "query" --context --llm-reshape
```

---

## Troubleshooting

### "Ollama not found" or Connection Refused

**Problem:** `--llm-reshape` requires Ollama running.

**Solution:**
```bash
# 1. Install Ollama: https://ollama.ai/
# 2. Start the server in a separate terminal:
ollama serve

# 3. Download a model:
ollama pull mistral

# 4. Try again:
hawp search "query" --context --llm-reshape
```

**Why:** LLM reshaping needs a running Ollama server on `localhost:11434`.

### Model Not Found

**Problem:** `error: model 'mistral' not found`

**Solution:**
```bash
# List available models:
ollama list

# Download missing model:
ollama pull mistral

# Verify it's there:
ollama list  # Should show 'mistral'
```

### ONNX Model Download Slow

**Problem:** First run downloads ~50MB for embeddings model.

**Solution:**
```bash
# Pre-download before first use:
go run ./cmd/hawp search "dummy" --context 2>&1 | head -20
# (This triggers ONNX model download in the background)

# Or wait — it only happens once
```

### Slow Performance (CPU Bound)

**Problem:** Reshaping takes >10 seconds.

**Options:**
1. **Use smaller Ollama model:**
   ```json
   {
     "llm": {
       "backend": "ollama",
       "model": "neural-chat:latest"  // Faster than mistral
     }
   }
   ```

2. **Reduce token budget:**
   ```json
   {
     "llm": {
       "backend": "ollama",
       "max_tokens": 256  // Default 512
     }
   }
   ```

3. **Use GPU (if available):**
   - Ollama auto-detects GPU if installed
   - Install CUDA/Metal drivers for your system

---

## CLI Flags

```bash
# Basic search
hawp search "query"

# Search + collect context blocks
hawp search "query" --context

# Search + reshape context via LLM (requires Ollama running)
hawp search "query" --context --llm-reshape

# Output to file
hawp search "query" --context --llm-reshape > output.md

# JSON format (for scripting)
hawp search "query" --context --llm-reshape --json

# Limit results
hawp search "query" --context --llm-reshape --limit 3

# Verbose logging
hawp search "query" --context --llm-reshape -vv
```

---

## Go API Usage

```go
package main

import (
	"context"
	"fmt"
	"log"

	appContext "github.com/sentzunhat/hawp/librarian/go/internal/application/context"
)

func main() {
	ctx := context.Background()

	// Configure backends
	config := appContext.ReshapingConfig{
		EmbeddingsBackend: "onnx",
		EmbeddingsModel:   "all-MiniLM-L6-v2",
		LLMBackend:        "ollama",
		LLMModel:          "mistral",
		MaxTokens:         512,
		TopK:              5,
	}

	// Create reshaper (initializes embedder + LLM)
	reshaper, err := appContext.NewContextReshaper(config)
	if err != nil {
		log.Fatalf("Failed to create reshaper: %v", err)
	}
	defer reshaper.Close()

	// Get a context block from your search results
	block := &appContext.ContextBlock{
		Title:   "Search Results",
		Query:   "kubernetes",
		Results: []appContext.FormattedResult{
			// ... populate from search
		},
		Metadata: make(map[string]string),
	}

	// Reshape it
	improved, err := reshaper.Reshape(ctx, block)
	if err != nil {
		log.Fatalf("Failed to reshape: %v", err)
	}

	fmt.Printf("Reshaped content:\n%s\n", improved.Content)
	fmt.Printf("Key concepts: %v\n", improved.KeyConcepts)
	fmt.Printf("Pipeline: %s\n", improved.Pipeline)
}
```

---

## Performance Expectations

### Typical Query (5 results, ~2000 tokens)

| Backend | Embeddings | LLM | Total |
|---------|---|---|---|
| ONNX + Ollama | 40ms | 1200ms | ~1.3s |
| Ollama + Ollama | 150ms | 1200ms | ~1.4s |

**Notes:**
- LLM time dominated by model (bigger = slower)
- Use `neural-chat:latest` for 30% faster LLM
- ONNX embeddings much faster than HTTP round-trips
- First run includes model downloads

---

## What Gets Reshaped?

The reshaper improves context readability:

**Before:**
```
1. Kubernetes Deployment best practices (relevance: 0.95)
   Container orchestration patterns for scaling services...
   
2. Docker multi-stage builds (relevance: 0.87)
   Reducing image size through layer caching...
```

**After:**
```
# Deployment & Scaling Strategies

Key concepts: Kubernetes, Container orchestration, Scaling, Deployment

The provided content covers container orchestration patterns and best practices.
Focus areas:
- Deployment patterns for horizontal scaling
- Service management in Kubernetes
- Efficient container imaging through Docker multi-stage builds

Recommended approach: Use Kubernetes Deployments with resource limits, combine 
with multi-stage Docker builds for image optimization.
```

---

## Next Steps

- See [BACKENDS.md](./BACKENDS.md) for detailed backend architecture
- See [TROUBLESHOOTING.md](./TROUBLESHOOTING.md) for more issues
- See `.hawp/kit/` for HAWP methodology documentation
