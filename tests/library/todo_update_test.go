package main

import (
	"errors"
	"fmt"
	"testing"
	"time"

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

func UpdateLibraryTodoCountWithIncAndDecrement(t *testing.T) {
	seed := &ToDo{
		Text:  mongorm.String(fmt.Sprintf("inc-dec-%d", time.Now().UnixNano())),
		Count: 10,
	}

	CreateLibraryTodo(t, seed)
	defer DeleteLibraryTodoByID(t, seed.ID)

	incModel := mongorm.New(&ToDo{})
	incModel.Where(ToDoFields.ID.Eq(*seed.ID)).IncData(ToDoFields.Count, int64(5))
	if err := incModel.Save(t.Context()); err != nil {
		t.Fatal(err)
	}

	verifyAfterInc := &ToDo{}
	if err := mongorm.New(verifyAfterInc).Where(ToDoFields.ID.Eq(*seed.ID)).First(t.Context()); err != nil {
		t.Fatal(err)
	}

	if verifyAfterInc.Count != 15 {
		t.Fatalf("expected count 15 after increment, got %d", verifyAfterInc.Count)
	}

	decModel := mongorm.New(&ToDo{})
	decModel.Where(ToDoFields.ID.Eq(*seed.ID)).DecData(ToDoFields.Count, 3)
	if err := decModel.Save(t.Context()); err != nil {
		t.Fatal(err)
	}

	verifyAfterDec := &ToDo{}
	if err := mongorm.New(verifyAfterDec).Where(ToDoFields.ID.Eq(*seed.ID)).First(t.Context()); err != nil {
		t.Fatal(err)
	}

	if verifyAfterDec.Count != 12 {
		t.Fatalf("expected count 12 after decrement, got %d", verifyAfterDec.Count)
	}
}

func UpdateLibraryTodoFindOneAndUpdateNotFound(t *testing.T) {
	targetID := bson.NewObjectID()

	todoModel := mongorm.New(&ToDo{})
	todoModel.Where(ToDoFields.ID.Eq(targetID)).Set(&ToDo{Text: mongorm.String("should-not-upsert")})

	err := todoModel.FindOneAndUpdate(t.Context())
	if err == nil {
		t.Fatal("expected not found error")
	}

	if !errors.Is(err, mongorm.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got: %v", err)
	}
}
