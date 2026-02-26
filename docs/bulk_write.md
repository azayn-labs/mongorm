# Bulk Write

MongORM supports high-throughput batch write operations using MongoDB bulk writes. This is ideal for ingestion jobs, migration scripts, and backend tasks that need efficient multi-operation execution.

## Execute Bulk Write

Use `BulkWrite(ctx, models, opts...)` to execute one or more write models in a single request.

```go
models := mongorm.NewBulkWriteBuilder[ToDo]().
    InsertOne(&ToDo{Text: mongorm.String("a")}).
    InsertOne(&ToDo{Text: mongorm.String("b")}).
    UpdateMany(
        ToDoFields.Text.Reg("^a|b"),
        bson.M{"$set": bson.M{ToDoFields.Done.BSONName(): true}},
        false,
    ).
    DeleteOne(ToDoFields.Text.Eq("a")).
    Models()

res, err := mongorm.New(&ToDo{}).BulkWrite(
    ctx,
    models,
    options.BulkWrite().SetOrdered(false),
)
if err != nil {
    panic(err)
}

fmt.Printf("inserted=%d matched=%d modified=%d deleted=%d\n",
    res.InsertedCount,
    res.MatchedCount,
    res.ModifiedCount,
    res.DeletedCount,
)
```

## Bulk Write In Transaction

Use `BulkWriteInTransaction(ctx, models, opts...)` when the batch must be committed atomically:

```go
models := mongorm.NewBulkWriteBuilder[ToDo]().
    InsertOne(&ToDo{Text: mongorm.String("tx-a")}).
    UpdateOne(
        ToDoFields.Text.Eq("tx-a"),
        bson.M{"$set": bson.M{ToDoFields.Done.BSONName(): true}},
        false,
    ).
    Models()

_, err := mongorm.New(&ToDo{}).BulkWriteInTransaction(
    ctx,
    models,
    options.BulkWrite().SetOrdered(true),
)
if err != nil {
    panic(err)
}
```

## Builder Helpers

`NewBulkWriteBuilder[T]()` provides fluent helpers:

- `InsertOne(document)`
- `UpdateOne(filter, update, upsert)`
- `UpdateMany(filter, update, upsert)`
- `ReplaceOne(filter, replacement, upsert)`
- `DeleteOne(filter)`
- `DeleteMany(filter)`
- `Models()`

## Notes

- Bulk writes execute directly against the MongoDB collection.
- The operation uses the provided context, so it can run inside `WithTransaction`.
- At least one model is required; otherwise `BulkWrite` returns an error.

---

[Back to Documentation Index](./index.md) | [README](../README.md)
