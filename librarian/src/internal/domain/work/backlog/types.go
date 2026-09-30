package backlog

import (
	"github.com/sentzunhat/hawp/librarian/src/internal/domain/work/identity"
	"github.com/sentzunhat/hawp/librarian/src/internal/domain/work/model"
)

type Backlog = model.Backlog
type BacklogRow = model.BacklogRow

func ExtractShortUUID(value string) string      { return identity.ExtractShortUUID(value) }
func ExtractIDFromFilename(value string) string { return identity.ExtractIDFromFilename(value) }
