# Updating Documents

MongORM supports two update modes: updating a **single** document with `Save()` / `Update()`, and updating **multiple** documents with `SaveMulti()`.

## Update a Single Document

Use `Where()` to filter to the target document, `Set()` to specify which fields to change, then call `Save()` or `Update()`.

```go
package main

import (
    "context"
    "fmt"

    "github.com/CdTgr/mongorm"
    "github.com/CdTgr/mongorm/primitives"
    "go.mongodb.org/mongo-driver/v2/bson"
)

type ToDo struct {
    ID   *bson.ObjectID `bson:"_id,omitempty" mongorm:"primary"`
    Text *string        `bson:"text,omitempty"`

    connectionString *string `mongorm:"mongodb://localhost:27017,connection:url"`
    database         *string `mongorm:"mydb,connection:database"`
    collection       *string `mongorm:"todos,connection:collection"`
}

type ToDoSchema struct {
    ID   *primitives.ObjectIDField
    Text *primitives.StringField
}

var ToDoFields = mongorm.FieldsOf[ToDo, ToDoSchema]()

func main() {
    ctx := context.Background()

    targetID := bson.NewObjectID() // use a real existing ID

    update := &ToDo{
        Text: mongorm.String("Updated task text"),
    }

    todo := &ToDo{}
    orm  := mongorm.New(todo)
    orm.Where(ToDoFields.ID.Eq(targetID)).Set(update)

    if err := orm.Save(ctx); err != nil {
        panic(err)
    }

    fmt.Println("Document updated")
}
```

`Update()` is an alias for `Save()` and behaves identically:

```go
if err := orm.Update(ctx); err != nil {
    panic(err)
}
```

## How Set() Works

`Set()` accepts a pointer to a partial model struct. Only non-nil pointer fields and non-zero value fields are included in the `$set` operation. The primary key and `readonly` fields are always skipped.

```go
// Only Text is updated; other fields are untouched in the database
update := &ToDo{Text: mongorm.String("New text")}
orm.Set(update)
```

## Unset Fields

Use `Unset()` to remove fields from the document (MongoDB `$unset`). Timestamp fields and the primary key are always protected from being unset.

```go
update := &ToDo{Text: mongorm.String("placeholder")}
orm.Unset(update) // removes the "text" field from the document
```

## Empty Update Guard

For single-document updates (`Save()` / `Update()`), MongORM returns an explicit configuration error when no update operators are present.

This can happen if you target an existing document (for example via `Where()` or primary key) but do not call `Set()` / `Unset()` and no automatic `UpdatedAt` field is available.

- Error: `mongorm: invalid configuration: no update operations specified`

`SaveMulti()` keeps its own validation and also requires at least one update operator.

## Update Multiple Documents

Use `SaveMulti()` to apply a `Set()` to every document matching the `Where()` filter. It returns a `*mongo.UpdateResult` with match and modification counts.

```go
// Update the text of ALL todos
update := &ToDo{Text: mongorm.String("Mass updated")}

todo := &ToDo{}
orm  := mongorm.New(todo)
orm.Set(update)

result, err := orm.SaveMulti(ctx)
if err != nil {
    panic(err)
}
fmt.Printf("Matched: %d, Modified: %d\n", result.MatchedCount, result.ModifiedCount)
```

You can combine with `Where()` to scope the update:

```go
orm.Where(ToDoFields.Text.Reg("old text")).Set(update)
result, err := orm.SaveMulti(ctx)
```

## With Timestamps

If timestamps are enabled, `UpdatedAt` is automatically set to `time.Now()` on every `Save()` / `Update()` call. `CreatedAt` is never modified after the initial insert.

See [Timestamps](./timestamps.md) for more details.

## Upsert Behaviour

`Save()` (and `Update()`) performs an **upsert**:

- If a matching document is found (via `Where()` filters or a set primary key), it is updated.
- If no matching document exists and no filter is set, a new document is inserted.

Use `SaveMulti()` when you explicitly want to update many documents and do not want insert behaviour.

## Optimistic Locking with `_version`

Add a version field with `mongorm:"version"` to enable optimistic locking on single-document updates:

```go
type ToDo struct {
    ID      *bson.ObjectID `bson:"_id,omitempty" mongorm:"primary"`
    Version int64          `bson:"_version,omitempty" mongorm:"version"`
    Text    *string        `bson:"text,omitempty"`
}
```

Behavior:

- On insert, if `_version` is unset, nil, or `<= 0`, MongORM initializes it to `1`. If you provide a positive `_version` value explicitly, that value is preserved.
- On update (`Save()` / `Update()`), MongORM matches by current `_version` and increments it atomically.
- If the version is stale, update fails with `ErrOptimisticLockConflict`.

---

[Back to Documentation Index](./index.md) | [README](../README.md)
