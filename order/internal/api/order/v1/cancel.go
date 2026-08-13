package v1

import (
	"context"
	"errors"

	"github.com/Mas4trt/microservices/order/internal/model"
	orderV1 "github.com/Mas4trt/microservices/shared/pkg/openapi/order/v1"
)

func (a *api) CancelOrder(ctx context.Context, params orderV1.CancelOrderParams) (orderV1.CancelOrderRes, error) {
	err := a.orderService.Cancel(ctx, params.OrderUUID)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrOrderNotFound):
			return &orderV1.CancelOrderNotFound{
				Code:    "ORDER_NOT_FOUND",
				Message: "order not found. orderUUID: " + params.OrderUUID.String(),
			}, nil
		case errors.Is(err, model.ErrOrderAlreadyPaid):
			return &orderV1.CancelOrderConflict{
				Code:    "ORDER_ALREADY_PAID",
				Message: "order is already paid and cannot be cancelled. orderUUID: " + params.OrderUUID.String(),
			}, nil
		case errors.Is(err, model.ErrOrderAlreadyCancelled):
			return &orderV1.CancelOrderConflict{
				Code:    "ORDER_ALREADY_CANCELLED",
				Message: "order is already cancelled. orderUUID: " + params.OrderUUID.String(),
			}, nil
		case errors.Is(err, model.ErrInvalidOrderStatus):
			return &orderV1.CancelOrderConflict{
				Code:    "INVALID_ORDER_STATUS",
				Message: "order cannot be cancelled in current status",
			}, nil
		default:
			return &orderV1.CancelOrderInternalServerError{
				Code:    "INTERNAL_ERROR",
				Message: "failed to cancel order",
			}, nil
		}
	}

	return &orderV1.CancelOrderNoContent{}, nil
}
