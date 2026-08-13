package order

import (
	"context"

	"github.com/Mas4trt/microservices/order/internal/model"
	repoConverter "github.com/Mas4trt/microservices/order/internal/repository/converter"
	"github.com/google/uuid"
)

func (r *repository) Update(ctx context.Context, uuid uuid.UUID, newOrder model.OrderDto) error {
	repoOrder := repoConverter.ServiceToRepoModel(newOrder)

	r.mu.Lock()
	defer r.mu.Unlock()

	_, ok := r.orders[uuid]
	if !ok {
		return model.ErrOrderNotFound
	}

	r.orders[uuid] = repoOrder

	return nil
}
