package validation

import (
	"errors"

	orderV1 "github.com/Mas4trt/microservices/shared/pkg/openapi/order/v1"
	"github.com/ogen-go/ogen/validate"
)

func ValidationError(err error) *orderV1.ValidationError {
	var validationErr *validate.Error

	if errors.As(err, &validationErr) {
		violations := make(
			[]orderV1.ValidationErrorViolationsItem,
			0,
			len(validationErr.Fields),
		)

		for _, field := range validationErr.Fields {
			violations = append(
				violations,
				orderV1.ValidationErrorViolationsItem{
					Field:   field.Name,
					Message: field.Error.Error(),
				},
			)
		}

		return &orderV1.ValidationError{
			Code:       "VALIDATION_ERROR",
			Message:    "request validation failed",
			Violations: violations,
		}
	}

	return &orderV1.ValidationError{
		Code:    "VALIDATION_ERROR",
		Message: err.Error(),
	}
}
