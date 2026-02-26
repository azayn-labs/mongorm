package mongorm

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// BulkWrite executes multiple write models in a single request.
func (m *MongORM[T]) BulkWrite(
	ctx context.Context,
	models []mongo.WriteModel,
	opts ...options.Lister[options.BulkWriteOptions],
) (*mongo.BulkWriteResult, error) {
	if err := m.ensureReady(); err != nil {
		return nil, err
	}

	if len(models) == 0 {
		return nil, fmt.Errorf("bulk write requires at least one write model")
	}

	return m.info.collection.BulkWrite(ctx, models, opts...)
}

// BulkWriteInTransaction executes a bulk write operation inside a transaction.
func (m *MongORM[T]) BulkWriteInTransaction(
	ctx context.Context,
	models []mongo.WriteModel,
	bulkOpts ...options.Lister[options.BulkWriteOptions],
) (*mongo.BulkWriteResult, error) {
	var result *mongo.BulkWriteResult

	err := m.WithTransaction(ctx, func(txCtx context.Context) error {
		res, err := m.BulkWrite(txCtx, models, bulkOpts...)
		if err != nil {
			return err
		}

		result = res

		return nil
	})
	if err != nil {
		return nil, err
	}

	return result, nil
}

// BulkWriteBuilder builds typed MongoDB write models for bulk operations.
type BulkWriteBuilder[T any] struct {
	models []mongo.WriteModel
}

// NewBulkWriteBuilder creates a new bulk write model builder.
func NewBulkWriteBuilder[T any]() *BulkWriteBuilder[T] {
	return &BulkWriteBuilder[T]{
		models: []mongo.WriteModel{},
	}
}

// InsertOne appends an insert-one write model.
func (b *BulkWriteBuilder[T]) InsertOne(document *T) *BulkWriteBuilder[T] {
	if document == nil {
		return b
	}

	b.models = append(b.models, mongo.NewInsertOneModel().SetDocument(document))

	return b
}

// UpdateOne appends an update-one write model.
func (b *BulkWriteBuilder[T]) UpdateOne(
	filter bson.M,
	update bson.M,
	upsert bool,
) *BulkWriteBuilder[T] {
	b.models = append(
		b.models,
		mongo.NewUpdateOneModel().
			SetFilter(filter).
			SetUpdate(update).
			SetUpsert(upsert),
	)

	return b
}

// UpdateMany appends an update-many write model.
func (b *BulkWriteBuilder[T]) UpdateMany(
	filter bson.M,
	update bson.M,
	upsert bool,
) *BulkWriteBuilder[T] {
	b.models = append(
		b.models,
		mongo.NewUpdateManyModel().
			SetFilter(filter).
			SetUpdate(update).
			SetUpsert(upsert),
	)

	return b
}

// ReplaceOne appends a replace-one write model.
func (b *BulkWriteBuilder[T]) ReplaceOne(
	filter bson.M,
	replacement *T,
	upsert bool,
) *BulkWriteBuilder[T] {
	if replacement == nil {
		return b
	}

	b.models = append(
		b.models,
		mongo.NewReplaceOneModel().
			SetFilter(filter).
			SetReplacement(replacement).
			SetUpsert(upsert),
	)

	return b
}

// DeleteOne appends a delete-one write model.
func (b *BulkWriteBuilder[T]) DeleteOne(filter bson.M) *BulkWriteBuilder[T] {
	b.models = append(b.models, mongo.NewDeleteOneModel().SetFilter(filter))

	return b
}

// DeleteMany appends a delete-many write model.
func (b *BulkWriteBuilder[T]) DeleteMany(filter bson.M) *BulkWriteBuilder[T] {
	b.models = append(b.models, mongo.NewDeleteManyModel().SetFilter(filter))

	return b
}

// Models returns the accumulated write models.
func (b *BulkWriteBuilder[T]) Models() []mongo.WriteModel {
	if b == nil {
		return nil
	}

	return b.models
}
