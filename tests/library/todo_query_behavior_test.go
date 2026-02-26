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

func CountLibraryTodoByText(t *testing.T, text string) {
	toDo := &ToDo{}
	todoModel := mongorm.New(toDo)
	todoModel.WhereBy(ToDoFields.Text, text)

	count, err := todoModel.Count(t.Context())
	if err != nil {
		t.Fatal(err)
	}

	if count < 2 {
		t.Fatalf("expected count >= 2, got %d", count)
	}
}

func DistinctLibraryTodoTextByPrefix(t *testing.T, prefix string) {
	toDo := &ToDo{}
	todoModel := mongorm.New(toDo)
	todoModel.Where(ToDoFields.Text.Reg("^" + prefix))

	values, err := todoModel.Distinct(t.Context(), ToDoFields.Text)
	if err != nil {
		t.Fatal(err)
	}

	if len(values) < 2 {
		t.Fatalf("expected at least 2 distinct values, got %d", len(values))
	}
}

func DistinctLibraryTodoTypedHelpers(t *testing.T, prefix string) {
	textModel := mongorm.New(&ToDo{})
	textModel.Where(ToDoFields.Text.Reg("^" + prefix))

	texts, err := textModel.DistinctStrings(t.Context(), ToDoFields.Text)
	if err != nil {
		t.Fatal(err)
	}

	if len(texts) < 2 {
		t.Fatalf("expected at least 2 distinct strings, got %d", len(texts))
	}

	countModel := mongorm.New(&ToDo{})
	countModel.Where(ToDoFields.Text.Reg("^" + prefix))

	counts, err := countModel.DistinctInt64(t.Context(), ToDoFields.Count)
	if err != nil {
		t.Fatal(err)
	}

	if len(counts) < 2 {
		t.Fatalf("expected at least 2 distinct counts, got %d", len(counts))
	}

	boolModel := mongorm.New(&ToDo{})
	boolModel.Where(ToDoFields.Text.Reg("^" + prefix))

	bools, err := boolModel.DistinctBool(t.Context(), ToDoFields.Done)
	if err != nil {
		t.Fatal(err)
	}

	if len(bools) < 2 {
		t.Fatalf("expected at least 2 distinct booleans, got %d", len(bools))
	}

	floatModel := mongorm.New(&ToDo{})
	floatModel.Where(ToDoFields.Text.Reg("^" + prefix))

	floats, err := floatModel.DistinctFloat64(t.Context(), ToDoFields.Count)
	if err != nil {
		t.Fatal(err)
	}

	if len(floats) < 2 {
		t.Fatalf("expected at least 2 distinct floats, got %d", len(floats))
	}

	idModel := mongorm.New(&ToDo{})
	idModel.Where(ToDoFields.Text.Reg("^" + prefix))

	ids, err := idModel.DistinctObjectIDs(t.Context(), ToDoFields.ID)
	if err != nil {
		t.Fatal(err)
	}

	if len(ids) < 2 {
		t.Fatalf("expected at least 2 distinct object ids, got %d", len(ids))
	}

	timeModel := mongorm.New(&ToDo{})
	timeModel.Where(ToDoFields.Text.Reg("^" + prefix))

	times, err := timeModel.DistinctTimes(t.Context(), ToDoFields.CreatedAt)
	if err != nil {
		t.Fatal(err)
	}

	if len(times) < 2 {
		t.Fatalf("expected at least 2 distinct times, got %d", len(times))
	}
}

func DistinctLibraryTodoGenericHelper(t *testing.T, prefix string) {
	baseModel := mongorm.New(&ToDo{})
	baseModel.Where(ToDoFields.Text.Reg("^" + prefix))

	texts, err := mongorm.DistinctFieldAs[ToDo, string](baseModel, t.Context(), ToDoFields.Text)
	if err != nil {
		t.Fatal(err)
	}
	if len(texts) < 2 {
		t.Fatalf("expected at least 2 distinct strings via DistinctAs, got %d", len(texts))
	}

	countModel := mongorm.New(&ToDo{})
	countModel.Where(ToDoFields.Text.Reg("^" + prefix))

	counts, err := mongorm.DistinctFieldAs[ToDo, int64](countModel, t.Context(), ToDoFields.Count)
	if err != nil {
		t.Fatal(err)
	}
	if len(counts) < 2 {
		t.Fatalf("expected at least 2 distinct counts via DistinctAs, got %d", len(counts))
	}

	boolModel := mongorm.New(&ToDo{})
	boolModel.Where(ToDoFields.Text.Reg("^" + prefix))

	flags, err := mongorm.DistinctFieldAs[ToDo, bool](boolModel, t.Context(), ToDoFields.Done)
	if err != nil {
		t.Fatal(err)
	}
	if len(flags) < 2 {
		t.Fatalf("expected at least 2 distinct bools via DistinctAs, got %d", len(flags))
	}

	idModel := mongorm.New(&ToDo{})
	idModel.Where(ToDoFields.Text.Reg("^" + prefix))

	ids, err := mongorm.DistinctFieldAs[ToDo, bson.ObjectID](idModel, t.Context(), ToDoFields.ID)
	if err != nil {
		t.Fatal(err)
	}
	if len(ids) < 2 {
		t.Fatalf("expected at least 2 distinct object ids via DistinctAs, got %d", len(ids))
	}

	timeModel := mongorm.New(&ToDo{})
	timeModel.Where(ToDoFields.Text.Reg("^" + prefix))

	times, err := mongorm.DistinctFieldAs[ToDo, time.Time](timeModel, t.Context(), ToDoFields.CreatedAt)
	if err != nil {
		t.Fatal(err)
	}
	if len(times) < 2 {
		t.Fatalf("expected at least 2 distinct times via DistinctAs, got %d", len(times))
	}
}

func FindLibraryTodoWithKeysetPagination(t *testing.T) {
	prefix := fmt.Sprintf("keyset-check-%d", time.Now().UnixNano())

	todos := []*ToDo{
		{Text: mongorm.String(prefix), Count: 10},
		{Text: mongorm.String(prefix), Count: 20},
		{Text: mongorm.String(prefix), Count: 30},
	}

	for _, item := range todos {
		CreateLibraryTodo(t, item)
		defer DeleteLibraryTodoByID(t, item.ID)
	}

	toDo := &ToDo{}
	todoModel := mongorm.New(toDo)
	todoModel.
		WhereBy(ToDoFields.Text, prefix).
		PaginateAfter(ToDoFields.Count, int64(10), 1)

	if err := todoModel.First(t.Context()); err != nil {
		t.Fatal(err)
	}

	if toDo.Count != 20 {
		t.Fatalf("expected first keyset page item with count 20, got %d", toDo.Count)
	}
}
