package mongorm

import (
	"context"
	"encoding/json"
	"fmt"
	"maps"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// applySchema applies the given document to the schema of the MongORM instance, calling
// the BeforeFinalize and AfterFinalize hooks if they are implemented by the schema. This
// function is used internally to update the schema with the results of database operations,
// ensuring that any necessary hooks are executed in the correct order. It returns an error
// if any of the hooks fail or if there is an issue applying the document to the schema.
//
// > NOTE: This method is internal only.
func (m *MongORM[T]) applySchema(doc *T) error {
	schema := any(m.schema)
	if hook, ok := schema.(BeforeFinalizeHook[T]); ok {
		if err := hook.BeforeFinalize(m); err != nil {
			return err
		}
	}

	*m.schema = *doc

	if hook, ok := schema.(AfterFinalizeHook[T]); ok {
		if err := hook.AfterFinalize(m); err != nil {
			return err
		}
	}

	// clear all operations
	m.operations.reset()

	return nil
}

// updateSchema retrieves the document with the specified ID from the database and updates
// the schema of the MongORM instance with the retrieved document. It calls the BeforeFind
// and AfterFind hooks if they are implemented by the schema. This function is used internally
// after an insert operation to ensure that the schema is updated with the newly created
// document, including any fields that may have been set by the database (such as the primary
// key). It returns an error if any of the hooks fail or if there is an issue retrieving or
// applying the document to the schema.
//
// > NOTE: This method is internal only.
func (m *MongORM[T]) updateSchema(ctx context.Context, id *bson.ObjectID) error {
	_, primaryField, err := m.getFieldByTag(ModelTagPrimary)
	if err != nil {
		return err
	}

	filter := bson.M{primaryField: id}

	schema := any(m.schema)
	if hook, ok := schema.(BeforeFindHook[T]); ok {
		if err := hook.BeforeFind(m, &filter); err != nil {
			return err
		}
	}

	var doc T
	if err := m.info.collection.FindOne(ctx, filter).Decode(&doc); err != nil {
		return err
	}

	if hook, ok := schema.(AfterFindHook[T]); ok {
		if err := hook.AfterFind(m); err != nil {
			return err
		}
	}

	return m.applySchema(&doc)
}

// insertOne inserts a new document into the collection based on the current schema of the
// MongORM instance. It applies any necessary timestamps and executes any defined hooks
// before and after the insert operation. This method is used internally by the Save
// method when inserting a new document, and it returns an error if the operation fails.
//
// Example usage:
//
//	err := mongormInstance.InsertOne(ctx)
//	if err != nil {
//	    // Handle error
//	} else {
//	    // Document inserted successfully
//	}
func (m *MongORM[T]) insertOne(ctx context.Context) error {
	schema := any(m.schema)
	if hook, ok := schema.(BeforeCreateHook[T]); ok {
		if err := hook.BeforeCreate(m); err != nil {
			return err
		}
	}

	ins, err := m.info.collection.InsertOne(ctx, m.schema)
	if err != nil {
		return err
	}

	id, ok := ins.InsertedID.(bson.ObjectID)
	if !ok {
		return fmt.Errorf("Invalid document from database: missing identifier")
	}

	if err := m.updateSchema(ctx, &id); err != nil {
		return err
	}

	if hook, ok := schema.(AfterCreateHook[T]); ok {
		if err := hook.AfterCreate(m); err != nil {
			return err
		}
	}

	return nil
}

// updateOne updates an existing document in the collection based on the provided filter
// and update documents. It applies any necessary timestamps and executes any defined hooks
// before and after the update operation. This method is used internally by the Save method
// when updating an existing document, and it returns an error if the operation fails.
//
// Example usage:
//
//	err := mongormInstance.updateOne(ctx, &filter, &update)
//	if err != nil {
//	    // Handle error
//	} else {
//	    // Document updated successfully
//	}
//
// > NOTE: This method is internal only.
func (m *MongORM[T]) updateOne(
	ctx context.Context,
	filter *bson.M,
	update *bson.M,
) error {
	var doc T

	schema := any(m.schema)
	if hook, ok := schema.(BeforeUpdateHook[T]); ok {
		if err := hook.BeforeUpdate(m, filter, update); err != nil {
			return err
		}
	}

	if err := m.info.collection.FindOneAndUpdate(
		ctx,
		filter,
		update,
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	).Decode(&doc); err != nil {
		return err
	}

	return m.applySchema(&doc)
}

// findOne retrieves a single document from the collection based on the provided filter and
// options. It applies any necessary timestamps and executes any defined hooks before and after
// the find operation. This method is used internally by various query methods to retrieve
// a single document and update the schema of the MongORM instance with the retrieved document.
// It returns an error if the operation fails or if no document is found.
//
// Example usage:
//
//	err := mongormInstance.findOne(ctx, &filter, opts...)
//	if err != nil {
//	    // Handle error
//	} else {
//	    // Document found and schema updated successfully
//	}
//
// > NOTE: This method is internal only.
func (m *MongORM[T]) findOne(
	ctx context.Context,
	filter *bson.M,
	opts ...options.Lister[options.FindOneOptions],
) error {
	schema := any(m.schema)
	if hook, ok := schema.(BeforeFindHook[T]); ok {
		if err := hook.BeforeFind(m, filter); err != nil {
			return err
		}
	}

	var doc T
	if err := m.info.collection.FindOne(
		ctx,
		filter,
		opts...,
	).Decode(&doc); err != nil {
		return err
	}

	if err := m.applySchema(&doc); err != nil {
		return err
	}

	if hook, ok := schema.(AfterFindHook[T]); ok {
		if err := hook.AfterFind(m); err != nil {
			return err
		}
	}

	return nil
}

// deleteOne deletes a single document from the collection based on the provided filter.
// It applies any necessary timestamps and executes any defined hooks before and after
// the delete operation.
//
// Example usage:
//
//	err := mongormInstance.deleteOne(ctx, &filter)
//	if err != nil {
//	    // Handle error
//	} else {
//	    // Document deleted successfully
//	}
//
// > NOTE: This method is internal only.
func (m *MongORM[T]) deleteOne(
	ctx context.Context,
	filter *bson.M,
) error {
	schema := any(m.schema)
	if hook, ok := schema.(BeforeDeleteHook[T]); ok {
		if err := hook.BeforeDelete(m, filter); err != nil {
			return err
		}
	}

	res, err := m.info.collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	if res.DeletedCount == 0 {
		return fmt.Errorf("no document found to delete")
	}

	m.reset() // clear all

	if hook, ok := schema.(AfterDeleteHook[T]); ok {
		if err := hook.AfterDelete(m); err != nil {
			return err
		}
	}

	return nil
}

// withPrimaryFilters constructs a filter that includes the primary key field based on
// the current state of the MongORM instance. It retrieves the primary field from the
// schema and checks if it exists and is not nil. If the primary field exists, it is
// added to the filter, and the method returns the updated filter along with the primary
// key value. If the primary field does not exist or is nil, the method returns the
// original filter without modification.
//
// > NOTE: This method is internal only.
func (m *MongORM[T]) withPrimaryFilters() (bson.M, *bson.ObjectID, error) {
	m.operations.fixQuery()

	filters := bson.M{}
	maps.Copy(filters, m.operations.query)

	_, primaryFieldName, err := m.getFieldByTag(ModelTagPrimary)
	if err != nil {
		return nil, nil, err
	}
	jsonSchema, err := json.Marshal(m.schema)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal schema: %v", err)
	}

	if val, ok := jsonContainsField(jsonSchema, primaryFieldName); ok && val != nil {
		str, ok := val.(string)
		if !ok {
			return nil, nil, fmt.Errorf("primary field cannot be converted to string")
		}

		id, err := bson.ObjectIDFromHex(str)
		if err != nil {
			return nil, nil, err
		}

		// Add the primary field to the filter
		filters[primaryFieldName] = id

		return filters, &id, nil
	}

	return filters, nil, nil
}
