package main

import (
	"fmt"
	"testing"

	"github.com/CdTgr/mongorm"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func DeleteLibraryTodoByID(t *testing.T, id *bson.ObjectID) {
	logger(t, fmt.Sprintf("[TODO] Deleting by id %s\n", id.Hex()))

	toDo := &ToDo{ID: id}
	todoModel := mongorm.New(toDo)

	if err := todoModel.Delete(t.Context()); err != nil {
		t.Fatal(err)
	}

	logger(t, fmt.Sprintf("[TODO] Deleted by ID %s\n", id.Hex()))
}
