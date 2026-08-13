package v1

import (
	"context"

	clientConverter "github.com/Mas4trt/microservices/order/internal/client/converter"
	"github.com/Mas4trt/microservices/order/internal/model"
	generatedInventoryV1 "github.com/Mas4trt/microservices/shared/pkg/proto/inventory/v1"
)

func (c *client) ListParts(ctx context.Context, filter model.PartsFilter) ([]model.Part, error) {
	parts, err := c.generatedClient.ListParts(ctx, &generatedInventoryV1.ListPartsRequest{
		Filter: clientConverter.PartsFilterToProto(filter),
	})
	if err != nil {
		return nil, err
	}

	return clientConverter.PartListToModel(parts.Parts), nil
}
