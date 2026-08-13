package converter

import (
	"github.com/Mas4trt/microservices/order/internal/model"
	orderV1 "github.com/Mas4trt/microservices/shared/pkg/openapi/order/v1"
)

func ToProtoOrder(order model.OrderDto) *orderV1.OrderDto {
	dto := &orderV1.OrderDto{
		OrderUUID:  order.OrderUUID,
		UserUUID:   order.UserUUID,
		PartUuids:  order.PartUuids,
		TotalPrice: order.TotalPrice,
		Status:     OrderStatusToProto(order.Status),
	}

	if order.TransactionUUID != nil {
		dto.TransactionUUID = orderV1.NewOptNilUUID(*order.TransactionUUID)
	} else {
		dto.TransactionUUID.SetToNull()
	}

	if order.PaymentMethod != "" && order.PaymentMethod != model.PaymentMethodUNKNOWN {
		dto.PaymentMethod = orderV1.NewOptNilPaymentMethod(PaymentMethodToProtoDto(order.PaymentMethod))
	} else {
		dto.PaymentMethod.SetToNull()
	}

	return dto
}

func OrderStatusToProto(status model.OrderStatus) orderV1.OrderStatus {
	switch status {
	case model.OrderStatusPENDINGPAYMENT:
		return orderV1.OrderStatusPENDINGPAYMENT
	case model.OrderStatusPAID:
		return orderV1.OrderStatusPAID
	case model.OrderStatusCANCELLED:
		return orderV1.OrderStatusCANCELLED
	default:
		return orderV1.OrderStatusPENDINGPAYMENT
	}
}

func PaymentMethodToProtoDto(method model.PaymentMethod) orderV1.PaymentMethod {
	switch method {
	case model.PaymentMethodCARD:
		return orderV1.PaymentMethodCARD
	case model.PaymentMethodSBP:
		return orderV1.PaymentMethodSBP
	case model.PaymentMethodCREDITCARD:
		return orderV1.PaymentMethodCREDITCARD
	case model.PaymentMethodINVESTORMONEY:
		return orderV1.PaymentMethodINVESTORMONEY
	default:
		return orderV1.PaymentMethodUNKNOWN
	}
}
