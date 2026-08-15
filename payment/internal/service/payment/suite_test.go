package payment

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
)

type ServiceSuite struct {
	suite.Suite

	//nolint:containedctx // context is intentionally stored for test suite lifecycle
	ctx context.Context
	svc *service
}

func (s *ServiceSuite) SetupTest() {
	s.ctx = context.Background()
	s.svc = NewService()
}

func TestService(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
