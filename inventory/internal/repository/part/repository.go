package part

import (
	"go.mongodb.org/mongo-driver/mongo"

	def "github.com/Mas4trt/microservices/inventory/internal/repository"
)

var _ def.InventoryRepository = (*repository)(nil)

type repository struct {
	client     *mongo.Client
	collection *mongo.Collection
}

func NewRepository(client *mongo.Client, dbName string) *repository {
	return &repository{
		client:     client,
		collection: client.Database(dbName).Collection("parts"),
	}
}
