package main

import (
	"testing"
	"time"

	"github.com/azayn-labs/mongorm"
)

func TestSaveUsesSchemaFiltersForUpdate(t *testing.T) {
	text := "save-schema-filter-" + time.Now().Format(time.RFC3339Nano)

	target := &ToDo{Text: mongorm.String(text), Done: mongorm.Bool(true), Count: 1}
	decoy := &ToDo{Text: mongorm.String(text), Done: mongorm.Bool(false), Count: 1}

	CreateLibraryTodo(t, target)
	CreateLibraryTodo(t, decoy)
	defer DeleteAllLibraryTodoByText(t, text)

	model := mongorm.New(&ToDo{Text: mongorm.String(text), Done: mongorm.Bool(true)})
	if err := model.Set(&ToDo{Count: 99}).Save(t.Context()); err != nil {
		t.Fatalf("expected Save(update by schema filter) to succeed, got: %v", err)
	}

	updated := model.Document()
	if updated == nil || updated.ID == nil {
		t.Fatal("expected updated document with id after Save")
	}

	if *updated.ID != *target.ID {
		t.Fatalf("expected Save to update target id=%s, got id=%s", target.ID.Hex(), updated.ID.Hex())
	}

	targetCheck := &ToDo{}
	if err := mongorm.New(targetCheck).WhereBy(ToDoFields.ID, *target.ID).First(t.Context()); err != nil {
		t.Fatal(err)
	}

	decoyCheck := &ToDo{}
	if err := mongorm.New(decoyCheck).WhereBy(ToDoFields.ID, *decoy.ID).First(t.Context()); err != nil {
		t.Fatal(err)
	}

	if targetCheck.Count != 99 {
		t.Fatalf("expected target document to be updated to count=99, got: %d", targetCheck.Count)
	}

	if decoyCheck.Count != 1 {
		t.Fatalf("expected decoy document to remain unchanged with count=1, got: %d", decoyCheck.Count)
	}

	countModel := mongorm.New(&ToDo{})
	total, err := countModel.WhereBy(ToDoFields.Text, text).Count(t.Context())
	if err != nil {
		t.Fatal(err)
	}

	if total != 2 {
		t.Fatalf("expected no extra upsert insert, total docs should stay 2, got: %d", total)
	}
}
