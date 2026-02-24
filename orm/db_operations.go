package orm

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func (m *Model[T]) Save(
	ctx context.Context,
) error {
	schema := any(m.clone)

	if hook, ok := schema.(BeforeSaveHook); ok {
		if err := hook.BeforeSave(); err != nil {
			return err
		}
	}

	if m.options.Timestamps {
		m.applyTimestamps()
	}

	idSchema, err := convertToStruct[BaseModel](m.schema)
	if err != nil {
		return fmt.Errorf("Model type is not extending the BaseModel")
	}

	if idSchema.ID != nil {
		if err := m.updateOne(ctx, idSchema.ID, schema); err != nil {
			return err
		}
	} else {
		// Insert new document
		if err := m.insertOne(ctx, schema); err != nil {
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

func (m *Model[T]) insertOne(ctx context.Context, document any) error {
	ins, err := m.collection.InsertOne(ctx, document)
	if err != nil {
		return err
	}

	id, ok := ins.InsertedID.(bson.ObjectID)
	if !ok {
		return mongo.ErrInvalidIndexValue
	}

	return m.updateSchema(ctx, id)
}

func (m *Model[T]) updateOne(ctx context.Context, id *bson.ObjectID, update any) error {
	var doc T
	set := m.getInformationToSet()
	unset := m.getInformationToUnset()

	toDo := bson.M{}
	if len(set) > 0 {
		toDo["$set"] = set
	}
	if len(unset) > 0 {
		toDo["$unset"] = unset
	}

	if len(toDo) == 0 {
		fmt.Println("No operations to do")
		return nil
	}

	if err := m.collection.FindOneAndUpdate(ctx, BaseModel{
		ID: id,
	}, bson.M{
		"$set": set,
	},
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	).Decode(&doc); err != nil {
		return err
	}

	return m.applySchema(&doc)
}
