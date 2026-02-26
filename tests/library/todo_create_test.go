package main

import (
	"fmt"
	"testing"

	"github.com/azayn-labs/mongorm"
)

func CreateLibraryTodo(t *testing.T, toDo *ToDo) {
	logger(t, "[TODO] Creating")

	todoModel := mongorm.New(toDo)
	if err := todoModel.Save(t.Context()); err != nil {
		t.Fatal(err)
	}

	if toDo.ID == nil {
		t.Fatal("expected todo id after create")
	}

	logger(t, fmt.Sprintf("[TODO] Created: %+v\n", toDo))
}
