package db

import (
	"fmt"
	"path/filepath"

	"github.com/sentzunhat/hawp/librarian/go/internal/infrastructure/filesystem"
)

type InitResult struct {
	DBPath        string
	ModelsPath    string
	DownloadsPath string
}

type InitService struct {
	layout filesystem.LayoutService
}

func NewInitService(layout filesystem.LayoutService) InitService {
	return InitService{layout: layout}
}

func (service InitService) Execute() (InitResult, error) {
	home, err := service.layout.HomeDir()
	if err != nil {
		return InitResult{}, err
	}

	root := filepath.Join(home, ".hawp")
	result := InitResult{
		DBPath:        filepath.Join(root, "index", "librarian.db"),
		ModelsPath:    filepath.Join(root, "models"),
		DownloadsPath: filepath.Join(root, "cache", "downloads"),
	}

	return result, nil
}

func (result InitResult) String() string {
	return fmt.Sprintf(
		"db:init\n=======\ndb: %s\nmodels: %s\ndownloads: %s",
		result.DBPath,
		result.ModelsPath,
		result.DownloadsPath,
	)
}
