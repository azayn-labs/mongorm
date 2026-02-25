package mongorm

import (
	"context"
	"encoding/json"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func (m *Model[T]) Save(
	ctx context.Context,
) error {
	schema := any(m.clone)
	fmt.Printf("Document %+v\n", m.clone)

	if hook, ok := schema.(BeforeSaveHook); ok {
		if err := hook.BeforeSave(); err != nil {
			return err
		}
	}

	if m.options.Timestamps {
		m.applyTimestamps()
	}

	_, fieldID, err := m.getPrimaryField()
	if err != nil {
		return err
	}
	jsonSchema, err := json.Marshal(m.schema)
	if err != nil {
		return fmt.Errorf("failed to marshal schema: %v", err)
	}

	if val, ok := jsonContainsField(jsonSchema, fieldID); ok && val != nil {
		id, err := bson.ObjectIDFromHex(val.(string))
		if err != nil {
			return err
		}

		if err := m.updateOne(ctx, &id); err != nil {
			return err
		}
	} else {
		// Insert new document
		fmt.Println("Creating new document")
		if err := m.insertOne(ctx); err != nil {
			return err
		}
	}

	if hook, ok := schema.(AfterSaveHook); ok {
		if err := hook.AfterSave(); err != nil {
			return err
		}
	}

	return nil
}

func (m *Model[T]) FindById(ctx context.Context) error {
	return m.First(ctx)
}

func (m *Model[T]) First(
	ctx context.Context,
) error {
	structFieldID, fieldID, err := m.getPrimaryField()
	if err != nil {
		return err
	}

	jsonSchema, err := json.Marshal(m.schema)
	if err != nil {
		return fmt.Errorf("failed to marshal schema: %v", err)
	}
	if val, ok := jsonContainsField(jsonSchema, fieldID); ok && val != nil {
		id, err := bson.ObjectIDFromHex(val.(string))
		if err != nil {
			return err
		}

		m.where(m.fields[structFieldID], id)
	}

	return m.findOne(ctx)
}
