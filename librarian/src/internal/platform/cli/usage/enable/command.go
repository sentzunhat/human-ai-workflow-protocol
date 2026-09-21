package enable

import (
	"fmt"
	"os"

	"github.com/sentzunhat/hawp/librarian/src/internal/infrastructure/filesystem"
	usageinfra "github.com/sentzunhat/hawp/librarian/src/internal/infrastructure/repositories/usage"
)

func Run(args []string) error {
	opts, err := parseArgs(args)
	if err != nil {
		return err
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("could not resolve home directory: %w", err)
	}
	h := filesystem.ResolveHawpHome(home)
	cfg := usageinfra.LoadConfig(h.UsageConfigFile)
	cfg.Enabled = true
	cfg.LogBodies = opts.logBodies
	if err := usageinfra.SaveConfig(h.UsageConfigFile, cfg); err != nil {
		return err
	}
	msg := "Usage logging enabled."
	if cfg.LogBodies {
		msg += " Body capture on (raw input/output stored)."
	}
	fmt.Println(msg)
	return nil
}
