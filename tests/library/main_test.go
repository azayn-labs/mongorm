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

	t.Run("Count TODOs by text", func(t *testing.T) {
		countText := "count-check-" + time.Now().Format(time.RFC3339Nano)
		CreateLibraryTodo(t, &ToDo{Text: mongorm.String(countText), Count: 1})
		CreateLibraryTodo(t, &ToDo{Text: mongorm.String(countText), Count: 2})
		defer DeleteAllLibraryTodoByText(t, countText)

		CountLibraryTodoByText(t, countText)
	})

	t.Run("Distinct TODO texts by prefix", func(t *testing.T) {
		prefix := "distinctcheck" + time.Now().Format("20060102150405")
		base := time.Now().UTC().Truncate(time.Millisecond)
		CreateLibraryTodo(t, &ToDo{Text: mongorm.String(prefix + "-a"), Done: mongorm.Bool(false), Count: 1, CreatedAt: mongorm.Timestamp(base)})
		CreateLibraryTodo(t, &ToDo{Text: mongorm.String(prefix + "-b"), Done: mongorm.Bool(true), Count: 2, CreatedAt: mongorm.Timestamp(base.Add(1 * time.Second))})
		CreateLibraryTodo(t, &ToDo{Text: mongorm.String(prefix + "-c"), Done: mongorm.Bool(true), Count: 3, CreatedAt: mongorm.Timestamp(base.Add(2 * time.Second))})
		defer DeleteAllLibraryTodoByText(t, prefix+"-a")
		defer DeleteAllLibraryTodoByText(t, prefix+"-b")
		defer DeleteAllLibraryTodoByText(t, prefix+"-c")

		DistinctLibraryTodoTextByPrefix(t, prefix)
		DistinctLibraryTodoTypedHelpers(t, prefix)
		DistinctLibraryTodoGenericHelper(t, prefix)
	})

	t.Run("Find with keyset pagination", func(t *testing.T) {
		FindLibraryTodoWithKeysetPagination(t)
	})

	t.Run("Ensure indexes", func(t *testing.T) {
		EnsureLibraryIndexes(t)
	})

	t.Run("Transactions", func(t *testing.T) {
		ValidateLibraryTransactions(t)
	})

	t.Run("Aggregate TODOs by text", func(t *testing.T) {
		aggText := "aggregate-check-" + time.Now().Format(time.RFC3339Nano)
		CreateLibraryTodo(t, &ToDo{Text: mongorm.String(aggText), Done: mongorm.Bool(false), Count: 1})
		CreateLibraryTodo(t, &ToDo{Text: mongorm.String(aggText), Done: mongorm.Bool(true), Count: 2})
		CreateLibraryTodo(t, &ToDo{Text: mongorm.String(aggText), Done: mongorm.Bool(true), Count: 3})
		defer DeleteAllLibraryTodoByText(t, aggText)

		AggregateLibraryTodoByText(t, aggText)
		AggregateLibraryTodoGroups(t, aggText)
		AggregateLibraryTodoByBuilder(t, aggText)
		AggregateLibraryTodoGroupsByBuilder(t, aggText)
		AggregateLibraryTodoByBuilderOperators(t, aggText)
		AggregateLibraryTodoAddFieldsAndFacet(t, aggText)
		AggregateLibraryTodoGroupSumByBuilder(t, aggText)
	})
}
