package v1

import (
	"context"
	"log"

	orderV1 "github.com/Mas4trt/microservices/shared/pkg/openapi/order/v1"
)

func (a *api) NewError(ctx context.Context, err error) *orderV1.GenericErrorStatusCode {
	log.Printf("internal error: %v", err)

	return &orderV1.GenericErrorStatusCode{
		StatusCode: 500,
		Response: orderV1.GenericError{
			Code:    "INTERNAL_ERROR",
			Message: "internal server error",
		},
	}
}
