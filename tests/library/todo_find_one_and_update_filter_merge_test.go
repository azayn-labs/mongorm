package main

import (
	"testing"

	"github.com/azayn-labs/mongorm"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func TestFindOneAndUpdateMergesSchemaAndWhereFilters(t *testing.T) {
	first := &ToDo{Text: mongorm.String("merge-filter-first"), Count: 7, Done: mongorm.Bool(false)}
	second := &ToDo{Text: mongorm.String("merge-filter-second"), Count: 7, Done: mongorm.Bool(false)}
	CreateLibraryTodo(t, first)
	CreateLibraryTodo(t, second)
	defer DeleteLibraryTodoByID(t, first.ID)
	defer DeleteLibraryTodoByID(t, second.ID)

	model := &ToDo{Text: mongorm.String("merge-filter-second")}
	err := mongorm.New(model).
		WhereBy(ToDoFields.Count, int64(7)).
		Set(&ToDo{Done: mongorm.Bool(true)}).
		FindOneAndUpdate(
			t.Context(),
			options.FindOneAndUpdate().SetSort(bson.D{{Key: "_id", Value: 1}}),
		)
	if err != nil {
		t.Fatalf("expected FindOneAndUpdate to succeed, got: %v", err)
	}

	firstCheck := &ToDo{}
	if err := mongorm.New(firstCheck).WhereBy(ToDoFields.ID, *first.ID).First(t.Context()); err != nil {
		t.Fatal(err)
	}

	secondCheck := &ToDo{}
	if err := mongorm.New(secondCheck).WhereBy(ToDoFields.ID, *second.ID).First(t.Context()); err != nil {
		t.Fatal(err)
	}

	if firstCheck.Done != nil && *firstCheck.Done {
		t.Fatalf("expected first document to remain unchanged, got done=%v", *firstCheck.Done)
	}

	if secondCheck.Done == nil || !*secondCheck.Done {
		t.Fatalf("expected second document to be updated by merged schema+where filter, got done=%v", secondCheck.Done)
	}
}
