package order

import (
	"context"
	"errors"

	"github.com/Mas4trt/microservices/order/internal/model"
	"github.com/google/uuid"
)

func (s *service) Cancel(ctx context.Context, orderUUID uuid.UUID) error {
	order, err := s.orderRepository.Get(ctx, orderUUID)
	if err != nil {
		if errors.Is(err, model.ErrOrderNotFound) {
			return model.ErrOrderNotFound
		}
		return err
	}

	switch order.Status {
	case model.OrderStatusPENDINGPAYMENT:
		order.Status = model.OrderStatusCANCELLED
		return s.orderRepository.Update(ctx, orderUUID, order)
	case model.OrderStatusPAID:
		return model.ErrOrderAlreadyPaid
	case model.OrderStatusCANCELLED:
		return model.ErrOrderAlreadyCancelled
	default:
		// Защитная ветка: срабатывает только если в OrderStatus добавят
		// новое значение и забудут обработать его здесь явно.
		return model.ErrInvalidOrderStatus
	}
}
