package filesystem

import "os"

type LayoutService struct{}

func NewLayoutService() LayoutService {
	return LayoutService{}
}

func (service LayoutService) HomeDir() (string, error) {
	return os.UserHomeDir()
}
