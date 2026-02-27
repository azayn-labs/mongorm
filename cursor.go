package mongorm

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

// MongORMCursor provides an abstraction over the MongoDB cursor, allowing for easy
// iteration and retrieval of documents as MongORM instances. It encapsulates the MongoDB
// cursor and a reference to the MongORM instance to facilitate decoding documents into
// the appropriate schema. The cursor should be closed when done to free up resources.
//
// Example usage:
//
//	cursor, err := mongormInstance.FindAll(ctx)
//	if err != nil {
//	    // Handle error
//	} else {
//	    mongormCursor := &MongORMCursor[ToDo]{
//	        MongoCursor: cursor,
//	        m:           mongormInstance,
//	    }
//	    // Use mongormCursor to iterate over results
//	}
type MongORMCursor[T any] struct {
	MongoCursor *mongo.Cursor `json:"-"`
	m           *MongORM[T]   `json:"-"`
	current     *T            `json:"-"`
	err         error         `json:"-"`
}

// Next advances the cursor to the next document and decodes it into a new MongORM instance.
// It returns false when there are no more documents to read or when an error occurs.
// Call Err to inspect the last cursor error after iteration ends. The caller is responsible
// for closing the cursor when done. The context can be used to cancel the operation if needed.
//
// Example usage:
//
//	for mongormCursor.Next(ctx) {
//	    current := mongormCursor.Current()
//	    if current != nil {
//	        // Use current
//	    }
//	}
func (c *MongORMCursor[T]) Next(ctx context.Context) bool {
	if c == nil || c.MongoCursor == nil {
		if c != nil {
			c.current = nil
			c.err = configErrorf("cursor is nil")
		}
		return false
	}

	c.current = nil
	c.err = nil

	if !c.MongoCursor.Next(ctx) {
		c.err = normalizeError(c.MongoCursor.Err())
		return false
	}

	var u T
	if err := c.MongoCursor.Decode(&u); err != nil {
		c.err = normalizeError(err)
		return false
	}
	c.current = &u

	return true
}

// Err returns the most recent cursor error observed by Next or All.
func (c *MongORMCursor[T]) Err() error {
	if c == nil {
		return configErrorf("cursor is nil")
	}

	return c.err
}

// Current returns the current document as a MongORM instance. It should be called after
// a successful call to Next. If there is no current document (e.g., before the first call
// to Next or after reaching the end of the cursor), it returns nil.
//
// Example usage:
//
//	if mongormCursor.Next(ctx) {
//	    current := mongormCursor.Current()
//	    if current != nil {
//	        // Use current
//	    }
//	}
//	if err := mongormCursor.Err(); err != nil {
//	    // Handle error
//	}
func (c *MongORMCursor[T]) Current() *MongORM[T] {
	if c == nil || c.current == nil || c.m == nil {
		return nil
	}

	clone := c.m.clone()
	clone.schema = c.current
	return clone
}

// All retrieves all remaining documents from the cursor and decodes them into a slice
// of MongORM instances. The caller is responsible for closing the cursor when done.
// The context can be used to cancel the operation if needed. It is recommended to use
// the cursor's Next method for large result sets to avoid loading all documents into
// memory at once.
//
// Example usage:
//
//	cursors, err := mongormCursor.All(ctx)
//	if err != nil {
//	    // Handle error
//	} else {
//	    // Use the cursors
//	}
func (c *MongORMCursor[T]) All(ctx context.Context) ([]*MongORM[T], error) {
	if c == nil || c.MongoCursor == nil {
		if c != nil {
			c.current = nil
			c.err = configErrorf("cursor is nil")
		}
		return nil, configErrorf("cursor is nil")
	}

	if c.m == nil {
		c.current = nil
		c.err = configErrorf("cursor model is nil")
		return nil, c.err
	}

	c.current = nil
	c.err = nil

	var results []T
	if err := c.MongoCursor.All(ctx, &results); err != nil {
		c.err = normalizeError(err)
		return nil, c.err
	}

	clones := make([]*MongORM[T], len(results))
	for i := range results {
		clone := c.m.clone()
		clone.schema = &results[i]
		clones[i] = clone
	}

	return clones, nil
}

// Close closes the cursor and releases any resources associated with it. The context
// can be used to cancel the operation if needed. It is important to close the cursor
// when done to avoid resource leaks.
//
// Example usage:
//
//	err := mongormCursor.Close(ctx)
//	if err != nil {
//	    // Handle error
//	}
func (c *MongORMCursor[T]) Close(ctx context.Context) error {
	if c == nil || c.MongoCursor == nil {
		if c != nil {
			c.current = nil
			c.err = configErrorf("cursor is nil")
		}
		return configErrorf("cursor is nil")
	}

	c.current = nil
	c.err = nil

	return normalizeError(c.MongoCursor.Close(ctx))
}
