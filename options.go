package mongorm

import (
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// MongORMOptions holds the configuration options for a MongORM instance, including
// settings for timestamps, collection and database names, and the MongoDB client. This
// struct is used to customize the behavior of the MongORM instance when connecting to the
// database and performing operations.
//
// Example usage:
//
//	type ToDo struct {
//	  Text *string `bson:"text"`
//	}
//	options := &mongorm.MongORMOptions{
//	    Timestamps: true,
//	    CollectionName: ptr("todos"),
//	    DatabaseName: ptr("mydb"),
//	    MongoClient: myMongoClient,
//	}
//	orm := mongorm.FromOptions(&ToDo{}, options)
type MongORMOptions struct {
	Timestamps     bool          `json:"-"`
	CollectionName *string       `json:"-"`
	DatabaseName   *string       `json:"-"`
	MongoClient    *mongo.Client `json:"-"`
}

// FromOptions creates a new MongORM instance with the provided schema and options. This function
// allows you to customize the MongORM instance with specific settings for timestamps, collection
// and database names, and the MongoDB client. If options is nil, default settings will be used.
//
// Example usage:
//
//	type ToDo struct {
//	  Text *string `bson:"text"`
//	}
//	options := &mongorm.MongORMOptions{
//	    Timestamps: true,
//	    CollectionName: ptr("todos"),
//	    DatabaseName: ptr("mydb"),
//	    MongoClient: myMongoClient,
//	}
//	orm := mongorm.FromOptions(&ToDo{}, options)
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
	m.operations.reset()

	return m
}
