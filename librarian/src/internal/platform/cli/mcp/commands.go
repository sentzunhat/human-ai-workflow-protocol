package mcp

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"

	domainupdate "github.com/sentzunhat/hawp/librarian/src/internal/domain/update"
	"github.com/sentzunhat/hawp/librarian/src/internal/infrastructure/repo"
	"github.com/sentzunhat/hawp/librarian/src/internal/platform/cli/mcp/configure"
	appmcp "github.com/sentzunhat/hawp/librarian/src/internal/platform/mcp/server"
)

func RunServer(args []string, cwd string) error {
	root, err := resolveRoot(args, cwd)
	if err != nil {
		return err
	}
	return appmcp.Serve(root, domainupdate.Version)
}

func RunConfigure(args []string, cwd string) error {
	return configure.Run(args, cwd)
}

func resolveRoot(args []string, cwd string) (string, error) {
	flags := flag.NewFlagSet("mcp", flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	var root string
	flags.StringVar(&root, "repo-root", "", "repository served by MCP")
	flags.Bool("no-update-check", false, "suppress update notification")
	if err := flags.Parse(args); err != nil {
		return "", err
	}
	if flags.NArg() != 0 {
		return "", fmt.Errorf("unexpected mcp argument %q", flags.Arg(0))
	}
	explicit := false
	flags.Visit(func(f *flag.Flag) {
		if f.Name == "repo-root" {
			explicit = true
		}
	})
	if explicit {
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
	if resolved, err := repo.FindBacklogRepoRoot(cwd); err == nil {
		return resolved, nil
	}
	return cwd, nil
}
