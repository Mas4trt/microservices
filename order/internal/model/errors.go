package model

import "errors"

var (
	ErrOrderNotFound         = errors.New("order not found")
	ErrPartNotFound          = errors.New("part not found")
	ErrOrderAlreadyPaid      = errors.New("order is already paid")
	ErrOrderAlreadyCancelled = errors.New("order is already cancelled")
	ErrInvalidOrderStatus    = errors.New("order cannot be processed in current status")
	ErrInvalidPaymentMethod  = errors.New("invalid payment method")
	ErrInventoryUnavailable  = errors.New("inventory service is unavailable")
	ErrPaymentUnavailable    = errors.New("payment service is unavailable")
	ErrEmptyPartUUIDs        = errors.New("part uuids list is empty")
	ErrDuplicatePartUUID     = errors.New("duplicate part uuid in request")
	ErrOrderCreateFailed     = errors.New("failed to create order")
)
