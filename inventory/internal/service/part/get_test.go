package part

import (
	"github.com/Mas4trt/microservices/inventory/internal/model"
	"github.com/brianvoe/gofakeit/v7"
)

func (s *ServiceSuite) TestGetSuccess() {
	partUUID := gofakeit.UUID()
	expected := newTestPart()

	s.inventoryRepository.
		On("Get", s.ctx, partUUID).
		Return(expected, nil).
		Once()

	actual, err := s.srv.Get(s.ctx, partUUID)

	s.Require().NoError(err)
	s.Require().Equal(expected, actual)

	s.inventoryRepository.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestGetRepoError() {
	partUUID := gofakeit.UUID()
	repoErr := gofakeit.Error()

	s.inventoryRepository.
		On("Get", s.ctx, partUUID).
		Return(model.Part{}, repoErr).
		Once()

	actual, err := s.srv.Get(s.ctx, partUUID)

	s.Require().ErrorIs(err, repoErr)
	s.Require().Empty(actual)

	s.inventoryRepository.AssertExpectations(s.T())
}
