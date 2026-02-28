# Updating Documents

MongORM supports robust MongoDB update workflows for both **single** and **multiple** documents. Use `Save()` / `Update()` for focused writes and `SaveMulti()` for efficient batch-style updates with the same fluent, type-safe developer experience.

## Update a Single Document

Use `Where()` to filter to the target document, `Set()` to specify which fields to change, then call `Save()` or `Update()`.

```go
package main

import (
    "context"
    "fmt"

    "github.com/azayn-labs/mongorm"
    "github.com/azayn-labs/mongorm/primitives"
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

## FindOneAndUpdate (No Upsert)

Use `FindOneAndUpdate()` when you want strict update-only behavior.

- It updates exactly one matching document.
- It **does not upsert**.
- It returns `ErrNotFound` when no document matches.

```go
err := mongorm.New(&ToDo{}).
    Where(ToDoFields.ID.Eq(targetID)).
    Set(&ToDo{Text: mongorm.String("Run task")}).
    FindOneAndUpdate(ctx)

if err != nil {
    if errors.Is(err, mongorm.ErrNotFound) {
        // no matching document to update
    }
    panic(err)
}
```

`FindOneAndUpdate()` also requires a selector (`Where()` or primary key on schema) and at least one update operator (`Set()`, `Unset()`, `IncData()`, etc.).

## How Set() Works

`Set()` accepts a pointer to a partial model struct. Only non-nil pointer fields and non-zero value fields are included in the `$set` operation. The primary key and `readonly` fields are always skipped.

```go
// Only Text is updated; other fields are untouched in the database
update := &ToDo{Text: mongorm.String("New text")}
orm.Set(update)
```

## Set a Single Field by Schema Field

Use `SetData(field, value)` when you want to update one field directly from your generated schema fields.

```go
orm.SetData(ToDoFields.Text, "Updated task text")
orm.SetData(ToDoFields.User.Email, "john@example.com") // nested path
```

This is especially useful for dynamic updates or deep nested fields where building a partial struct is cumbersome.

## Increment / Decrement Numeric Fields

Use `IncData(field, value)` (or alias `IncrementData`) for MongoDB `$inc` updates.
Use `DecData(field, amount)` (or alias `DecrementData`) for decrement operations.

```go
orm.
    WhereBy(ToDoFields.ID, targetID).
    IncData(ToDoFields.Count, int64(3)).
    Save(ctx)

orm.
    WhereBy(ToDoFields.ID, targetID).
    DecData(ToDoFields.Count, 1).
    Save(ctx)
```

You can also build strict field-only `$inc` documents for bulk updates:

```go
update := mongorm.IncUpdateFromPairs(
    mongorm.FieldValuePair{Field: ToDoFields.Count, Value: int64(2)},
)
```

## Array Update Operators

Use field-safe helpers for MongoDB array update operators:

```go
orm.
    WhereBy(ToDoFields.ID, targetID).
    PushData(ToDoFields.Tags, "urgent").
    AddToSetData(ToDoFields.Tags, "backend").
    PullData(ToDoFields.Tags, "deprecated").
    PopLastData(ToDoFields.Tags).
    Save(ctx)
```

Batch variants using `$each` are also available:

```go
orm.PushEachData(ToDoFields.Tags, []any{"urgent", "backend"})
orm.AddToSetEachData(ToDoFields.Tags, []any{"urgent", "backend"})
```

If you need explicit direction control for `$pop`, use:

```go
orm.PopData(ToDoFields.Tags, -1) // first
orm.PopData(ToDoFields.Tags, 1)  // last
```

You can build strict field-only update docs for bulk workflows as well:

```go
push := mongorm.PushUpdateFromPairs(
    mongorm.FieldValuePair{Field: ToDoFields.Tags, Value: "urgent"},
)

addToSet := mongorm.AddToSetUpdateFromPairs(
    mongorm.FieldValuePair{Field: ToDoFields.Tags, Value: bson.M{"$each": []any{"urgent", "backend"}}},
)

pull := mongorm.PullUpdateFromPairs(
    mongorm.FieldValuePair{Field: ToDoFields.Tags, Value: "deprecated"},
)

pop := mongorm.PopUpdateFromPairs(
    mongorm.FieldValuePair{Field: ToDoFields.Tags, Value: 1},
)
```

## Strict field-only workflow (no BSON field keys)

Use field-based helpers end-to-end:

```go
orm.
    WhereBy(ToDoFields.ID, targetID).
    SortDesc(ToDoFields.CreatedAt).
    ProjectionInclude(ToDoFields.ID, ToDoFields.Text)

builder := mongorm.NewBulkWriteBuilder[ToDo]().
    UpdateOneBy(
        ToDoFields.ID,
        targetID,
        mongorm.SetUpdateFromPairs(
            mongorm.FieldValuePair{Field: ToDoFields.Text, Value: "updated"},
            mongorm.FieldValuePair{Field: ToDoFields.User.Email, Value: "john@example.com"},
        ),
        false,
    )
```

## Unset Fields

Use `Unset()` to remove fields from the document (MongoDB `$unset`). Timestamp fields and the primary key are always protected from being unset.

```go
update := &ToDo{Text: mongorm.String("placeholder")}
orm.Unset(update) // removes the "text" field from the document
```

You can also unset a single field directly from schema fields:

```go
orm.UnsetData(ToDoFields.Text)
orm.UnsetData(ToDoFields.User.Email) // nested path
```

## Array positional updates

You can build positional paths and use them with `SetData` / `UnsetData`.

```go
import "go.mongodb.org/mongo-driver/v2/mongo/options"

// items.$[item].name = "updated"
path := mongorm.FieldPath(
    mongorm.PositionalFiltered(ToDoFields.Items, "item"),
    "name",
)

orm.SetData(path, "updated")

err := orm.Save(
    ctx,
    options.FindOneAndUpdate().SetArrayFilters(
        options.ArrayFilters{Filters: []any{bson.M{"item.id": targetID}}},
    ),
)

// items.$.name unset
orm.UnsetData(mongorm.FieldPath(mongorm.Positional(ToDoFields.Items), "name"))
```

Available helpers:

- `mongorm.Positional(field)` => `x.$`
- `mongorm.PositionalAll(field)` => `x.$[]`
- `mongorm.PositionalFiltered(field, "id")` => `x.$[id]`
- `mongorm.Indexed(field, 2)` => `x.2`
- `mongorm.FieldPath(base, "y.z")` => append nested path

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

Use `FindOneAndUpdate()` when you need strict single-document update without insert behavior. Use `SaveMulti()` when you explicitly want to update many documents and do not want insert behavior.

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
