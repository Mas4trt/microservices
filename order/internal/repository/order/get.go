package order

import (
	"context"

	"github.com/google/uuid"

	"github.com/Mas4trt/microservices/order/internal/model"
	repoConverter "github.com/Mas4trt/microservices/order/internal/repository/converter"
)

func (r *repository) Get(ctx context.Context, orderUUID uuid.UUID) (model.OrderDto, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	order, ok := r.orders[orderUUID]
	if !ok {
		return model.OrderDto{}, model.ErrOrderNotFound
	}

	return repoConverter.RepoToServiceModel(order), nil
}
