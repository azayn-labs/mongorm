# MongORM

> Open source by **Azayn Labs** — Home: [azayn.com](https://azayn.com)

MongORM is a production-ready, type-safe MongoDB ORM for Go that combines generics, fluent query building, and predictable data modeling. It helps Go teams ship faster with clean CRUD workflows, typed filters, aggregation support, and transaction-safe operations.

If you are building backend APIs, SaaS products, internal tools, or data-heavy services on MongoDB, MongORM gives you a high-signal developer experience without hiding core MongoDB power.

## Features

- Type-safe model and schema definitions using Go generics
- Fluent API for building queries with `Where()`, `Sort()`, `Limit()`, `Skip()`, `Projection()`, keyset pagination helpers, `Set()`, and `Unset()`
- Typed projection decoding for DTO targets via `FindOneAs[T, R]()` and `FindAllAs[T, R]()`
- Full CRUD support: create, find, update (single and multi), and delete (single and multi)
- Bulk write support with typed builder helpers via `BulkWrite()`, `BulkWriteInTransaction()`, and `NewBulkWriteBuilder[T]()`
- Aggregation support: raw pipelines and fluent stage builder via `Aggregate()`, `AggregateRaw()`, `AggregateAs[T, R]()`, and `AggregatePipeline()`
- Query utilities: `Count()`, `Distinct()`, `DistinctFieldAs[T, V]()`, `DistinctStrings()`, `DistinctInt64()`, `DistinctBool()`, `DistinctFloat64()`, `DistinctObjectIDs()`, and `DistinctTimes()`
- Geospatial support: `GeoField` with `Near`, `Within`, and `Intersects` query helpers
- Index support: field-driven builders, `Ensure2DSphereIndex()`, and `EnsureGeoDefaults()`
- Transactions: `WithTransaction()` for atomic multi-operation workflows
- Optimistic locking via `mongorm:"version"` (`_version`) and `ErrOptimisticLockConflict`
- Error taxonomy with sentinels: `ErrNotFound`, `ErrDuplicateKey`, `ErrInvalidConfig`, `ErrTransactionUnsupported`
- Lifecycle hooks for every operation (Before/After Create, Save, Update, Find, Delete, Finalize)
- Automatic `CreatedAt` / `UpdatedAt` timestamp management
- Flexible configuration: struct tags, options struct, or both
- Cursor-based iteration for large result sets
- Connection pooling — clients are reused by connection string
- Lightweight with a single dependency: the official MongoDB Go driver

## Why teams choose MongORM

- **Type-safe by default:** schema primitives and generics reduce runtime query mistakes.
- **Fast developer velocity:** fluent APIs for filtering, updates, aggregation, and bulk workflows.
- **Production-focused reliability:** transactions, optimistic locking, timestamps, hooks, and clear error taxonomy.
- **MongoDB-native flexibility:** raw BSON compatibility when you need full control.
- **Clean architecture fit:** works naturally with service layers, repository patterns, and domain models.

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
| [Bulk Write](./docs/bulk_write.md) | Batch insert/update/replace/delete operations |
| [Indexes](./docs/indexes.md) | Field-based index builders and geo index setup |
| [Aggregation](./docs/aggregate.md) | Aggregation pipelines with fluent builder and typed decoding |
| [Cursors](./docs/cursors.md) | Iterating with `FindAll()` |
| [Query Building](./docs/query_building.md) | `Where()`, find modifiers, pagination helpers, `Set()`, `Unset()` |
| [Primitives](./docs/primitives.md) | Type-safe field query methods |
| [Hooks](./docs/hooks.md) | Lifecycle hooks |
| [Transactions](./docs/transactions.md) | Running operations in MongoDB transactions |
| [Errors](./docs/errors.md) | Sentinel errors and handling patterns |
| [Timestamps](./docs/timestamps.md) | Automatic `CreatedAt` / `UpdatedAt` |
| [Utility Types](./docs/types.md) | Pointer helpers |

HTML documentation is available at [`html_docs/index.html`](./html_docs/index.html).

## Keywords

Go MongoDB ORM, Golang MongoDB ORM, type-safe MongoDB query builder for Go, Go generics ORM, MongoDB CRUD library for Go, MongoDB transactions in Go, MongoDB aggregation in Go, MongoDB bulk write Go, MongoDB hooks and timestamps, lightweight Go ORM.

## GitHub SEO Pack

Use these when publishing the repository to maximize discoverability in GitHub search and community feeds.

### Repository Description (pick one)

- Production-ready, type-safe MongoDB ORM for Go with generics, fluent query building, transactions, aggregation, and bulk operations.
- Type-safe MongoDB ORM for Go: fluent CRUD, typed filters, aggregation pipelines, bulk writes, hooks, and transaction support.
- Lightweight Golang MongoDB ORM with generics, query builder, optimistic locking, timestamps, and production-friendly data workflows.

### Suggested GitHub Topics

`go`, `golang`, `mongodb`, `orm`, `mongo-orm`, `golang-library`, `backend`, `query-builder`, `type-safe`, `generics`, `aggregation`, `transactions`, `bulk-write`, `developer-tools`, `data-access`

### Launch Checklist

- Set repository description using one of the options above.
- Add the suggested topic tags in GitHub repository settings.
- Publish `v1.0.0` as the first stable release using `.github/RELEASE_TEMPLATE.md`.
- Pin a concise usage example in the release body (create + find + update).
- Share release in Go and MongoDB communities with keywords from this README.

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
