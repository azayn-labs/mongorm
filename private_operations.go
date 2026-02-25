package mongorm

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func (m *MongORM[T]) applySchema(doc *T) error {
	*m.schema = *doc

	// clear all operations
	m.operations.Clear()

	return nil
}

func (m *MongORM[T]) updateSchema(ctx context.Context, id *bson.ObjectID) error {
	_, primaryField, err := m.getFieldByTag(ModelTagPrimary)
	if err != nil {
		return err
	}

	if err := m.info.collection.FindOne(ctx, bson.M{
		primaryField: id,
	}).Decode(m.schema); err != nil {
		return err
	}

	// clear all operations
	m.operations.Clear()

	return nil
}

func (m *MongORM[T]) insertOne(ctx context.Context) error {
	ins, err := m.info.collection.InsertOne(ctx, m.schema)
	if err != nil {
		return err
	}

	id, ok := ins.InsertedID.(bson.ObjectID)
	if !ok {
		return fmt.Errorf("Invalid document from database: missing identifier")
	}

	return m.updateSchema(ctx, &id)
}

func (m *MongORM[T]) updateOne(ctx context.Context, id *bson.ObjectID) error {
	var doc T
	toDo := bson.M{}

	set, ok := m.operations.update["$set"].(bson.M)
	if ok && len(set) > 0 {
		toDo["$set"] = set
	}

	unset, ok := m.operations.update["$unset"].(bson.M)
	if ok && len(unset) > 0 {
		toDo["$unset"] = unset
	}

	if len(toDo) == 0 {
		return nil
	}

	_, fieldID, err := m.getFieldByTag(ModelTagPrimary)
	if err != nil {
		return err
	}

	if err := m.info.collection.FindOneAndUpdate(
		ctx,
		bson.M{fieldID: id},
		toDo,
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	).Decode(&doc); err != nil {
		return err
	}

	return m.applySchema(&doc)
}

func (m *MongORM[T]) findOne(ctx context.Context) error {
	var doc T
	if err := m.info.collection.FindOne(ctx, m.operations.query).Decode(&doc); err != nil {
		return err
	}

	return m.applySchema(&doc)
}
