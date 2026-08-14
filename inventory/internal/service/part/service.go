package part

import (
	"github.com/Mas4trt/microservices/inventory/internal/repository"
	def "github.com/Mas4trt/microservices/inventory/internal/service"
)

var _ def.InventoryService = (*service)(nil)

type service struct {
	inventoryRepository repository.InventoryRepository
}

func NewService(inventoryRepository repository.InventoryRepository) *service {
	return &service{
		inventoryRepository: inventoryRepository,
	}
}
