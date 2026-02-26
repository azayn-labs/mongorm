package main

import (
	"errors"
	"fmt"
	"testing"

	"github.com/CdTgr/mongorm"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
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

func DeleteAllLibraryTodoByText(t *testing.T, text string) {
	logger(t, fmt.Sprintf("[TODO] Deleting all by text %s\n", text))

	toDo := &ToDo{}
	todoModel := mongorm.New(toDo)
	todoModel.WhereBy(ToDoFields.Text, text)

	res, err := todoModel.DeleteMulti(t.Context())
	if err != nil {
		t.Fatal(err)
	}

	if res.DeletedCount < 1 {
		t.Fatalf("expected at least 1 deleted document, got %d", res.DeletedCount)
	}

	verify := &ToDo{}
	verifyModel := mongorm.New(verify)
	verifyModel.WhereBy(ToDoFields.Text, text)

	err = verifyModel.First(t.Context())
	if !errors.Is(err, mongo.ErrNoDocuments) {
		t.Fatalf("expected no documents after DeleteMulti, got: %v", err)
	}
}
