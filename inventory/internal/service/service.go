package service

import (
	"context"

	"github.com/Mas4trt/microservices/inventory/internal/model"
)

type InventoryService interface {
	Get(ctx context.Context, partUUID string) (model.Part, error)
	List(ctx context.Context, filter model.PartsFilter) ([]model.Part, error)
}
