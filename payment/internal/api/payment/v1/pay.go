package v1

import (
	"context"
	"errors"

	"github.com/Mas4trt/microservices/payment/internal/converter"
	"github.com/Mas4trt/microservices/payment/internal/model"
	paymentV1 "github.com/Mas4trt/microservices/shared/pkg/proto/payment/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a *api) PayOrder(ctx context.Context, req *paymentV1.PayOrderRequest) (*paymentV1.PayOrderResponse, error) {
	transactionUUID, err := a.paymentService.PayOrder(ctx, converter.PayOrderRequestToModel(req))
	if err != nil {
		if errors.Is(err, model.ErrInvalidPaymentMethod) {
			return nil, status.Error(codes.InvalidArgument, "invalid payment method")
		}
		return nil, status.Error(codes.Internal, "internal server error")
	}

	return &paymentV1.PayOrderResponse{
		TransactionUuid: transactionUUID,
	}, nil
}
