# MongORM Documentation

> Published under **Azayn Labs** — Home: [www.azayn.com](https://www.azayn.com)
>
> Azayn Labs is currently an operating open-source brand and not yet a registered legal entity.
>
> Note: The homepage currently shows a coming-soon page.

MongORM is a production-ready, type-safe MongoDB ORM for Go focused on high developer velocity, clean data access patterns, and MongoDB-native flexibility.

This documentation covers everything from first setup to advanced patterns for Go + MongoDB applications: type-safe query building, CRUD operations, aggregation pipelines, transactions, hooks, indexes, and performance-friendly cursor workflows.

## Best For

- Go backend services using MongoDB as a primary datastore
- Teams that want typed query builders with minimal runtime surprises
- API and platform projects that need fast iteration plus production safety
- Developers who want an ORM-like DX without losing direct MongoDB control

## Menu

### Getting Started

- [Getting Started](./getting_started.md) — Installation, model definition, schema setup, and first steps
- [Configuration](./configuration.md) — Struct tag configuration, options struct, and mixed mode

### CRUD Operations

- [Creating Documents](./create.md) — Inserting new documents with `Save()`
- [Finding Documents](./find.md) — Querying with `First()` / `Find()`, typed projection decoding, count, and distinct
- [Updating Documents](./update.md) — Updating single or multiple documents
- [Deleting Documents](./delete.md) — Removing documents from a collection
- [Bulk Write](./bulk_write.md) — Executing batch insert/update/replace/delete operations
- [Indexes](./indexes.md) — Field-based index builders and geo index setup
- [Aggregation](./aggregate.md) — MongoDB aggregation pipelines with fluent stages and typed decoding
- [Cursors](./cursors.md) — Iterating over multiple results with `FindAll()`

### Query Building

- [Query Building](./query_building.md) — Type-safe filters, pagination helpers, projection, and update operators
- [Primitives](./primitives.md) — Type-safe field types (including geospatial) and their query methods

### Advanced

- [Hooks](./hooks.md) — Lifecycle hooks for all CRUD operations
- [Transactions](./transactions.md) — Execute operations in an atomic multi-step transaction
- [Errors](./errors.md) — Sentinel error taxonomy for consistent application handling
- [Timestamps](./timestamps.md) — Automatic `CreatedAt` / `UpdatedAt` management
- [Utility Types](./types.md) — Pointer helpers: `String()`, `Bool()`, `Int64()`, `Timestamp()`

## Discoverability Keywords

Go MongoDB ORM, type-safe query builder, Golang MongoDB library, MongoDB CRUD in Go, MongoDB aggregation Go, MongoDB transactions Go, MongoDB bulk operations Go, Go backend data layer.

---

[README](../README.md)
