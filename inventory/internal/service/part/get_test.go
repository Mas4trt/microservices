package part

import (
	"context"

	"github.com/brianvoe/gofakeit/v7"

	"github.com/Mas4trt/microservices/inventory/internal/model"
)

func (s *ServiceSuite) TestGetSuccess() {
	partUUID := gofakeit.UUID()
	expected := newTestPart()
	ctx := context.Background()

	s.inventoryRepository.
		On("Get", ctx, partUUID).
		Return(expected, nil).
		Once()

	actual, err := s.srv.Get(ctx, partUUID)

	s.Require().NoError(err)
	s.Require().Equal(expected, actual)

	s.inventoryRepository.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestGetRepoError() {
	partUUID := gofakeit.UUID()
	repoErr := gofakeit.Error()
	ctx := context.Background()

	s.inventoryRepository.
		On("Get", ctx, partUUID).
		Return(model.Part{}, repoErr).
		Once()

	actual, err := s.srv.Get(ctx, partUUID)

	s.Require().ErrorIs(err, repoErr)
	s.Require().Empty(actual)

	s.inventoryRepository.AssertExpectations(s.T())
}
