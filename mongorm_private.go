package mongorm

import "go.mongodb.org/mongo-driver/v2/bson"

func (m *MongORM[T]) ensureReady() error {
	if m == nil {
		return configErrorf("mongorm instance is nil")
	}

	if m.initErr != nil {
		return m.initErr
	}

	if m.info == nil || m.info.collection == nil {
		return configErrorf("mongodb collection is not initialized")
	}

	return nil
}

// clone creates a deep copy of the MongORM instance, including its schema and operations.
// This is useful for creating a new instance with the same connection and collection
// information, but without any accumulated state from previous operations.
//
// > NOTE: This method is not intended for public use.
func (m *MongORM[T]) clone() *MongORM[T] {
	p := clonePtr(m, false)
	p.schema = clonePtr(m.schema, true)
	p.operations = &MongORMOperations{
		query:  bson.M{},
		update: bson.M{},
	}

	return p
}

// Resets the MongORM instance to its initial state. This is useful for reusing the same
// instance for multiple operations without retaining any previous state. This will
// still preseve the connection and collection information, but will clear any accumulated
// operations or schema information.
//
// > NOTE: This method is not intended for public use.
func (m *MongORM[T]) reset() {
	m.operations = &MongORMOperations{
		query:  bson.M{},
		update: bson.M{},
	}
	m.schema = nil
}
