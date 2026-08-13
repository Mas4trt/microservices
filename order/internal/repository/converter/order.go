package converter

import (
	"github.com/Mas4trt/microservices/order/internal/model"
	repoModel "github.com/Mas4trt/microservices/order/internal/repository/model"
)

func ServiceToRepoModel(order model.OrderDto) repoModel.OrderDto {
	return repoModel.OrderDto{
		OrderUUID:       order.OrderUUID,
		UserUUID:        order.UserUUID,
		PartUuids:       order.PartUuids,
		TotalPrice:      order.TotalPrice,
		TransactionUUID: order.TransactionUUID,
		PaymentMethod:   repoModel.PaymentMethod(order.PaymentMethod),
		Status:          repoModel.OrderStatus(order.Status),
	}
}

func RepoToServiceModel(order repoModel.OrderDto) model.OrderDto {
	return model.OrderDto{
		OrderUUID:       order.OrderUUID,
		UserUUID:        order.UserUUID,
		PartUuids:       order.PartUuids,
		TotalPrice:      order.TotalPrice,
		TransactionUUID: order.TransactionUUID,
		PaymentMethod:   model.PaymentMethod(order.PaymentMethod),
		Status:          model.OrderStatus(order.Status),
	}
}
