package part

import (
	"sync"

	def "github.com/Mas4trt/microservices/inventory/internal/repository"
	repoModel "github.com/Mas4trt/microservices/inventory/internal/repository/model"
)

var _ def.InventoryRepository = (*repository)(nil)

type repository struct {
	mu    sync.RWMutex
	parts map[string]repoModel.Part
}

func NewRepository() *repository {
	return &repository{
		parts: make(map[string]repoModel.Part),
	}
}
