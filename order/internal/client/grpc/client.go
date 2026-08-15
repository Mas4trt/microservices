package grpc

import (
	"context"

	"github.com/Mas4trt/microservices/order/internal/model"
)

type InventoryClient interface {
	ListParts(ctx context.Context, filter model.PartsFilter) ([]model.PricedPart, error)
}

type PaymentClient interface {
	PayOrder(ctx context.Context, userUuid, orderUuid string, paymentMethod model.PaymentMethod) (string, error)
}
