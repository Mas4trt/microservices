package order

import (
	"sync"

	def "github.com/Mas4trt/microservices/order/internal/repository"
	repoModel "github.com/Mas4trt/microservices/order/internal/repository/model"
	"github.com/google/uuid"
)

var _ def.OrderRepository = (*repository)(nil)

type repository struct {
	mu     sync.RWMutex
	orders map[uuid.UUID]repoModel.OrderDto
}

func NewRepository() *repository {
	return &repository{
		orders: make(map[uuid.UUID]repoModel.OrderDto),
	}
}
