# Query Building

MongORM provides a high-productivity, type-safe query builder for MongoDB in Go. Filters and modifiers are composed with fluent chaining, then executed through the same ORM instance for clear, readable, and maintainable data access code.

## Where()

`Where()` adds a filter expression to the query. Multiple calls are combined with MongoDB's `$and` operator.

```go
// Signature
func (m *MongORM[T]) Where(expr bson.M) *MongORM[T]
```

Pass any `bson.M` value, or use a primitive field method to get a type-safe `bson.M`:

```go
orm.Where(ToDoFields.ID.Eq(someID))
orm.Where(ToDoFields.Text.Reg("groceries"))
orm.Where(bson.M{"status": "active"})  // raw bson.M also accepted
```

### Chaining

```go
orm.
    Where(ToDoFields.Text.Reg("buy")).
    Where(ToDoFields.CreatedAt.Gte(cutoff)).
    First(ctx)
```

Each `Where()` returns the same `*MongORM[T]` instance, so calls can be chained.

## WhereBy()

`WhereBy()` is a lower-level alternative that takes a `Field` and a raw value:

```go
// Signature
func (m *MongORM[T]) WhereBy(field Field, value any) *MongORM[T]
```

```go
orm.WhereBy(ToDoFields.Text, "Buy groceries")
// equivalent to: orm.Where(bson.M{"text": "Buy groceries"})
```

## OrWhere()

`OrWhere()` adds a filter expression to a shared MongoDB `$or` group.

```go
// Signature
func (m *MongORM[T]) OrWhere(expr bson.M) *MongORM[T]
```

```go
orm.
    OrWhere(ToDoFields.Text.Eq("Buy groceries")).
    OrWhere(ToDoFields.Text.Eq("Pay bills"))
```

## OrWhereBy()

`OrWhereBy()` is the field/value variant of `OrWhere()`.

```go
// Signature
func (m *MongORM[T]) OrWhereBy(field Field, value any) *MongORM[T]
```

```go
orm.
    OrWhereBy(ToDoFields.Text, "Buy groceries").
    OrWhereBy(ToDoFields.Text, "Pay bills")
```

## OrWhereAnd()

`OrWhereAnd()` lets you create one `$or` branch from multiple expressions combined with `$and`.
This avoids writing raw BSON field names for grouped OR logic.

```go
// Signature
func (m *MongORM[T]) OrWhereAnd(exprs ...bson.M) *MongORM[T]
```

```go
orm.
    Where(ToDoFields.Status.Eq(models.TaskRunnerStatusPending)).
    Where(ToDoFields.RunAt.Lte(now)).
    OrWhereAnd(ToDoFields.LockedUntil.NotExists()).
    OrWhereAnd(ToDoFields.LockedUntil.Lte(now))
```

For multi-condition OR branches:

```go
orm.
    OrWhereAnd(ToDoFields.Done.Eq(true), ToDoFields.Count.Gte(10)).
    OrWhereAnd(ToDoFields.Done.Eq(false), ToDoFields.Count.Lte(3))
```

## Combining Where() and OrWhere()

`Where()` clauses are grouped under `$and`, and `OrWhere()` clauses are grouped under `$or`.
When both are present, MongoDB applies both groups together (logical AND between groups).

```go
orm.
    Where(ToDoFields.Done.Eq(true)).
    OrWhere(ToDoFields.Text.Eq("Buy groceries")).
    OrWhere(ToDoFields.Text.Eq("Pay bills"))
// equivalent filter shape:
// {
//   "$and": [{"done": true}],
//   "$or":  [{"text": "Buy groceries"}, {"text": "Pay bills"}]
// }
```

## Sort()

`Sort()` sets sort order for find operations.

```go
// Signature
func (m *MongORM[T]) Sort(value any) *MongORM[T]
```

Examples:

```go
orm.Sort(bson.D{{"createdAt", -1}})   // newest first
orm.Sort(bson.M{"count": 1})          // ascending by count
```

## Limit()

`Limit()` caps the number of returned documents for `FindAll()`.

```go
// Signature
func (m *MongORM[T]) Limit(value int64) *MongORM[T]
```

```go
orm.Limit(10)
```

> `First()`/`Find()` always return one document; `Limit()` is primarily useful with `FindAll()`.

## Skip()

`Skip()` skips a number of matching documents before returning results.

```go
// Signature
func (m *MongORM[T]) Skip(value int64) *MongORM[T]
```

```go
orm.Skip(20)
```

## Projection()

`Projection()` controls which fields are returned by find operations.

```go
// Signature
func (m *MongORM[T]) Projection(value any) *MongORM[T]
```

```go
orm.Projection(bson.M{"text": 1, "count": 1})
```

## Combining Find Modifiers

```go
cursor, err := orm.
    Where(ToDoFields.Text.Reg("groceries")).
    Sort(bson.D{{"createdAt", -1}}).
    Skip(10).
    Limit(10).
    Projection(bson.M{"text": 1, "createdAt": 1}).
    FindAll(ctx)
```

## Cursor-Style Pagination Helpers

For keyset pagination, use `After()` / `Before()` or the convenience methods `PaginateAfter()` / `PaginateBefore()`.

```go
// low-level keyset filters
orm.After(ToDoFields.Count, int64(100)).Limit(20)
orm.Before(ToDoFields.Count, int64(200)).Limit(20)
```

```go
// convenience helpers (includes sort + page size)
orm.PaginateAfter(ToDoFields.Count, int64(100), 20)  // count > 100, sort asc
orm.PaginateBefore(ToDoFields.Count, int64(200), 20) // count < 200, sort desc
```

## Generic Distinct Query

When you need typed distinct values without using a dedicated helper, use `DistinctFieldAs[T, V]`:

```go
texts, err := mongorm.DistinctFieldAs[ToDo, string](orm, ctx, ToDoFields.Text)
if err != nil {
    panic(err)
}

counts, err := mongorm.DistinctFieldAs[ToDo, int64](orm, ctx, ToDoFields.Count)
if err != nil {
    panic(err)
}
```

## Set()

`Set()` specifies which fields to write during an update (`$set`). Pass a partial model struct with only the fields you want to change.

```go
// Signature
func (m *MongORM[T]) Set(value *T) *MongORM[T]
```

Rules:

- Only non-nil pointer fields and non-zero value fields are included.
- The primary key field is always skipped.
- `readonly` fields are always skipped.
- `UpdatedAt` is automatically updated (if timestamps are enabled).

```go
update := &ToDo{Text: mongorm.String("New text")}
orm.Where(ToDoFields.ID.Eq(id)).Set(update).Save(ctx)
```

## SetOnInsert()

`SetOnInsert()` adds fields to MongoDB `$setOnInsert`. Values are only written when an upsert inserts a new document.

```go
// Signature
func (m *MongORM[T]) SetOnInsert(value *T) *MongORM[T]
```

Rules:

- Only non-nil pointer fields and non-zero value fields are included.
- The primary key field is always skipped.
- `readonly` fields are always skipped.

```go
insertDefaults := &ToDo{Done: mongorm.Bool(false), Count: 1}
orm.
    WhereBy(ToDoFields.Text, "new-task").
    SetOnInsert(insertDefaults).
    Save(ctx)
```

Use `SetOnInsertData(field, value)` for direct field-based paths:

```go
orm.SetOnInsertData(ToDoFields.Text, "created-on-upsert")
orm.SetOnInsertData(ToDoFields.User.Email, "john@example.com")
```

## IncData() / DecData()

Use `IncData()` for `$inc` updates and `DecData()` for decrement operations.

```go
orm.Where(ToDoFields.ID.Eq(id)).IncData(ToDoFields.Count, int64(2)).Save(ctx)
orm.Where(ToDoFields.ID.Eq(id)).DecData(ToDoFields.Count, 1).Save(ctx)
```

## Unset()

`Unset()` removes the specified fields from the document (`$unset`).

```go
// Signature
func (m *MongORM[T]) Unset(value *T) *MongORM[T]
```

Rules:

- The primary key field is always protected.
- `readonly` fields are always protected.
- Timestamp fields (`CreatedAt`, `UpdatedAt`) are always protected.

```go
remove := &ToDo{Text: mongorm.String("placeholder")}
orm.Where(ToDoFields.ID.Eq(id)).Unset(remove).Save(ctx)
// removes the "text" field from the matched document
```

## Combining Set() and Unset()

Both can be used in the same operation:

```go
setFields   := &ToDo{Text: mongorm.String("Updated")}
unsetFields := &ToDo{SomeField: mongorm.String("x")}

orm.Where(ToDoFields.ID.Eq(id)).
    Set(setFields).
    Unset(unsetFields).
    Save(ctx)
```

## Full Example

```go
package main

import (
    "context"
    "time"

    "github.com/azayn-labs/mongorm"
    "github.com/azayn-labs/mongorm/primitives"
    "go.mongodb.org/mongo-driver/v2/bson"
)

type ToDo struct {
    ID        *bson.ObjectID `bson:"_id,omitempty"       mongorm:"primary"`
    Text      *string        `bson:"text,omitempty"`
    CreatedAt *time.Time     `bson:"createdAt,omitempty" mongorm:"true,timestamp:created_at"`
    UpdatedAt *time.Time     `bson:"updatedAt,omitempty" mongorm:"true,timestamp:updated_at"`

    connectionString *string `mongorm:"mongodb://localhost:27017,connection:url"`
    database         *string `mongorm:"mydb,connection:database"`
    collection       *string `mongorm:"todos,connection:collection"`
}

type ToDoSchema struct {
    ID        *primitives.ObjectIDField
    Text      *primitives.StringField
    CreatedAt *primitives.TimestampField
    UpdatedAt *primitives.TimestampField
}

var ToDoFields = mongorm.FieldsOf[ToDo, ToDoSchema]()

func main() {
    ctx := context.Background()
    cutoff := time.Now().Add(-7 * 24 * time.Hour) // 1 week ago

    // Find todos created in the last week containing "groceries"
    todo := &ToDo{}
    orm  := mongorm.New(todo)
    orm.
        Where(ToDoFields.Text.Reg("groceries")).
        Where(ToDoFields.CreatedAt.Gte(cutoff))

    if err := orm.First(ctx); err != nil {
        panic(err)
    }

    // Update matching todos' text
    update := &ToDo{Text: mongorm.String("Buy organic groceries")}
    orm2   := mongorm.New(&ToDo{})
    orm2.
        Where(ToDoFields.Text.Reg("groceries")).
        Where(ToDoFields.CreatedAt.Gte(cutoff)).
        Sort(bson.D{{"createdAt", -1}}).
        Limit(50).
        Set(update)

    if _, err := orm2.SaveMulti(ctx); err != nil {
        panic(err)
    }
}
```

---

[Back to Documentation Index](./index.md) | [README](../README.md)
