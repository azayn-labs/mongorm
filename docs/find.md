# Finding Documents

Use `First()` (or its alias `Find()`) to retrieve a single document matching your filter. Build the filter using `Where()` with type-safe field methods from the schema.

## Find by Primary Key

The most common pattern — query by `ObjectID`:

```go
package main

import (
    "context"
    "errors"
    "fmt"

    "github.com/CdTgr/mongorm"
    "github.com/CdTgr/mongorm/primitives"
    "go.mongodb.org/mongo-driver/v2/bson"
    "go.mongodb.org/mongo-driver/v2/mongo"
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

    targetID := bson.NewObjectID() // use a real ID in production

    todo := &ToDo{}
    orm  := mongorm.New(todo)
    orm.Where(ToDoFields.ID.Eq(targetID))

    if err := orm.First(ctx); err != nil {
        if errors.Is(err, mongo.ErrNoDocuments) {
            fmt.Println("Document not found")
        } else {
            panic(err)
        }
        return
    }

    fmt.Printf("Found: %+v\n", todo)
}
```

After a successful call, the `todo` struct is populated with the data from the database.

## Find by Field Value

```go
todo := &ToDo{}
orm  := mongorm.New(todo)
orm.Where(ToDoFields.Text.Eq("Buy groceries"))

if err := orm.First(ctx); err != nil {
    panic(err)
}
```

## Find by Regex Pattern

```go
todo := &ToDo{}
orm  := mongorm.New(todo)
orm.Where(ToDoFields.Text.Reg("groceries$")) // matches text ending with "groceries"

if err := orm.First(ctx); err != nil {
    panic(err)
}
```

## Chaining Multiple Filters

Multiple `Where()` calls are combined with `$and`:

```go
orm.Where(ToDoFields.Text.Reg("buy")).
    Where(ToDoFields.CreatedAt.Gte(cutoff))
```

## Find() vs First()

`Find()` is an alias for `First()`. Both retrieve one document and populate the schema pointer.

```go
err := orm.Find(ctx)   // same as orm.First(ctx)
err := orm.First(ctx)
```

## WhereBy

`WhereBy` is a lower-level alternative that accepts any `Field` and a raw value:

```go
orm.WhereBy(ToDoFields.Text, "Buy groceries")
```

## Find All Documents

To retrieve multiple documents, use `FindAll()` which returns a cursor. See [Cursors](./cursors.md) for details.

## Count Documents

Use `Count()` to get the number of documents matching current filters.

```go
todo := &ToDo{}
orm  := mongorm.New(todo)

count, err := orm.
    Where(ToDoFields.Text.Reg("groceries")).
    Count(ctx)
if err != nil {
    panic(err)
}

fmt.Printf("Matched: %d\n", count)
```

## Distinct Values

Use `Distinct()` to return unique values for a field among matched documents.

```go
todo := &ToDo{}
orm  := mongorm.New(todo)

values, err := orm.
    Where(ToDoFields.Text.Reg("^buy")).
    Distinct(ctx, ToDoFields.Text)
if err != nil {
    panic(err)
}

fmt.Printf("Distinct text values: %d\n", len(values))
```

## Typed Distinct Helpers

Use typed helpers to avoid manual casting from `[]any`:

```go
texts, err := orm.DistinctStrings(ctx, ToDoFields.Text)
if err != nil {
    panic(err)
}

counts, err := orm.DistinctInt64(ctx, ToDoFields.Count)
if err != nil {
    panic(err)
}
```

---

[Back to Documentation Index](./index.md) | [README](../README.md)
