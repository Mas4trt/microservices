package order

import (
	"errors"
	"fmt"

	"github.com/Mas4trt/microservices/order/internal/model"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
)

func (s *ServiceSuite) TestGetSuccess() {
	orderUUID := uuid.New()
	userUUID := uuid.New()
	partUUID1 := uuid.New()
	partUUID2 := uuid.New()
	transactionUUID := uuid.New()

	expected := model.OrderDto{
		OrderUUID: orderUUID,
		UserUUID:  userUUID,
		PartUuids: []uuid.UUID{
			partUUID1,
			partUUID2,
		},
		TotalPrice:      gofakeit.Float64Range(0, 10000),
		TransactionUUID: &transactionUUID,
		PaymentMethod:   model.PaymentMethodCARD,
		Status:          model.OrderStatusPAID,
	}

	s.orderRepository.
		On("Get", s.ctx, orderUUID).
		Return(expected, nil).
		Once()

	actual, err := s.srv.Get(s.ctx, orderUUID)

	s.Require().NoError(err)
	s.Require().Equal(expected, actual)

	s.orderRepository.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestServiceGetOrderNotFound() {
	orderUUID := uuid.New()

	wrappedErr := fmt.Errorf("get order: %w", model.ErrOrderNotFound)

	s.orderRepository.
		On("Get", s.ctx, orderUUID).
		Return(model.OrderDto{}, wrappedErr).
		Once()

	actual, err := s.srv.Get(s.ctx, orderUUID)

	s.Require().ErrorIs(err, model.ErrOrderNotFound)
	s.Require().Equal(model.OrderDto{}, actual)

	s.orderRepository.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestGetRepositoryError() {
	orderUUID := uuid.New()

	expectedErr := errors.New("database unavailable")

	s.orderRepository.
		On("Get", s.ctx, orderUUID).
		Return(model.OrderDto{}, expectedErr).
		Once()

	actual, err := s.srv.Get(s.ctx, orderUUID)

	s.Require().ErrorIs(err, expectedErr)
	s.Require().Equal(model.OrderDto{}, actual)

	s.orderRepository.AssertExpectations(s.T())
}
