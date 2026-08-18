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
	OrderUUID       uuid.UUID     `db:"order_uuid"`
	UserUUID        uuid.UUID     `db:"user_uuid"`
	TotalPrice      float64       `db:"total_price"`
	TransactionUUID *uuid.UUID    `db:"transaction_uuid"`
	PaymentMethod   PaymentMethod `db:"payment_method"`
	Status          OrderStatus   `db:"status"`
}

type OrderPart struct {
	OrderUUID uuid.UUID `db:"order_uuid"`
	PartUUID  uuid.UUID `db:"part_uuid"`
}
