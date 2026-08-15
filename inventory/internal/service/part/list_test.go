package part

import (
	"context"

	"github.com/brianvoe/gofakeit/v7"
)

func (s *ServiceSuite) TestListSuccess() {
	filter := newTestPartsFilter()
	expected := newTestParts(3)
	ctx := context.Background()

	s.inventoryRepository.
		On("List", ctx, filter).
		Return(expected, nil).
		Once()

	actual, err := s.srv.List(ctx, filter)

	s.Require().NoError(err)
	s.Require().Equal(expected, actual)

	s.inventoryRepository.AssertExpectations(s.T())
}

func (s *ServiceSuite) TestListRepoError() {
	filter := newTestPartsFilter()
	repoErr := gofakeit.Error()
	ctx := context.Background()

	s.inventoryRepository.
		On("List", ctx, filter).
		Return(nil, repoErr).
		Once()

	actual, err := s.srv.List(ctx, filter)

	s.Require().ErrorIs(err, repoErr)
	s.Require().Nil(actual)

	s.inventoryRepository.AssertExpectations(s.T())
}
