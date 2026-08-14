package payment

import (
	"context"
	"log"

	"github.com/Mas4trt/microservices/payment/internal/model"
	"github.com/google/uuid"
)

func (s *service) PayOrder(_ context.Context, req model.PayOrderRequest) (string, error) {
	if req.PaymentMethod == model.PaymentMethodUnspecified {
		return "", model.ErrInvalidPaymentMethod
	}

	transactionUUID := uuid.NewString()

	log.Printf(
		"Processed payment for order %s, transaction ID: %s",
		req.OrderUUID,
		transactionUUID,
	)

	return transactionUUID, nil
}
