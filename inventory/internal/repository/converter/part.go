package converter

import (
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
		Uuid:          part.Uuid,
		Name:          part.Name,
		Description:   part.Description,
		Price:         part.Price,
		StockQuantity: part.StockQuantity,
		Category:      repoModel.Category(part.Category),
		Dimensions:    repoModel.Dimensions(part.Dimensions),
		Manufacturer:  repoModel.Manufacturer(part.Manufacturer),
		Tags:          part.Tags,
		Metadata:      part.Metadata,
		CreatedAt:     part.CreatedAt,
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
		Uuid:          part.Uuid,
		Name:          part.Name,
		Description:   part.Description,
		Price:         part.Price,
		StockQuantity: part.StockQuantity,
		Category:      model.Category(part.Category),
		Dimensions:    model.Dimensions(part.Dimensions),
		Manufacturer:  model.Manufacturer(part.Manufacturer),
		Tags:          part.Tags,
		Metadata:      part.Metadata,
		CreatedAt:     part.CreatedAt,
	}
}

func FilterToRepoModel(filter model.PartsFilter) repoModel.PartsFilter {
	categories := make([]repoModel.Category, len(filter.Categories))
	for i, category := range filter.Categories {
		categories[i] = repoModel.Category(category)
	}

	return repoModel.PartsFilter{
		UUIDs:                 filter.UUIDs,
		Names:                 filter.Names,
		Categories:            categories,
		ManufacturerCountries: filter.ManufacturerCountries,
		Tags:                  filter.Tags,
	}
}
