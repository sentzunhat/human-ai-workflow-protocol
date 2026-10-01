// Package cli routes hawp commands to their application services.
package cli

import (
	"errors"
	"fmt"
	"os"

	appcheck "github.com/sentzunhat/hawp/librarian/src/internal/application/check"
	appdb "github.com/sentzunhat/hawp/librarian/src/internal/application/db"
	appupdate "github.com/sentzunhat/hawp/librarian/src/internal/application/update"
	appuuid "github.com/sentzunhat/hawp/librarian/src/internal/application/uuidgen"
	domainupdate "github.com/sentzunhat/hawp/librarian/src/internal/domain/update"
	"github.com/sentzunhat/hawp/librarian/src/internal/infrastructure/filesystem"
	"github.com/sentzunhat/hawp/librarian/src/internal/infrastructure/repo"
	distributioncmd "github.com/sentzunhat/hawp/librarian/src/internal/platform/cli/distribution"
	indexcmd "github.com/sentzunhat/hawp/librarian/src/internal/platform/cli/index"
	initcmd "github.com/sentzunhat/hawp/librarian/src/internal/platform/cli/init"
	kitcmd "github.com/sentzunhat/hawp/librarian/src/internal/platform/cli/kit"
	linkscmd "github.com/sentzunhat/hawp/librarian/src/internal/platform/cli/links"
	mcpcmd "github.com/sentzunhat/hawp/librarian/src/internal/platform/cli/mcp"
	modelcmd "github.com/sentzunhat/hawp/librarian/src/internal/platform/cli/model"
	providerscmd "github.com/sentzunhat/hawp/librarian/src/internal/platform/cli/providers"
	searchcmd "github.com/sentzunhat/hawp/librarian/src/internal/platform/cli/search"
	updatecmd "github.com/sentzunhat/hawp/librarian/src/internal/platform/cli/update"
	usagecmd "github.com/sentzunhat/hawp/librarian/src/internal/platform/cli/usage"
	workcmd "github.com/sentzunhat/hawp/librarian/src/internal/platform/cli/work"
	"github.com/sentzunhat/hawp/librarian/src/internal/platform/exitcode"
)

// ExitError carries a non-zero exit code from a command that already
// reported its findings.
type ExitError = exitcode.Error

func Run(args []string) error {
	if len(args) == 0 || args[0] == "--help" || args[0] == "-h" {
		fmt.Println(helpText())
		return nil
	}

	command := args[0]
	sub := ""
	if len(args) >= 2 {
		sub = args[1]
	}

	// Emit a non-blocking update notice after the command completes.
	// Skip for mcp (long-running server), update (already handling updates),
	// version (informational), and --no-update-check.
	skipNotify := command == "mcp" || command == "update" || command == "version" ||
		containsArg(args, "--no-update-check")
	defer appupdate.CheckAndNotify(domainupdate.Version, skipNotify)

	switch {
	case command == "uuid":
		uuid, err := appuuid.New()
		if err != nil {
			return err
		}
		if len(args) >= 2 && args[1] == "--short" {
			fmt.Println(appuuid.Short(uuid))
		} else {
			fmt.Println(uuid)
		}
		return nil

	case command == "links" && sub == "check":
		return linkscmd.RunCheck(mustGetwd())

	case command == "links" && sub == "clean":
		return linkscmd.RunClean(args[2:], mustGetwd())

	case command == "kit" && sub == "validate":
		return kitcmd.RunValidate(args[2:], mustGetwd())

	case command == "kit" && sub == "normalize":
		return kitcmd.RunNormalize(args[2:], mustGetwd())

	case command == "work" && sub == "validate":
		return workcmd.RunValidate(args[2:], mustGetwd())

	case command == "work" && sub == "normalize":
		return workcmd.RunNormalize(args[2:], mustGetwd())

	case command == "work" && sub == "new":
		return workcmd.RunNew(args[2:], mustGetwd())

	case command == "work" && (sub == "status" || sub == "evidence" || sub == "decision" || sub == "note"):
		return workcmd.RunDoc(sub, args[2:], mustGetwd())

	case command == "providers" && sub == "materialize":
		return providerscmd.RunMaterialize(mustGetwd())

	case command == "providers" && sub == "validate":
		return providerscmd.RunValidate(mustGetwd())

	case command == "providers" && sub == "sync":
		return providerscmd.RunSync(mustGetwd())

	case command == "distribution" && sub == "build":
		return distributioncmd.RunBuild(mustGetwd())

	case command == "distribution" && sub == "validate":
		return distributioncmd.RunValidate(mustGetwd())

	case command == "distribution" && sub == "sync":
		return distributioncmd.RunSync(mustGetwd())

	case command == "check":
		root, err := repo.FindBacklogRepoRoot(mustGetwd())
		if err != nil {
			return err
		}
		if code := appcheck.Run(os.Stdout, os.Stderr, root); code != 0 {
			return ExitError{Code: code}
		}
		return nil

	case command == "init":
		return initcmd.Run(args[1:], mustGetwd())

	case command == "mcp" && sub == "configure":
		return mcpcmd.RunConfigure(args[2:], mustGetwd())

	case command == "mcp":
		return mcpcmd.RunServer(args[1:], mustGetwd())

	case command == "version":
		fmt.Println(domainupdate.Version)
		return nil

	case command == "update" && sub == "latest":
		return updatecmd.RunLatest()

	case command == "update" && sub == "sync":
		return updatecmd.RunSync(args[2:], mustGetwd())

	case command == "update" && sub == "verify":
		return updatecmd.RunVerify()

	case command == "update" && containsArg(args[1:], "--check"):
		return updatecmd.RunVerify()

	case command == "update" && containsArg(args[1:], "--disable-auto"):
		return updatecmd.RunAutoConfig(false)

	case command == "update" && containsArg(args[1:], "--enable-auto"):
		return updatecmd.RunAutoConfig(true)

	case command == "update":
		return updatecmd.RunFull(args[1:], mustGetwd())

	case command == "commands":
		return runCommands(args[1:])

	case command == "backlog" && sub == "validate":
		root, err := repo.FindBacklogRepoRoot(mustGetwd())
		if err != nil {
			return err
		}
		if code := appcheck.Run(os.Stdout, os.Stderr, root); code != 0 {
			return ExitError{Code: code}
		}
		return nil

	case command == "backlog" && sub == "upgrade":
		return workcmd.RunNormalize(args[2:], mustGetwd())

	case command == "db" && sub == "init":
		service := appdb.NewInitService(filesystem.NewLayoutService())
		result, err := service.Execute()
		if err != nil {
			return err
		}
		fmt.Println(result.String())
		return nil

	case command == "index" && sub == "build":
		return indexcmd.RunBuild(args[2:], mustGetwd())

	case command == "model" && sub == "pull":
		return modelcmd.RunPull(args[2:])

	case command == "embed":
		return modelcmd.RunEmbed(args[1:])

	case command == "search" && sub == "index":
		return indexcmd.RunSearchIndex(mustGetwd())

	case command == "search" && sub == "embed":
		return modelcmd.RunSearchEmbed(args[2:], mustGetwd())

	case command == "search" && sub == "benchmark":
		return searchcmd.RunBenchmark(args[2:], mustGetwd())

	case command == "search":
		return searchcmd.RunQuery(args[1:], mustGetwd())

	case command == "usage" && sub == "enable":
		return usagecmd.RunEnable(args[2:])

	case command == "usage" && sub == "disable":
		return usagecmd.RunDisable()

	case command == "usage" && sub == "log":
		return usagecmd.RunLog()

	case command == "usage" && sub == "report":
		return usagecmd.RunReport(args[2:])

	case command == "usage" && sub == "clear":
		return usagecmd.RunClear()

	case command == "usage":
		return usagecmd.RunTotals()
	}

	return errors.New("unknown command\n\n" + helpText())
}

func containsArg(args []string, flag string) bool {
	for _, a := range args {
		if a == flag {
			return true
		}
	}
	return false
}

func mustGetwd() string {
	wd, err := os.Getwd()
	if err != nil {
		return "."
	}
	return wd
}
