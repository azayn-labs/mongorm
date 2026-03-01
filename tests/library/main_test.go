package main

import (
	"testing"
	"time"

	"github.com/azayn-labs/mongorm"
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

	t.Run("Advanced fields generation", func(t *testing.T) {
		ValidateAdvancedFieldsOf(t)
	})

	t.Run("Primitive queries", func(t *testing.T) {
		ValidatePrimitiveQueries(t)
	})

	t.Run("Advanced generic queries", func(t *testing.T) {
		ValidateAdvancedGenericQueries(t)
	})

	t.Run("Create TODO", func(t *testing.T) {
		CreateLibraryTodo(t, toDo)
	})

	t.Run("Save create with Set merges fields", func(t *testing.T) {
		TestSaveCreateWithSetMergesSchemaFields(t)
	})

	t.Run("Save create with SetOnInsert merges fields", func(t *testing.T) {
		TestSaveCreateWithSetOnInsertMergesSchemaFields(t)
	})

	t.Run("Save upsert SetOnInsert insert and match noop", func(t *testing.T) {
		TestSaveWithSetOnInsertUpsertInsertAndMatchNoop(t)
	})

	t.Run("Save merges schema filters on update", func(t *testing.T) {
		TestSaveUsesSchemaFiltersForUpdate(t)
	})

	t.Run("Find TODO by ID", func(t *testing.T) {
		FindLibraryTodoByID(t, toDo.ID)
	})

	t.Run("Find TODO by text with WhereBy", func(t *testing.T) {
		FindLibraryTodoByTextWhereBy(t, *toDo.Text)
	})

	t.Run("Find TODO with OrWhere", func(t *testing.T) {
		FindLibraryTodoByOrWhere(t)
	})

	t.Run("Find TODO with OrWhereAnd", func(t *testing.T) {
		FindLibraryTodoByOrWhereAnd(t)
	})

	t.Run("Update TODO by ID", func(t *testing.T) {
		update := &ToDo{
			Text:  mongorm.String("This is an updated functional library test todo"),
			Done:  mongorm.Bool(true),
			Count: 2,
		}
		UpdateLibraryTodoByID(t, toDo.ID, update)
	})

	t.Run("Update TODO count with inc and decrement", func(t *testing.T) {
		UpdateLibraryTodoCountWithIncAndDecrement(t)
	})

	t.Run("BeforeSave hook can add updates", func(t *testing.T) {
		TestBeforeSaveCanAddUpdateOperations(t)
	})

	t.Run("BeforeCreate hook recalculates modified", func(t *testing.T) {
		TestBeforeCreateRebuildsModifiedAfterHook(t)
	})

	t.Run("BeforeUpdate hook recalculates modified", func(t *testing.T) {
		TestBeforeUpdateRebuildsModifiedAfterHook(t)
	})

	t.Run("BeforeFind hook can mutate filter", func(t *testing.T) {
		TestBeforeFindCanMutateFilter(t)
	})

	t.Run("BeforeDelete hook can mutate filter", func(t *testing.T) {
		TestBeforeDeleteCanMutateFilter(t)
	})

	t.Run("FindOneAndUpdate returns not found", func(t *testing.T) {
		UpdateLibraryTodoFindOneAndUpdateNotFound(t)
	})

	t.Run("FindOneAndUpdate merges schema and where filters", func(t *testing.T) {
		TestFindOneAndUpdateMergesSchemaAndWhereFilters(t)
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

	t.Run("Cursor All returns distinct documents", func(t *testing.T) {
		CursorAllReturnsDistinctDocuments(t)
	})

	t.Run("Cursor Current clears after exhaustion", func(t *testing.T) {
		CursorCurrentClearedAfterExhaustion(t)
	})

	t.Run("Find with sort/limit/skip/projection", func(t *testing.T) {
		FindLibraryTodoWithSortLimitSkipProjection(t)
	})

	t.Run("Find with typed projection", func(t *testing.T) {
		FindLibraryTodoWithTypedProjection(t)
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

	t.Run("Optimistic locking", func(t *testing.T) {
		ValidateLibraryOptimisticLocking(t)
	})

	t.Run("Error taxonomy", func(t *testing.T) {
		ValidateLibraryErrorTaxonomy(t)
	})

	t.Run("Versioning configuration safety", func(t *testing.T) {
		ValidateLibraryVersioningConfigSafety(t)
	})

	t.Run("Timestamps", func(t *testing.T) {
		ValidateLibraryTimestamps(t)
	})

	t.Run("Bulk write", func(t *testing.T) {
		ValidateLibraryBulkWrite(t)
	})

	t.Run("Aggregate TODOs by text", func(t *testing.T) {
		aggText := "aggregate-check-" + time.Now().Format(time.RFC3339Nano)
		CreateLibraryTodo(t, &ToDo{Text: mongorm.String(aggText), Done: mongorm.Bool(false), Count: 1})
		CreateLibraryTodo(t, &ToDo{Text: mongorm.String(aggText), Done: mongorm.Bool(true), Count: 2})
		CreateLibraryTodo(t, &ToDo{Text: mongorm.String(aggText), Done: mongorm.Bool(true), Count: 3})
		defer DeleteAllLibraryTodoByText(t, aggText)

		AggregateLibraryTodoByText(t, aggText)
		AggregateLibraryTodoGroups(t, aggText)
		AggregateLibraryTodoByBuilderOperators(t, aggText)
		AggregateLibraryTodoAddFieldsAndFacet(t, aggText)
		AggregateLibraryTodoGroupSumByBuilder(t, aggText)
	})
}
