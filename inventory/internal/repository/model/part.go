package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
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
	Length float64 `bson:"length"`
	Width  float64 `bson:"width"`
	Height float64 `bson:"height"`
	Weight float64 `bson:"weight"`
}

type Manufacturer struct {
	Name    string `bson:"name"`
	Country string `bson:"country"`
	Website string `bson:"website"`
}

type Part struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	UUID          string             `bson:"uuid"`
	Name          string             `bson:"name"`
	Description   string             `bson:"description"`
	Price         float64            `bson:"price"`
	StockQuantity int64              `bson:"stock_quantity"`
	Category      string             `bson:"category"`
	Dimensions    Dimensions         `bson:"dimensions"`
	Manufacturer  Manufacturer       `bson:"manufacturer"`
	Tags          []string           `bson:"tags,omitempty"`
	Metadata      map[string]any     `bson:"metadata,omitempty"`
	CreatedAt     time.Time          `bson:"created_at"`
	UpdatedAt     *time.Time         `bson:"updated_at,omitempty"`
}
