package order

import (
	"errors"

	"github.com/Mas4trt/microservices/order/internal/model"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
)

func (s *ServiceSuite) TestCreateEmptyPartUUIDs() {
	userUUID := uuid.New()

	orderUUID, totalPrice, err := s.srv.Create(
		s.ctx,
		userUUID,
		[]uuid.UUID{},
	)

	s.Require().ErrorIs(err, model.ErrEmptyPartUUIDs)
	s.Require().Equal(uuid.Nil, orderUUID)
	s.Require().Equal(float64(0), totalPrice)

	s.orderRepository.AssertNotCalled(s.T(), "Create")
	s.inventoryClient.AssertNotCalled(s.T(), "ListParts")
}

func (s *ServiceSuite) TestCreateSuccess() {
	userUUID := uuid.New()
	partUUID := uuid.New()
	expectedOrderUUID := uuid.New()
	price := gofakeit.Float64Range(0, 10000)

	resp := []model.PricedPart{
		{
			UUID:  partUUID,
			Price: price,
		},
	}

	s.inventoryClient.
		On("ListParts", s.ctx, model.PartsFilter{
			UUIDs: []string{partUUID.String()},
		}).
		Return(resp, nil).
		Once()

	expectedOrder := model.OrderDto{
		UserUUID:      userUUID,
		PartUuids:     []uuid.UUID{partUUID},
		TotalPrice:    price,
		PaymentMethod: model.PaymentMethodUNKNOWN,
		Status:        model.OrderStatusPENDINGPAYMENT,
	}

	s.orderRepository.
		On("Create", s.ctx, expectedOrder).
		Return(expectedOrderUUID, nil).
		Once()

	orderUUID, totalPrice, err := s.srv.Create(
		s.ctx,
		userUUID,
		[]uuid.UUID{partUUID},
	)

	s.Require().NoError(err)
	s.Require().Equal(expectedOrderUUID, orderUUID)
	s.Require().Equal(price, totalPrice)

	s.inventoryClient.AssertExpectations(s.T())
	s.orderRepository.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestCreateDuplicatePartUUID() {
	userUUID := uuid.New()
	partUUID := uuid.New()

	partUUIDs := []uuid.UUID{
		partUUID,
		partUUID,
	}

	orderUUID, totalPrice, err := s.srv.Create(
		s.ctx,
		userUUID,
		partUUIDs,
	)

	s.Require().ErrorIs(err, model.ErrDuplicatePartUUID)
	s.Require().Equal(uuid.Nil, orderUUID)
	s.Require().Equal(float64(0), totalPrice)

	s.inventoryClient.AssertNotCalled(s.T(), "ListParts")
	s.orderRepository.AssertNotCalled(s.T(), "Create")
}

func (s *ServiceSuite) TestCreateListPartsError() {
	userUUID := uuid.New()
	partUUID := uuid.New()

	expectedErr := errors.New("inventory service unavailable")

	s.inventoryClient.
		On("ListParts", s.ctx, model.PartsFilter{
			UUIDs: []string{partUUID.String()},
		}).
		Return(nil, expectedErr).
		Once()

	orderUUID, totalPrice, err := s.srv.Create(
		s.ctx,
		userUUID,
		[]uuid.UUID{partUUID},
	)

	s.Require().ErrorIs(err, expectedErr)
	s.Require().Equal(uuid.Nil, orderUUID)
	s.Require().Equal(float64(0), totalPrice)

	s.orderRepository.AssertNotCalled(s.T(), "Create")
	s.inventoryClient.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestCreatePartNotFound() {
	userUUID := uuid.New()
	partUUID1 := uuid.New()
	partUUID2 := uuid.New()

	resp := []model.PricedPart{
		{
			UUID:  partUUID1,
			Price: gofakeit.Float64Range(0, 10000),
		},
	}

	s.inventoryClient.
		On("ListParts", s.ctx, model.PartsFilter{
			UUIDs: []string{
				partUUID1.String(),
				partUUID2.String(),
			},
		}).
		Return(resp, nil).
		Once()

	orderUUID, totalPrice, err := s.srv.Create(
		s.ctx,
		userUUID,
		[]uuid.UUID{
			partUUID1,
			partUUID2,
		},
	)

	s.Require().ErrorIs(err, model.ErrPartNotFound)
	s.Require().Equal(uuid.Nil, orderUUID)
	s.Require().Equal(float64(0), totalPrice)

	s.orderRepository.AssertNotCalled(s.T(), "Create")
	s.inventoryClient.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestCreateMultipleParts() {
	userUUID := uuid.New()
	partUUID1 := uuid.New()
	partUUID2 := uuid.New()
	partUUID3 := uuid.New()
	expectedOrderUUID := uuid.New()
	price1 := gofakeit.Float64Range(0, 10000)
	price2 := gofakeit.Float64Range(0, 10000)
	price3 := gofakeit.Float64Range(0, 10000)
	totalPriceTest := price1 + price2 + price3

	resp := []model.PricedPart{
		{
			UUID:  partUUID1,
			Price: price1,
		},
		{
			UUID:  partUUID2,
			Price: price2,
		},
		{
			UUID:  partUUID3,
			Price: price3,
		},
	}

	partUUIDs := []uuid.UUID{
		partUUID1,
		partUUID2,
		partUUID3,
	}

	s.inventoryClient.
		On("ListParts", s.ctx, model.PartsFilter{
			UUIDs: []string{
				partUUID1.String(),
				partUUID2.String(),
				partUUID3.String(),
			},
		}).
		Return(resp, nil).
		Once()

	expectedOrder := model.OrderDto{
		UserUUID:   userUUID,
		PartUuids:  partUUIDs,
		TotalPrice: totalPriceTest,

		PaymentMethod: model.PaymentMethodUNKNOWN,
		Status:        model.OrderStatusPENDINGPAYMENT,
	}

	s.orderRepository.
		On("Create", s.ctx, expectedOrder).
		Return(expectedOrderUUID, nil).
		Once()

	orderUUID, totalPrice, err := s.srv.Create(
		s.ctx,
		userUUID,
		partUUIDs,
	)

	s.Require().NoError(err)
	s.Require().Equal(expectedOrderUUID, orderUUID)
	s.Require().Equal(totalPriceTest, totalPrice)

	s.inventoryClient.AssertExpectations(s.T())
	s.orderRepository.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestCreateOrderRepositoryError() {
	userUUID := uuid.New()
	partUUID := uuid.New()
	price := gofakeit.Float64Range(0, 10000)

	createErr := errors.New("database unavailable")

	resp := []model.PricedPart{
		{
			UUID:  partUUID,
			Price: price,
		},
	}

	s.inventoryClient.
		On("ListParts", s.ctx, model.PartsFilter{
			UUIDs: []string{partUUID.String()},
		}).
		Return(resp, nil).
		Once()

	expectedOrder := model.OrderDto{
		UserUUID:      userUUID,
		PartUuids:     []uuid.UUID{partUUID},
		TotalPrice:    price,
		PaymentMethod: model.PaymentMethodUNKNOWN,
		Status:        model.OrderStatusPENDINGPAYMENT,
	}

	s.orderRepository.
		On("Create", s.ctx, expectedOrder).
		Return(uuid.Nil, createErr).
		Once()

	orderUUID, totalPrice, err := s.srv.Create(
		s.ctx,
		userUUID,
		[]uuid.UUID{partUUID},
	)

	s.Require().ErrorIs(err, model.ErrOrderCreateFailed)
	s.Require().Equal(uuid.Nil, orderUUID)
	s.Require().Equal(float64(0), totalPrice)

	s.inventoryClient.AssertExpectations(s.T())
	s.orderRepository.AssertExpectations(s.T())
}
