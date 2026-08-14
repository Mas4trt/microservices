package order

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"github.com/Mas4trt/microservices/order/internal/model"
)

func (s *service) Get(ctx context.Context, orderUUID uuid.UUID) (model.OrderDto, error) {
	order, err := s.orderRepository.Get(ctx, orderUUID)
	if err != nil {
		if errors.Is(err, model.ErrOrderNotFound) {
			return model.OrderDto{}, model.ErrOrderNotFound
		}
		return model.OrderDto{}, err
	}

	return order, nil
}
