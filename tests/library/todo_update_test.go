package main

import (
	"fmt"
	"testing"

	"github.com/azayn-labs/mongorm"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func UpdateLibraryTodoByID(t *testing.T, id *bson.ObjectID, update *ToDo) {
	logger(t, fmt.Sprintf("[TODO] Using id %s for update\n", id.Hex()))

	toDo := &ToDo{}
	todoModel := mongorm.New(toDo)
	todoModel.Where(ToDoFields.ID.Eq(*id)).Set(update)

	if err := todoModel.Save(t.Context()); err != nil {
		t.Fatal(err)
	}

	FindLibraryTodoByID(t, id)

	logger(t, fmt.Sprintf("[TODO] Updated with ID %s\n", id.Hex()))
}

func UpdateAllLibraryTodo(t *testing.T, update *ToDo) {
	logger(t, "[TODO] Trying to update all")

	toDo := &ToDo{}
	todoModel := mongorm.New(toDo)
	todoModel.Set(update)

	if _, err := todoModel.SaveMulti(t.Context()); err != nil {
		t.Fatal(err)
	}

	logger(t, "[TODO] Updated all")
}
