package main

import (
	"context"
	"fmt"
	"testing"

	"github.com/azayn-labs/mongorm"
)

func CreateTodo(t *testing.T, toDo *ToDo) {
	logger(t, "[TODO] Creating")
	client, err := mongorm.NewClient("mongodb://localhost:27017")
	if err != nil {
		t.Fatal(err)
	}

	modelOptions := &mongorm.MongORMOptions{
		MongoClient: client,
	}

	todoModel := mongorm.FromOptions(toDo, modelOptions)
	if err := todoModel.Save(context.TODO()); err != nil {
		t.Fatal(err)
	}

	logger(t, fmt.Sprintf("[TODO] Created: %+v\n", toDo))
}
