package mongorm

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// Find retrieves a single document that matches the specified query criteria and decodes
// it into the schema of the MongORM instance. It uses the FindOne method of the MongoDB
// collection to execute the query and obtain the matching document. The query is constructed
// based on the state of the MongORM instance, and additional options can be provided using
// the opts parameter. It returns an error if the operation fails or if no document matches
// the query criteria.
//
// Example usage:
//
//	err := mongormInstance.Find(ctx)
//	if err != nil {
//	    if errors.Is(err, mongo.ErrNoDocuments) {
//	        // No document found
//	    } else {
//	        // Handle error
//	    }
//	} else {
//	    // Document found and decoded into mongormInstance.schema
//	}
func (m *MongORM[T]) Find(
	ctx context.Context,
	opts ...options.Lister[options.FindOneOptions],
) error {
	return m.First(ctx, opts...)
}

// First retrieves the first document that matches the specified query criteria and decodes
// it into the schema of the MongORM instance. It uses the FindOne method of the MongoDB
// collection to execute the query and obtain the matching document. The query is constructed
// based on the state of the MongORM instance, and additional options can be provided using
// the opts parameter. It returns an error if the operation fails or if no document matches
// the query criteria.
//
// Example usage:
//
//	err := mongormInstance.First(ctx)
//	if err != nil {
//	    if errors.Is(err, mongo.ErrNoDocuments) {
//	        // No document found
//	    } else {
//	        // Handle error
//	    }
//	} else {
//	    // Document found and decoded into mongormInstance.schema
//	}
func (m *MongORM[T]) First(
	ctx context.Context,
	opts ...options.Lister[options.FindOneOptions],
) error {
	if err := m.ensureReady(); err != nil {
		return err
	}

	filter, _, err := m.withPrimaryFilters()
	if err != nil {
		return err
	}

	allOpts := []options.Lister[options.FindOneOptions]{
		m.operations.findOneOptions(),
	}
	allOpts = append(allOpts, opts...)

	return m.findOne(ctx, &filter, allOpts...)
}

// Save performs an upsert operation on a single document based on the state of the
// MongORM instance. It checks if the document already exists by looking for a primary
// key or other unique identifier in the query filters. If the document exists, it updates
// the existing document with the new values from the MongORM instance. If the document
// does not exist, it inserts a new document into the collection. The method also applies
// any necessary timestamps and executes any defined hooks before and after the save
// operation. It returns an error if the operation fails.
//
// Example usage:
//
//	err := mongormInstance.Save(ctx)
//	if err != nil {
//	    // Handle error
//	} else {
//	    // Document saved successfully
//	}
func (m *MongORM[T]) Save(
	ctx context.Context,
	opts ...options.Lister[options.FindOneAndUpdateOptions],
) error {
	if err := m.ensureReady(); err != nil {
		return err
	}

	schema := any(m.schema)
	m.applyTimestamps()
	m.operations.fixUpdate()

	filter, id, err := m.withPrimaryFilters()
	if err != nil {
		return err
	}

	if id != nil || len(filter) > 0 {
		optimisticLockEnabled, err := m.applyOptimisticLock(&filter, &m.operations.update)
		if err != nil {
			return err
		}

		// Update existing document
		if hook, ok := schema.(BeforeSaveHook[T]); ok {
			if err := hook.BeforeSave(m, &filter); err != nil {
				return err
			}
		}

		if err := m.updateOne(
			ctx,
			&filter,
			&m.operations.update,
			optimisticLockEnabled,
			opts...,
		); err != nil {
			return err
		}
	} else {
		// Insert new document
		if hook, ok := schema.(BeforeSaveHook[T]); ok {
			if err := hook.BeforeSave(m, nil); err != nil {
				return err
			}
		}

		if err := m.insertOne(ctx); err != nil {
			return err
		}
	}

	if hook, ok := schema.(AfterSaveHook[T]); ok {
		if err := hook.AfterSave(m); err != nil {
			return err
		}
	}

	return nil
}

// Update is an alias for the Save method, providing a more intuitive name for updating
// an existing document. It performs the same upsert operation as Save, checking for the
// existence of the document and either updating it or inserting a new one. The method also
// applies any necessary timestamps and executes any defined hooks before and after the
// update operation. It returns an error if the operation fails.
//
// Example usage:
//
//	err := mongormInstance.Update(ctx)
//	if err != nil {
//	    // Handle error
//	} else {
//	    // Document updated successfully
//	}
func (m *MongORM[T]) Update(
	ctx context.Context,
) error {
	return m.Save(ctx)
}

// Delete removes a single document that matches the specified query criteria from the
// collection. It uses the DeleteOne method of the MongoDB collection to execute the delete
// operation. The query is constructed based on the state of the MongORM instance, and it
// typically includes a primary key or other unique identifier to ensure that only one
// document is deleted. It returns an error if the operation fails or if no document
// matches the query criteria.
//
// Example usage:
//
//	err := mongormInstance.Delete(ctx)
//	if err != nil {
//	    // Handle error
//	} else {
//	    // Document deleted successfully
//	}
func (m *MongORM[T]) Delete(
	ctx context.Context,
) error {
	if err := m.ensureReady(); err != nil {
		return err
	}

	filter, _, err := m.withPrimaryFilters()
	if err != nil {
		return err
	}

	return m.deleteOne(ctx, &filter)
}
