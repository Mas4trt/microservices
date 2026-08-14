package order

import (
	"context"

	"github.com/google/uuid"

	"github.com/Mas4trt/microservices/order/internal/model"
	repoConverter "github.com/Mas4trt/microservices/order/internal/repository/converter"
)

func (r *repository) Update(ctx context.Context, orderUUID uuid.UUID, newOrder model.OrderDto) error {
	repoOrder := repoConverter.ServiceToRepoModel(newOrder)

	r.mu.Lock()
	defer r.mu.Unlock()

	_, ok := r.orders[orderUUID]
	if !ok {
		return model.ErrOrderNotFound
	}

	r.orders[orderUUID] = repoOrder

	return nil
}
