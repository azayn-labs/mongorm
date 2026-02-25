package main

import (
	"context"
	"testing"

	"github.com/CdTgr/mongorm"
)

func CreateTodo(t *testing.T) *ToDo {
	t.Log("Creating new TODO using options")
	client, err := mongorm.NewClient("mongodb://localhost:27017")
	if err != nil {
		t.Fatal(err)
	}

	modelOptions := &mongorm.MongORMOptions{
		MongoClient:    client,
		CollectionName: mongorm.String("todo"),
		Timestamps:     true,
		DatabaseName:   mongorm.String("orm-test"),
	}
	toDo := &ToDo{
		Text: mongorm.String("This is an example todo created with options"),
	}

	todoModel := mongorm.FromOptions(toDo, modelOptions)
	if err := todoModel.Save(context.TODO()); err != nil {
		t.Fatal(err)
	}

	t.Logf("TODO created with options: %+v\n", toDo)
	return toDo
}
