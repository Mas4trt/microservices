package converter

import (
	"github.com/google/uuid"

	"github.com/Mas4trt/microservices/order/internal/model"
	repoModel "github.com/Mas4trt/microservices/order/internal/repository/model"
)

func ServiceToRepoModel(order model.OrderDto) repoModel.OrderDto {
	return repoModel.OrderDto{
		OrderUUID:       order.OrderUUID,
		UserUUID:        order.UserUUID,
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
		TotalPrice:      order.TotalPrice,
		TransactionUUID: order.TransactionUUID,
		PaymentMethod:   model.PaymentMethod(order.PaymentMethod),
		Status:          model.OrderStatus(order.Status),
	}
}

func ServiceToRepoOrderParts(orderUUID uuid.UUID, partUUIDs []uuid.UUID) []repoModel.OrderPart {
	result := make([]repoModel.OrderPart, 0, len(partUUIDs))
	for _, partUUID := range partUUIDs {
		result = append(result, repoModel.OrderPart{
			OrderUUID: orderUUID,
			PartUUID:  partUUID,
		})
	}
	return result
}

func RepoOrderPartsToUUIDs(parts []repoModel.OrderPart) []uuid.UUID {
	result := make([]uuid.UUID, 0, len(parts))
	for _, p := range parts {
		result = append(result, p.PartUUID)
	}
	return result
}
