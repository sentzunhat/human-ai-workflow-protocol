package work

import (
	"os"

	domainwork "github.com/sentzunhat/hawp/librarian/src/internal/domain/work"
)

// ReadBacklog reads BACKLOG.md at path and parses it into a Backlog.
func ReadBacklog(path string) (*domainwork.Backlog, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return domainwork.ParseBacklogMarkdown(string(content)), nil
}
