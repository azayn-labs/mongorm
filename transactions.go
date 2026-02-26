package mongorm

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// WithTransaction executes fn within a MongoDB transaction.
//
// The callback receives a transaction-bound context. Any MongORM operation
// using that context will run in the same transaction session.
func (m *MongORM[T]) WithTransaction(
	ctx context.Context,
	fn func(txCtx context.Context) error,
	opts ...options.Lister[options.TransactionOptions],
) error {
	if err := m.ensureReady(); err != nil {
		return err
	}

	if fn == nil {
		return configErrorf("transaction callback cannot be nil")
	}

	client, err := m.client()
	if err != nil {
		return err
	}

	session, err := client.StartSession()
	if err != nil {
		return normalizeError(err)
	}
	defer session.EndSession(ctx)

	_, err = session.WithTransaction(
		ctx,
		func(txCtx context.Context) (any, error) {
			if err := fn(txCtx); err != nil {
				return nil, err
			}

			return nil, nil
		},
		opts...,
	)

	return normalizeError(err)
}

func (m *MongORM[T]) client() (*mongo.Client, error) {
	if m == nil {
		return nil, configErrorf("mongorm instance is nil")
	}

	if m.options != nil && m.options.MongoClient != nil {
		return m.options.MongoClient, nil
	}

	if m.info != nil && m.info.db != nil {
		return m.info.db.Client(), nil
	}

	return nil, configErrorf("mongodb client is not initialized")
}
