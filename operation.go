package mongorm

import (
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// MongORMOperations holds the accumulated operations for a MongORM instance, including
// query filters, update documents, and other operation-specific information. This struct
// is used internally by the MongORM instance to manage the state of ongoing operations.
//
// > NOTE: This struct is not intended for public use.
type MongORMOperations struct {
	query      bson.M `json:"-"`
	update     bson.M `json:"-"`
	sort       any    `json:"-"`
	projection any    `json:"-"`
	limit      *int64 `json:"-"`
	skip       *int64 `json:"-"`
	pipeline   bson.A `json:"-"`
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
	o.sort = nil
	o.projection = nil
	o.limit = nil
	o.skip = nil
	o.pipeline = nil
}

// fixUpdate ensures that the update document is properly structured for MongoDB operations.
// It checks if the update document is nil and initializes it if necessary. It also removes
// any empty update operators (such as $set, $unset, $inc, $push, $addToSet, $pull, and $pop)
// to prevent sending unnecessary updates to the database.
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

	setOnInsert, ok := o.update["$setOnInsert"].(bson.M)
	if ok && len(setOnInsert) == 0 {
		delete(o.update, "$setOnInsert")
	}

	unset, ok := o.update["$unset"].(bson.M)
	if ok && len(unset) == 0 {
		delete(o.update, "$unset")
	}

	inc, ok := o.update["$inc"].(bson.M)
	if ok && len(inc) == 0 {
		delete(o.update, "$inc")
	}

	push, ok := o.update["$push"].(bson.M)
	if ok && len(push) == 0 {
		delete(o.update, "$push")
	}

	addToSet, ok := o.update["$addToSet"].(bson.M)
	if ok && len(addToSet) == 0 {
		delete(o.update, "$addToSet")
	}

	pull, ok := o.update["$pull"].(bson.M)
	if ok && len(pull) == 0 {
		delete(o.update, "$pull")
	}

	pop, ok := o.update["$pop"].(bson.M)
	if ok && len(pop) == 0 {
		delete(o.update, "$pop")
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

func (o *MongORMOperations) findOptions() options.Lister[options.FindOptions] {
	findOpts := options.Find()

	if o.sort != nil {
		findOpts.SetSort(o.sort)
	}

	if o.projection != nil {
		findOpts.SetProjection(o.projection)
	}

	if o.limit != nil {
		findOpts.SetLimit(*o.limit)
	}

	if o.skip != nil {
		findOpts.SetSkip(*o.skip)
	}

	return findOpts
}

func (o *MongORMOperations) findOneOptions() options.Lister[options.FindOneOptions] {
	findOneOpts := options.FindOne()

	if o.sort != nil {
		findOneOpts.SetSort(o.sort)
	}

	if o.projection != nil {
		findOneOpts.SetProjection(o.projection)
	}

	if o.skip != nil {
		findOneOpts.SetSkip(*o.skip)
	}

	return findOneOpts
}
