package main

import (
	"testing"

	"github.com/azayn-labs/mongorm"
	"github.com/azayn-labs/mongorm/primitives"
)

type todoCompatModel struct {
	History []map[string]any `bson:"history,omitempty"`

	connectionString *string `mongorm:"mongodb://localhost:27017,connection:url"`
	database         *string `mongorm:"orm-test,connection:database"`
	collection       *string `mongorm:"todo_library,connection:collection"`
}

type todoCompatSchema struct {
	History []primitives.GenericField
}

var todoCompatFields = mongorm.FieldsOf[todoCompatModel, todoCompatSchema]()

func TestArraySchemaFieldCanBeUsedInPushData(t *testing.T) {
	if len(todoCompatFields.History) == 0 {
		t.Fatal("expected generated history field entry for []primitives.GenericField schema")
	}

	if todoCompatFields.History[0].BSONName() != "history" {
		t.Fatalf("expected history bson path, got %q", todoCompatFields.History[0].BSONName())
	}

	m := mongorm.New(&todoCompatModel{})
	todoHistoryEntry := map[string]any{"id": "todo_123"}
	m.PushData(todoCompatFields.History, todoHistoryEntry)

	if !m.IsModified("history") {
		t.Fatal("expected history to be marked as modified after PushData")
	}
}
