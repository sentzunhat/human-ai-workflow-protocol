package context

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// ContextConfig represents the full configuration for context reshaping.
// Priority (highest to lowest): CLI flags → project config → home config → defaults
type ContextConfig struct {
	Embeddings EmbeddingsConfig `json:"embeddings"`
	LLM        LLMConfig        `json:"llm"`
	Backends   BackendsConfig   `json:"backends"`
	Security   SecurityConfig   `json:"security"`
}

// EmbeddingsConfig specifies which embedding backend to use and which model.
type EmbeddingsConfig struct {
	Backend  string `json:"backend"`  // "onnx" | "ollama" | "openai" | "anthropic"
	Model    string `json:"model"`    // Model name/ID per backend
	CacheDir string `json:"cacheDir"` // Where to store cached embeddings (optional)
}

// LLMConfig specifies which LLM backend to use and parameters.
type LLMConfig struct {
	Backend     string  `json:"backend"`     // "onnx" | "ollama" | "openai" | "anthropic"
	Model       string  `json:"model"`       // Model name per backend
	Temperature float32 `json:"temperature"` // 0.0-1.0, default 0.3
	MaxTokens   int     `json:"maxTokens"`   // Output token limit, default 2000
}

// BackendsConfig holds API URLs and credentials for optional backends.
type BackendsConfig struct {
	Ollama    OllamaBackend    `json:"ollama"`
	OpenAI    OpenAIBackend    `json:"openai"`
	Anthropic AnthropicBackend `json:"anthropic"`
}

type OllamaBackend struct {
	URL string `json:"url"` // http://localhost:11434 (default)
}

type OpenAIBackend struct {
	APIKey  string `json:"apiKey"`  // Encrypted in file, decrypted at runtime
	BaseURL string `json:"baseUrl"` // https://api.openai.com/v1 (default)
}

type AnthropicBackend struct {
	APIKey  string `json:"apiKey"`  // Encrypted in file, decrypted at runtime
	BaseURL string `json:"baseUrl"` // https://api.anthropic.com (default)
}

// SecurityConfig for encryption and key management.
type SecurityConfig struct {
	EncryptionMethod string `json:"encryptionMethod"` // "base64" or "aes256"
	KeyRotationDays  int    `json:"keyRotationDays"`  // 0 = disabled
}

// DefaultConfig returns the built-in default configuration.
// Uses: ONNX embeddings (bge-base-en-v1.5, local/offline) + Ollama LLM (mistral).
// ONNX LLM is blocked on hugot's CGO ORT backend (see llm.ErrGenerativeRequiresORT)
// — defaulting LLM.Backend to "onnx" would make every reshape fail out of the
// box, so the default LLM backend is "ollama" until that infra ships.
func DefaultConfig() ContextConfig {
	return ContextConfig{
		Embeddings: EmbeddingsConfig{
			Backend: "onnx",
			Model:   "bge-base-en-v1.5",
		},
		LLM: LLMConfig{
			Backend:     "ollama",
			Model:       "mistral",
			Temperature: 0.3,
			MaxTokens:   2000,
		},
		Backends: BackendsConfig{
			Ollama: OllamaBackend{
				URL: "http://localhost:11434",
			},
			OpenAI: OpenAIBackend{
				BaseURL: "https://api.openai.com/v1",
			},
			Anthropic: AnthropicBackend{
				BaseURL: "https://api.anthropic.com",
			},
		},
		Security: SecurityConfig{
			EncryptionMethod: "base64",
			KeyRotationDays:  0,
		},
	}
}

// LoadContextConfig loads configuration with proper priority:
// 1. Built-in defaults
// 2. Home config (~/.hawp/config/context.json)
// 3. Project config (.hawp/config/context.json)
// 4. Environment variables (HAWP_*)
// CLI flags override everything (handled in cli/run.go)
func LoadContextConfig(hawpHome, projectRoot string) (ContextConfig, error) {
	cfg := DefaultConfig()

	// Step 1: Load home config if it exists (~/.hawp/config/context.json)
	if hawpHome != "" {
		homeConfigPath := filepath.Join(hawpHome, "config", "context.json")
		if err := loadConfigFile(homeConfigPath, &cfg); err != nil && !os.IsNotExist(err) {
			return cfg, fmt.Errorf("failed to load home config: %w", err)
		}
	}

	// Step 2: Load project config if it exists (.hawp/config/context.json)
	// Project config overrides home config
	if projectRoot != "" {
		projectConfigPath := filepath.Join(projectRoot, ".hawp", "config", "context.json")
		if err := loadConfigFile(projectConfigPath, &cfg); err != nil && !os.IsNotExist(err) {
			return cfg, fmt.Errorf("failed to load project config: %w", err)
		}
	}

	// Step 3: Load environment variables (override file configs)
	if err := loadEnvConfig(&cfg); err != nil {
		return cfg, fmt.Errorf("failed to load env config: %w", err)
	}

	// Step 4: Validate
	if err := cfg.Validate(); err != nil {
		return cfg, err
	}

	return cfg, nil
}

// loadConfigFile merges a JSON config file into the current config.
func loadConfigFile(path string, cfg *ContextConfig) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	// Unmarshal into a temporary config
	var fileCfg ContextConfig
	if err := json.Unmarshal(data, &fileCfg); err != nil {
		return fmt.Errorf("invalid JSON in config file %s: %w", path, err)
	}

	// Merge: non-empty fields in fileCfg override defaults
	mergeConfig(cfg, &fileCfg)

	return nil
}

// mergeConfig merges fileCfg into baseCfg, with fileCfg taking precedence.
func mergeConfig(baseCfg, fileCfg *ContextConfig) {
	// Embeddings
	if fileCfg.Embeddings.Backend != "" {
		baseCfg.Embeddings.Backend = fileCfg.Embeddings.Backend
	}
	if fileCfg.Embeddings.Model != "" {
		baseCfg.Embeddings.Model = fileCfg.Embeddings.Model
	}
	if fileCfg.Embeddings.CacheDir != "" {
		baseCfg.Embeddings.CacheDir = fileCfg.Embeddings.CacheDir
	}

	// LLM
	if fileCfg.LLM.Backend != "" {
		baseCfg.LLM.Backend = fileCfg.LLM.Backend
	}
	if fileCfg.LLM.Model != "" {
		baseCfg.LLM.Model = fileCfg.LLM.Model
	}
	if fileCfg.LLM.Temperature != 0 {
		baseCfg.LLM.Temperature = fileCfg.LLM.Temperature
	}
	if fileCfg.LLM.MaxTokens != 0 {
		baseCfg.LLM.MaxTokens = fileCfg.LLM.MaxTokens
	}

	// Backends
	if fileCfg.Backends.Ollama.URL != "" {
		baseCfg.Backends.Ollama.URL = fileCfg.Backends.Ollama.URL
	}
	if fileCfg.Backends.OpenAI.APIKey != "" {
		baseCfg.Backends.OpenAI.APIKey = fileCfg.Backends.OpenAI.APIKey
	}
	if fileCfg.Backends.OpenAI.BaseURL != "" {
		baseCfg.Backends.OpenAI.BaseURL = fileCfg.Backends.OpenAI.BaseURL
	}
	if fileCfg.Backends.Anthropic.APIKey != "" {
		baseCfg.Backends.Anthropic.APIKey = fileCfg.Backends.Anthropic.APIKey
	}
	if fileCfg.Backends.Anthropic.BaseURL != "" {
		baseCfg.Backends.Anthropic.BaseURL = fileCfg.Backends.Anthropic.BaseURL
	}

	// Security
	if fileCfg.Security.EncryptionMethod != "" {
		baseCfg.Security.EncryptionMethod = fileCfg.Security.EncryptionMethod
	}
	if fileCfg.Security.KeyRotationDays != 0 {
		baseCfg.Security.KeyRotationDays = fileCfg.Security.KeyRotationDays
	}
}

// loadEnvConfig loads configuration from environment variables (HAWP_*).
// Environment variables override file configs.
func loadEnvConfig(cfg *ContextConfig) error {
	// Embeddings
	if backend := os.Getenv("HAWP_EMBEDDINGS_BACKEND"); backend != "" {
		cfg.Embeddings.Backend = backend
	}
	if model := os.Getenv("HAWP_EMBEDDINGS_MODEL"); model != "" {
		cfg.Embeddings.Model = model
	}

	// LLM
	if backend := os.Getenv("HAWP_LLM_BACKEND"); backend != "" {
		cfg.LLM.Backend = backend
	}
	if model := os.Getenv("HAWP_LLM_MODEL"); model != "" {
		cfg.LLM.Model = model
	}

	// Backends
	if url := os.Getenv("HAWP_OLLAMA_URL"); url != "" {
		cfg.Backends.Ollama.URL = url
	}
	if key := os.Getenv("HAWP_OPENAI_API_KEY"); key != "" {
		cfg.Backends.OpenAI.APIKey = key
	}
	if key := os.Getenv("HAWP_ANTHROPIC_API_KEY"); key != "" {
		cfg.Backends.Anthropic.APIKey = key
	}

	return nil
}

// Validate checks configuration consistency and backend availability.
func (c *ContextConfig) Validate() error {
	// Validate embeddings backend
	if err := validateBackend(c.Embeddings.Backend, "embeddings"); err != nil {
		return err
	}

	// Validate LLM backend
	if err := validateBackend(c.LLM.Backend, "llm"); err != nil {
		return err
	}

	// Validate temperature
	if c.LLM.Temperature < 0 || c.LLM.Temperature > 1 {
		return fmt.Errorf("LLM temperature must be between 0 and 1, got %f", c.LLM.Temperature)
	}

	// Validate max tokens
	if c.LLM.MaxTokens <= 0 {
		return fmt.Errorf("LLM maxTokens must be positive, got %d", c.LLM.MaxTokens)
	}

	// Validate API keys for API backends (if present)
	if c.LLM.Backend == "openai" && c.Backends.OpenAI.APIKey == "" {
		return fmt.Errorf("openai backend selected but no API key provided (set HAWP_OPENAI_API_KEY or config)")
	}
	if c.LLM.Backend == "anthropic" && c.Backends.Anthropic.APIKey == "" {
		return fmt.Errorf("anthropic backend selected but no API key provided (set HAWP_ANTHROPIC_API_KEY or config)")
	}

	// aes256 is not implemented (EncryptKey/DecryptKey silently fall back to
	// base64 for it, which is encoding, not encryption). Reject explicitly
	// rather than let a user believe their keys are AES-encrypted at rest.
	if c.Security.EncryptionMethod == "aes256" {
		return fmt.Errorf("security.encryptionMethod \"aes256\" is not implemented yet; use \"base64\"")
	}

	return nil
}

// validateBackend checks if a backend name is valid.
func validateBackend(backend, kind string) error {
	validBackends := map[string]bool{
		"onnx":      true,
		"ollama":    true,
		"openai":    true,
		"anthropic": true,
		"none":      true,
	}

	if !validBackends[backend] {
		return fmt.Errorf("%s backend %q is not supported; choose from: onnx, ollama, openai, anthropic, none", kind, backend)
	}

	return nil
}

// BackendCategory classifies a backend by whether it needs a live network
// call to a third party per request:
//
//   - "none": no model, no network, no computation at all. A deliberate
//     passthrough — embeddings return empty vectors, LLM returns input
//     unchanged. Use this when you want structured search/reference output
//     without requiring any local model download or running server.
//   - "offline": runs entirely on this machine once its model is present
//     (ONNX in-process; Ollama against a local/self-hosted server). No
//     per-request API key or internet dependency.
//   - "online": calls a third-party API per request (OpenAI, Anthropic).
//     Requires an API key and network access; incurs per-request cost.
//
// v0.1.0 ships "none" and "offline" backends (ONNX + Ollama + none).
// "online" backends are defined in ContextConfig.Backends but not yet wired
// into the embeddings/llm factories — see librarian/docs/v0.1.0-vision.md.
func BackendCategory(backend string) string {
	switch backend {
	case "none":
		return "none"
	case "onnx", "ollama":
		return "offline"
	case "openai", "anthropic":
		return "online"
	default:
		return "unknown"
	}
}

// String returns a human-readable summary of the configuration.
func (c *ContextConfig) String() string {
	return fmt.Sprintf(`Context Config:
  Embeddings: %s (%s) [%s]
  LLM: %s (%s) [%s, temp=%.1f, maxTokens=%d]
  Encryption: %s`,
		c.Embeddings.Backend, c.Embeddings.Model, BackendCategory(c.Embeddings.Backend),
		c.LLM.Backend, c.LLM.Model, BackendCategory(c.LLM.Backend), c.LLM.Temperature, c.LLM.MaxTokens,
		c.Security.EncryptionMethod)
}
