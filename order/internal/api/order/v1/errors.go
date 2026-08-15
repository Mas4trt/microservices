package v1

import (
	"errors"

	"github.com/Mas4trt/microservices/order/internal/model"
	orderV1 "github.com/Mas4trt/microservices/shared/pkg/openapi/order/v1"
)

//nolint:dupl
func HandleCreateOrderError(err error) orderV1.CreateOrderRes {
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

//nolint:dupl
func HandlePayOrderError(err error) orderV1.PayOrderRes {
	switch {
	case errors.Is(err, model.ErrOrderNotFound):
		return &orderV1.PayOrderNotFound{
			Code:    "ORDER_NOT_FOUND",
			Message: "order not found",
		}
	case errors.Is(err, model.ErrOrderAlreadyPaid):
		return &orderV1.PayOrderConflict{
			Code:    "ORDER_ALREADY_PAID",
			Message: "order is already paid",
		}
	case errors.Is(err, model.ErrOrderAlreadyCancelled):
		return &orderV1.PayOrderConflict{
			Code:    "ORDER_ALREADY_CANCELLED",
			Message: "order is already cancelled",
		}
	case errors.Is(err, model.ErrInvalidOrderStatus):
		return &orderV1.PayOrderConflict{
			Code:    "INVALID_ORDER_STATUS",
			Message: "order cannot be paid in its current status",
		}
	case errors.Is(err, model.ErrInvalidPaymentMethod):
		return &orderV1.ValidationError{
			Code:    "INVALID_PAYMENT_METHOD",
			Message: "invalid payment method",
		}
	case errors.Is(err, model.ErrPaymentUnavailable):
		return &orderV1.PayOrderBadGateway{
			Code:    "UPSTREAM_ERROR",
			Message: "payment service is unavailable",
		}
	case errors.Is(err, model.ErrPaymentInternal):
		return &orderV1.PayOrderInternalServerError{
			Code:    "INTERNAL_ERROR",
			Message: "internal server error",
		}
	default:
		return &orderV1.PayOrderInternalServerError{
			Code:    "INTERNAL_ERROR",
			Message: "internal server error",
		}
	}
}
