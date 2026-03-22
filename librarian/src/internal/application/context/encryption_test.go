package context

import (
	"strings"
	"testing"
)

func TestEncryptDecryptBase64(t *testing.T) {
	plaintext := "sk-test-api-key-12345"

	// Encrypt
	encrypted, err := EncryptKey(plaintext, "base64")
	if err != nil {
		t.Fatalf("encryption failed: %v", err)
	}

	if encrypted == plaintext {
		t.Error("encrypted key should not be same as plaintext")
	}

	// Decrypt
	decrypted, err := DecryptKey(encrypted, "base64")
	if err != nil {
		t.Fatalf("decryption failed: %v", err)
	}

	if decrypted != plaintext {
		t.Errorf("decrypted key should match plaintext, got %s", decrypted)
	}
}

func TestEncryptKeyEmpty(t *testing.T) {
	encrypted, err := EncryptKey("", "base64")
	if err != nil {
		t.Fatalf("encrypting empty string failed: %v", err)
	}

	if encrypted != "" {
		t.Errorf("encrypted empty string should be empty, got %s", encrypted)
	}
}

func TestDecryptKeyEmpty(t *testing.T) {
	decrypted, err := DecryptKey("", "base64")
	if err != nil {
		t.Fatalf("decrypting empty string failed: %v", err)
	}

	if decrypted != "" {
		t.Errorf("decrypted empty string should be empty, got %s", decrypted)
	}
}

func TestDecryptKeyInvalid(t *testing.T) {
	// Invalid base64
	_, err := DecryptKey("!!!invalid base64!!!", "base64")
	if err == nil {
		t.Error("decrypting invalid base64 should fail")
	}
}

func TestEncryptUnsupportedMethod(t *testing.T) {
	_, err := EncryptKey("test", "unsupported-method")
	if err == nil {
		t.Error("encrypting with unsupported method should fail")
	}
}

func TestDecryptUnsupportedMethod(t *testing.T) {
	_, err := DecryptKey("test", "unsupported-method")
	if err == nil {
		t.Error("decrypting with unsupported method should fail")
	}
}

func TestRedactKey(t *testing.T) {
	tests := []struct {
		key      string
		expected string
	}{
		{"", "[EMPTY]"},
		{"a", "[REDACTED:a...]"},
		{"sk", "[REDACTED:sk...]"},
		{"sk-", "[REDACTED:sk-...]"},
		{"sk-test-key-123", "[REDACTED:sk-...]"},
	}

	for _, tc := range tests {
		result := RedactKey(tc.key)
		if result != tc.expected {
			t.Errorf("RedactKey(%q) = %q, want %q", tc.key, result, tc.expected)
		}
	}
}

func TestGetDecryptedAPIKeyPlaintext(t *testing.T) {
	// Key that looks like plaintext (starts with sk-)
	key := "sk-test-api-key"

	decrypted, err := GetDecryptedAPIKey(key, "base64")
	if err != nil {
		t.Fatalf("GetDecryptedAPIKey failed: %v", err)
	}

	if decrypted != key {
		t.Errorf("plaintext key should be returned as-is, got %s", decrypted)
	}
}

func TestGetDecryptedAPIKeyEncrypted(t *testing.T) {
	plaintext := "sk-test-api-key"

	// Encrypt first
	encrypted, err := EncryptKey(plaintext, "base64")
	if err != nil {
		t.Fatalf("encryption failed: %v", err)
	}

	// Now decrypt using GetDecryptedAPIKey
	decrypted, err := GetDecryptedAPIKey(encrypted, "base64")
	if err != nil {
		t.Fatalf("GetDecryptedAPIKey failed: %v", err)
	}

	if decrypted != plaintext {
		t.Errorf("decrypted key should match plaintext, got %s", decrypted)
	}
}

func TestGetDecryptedAPIKeyEmpty(t *testing.T) {
	decrypted, err := GetDecryptedAPIKey("", "base64")
	if err != nil {
		t.Fatalf("GetDecryptedAPIKey on empty string failed: %v", err)
	}

	if decrypted != "" {
		t.Errorf("decrypted empty string should be empty, got %s", decrypted)
	}
}

func TestMultipleKeysEncryption(t *testing.T) {
	keys := []string{
		"sk-openai-key-12345",
		"sk-ant-anthropic-key-98765",
		"ollama-url",
	}

	for _, key := range keys {
		encrypted, err := EncryptKey(key, "base64")
		if err != nil {
			t.Fatalf("encryption of %s failed: %v", key, err)
		}

		decrypted, err := DecryptKey(encrypted, "base64")
		if err != nil {
			t.Fatalf("decryption of %s failed: %v", key, err)
		}

		if decrypted != key {
			t.Errorf("encryption/decryption cycle failed for %s", key)
		}
	}
}

func TestRedactKeyDoesNotLeakValue(t *testing.T) {
	key := "sk-super-secret-api-key-12345"
	redacted := RedactKey(key)

	// Should not contain the full key
	if strings.Contains(redacted, "secret") {
		t.Error("redacted key should not contain sensitive parts of original key")
	}
	if strings.Contains(redacted, "12345") {
		t.Error("redacted key should not contain sensitive parts of original key")
	}

	// Should contain redaction marker
	if !strings.Contains(redacted, "[REDACTED") {
		t.Error("redacted key should contain [REDACTED marker")
	}
}
