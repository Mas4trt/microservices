package converter

import (
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/Mas4trt/microservices/inventory/internal/model"
	inventoryV1 "github.com/Mas4trt/microservices/shared/pkg/proto/inventory/v1"
)

func PartsFilterToModel(filter *inventoryV1.PartsFilter) model.PartsFilter {
	if filter == nil {
		return model.PartsFilter{}
	}

	return model.PartsFilter{
		UUIDs:                 StringValuesToStrings(filter.GetUuids()),
		Names:                 StringValuesToStrings(filter.GetNames()),
		Categories:            CategoriesToModel(filter.GetCategories()),
		ManufacturerCountries: StringValuesToStrings(filter.GetManufacturerCountries()),
		Tags:                  StringValuesToStrings(filter.GetTags()),
	}
}

func StringValuesToStrings(values []*wrapperspb.StringValue) []string {
	if len(values) == 0 {
		return nil
	}

	result := make([]string, len(values))

	for i, value := range values {
		result[i] = value.GetValue()
	}

	return result
}

func CategoriesToModel(categories []inventoryV1.Category) []model.Category {
	if len(categories) == 0 {
		return nil
	}

	result := make([]model.Category, len(categories))

	for i, category := range categories {
		result[i] = model.Category(category)
	}

	return result
}

func PartToProto(part model.Part) *inventoryV1.Part {
	return &inventoryV1.Part{
		Uuid:          part.UUID,
		Name:          part.Name,
		Description:   part.Description,
		Price:         part.Price,
		StockQuantity: part.StockQuantity,
		Category:      inventoryV1.Category(part.Category),
		Dimensions:    DimensionsToProto(part.Dimensions),
		Manufacturer:  ManufacturerToProto(part.Manufacturer),
		Tags:          part.Tags,
		Metadata:      MetadataToProto(part.Metadata),
		CreatedAt:     timestamppb.New(part.CreatedAt),
		UpdatedAt:     TimeToTimestamp(part.UpdatedAt),
	}
}

func DimensionsToProto(dimensions model.Dimensions) *inventoryV1.Dimensions {
	return &inventoryV1.Dimensions{
		Length: dimensions.Length,
		Width:  dimensions.Width,
		Height: dimensions.Height,
		Weight: dimensions.Weight,
	}
}

func ManufacturerToProto(manufacturer model.Manufacturer) *inventoryV1.Manufacturer {
	return &inventoryV1.Manufacturer{
		Name:    manufacturer.Name,
		Country: manufacturer.Country,
		Website: manufacturer.Website,
	}
}

func TimeToTimestamp(t *time.Time) *timestamppb.Timestamp {
	if t == nil {
		return nil
	}

	return timestamppb.New(*t)
}

func MetadataToProto(metadata map[string]any) map[string]*inventoryV1.Value {
	if metadata == nil {
		return nil
	}

	result := make(map[string]*inventoryV1.Value, len(metadata))

	for key, value := range metadata {
		if protoVal := AnyToProtoValue(value); protoVal != nil {
			result[key] = protoVal
		}
	}

	return result
}

func AnyToProtoValue(value any) *inventoryV1.Value {
	switch v := value.(type) {
	case string:
		return &inventoryV1.Value{
			Value: &inventoryV1.Value_StringValue{
				StringValue: v,
			},
		}

	case int64:
		return &inventoryV1.Value{
			Value: &inventoryV1.Value_Int64Value{
				Int64Value: v,
			},
		}

	case float64:
		return &inventoryV1.Value{
			Value: &inventoryV1.Value_DoubleValue{
				DoubleValue: v,
			},
		}

	case bool:
		return &inventoryV1.Value{
			Value: &inventoryV1.Value_BoolValue{
				BoolValue: v,
			},
		}

	default:
		return nil
	}
}

func PartsToProto(parts []model.Part) []*inventoryV1.Part {
	if parts == nil {
		return nil
	}

	result := make([]*inventoryV1.Part, len(parts))

	for i, part := range parts {
		result[i] = PartToProto(part)
	}

	return result
}
