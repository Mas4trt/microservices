package payment

import (
	"github.com/Mas4trt/microservices/payment/internal/model"
	"github.com/google/uuid"
)

func (s *ServiceSuite) TestPayOrder() {
	tests := []struct {
		name              string
		paymentMethod     model.PaymentMethod
		wantErr           error
		wantTransactionID bool
	}{
		{
			name:              "success",
			paymentMethod:     model.PaymentMethodCard,
			wantTransactionID: true,
		},
		{
			name:          "invalid payment method",
			paymentMethod: model.PaymentMethodUnspecified,
			wantErr:       model.ErrInvalidPaymentMethod,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			req := model.PayOrderRequest{
				OrderUUID:     uuid.NewString(),
				UserUUID:      uuid.NewString(),
				PaymentMethod: tt.paymentMethod,
			}

			transactionUUID, err := s.svc.PayOrder(s.ctx, req)

			if tt.wantErr != nil {
				s.Require().ErrorIs(err, tt.wantErr)
				s.Require().Empty(transactionUUID)
				return
			}

			s.Require().NoError(err)
			s.Require().NotEmpty(transactionUUID)
			s.Require().NoError(uuid.Validate(transactionUUID))
		})
	}
}
