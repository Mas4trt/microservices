package v1

import (
	"context"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	clientConverter "github.com/Mas4trt/microservices/order/internal/client/converter"
	"github.com/Mas4trt/microservices/order/internal/model"
	inventoryV1 "github.com/Mas4trt/microservices/shared/pkg/proto/inventory/v1"
)

func (c *client) ListParts(ctx context.Context, filter model.PartsFilter) ([]model.PricedPart, error) {
	parts, err := c.generatedClient.ListParts(ctx, &inventoryV1.ListPartsRequest{
		Filter: clientConverter.PartsFilterToProto(filter),
	})
	if err != nil {
		return nil, mapInventoryErr(err)
	}

	res, err := clientConverter.PartListToModel(parts.Parts)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", model.ErrInventoryInternal, err)
	}

	return res, nil
}

func mapInventoryErr(err error) error {
	st, ok := status.FromError(err)
	if !ok {
		return fmt.Errorf("%w: %w", model.ErrInventoryUnavailable, err)
	}

	switch st.Code() {
	case codes.Unavailable, codes.DeadlineExceeded, codes.ResourceExhausted, codes.Canceled:
		return fmt.Errorf("%w: %w", model.ErrInventoryUnavailable, err)
	case codes.InvalidArgument:
		return fmt.Errorf("%w: %w", model.ErrInventoryInvalidArgument, err)
	default:
		return fmt.Errorf("%w: %w", model.ErrInventoryInternal, err)
	}
}
