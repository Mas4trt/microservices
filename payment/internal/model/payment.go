package model

// PaymentMethod — доменное представление способа оплаты,
// независимое от транспортного (proto) типа.
type PaymentMethod int32

const (
	PaymentMethodUnspecified PaymentMethod = iota
	PaymentMethodCard
	PaymentMethodSBP
	PaymentMethodCreditCard
	PaymentMethodInvestorMoney
)

// PayOrderRequest — доменная модель команды на оплату.
type PayOrderRequest struct {
	OrderUUID     string
	UserUUID      string
	PaymentMethod PaymentMethod
}
