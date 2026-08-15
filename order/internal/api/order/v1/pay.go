package v1

import (
	"context"

	"github.com/Mas4trt/microservices/order/internal/model"
	orderV1 "github.com/Mas4trt/microservices/shared/pkg/openapi/order/v1"
)

func (a *api) PayOrder(ctx context.Context, req *orderV1.PayOrderRequest, params orderV1.PayOrderParams) (orderV1.PayOrderRes, error) {
	resp, err := a.orderService.Pay(ctx, params.OrderUUID, model.PaymentMethod(req.PaymentMethod))
	if err != nil {
		return HandlePayOrderError(err), nil
	}

	return &orderV1.PayOrderResponse{
		TransactionUUID: resp,
	}, nil
}
