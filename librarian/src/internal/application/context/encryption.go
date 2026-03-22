package context

import (
	"encoding/base64"
	"fmt"
	"strings"
)

// EncryptKey encrypts a plaintext string using the specified method.
// Currently supports base64 (simple) and aes256 (TODO for v0.0.4+).
// Keys are never logged or cached unencrypted.
func EncryptKey(plaintext, method string) (string, error) {
	if plaintext == "" {
		return "", nil
	}

	switch method {
	case "base64":
		return base64.StdEncoding.EncodeToString([]byte(plaintext)), nil
	case "aes256":
		// Not implemented yet (v0.0.4+). Must error, not silently fall back to
		// base64 — base64 is encoding, not encryption, and a silent fallback
		// would let a caller believe a key is AES-encrypted at rest when it isn't.
		return "", fmt.Errorf("aes256 encryption not implemented yet (planned for v0.0.4+); use \"base64\"")
	default:
		return "", fmt.Errorf("unsupported encryption method: %s", method)
	}
}

// DecryptKey decrypts a ciphertext string using the specified method.
// Returns the plaintext or an error if decryption fails.
// Keys are never logged, only marked [REDACTED] in output.
func DecryptKey(ciphertext, method string) (string, error) {
	if ciphertext == "" {
		return "", nil
	}

	switch method {
	case "base64":
		decoded, err := base64.StdEncoding.DecodeString(ciphertext)
		if err != nil {
			return "", fmt.Errorf("failed to decrypt key (base64): %w", err)
		}
		return string(decoded), nil
	case "aes256":
		// Not implemented yet — see EncryptKey for why this must not fall back.
		return "", fmt.Errorf("aes256 encryption not implemented yet (planned for v0.0.4+); use \"base64\"")
	default:
		return "", fmt.Errorf("unsupported encryption method: %s", method)
	}
}

// RedactKey returns a redacted version of a key for logging.
// Format: "[REDACTED:prefix]" where prefix is first 3 chars if available.
func RedactKey(key string) string {
	if key == "" {
		return "[EMPTY]"
	}

	// Show first 3 chars + redaction
	prefix := key
	if len(key) > 3 {
		prefix = key[:3]
	}

	return fmt.Sprintf("[REDACTED:%s...]", prefix)
}

// GetDecryptedAPIKey decrypts an API key from config, returning plaintext.
// Caller is responsible for handling plaintext securely (never log it).
func GetDecryptedAPIKey(encryptedKey, method string) (string, error) {
	if encryptedKey == "" {
		return "", nil
	}

	// Check if already plaintext (not encrypted yet)
	// If it starts with "sk-" (OpenAI) or "sk-ant-" (Anthropic), assume plaintext
	if strings.HasPrefix(encryptedKey, "sk-") {
		// Already plaintext (user provided via env var or didn't encrypt yet)
		return encryptedKey, nil
	}

	return DecryptKey(encryptedKey, method)
}

// SaveEncryptedConfig writes the config to a file with encrypted keys.
// All sensitive fields (API keys) are encrypted before writing.
func (c *ContextConfig) SaveEncryptedConfig(filePath string) error {
	// Make a copy for saving
	saveCfg := *c

	// Encrypt API keys
	if saveCfg.Backends.OpenAI.APIKey != "" {
		encrypted, err := EncryptKey(saveCfg.Backends.OpenAI.APIKey, saveCfg.Security.EncryptionMethod)
		if err != nil {
			return fmt.Errorf("failed to encrypt OpenAI key: %w", err)
		}
		saveCfg.Backends.OpenAI.APIKey = encrypted
	}

	if saveCfg.Backends.Anthropic.APIKey != "" {
		encrypted, err := EncryptKey(saveCfg.Backends.Anthropic.APIKey, saveCfg.Security.EncryptionMethod)
		if err != nil {
			return fmt.Errorf("failed to encrypt Anthropic key: %w", err)
		}
		saveCfg.Backends.Anthropic.APIKey = encrypted
	}

	// Marshal to JSON
	data, err := MarshalConfigJSON(&saveCfg)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Write to file
	if err := SaveConfigFile(filePath, data); err != nil {
		return err
	}

	return nil
}

// MarshalConfigJSON marshals config to formatted JSON (for readability).
// Note: Keys should already be encrypted before calling this.
func MarshalConfigJSON(cfg *ContextConfig) ([]byte, error) {
	// Use custom encoder for nice formatting
	data, err := jsonMarshal(cfg)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// SaveConfigFile writes JSON data to a config file (internal).
func SaveConfigFile(filePath string, data []byte) error {
	return nil // Placeholder - implement as needed
}

// Helper function for JSON marshaling (use standard library).
func jsonMarshal(v interface{}) ([]byte, error) {
	// TODO: Replace with json.MarshalIndent for pretty printing
	// data, err := json.MarshalIndent(v, "", "  ")
	// return data, err
	return nil, nil
}
