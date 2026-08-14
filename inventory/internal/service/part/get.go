package part

import (
	"context"

	"github.com/Mas4trt/microservices/inventory/internal/model"
)

func (s *service) Get(ctx context.Context, partUUID string) (model.Part, error) {
	part, err := s.inventoryRepository.Get(ctx, partUUID)
	if err != nil {
		return model.Part{}, err
	}

	return part, nil
}
