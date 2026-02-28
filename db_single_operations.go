package mongorm

import (
	"context"
	"strings"

	"go.mongodb.org/mongo-driver/v2/bson"
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

	filter, _, err := m.withPrimaryAndSchemaFilters()
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
	hasExplicitUpdate := len(m.operations.update) > 0
	m.clearModified()
	m.applyTimestamps()
	m.operations.fixUpdate()

	filter, id, err := m.withPrimaryFilters()
	if err != nil {
		return err
	}

	hasSelector := len(filter) > 0

	if hasSelector && (hasExplicitUpdate || len(m.operations.query) > 0) {
		m.rebuildModifiedFromUpdate(m.operations.update)

		if hook, ok := schema.(BeforeSaveHook[T]); ok {
			if err := hook.BeforeSave(m, &filter); err != nil {
				return err
			}

			// Re-normalize updates only when hook mutations may have happened.
			m.operations.fixUpdate()
			m.rebuildModifiedFromUpdate(m.operations.update)
		}

		optimisticLockEnabled := false
		if id != nil {
			optimisticLockEnabled, err = m.applyOptimisticLock(&filter, &m.operations.update)
			if err != nil {
				return err
			}
			m.operations.fixUpdate()
			m.rebuildModifiedFromUpdate(m.operations.update)
		}

		if len(m.operations.update) == 0 {
			return configErrorf("no update operations specified")
		}

		// Save keeps upsert behavior by default, but must disable upsert when
		// optimistic locking is active to avoid stale-version inserts.
		upsertEnabled := !optimisticLockEnabled
		if len(opts) == 0 {
			opts = []options.Lister[options.FindOneAndUpdateOptions]{
				options.FindOneAndUpdate().SetUpsert(upsertEnabled),
			}
		} else {
			opts = append(opts, options.FindOneAndUpdate().SetUpsert(upsertEnabled))
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

		m.operations.fixUpdate()
		if len(m.operations.update) > 0 {
			if err := m.applyUpdateOpsToSchemaForInsert(); err != nil {
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

func (m *MongORM[T]) applyUpdateOpsToSchemaForInsert() error {
	if m == nil || m.schema == nil || len(m.operations.update) == 0 {
		return nil
	}

	raw, err := bson.Marshal(m.schema)
	if err != nil {
		return err
	}

	doc := bson.M{}
	if err := bson.Unmarshal(raw, &doc); err != nil {
		return err
	}

	if setDoc, ok := m.operations.update["$set"].(bson.M); ok {
		for path, value := range setDoc {
			setBSONPathValue(doc, path, value)
		}
	}

	if unsetDoc, ok := m.operations.update["$unset"].(bson.M); ok {
		for path := range unsetDoc {
			deleteBSONPathValue(doc, path)
		}
	}

	updatedRaw, err := bson.Marshal(doc)
	if err != nil {
		return err
	}

	var updated T
	if err := bson.Unmarshal(updatedRaw, &updated); err != nil {
		return err
	}

	*m.schema = updated
	m.operations.update = bson.M{}
	return nil
}

func setBSONPathValue(doc bson.M, path string, value any) {
	segments := splitBSONPath(path)
	if len(segments) == 0 {
		return
	}

	current := doc
	for i := 0; i < len(segments)-1; i++ {
		key := segments[i]
		next, ok := current[key].(bson.M)
		if !ok || next == nil {
			next = bson.M{}
			current[key] = next
		}
		current = next
	}

	current[segments[len(segments)-1]] = value
}

func deleteBSONPathValue(doc bson.M, path string) {
	segments := splitBSONPath(path)
	if len(segments) == 0 {
		return
	}

	current := doc
	for i := 0; i < len(segments)-1; i++ {
		next, ok := current[segments[i]].(bson.M)
		if !ok || next == nil {
			return
		}
		current = next
	}

	delete(current, segments[len(segments)-1])
}

func splitBSONPath(path string) []string {
	parts := strings.Split(strings.TrimSpace(path), ".")
	segments := make([]string, 0, len(parts))
	for _, part := range parts {
		normalized := strings.TrimSpace(part)
		if normalized == "" {
			continue
		}
		segments = append(segments, normalized)
	}

	return segments
}

// FindOneAndUpdate updates a single existing document that matches the current
// query criteria and decodes the updated document back into the schema.
//
// Unlike Save/Update, this method never performs an upsert. If no document
// matches the query criteria, it returns ErrNotFound.
//
// Example usage:
//
//	err := mongormInstance.Where(...).Set(...).FindOneAndUpdate(ctx)
//	if err != nil {
//	    // Handle error
//	}
func (m *MongORM[T]) FindOneAndUpdate(
	ctx context.Context,
	opts ...options.Lister[options.FindOneAndUpdateOptions],
) error {
	if err := m.ensureReady(); err != nil {
		return err
	}

	m.clearModified()
	m.applyTimestamps()
	m.operations.fixUpdate()

	filter, id, err := m.withPrimaryAndSchemaFilters()
	if err != nil {
		return err
	}

	if len(filter) == 0 {
		return configErrorf("findOneAndUpdate requires a filter or primary key")
	}

	m.rebuildModifiedFromUpdate(m.operations.update)

	optimisticLockEnabled := false
	if id != nil {
		optimisticLockEnabled, err = m.applyOptimisticLock(&filter, &m.operations.update)
		if err != nil {
			return err
		}
		m.rebuildModifiedFromUpdate(m.operations.update)
	}

	if len(m.operations.update) == 0 {
		return configErrorf("no update operations specified")
	}

	if len(opts) == 0 {
		opts = []options.Lister[options.FindOneAndUpdateOptions]{
			options.FindOneAndUpdate().SetUpsert(false),
		}
	} else {
		opts = append(opts, options.FindOneAndUpdate().SetUpsert(false))
	}

	return m.updateOne(
		ctx,
		&filter,
		&m.operations.update,
		optimisticLockEnabled,
		opts...,
	)
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

	filter, _, err := m.withPrimaryAndSchemaFilters()
	if err != nil {
		return err
	}

	return m.deleteOne(ctx, &filter)
}
