package order

import (
	"context"

	"github.com/google/uuid"

	"github.com/Mas4trt/microservices/order/internal/model"
	repoConverter "github.com/Mas4trt/microservices/order/internal/repository/converter"
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
