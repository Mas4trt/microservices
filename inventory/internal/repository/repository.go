package repository

import (
	"context"

	"github.com/Mas4trt/microservices/inventory/internal/model"
)

type InventoryRepository interface {
	Get(ctx context.Context, partUUID string) (model.Part, error)
	List(ctx context.Context, filter model.PartsFilter) ([]model.Part, error)
	Init(ctx context.Context, parts map[string]model.Part) error
}
