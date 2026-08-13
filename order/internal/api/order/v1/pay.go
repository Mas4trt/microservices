package v1

import (
	"context"
	"errors"

	"github.com/Mas4trt/microservices/order/internal/api/validation"
	"github.com/Mas4trt/microservices/order/internal/model"
	orderV1 "github.com/Mas4trt/microservices/shared/pkg/openapi/order/v1"
)

func (a *api) PayOrder(ctx context.Context, req *orderV1.PayOrderRequest, params orderV1.PayOrderParams) (orderV1.PayOrderRes, error) {
	// If req = nil Validate return validate.ErrNilPointer
	if err := req.Validate(); err != nil {
		return validation.ValidationError(err), nil
	}

	resp, err := a.orderService.Pay(ctx, params.OrderUUID, model.PaymentMethod(req.PaymentMethod))
	if err != nil {
		return payOrderError(err), nil
	}

	return &orderV1.PayOrderResponse{
		TransactionUUID: resp,
	}, nil
}

func payOrderError(err error) orderV1.PayOrderRes {
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

	default:
		return &orderV1.PayOrderInternalServerError{
			Code:    "INTERNAL_ERROR",
			Message: "internal server error",
		}
	}
}
