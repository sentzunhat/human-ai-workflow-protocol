package work

import "github.com/sentzunhat/hawp/librarian/src/internal/domain/work/backlog"

func ParseBacklogMarkdown(content string) *Backlog {
	return backlog.ParseBacklogMarkdown(content)
}
