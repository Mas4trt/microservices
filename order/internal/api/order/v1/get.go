package v1

import (
	"context"
	"errors"

	"github.com/Mas4trt/microservices/order/internal/converter"
	"github.com/Mas4trt/microservices/order/internal/model"
	orderV1 "github.com/Mas4trt/microservices/shared/pkg/openapi/order/v1"
)

func (a *api) GetOrder(ctx context.Context, params orderV1.GetOrderParams) (orderV1.GetOrderRes, error) {
	order, err := a.orderService.Get(ctx, params.OrderUUID)
	if err != nil {
		if errors.Is(err, model.ErrOrderNotFound) {
			return &orderV1.GetOrderNotFound{
				Code:    "ORDER_NOT_FOUND",
				Message: "order not found. orderUUID: " + params.OrderUUID.String(),
			}, nil
		}

		return &orderV1.GetOrderInternalServerError{
			Code:    "INTERNAL_ERROR",
			Message: "failed to get order",
		}, nil
	}

	return converter.ToProtoOrder(order), nil
}
