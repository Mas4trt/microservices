package payment

import (
	"context"
	"log"

	"github.com/google/uuid"
)

func (s *service) PayOrder(_ context.Context, orderUUID string) (string, error) {
	transactionUUID := uuid.NewString()

	log.Printf(
		"Processed payment for order %s, transaction ID: %s",
		orderUUID,
		transactionUUID,
	)

	return transactionUUID, nil
}
