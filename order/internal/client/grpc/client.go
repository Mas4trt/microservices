package grpc

import (
	"context"

	"github.com/Mas4trt/microservices/order/internal/model"
)

type InventoryClient interface {
	ListParts(ctx context.Context, filter model.PartsFilter) ([]model.Part, error)
}

type PaymentClient interface {
	PayOrder(ctx context.Context, userUuid string, orderUuid string, paymentMethod model.PaymentMethod) (string, error)
}
