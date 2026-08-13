package v1

import (
	"github.com/Mas4trt/microservices/order/internal/service"
)

type api struct {
	orderService service.OrderService
}

func NewOrderHandler(orderService service.OrderService) *api {
	return &api{
		orderService: orderService,
	}
}
