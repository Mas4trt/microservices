package part

import (
	"go.mongodb.org/mongo-driver/mongo"

	def "github.com/Mas4trt/microservices/inventory/internal/repository"
)

var _ def.InventoryRepository = (*repository)(nil)

type repository struct {
	coll *mongo.Collection
}

func NewRepository(db *mongo.Database) *repository {
	return &repository{
		coll: db.Collection("parts"),
	}
}
