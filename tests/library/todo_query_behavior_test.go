package main

import (
	"fmt"
	"testing"
	"time"

	"github.com/azayn-labs/mongorm"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type toDoProjection struct {
	Text  *string `bson:"text,omitempty"`
	Count int64   `bson:"count,omitempty"`
}

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

func FindLibraryTodoByOrWhere(t *testing.T) {
	prefix := fmt.Sprintf("or-where-%d", time.Now().UnixNano())

	a := &ToDo{Text: mongorm.String(prefix + "-a"), Done: mongorm.Bool(false), Count: 1}
	b := &ToDo{Text: mongorm.String(prefix + "-b"), Done: mongorm.Bool(true), Count: 2}
	c := &ToDo{Text: mongorm.String(prefix + "-c"), Done: mongorm.Bool(true), Count: 3}

	CreateLibraryTodo(t, a)
	CreateLibraryTodo(t, b)
	CreateLibraryTodo(t, c)
	defer DeleteLibraryTodoByID(t, a.ID)
	defer DeleteLibraryTodoByID(t, b.ID)
	defer DeleteLibraryTodoByID(t, c.ID)

	model := mongorm.New(&ToDo{})
	model.
		Where(ToDoFields.Done.Eq(true)).
		OrWhereBy(ToDoFields.Text, prefix+"-a").
		OrWhere(ToDoFields.Text.Eq(prefix + "-b")).
		SortAsc(ToDoFields.Count)

	cursor, err := model.FindAll(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	defer cursor.Close(t.Context())

	items, err := cursor.All(t.Context())
	if err != nil {
		t.Fatal(err)
	}

	if len(items) != 1 {
		t.Fatalf("expected exactly one result from Where+OrWhere query, got %d", len(items))
	}

	if items[0] == nil || items[0].Document() == nil || items[0].Document().Text == nil {
		t.Fatal("expected one decoded document from OrWhere query")
	}

	if *items[0].Document().Text != prefix+"-b" {
		t.Fatalf("expected matched text %q, got %q", prefix+"-b", *items[0].Document().Text)
	}
}

func FindLibraryTodoByOrWhereAnd(t *testing.T) {
	prefix := fmt.Sprintf("or-where-and-%d", time.Now().UnixNano())

	a := &ToDo{Text: mongorm.String(prefix + "-a"), Done: mongorm.Bool(false), Count: 1}
	b := &ToDo{Text: mongorm.String(prefix + "-b"), Done: mongorm.Bool(true), Count: 2}
	c := &ToDo{Text: mongorm.String(prefix + "-c"), Done: mongorm.Bool(true), Count: 9}

	CreateLibraryTodo(t, a)
	CreateLibraryTodo(t, b)
	CreateLibraryTodo(t, c)
	defer DeleteLibraryTodoByID(t, a.ID)
	defer DeleteLibraryTodoByID(t, b.ID)
	defer DeleteLibraryTodoByID(t, c.ID)

	model := mongorm.New(&ToDo{})
	model.
		Where(ToDoFields.Text.Reg("^"+prefix)).
		OrWhereAnd(ToDoFields.Done.Eq(false), ToDoFields.Count.Eq(int64(1))).
		OrWhereAnd(ToDoFields.Done.Eq(true), ToDoFields.Count.Eq(int64(2))).
		SortAsc(ToDoFields.Count)

	cursor, err := model.FindAll(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	defer cursor.Close(t.Context())

	items, err := cursor.All(t.Context())
	if err != nil {
		t.Fatal(err)
	}

	if len(items) != 2 {
		t.Fatalf("expected 2 results from OrWhereAnd query, got %d", len(items))
	}

	if items[0] == nil || items[1] == nil || items[0].Document() == nil || items[1].Document() == nil {
		t.Fatal("expected decoded documents from OrWhereAnd query")
	}

	if items[0].Document().Count != 1 || items[1].Document().Count != 2 {
		t.Fatalf("expected sorted counts [1, 2], got [%d, %d]", items[0].Document().Count, items[1].Document().Count)
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

	if !cursor.Next(t.Context()) {
		if err := cursor.Err(); err != nil {
			t.Fatal(err)
		}
		t.Fatal("expected at least one cursor document")
	}

	first := cursor.Current()
	if first == nil {
		t.Fatal("expected current cursor item after successful Next")
	}

	if first.Document() == nil || first.Document().Text == nil || *first.Document().Text != text {
		t.Fatal("expected first cursor document with requested text")
	}
	if first.Document().ID == nil {
		t.Fatal("expected first cursor document id")
	}

	count, err := first.WhereBy(ToDoFields.ID, *first.Document().ID).Count(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Fatalf("expected cloned cursor model count to be 1 by id, got %d", count)
	}

	if cursor.Next(t.Context()) {
		t.Fatal("expected cursor to be exhausted on second Next")
	}

	if err := cursor.Err(); err != nil {
		t.Fatal(err)
	}

	if item := cursor.Current(); item != nil {
		t.Fatal("expected current item to be nil after exhaustion")
	}
}

func CursorAllReturnsDistinctDocuments(t *testing.T) {
	prefix := fmt.Sprintf("cursor-all-distinct-%d", time.Now().UnixNano())

	first := &ToDo{Text: mongorm.String(prefix + "-a"), Count: 101}
	second := &ToDo{Text: mongorm.String(prefix + "-b"), Count: 202}

	CreateLibraryTodo(t, first)
	CreateLibraryTodo(t, second)
	defer DeleteLibraryTodoByID(t, first.ID)
	defer DeleteLibraryTodoByID(t, second.ID)

	model := mongorm.New(&ToDo{})
	model.
		Where(ToDoFields.Text.Reg("^" + prefix)).
		Sort(bson.D{{Key: ToDoFields.Count.BSONName(), Value: 1}})

	cursor, err := model.FindAll(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	defer cursor.Close(t.Context())

	items, err := cursor.All(t.Context())
	if err != nil {
		t.Fatal(err)
	}

	if len(items) != 2 {
		t.Fatalf("expected exactly 2 cursor items, got %d", len(items))
	}

	if items[0] == nil || items[1] == nil || items[0].Document() == nil || items[1].Document() == nil {
		t.Fatal("expected decoded documents from cursor.All")
	}

	if items[0].Document() == items[1].Document() {
		t.Fatal("expected distinct document pointers from cursor.All")
	}

	if items[0].Document().Count == items[1].Document().Count {
		t.Fatalf("expected distinct decoded document values, got identical counts: %d", items[0].Document().Count)
	}

	if items[0].Document().Count != 101 || items[1].Document().Count != 202 {
		t.Fatalf("expected sorted counts [101, 202], got [%d, %d]", items[0].Document().Count, items[1].Document().Count)
	}
}

func CursorCurrentClearedAfterExhaustion(t *testing.T) {
	text := fmt.Sprintf("cursor-current-clear-%d", time.Now().UnixNano())

	todo := &ToDo{Text: mongorm.String(text), Count: 1}
	CreateLibraryTodo(t, todo)
	defer DeleteLibraryTodoByID(t, todo.ID)

	model := mongorm.New(&ToDo{})
	model.WhereBy(ToDoFields.Text, text)

	cursor, err := model.FindAll(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	defer cursor.Close(t.Context())

	if !cursor.Next(t.Context()) {
		if err := cursor.Err(); err != nil {
			t.Fatal(err)
		}
		t.Fatal("expected first cursor document")
	}

	first := cursor.Current()
	if first == nil || first.Document() == nil || first.Document().Text == nil || *first.Document().Text != text {
		t.Fatal("expected current document after successful Next")
	}

	if cursor.Next(t.Context()) {
		t.Fatal("expected cursor to be exhausted on second Next")
	}

	if err := cursor.Err(); err != nil {
		t.Fatal(err)
	}

	if current := cursor.Current(); current != nil {
		t.Fatal("expected Current() to be nil after exhaustion")
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
		Sort(bson.D{{Key: ToDoFields.Count.BSONName(), Value: -1}}).
		Skip(1).
		Limit(1).
		Projection(bson.M{ToDoFields.Text.BSONName(): 1, ToDoFields.Count.BSONName(): 1})

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

func FindLibraryTodoWithTypedProjection(t *testing.T) {
	prefix := fmt.Sprintf("projection-dto-%d", time.Now().UnixNano())

	first := &ToDo{Text: mongorm.String(prefix + "-a"), Count: 11}
	second := &ToDo{Text: mongorm.String(prefix + "-b"), Count: 22}

	CreateLibraryTodo(t, first)
	CreateLibraryTodo(t, second)
	defer DeleteLibraryTodoByID(t, first.ID)
	defer DeleteLibraryTodoByID(t, second.ID)

	oneModel := mongorm.New(&ToDo{})
	oneModel.
		WhereBy(ToDoFields.Text, prefix+"-b").
		Projection(bson.M{ToDoFields.Text.BSONName(): 1, ToDoFields.Count.BSONName(): 1})

	one, err := mongorm.FindOneAs[ToDo, toDoProjection](oneModel, t.Context())
	if err != nil {
		t.Fatal(err)
	}

	if one.Text == nil || *one.Text != prefix+"-b" {
		t.Fatal("expected projected dto text for single result")
	}

	allModel := mongorm.New(&ToDo{})
	allModel.
		Where(ToDoFields.Text.Reg("^" + prefix)).
		Sort(bson.D{{Key: ToDoFields.Count.BSONName(), Value: 1}}).
		Projection(bson.M{ToDoFields.Text.BSONName(): 1, ToDoFields.Count.BSONName(): 1})

	rows, err := mongorm.FindAllAs[ToDo, toDoProjection](allModel, t.Context())
	if err != nil {
		t.Fatal(err)
	}

	if len(rows) < 2 {
		t.Fatalf("expected at least 2 projected rows, got %d", len(rows))
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

func AggregateLibraryTodoByText(t *testing.T, text string) {
	toDo := &ToDo{}
	todoModel := mongorm.New(toDo)

	pipeline := bson.A{
		bson.M{"$match": bson.M{ToDoFields.Text.BSONName(): text}},
		bson.M{"$sort": bson.M{ToDoFields.Count.BSONName(): -1}},
		bson.M{"$limit": 1},
	}

	cursor, err := todoModel.Aggregate(t.Context(), pipeline)
	if err != nil {
		t.Fatal(err)
	}
	defer cursor.Close(t.Context())

	if !cursor.Next(t.Context()) {
		if err := cursor.Err(); err != nil {
			t.Fatal(err)
		}
		t.Fatal("expected aggregate cursor item")
	}

	item := cursor.Current()
	if item == nil {
		t.Fatal("expected aggregate current item after successful Next")
	}

	if item.Document() == nil || item.Document().Count != 3 {
		t.Fatalf("expected aggregated top count 3, got %+v", item.Document())
	}
}

func AggregateLibraryTodoGroups(t *testing.T, text string) {
	totalAlias := mongorm.Alias("total")

	type GroupResult struct {
		Done  bool  `bson:"_id"`
		Total int64 `bson:"total"`
	}

	baseModel := mongorm.New(&ToDo{})
	baseModel.WhereBy(ToDoFields.Text, text)

	pipeline := bson.A{
		bson.M{"$group": bson.M{"_id": mongorm.FieldRef(ToDoFields.Done), totalAlias.Key(): bson.M{"$sum": 1}}},
		bson.M{"$sort": bson.M{"_id": 1}},
	}

	groups, err := mongorm.AggregateAs[ToDo, GroupResult](baseModel, t.Context(), pipeline)
	if err != nil {
		t.Fatal(err)
	}

	if len(groups) != 2 {
		t.Fatalf("expected 2 aggregate groups, got %d", len(groups))
	}
}

func AggregateLibraryTodoByBuilderOperators(t *testing.T, text string) {
	toDo := &ToDo{}
	todoModel := mongorm.New(toDo)

	todoModel.
		MatchExpr(ToDoFields.Text.Eq(text)).
		SortByStage(ToDoFields.Count, -1).
		LimitStage(1)

	cursor, err := todoModel.AggregatePipeline(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	defer cursor.Close(t.Context())

	if !cursor.Next(t.Context()) {
		if err := cursor.Err(); err != nil {
			t.Fatal(err)
		}
		t.Fatal("expected aggregate pipeline cursor item")
	}

	item := cursor.Current()
	if item == nil {
		t.Fatal("expected aggregate pipeline current item after successful Next")
	}

	if item.Document() == nil || item.Document().Count != 3 {
		t.Fatalf("expected operator aggregate top count 3, got %+v", item.Document())
	}
}

func AggregateLibraryTodoAddFieldsAndFacet(t *testing.T, text string) {
	rankAlias := mongorm.Alias("rank")
	doneTrueAlias := mongorm.Alias("doneTrue")
	doneFalseAlias := mongorm.Alias("doneFalse")

	type AddFieldsResult struct {
		Text  *string `bson:"text"`
		Count int64   `bson:"count"`
		Rank  int64   `bson:"rank"`
	}

	baseModel := mongorm.New(&ToDo{})
	baseModel.
		MatchBy(ToDoFields.Text, text).
		SortByStage(ToDoFields.Count, -1).
		LimitStage(1).
		AddFieldStage(rankAlias, 1)

	rows, err := mongorm.AggregatePipelineAs[ToDo, AddFieldsResult](baseModel, t.Context())
	if err != nil {
		t.Fatal(err)
	}

	if len(rows) != 1 || rows[0].Rank != 1 {
		t.Fatalf("expected one row with rank 1, got %+v", rows)
	}

	type FacetResult struct {
		DoneTrue  []bson.M `bson:"doneTrue"`
		DoneFalse []bson.M `bson:"doneFalse"`
	}

	facetModel := mongorm.New(&ToDo{})
	facetModel.
		MatchBy(ToDoFields.Text, text).
		FacetStageEntries(
			mongorm.Facet(doneTrueAlias, bson.A{
				bson.M{"$match": bson.M{ToDoFields.Done.BSONName(): true}},
				bson.M{"$limit": 1},
			}),
			mongorm.Facet(doneFalseAlias, bson.A{
				bson.M{"$match": bson.M{ToDoFields.Done.BSONName(): false}},
				bson.M{"$limit": 1},
			}),
		)

	facets, err := mongorm.AggregatePipelineAs[ToDo, FacetResult](facetModel, t.Context())
	if err != nil {
		t.Fatal(err)
	}

	if len(facets) != 1 || len(facets[0].DoneTrue) == 0 || len(facets[0].DoneFalse) == 0 {
		t.Fatalf("expected populated facet buckets, got %+v", facets)
	}
}

func AggregateLibraryTodoGroupSumByBuilder(t *testing.T, text string) {
	sumCountAlias := mongorm.Alias("sumCount")

	type SumResult struct {
		Done     bool  `bson:"_id"`
		SumCount int64 `bson:"sumCount"`
	}

	baseModel := mongorm.New(&ToDo{})
	baseModel.
		WhereBy(ToDoFields.Text, text).
		MatchWhere().
		GroupSumByAlias(ToDoFields.Done, ToDoFields.Count, sumCountAlias).
		SortStage(bson.M{"_id": 1})

	rows, err := mongorm.AggregatePipelineAs[ToDo, SumResult](baseModel, t.Context())
	if err != nil {
		t.Fatal(err)
	}

	if len(rows) != 2 {
		t.Fatalf("expected 2 grouped sums, got %d", len(rows))
	}

	total := int64(0)
	for _, row := range rows {
		total += row.SumCount
	}

	if total != 6 {
		t.Fatalf("expected grouped sum total 6, got %d", total)
	}
}

func EnsureLibraryIndexes(t *testing.T) {
	model := mongorm.New(&ToDo{})

	geoIndex, err := model.Ensure2DSphereIndex(t.Context(), ToDoFields.Location)
	if err != nil {
		t.Fatal(err)
	}
	if geoIndex == "" {
		t.Fatal("expected geo index name")
	}

	compound := mongorm.NamedIndexModelFromKeys(
		"todo_text_count_idx",
		mongorm.Asc(ToDoFields.Text),
		mongorm.Desc(ToDoFields.Count),
	)

	compoundIndex, err := model.EnsureIndex(t.Context(), compound)
	if err != nil {
		t.Fatal(err)
	}
	if compoundIndex == "" {
		t.Fatal("expected compound index name")
	}

	indexes, err := model.EnsureIndexes(t.Context(), []mongo.IndexModel{
		mongorm.IndexModelFromKeys(mongorm.Text(ToDoFields.Text)),
		mongorm.UniqueIndexModelFromKeys(mongorm.Asc(ToDoFields.Text), mongorm.Asc(ToDoFields.ID)),
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(indexes) != 2 {
		t.Fatalf("expected 2 created index names, got %d", len(indexes))
	}

	geoDefaults, err := model.EnsureGeoDefaults(
		t.Context(),
		ToDoFields.Location,
		[]bson.E{mongorm.Asc(ToDoFields.CreatedAt)},
	)
	if err != nil {
		t.Fatal(err)
	}
	if len(geoDefaults) != 2 {
		t.Fatalf("expected 2 geo default index names, got %d", len(geoDefaults))
	}
}
