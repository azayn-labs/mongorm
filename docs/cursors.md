# Cursors

Use `FindAll()` to retrieve multiple documents. It returns a `*MongORMCursor[T]` that wraps the MongoDB cursor with MongORM's type system.

## Basic Usage

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

    todo := &ToDo{}
    orm  := mongorm.New(todo)

    cursor, err := orm.FindAll(ctx)
    if err != nil {
        panic(err)
    }
    defer cursor.Close(ctx)

    for cursor.Next(ctx) {
        item := cursor.Current()
        if item == nil {
            continue
        }

        fmt.Printf("Document: %+v\n", item.Document())
    }

    if err := cursor.Err(); err != nil {
        panic(err)
    }
}
```

## Filtering with FindAll()

Combine `Where()` with `FindAll()` to retrieve only matching documents:

```go
cursor, err := orm.
    Where(ToDoFields.Text.Reg("groceries")).
    FindAll(ctx)
```

## Load All Into Memory

Use `All()` on the cursor to load every result into a slice at once. Only recommended for small result sets.

```go
cursor, err := orm.FindAll(ctx)
if err != nil {
    panic(err)
}
defer cursor.Close(ctx)

items, err := cursor.All(ctx)
if err != nil {
    panic(err)
}

for _, item := range items {
    fmt.Printf("Document: %+v\n", item.Document())
}
```

## Cursor Methods

| Method | Returns | Description |
| --- | --- | --- |
| `Next(ctx)` | `bool` | Advance and decode the next document. Returns `false` when exhausted or on error. |
| `Current()` | `*MongORM[T]` | Return the current decoded document after a successful `Next(ctx)`. |
| `Err()` | `error` | Return the last cursor error after iteration ends. |
| `All(ctx)` | `([]*MongORM[T], error)` | Decode all remaining documents into a slice. |
| `Close(ctx)` | `error` | Close cursor and release server-side resources. |

## Accessing Documents

Each item returned by `Current()` or `All()` is a `*MongORM[T]` instance. Access the decoded struct via `Document()`:

```go
if cursor.Next(ctx) {
    item := cursor.Current()
    doc := item.Document()  // *ToDo
    fmt.Println(*doc.Text)
}

if err := cursor.Err(); err != nil {
    panic(err)
}
```

## Disk-Use

`FindAll()` automatically enables `allowDiskUse` on the MongoDB query, which allows large sorts and aggregations to use temporary storage rather than fail.

## Closing the Cursor

Always close the cursor when done to avoid server-side resource leaks:

```go
defer cursor.Close(ctx)
```

---

[Back to Documentation Index](./index.md) | [README](../README.md)
