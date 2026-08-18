package part

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/Mas4trt/microservices/inventory/internal/repository/mocks"
)

type ServiceSuite struct {
	suite.Suite

	//nolint:containedctx // context is intentionally stored for test suite lifecycle
	ctx context.Context

	inventoryRepository *mocks.InventoryRepository

	srv *service
}

func (s *ServiceSuite) SetupTest() {
	s.ctx = context.Background()

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
