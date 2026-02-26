package mongorm

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// Asc returns an ascending index key for the given field.
func Asc(field Field) bson.E {
	return bson.E{Key: field.BSONName(), Value: 1}
}

// Desc returns a descending index key for the given field.
func Desc(field Field) bson.E {
	return bson.E{Key: field.BSONName(), Value: -1}
}

// Text returns a text index key for the given field.
func Text(field Field) bson.E {
	return bson.E{Key: field.BSONName(), Value: "text"}
}

// Geo2DSphere returns a 2dsphere index key for the given field.
func Geo2DSphere(field Field) bson.E {
	return bson.E{Key: field.BSONName(), Value: "2dsphere"}
}

// Geo2D returns a 2d index key for the given field.
func Geo2D(field Field) bson.E {
	return bson.E{Key: field.BSONName(), Value: "2d"}
}

// IndexModelFromKeys builds an index model from field-based keys.
func IndexModelFromKeys(keys ...bson.E) mongo.IndexModel {
	return mongo.IndexModel{Keys: bson.D(keys)}
}

// UniqueIndexModelFromKeys builds a unique index model from field-based keys.
func UniqueIndexModelFromKeys(keys ...bson.E) mongo.IndexModel {
	return mongo.IndexModel{
		Keys:    bson.D(keys),
		Options: options.Index().SetUnique(true),
	}
}

// NamedIndexModelFromKeys builds a named index model from field-based keys.
func NamedIndexModelFromKeys(name string, keys ...bson.E) mongo.IndexModel {
	return mongo.IndexModel{
		Keys:    bson.D(keys),
		Options: options.Index().SetName(name),
	}
}

// EnsureIndex creates an index using the provided model.
func (m *MongORM[T]) EnsureIndex(
	ctx context.Context,
	model mongo.IndexModel,
	opts ...options.Lister[options.CreateIndexesOptions],
) (string, error) {
	if err := m.ensureReady(); err != nil {
		return "", err
	}

	return m.info.collection.Indexes().CreateOne(ctx, model, opts...)
}

// EnsureIndexes creates multiple indexes in one call.
func (m *MongORM[T]) EnsureIndexes(
	ctx context.Context,
	models []mongo.IndexModel,
	opts ...options.Lister[options.CreateIndexesOptions],
) ([]string, error) {
	if err := m.ensureReady(); err != nil {
		return nil, err
	}

	return m.info.collection.Indexes().CreateMany(ctx, models, opts...)
}

// Ensure2DSphereIndex creates a 2dsphere index for the provided field.
func (m *MongORM[T]) Ensure2DSphereIndex(
	ctx context.Context,
	field Field,
	opts ...options.Lister[options.CreateIndexesOptions],
) (string, error) {
	model := IndexModelFromKeys(Geo2DSphere(field))

	return m.EnsureIndex(ctx, model, opts...)
}

// EnsureGeoDefaults creates the baseline geospatial indexes for a field:
// - required 2dsphere index on geoField
// - optional supporting index (if supportingKeys are provided)
func (m *MongORM[T]) EnsureGeoDefaults(
	ctx context.Context,
	geoField Field,
	supportingKeys []bson.E,
	opts ...options.Lister[options.CreateIndexesOptions],
) ([]string, error) {
	models := []mongo.IndexModel{
		IndexModelFromKeys(Geo2DSphere(geoField)),
	}

	if len(supportingKeys) > 0 {
		models = append(models, IndexModelFromKeys(supportingKeys...))
	}

	return m.EnsureIndexes(ctx, models, opts...)
}
