package service

import (
	"context"

	"github.com/google/uuid"

	"github.com/Mas4trt/microservices/order/internal/model"
)

type OrderService interface {
	Create(ctx context.Context, userUUID uuid.UUID, partUUIDs []uuid.UUID) (uuid.UUID, float64, error)
	Get(ctx context.Context, orderUUID uuid.UUID) (model.OrderDto, error)
	Pay(ctx context.Context, orderUUID uuid.UUID, method model.PaymentMethod) (uuid.UUID, error)
	Cancel(ctx context.Context, orderUUID uuid.UUID) error
}
