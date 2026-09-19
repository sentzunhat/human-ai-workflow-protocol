package usage

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	appusage "github.com/sentzunhat/hawp/librarian/src/internal/application/usage"
	domainusage "github.com/sentzunhat/hawp/librarian/src/internal/domain/usage"
	"github.com/sentzunhat/hawp/librarian/src/internal/infrastructure/filesystem"
	usageinfra "github.com/sentzunhat/hawp/librarian/src/internal/infrastructure/repositories/usage"
	"github.com/sentzunhat/hawp/librarian/src/internal/platform/cli/usage/enable"
)

func RunEnable(args []string) error {
	return enable.Run(args)
}

func RunDisable() error {
	h, err := home()
	if err != nil {
		return err
	}
	cfg := usageinfra.LoadConfig(h.UsageConfigFile)
	cfg.Enabled = false
	if err := usageinfra.SaveConfig(h.UsageConfigFile, cfg); err != nil {
		return err
	}
	fmt.Println("Usage logging disabled.")
	return nil
}

func RunLog() error {
	h, err := home()
	if err != nil {
		return err
	}
	entries, err := appusage.RecentLog(h.UsageDB, 20)
	if err != nil {
		return fmt.Errorf("usage db: %w", err)
	}
	if len(entries) == 0 {
		fmt.Println("No entries recorded yet.")
		return nil
	}
	for _, e := range entries {
		fmt.Printf("%s  %-22s  in=%-5d out=%-5d  %s\n",
			e.TS.Format("2006-01-02 15:04:05"),
			e.Tool, e.TokensIn, e.TokensOut, domainusage.EntrySummary(e))
	}
	return nil
}

func RunReport(args []string) error {
	exportPath := ""
	for i := 0; i < len(args); i++ {
		if args[i] == "--export" && i+1 < len(args) {
			exportPath = args[i+1]
			i++
		}
	}
	h, err := home()
	if err != nil {
		return err
	}
	rep, err := appusage.GetReport(h.UsageDB)
	if err != nil {
		return fmt.Errorf("usage db: %w", err)
	}
	out := domainusage.FormatReport(rep)
	fmt.Print(out)
	if exportPath != "" {
		if err := os.MkdirAll(filepath.Dir(exportPath), 0o755); err != nil {
			return fmt.Errorf("create export dir: %w", err)
		}
		if err := os.WriteFile(exportPath, []byte(out), 0o644); err != nil {
			return fmt.Errorf("write report: %w", err)
		}
		fmt.Printf("\nReport written to %s\n", exportPath)
	}
	return nil
}

func RunClear() error {
	h, err := home()
	if err != nil {
		return err
	}
	fmt.Print("Delete all usage log entries? [y/N] ")
	var answer string
	fmt.Scanln(&answer)
	if strings.ToLower(strings.TrimSpace(answer)) != "y" {
		fmt.Println("Cancelled.")
		return nil
	}
	if err := appusage.ClearLog(h.UsageDB); err != nil {
		return fmt.Errorf("usage db: %w", err)
	}
	fmt.Println("Usage log cleared.")
	return nil
}

func RunTotals() error {
	h, err := home()
	if err != nil {
		return err
	}
	cfg := usageinfra.LoadConfig(h.UsageConfigFile)
	if !cfg.Enabled {
		fmt.Println("Usage logging is disabled. Run `hawp usage enable` to start recording calls.")
		return nil
	}
	totals, err := appusage.GetTotals(h.UsageDB)
	if err != nil {
		return fmt.Errorf("usage db: %w", err)
	}
	fmt.Print(domainusage.FormatTotals(totals))
	return nil
}

func home() (filesystem.HawpHome, error) {
	dir, err := os.UserHomeDir()
	if err != nil {
		return filesystem.HawpHome{}, fmt.Errorf("could not resolve home directory: %w", err)
	}
	return filesystem.ResolveHawpHome(dir), nil
}
