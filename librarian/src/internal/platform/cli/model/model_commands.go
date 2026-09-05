package modelcmd

import (
	"github.com/sentzunhat/hawp/librarian/src/internal/platform/cli/model/embed"
	"github.com/sentzunhat/hawp/librarian/src/internal/platform/cli/model/pull"
	searchembed "github.com/sentzunhat/hawp/librarian/src/internal/platform/cli/model/search-embed"
)

func RunPull(args []string) error                    { return pull.Run(args) }
func RunEmbed(args []string) error                   { return embed.Run(args) }
func RunSearchEmbed(args []string, cwd string) error { return searchembed.Run(args, cwd) }
