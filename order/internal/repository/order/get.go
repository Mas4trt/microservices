package order

import (
	"context"

	"github.com/Mas4trt/microservices/order/internal/model"
	repoConverter "github.com/Mas4trt/microservices/order/internal/repository/converter"
	"github.com/google/uuid"
)

func (r *repository) Get(ctx context.Context, uuid uuid.UUID) (model.OrderDto, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	order, ok := r.orders[uuid]
	if !ok {
		return model.OrderDto{}, model.ErrOrderNotFound
	}

	return repoConverter.RepoToServiceModel(order), nil
}
