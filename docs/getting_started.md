# Getting Started

This guide helps you launch a production-ready Go + MongoDB data layer with MongORM in minutes. You will set up strongly typed models, schema primitives, and fluent query accessors that scale cleanly from MVP to larger backend services.

## Installation

```bash
go get github.com/azayn-labs/mongorm
```

## Define a Model

A model is a regular Go struct with `bson` tags (for MongoDB field mapping) and `mongorm` tags (for ORM behavior). Pointer fields are recommended so MongORM can distinguish zero values from intentionally set values, which is especially useful for partial updates and precise write semantics.

```go
package main

import (
    "time"
    "go.mongodb.org/mongo-driver/v2/bson"
)

type ToDo struct {
    ID   *bson.ObjectID `json:"id"        bson:"_id,omitempty" mongorm:"primary"`
    Text *string        `json:"text"       bson:"text,omitempty"`

    // Optional automatic timestamps
    CreatedAt *time.Time `json:"createdAt" bson:"createdAt,omitempty" mongorm:"true,timestamp:created_at"`
    UpdatedAt *time.Time `json:"updatedAt" bson:"updatedAt,omitempty" mongorm:"true,timestamp:updated_at"`
}
```

> **Note:** The `mongorm:"primary"` tag is required on the primary key field. Without it, MongORM cannot distinguish inserts from updates.

For the full list of available `mongorm` struct tags see [Configuration](./configuration.md).

## Define the Schema

A schema struct mirrors your model but uses `primitives.FieldType` values instead of data types. It provides type-safe field references for building queries.

Field names in the schema must match the corresponding field names in the model struct exactly (case-sensitive).

```go
package main

import "github.com/azayn-labs/mongorm/primitives"

type ToDoSchema struct {
    ID        *primitives.ObjectIDField
    Text      *primitives.StringField
    Location  *primitives.GeoField
    CreatedAt *primitives.TimestampField
    UpdatedAt *primitives.TimestampField
}
```

See [Primitives](./primitives.md) for the full list of available field types.

## Generate Field Accessors

Use `FieldsOf[Model, Schema]()` to generate a populated schema instance. The returned value gives you type-safe query builder methods.

```go
package main

import "github.com/azayn-labs/mongorm"

var ToDoFields = mongorm.FieldsOf[ToDo, ToDoSchema]()
```

You can now use `ToDoFields.ID.Eq(someID)`, `ToDoFields.Text.Reg("pattern")`, etc. when calling `Where()`.

## Create an ORM Instance

### Option A — Struct Tags (static configuration)

Embed private fields with `mongorm` connection tags directly in your model:

```go
type ToDo struct {
    ID   *bson.ObjectID `bson:"_id,omitempty" mongorm:"primary"`
    Text *string        `bson:"text,omitempty"`

    // Connection configuration embedded in the struct
    connectionString *string `mongorm:"mongodb://localhost:27017,connection:url"`
    database         *string `mongorm:"mydb,connection:database"`
    collection       *string `mongorm:"todos,connection:collection"`
}

todo := &ToDo{Text: mongorm.String("Buy groceries")}
orm  := mongorm.New(todo)
```

### Option B — Options Struct (runtime configuration)

```go
import (
    "github.com/azayn-labs/mongorm"
    "go.mongodb.org/mongo-driver/v2/mongo"
    "go.mongodb.org/mongo-driver/v2/mongo/options"
)

client, _ := mongo.Connect(
    options.Client().ApplyURI("mongodb://localhost:27017"),
)

opts := &mongorm.MongORMOptions{
    Timestamps:     true,
    CollectionName: mongorm.String("todos"),
    DatabaseName:   mongorm.String("mydb"),
    MongoClient:    client,
}

todo := &ToDo{Text: mongorm.String("Buy groceries")}
orm  := mongorm.FromOptions(todo, opts)
```

### Option C — Mixed

Struct tags and `MongORMOptions` can be combined. `MongORMOptions` values take precedence.

See [Configuration](./configuration.md) for full details.

## Geo Index Defaults

For geo-enabled models, create a baseline `2dsphere` index and an optional supporting index in one call:

```go
_, err := mongorm.New(&ToDo{}).EnsureGeoDefaults(
    ctx,
    ToDoFields.Location,
    []bson.E{mongorm.Asc(ToDoFields.CreatedAt)},
)
if err != nil {
    panic(err)
}
```

## Quick Example

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

    // Create
    todo := &ToDo{Text: mongorm.String("Buy groceries")}
    orm  := mongorm.New(todo)
    if err := orm.Save(ctx); err != nil {
        panic(err)
    }
    fmt.Printf("Created: %+v\n", todo)

    // Find by ID
    found := &ToDo{}
    orm2  := mongorm.New(found)
    orm2.Where(ToDoFields.ID.Eq(*todo.ID))
    if err := orm2.First(ctx); err != nil {
        panic(err)
    }
    fmt.Printf("Found: %+v\n", found)
}
```

---

[Back to Documentation Index](./index.md) | [README](../README.md)
