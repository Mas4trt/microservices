package order

import (
	"context"

	"github.com/Mas4trt/microservices/order/internal/model"
	repoConverter "github.com/Mas4trt/microservices/order/internal/repository/converter"
	"github.com/google/uuid"
)

func (r *repository) Create(_ context.Context, order model.OrderDto) (uuid.UUID, error) {
	newUUID, err := uuid.NewRandom()
	if err != nil {
		return uuid.Nil, err
	}

	order.OrderUUID = newUUID
	repoOrder := repoConverter.ServiceToRepoModel(order)

	r.mu.Lock()
	defer r.mu.Unlock()

	r.orders[newUUID] = repoOrder

	return newUUID, nil
}
