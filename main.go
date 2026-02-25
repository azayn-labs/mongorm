package mongorm

import (
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type MongORMOperations struct {
	query  bson.M `json:"-"`
	update bson.M `json:"-"`
}

type MongORMInfo struct {
	dbName     *string           `json:"-"`
	db         *mongo.Database   `json:"-"`
	collection *mongo.Collection `json:"-"`
	fields     map[string]Field  `json:"-"`
}

type MongORMOptions struct {
	Timestamps     bool          `json:"-"`
	CollectionName *string       `json:"-"`
	DatabaseName   *string       `json:"-"`
	MongoClient    *mongo.Client `json:"-"`
}

type MongORM[T any] struct {
	schema *T

	options    *MongORMOptions    `json:"-"`
	info       *MongORMInfo       `json:"-"`
	operations *MongORMOperations `json:"-"`
}

func New[T any](schema *T) *MongORM[T] {
	return FromOptions(schema, nil)
}

func FromOptions[T any](schema *T, options *MongORMOptions) *MongORM[T] {
	m := &MongORM[T]{
		schema: schema,
		info: &MongORMInfo{
			fields: buildFields[T](),
		},
		operations: &MongORMOperations{
			query:  bson.M{},
			update: bson.M{},
		},
	}

	m.options = &MongORMOptions{}
	if options != nil {
		m.options = options
	}

	m.initializeClient()
	m.operations.Clear()

	return m
}

func NewClient(connectionString string) (*mongo.Client, error) {
	client, err := mongo.Connect(
		options.Client().
			ApplyURI(connectionString).
			SetServerAPIOptions(options.ServerAPI(options.ServerAPIVersion1)),
	)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func (o *MongORMOperations) Clear() {
	o.query = bson.M{}
	o.update = bson.M{}
}
