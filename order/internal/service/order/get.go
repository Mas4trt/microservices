package order

import (
	"context"

	"github.com/Mas4trt/microservices/order/internal/model"
	"github.com/google/uuid"
)

func (s *service) Get(ctx context.Context, orderUUID uuid.UUID) (model.OrderDto, error) {
	order, err := s.orderRepository.Get(ctx, orderUUID)
	if err != nil {
		return model.OrderDto{}, model.ErrOrderNotFound
	}

	return order, nil
}
