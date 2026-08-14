package v1

import (
	"context"
	"errors"

	"github.com/Mas4trt/microservices/order/internal/model"
	orderV1 "github.com/Mas4trt/microservices/shared/pkg/openapi/order/v1"
)

func (a *api) CreateOrder(ctx context.Context, req *orderV1.CreateOrderRequest) (orderV1.CreateOrderRes, error) {
	orderUUID, totalPrice, err := a.orderService.Create(ctx, req.UserUUID, req.PartUuids)
	if err != nil {
		return createOrderError(err), nil
	}

	return &orderV1.CreateOrderResponse{
		OrderUUID:  orderUUID,
		TotalPrice: totalPrice,
	}, nil
}

func createOrderError(err error) orderV1.CreateOrderRes {
	switch {
	case errors.Is(err, model.ErrInventoryUnavailable):
		return &orderV1.CreateOrderBadGateway{
			Code:    "UPSTREAM_ERROR",
			Message: "inventory service is unavailable",
		}

	case errors.Is(err, model.ErrInventoryInvalidArgument):
		return &orderV1.ValidationError{
			Code:    "INVALID_PART_FILTER",
			Message: "invalid request sent to inventory service",
		}

	case errors.Is(err, model.ErrInventoryInternal):
		return &orderV1.CreateOrderInternalServerError{
			Code:    "INTERNAL_ERROR",
			Message: "internal server error",
		}

	case errors.Is(err, model.ErrPartNotFound):
		return &orderV1.CreateOrderNotFound{
			Code:    "PART_NOT_FOUND",
			Message: "part not found",
		}

	case errors.Is(err, model.ErrDuplicatePartUUID):
		return &orderV1.ValidationError{
			Code:    "DUPLICATE_PART_UUID",
			Message: "duplicate part UUID",
		}

	case errors.Is(err, model.ErrEmptyPartUUIDs):
		return &orderV1.ValidationError{
			Code:    "EMPTY_PART_UUIDS",
			Message: "part UUIDs must not be empty",
		}

	case errors.Is(err, model.ErrOrderCreateFailed):
		return &orderV1.CreateOrderInternalServerError{
			Code:    "INTERNAL_ERROR",
			Message: "failed to create order",
		}

	default:
		return &orderV1.CreateOrderInternalServerError{
			Code:    "INTERNAL_ERROR",
			Message: "internal server error",
		}
	}
}
