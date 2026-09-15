package init

import (
	"flag"
	"fmt"
	"io"
	"os"
	"sort"

	appkitsync "github.com/sentzunhat/hawp/librarian/src/internal/application/kitsync"
	appprovision "github.com/sentzunhat/hawp/librarian/src/internal/application/provision"
	domainupdate "github.com/sentzunhat/hawp/librarian/src/internal/domain/update"
	download "github.com/sentzunhat/hawp/librarian/src/internal/infrastructure/clients/download"
	githubrelease "github.com/sentzunhat/hawp/librarian/src/internal/infrastructure/clients/githubrelease"
	"github.com/sentzunhat/hawp/librarian/src/internal/infrastructure/repo"
	"github.com/sentzunhat/hawp/librarian/src/internal/platform/exitcode"
	appmcp "github.com/sentzunhat/hawp/librarian/src/internal/platform/mcp/configure"
)

func Run(args []string, cwd string) error {
	parsed, err := parseInitArgs(args)
	if err != nil {
		return err
	}
	providers := parsed.providers

	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	result := appprovision.Run(download.NewHTTPFetcher(), home, appprovision.DefaultRegistry())
	fmt.Print(result.String())
	// Asset failures are non-blocking: kit sync and provider config writes are
	// independent of model/runtime downloads and must proceed regardless.
	// We carry the failure forward and return exit 1 at the very end.

	client := githubrelease.NewClient()
	if err := doKitSync(client, providers, cwd); err != nil {
		return err
	}

	if len(providers) > 0 {
		root, rerr := repo.FindBacklogRepoRoot(cwd)
		if rerr != nil {
			fmt.Println("Not in a HAWP repo; skipping MCP config write.")
		} else if err := appmcp.WriteProviderConfigs(root, providers); err != nil {
			return err
		}
	}

	if result.Failed() {
		return exitcode.Error{Code: 1}
	}
	return nil
}

type initArgs struct {
	providers []string
}

func parseInitArgs(args []string) (initArgs, error) {
	opts := initArgs{}
	flags := flag.NewFlagSet("init", flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	flags.Func("provider", "write MCP config for this provider (repeatable)", func(value string) error {
		if value == "" {
			return fmt.Errorf("--provider requires a non-empty name")
		}
		opts.providers = append(opts.providers, value)
		return nil
	})
	flags.Bool("no-update-check", false, "suppress update notice")
	if err := flags.Parse(args); err != nil {
		return initArgs{}, fmt.Errorf("init: %w", err)
	}
	if flags.NArg() != 0 {
		return initArgs{}, fmt.Errorf("init: unexpected argument %q", flags.Arg(0))
	}
	return opts, nil
}

func sortedProviderNames(m map[string]int) []string {
	names := make([]string, 0, len(m))
	for name := range m {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
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
