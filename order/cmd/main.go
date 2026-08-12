package cmd

import (
	"errors"
	"fmt"
	"sync"

	"github.com/google/uuid"

	orderV1 "github.com/Mas4trt/microservices/shared/pkg/openapi/order/v1"
)

// const (
// 	httpPort          = "8080"
// 	readHeaderTimeout = 5 * time.Second
// 	shutdownTimeout   = 10 * time.Second
// )

type OrderStorage struct {
	mu     sync.RWMutex
	orders map[string]*orderV1.OrderDto
}

func NewOrderStorage() *OrderStorage {
	return &OrderStorage{
		orders: make(map[string]*orderV1.OrderDto),
	}
}

func (s *OrderStorage) CreateOrder(uuid string, order *orderV1.OrderDto) error {
	if order == nil {
		return errors.New("order cannot be nil")
	}

	orderCopy := cloneOrder(order)

	s.mu.Lock()
	defer s.mu.Unlock()

	s.orders[uuid] = orderCopy

	return nil
}

func (s *OrderStorage) GetOrder(uuid string) *orderV1.OrderDto {
	s.mu.RLock()
	defer s.mu.RUnlock()

	order, ok := s.orders[uuid]
	if !ok {
		return nil
	}

	return cloneOrder(order)
}

func (s *OrderStorage) UpdateOrder(uuid string, order *orderV1.OrderDto) error {
	if order == nil {
		return errors.New("order cannot be nil")
	}

	orderCopy := cloneOrder(order)

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.orders[uuid]; !exists {
		return fmt.Errorf("order %s not found", uuid)
	}

	s.orders[uuid] = orderCopy

	return nil
}

func cloneOrder(order *orderV1.OrderDto) *orderV1.OrderDto {
	if order == nil {
		return nil
	}

	clone := *order

	if order.PartUuids != nil {
		clone.PartUuids = append([]uuid.UUID(nil), order.PartUuids...)
	}

	return &clone
}

func main() {

}
