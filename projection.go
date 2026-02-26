package mongorm

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// FindOneAs finds a single document and decodes it into a projection DTO type.
func FindOneAs[T any, R any](
	m *MongORM[T],
	ctx context.Context,
	opts ...options.Lister[options.FindOneOptions],
) (*R, error) {
	if err := m.ensureReady(); err != nil {
		return nil, err
	}

	filter, _, err := m.withPrimaryFilters()
	if err != nil {
		return nil, err
	}

	allOpts := []options.Lister[options.FindOneOptions]{
		m.operations.findOneOptions(),
	}
	allOpts = append(allOpts, opts...)

	var result R
	if err := m.info.collection.FindOne(ctx, filter, allOpts...).Decode(&result); err != nil {
		return nil, normalizeError(err)
	}

	m.operations.reset()

	return &result, nil
}

// FindAllAs finds documents and decodes them into a projection DTO slice.
func FindAllAs[T any, R any](
	m *MongORM[T],
	ctx context.Context,
	opts ...options.Lister[options.FindOptions],
) ([]R, error) {
	if err := m.ensureReady(); err != nil {
		return nil, err
	}

	filter, _, err := m.withPrimaryFilters()
	if err != nil {
		return nil, err
	}

	allOpts := []options.Lister[options.FindOptions]{
		m.operations.findOptions(),
	}
	allOpts = append(allOpts, opts...)
	allOpts = append(allOpts, options.Find().SetAllowDiskUse(true))

	cursor, err := m.info.collection.Find(ctx, filter, allOpts...)
	if err != nil {
		return nil, normalizeError(err)
	}

	results := []R{}
	if err := cursor.All(ctx, &results); err != nil {
		if closeErr := cursor.Close(ctx); closeErr != nil {
			return nil, errors.Join(normalizeError(err), normalizeError(closeErr))
		}
		return nil, normalizeError(err)
	}

	if err := cursor.Close(ctx); err != nil {
		return nil, normalizeError(err)
	}

	m.operations.reset()

	return results, nil
}
