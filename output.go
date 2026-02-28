package mongorm

import (
	"encoding/json"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// Document returns the current document schema of the MongORM instance. This is useful for
// accessing the underlying document structure for queries and updates. The returned value is
// a pointer to the schema struct, allowing the caller to modify the document fields directly.
//
// Example usage:
//
//	doc := mongormInstance.Document()
//	doc.Text = ptr("New ToDo Item")
func (m *MongORM[T]) Document() *T {
	return m.schema
}

// JSON converts the provided document struct into a map[string]any representation. This is
// useful for constructing query filters and update documents in a flexible format that can
// be easily manipulated. The method uses JSON marshaling and unmarshaling to perform the
// conversion, which allows it to handle nested structs and complex field types.
//
// Example usage:
//
//	doc := &ToDo{Text: ptr("New ToDo Item")}
//	docMap, err := mongormInstance.JSON(doc)
//	if err != nil {
//	    // Handle error
//	} else {
//	    // Use docMap for queries or updates
//	}
func (m *MongORM[T]) JSON(doc *T) (map[string]any, error) {
	b, err := json.Marshal(doc)
	if err != nil {
		return nil, err
	}

	var docMap map[string]any
	if err := json.Unmarshal(b, &docMap); err != nil {
		return nil, err
	}

	return docMap, nil
}

// GetRawQuery returns a debug copy of the currently accumulated query document.
// This is useful for logging and troubleshooting query builder state before execution.
func (m *MongORM[T]) GetRawQuery() bson.M {
	if m == nil || m.operations == nil {
		return bson.M{}
	}

	m.operations.fixQuery()
	return cloneBSONMap(m.operations.query)
}

// GetRawUpdate returns a debug copy of the currently accumulated update document.
// This is useful for logging and troubleshooting update builder state before execution.
func (m *MongORM[T]) GetRawUpdate() bson.M {
	if m == nil || m.operations == nil {
		return bson.M{}
	}

	m.operations.fixUpdate()
	return cloneBSONMap(m.operations.update)
}

func cloneBSONMap(source bson.M) bson.M {
	if source == nil {
		return bson.M{}
	}

	raw, err := bson.Marshal(source)
	if err != nil {
		copyMap := bson.M{}
		for key, value := range source {
			copyMap[key] = value
		}
		return copyMap
	}

	cloned := bson.M{}
	if err := bson.Unmarshal(raw, &cloned); err != nil {
		copyMap := bson.M{}
		for key, value := range source {
			copyMap[key] = value
		}
		return copyMap
	}

	return cloned
}
