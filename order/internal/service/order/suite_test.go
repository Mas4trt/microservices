package order

import (
	"context"
	"testing"

	clientMocks "github.com/Mas4trt/microservices/order/internal/client/grpc/mocks"
	repoMocks "github.com/Mas4trt/microservices/order/internal/repository/mocks"
	"github.com/stretchr/testify/suite"
)

type ServiceSuite struct {
	suite.Suite

	ctx context.Context

	orderRepository *repoMocks.OrderRepository
	inventoryClient *clientMocks.InventoryClient
	paymentClient   *clientMocks.PaymentClient

	srv *service
}

func (s *ServiceSuite) SetupTest() {
	s.ctx = context.Background()

	s.orderRepository = repoMocks.NewOrderRepository(s.T())
	s.inventoryClient = clientMocks.NewInventoryClient(s.T())
	s.paymentClient = clientMocks.NewPaymentClient(s.T())

	s.srv = NewService(
		s.orderRepository,
		s.inventoryClient,
		s.paymentClient,
	)
}

func (s *ServiceSuite) TearDownTest() {
}

func TestAPIIntegration(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
