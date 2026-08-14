package repository

import (
	"context"

	"github.com/google/uuid"

	"github.com/Mas4trt/microservices/order/internal/model"
)

type OrderRepository interface {
	Create(ctx context.Context, order model.OrderDto) (uuid.UUID, error)
	Get(ctx context.Context, orderUUID uuid.UUID) (model.OrderDto, error)
	Update(ctx context.Context, orderUUID uuid.UUID, newOrder model.OrderDto) error
}
