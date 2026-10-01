package enable

import (
	"path/filepath"
	"testing"

	domainusage "github.com/sentzunhat/hawp/librarian/src/internal/domain/usage"
	"github.com/sentzunhat/hawp/librarian/src/internal/infrastructure/filesystem"
	usageinfra "github.com/sentzunhat/hawp/librarian/src/internal/infrastructure/repositories/usage"
)

func TestRunPreservesExistingLogBodiesPreference(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	h := filesystem.ResolveHawpHome(home)
	if err := usageinfra.SaveConfig(h.UsageConfigFile, domainusage.Config{
		Enabled:   false,
		LogBodies: true,
	}); err != nil {
		t.Fatal(err)
	}

	if err := Run(nil); err != nil {
		t.Fatal(err)
	}

	cfg := usageinfra.LoadConfig(filepath.Clean(h.UsageConfigFile))
	if !cfg.Enabled {
		t.Fatal("usage logging was not enabled")
	}
	if !cfg.LogBodies {
		t.Fatal("existing log_bodies preference was cleared")
	}
}
