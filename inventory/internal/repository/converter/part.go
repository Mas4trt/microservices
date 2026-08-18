package converter

import (
	"go.mongodb.org/mongo-driver/bson"

	"github.com/Mas4trt/microservices/inventory/internal/model"
	repoModel "github.com/Mas4trt/microservices/inventory/internal/repository/model"
)

func PartsToRepoModel(parts map[string]model.Part) map[string]repoModel.Part {
	if parts == nil {
		return nil
	}

	result := make(map[string]repoModel.Part, len(parts))

	for key, part := range parts {
		result[key] = PartToRepoModel(part)
	}

	return result
}

func PartToRepoModel(part model.Part) repoModel.Part {
	return repoModel.Part{
		UUID:          part.UUID,
		Name:          part.Name,
		Description:   part.Description,
		Price:         part.Price,
		StockQuantity: part.StockQuantity,
		Category:      part.Category.String(),
		Dimensions:    repoModel.Dimensions(part.Dimensions),
		Manufacturer:  repoModel.Manufacturer(part.Manufacturer),
		Tags:          part.Tags,
		Metadata:      cloneMetadata(part.Metadata),
		CreatedAt:     part.CreatedAt,
		UpdatedAt:     part.UpdatedAt,
	}
}

func PartsToModel(parts []repoModel.Part) []model.Part {
	result := make([]model.Part, 0, len(parts))

	for _, part := range parts {
		result = append(result, PartToModel(part))
	}

	return result
}

func PartToModel(part repoModel.Part) model.Part {
	return model.Part{
		UUID:          part.UUID,
		Name:          part.Name,
		Description:   part.Description,
		Price:         part.Price,
		StockQuantity: part.StockQuantity,
		Category:      categoryToModel(part.Category),
		Dimensions:    model.Dimensions(part.Dimensions),
		Manufacturer:  model.Manufacturer(part.Manufacturer),
		Tags:          part.Tags,
		Metadata:      cloneMetadata(part.Metadata),
		CreatedAt:     part.CreatedAt,
		UpdatedAt:     part.UpdatedAt,
	}
}

func categoryToModel(category string) model.Category {
	switch category {
	case model.CategoryEngine.String():
		return model.CategoryEngine
	case model.CategoryFuel.String():
		return model.CategoryFuel
	case model.CategoryPorthole.String():
		return model.CategoryPorthole
	case model.CategoryWing.String():
		return model.CategoryWing
	default:
		return model.CategoryUnspecified
	}
}

func cloneMetadata(metadata map[string]any) map[string]any {
	if metadata == nil {
		return nil
	}

	result := make(map[string]any, len(metadata))

	for key, value := range metadata {
		result[key] = value
	}

	return result
}

func FilterToBSON(filter model.PartsFilter) bson.M {
	query := bson.M{}

	if len(filter.UUIDs) > 0 {
		query["uuid"] = bson.M{"$in": filter.UUIDs}
	}

	if len(filter.Names) > 0 {
		query["name"] = bson.M{"$in": filter.Names}
	}

	if len(filter.Categories) > 0 {
		categories := make([]string, 0, len(filter.Categories))
		for _, cat := range filter.Categories {
			categories = append(categories, cat.String())
		}
		query["category"] = bson.M{"$in": categories}
	}

	if len(filter.ManufacturerCountries) > 0 {
		query["manufacturer.country"] = bson.M{"$in": filter.ManufacturerCountries}
	}

	if len(filter.Tags) > 0 {
		query["tags"] = bson.M{"$in": filter.Tags}
	}

	return query
}
