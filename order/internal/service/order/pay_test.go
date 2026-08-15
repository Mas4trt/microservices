package order

import (
	"errors"

	"github.com/Mas4trt/microservices/order/internal/model"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
)

func (s *ServiceSuite) TestPaySuccess() {
	orderUUID := uuid.New()
	userUUID := uuid.New()
	transactionUUID := uuid.New()

	method := model.PaymentMethodCARD

	order := model.OrderDto{
		OrderUUID:     orderUUID,
		UserUUID:      userUUID,
		PartUuids:     []uuid.UUID{uuid.New()},
		TotalPrice:    gofakeit.Float64Range(0, 10000),
		PaymentMethod: model.PaymentMethodUNKNOWN,
		Status:        model.OrderStatusPENDINGPAYMENT,
	}

	s.orderRepository.
		On("Get", s.ctx, orderUUID).
		Return(order, nil).
		Once()

	s.paymentClient.
		On(
			"PayOrder",
			s.ctx,
			userUUID.String(),
			orderUUID.String(),
			method,
		).
		Return(transactionUUID.String(), nil).
		Once()

	expectedOrder := order
	expectedOrder.TransactionUUID = &transactionUUID
	expectedOrder.PaymentMethod = method
	expectedOrder.Status = model.OrderStatusPAID

	s.orderRepository.
		On("Update", s.ctx, orderUUID, expectedOrder).
		Return(nil).
		Once()

	actualTransactionUUID, err := s.srv.Pay(
		s.ctx,
		orderUUID,
		method,
	)

	s.Require().NoError(err)
	s.Require().Equal(transactionUUID, actualTransactionUUID)

	s.orderRepository.AssertExpectations(s.T())
	s.paymentClient.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestPayOrderNotFound() {
	orderUUID := uuid.New()

	s.orderRepository.
		On("Get", s.ctx, orderUUID).
		Return(model.OrderDto{}, model.ErrOrderNotFound).
		Once()

	actualTransactionUUID, err := s.srv.Pay(
		s.ctx,
		orderUUID,
		model.PaymentMethodCARD,
	)

	s.Require().ErrorIs(err, model.ErrOrderNotFound)
	s.Require().Equal(uuid.Nil, actualTransactionUUID)

	s.orderRepository.AssertExpectations(s.T())
	s.paymentClient.AssertNotCalled(s.T(), "PayOrder")
}

func (s *ServiceSuite) TestPayRepositoryGetError() {
	orderUUID := uuid.New()

	expectedErr := errors.New("database unavailable")

	s.orderRepository.
		On("Get", s.ctx, orderUUID).
		Return(model.OrderDto{}, expectedErr).
		Once()

	actualTransactionUUID, err := s.srv.Pay(
		s.ctx,
		orderUUID,
		model.PaymentMethodCARD,
	)

	s.Require().ErrorIs(err, expectedErr)
	s.Require().Equal(uuid.Nil, actualTransactionUUID)

	s.orderRepository.AssertExpectations(s.T())
	s.paymentClient.AssertNotCalled(s.T(), "PayOrder")
}

func (s *ServiceSuite) TestPayAlreadyPaid() {
	orderUUID := uuid.New()
	userUUID := uuid.New()

	order := model.OrderDto{
		OrderUUID: orderUUID,
		UserUUID:  userUUID,
		Status:    model.OrderStatusPAID,
	}

	s.orderRepository.
		On("Get", s.ctx, orderUUID).
		Return(order, nil).
		Once()

	actualTransactionUUID, err := s.srv.Pay(
		s.ctx,
		orderUUID,
		model.PaymentMethodCARD,
	)

	s.Require().ErrorIs(err, model.ErrOrderAlreadyPaid)
	s.Require().Equal(uuid.Nil, actualTransactionUUID)

	s.paymentClient.AssertNotCalled(s.T(), "PayOrder")
	s.orderRepository.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestPayAlreadyCancelled() {
	orderUUID := uuid.New()
	userUUID := uuid.New()

	order := model.OrderDto{
		OrderUUID: orderUUID,
		UserUUID:  userUUID,
		Status:    model.OrderStatusCANCELLED,
	}

	s.orderRepository.
		On("Get", s.ctx, orderUUID).
		Return(order, nil).
		Once()

	actualTransactionUUID, err := s.srv.Pay(
		s.ctx,
		orderUUID,
		model.PaymentMethodCARD,
	)

	s.Require().ErrorIs(err, model.ErrOrderAlreadyCancelled)
	s.Require().Equal(uuid.Nil, actualTransactionUUID)

	s.paymentClient.AssertNotCalled(s.T(), "PayOrder")
	s.orderRepository.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestPayInvalidOrderStatus() {
	orderUUID := uuid.New()
	userUUID := uuid.New()

	order := model.OrderDto{
		OrderUUID: orderUUID,
		UserUUID:  userUUID,
		Status:    model.OrderStatus("UNKNOWN_STATUS"),
	}

	s.orderRepository.
		On("Get", s.ctx, orderUUID).
		Return(order, nil).
		Once()

	actualTransactionUUID, err := s.srv.Pay(
		s.ctx,
		orderUUID,
		model.PaymentMethodCARD,
	)

	s.Require().ErrorIs(err, model.ErrInvalidOrderStatus)
	s.Require().Equal(uuid.Nil, actualTransactionUUID)

	s.paymentClient.AssertNotCalled(s.T(), "PayOrder")
	s.orderRepository.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestPayPaymentClientError() {
	orderUUID := uuid.New()
	userUUID := uuid.New()

	expectedErr := errors.New("payment service unavailable")

	order := model.OrderDto{
		OrderUUID: orderUUID,
		UserUUID:  userUUID,
		Status:    model.OrderStatusPENDINGPAYMENT,
	}

	method := model.PaymentMethodCARD

	s.orderRepository.
		On("Get", s.ctx, orderUUID).
		Return(order, nil).
		Once()

	s.paymentClient.
		On(
			"PayOrder",
			s.ctx,
			userUUID.String(),
			orderUUID.String(),
			method,
		).
		Return("", expectedErr).
		Once()

	actualTransactionUUID, err := s.srv.Pay(
		s.ctx,
		orderUUID,
		method,
	)

	s.Require().ErrorIs(err, expectedErr)
	s.Require().Equal(uuid.Nil, actualTransactionUUID)

	s.orderRepository.AssertNotCalled(s.T(), "Update")

	s.orderRepository.AssertExpectations(s.T())
	s.paymentClient.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestPayInvalidTransactionUUID() {
	orderUUID := uuid.New()
	userUUID := uuid.New()

	invalidTransactionUUID := "not-a-valid-uuid"

	order := model.OrderDto{
		OrderUUID: orderUUID,
		UserUUID:  userUUID,
		Status:    model.OrderStatusPENDINGPAYMENT,
	}

	method := model.PaymentMethodCARD

	s.orderRepository.
		On("Get", s.ctx, orderUUID).
		Return(order, nil).
		Once()

	s.paymentClient.
		On(
			"PayOrder",
			s.ctx,
			userUUID.String(),
			orderUUID.String(),
			method,
		).
		Return(invalidTransactionUUID, nil).
		Once()

	actualTransactionUUID, err := s.srv.Pay(
		s.ctx,
		orderUUID,
		method,
	)

	s.Require().Error(err)
	s.Require().Equal(uuid.Nil, actualTransactionUUID)

	s.orderRepository.AssertNotCalled(s.T(), "Update")

	s.orderRepository.AssertExpectations(s.T())
	s.paymentClient.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestPayUpdateError() {
	orderUUID := uuid.New()
	userUUID := uuid.New()
	transactionUUID := uuid.New()

	method := model.PaymentMethodCARD

	order := model.OrderDto{
		OrderUUID:     orderUUID,
		UserUUID:      userUUID,
		Status:        model.OrderStatusPENDINGPAYMENT,
		PaymentMethod: model.PaymentMethodUNKNOWN,
	}

	s.orderRepository.
		On("Get", s.ctx, orderUUID).
		Return(order, nil).
		Once()

	s.paymentClient.
		On(
			"PayOrder",
			s.ctx,
			userUUID.String(),
			orderUUID.String(),
			method,
		).
		Return(transactionUUID.String(), nil).
		Once()

	expectedOrder := order
	expectedOrder.TransactionUUID = &transactionUUID
	expectedOrder.PaymentMethod = method
	expectedOrder.Status = model.OrderStatusPAID

	expectedErr := errors.New("database unavailable")

	s.orderRepository.
		On("Update", s.ctx, orderUUID, expectedOrder).
		Return(expectedErr).
		Once()

	actualTransactionUUID, err := s.srv.Pay(
		s.ctx,
		orderUUID,
		method,
	)

	s.Require().ErrorIs(err, expectedErr)
	s.Require().Equal(uuid.Nil, actualTransactionUUID)

	s.orderRepository.AssertExpectations(s.T())
	s.paymentClient.AssertExpectations(s.T())
}
