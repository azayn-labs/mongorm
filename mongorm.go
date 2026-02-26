package mongorm

import (
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// MongORMInfo holds the internal information for a MongORM instance, including database
// and collection references, as well as field mappings. This struct is used internally
// by the MongORM instance to manage its connection and schema information.
//
// > NOTE: This struct is not intended for public use.
type MongORMInfo struct {
	dbName     *string
	db         *mongo.Database
	collection *mongo.Collection
	fields     map[string]Field
}

// MongORMOperations holds the accumulated operations for a MongORM instance, including
// query filters, update documents, and other operation-specific information. This struct
// is used internally by the MongORM instance to manage the state of ongoing operations.
//
// > NOTE: This struct is not intended for public use.
type MongORM[T any] struct {
	schema     *T
	options    *MongORMOptions
	info       *MongORMInfo
	operations *MongORMOperations
}

// Creates a new MongORM instance with the provided schema and options.
//
// Example usage:
//
//	type ToDo struct {
//	  Text *string `bson:"text"`
//		// MongORM options
//		connectionString *string `mongorm:"[connection_string],connection:url"`
//		database         *string `mongorm:"[database],connection:database"`
//		collection       *string `mongorm:"[collection],connection:collection"`
//	}
//	orm := mongorm.New(&ToDo{})
func New[T any](schema *T) *MongORM[T] {
	return FromOptions(schema, nil)
}
