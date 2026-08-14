package v1

import (
	"context"
	"fmt"

	"github.com/Mas4trt/microservices/order/internal/model"
	generatedPaymentV1 "github.com/Mas4trt/microservices/shared/pkg/proto/payment/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
		return "", mapPaymentErr(err)
	}

	return res.TransactionUuid, nil
}

func mapPaymentErr(err error) error {
	st, ok := status.FromError(err)
	if !ok {
		return fmt.Errorf("%w: %v", model.ErrPaymentUnavailable, err)
	}

	switch st.Code() {
	case codes.Unavailable, codes.DeadlineExceeded, codes.ResourceExhausted, codes.Canceled:
		return fmt.Errorf("%w: %v", model.ErrPaymentUnavailable, err)
	case codes.InvalidArgument:
		return fmt.Errorf("%w: %v", model.ErrInvalidPaymentMethod, err)
	default:
		return fmt.Errorf("%w: %v", model.ErrPaymentInternal, err)
	}
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
