package order

import (
	"context"
	"errors"

	"github.com/google/uuid"

	"github.com/Mas4trt/microservices/order/internal/model"
)

func (s *service) Pay(ctx context.Context, orderUUID uuid.UUID, method model.PaymentMethod) (uuid.UUID, error) {
	order, err := s.orderRepository.Get(ctx, orderUUID)
	if err != nil {
		if errors.Is(err, model.ErrOrderNotFound) {
			return uuid.Nil, model.ErrOrderNotFound
		}
		return uuid.Nil, err
	}

	switch order.Status {
	case model.OrderStatusPAID:
		return uuid.Nil, model.ErrOrderAlreadyPaid
	case model.OrderStatusCANCELLED:
		return uuid.Nil, model.ErrOrderAlreadyCancelled
	case model.OrderStatusPENDINGPAYMENT:
		// Order can be paid.
	default:
		return uuid.Nil, model.ErrInvalidOrderStatus
	}

	resp, err := s.paymentClient.PayOrder(ctx, order.UserUUID.String(), order.OrderUUID.String(), method)
	if err != nil {
		return uuid.Nil, err
	}

	transactionUUID, err := uuid.Parse(resp)
	if err != nil {
		return uuid.Nil, err
	}

	order.TransactionUUID = &transactionUUID
	order.PaymentMethod = method
	order.Status = model.OrderStatusPAID

	err = s.orderRepository.Update(ctx, orderUUID, order)
	if err != nil {
		return uuid.Nil, err
	}

	return transactionUUID, nil
}
