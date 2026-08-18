package part

import (
	"context"
	"fmt"

	"github.com/Mas4trt/microservices/inventory/internal/model"
	repoConverter "github.com/Mas4trt/microservices/inventory/internal/repository/converter"
	repoModel "github.com/Mas4trt/microservices/inventory/internal/repository/model"
)

func (r *repository) List(ctx context.Context, filter model.PartsFilter) (res []model.Part, err error) {
	query := repoConverter.FilterToBSON(filter)

	cursor, err := r.coll.Find(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to find parts: %w", err)
	}
	defer func() {
		if closeErr := cursor.Close(ctx); closeErr != nil && err == nil {
			err = fmt.Errorf("failed to close cursor: %w", closeErr)
		}
	}()

	var repoParts []repoModel.Part
	if err := cursor.All(ctx, &repoParts); err != nil {
		return nil, fmt.Errorf("failed to decode parts: %w", err)
	}

	if len(repoParts) == 0 {
		return []model.Part{}, nil
	}

	return repoConverter.PartsToModel(repoParts), nil
}
