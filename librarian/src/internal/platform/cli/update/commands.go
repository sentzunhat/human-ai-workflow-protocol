package update

import (
	"flag"
	"fmt"
	"io"
	"os"
	"sort"

	appkitsync "github.com/sentzunhat/hawp/librarian/src/internal/application/kitsync"
	appupdate "github.com/sentzunhat/hawp/librarian/src/internal/application/update"
	domainupdate "github.com/sentzunhat/hawp/librarian/src/internal/domain/update"
	download "github.com/sentzunhat/hawp/librarian/src/internal/infrastructure/clients/download"
	githubrelease "github.com/sentzunhat/hawp/librarian/src/internal/infrastructure/clients/githubrelease"
	"github.com/sentzunhat/hawp/librarian/src/internal/infrastructure/repo"
	"github.com/sentzunhat/hawp/librarian/src/internal/platform/exitcode"
)

func RunVerify() error {
	client := githubrelease.NewClient()
	status, err := appupdate.Check(client, domainupdate.Repo, domainupdate.Version)
	if err != nil {
		return err
	}
	if status.NoReleases {
		fmt.Println("No published releases found yet.")
		return nil
	}
	fmt.Printf("current: %s\nlatest:  %s\n", status.Current, status.Latest)
	if status.UpdateAvailable {
		fmt.Printf("Update available: run `hawp update` to install %s.\n", status.Latest)
		return exitcode.Error{Code: 1}
	}
	fmt.Println("Already up to date.")
	return nil
}

func RunLatest() error {
	client := githubrelease.NewClient()
	status, err := appupdate.Check(client, domainupdate.Repo, domainupdate.Version)
	if err != nil {
		return err
	}
	if status.NoReleases {
		fmt.Println("No published releases found yet; nothing to update to.")
		return nil
	}
	fmt.Printf("current: %s\nlatest:  %s\n", status.Current, status.Latest)
	if !status.UpdateAvailable {
		fmt.Println("Already up to date.")
		return nil
	}
	execPath, err := os.Executable()
	if err != nil {
		return err
	}
	applied, err := appupdate.Apply(download.NewHTTPFetcher(), client, domainupdate.Repo, execPath)
	if err != nil {
		return err
	}
	fmt.Printf("Updated binary to %s.\n", applied)
	return nil
}

func RunSync(args []string, cwd string) error {
	parsed, err := parseUpdateSyncArgs(args)
	if err != nil {
		return err
	}
	providers := parsed.providers
	client := githubrelease.NewClient()
	return doKitSync(client, providers, cwd)
}

func RunFull(args []string, cwd string) error {
	parsed, err := parseUpdateFullArgs(args)
	if err != nil {
		return err
	}
	providers := parsed.providers
	if len(providers) == 0 && !parsed.noProviders {
		providers = []string{"all"}
	}
	client := githubrelease.NewClient()
	status, err := appupdate.Check(client, domainupdate.Repo, domainupdate.Version)
	if err != nil {
		return err
	}
	if status.NoReleases {
		fmt.Println("No published releases found yet; nothing to update to.")
		return nil
	}
	fmt.Printf("current: %s\nlatest:  %s\n", status.Current, status.Latest)
	if !status.UpdateAvailable {
		fmt.Println("Already up to date.")
		return nil
	}
	execPath, err := os.Executable()
	if err != nil {
		return err
	}
	applied, err := appupdate.Apply(download.NewHTTPFetcher(), client, domainupdate.Repo, execPath)
	if err != nil {
		return err
	}
	fmt.Printf("Updated binary to %s.\n", applied)
	return doKitSync(client, providers, cwd)
}

func RunAutoConfig(enabled bool) error {
	if err := appupdate.SetAutoUpdate(enabled); err != nil {
		return fmt.Errorf("could not write update config: %w", err)
	}
	if enabled {
		fmt.Println("Auto-update enabled. hawp will self-install after the 21-minute countdown.")
	} else {
		fmt.Println("Auto-update disabled. hawp will still notify you about new versions.")
		fmt.Println("Run `hawp update --enable-auto` to re-enable, or `hawp update` to install manually.")
	}
	return nil
}

type updateSyncArgs struct {
	providers []string
}

func parseUpdateSyncArgs(args []string) (updateSyncArgs, error) {
	opts := updateSyncArgs{}
	flags := flag.NewFlagSet("update sync", flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	flags.Func("provider", "also sync this provider overlay (repeatable)", func(value string) error {
		if value == "" {
			return fmt.Errorf("--provider requires a non-empty name")
		}
		opts.providers = append(opts.providers, value)
		return nil
	})
	flags.Bool("no-update-check", false, "suppress update notice")
	if err := flags.Parse(args); err != nil {
		return updateSyncArgs{}, fmt.Errorf("update sync: %w", err)
	}
	if flags.NArg() != 0 {
		return updateSyncArgs{}, fmt.Errorf("update sync: unexpected argument %q", flags.Arg(0))
	}
	return opts, nil
}

func doKitSync(client githubrelease.Client, providers []string, cwd string) error {
	root, err := repo.FindBacklogRepoRoot(cwd)
	if err != nil {
		return nil // not in a HAWP repo; skip gracefully
	}
	result, err := appkitsync.Sync(download.NewHTTPFetcher(), client, domainupdate.Repo, root, providers)
	if err != nil {
		return err
	}
	if result.NoBundleAsset {
		fmt.Println("This release has no kit bundle; .hawp/kit/ and providers left unchanged.")
		return nil
	}
	fmt.Printf("Kit refreshed: %d file(s).\n", result.KitFilesWritten)
	for _, name := range sortedProviderNames(result.Providers) {
		action := "refreshed"
		if result.ProviderInstalls[name] {
			action = "installed"
		}
		fmt.Printf("Provider %s %s: %d file(s).\n", name, action, result.Providers[name])
	}
	if len(result.Providers) == 0 {
		fmt.Println("No provider overlays synced (use --provider <name>|all or run bare `hawp update`).")
	}
	return nil
}

func sortedProviderNames(m map[string]int) []string {
	names := make([]string, 0, len(m))
	for name := range m {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

type updateFullArgs struct {
	providers     []string
	noProviders   bool
	noUpdateCheck bool
}

func parseUpdateFullArgs(args []string) (updateFullArgs, error) {
	opts := updateFullArgs{}
	flags := flag.NewFlagSet("update", flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	flags.Func("provider", "also sync this provider overlay (repeatable)", func(value string) error {
		if value == "" {
			return fmt.Errorf("--provider requires a non-empty name")
		}
		opts.providers = append(opts.providers, value)
		return nil
	})
	flags.BoolVar(&opts.noProviders, "no-providers", false, "skip provider sync after update")
	flags.BoolVar(&opts.noUpdateCheck, "no-update-check", false, "suppress update notice")
	if err := flags.Parse(args); err != nil {
		return updateFullArgs{}, fmt.Errorf("update: %w", err)
	}
	if flags.NArg() != 0 {
		return updateFullArgs{}, fmt.Errorf("update: unexpected argument %q", flags.Arg(0))
	}
	return opts, nil
}
