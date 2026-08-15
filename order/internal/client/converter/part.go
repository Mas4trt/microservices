package converter

import (
	"fmt"

	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/Mas4trt/microservices/order/internal/model"
	inventoryV1 "github.com/Mas4trt/microservices/shared/pkg/proto/inventory/v1"
)

func PartListToModel(parts []*inventoryV1.Part) ([]model.PricedPart, error) {
	res := make([]model.PricedPart, 0, len(parts))
	for _, part := range parts {
		p, err := PartToModel(part)
		if err != nil {
			return nil, err
		}
		res = append(res, p)
	}

	return res, nil
}

func PartToModel(part *inventoryV1.Part) (model.PricedPart, error) {
	partUUID, err := uuid.Parse(part.GetUuid())
	if err != nil {
		return model.PricedPart{}, fmt.Errorf("parse part uuid %q: %w", part.GetUuid(), err)
	}

	return model.PricedPart{
		UUID:  partUUID,
		Price: part.GetPrice(),
	}, nil
}

func PartsFilterToProto(filter model.PartsFilter) *inventoryV1.PartsFilter {
	categories := make([]inventoryV1.Category, 0, len(filter.Categories))
	for _, category := range filter.Categories {
		categories = append(categories, inventoryV1.Category(category))
	}

	uuids := make([]*wrapperspb.StringValue, 0, len(filter.UUIDs))
	for _, id := range filter.UUIDs {
		uuids = append(uuids, wrapperspb.String(id))
	}

	names := make([]*wrapperspb.StringValue, 0, len(filter.Names))
	for _, name := range filter.Names {
		names = append(names, wrapperspb.String(name))
	}

	manufacturerCountries := make([]*wrapperspb.StringValue, 0, len(filter.ManufacturerCountries))
	for _, country := range filter.ManufacturerCountries {
		manufacturerCountries = append(manufacturerCountries, wrapperspb.String(country))
	}

	tags := make([]*wrapperspb.StringValue, 0, len(filter.Tags))
	for _, tag := range filter.Tags {
		tags = append(tags, wrapperspb.String(tag))
	}

	return &inventoryV1.PartsFilter{
		Uuids:                 uuids,
		Names:                 names,
		Categories:            categories,
		ManufacturerCountries: manufacturerCountries,
		Tags:                  tags,
	}
}
