package payment

import (
	"context"
	"log"

	"github.com/google/uuid"

	"github.com/Mas4trt/microservices/payment/internal/model"
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
