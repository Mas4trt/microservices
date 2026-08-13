package converter

import (
	"log"
	"time"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/Mas4trt/microservices/order/internal/model"
	inventoryV1 "github.com/Mas4trt/microservices/shared/pkg/proto/inventory/v1"
)

func PartListToModel(parts []*inventoryV1.Part) []model.Part {
	res := make([]model.Part, 0, len(parts))
	for _, part := range parts {
		res = append(res, PartToModel(part))
	}

	return res
}

func PartToModel(part *inventoryV1.Part) model.Part {
	partUUID, err := uuid.Parse(part.GetUuid())
	if err != nil {
		return model.Part{}
	}

	var updatedAt *time.Time
	if part.UpdatedAt != nil {
		t := part.UpdatedAt.AsTime()
		updatedAt = &t
	}

	return model.Part{
		UUID:          partUUID,
		Name:          part.Name,
		Description:   part.Description,
		Price:         part.Price,
		StockQuantity: part.StockQuantity,
		Category:      model.Category(part.GetCategory()),
		Dimensions:    DimensionsToModel(part.Dimensions),
		Manufacturer:  ManufacturerToModel(part.Manufacturer),
		Tags:          part.Tags,
		Metadata:      PartMetadataToModel(part.Metadata),
		CreatedAt:     part.CreatedAt.AsTime(),
		UpdatedAt:     updatedAt,
	}
}

func DimensionsToModel(dimensions *inventoryV1.Dimensions) model.Dimensions {
	return model.Dimensions{
		Length: dimensions.Length,
		Width:  dimensions.Width,
		Height: dimensions.Height,
		Weight: dimensions.Weight,
	}
}

func ManufacturerToModel(manufacturer *inventoryV1.Manufacturer) model.Manufacturer {
	return model.Manufacturer{
		Name:    manufacturer.Name,
		Country: manufacturer.Country,
		Website: manufacturer.Website,
	}
}

func PartMetadataToModel(metadata map[string]*inventoryV1.Value) map[string]any {
	res := make(map[string]any, len(metadata))

	for key, value := range metadata {
		if value == nil || value.Value == nil {
			continue
		}

		switch v := value.Value.(type) {
		case *inventoryV1.Value_StringValue:
			res[key] = v.StringValue
		case *inventoryV1.Value_BoolValue:
			res[key] = v.BoolValue
		case *inventoryV1.Value_DoubleValue:
			res[key] = v.DoubleValue
		case *inventoryV1.Value_Int64Value:
			res[key] = v.Int64Value
		default:
			log.Printf("PartMetadataToModel: unsupported type %T for key %q, value: %+v", v, key, value)
		}
	}

	return res
}

func PartsFilterToProto(filter model.PartsFilter) *inventoryV1.PartsFilter {
	categories := make([]inventoryV1.Category, 0, len(filter.Categories))
	for _, category := range filter.Categories {
		categories = append(
			categories,
			inventoryV1.Category(category),
		)
	}

	uuids := make([]*wrapperspb.StringValue, 0, len(filter.UUIDs))
	for _, id := range filter.UUIDs {
		uuids = append(
			uuids,
			wrapperspb.String(id),
		)
	}

	names := make([]*wrapperspb.StringValue, 0, len(filter.Names))
	for _, name := range filter.Names {
		names = append(
			names,
			wrapperspb.String(name),
		)
	}

	manufacturerCountries := make([]*wrapperspb.StringValue, 0, len(filter.ManufacturerCountries))
	for _, country := range filter.ManufacturerCountries {
		manufacturerCountries = append(
			manufacturerCountries,
			wrapperspb.String(country),
		)
	}

	tags := make([]*wrapperspb.StringValue, 0, len(filter.Tags))
	for _, tag := range filter.Tags {
		tags = append(
			tags,
			wrapperspb.String(tag),
		)
	}

	return &inventoryV1.PartsFilter{
		Uuids:                 uuids,
		Names:                 names,
		Categories:            categories,
		ManufacturerCountries: manufacturerCountries,
		Tags:                  tags,
	}
}
