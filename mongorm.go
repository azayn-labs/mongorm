package mongorm

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

type MongoORM struct {
	MongoClient *mongo.Client
}

func NewMongoORM(
	ctx context.Context,
	client *mongo.Client,
) *MongoORM {
	return &MongoORM{
		MongoClient: client,
	}
}
