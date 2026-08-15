package v1

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/Mas4trt/microservices/order/internal/model"
	paymentV1 "github.com/Mas4trt/microservices/shared/pkg/proto/payment/v1"
)

func (c *client) PayOrder(ctx context.Context, userUuid, orderUuid string, paymentMethod model.PaymentMethod) (string, error) {
	protoMethod, err := PaymentMethodToProto(paymentMethod)
	if err != nil {
		return "", err
	}

	res, err := c.generatedClient.PayOrder(ctx, &paymentV1.PayOrderRequest{
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
		return fmt.Errorf("%w: %w", model.ErrPaymentUnavailable, err)
	}

	switch st.Code() {
	case codes.Unavailable, codes.DeadlineExceeded, codes.ResourceExhausted, codes.Canceled:
		return fmt.Errorf("%w: %w", model.ErrPaymentUnavailable, err)
	case codes.InvalidArgument:
		return fmt.Errorf("%w: %w", model.ErrInvalidPaymentMethod, err)
	default:
		return fmt.Errorf("%w: %w", model.ErrPaymentInternal, err)
	}
}

func PaymentMethodToProto(method model.PaymentMethod) (paymentV1.PaymentMethod, error) {
	switch method {
	case model.PaymentMethodCARD:
		return paymentV1.PaymentMethod_PAYMENT_METHOD_CARD, nil
	case model.PaymentMethodSBP:
		return paymentV1.PaymentMethod_PAYMENT_METHOD_SBP, nil
	case model.PaymentMethodCREDITCARD:
		return paymentV1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD, nil
	case model.PaymentMethodINVESTORMONEY:
		return paymentV1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY, nil
	default:
		return paymentV1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED, model.ErrInvalidPaymentMethod
	}
}
