package v1

import (
	"context"

	orderV1 "github.com/Mas4trt/microservices/shared/pkg/openapi/order/v1"
)

func (a *api) CreateOrder(ctx context.Context, req *orderV1.CreateOrderRequest) (orderV1.CreateOrderRes, error) {
	orderUUID, totalPrice, err := a.orderService.Create(ctx, req.UserUUID, req.PartUuids)
	if err != nil {
		return HandleCreateOrderError(err), nil
	}

	return &orderV1.CreateOrderResponse{
		OrderUUID:  orderUUID,
		TotalPrice: totalPrice,
	}, nil
}
