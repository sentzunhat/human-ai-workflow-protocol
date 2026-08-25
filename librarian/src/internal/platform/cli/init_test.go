package cli

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	domainprovision "github.com/sentzunhat/hawp/librarian/src/internal/domain/provision"
	appprovision "github.com/sentzunhat/hawp/librarian/src/internal/application/provision"
	appmcp "github.com/sentzunhat/hawp/librarian/src/internal/platform/mcp"
	"github.com/sentzunhat/hawp/librarian/src/internal/infrastructure/download"
)

// TestProviderConfigWrittenAfterProvisionFailure is a regression test for the
// runInit early-exit bug: when any asset download fails, the original code
// returned ExitError{Code: 1} before reaching WriteProviderConfigs, leaving
// Codex (and other providers) without MCP config.
//
// The fix: provision failure is now non-blocking. Kit sync and provider config
// writes always run; the exit code reflects provision outcome at the very end.
//
// This test directly exercises the two formerly-coupled steps in isolation:
// provision with guaranteed failures followed by provider config write.
func TestProviderConfigWrittenAfterProvisionFailure(t *testing.T) {
	// Run provision with assets that will always fail so we reproduce the bug condition.
	home := t.TempDir()
	registry := appprovision.Registry{
		RuntimeAssetErr: errors.New("unsupported platform (simulated)"),
		ModelAssets: []domainprovision.Asset{
			{
				Name:     "bad_model",
				URL:      "http://localhost:1/nonexistent", // unreachable
				SHA256:   strings.Repeat("0", 64),         // wrong checksum
				Size:     100,
				DestName: "bad/model.bin",
			},
		},
	}
	result := appprovision.Run(download.NewHTTPFetcher(), home, registry)
	if !result.Failed() {
		t.Fatal("expected provision to fail for this regression test")
	}

	// Despite provision failure, provider config write must succeed.
	repoRoot := t.TempDir()
	if err := appmcp.WriteProviderConfigs(repoRoot, []string{"codex"}); err != nil {
		t.Fatalf("WriteProviderConfigs failed after provision failure: %v", err)
	}

	codexTOML := filepath.Join(repoRoot, "codex.toml")
	data, err := os.ReadFile(codexTOML)
	if err != nil {
		t.Fatalf("codex.toml not written: %v — MCP integration is missing despite successful provider spec", err)
	}
	if !strings.Contains(string(data), "[mcp_servers.hawp]") {
		t.Errorf("codex.toml missing hawp MCP block:\n%s", data)
	}
}

// TestProviderConfigWrittenAfterProvisionFailure_AllProviders verifies the
// same fix for the "all" expansion: claude, cursor, codex each write their
// respective config files when provision fails.
func TestProviderConfigWrittenAfterProvisionFailure_AllProviders(t *testing.T) {
	home := t.TempDir()
	registry := appprovision.Registry{
		RuntimeAssetErr: errors.New("unsupported platform (simulated)"),
	}
	result := appprovision.Run(download.NewHTTPFetcher(), home, registry)
	if !result.Failed() {
		t.Fatal("expected provision to fail for this regression test")
	}

	repoRoot := t.TempDir()
	if err := appmcp.WriteProviderConfigs(repoRoot, []string{"claude", "cursor", "codex"}); err != nil {
		t.Fatalf("WriteProviderConfigs failed: %v", err)
	}

	for _, want := range []struct {
		path string
		text string
	}{
		{".mcp.json", `"hawp"`},
		{".cursor/mcp.json", `"hawp"`},
		{"codex.toml", "[mcp_servers.hawp]"},
	} {
		data, err := os.ReadFile(filepath.Join(repoRoot, want.path))
		if err != nil {
			t.Errorf("%s not written after provision failure: %v", want.path, err)
			continue
		}
		if !strings.Contains(string(data), want.text) {
			t.Errorf("%s missing expected content %q:\n%s", want.path, want.text, data)
		}
	}
}
