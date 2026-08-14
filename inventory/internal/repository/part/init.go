package part

import (
	"context"

	"github.com/Mas4trt/microservices/inventory/internal/model"
	repoConverter "github.com/Mas4trt/microservices/inventory/internal/repository/converter"
)

func (r *repository) Init(_ context.Context, parts map[string]model.Part) error {
	convParts := repoConverter.PartsToRepoModel(parts)

	r.parts = convParts

	return nil
}
