package orm

import (
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type ModelOptions struct {
	CollectionName string          `json:"-"`
	Database       string          `json:"-"`
	Timestamps     bool            `json:"-"`
	DB             *mongo.Database `json:"-"`
}

type Model[T any] struct {
	schema     *T
	clone      *T
	options    ModelOptions
	db         *mongo.Database   `json:"-"`
	collection *mongo.Collection `json:"-"`
}

func NewModel[T any](schema *T, options ModelOptions) *Model[T] {

	m := &Model[T]{
		schema:  schema,
		options: options,
		clone:   clonePtr(schema, false),
	}

	if options.DB != nil {
		m.db = options.DB
		m.collection = m.db.Collection(options.CollectionName)
	}

	return m
}
