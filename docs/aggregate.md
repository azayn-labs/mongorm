# Aggregation

MongORM supports MongoDB aggregation pipelines with both cursor-based and typed-result APIs.

## Aggregate()

`Aggregate()` runs a pipeline and decodes each output document into your model type `T`.

```go
todo := &ToDo{}
orm  := mongorm.New(todo)

pipeline := bson.A{
    bson.M{"$match": ToDoFields.Text.Reg("^todo")},
    bson.M{"$sort": bson.M{ToDoFields.Count.BSONName(): -1}},
    bson.M{"$limit": 10},
}

cursor, err := orm.Aggregate(ctx, pipeline)
if err != nil {
    panic(err)
}
defer cursor.Close(ctx)

items, err := cursor.All(ctx)
if err != nil {
    panic(err)
}
```

## AggregateRaw()

`AggregateRaw()` returns the underlying MongoDB cursor when you want manual decode flow.

```go
cursor, err := orm.AggregateRaw(ctx, pipeline)
if err != nil {
    panic(err)
}
defer cursor.Close(ctx)
```

## Fluent Pipeline Builder

You can build pipelines fluently on `MongORM` without manually constructing `bson.A`.

```go
todo := &ToDo{}
orm  := mongorm.New(todo)

cursor, err := orm.
    MatchBy(ToDoFields.Text, "todo").
    SortByStage(ToDoFields.Count, -1).
    LimitStage(10).
    AggregatePipeline(ctx)
if err != nil {
    panic(err)
}
defer cursor.Close(ctx)
```

Available stage helpers:

- `Alias(name)`
- `FieldRef(field)`
- `Facet(alias, pipeline)`
- `Pipeline(stages...)`
- `ResetPipeline()`
- `MatchStage()`
- `MatchExpr()`
- `MatchBy()`
- `MatchWhere()`
- `GroupStage()`
- `GroupCountBy()`
- `GroupCountByAlias()`
- `GroupSumBy()`
- `GroupSumByAlias()`
- `ProjectStage()`
- `ProjectIncludeFields()`
- `SortStage()`
- `SortByStage()`
- `LimitStage()`
- `SkipStage()`
- `UnwindStage()`
- `AddFieldsStage()`
- `AddFieldStage()`
- `FacetStage()`
- `FacetStageEntries()`
- `LookupStage()`
- `LookupPipelineStage()`

Typed execution for fluent pipelines:

```go
type DoneSummary struct {
    Done  bool  `bson:"_id"`
    Total int64 `bson:"total"`
}

orm := mongorm.New(&ToDo{})
totalAlias := mongorm.Alias("total")

orm.WhereBy(ToDoFields.Text, "todo").
    GroupCountByAlias(ToDoFields.Done, totalAlias)

groups, err := mongorm.AggregatePipelineAs[ToDo, DoneSummary](orm, ctx)
if err != nil {
    panic(err)
}
```

Facet aliases can be centralized too:

```go
doneTrueAlias := mongorm.Alias("doneTrue")
doneFalseAlias := mongorm.Alias("doneFalse")

orm.FacetStageEntries(
    mongorm.Facet(doneTrueAlias, bson.A{bson.M{"$match": bson.M{ToDoFields.Done.BSONName(): true}}}),
    mongorm.Facet(doneFalseAlias, bson.A{bson.M{"$match": bson.M{ToDoFields.Done.BSONName(): false}}}),
)
```

Recommended team pattern (centralized alias registry):

```go
var AggAliases = struct {
    Total     mongorm.AggregateAlias
    SumCount  mongorm.AggregateAlias
    Rank      mongorm.AggregateAlias
    DoneTrue  mongorm.AggregateAlias
    DoneFalse mongorm.AggregateAlias
}{
    Total:     mongorm.Alias("total"),
    SumCount:  mongorm.Alias("sumCount"),
    Rank:      mongorm.Alias("rank"),
    DoneTrue:  mongorm.Alias("doneTrue"),
    DoneFalse: mongorm.Alias("doneFalse"),
}

// usage
orm.GroupCountByAlias(ToDoFields.Done, AggAliases.Total)
orm.AddFieldStage(AggAliases.Rank, 1)
```

## AggregateAs

Use `AggregateAs` with type parameters to decode aggregation output into a custom result struct.

```go
type DoneSummary struct {
    Done  bool  `bson:"_id"`
    Total int64 `bson:"total"`
}

pipeline := bson.A{
    bson.M{"$group": bson.M{"_id": mongorm.FieldRef(ToDoFields.Done), "total": bson.M{"$sum": 1}}},
    bson.M{"$sort": bson.M{"_id": 1}},
}

orm := mongorm.New(&ToDo{})
groups, err := mongorm.AggregateAs[ToDo, DoneSummary](orm, ctx, pipeline)
if err != nil {
    panic(err)
}
```

## Interaction with Where()

If you call `Where()` / `WhereBy()` before `Aggregate()`, MongORM automatically prepends a `$match` stage using those filters.

```go
orm := mongorm.New(&ToDo{})
orm.WhereBy(ToDoFields.Text, "todo")

// Effective pipeline starts with: {"$match": {"$and": [{"text": "todo"}]}}
cursor, err := orm.Aggregate(ctx, bson.A{
    bson.M{"$group": bson.M{"_id": mongorm.FieldRef(ToDoFields.Done), "total": bson.M{"$sum": 1}}},
})
```

You can also reuse field operators directly in fluent aggregate stages:

```go
cursor, err := orm.
    MatchExpr(ToDoFields.Text.Eq("todo")).
    SortByStage(ToDoFields.Count, -1).
    LimitStage(10).
    AggregatePipeline(ctx)
```

`MatchStage` accepts `bson.M` intentionally because MongoDB aggregation stages can include arbitrary complex expressions.
For ergonomic reuse of existing query helpers, use `MatchExpr`, `MatchBy`, or build filters with `Where`/`WhereBy` and let MongORM prepend `$match` automatically.

## Notes

- `allowDiskUse` is enabled by default for aggregation queries.
- `Aggregate()` is best when output documents map to your model `T`.
- `AggregateAs` is best for transformed/grouped outputs.

---

[Back to Documentation Index](./index.md) | [README](../README.md)
