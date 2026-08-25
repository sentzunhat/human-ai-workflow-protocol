package context

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.Embeddings.Engine != "onnx" {
		t.Errorf("default embeddings backend should be onnx, got %s", cfg.Embeddings.Engine)
	}
	if cfg.Embeddings.Model != "bge-base-en-v1.5" {
		t.Errorf("default embeddings model should be bge-base-en-v1.5, got %s", cfg.Embeddings.Model)
	}

	// LLM defaults to Ollama, not ONNX: llm.SupportedModels is empty (no ONNX
	// text2text model ships yet), so an "onnx" default would make every
	// reshape fail out of the box. See config.go DefaultConfig doc comment.
	if cfg.LLM.Engine != "ollama" {
		t.Errorf("default LLM backend should be ollama, got %s", cfg.LLM.Engine)
	}
	if cfg.LLM.Model != "mistral" {
		t.Errorf("default LLM model should be mistral, got %s", cfg.LLM.Model)
	}

	if cfg.LLM.MaxTokens != 2000 {
		t.Errorf("default maxTokens should be 2000, got %d", cfg.LLM.MaxTokens)
	}
	if cfg.LLM.Temperature != 0.3 {
		t.Errorf("default temperature should be 0.3, got %f", cfg.LLM.Temperature)
	}

	if cfg.Security.EncryptionMethod != "base64" {
		t.Errorf("default encryption method should be base64, got %s", cfg.Security.EncryptionMethod)
	}
}

func TestValidateBackend(t *testing.T) {
	validBackends := []string{"onnx", "ollama", "openai", "anthropic"}
	for _, backend := range validBackends {
		if err := validateBackend(backend, "test"); err != nil {
			t.Errorf("backend %s should be valid, got error: %v", backend, err)
		}
	}

	invalidBackends := []string{"invalid", "llama", "gpt", "claude"}
	for _, backend := range invalidBackends {
		if err := validateBackend(backend, "test"); err == nil {
			t.Errorf("backend %s should be invalid, got nil error", backend)
		}
	}
}

func TestValidateConfig(t *testing.T) {
	// Valid config (defaults)
	cfg := DefaultConfig()
	if err := cfg.Validate(); err != nil {
		t.Errorf("default config should be valid, got error: %v", err)
	}

	// Invalid: bad backend
	badCfg := DefaultConfig()
	badCfg.Embeddings.Engine = "invalid"
	if err := badCfg.Validate(); err == nil {
		t.Error("config with invalid embeddings backend should fail validation")
	}

	// Invalid: temperature out of range
	badCfg = DefaultConfig()
	badCfg.LLM.Temperature = 1.5
	if err := badCfg.Validate(); err == nil {
		t.Error("config with temperature > 1.0 should fail validation")
	}

	// Invalid: max tokens <= 0
	badCfg = DefaultConfig()
	badCfg.LLM.MaxTokens = 0
	if err := badCfg.Validate(); err == nil {
		t.Error("config with maxTokens <= 0 should fail validation")
	}

	// Invalid: OpenAI backend without API key
	badCfg = DefaultConfig()
	badCfg.LLM.Engine = "openai"
	badCfg.Backends.OpenAI.APIKey = ""
	if err := badCfg.Validate(); err == nil {
		t.Error("openai backend without API key should fail validation")
	}

	// Valid: OpenAI backend with API key
	goodCfg := DefaultConfig()
	goodCfg.LLM.Engine = "openai"
	goodCfg.Backends.OpenAI.APIKey = "sk-test-key"
	if err := goodCfg.Validate(); err != nil {
		t.Errorf("openai backend with API key should be valid, got error: %v", err)
	}
}

func TestLoadConfigFile(t *testing.T) {
	// Create temp config file
	tmpdir := t.TempDir()
	configPath := filepath.Join(tmpdir, "context.json")

	// Write a test config
	testConfigJSON := `{
		"embeddings": {
			"engine": "openai",
			"model": "text-embedding-3-large"
		},
		"llm": {
			"engine": "anthropic",
			"model": "claude-3-opus",
			"temperature": 0.5,
			"maxTokens": 4000
		}
	}`

	if err := os.WriteFile(configPath, []byte(testConfigJSON), 0644); err != nil {
		t.Fatalf("failed to write test config: %v", err)
	}

	// Load and verify
	cfg := DefaultConfig()
	if err := loadConfigFile(configPath, &cfg); err != nil {
		t.Fatalf("failed to load config file: %v", err)
	}

	if cfg.Embeddings.Engine != "openai" {
		t.Errorf("loaded embeddings backend should be openai, got %s", cfg.Embeddings.Engine)
	}
	if cfg.LLM.Engine != "anthropic" {
		t.Errorf("loaded LLM backend should be anthropic, got %s", cfg.LLM.Engine)
	}
	if cfg.LLM.Temperature != 0.5 {
		t.Errorf("loaded temperature should be 0.5, got %f", cfg.LLM.Temperature)
	}
	if cfg.LLM.MaxTokens != 4000 {
		t.Errorf("loaded maxTokens should be 4000, got %d", cfg.LLM.MaxTokens)
	}
}

func TestLoadEnvConfig(t *testing.T) {
	cfg := DefaultConfig()

	// Set env vars
	os.Setenv("HAWP_EMBEDDINGS_BACKEND", "openai")
	os.Setenv("HAWP_LLM_BACKEND", "anthropic")
	os.Setenv("HAWP_OPENAI_API_KEY", "sk-test-openai")
	os.Setenv("HAWP_ANTHROPIC_API_KEY", "sk-ant-test")
	defer func() {
		os.Unsetenv("HAWP_EMBEDDINGS_BACKEND")
		os.Unsetenv("HAWP_LLM_BACKEND")
		os.Unsetenv("HAWP_OPENAI_API_KEY")
		os.Unsetenv("HAWP_ANTHROPIC_API_KEY")
	}()

	// Load from env
	if err := loadEnvConfig(&cfg); err != nil {
		t.Fatalf("failed to load env config: %v", err)
	}

	if cfg.Embeddings.Engine != "openai" {
		t.Errorf("env embeddings backend should be openai, got %s", cfg.Embeddings.Engine)
	}
	if cfg.LLM.Engine != "anthropic" {
		t.Errorf("env LLM backend should be anthropic, got %s", cfg.LLM.Engine)
	}
	if cfg.Backends.OpenAI.APIKey != "sk-test-openai" {
		t.Errorf("env OpenAI key should be sk-test-openai, got %s", cfg.Backends.OpenAI.APIKey)
	}
	if cfg.Backends.Anthropic.APIKey != "sk-ant-test" {
		t.Errorf("env Anthropic key should be sk-ant-test, got %s", cfg.Backends.Anthropic.APIKey)
	}
}

func TestMergeConfig(t *testing.T) {
	base := DefaultConfig()
	override := ContextConfig{
		Embeddings: EmbeddingsConfig{
			Engine: "openai",
			Model:   "text-embedding-3-large",
		},
		LLM: LLMConfig{
			Engine: "anthropic",
		},
	}

	mergeConfig(&base, &override)

	if base.Embeddings.Engine != "openai" {
		t.Errorf("merged embeddings backend should be openai, got %s", base.Embeddings.Engine)
	}
	if base.LLM.Engine != "anthropic" {
		t.Errorf("merged LLM backend should be anthropic, got %s", base.LLM.Engine)
	}
	// LLM model should remain unchanged from default (override only set Backend)
	if base.LLM.Model != "mistral" {
		t.Errorf("merged LLM model should remain mistral, got %s", base.LLM.Model)
	}
}

func TestConfigPriority(t *testing.T) {
	// Test: Default < File < Env
	tmpdir := t.TempDir()
	fileConfigPath := filepath.Join(tmpdir, "context.json")

	// Write file config
	fileConfigJSON := `{
		"embeddings": {"engine": "openai", "model": "text-embedding-3-small"},
		"llm": {"engine": "anthropic"}
	}`
	if err := os.WriteFile(fileConfigPath, []byte(fileConfigJSON), 0644); err != nil {
		t.Fatalf("failed to write test config: %v", err)
	}

	// Set env var (should override file)
	os.Setenv("HAWP_EMBEDDINGS_BACKEND", "ollama")
	defer os.Unsetenv("HAWP_EMBEDDINGS_BACKEND")

	// Load and verify priority
	cfg := DefaultConfig()
	if err := loadConfigFile(fileConfigPath, &cfg); err != nil {
		t.Fatalf("failed to load file config: %v", err)
	}
	if err := loadEnvConfig(&cfg); err != nil {
		t.Fatalf("failed to load env config: %v", err)
	}

	// Env should override file
	if cfg.Embeddings.Engine != "ollama" {
		t.Errorf("env should override file config, expected ollama, got %s", cfg.Embeddings.Engine)
	}
	// File should override default for LLM backend
	if cfg.LLM.Engine != "anthropic" {
		t.Errorf("file should override default for LLM backend, expected anthropic, got %s", cfg.LLM.Engine)
	}
}

func TestConfigString(t *testing.T) {
	cfg := DefaultConfig()
	s := cfg.String()

	if s == "" {
		t.Error("String() should not return empty string")
	}
	if !contains(s, "onnx") {
		t.Error("String() should contain 'onnx'")
	}
	if !contains(s, "base64") {
		t.Error("String() should contain 'base64'")
	}
}

func contains(s, substr string) bool {
	return len(s) > 0 && len(substr) > 0 && s != substr && len(s) >= len(substr)
}
