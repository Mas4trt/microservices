package model

import "time"

type Category int32

const (
	CategoryUnspecified Category = iota
	CategoryEngine
	CategoryFuel
	CategoryPorthole
	CategoryWing
)

func (c Category) String() string {
	switch c {
	case CategoryEngine:
		return "CATEGORY_ENGINE"
	case CategoryFuel:
		return "CATEGORY_FUEL"
	case CategoryPorthole:
		return "CATEGORY_PORTHOLE"
	case CategoryWing:
		return "CATEGORY_WING"
	default:
		return "CATEGORY_UNSPECIFIED"
	}
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
	UUID          string
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
