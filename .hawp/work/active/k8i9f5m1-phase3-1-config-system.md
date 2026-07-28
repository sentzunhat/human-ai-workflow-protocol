---
work-item: k8i9f5m1-p3-1
type: feature
title: "v0.0.3 Phase 3.1: Config System & Secure Key Storage (Foundation)"
status: plan-ready
owner: unassigned
created: 2026-07-24
updated: 2026-07-24
---

# Phase 3.1: Config System & Secure Key Storage

## Mission

Build the foundation config system for v0.0.3 that enables all downstream embedding and LLM backends to be configurable, secure, and composable.

This phase **must complete first** — it unblocks Phase 3.2 (embedding backends) and Phase 3.3 (LLM backends) to run in parallel.

---

## Context

**Phase 3 Overview:**
- Context reshaping via embeddings + LLM inference
- Both embeddings and LLM backends are fully configurable (user can mix and match)
- Primary: ONNX local (no API costs, fast, private)
- Optional: Ollama, OpenAI, Anthropic

**Why Config System First?**
- Embedding backends (3.2) need to know: which backend to use, which model, where to fetch from
- LLM backends (3.3) need to know: which backend to use, which model, API credentials
- Reshaping pipeline (3.4) needs to instantiate backends based on config
- Without this layer, 3.2 and 3.3 can't be unified or tested independently

---

## Design

### Config File Structure

**Location:** `~/.hawp/config/context.json` (user-creatable, gitignored)

**Format (JSON with encrypted secrets):**

```json
{
  "embeddings": {
    "backend": "onnx",
    "model": "bge-base-en-v1.5",
    "cache_dir": "~/.hawp/models"
  },
  "llm": {
    "backend": "onnx",
    "model": "none (onnx uses local model)",
    "temperature": 0.3,
    "max_tokens": 2000
  },
  "backends": {
    "ollama": {
      "url": "http://localhost:11434"
    },
    "openai": {
      "api_key": "[BASE64-ENCRYPTED]",
      "base_url": "https://api.openai.com/v1"
    },
    "anthropic": {
      "api_key": "[BASE64-ENCRYPTED]",
      "base_url": "https://api.anthropic.com"
    }
  },
  "security": {
    "encryption_method": "base64",
    "key_rotation_enabled": false
  }
}
```

### Environment File Support

**Location:** `.env` in repo root or `~/.hawp/.env`

**Format:**
```bash
# Embeddings
HAWP_EMBEDDINGS_BACKEND=onnx
HAWP_EMBEDDINGS_MODEL=bge-base-en-v1.5

# LLM
HAWP_LLM_BACKEND=onnx

# API Keys (env-based, not in config.json)
HAWP_OPENAI_API_KEY=sk-...
HAWP_ANTHROPIC_API_KEY=sk-ant-...
HAWP_OLLAMA_URL=http://localhost:11434
```

### Config Priority (Highest to Lowest)

1. CLI flags (e.g., `hawp search --context --embeddings-backend openai`)
2. Environment variables (`HAWP_*`)
3. Config file (`~/.hawp/config/context.json`)
4. Defaults (ONNX + BGE-base-en-v1.5)

### Go Type Structure

```go
// ContextConfig represents the full configuration for context reshaping
type ContextConfig struct {
    Embeddings EmbeddingsConfig `json:"embeddings"`
    LLM        LLMConfig        `json:"llm"`
    Backends   BackendsConfig   `json:"backends"`
    Security   SecurityConfig   `json:"security"`
}

// EmbeddingsConfig specifies which embedding backend to use
type EmbeddingsConfig struct {
    Backend  string // "onnx" | "ollama" | "openai" | "anthropic"
    Model    string // Model name/ID per backend
    CacheDir string // Where to store downloaded models
}

// LLMConfig specifies which LLM backend to use
type LLMConfig struct {
    Backend      string  // "onnx" | "ollama" | "openai" | "anthropic"
    Model        string  // Model name per backend
    Temperature  float32 // 0.0-1.0
    MaxTokens    int     // Output token limit
}

// BackendsConfig holds API URLs and encrypted credentials
type BackendsConfig struct {
    Ollama    OllamaBackend    `json:"ollama"`
    OpenAI    OpenAIBackend    `json:"openai"`
    Anthropic AnthropicBackend `json:"anthropic"`
}

type OllamaBackend struct {
    URL string `json:"url"` // http://localhost:11434
}

type OpenAIBackend struct {
    APIKey  string `json:"api_key"` // [ENCRYPTED]
    BaseURL string `json:"base_url"`
}

type AnthropicBackend struct {
    APIKey  string `json:"api_key"` // [ENCRYPTED]
    BaseURL string `json:"base_url"`
}

// SecurityConfig for encryption/key management
type SecurityConfig struct {
    EncryptionMethod string `json:"encryption_method"` // "base64" or "aes256"
    KeyRotationDays  int    `json:"key_rotation_days"` // 0 = disabled
}

// ConfigLoader loads config from all sources with proper priority
func LoadContextConfig() (ContextConfig, error) {
    // 1. Start with defaults
    // 2. Merge from config file (~/.hawp/config/context.json)
    // 3. Override with env vars (HAWP_*)
    // 4. CLI flags override everything (handled in cli/run.go)
}

// Validate checks config consistency and backend availability
func (c ContextConfig) Validate() error {
    // Check backend exists
    // Check model name format
    // Check API keys present if needed
    // Return friendly errors
}

// Encrypt/Decrypt for API keys
func EncryptKey(plaintext, method string) (string, error)
func DecryptKey(ciphertext, method string) (string, error)
```

---

## Acceptance Criteria

- [ ] Config struct defined with all backends
- [ ] Config file loading from `~/.hawp/config/context.json`
- [ ] Environment variable override (HAWP_*)
- [ ] CLI flag override (stubbed for Phase 3.4)
- [ ] Config validation (backend exists, required fields present)
- [ ] API key encryption/decryption (base64 for v0.0.3, upgrade to AES256 later)
- [ ] Tests: 8+ for loading, 5+ for validation, 3+ for encryption
- [ ] Documentation: config format + examples in help text

---

## Implementation Steps

1. **Define types** (`context_config.go`)
   - ContextConfig, EmbeddingsConfig, LLMConfig, BackendsConfig, etc.

2. **Config file loading** (`config_loader.go`)
   - Read from `~/.hawp/config/context.json`
   - Parse JSON with defaults
   - Return ContextConfig struct

3. **Environment variable support** (`config_env.go`)
   - Override from HAWP_* vars
   - Priority: env > config file > defaults

4. **Validation** (`config_validation.go`)
   - Check backend name validity
   - Check model name format
   - Check required fields for chosen backend

5. **Encryption/Decryption** (`config_encryption.go`)
   - Base64 encode/decode for v0.0.3
   - Simple symmetric cipher (no key derivation needed yet)
   - Never log keys, only [REDACTED]

6. **CLI integration stub** (`cli/context_config_flags.go`)
   - Parse `--embeddings-backend`, `--llm-backend`, `--openai-key` flags
   - Pass to config system
   - Full integration in Phase 3.4

7. **Tests** (`*_test.go`)
   - Config loading from file
   - Env var override
   - Validation (good and bad configs)
   - Encryption round-trip
   - Priority ordering (env > file > default)

8. **Documentation** (help text + examples)
   - Config file format with examples
   - Environment variables list
   - CLI flags preview

---

## Effort Estimate

| Task | Est. Time |
|------|-----------|
| Type definitions | 30 min |
| Config file loading | 1 hour |
| Env var support | 30 min |
| Validation | 45 min |
| Encryption | 45 min |
| CLI stub | 30 min |
| Tests (16+) | 2 hours |
| Docs | 30 min |
| **Total** | **~6-7 hours** |

---

## Dependencies

- ✅ v0.0.2 released (context packing CLI working)
- ⏳ No external dependencies (use Go stdlib: encoding/json, os, ioutil)

---

## Unblocks

- Phase 3.2: Embedding backends (need config to know which backend/model to use)
- Phase 3.3: LLM backends (need config to know which backend/model/credentials to use)
- Phase 3.4: Context reshaping (needs both 3.2 and 3.3 to be available)

---

## Success Metrics

✅ **Config system works:**
- Load from file + env vars + CLI flags with proper priority
- Validate backend selection
- Encrypt/decrypt API keys safely

✅ **Ready for downstream:**
- Phase 3.2 can instantiate embedding backends based on config
- Phase 3.3 can instantiate LLM backends based on config
- No config duplication or tight coupling

✅ **User-friendly:**
- Clear error messages for bad configs
- Sensible defaults (ONNX + BGE-base)
- Works with no config file (env vars or defaults)

---

## Notes

- This phase is **foundation only** — it doesn't actually call embeddings or LLM backends
- Phase 3.2 and 3.3 will **import and use** this config system
- Config changes in future releases (v0.0.4+) should be backward compatible
- Start simple (base64 encryption), upgrade to AES256 when needed

---

**Status:** Ready to implement. No blockers.

This unblocks Phase 3.2 (embeddings) and Phase 3.3 (LLM) to run in parallel.
