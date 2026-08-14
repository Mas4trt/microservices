package model

import (
	"time"
)

type Category int32

const (
	CategoryUnspecified Category = 0
	CategoryEngine      Category = 1
	CategoryFuel        Category = 2
	CategoryPorthole    Category = 3
	CategoryWing        Category = 4
)

var categoryName = map[Category]string{
	CategoryUnspecified: "CATEGORY_UNSPECIFIED",
	CategoryEngine:      "CATEGORY_ENGINE",
	CategoryFuel:        "CATEGORY_FUEL",
	CategoryPorthole:    "CATEGORY_PORTHOLE",
	CategoryWing:        "CATEGORY_WING",
}

func (c Category) String() string {
	if name, ok := categoryName[c]; ok {
		return name
	}
	return "CATEGORY_UNSPECIFIED"
}

type Dimensions struct {
	Length float64
	Width  float64
	Height float64
	Weight float64
}

type Manufacturer struct {
	Name    string
	Country string
	Website string
}

type Value struct {
	StringValue *string
	Int64Value  *int64
	DoubleValue *float64
	BoolValue   *bool
}

type Part struct {
	Uuid          string
	Name          string
	Description   string
	Price         float64
	StockQuantity int64
	Category      Category
	Dimensions    Dimensions
	Manufacturer  Manufacturer
	Tags          []string
	Metadata      map[string]any
	CreatedAt     time.Time
	UpdatedAt     *time.Time
}

type PartsFilter struct {
	UUIDs                 []string
	Names                 []string
	Categories            []Category
	ManufacturerCountries []string
	Tags                  []string
}
