package order

import (
	"errors"

	"github.com/Mas4trt/microservices/order/internal/model"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
)

func (s *ServiceSuite) TestCancelSuccess() {
	orderUUID := uuid.New()
	userUUID := uuid.New()
	partUUID := uuid.New()

	order := model.OrderDto{
		OrderUUID:     orderUUID,
		UserUUID:      userUUID,
		PartUuids:     []uuid.UUID{partUUID},
		TotalPrice:    gofakeit.Float64Range(0, 10000),
		PaymentMethod: model.PaymentMethodCARD,
		Status:        model.OrderStatusPENDINGPAYMENT,
	}

	expectedOrder := order
	expectedOrder.Status = model.OrderStatusCANCELLED

	s.orderRepository.
		On("Get", s.ctx, orderUUID).
		Return(order, nil).
		Once()

	s.orderRepository.
		On("Update", s.ctx, orderUUID, expectedOrder).
		Return(nil).
		Once()

	err := s.srv.Cancel(s.ctx, orderUUID)

	s.Require().NoError(err)

	s.orderRepository.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestCancelOrderNotFound() {
	orderUUID := uuid.New()

	s.orderRepository.
		On("Get", s.ctx, orderUUID).
		Return(model.OrderDto{}, model.ErrOrderNotFound).
		Once()

	err := s.srv.Cancel(s.ctx, orderUUID)

	s.Require().ErrorIs(err, model.ErrOrderNotFound)

	s.orderRepository.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestCancelRepositoryGetError() {
	orderUUID := uuid.New()

	expectedErr := errors.New("database unavailable")

	s.orderRepository.
		On("Get", s.ctx, orderUUID).
		Return(model.OrderDto{}, expectedErr).
		Once()

	err := s.srv.Cancel(s.ctx, orderUUID)

	s.Require().ErrorIs(err, expectedErr)

	s.orderRepository.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestCancelAlreadyPaid() {
	orderUUID := uuid.New()

	order := model.OrderDto{
		OrderUUID: orderUUID,
		Status:    model.OrderStatusPAID,
	}

	s.orderRepository.
		On("Get", s.ctx, orderUUID).
		Return(order, nil).
		Once()

	err := s.srv.Cancel(s.ctx, orderUUID)

	s.Require().ErrorIs(err, model.ErrOrderAlreadyPaid)

	s.orderRepository.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestCancelAlreadyCancelled() {
	orderUUID := uuid.New()

	order := model.OrderDto{
		OrderUUID: orderUUID,
		Status:    model.OrderStatusCANCELLED,
	}

	s.orderRepository.
		On("Get", s.ctx, orderUUID).
		Return(order, nil).
		Once()

	err := s.srv.Cancel(s.ctx, orderUUID)

	s.Require().ErrorIs(err, model.ErrOrderAlreadyCancelled)

	s.orderRepository.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestCancelInvalidStatus() {
	orderUUID := uuid.New()

	order := model.OrderDto{
		OrderUUID: orderUUID,
		Status:    model.OrderStatus("UNKNOWN_STATUS"),
	}

	s.orderRepository.
		On("Get", s.ctx, orderUUID).
		Return(order, nil).
		Once()

	err := s.srv.Cancel(s.ctx, orderUUID)

	s.Require().ErrorIs(err, model.ErrInvalidOrderStatus)

	s.orderRepository.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestCancelUpdateError() {
	orderUUID := uuid.New()

	order := model.OrderDto{
		OrderUUID: orderUUID,
		Status:    model.OrderStatusPENDINGPAYMENT,
	}

	expectedErr := errors.New("update order failed")

	s.orderRepository.
		On("Get", s.ctx, orderUUID).
		Return(order, nil).
		Once()

	expectedUpdatedOrder := order
	expectedUpdatedOrder.Status = model.OrderStatusCANCELLED

	s.orderRepository.
		On("Update", s.ctx, orderUUID, expectedUpdatedOrder).
		Return(expectedErr).
		Once()

	err := s.srv.Cancel(s.ctx, orderUUID)

	s.Require().ErrorIs(err, expectedErr)

	s.orderRepository.AssertExpectations(s.T())
}
