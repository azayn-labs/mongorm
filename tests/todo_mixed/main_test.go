package main

import (
	"testing"

	"github.com/azayn-labs/mongorm"
)

func logger(t *testing.T, message string) {
	t.Logf("TODO [struct options] %s\n", message)
}

func TestMain(t *testing.T) {
	var toDo = &ToDo{
		Text: mongorm.String("This is an example todo created with struct and options"),
	}

	t.Run("Create TODO", func(t *testing.T) {
		CreateTodo(t, toDo)
	})
}
