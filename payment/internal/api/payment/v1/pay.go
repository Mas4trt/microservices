package v1

import (
	"context"

	paymentV1 "github.com/Mas4trt/microservices/shared/pkg/proto/payment/v1"
)

func (a *api) PayOrder(ctx context.Context, req *paymentV1.PayOrderRequest) (*paymentV1.PayOrderResponse, error) {
	transactionUUID, err := a.paymentService.PayOrder(ctx, req.OrderUuid)
	if err != nil {
		return nil, err
	}

	return &paymentV1.PayOrderResponse{
		TransactionUuid: transactionUUID,
	}, nil
}
