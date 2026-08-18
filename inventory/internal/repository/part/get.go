package part

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/Mas4trt/microservices/inventory/internal/model"
	repoConverter "github.com/Mas4trt/microservices/inventory/internal/repository/converter"
	repoModel "github.com/Mas4trt/microservices/inventory/internal/repository/model"
)

func (r *repository) Get(ctx context.Context, partUUID string) (model.Part, error) {
	var part repoModel.Part

	err := r.coll.FindOne(ctx, bson.M{"uuid": partUUID}).Decode(&part)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return model.Part{}, model.ErrPartNotFound
		}
		return model.Part{}, err
	}

	return repoConverter.PartToModel(part), nil
}
