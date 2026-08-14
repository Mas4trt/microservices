package service

import (
	"context"

	"github.com/Mas4trt/microservices/payment/internal/model"
)

type PaymentService interface {
	PayOrder(ctx context.Context, req model.PayOrderRequest) (string, error)
}
