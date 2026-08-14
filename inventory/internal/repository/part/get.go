package part

import (
	"context"

	"github.com/Mas4trt/microservices/inventory/internal/model"
	repoConverter "github.com/Mas4trt/microservices/inventory/internal/repository/converter"
)

func (r *repository) Get(_ context.Context, partUUID string) (model.Part, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	part, ok := r.parts[partUUID]
	if !ok {
		return model.Part{}, model.ErrPartNotFound
	}

	return repoConverter.PartToModel(part), nil
}
