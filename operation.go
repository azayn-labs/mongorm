package mongorm

import "go.mongodb.org/mongo-driver/v2/bson"

// MongORMOperations holds the accumulated operations for a MongORM instance, including
// query filters, update documents, and other operation-specific information. This struct
// is used internally by the MongORM instance to manage the state of ongoing operations.
//
// > NOTE: This struct is not intended for public use.
type MongORMOperations struct {
	query  bson.M `json:"-"`
	update bson.M `json:"-"`
}

// Resets the MongORMOperations instance to its initial state. This is useful for reusing
// the same instance for multiple operations without retaining any previous state. This will
// clear any accumulated query filters or update documents, allowing the caller to start
// fresh with a new set of operations.
//
// > NOTE: This method is not intended for public use.
func (o *MongORMOperations) reset() {
	o.query = bson.M{}
	o.update = bson.M{}
}

// fixUpdate ensures that the update document is properly structured for MongoDB operations.
// It checks if the update document is nil and initializes it if necessary. It also removes
// any empty $set or $unset operations to prevent sending unnecessary updates to the database.
// This method should be called before executing an update operation to ensure that the
// update document is in the correct format.
//
// > NOTE: This method is not intended for public use.
func (o *MongORMOperations) fixUpdate() {
	if o.update == nil {
		o.update = bson.M{}
	}

	set, ok := o.update["$set"].(bson.M)
	if ok && len(set) == 0 {
		delete(o.update, "$set")
	}

	unset, ok := o.update["$unset"].(bson.M)
	if ok && len(unset) == 0 {
		delete(o.update, "$unset")
	}
}

// fixQuery ensures that the query document is properly initialized for MongoDB operations.
// It checks if the query document is nil and initializes it if necessary. This method should
// be called before executing a find or update operation to ensure that the query document is
// in the correct format.
//
// > NOTE: This method is not intended for public use.
func (o *MongORMOperations) fixQuery() {
	if o.query == nil {
		o.query = bson.M{}
	}
}
