package mongorm

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// SaveMulti performs an update operation on multiple documents that match the specified query criteria.
// It uses the UpdateMany method of the MongoDB collection to apply the update operations to all matching documents.
// The query and update operations are constructed based on the state of the MongORM instance.
// The caller can provide additional options for the UpdateMany operation using the opts parameter.
// It returns an UpdateResult containing information about the operation, such as the number of documents matched and modified, or an error if the operation fails.
//
// Example usage:
//
//	updateResult, err := mongormInstance.SaveMulti(ctx)
//	if err != nil {
//	    // Handle error
//	} else {
//	    // Use updateResult
//	}
func (m *MongORM[T]) SaveMulti(
	ctx context.Context,
	opts ...options.Lister[options.UpdateManyOptions],
) (*mongo.UpdateResult, error) {
	if err := m.ensureReady(); err != nil {
		return nil, err
	}

	m.operations.fixQuery()
	m.operations.fixUpdate()

	if len(m.operations.update) == 0 {
		return nil, fmt.Errorf("no update operations specified")
	}

	res, err := m.info.collection.UpdateMany(
		ctx,
		m.operations.query,
		m.operations.update,
		opts...,
	)
	if err != nil {
		return nil, normalizeError(err)
	}

	return res, err
}

// FindAll retrieves all documents that match the specified query criteria and returns
// a cursor for iterating over the results. It uses the Find method of the MongoDB
// collection to execute the query and obtain a cursor for the matching documents.
// The query is constructed based on the state of the MongORM instance, and additional
// options can be provided using the opts parameter. The caller is responsible for
// closing the cursor when done to free up resources. It is recommended to use the cursor's
// Next method for large result sets to avoid loading all documents into memory at once.
//
// Example usage:
//
//	cursor, err := mongormInstance.FindAll(ctx)
//	if err != nil {
//	    // Handle error
//	} else {
//	    // Use cursor
//	}
func (m *MongORM[T]) FindAll(
	ctx context.Context,
	opts ...options.Lister[options.FindOptions],
) (*MongORMCursor[T], error) {
	if err := m.ensureReady(); err != nil {
		return nil, err
	}

	filters, _, err := m.withPrimaryFilters()
	if err != nil {
		return nil, normalizeError(err)
	}

	allOpts := []options.Lister[options.FindOptions]{
		m.operations.findOptions(),
	}
	allOpts = append(allOpts, opts...)
	allOpts = append(allOpts, options.Find().SetAllowDiskUse(true))

	cursor, err := m.info.collection.Find(
		ctx,
		filters,
		allOpts...,
	)
	if err != nil {
		return nil, err
	}

	return &MongORMCursor[T]{MongoCursor: cursor, m: m}, nil
}

// DeleteMulti removes all documents that match the current filters.
// It returns a DeleteResult containing the number of removed documents.
func (m *MongORM[T]) DeleteMulti(
	ctx context.Context,
	opts ...options.Lister[options.DeleteManyOptions],
) (*mongo.DeleteResult, error) {
	if err := m.ensureReady(); err != nil {
		return nil, err
	}

	filter, _, err := m.withPrimaryFilters()
	if err != nil {
		return nil, err
	}

	schema := any(m.schema)
	if hook, ok := schema.(BeforeDeleteHook[T]); ok {
		if err := hook.BeforeDelete(m, &filter); err != nil {
			return nil, err
		}
	}

	res, err := m.info.collection.DeleteMany(ctx, filter, opts...)
	if err != nil {
		return nil, normalizeError(err)
	}

	m.operations.reset()

	if hook, ok := schema.(AfterDeleteHook[T]); ok {
		if err := hook.AfterDelete(m); err != nil {
			return nil, err
		}
	}

	return res, nil
}

// Count returns the number of documents that match the current filters.
func (m *MongORM[T]) Count(
	ctx context.Context,
	opts ...options.Lister[options.CountOptions],
) (int64, error) {
	if err := m.ensureReady(); err != nil {
		return 0, err
	}

	filter, _, err := m.withPrimaryFilters()
	if err != nil {
		return 0, err
	}

	count, err := m.info.collection.CountDocuments(ctx, filter, opts...)
	if err != nil {
		return 0, normalizeError(err)
	}

	return count, nil
}

// Distinct returns all unique values of the given field among documents
// that match the current filters.
func (m *MongORM[T]) Distinct(
	ctx context.Context,
	field Field,
	opts ...options.Lister[options.DistinctOptions],
) ([]any, error) {
	if err := m.ensureReady(); err != nil {
		return nil, err
	}

	filter, _, err := m.withPrimaryFilters()
	if err != nil {
		return nil, err
	}

	result := m.info.collection.Distinct(ctx, field.BSONName(), filter, opts...)

	values := []any{}
	if err := result.Decode(&values); err != nil {
		return nil, normalizeError(err)
	}

	return values, nil
}

// DistinctFieldAs returns distinct values cast/converted to the requested type.
func DistinctFieldAs[T any, V any](
	m *MongORM[T],
	ctx context.Context,
	field Field,
	opts ...options.Lister[options.DistinctOptions],
) ([]V, error) {
	values, err := m.Distinct(ctx, field, opts...)
	if err != nil {
		return nil, err
	}

	return castDistinctValues[V](values)
}

// DistinctAs converts raw distinct values into a typed slice.
func DistinctAs[V any](values []any) ([]V, error) {
	return castDistinctValues[V](values)
}

// DistinctStrings returns distinct string values for the given field.
func (m *MongORM[T]) DistinctStrings(
	ctx context.Context,
	field Field,
	opts ...options.Lister[options.DistinctOptions],
) ([]string, error) {
	return DistinctFieldAs[T, string](m, ctx, field, opts...)
}

// DistinctInt64 returns distinct integer values (as int64) for the given field.
func (m *MongORM[T]) DistinctInt64(
	ctx context.Context,
	field Field,
	opts ...options.Lister[options.DistinctOptions],
) ([]int64, error) {
	return DistinctFieldAs[T, int64](m, ctx, field, opts...)
}

// DistinctBool returns distinct boolean values for the given field.
func (m *MongORM[T]) DistinctBool(
	ctx context.Context,
	field Field,
	opts ...options.Lister[options.DistinctOptions],
) ([]bool, error) {
	return DistinctFieldAs[T, bool](m, ctx, field, opts...)
}

// DistinctFloat64 returns distinct numeric values (as float64) for the given field.
func (m *MongORM[T]) DistinctFloat64(
	ctx context.Context,
	field Field,
	opts ...options.Lister[options.DistinctOptions],
) ([]float64, error) {
	return DistinctFieldAs[T, float64](m, ctx, field, opts...)
}

// DistinctObjectIDs returns distinct ObjectID values for the given field.
func (m *MongORM[T]) DistinctObjectIDs(
	ctx context.Context,
	field Field,
	opts ...options.Lister[options.DistinctOptions],
) ([]bson.ObjectID, error) {
	return DistinctFieldAs[T, bson.ObjectID](m, ctx, field, opts...)
}

// DistinctTimes returns distinct time values for the given field.
func (m *MongORM[T]) DistinctTimes(
	ctx context.Context,
	field Field,
	opts ...options.Lister[options.DistinctOptions],
) ([]time.Time, error) {
	return DistinctFieldAs[T, time.Time](m, ctx, field, opts...)
}

// Aggregate runs an aggregation pipeline and returns a MongORM cursor decoding
// each result document into T.
func (m *MongORM[T]) Aggregate(
	ctx context.Context,
	pipeline bson.A,
	opts ...options.Lister[options.AggregateOptions],
) (*MongORMCursor[T], error) {
	cursor, err := m.AggregateRaw(ctx, pipeline, opts...)
	if err != nil {
		return nil, err
	}

	return &MongORMCursor[T]{MongoCursor: cursor, m: m}, nil
}

// AggregateRaw runs an aggregation pipeline and returns a raw MongoDB cursor.
func (m *MongORM[T]) AggregateRaw(
	ctx context.Context,
	pipeline bson.A,
	opts ...options.Lister[options.AggregateOptions],
) (*mongo.Cursor, error) {
	if err := m.ensureReady(); err != nil {
		return nil, err
	}

	filters, _, err := m.withPrimaryFilters()
	if err != nil {
		return nil, err
	}

	finalPipeline := bson.A{}
	if len(filters) > 0 {
		finalPipeline = append(finalPipeline, bson.M{"$match": filters})
	}

	if pipeline != nil {
		finalPipeline = append(finalPipeline, pipeline...)
	}

	allOpts := append(opts, options.Aggregate().SetAllowDiskUse(true))

	cursor, err := m.info.collection.Aggregate(ctx, finalPipeline, allOpts...)
	if err != nil {
		return nil, normalizeError(err)
	}

	return cursor, nil
}

// AggregateAs runs an aggregation pipeline and decodes the results into a typed slice.
func AggregateAs[T any, R any](
	m *MongORM[T],
	ctx context.Context,
	pipeline bson.A,
	opts ...options.Lister[options.AggregateOptions],
) ([]R, error) {
	cursor, err := m.AggregateRaw(ctx, pipeline, opts...)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	results := []R{}
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}

func castDistinctValues[V any](values []any) ([]V, error) {
	var sample V
	targetType := reflect.TypeOf(sample)
	out := make([]V, len(values))

	for i, value := range values {
		casted, err := castDistinctValue[V](value, targetType)
		if err != nil {
			return nil, fmt.Errorf("distinct value at index %d: %w", i, err)
		}

		out[i] = casted
	}

	return out, nil
}
