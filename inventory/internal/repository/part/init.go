package part

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/Mas4trt/microservices/inventory/internal/model"
	repoConverter "github.com/Mas4trt/microservices/inventory/internal/repository/converter"
)

func (r *repository) Init(ctx context.Context, parts map[string]model.Part) error {
	if len(parts) == 0 {
		return model.ErrEmptyData
	}

	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "uuid", Value: 1}},
		Options: options.Index().SetUnique(true),
	}

	_, err := r.coll.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		return fmt.Errorf("failed to create index: %w", err)
	}

	count, err := r.coll.CountDocuments(ctx, bson.M{}, options.Count().SetLimit(1))
	if err != nil {
		return fmt.Errorf("failed to count parts: %w", err)
	}

	if count > 0 {
		return model.ErrAlreadyInitialized
	}

	docs := make([]any, 0, len(parts))
	for _, p := range parts {
		repoPart := repoConverter.PartToRepoModel(p)
		docs = append(docs, repoPart)
	}

	_, err = r.coll.InsertMany(ctx, docs)
	if err != nil {
		return fmt.Errorf("failed to seed parts: %w", err)
	}

	return nil
}
