# Query Building

MongORM provides a fluent API for building MongoDB query filters. Filters are accumulated via method chaining and executed when you call a database operation.

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

    "github.com/CdTgr/mongorm"
    "github.com/CdTgr/mongorm/primitives"
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
        Set(update)

    if _, err := orm2.SaveMulti(ctx); err != nil {
        panic(err)
    }
}
```

---

[Back to Documentation Index](./index.md) | [README](../README.md)
