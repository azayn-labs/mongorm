package mongorm

import (
	"context"
	"encoding/json"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func (m *MongORM[T]) Save(
	ctx context.Context,
) error {
	schema := any(m.schema)

	if hook, ok := schema.(BeforeSaveHook); ok {
		if err := hook.BeforeSave(); err != nil {
			return err
		}
	}

	m.applyTimestamps()

	_, primaryFieldName, err := m.getFieldByTag(ModelTagPrimary)
	if err != nil {
		return err
	}
	jsonSchema, err := json.Marshal(m.schema)
	if err != nil {
		return fmt.Errorf("failed to marshal schema: %v", err)
	}

	if val, ok := jsonContainsField(jsonSchema, primaryFieldName); ok && val != nil {
		id, err := bson.ObjectIDFromHex(val.(string))
		if err != nil {
			return err
		}

		if err := m.updateOne(ctx, &id); err != nil {
			return err
		}
	} else if len(m.operations.query) == 0 {
		// Insert new document
		if err := m.insertOne(ctx); err != nil {
			return err
		}
	} else {
		// Update existing document(s) based on query
		copy := clonePtr(m.operations, false)
		if err := m.First(ctx); err != nil {
			return err
		}
		m.operations = copy
		return m.Save(ctx)
	}

	if hook, ok := schema.(AfterSaveHook); ok {
		if err := hook.AfterSave(); err != nil {
			return err
		}
	}

	return nil
}

func (m *MongORM[T]) SaveMulti(
	ctx context.Context,
	opts ...options.Lister[options.UpdateManyOptions],
) (*mongo.UpdateResult, error) {
	toDo := bson.M{}

	set, ok := m.operations.update["$set"].(bson.M)
	if ok && len(set) > 0 {
		if len(set) > 0 {
			toDo["$set"] = set
		}
	}

	unset, ok := m.operations.update["$unset"].(bson.M)
	if ok && len(unset) > 0 {
		if len(unset) > 0 {
			toDo["$unset"] = unset
		}
	}

	if len(toDo) == 0 {
		return nil, nil
	}

	res, err := m.info.collection.UpdateMany(ctx, m.operations.query, toDo, opts...)
	if err != nil {
		return nil, err
	}

	return res, err
}

func (m *MongORM[T]) FindById(ctx context.Context) error {
	return m.First(ctx)
}

func (m *MongORM[T]) First(
	ctx context.Context,
) error {
	primaryField, primaryFieldName, err := m.getFieldByTag(ModelTagPrimary)
	if err != nil {
		return err
	}

	jsonSchema, err := json.Marshal(m.schema)
	if err != nil {
		return fmt.Errorf("failed to marshal schema: %v", err)
	}
	if val, ok := jsonContainsField(jsonSchema, primaryFieldName); ok && val != nil {
		id, err := bson.ObjectIDFromHex(val.(string))
		if err != nil {
			return err
		}

		m.where(m.info.fields[primaryField], id)
	}

	return m.findOne(ctx)
}
