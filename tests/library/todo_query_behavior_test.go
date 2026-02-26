package main

import (
	"errors"
	"fmt"
	"io"
	"testing"
	"time"

	"github.com/CdTgr/mongorm"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func FindLibraryTodoByTextWhereBy(t *testing.T, text string) {
	logger(t, fmt.Sprintf("[TODO] Finding by text using WhereBy: %s\n", text))

	toDo := &ToDo{}
	todoModel := mongorm.New(toDo)
	todoModel.WhereBy(ToDoFields.Text, text)

	if err := todoModel.First(t.Context()); err != nil {
		t.Fatal(err)
	}

	if toDo.Text == nil || *toDo.Text != text {
		t.Fatal("expected found todo with same text")
	}
}

func UnsetLibraryTodoByID(t *testing.T, id *bson.ObjectID) {
	logger(t, fmt.Sprintf("[TODO] Unsetting fields by id %s\n", id.Hex()))

	toDo := &ToDo{}
	todoModel := mongorm.New(toDo)
	todoModel.Where(ToDoFields.ID.Eq(*id)).Unset(&ToDo{
		Text:  mongorm.String("remove-text"),
		Done:  mongorm.Bool(true),
		Count: 1,
	})

	if err := todoModel.Save(t.Context()); err != nil {
		t.Fatal(err)
	}

	verify := &ToDo{}
	verifyModel := mongorm.New(verify)
	verifyModel.Where(ToDoFields.ID.Eq(*id))
	if err := verifyModel.First(t.Context()); err != nil {
		t.Fatal(err)
	}

	if verify.Text != nil {
		t.Fatal("expected text to be unset")
	}

	if verify.Done != nil {
		t.Fatal("expected done to be unset")
	}

	if verify.Count != 0 {
		t.Fatal("expected count to be unset and decoded as zero")
	}
}

func FindAllLibraryTodoByText(t *testing.T, text string) {
	logger(t, fmt.Sprintf("[TODO] Finding all by text %s\n", text))

	toDo := &ToDo{}
	todoModel := mongorm.New(toDo)
	todoModel.WhereBy(ToDoFields.Text, text)

	cursor, err := todoModel.FindAll(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	defer cursor.Close(t.Context())

	first, err := cursor.Next(t.Context())
	if err != nil {
		t.Fatal(err)
	}

	if first == nil || first.Document() == nil || first.Document().Text == nil || *first.Document().Text != text {
		t.Fatal("expected first cursor document with requested text")
	}

	_, err = cursor.Next(t.Context())
	if err != nil && !errors.Is(err, io.EOF) {
		t.Fatal(err)
	}
}

func FindLibraryTodoWithSortLimitSkipProjection(t *testing.T) {
	prefix := fmt.Sprintf("sorting-check-%d", time.Now().UnixNano())

	todos := []*ToDo{
		{Text: mongorm.String(prefix + "-1"), Count: 1},
		{Text: mongorm.String(prefix + "-2"), Count: 2},
		{Text: mongorm.String(prefix + "-3"), Count: 3},
	}

	for _, item := range todos {
		CreateLibraryTodo(t, item)
		defer DeleteLibraryTodoByID(t, item.ID)
	}

	toDo := &ToDo{}
	todoModel := mongorm.New(toDo)
	todoModel.
		Where(ToDoFields.Text.Reg(prefix)).
		Sort(bson.D{{Key: "count", Value: -1}}).
		Skip(1).
		Limit(1).
		Projection(bson.M{"text": 1, "count": 1})

	if err := todoModel.First(t.Context()); err != nil {
		t.Fatal(err)
	}

	if toDo.Count != 2 {
		t.Fatalf("expected count 2 after sort/skip/limit, got %d", toDo.Count)
	}

	if toDo.Done != nil {
		t.Fatal("expected done to be omitted by projection")
	}
}
