# MongORM

MongORM is a lightweight, type-safe ORM library for MongoDB in Go. It provides a convenient way to interact with MongoDB databases using Go structs, generics, and a fluent method API.

## Features

- Type-safe model and schema definitions using Go generics
- Fluent API for building queries with `Where()`, `Sort()`, `Limit()`, `Skip()`, `Projection()`, keyset pagination helpers, `Set()`, and `Unset()`
- Full CRUD support: create, find, update (single and multi), and delete (single and multi)
- Aggregation support: raw pipelines and fluent stage builder via `Aggregate()`, `AggregateRaw()`, `AggregateAs[T, R]()`, and `AggregatePipeline()`
- Query utilities: `Count()`, `Distinct()`, `DistinctFieldAs[T, V]()`, `DistinctStrings()`, `DistinctInt64()`, `DistinctBool()`, `DistinctFloat64()`, `DistinctObjectIDs()`, and `DistinctTimes()`
- Geospatial support: `GeoField` with `Near`, `Within`, and `Intersects` query helpers
- Index support: field-driven builders, `Ensure2DSphereIndex()`, and `EnsureGeoDefaults()`
- Lifecycle hooks for every operation (Before/After Create, Save, Update, Find, Delete, Finalize)
- Automatic `CreatedAt` / `UpdatedAt` timestamp management
- Flexible configuration: struct tags, options struct, or both
- Cursor-based iteration for large result sets
- Connection pooling — clients are reused by connection string
- Lightweight with a single dependency: the official MongoDB Go driver

## Installation

```bash
go get github.com/CdTgr/mongorm
```

## Quick Start

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

    // Connection embedded in struct tags
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
    fmt.Printf("Created: %s\n", todo.ID.Hex())

    // Find by ID
    found := &ToDo{}
    mongorm.New(found).Where(ToDoFields.ID.Eq(*todo.ID)).First(ctx)
    fmt.Printf("Found: %s\n", *found.Text)

    // Update
    mongorm.New(&ToDo{}).
        Where(ToDoFields.ID.Eq(*todo.ID)).
        Set(&ToDo{Text: mongorm.String("Buy organic groceries")}).
        Save(ctx)

    // Delete
    mongorm.New(&ToDo{}).Where(ToDoFields.ID.Eq(*todo.ID)).Delete(ctx)
}
```

## Documentation

Full documentation is in the [`docs/`](./docs/index.md) folder.

| Topic | Description |
| --- | --- |
| [Getting Started](./docs/getting_started.md) | Installation, model definition, schema setup |
| [Configuration](./docs/configuration.md) | Struct tags, options struct, mixed mode |
| [Creating Documents](./docs/create.md) | Inserting with `Save()` |
| [Finding Documents](./docs/find.md) | Querying with `First()` / `Find()`, `Count()`, and `Distinct()` |
| [Updating Documents](./docs/update.md) | Single and bulk updates |
| [Deleting Documents](./docs/delete.md) | Removing documents |
| [Indexes](./docs/indexes.md) | Field-based index builders and geo index setup |
| [Aggregation](./docs/aggregate.md) | Aggregation pipelines with fluent builder and typed decoding |
| [Cursors](./docs/cursors.md) | Iterating with `FindAll()` |
| [Query Building](./docs/query_building.md) | `Where()`, find modifiers, pagination helpers, `Set()`, `Unset()` |
| [Primitives](./docs/primitives.md) | Type-safe field query methods |
| [Hooks](./docs/hooks.md) | Lifecycle hooks |
| [Timestamps](./docs/timestamps.md) | Automatic `CreatedAt` / `UpdatedAt` |
| [Utility Types](./docs/types.md) | Pointer helpers |

HTML documentation is available at [`html_docs/index.html`](./html_docs/index.html).

## Geo Index Defaults Example

```go
ctx := context.Background()

err := mongorm.New(&ToDo{}).EnsureGeoDefaults(
    ctx,
    ToDoFields.Location,
    []bson.E{mongorm.Asc(ToDoFields.CreatedAt)},
)
if err != nil {
    panic(err)
}
```

## License

See [LICENSE](./LICENSE).
