package repository

import (
	"context"

	"github.com/Mas4trt/microservices/order/internal/model"
	"github.com/google/uuid"
)

type OrderRepository interface {
	Create(ctx context.Context, order model.OrderDto) (uuid.UUID, error)
	Get(ctx context.Context, uuid uuid.UUID) (model.OrderDto, error)
	Update(ctx context.Context, uuid uuid.UUID, newOrder model.OrderDto) error
}
