# Creating Documents

Use `Save()` to insert a new document or update an existing one (upsert). When no primary key is set and no `Where()` filter is active, MongORM performs an insert.

## Basic Insert

```go
package main

import (
    "context"
    "fmt"

    "github.com/azayn-labs/mongorm"
    "go.mongodb.org/mongo-driver/v2/bson"
)

// Model with connection embedded in struct tags
type ToDo struct {
    ID   *bson.ObjectID `bson:"_id,omitempty" mongorm:"primary"`
    Text *string        `bson:"text,omitempty"`

    connectionString *string `mongorm:"mongodb://localhost:27017,connection:url"`
    database         *string `mongorm:"mydb,connection:database"`
    collection       *string `mongorm:"todos,connection:collection"`
}

func main() {
    ctx := context.Background()

    todo := &ToDo{
        Text: mongorm.String("Buy groceries"),
    }

    orm := mongorm.New(todo)
    if err := orm.Save(ctx); err != nil {
        panic(err)
    }

    fmt.Printf("Created document with ID: %s\n", todo.ID.Hex())
}
```

After a successful insert, the `todo.ID` field is populated with the new `bson.ObjectID`.

## Insert Using Options

Use `FromOptions` when you want to configure the connection at runtime rather than via struct tags.

```go
import (
    "github.com/azayn-labs/mongorm"
    "go.mongodb.org/mongo-driver/v2/mongo"
    "go.mongodb.org/mongo-driver/v2/mongo/options"
)

client, err := mongo.Connect(
    options.Client().ApplyURI("mongodb://localhost:27017"),
)
if err != nil {
    panic(err)
}

opts := &mongorm.MongORMOptions{
    Timestamps:     true,
    CollectionName: mongorm.String("todos"),
    DatabaseName:   mongorm.String("mydb"),
    MongoClient:    client,
}

todo := &ToDo{Text: mongorm.String("Walk the dog")}
orm  := mongorm.FromOptions(todo, opts)

if err := orm.Save(ctx); err != nil {
    panic(err)
}
```

## With Timestamps

Add `CreatedAt` and `UpdatedAt` fields to your model and enable timestamps either via `MongORMOptions.Timestamps = true` or by tagging the fields. MongORM will set `CreatedAt` on first insert and update `UpdatedAt` on every save.

```go
import "time"

type ToDo struct {
    ID        *bson.ObjectID `bson:"_id,omitempty"   mongorm:"primary"`
    Text      *string        `bson:"text,omitempty"`
    CreatedAt *time.Time     `bson:"createdAt,omitempty" mongorm:"true,timestamp:created_at"`
    UpdatedAt *time.Time     `bson:"updatedAt,omitempty" mongorm:"true,timestamp:updated_at"`

    connectionString *string `mongorm:"mongodb://localhost:27017,connection:url"`
    database         *string `mongorm:"mydb,connection:database"`
    collection       *string `mongorm:"todos,connection:collection"`
}

todo := &ToDo{Text: mongorm.String("Buy groceries")}
orm  := mongorm.New(todo)

if err := orm.Save(ctx); err != nil {
    panic(err)
}

fmt.Printf("Created at: %s\n", todo.CreatedAt)
```

See [Timestamps](./timestamps.md) for more details.

## Using Utility Type Helpers

MongORM provides pointer helper functions to conveniently set field values:

```go
todo := &ToDo{
    Text: mongorm.String("Buy groceries"),
}
```

See [Utility Types](./types.md) for the full list of helpers.

---

[Back to Documentation Index](./index.md) | [README](../README.md)
