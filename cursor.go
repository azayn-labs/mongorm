package mongorm

import (
	"context"
	"io"

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
}

// Next advances the cursor to the next document and decodes it into a new MongORM instance.
// It returns io.EOF when there are no more documents to read. The caller is responsible
// for closing the cursor when done. The context can be used to cancel the operation if needed.
//
// Example usage:
//
//	cursor, err := mongormCursor.Next(ctx)
//	if err != nil {
//	    if errors.Is(err, io.EOF) {
//	        // No more documents
//	    } else {
//	        // Handle error
//	    }
//	} else {
//	    // Use the cursor
//	}
func (c *MongORMCursor[T]) Next(ctx context.Context) (*MongORM[T], error) {
	if !c.MongoCursor.Next(ctx) {
		if err := c.MongoCursor.Err(); err != nil {
			return nil, err
		}
		return nil, io.EOF
	}

	var u T
	if err := c.MongoCursor.Decode(&u); err != nil {
		return nil, err
	}

	clone := c.m.clone()
	clone.schema = &u

	return clone, nil
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
	var results []T
	if err := c.MongoCursor.All(ctx, &results); err != nil {
		return nil, err
	}

	clones := make([]*MongORM[T], len(results))
	for i, result := range results {
		clone := c.m.clone()
		clone.schema = &result
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
	return c.MongoCursor.Close(ctx)
}
