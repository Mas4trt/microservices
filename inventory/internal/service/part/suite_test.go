package part

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/Mas4trt/microservices/inventory/internal/repository/mocks"
)

type ServiceSuite struct {
	suite.Suite

	inventoryRepository *mocks.InventoryRepository

	srv *service
}

func (s *ServiceSuite) SetupTest() {
	s.inventoryRepository = mocks.NewInventoryRepository(s.T())

	s.srv = NewService(
		s.inventoryRepository,
	)
}

func (s *ServiceSuite) TearDownTest() {
}

func TestAPIIntegration(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
