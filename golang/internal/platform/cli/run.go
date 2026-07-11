package cli

import (
	"errors"
	"fmt"

	appdb "github.com/sentzunhat/hawp/golang/internal/application/db"
	appindex "github.com/sentzunhat/hawp/golang/internal/application/index"
	domainindex "github.com/sentzunhat/hawp/golang/internal/domain/index"
	"github.com/sentzunhat/hawp/golang/internal/infrastructure/filesystem"
)

func Run(args []string) error {
	if len(args) == 0 || args[0] == "--help" || args[0] == "-h" {
		fmt.Println(helpText())
		return nil
	}

	if len(args) >= 2 && args[0] == "db" && args[1] == "init" {
		service := appdb.NewInitService(filesystem.NewLayoutService())
		result, err := service.Execute()
		if err != nil {
			return err
		}

		fmt.Println(result.String())
		return nil
	}

	if len(args) >= 2 && args[0] == "index" && args[1] == "build" {
		scope := domainindex.ScopeAll
		if len(args) >= 4 && args[2] == "--scope" {
			scope = domainindex.DocumentScope(args[3])
		}

		service := appindex.NewBuildService()
		result := service.Execute(scope)
		fmt.Println(result.String())
		return nil
	}

	return errors.New("unknown command\n\n" + helpText())
}

func helpText() string {
	return `hawp

Minimal Zacatl-shaped Go CLI scaffold for the future librarian port.

USAGE
  hawp <command> [options]

COMMANDS
  db init
  index build [--scope all|work|kit]`
}
