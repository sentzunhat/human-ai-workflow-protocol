package configure

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/sentzunhat/hawp/librarian/src/internal/infrastructure/repo"
	appmcp "github.com/sentzunhat/hawp/librarian/src/internal/platform/mcp/configure"
)

func Run(args []string, cwd string) error {
	flags := flag.NewFlagSet("mcp configure", flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	var root string
	var providers []string
	flags.StringVar(&root, "repo-root", "", "repository to configure")
	flags.Func("provider", "client to configure (repeatable)", func(value string) error {
		if value == "" {
			return fmt.Errorf("--provider requires a name")
		}
		providers = append(providers, value)
		return nil
	})
	flags.Bool("no-update-check", false, "suppress update notice")
	if err := flags.Parse(args); err != nil {
		return err
	}
	if flags.NArg() > 0 {
		return fmt.Errorf("unexpected mcp configure argument %q", flags.Arg(0))
	}
	resolved, err := resolveConfigureRoot(root, cwd, isSet(flags, "repo-root"))
	if err != nil {
		return err
	}
	return appmcp.Configure(resolved, providers)
}

func isSet(flags *flag.FlagSet, name string) bool {
	found := false
	flags.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}

func resolveConfigureRoot(root, cwd string, explicit bool) (string, error) {
	if !explicit {
		if resolved, err := repo.FindBacklogRepoRoot(cwd); err == nil {
			return resolved, nil
		}
		return cwd, nil
	}
	if root == "" {
		return "", fmt.Errorf("--repo-root requires a non-empty directory")
	}
	if !filepath.IsAbs(root) {
		root = filepath.Join(cwd, root)
	}
	root = filepath.Clean(root)
	info, err := os.Stat(root)
	if err != nil {
		return "", fmt.Errorf("--repo-root: %w", err)
	}
	if !info.IsDir() {
		return "", fmt.Errorf("--repo-root must be a directory")
	}
	return root, nil
}
