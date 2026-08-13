package v1

import (
	"context"

	"github.com/Mas4trt/microservices/order/internal/model"
	generatedPaymentV1 "github.com/Mas4trt/microservices/shared/pkg/proto/payment/v1"
)

func (c *client) PayOrder(ctx context.Context, userUuid string, orderUuid string, paymentMethod model.PaymentMethod) (string, error) {
	protoMethod, err := PaymentMethodToProto(paymentMethod)
	if err != nil {
		return "", err
	}

	res, err := c.generatedClient.PayOrder(ctx, &generatedPaymentV1.PayOrderRequest{
		UserUuid:      userUuid,
		OrderUuid:     orderUuid,
		PaymentMethod: protoMethod,
	})
	if err != nil {
		return "", err
	}

	return res.TransactionUuid, nil
}

func PaymentMethodToProto(method model.PaymentMethod) (generatedPaymentV1.PaymentMethod, error) {
	switch method {
	case model.PaymentMethodCARD:
		return generatedPaymentV1.PaymentMethod_PAYMENT_METHOD_CARD, nil
	case model.PaymentMethodSBP:
		return generatedPaymentV1.PaymentMethod_PAYMENT_METHOD_SBP, nil
	case model.PaymentMethodCREDITCARD:
		return generatedPaymentV1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD, nil
	case model.PaymentMethodINVESTORMONEY:
		return generatedPaymentV1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY, nil
	default:
		return generatedPaymentV1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED, model.ErrInvalidPaymentMethod
	}
}
