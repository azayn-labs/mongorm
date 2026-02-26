package main

import (
	"errors"
	"fmt"
	"testing"

	"github.com/CdTgr/mongorm"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func FindLibraryTodoByID(t *testing.T, id *bson.ObjectID) {
	logger(t, fmt.Sprintf("[TODO] Finding by id %s\n", id.Hex()))

	toDo := &ToDo{}
	todoModel := mongorm.New(toDo)
	todoModel.Where(ToDoFields.ID.Eq(*id))

	if err := todoModel.First(t.Context()); err != nil {
		t.Fatal(err)
	}

	if toDo.ID == nil || toDo.ID.Hex() != id.Hex() {
		t.Fatal("expected found todo with same id")
	}

	logger(t, fmt.Sprintf("[TODO] Found using ID %s: %+v\n", id.Hex(), toDo))
}

func FindLibraryTodoByIDExpectNotFound(t *testing.T, id *bson.ObjectID) {
	logger(t, fmt.Sprintf("[TODO] Verifying deleted id %s\n", id.Hex()))

	toDo := &ToDo{}
	todoModel := mongorm.New(toDo)
	todoModel.Where(ToDoFields.ID.Eq(*id))

	err := todoModel.First(t.Context())
	if err == nil {
		t.Fatal("expected no document after delete")
	}

	if !errors.Is(err, mongo.ErrNoDocuments) {
		t.Fatalf("expected mongo.ErrNoDocuments, got: %v", err)
	}

	logger(t, fmt.Sprintf("[TODO] Confirmed deleted ID %s\n", id.Hex()))
}
