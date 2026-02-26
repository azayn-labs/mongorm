package main

import (
	"testing"
	"time"

	"github.com/CdTgr/mongorm"
)

func logger(t *testing.T, message string) {
	t.Logf("LIBRARY %s\n", message)
}

func TestMain(t *testing.T) {
	var toDo = &ToDo{
		Text:  mongorm.String("This is a functional library test todo"),
		Done:  mongorm.Bool(false),
		Count: 1,
	}

	t.Run("Types helpers", func(t *testing.T) {
		ValidateTypesHelpers(t)
	})

	t.Run("Fields generation", func(t *testing.T) {
		ValidateFieldsOf(t)
	})

	t.Run("Primitive queries", func(t *testing.T) {
		ValidatePrimitiveQueries(t)
	})

	t.Run("Create TODO", func(t *testing.T) {
		CreateLibraryTodo(t, toDo)
	})

	t.Run("Find TODO by ID", func(t *testing.T) {
		FindLibraryTodoByID(t, toDo.ID)
	})

	t.Run("Find TODO by text with WhereBy", func(t *testing.T) {
		FindLibraryTodoByTextWhereBy(t, *toDo.Text)
	})

	t.Run("Update TODO by ID", func(t *testing.T) {
		update := &ToDo{
			Text:  mongorm.String("This is an updated functional library test todo"),
			Done:  mongorm.Bool(true),
			Count: 2,
		}
		UpdateLibraryTodoByID(t, toDo.ID, update)
	})

	t.Run("Delete TODO by ID", func(t *testing.T) {
		DeleteLibraryTodoByID(t, toDo.ID)
	})

	t.Run("Verify TODO deleted", func(t *testing.T) {
		FindLibraryTodoByIDExpectNotFound(t, toDo.ID)
	})

	t.Run("Create and update all TODOs", func(t *testing.T) {
		allToDo := &ToDo{
			Text:  mongorm.String("This is for SaveMulti testing"),
			Done:  mongorm.Bool(false),
			Count: 1,
		}
		CreateLibraryTodo(t, allToDo)

		UnsetLibraryTodoByID(t, allToDo.ID)

		update := &ToDo{
			Text: mongorm.String("Updated all functional todos at " + time.Now().Format(time.RFC3339)),
		}
		UpdateAllLibraryTodo(t, update)

		DeleteLibraryTodoByID(t, allToDo.ID)
	})

	t.Run("FindAll TODOs by text", func(t *testing.T) {
		cursorToDo := &ToDo{
			Text: mongorm.String("cursor-check-" + time.Now().Format(time.RFC3339Nano)),
		}
		CreateLibraryTodo(t, cursorToDo)
		FindAllLibraryTodoByText(t, *cursorToDo.Text)
		DeleteLibraryTodoByID(t, cursorToDo.ID)
	})

	t.Run("Find with sort/limit/skip/projection", func(t *testing.T) {
		FindLibraryTodoWithSortLimitSkipProjection(t)
	})

	t.Run("DeleteMulti TODOs by text", func(t *testing.T) {
		bulkText := "bulk-delete-" + time.Now().Format(time.RFC3339Nano)

		CreateLibraryTodo(t, &ToDo{Text: mongorm.String(bulkText), Count: 10})
		CreateLibraryTodo(t, &ToDo{Text: mongorm.String(bulkText), Count: 20})

		DeleteAllLibraryTodoByText(t, bulkText)
	})
}
