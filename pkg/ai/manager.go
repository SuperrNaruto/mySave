package ai

import (
	"sync"

	"github.com/krau/SaveAny-Bot/config"
)

var (
	renameService *RenameService
	once          sync.Once
)

// GetRenameService returns the global rename service instance
func GetRenameService() *RenameService {
	once.Do(func() {
		renameService = NewRenameService(&config.Cfg.AIRename)
	})
	return renameService
}

// InitRenameService initializes the rename service with config
func InitRenameService(cfg *config.AIRename) {
	renameService = NewRenameService(cfg)
}
