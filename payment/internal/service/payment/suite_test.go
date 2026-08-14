package payment

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
)

type ServiceSuite struct {
	suite.Suite

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
