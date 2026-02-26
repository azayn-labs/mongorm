package mongorm

import (
	"context"
	"fmt"

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
		return nil, err
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
		return nil, err
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
		return nil, err
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

	return m.info.collection.CountDocuments(ctx, filter, opts...)
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
		return nil, err
	}

	return values, nil
}
