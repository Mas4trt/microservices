package model

import "github.com/google/uuid"

type PaymentMethod string

const (
	PaymentMethodUNKNOWN       PaymentMethod = "UNKNOWN"
	PaymentMethodCARD          PaymentMethod = "CARD"
	PaymentMethodSBP           PaymentMethod = "SBP"
	PaymentMethodCREDITCARD    PaymentMethod = "CREDIT_CARD"
	PaymentMethodINVESTORMONEY PaymentMethod = "INVESTOR_MONEY"
)

type OrderStatus string

const (
	OrderStatusPENDINGPAYMENT OrderStatus = "PENDING_PAYMENT"
	OrderStatusPAID           OrderStatus = "PAID"
	OrderStatusCANCELLED      OrderStatus = "CANCELLED"
)

type OrderDto struct {
	OrderUUID       uuid.UUID     `json:"order_uuid"`
	UserUUID        uuid.UUID     `json:"user_uuid"`
	PartUuids       []uuid.UUID   `json:"part_uuids"`
	TotalPrice      float64       `json:"total_price"`
	TransactionUUID *uuid.UUID    `json:"transaction_uuid"`
	PaymentMethod   PaymentMethod `json:"payment_method"`
	Status          OrderStatus   `json:"status"`
}
