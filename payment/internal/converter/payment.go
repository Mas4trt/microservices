package converter

import (
	"github.com/Mas4trt/microservices/payment/internal/model"
	paymentV1 "github.com/Mas4trt/microservices/shared/pkg/proto/payment/v1"
)

func PayOrderRequestToModel(req *paymentV1.PayOrderRequest) model.PayOrderRequest {
	return model.PayOrderRequest{
		OrderUUID:     req.GetOrderUuid(),
		UserUUID:      req.GetUserUuid(),
		PaymentMethod: PaymentMethodToModel(req.GetPaymentMethod()),
	}
}

func PaymentMethodToModel(method paymentV1.PaymentMethod) model.PaymentMethod {
	switch method {
	case paymentV1.PaymentMethod_PAYMENT_METHOD_CARD:
		return model.PaymentMethodCard
	case paymentV1.PaymentMethod_PAYMENT_METHOD_SBP:
		return model.PaymentMethodSBP
	case paymentV1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD:
		return model.PaymentMethodCreditCard
	case paymentV1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY:
		return model.PaymentMethodInvestorMoney
	default:
		return model.PaymentMethodUnspecified
	}
}
