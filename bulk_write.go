package mongorm

import (
	"context"

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
		return nil, configErrorf("bulk write requires at least one write model")
	}

	result, err := m.info.collection.BulkWrite(ctx, models, opts...)
	if err != nil {
		return nil, normalizeError(err)
	}

	return result, nil
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

// FieldValuePair is a field/value helper for field-only update builders.
type FieldValuePair struct {
	Field Field
	Value any
}

// FilterBy creates an equality filter using a schema field.
func FilterBy(field Field, value any) bson.M {
	if field == nil {
		return bson.M{}
	}

	return bson.M{field.BSONName(): value}
}

// SetUpdateFromPairs builds an update document with $set using schema fields.
func SetUpdateFromPairs(pairs ...FieldValuePair) bson.M {
	set := bson.M{}
	for _, pair := range pairs {
		if pair.Field == nil {
			continue
		}
		set[pair.Field.BSONName()] = pair.Value
	}

	if len(set) == 0 {
		return bson.M{}
	}

	return bson.M{"$set": set}
}

// SetOnInsertUpdateFromPairs builds an update document with $setOnInsert using schema fields.
func SetOnInsertUpdateFromPairs(pairs ...FieldValuePair) bson.M {
	setOnInsert := bson.M{}
	for _, pair := range pairs {
		if pair.Field == nil {
			continue
		}
		setOnInsert[pair.Field.BSONName()] = pair.Value
	}

	if len(setOnInsert) == 0 {
		return bson.M{}
	}

	return bson.M{"$setOnInsert": setOnInsert}
}

// UnsetUpdateFromFields builds an update document with $unset using schema fields.
func UnsetUpdateFromFields(fields ...Field) bson.M {
	unset := bson.M{}
	for _, field := range fields {
		if field == nil {
			continue
		}
		unset[field.BSONName()] = 1
	}

	if len(unset) == 0 {
		return bson.M{}
	}

	return bson.M{"$unset": unset}
}

// IncUpdateFromPairs builds an update document with $inc using schema fields.
func IncUpdateFromPairs(pairs ...FieldValuePair) bson.M {
	inc := bson.M{}
	for _, pair := range pairs {
		if pair.Field == nil {
			continue
		}
		inc[pair.Field.BSONName()] = pair.Value
	}

	if len(inc) == 0 {
		return bson.M{}
	}

	return bson.M{"$inc": inc}
}

// PushUpdateFromPairs builds an update document with $push using schema fields.
func PushUpdateFromPairs(pairs ...FieldValuePair) bson.M {
	push := bson.M{}
	for _, pair := range pairs {
		if pair.Field == nil {
			continue
		}
		push[pair.Field.BSONName()] = pair.Value
	}

	if len(push) == 0 {
		return bson.M{}
	}

	return bson.M{"$push": push}
}

// AddToSetUpdateFromPairs builds an update document with $addToSet using schema fields.
func AddToSetUpdateFromPairs(pairs ...FieldValuePair) bson.M {
	addToSet := bson.M{}
	for _, pair := range pairs {
		if pair.Field == nil {
			continue
		}
		addToSet[pair.Field.BSONName()] = pair.Value
	}

	if len(addToSet) == 0 {
		return bson.M{}
	}

	return bson.M{"$addToSet": addToSet}
}

// PullUpdateFromPairs builds an update document with $pull using schema fields.
func PullUpdateFromPairs(pairs ...FieldValuePair) bson.M {
	pull := bson.M{}
	for _, pair := range pairs {
		if pair.Field == nil {
			continue
		}
		pull[pair.Field.BSONName()] = pair.Value
	}

	if len(pull) == 0 {
		return bson.M{}
	}

	return bson.M{"$pull": pull}
}

// PopUpdateFromPairs builds an update document with $pop using schema fields.
func PopUpdateFromPairs(pairs ...FieldValuePair) bson.M {
	pop := bson.M{}
	for _, pair := range pairs {
		if pair.Field == nil {
			continue
		}
		pop[pair.Field.BSONName()] = pair.Value
	}

	if len(pop) == 0 {
		return bson.M{}
	}

	return bson.M{"$pop": pop}
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

// UpdateOneBy appends an update-one model using a schema field equality filter.
func (b *BulkWriteBuilder[T]) UpdateOneBy(
	field Field,
	value any,
	update bson.M,
	upsert bool,
) *BulkWriteBuilder[T] {
	return b.UpdateOne(FilterBy(field, value), update, upsert)
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

// UpdateManyBy appends an update-many model using a schema field equality filter.
func (b *BulkWriteBuilder[T]) UpdateManyBy(
	field Field,
	value any,
	update bson.M,
	upsert bool,
) *BulkWriteBuilder[T] {
	return b.UpdateMany(FilterBy(field, value), update, upsert)
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

// ReplaceOneBy appends a replace-one model using a schema field equality filter.
func (b *BulkWriteBuilder[T]) ReplaceOneBy(
	field Field,
	value any,
	replacement *T,
	upsert bool,
) *BulkWriteBuilder[T] {
	return b.ReplaceOne(FilterBy(field, value), replacement, upsert)
}

// DeleteOne appends a delete-one write model.
func (b *BulkWriteBuilder[T]) DeleteOne(filter bson.M) *BulkWriteBuilder[T] {
	b.models = append(b.models, mongo.NewDeleteOneModel().SetFilter(filter))

	return b
}

// DeleteOneBy appends a delete-one model using a schema field equality filter.
func (b *BulkWriteBuilder[T]) DeleteOneBy(field Field, value any) *BulkWriteBuilder[T] {
	return b.DeleteOne(FilterBy(field, value))
}

// DeleteMany appends a delete-many write model.
func (b *BulkWriteBuilder[T]) DeleteMany(filter bson.M) *BulkWriteBuilder[T] {
	b.models = append(b.models, mongo.NewDeleteManyModel().SetFilter(filter))

	return b
}

// DeleteManyBy appends a delete-many model using a schema field equality filter.
func (b *BulkWriteBuilder[T]) DeleteManyBy(field Field, value any) *BulkWriteBuilder[T] {
	return b.DeleteMany(FilterBy(field, value))
}

// Models returns the accumulated write models.
func (b *BulkWriteBuilder[T]) Models() []mongo.WriteModel {
	if b == nil {
		return nil
	}

	return b.models
}
