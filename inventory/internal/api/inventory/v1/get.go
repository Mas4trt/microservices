package v1

import (
	"context"
	"errors"
	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/Mas4trt/microservices/inventory/internal/converter"
	"github.com/Mas4trt/microservices/inventory/internal/model"
	inventoryV1 "github.com/Mas4trt/microservices/shared/pkg/proto/inventory/v1"
)

func (a *api) GetPart(ctx context.Context, req *inventoryV1.GetPartRequest) (*inventoryV1.GetPartResponse, error) {
	part, err := a.inventoryService.Get(ctx, req.Uuid)
	if err != nil {
		if errors.Is(err, model.ErrPartNotFound) {
			return nil, status.Errorf(
				codes.NotFound,
				"part with UUID %s not found",
				req.GetUuid(),
			)
		}

		log.Printf("GetPart: internal error: %v", err)
		return nil, status.Error(codes.Internal, "internal server error")
	}

	return &inventoryV1.GetPartResponse{
		Part: converter.PartToProto(part),
	}, nil
}
