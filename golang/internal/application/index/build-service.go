package index

import (
	"fmt"

	domainindex "github.com/sentzunhat/hawp/golang/internal/domain/index"
)

type BuildResult struct {
	Scope domainindex.DocumentScope
}

type BuildService struct{}

func NewBuildService() BuildService {
	return BuildService{}
}

func (service BuildService) Execute(scope domainindex.DocumentScope) BuildResult {
	return BuildResult{Scope: scope}
}

func (result BuildResult) String() string {
	return fmt.Sprintf(
		"index:build\n===========\nscope: %s\nstatus: poc scaffold only",
		result.Scope,
	)
}
