package part

import (
	"github.com/Mas4trt/microservices/inventory/internal/model"
	"github.com/brianvoe/gofakeit/v7"
)

func newTestParts(count int) []model.Part {
	parts := make([]model.Part, count)

	for i := range parts {
		parts[i] = newTestPart()
	}

	return parts
}

func newTestPart() model.Part {
	categories := []model.Category{
		model.CategoryEngine,
		model.CategoryFuel,
		model.CategoryPorthole,
		model.CategoryWing,
	}
	return model.Part{
		Uuid:          gofakeit.UUID(),
		Name:          gofakeit.Name(),
		Description:   gofakeit.Sentence(10),
		Price:         gofakeit.Price(100, 100000),
		StockQuantity: int64(gofakeit.IntRange(0, 10000)),

		Category: categories[gofakeit.IntRange(0, len(categories)-1)],

		Dimensions: model.Dimensions{
			Length: gofakeit.Float64Range(0.1, 10),
			Width:  gofakeit.Float64Range(0.1, 10),
			Height: gofakeit.Float64Range(0.1, 10),
			Weight: gofakeit.Float64Range(0.1, 5000),
		},

		Manufacturer: model.Manufacturer{
			Name:    gofakeit.Company(),
			Country: gofakeit.Country(),
			Website: gofakeit.URL(),
		},

		Tags: []string{
			gofakeit.Word(),
			gofakeit.Word(),
			gofakeit.Word(),
		},

		Metadata: map[string]any{
			"material":      gofakeit.Word(),
			"serial_number": gofakeit.UUID(),
			"certified":     gofakeit.Bool(),
			"max_rpm":       gofakeit.Int64(),
			"fuel_type":     gofakeit.Word(),
			"revision":      gofakeit.LetterN(1),
		},

		CreatedAt: gofakeit.Date(),

		UpdatedAt: timePtr(gofakeit.Date()),
	}
}

func newTestPartsFilter() model.PartsFilter {
	return model.PartsFilter{
		UUIDs: []string{
			gofakeit.UUID(),
			gofakeit.UUID(),
			gofakeit.UUID(),
		},
		Names: []string{
			gofakeit.Name(),
			gofakeit.Name(),
			gofakeit.Name(),
		},
		Categories: []model.Category{
			model.CategoryEngine,
			model.CategoryFuel,
			model.CategoryWing,
		},
		ManufacturerCountries: []string{
			gofakeit.Country(),
			gofakeit.Country(),
			gofakeit.Country(),
		},
		Tags: []string{
			gofakeit.Word(),
			gofakeit.Word(),
			gofakeit.Word(),
		},
	}
}
