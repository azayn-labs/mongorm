package main

import (
	"context"
	"testing"

	"github.com/CdTgr/mongorm"
)

func CreateTodo(t *testing.T) *ToDo {
	t.Log("Creating new TODO using options and struct")
	client, err := mongorm.NewClient("mongodb://localhost:27017")
	if err != nil {
		t.Fatal(err)
	}

	todo := &ToDo{
		Text: mongorm.String("This is an example todo created with options and struct"),
	}

	modelOptions := &mongorm.MongORMOptions{
		MongoClient: client,
	}

	todoModel := mongorm.FromOptions(todo, modelOptions)
	if err := todoModel.Save(context.TODO()); err != nil {
		t.Fatal(err)
	}

	t.Logf("TODO created with options and struct: %+v\n", todo)
	return todo
}
