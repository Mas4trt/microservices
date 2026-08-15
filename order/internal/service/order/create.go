package order

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"github.com/Mas4trt/microservices/order/internal/model"
)

func (s *service) Create(ctx context.Context, userUUID uuid.UUID, partUUIDs []uuid.UUID) (uuid.UUID, float64, error) {
	if len(partUUIDs) == 0 {
		return uuid.Nil, 0, model.ErrEmptyPartUUIDs
	}

	seen := make(map[uuid.UUID]struct{}, len(partUUIDs))
	for _, id := range partUUIDs {
		if _, ok := seen[id]; ok {
			return uuid.Nil, 0, model.ErrDuplicatePartUUID
		}
		seen[id] = struct{}{}
	}

	uuids := make([]string, 0, len(partUUIDs))
	for _, id := range partUUIDs {
		uuids = append(uuids, id.String())
	}

	resp, err := s.inventoryClient.ListParts(ctx, model.PartsFilter{
		UUIDs: uuids,
	})
	if err != nil {
		return uuid.Nil, 0, err
	}

	priceByUUID := make(map[uuid.UUID]float64, len(resp))
	for _, part := range resp {
		priceByUUID[part.UUID] = part.Price
	}

	var totalPrice float64
	for _, id := range partUUIDs {
		price, ok := priceByUUID[id]
		if !ok {
			return uuid.Nil, 0, model.ErrPartNotFound
		}
		totalPrice += price
	}

	newOrder := &model.OrderDto{
		UserUUID:      userUUID,
		PartUuids:     partUUIDs,
		TotalPrice:    totalPrice,
		PaymentMethod: model.PaymentMethodUNKNOWN,
		Status:        model.OrderStatusPENDINGPAYMENT,
	}

	orderUUID, err := s.orderRepository.Create(ctx, *newOrder)
	if err != nil {
		return uuid.Nil, 0, fmt.Errorf("%w: %w", model.ErrOrderCreateFailed, err)
	}

	return orderUUID, totalPrice, nil
}
