# MongORM Documentation

MongORM is a lightweight, type-safe ORM library for MongoDB in Go. It provides a fluent API for defining models and performing CRUD operations using Go generics and struct tags.

## Menu

### Getting Started

- [Getting Started](./getting_started.md) — Installation, model definition, schema setup, and first steps
- [Configuration](./configuration.md) — Struct tag configuration, options struct, and mixed mode

### CRUD Operations

- [Creating Documents](./create.md) — Inserting new documents with `Save()`
- [Finding Documents](./find.md) — Querying single documents with `First()` / `Find()`
- [Updating Documents](./update.md) — Updating single or multiple documents
- [Deleting Documents](./delete.md) — Removing documents from a collection
- [Indexes](./indexes.md) — Field-based index builders and geo index setup
- [Aggregation](./aggregate.md) — Running aggregation pipelines with typed decoding
- [Cursors](./cursors.md) — Iterating over multiple results with `FindAll()`

### Query Building

- [Query Building](./query_building.md) — Using `Where()`, `WhereBy()`, `Set()`, and `Unset()`
- [Primitives](./primitives.md) — Type-safe field types (including geospatial) and their query methods

### Advanced

- [Hooks](./hooks.md) — Lifecycle hooks for all CRUD operations
- [Transactions](./transactions.md) — Execute operations in an atomic multi-step transaction
- [Timestamps](./timestamps.md) — Automatic `CreatedAt` / `UpdatedAt` management
- [Utility Types](./types.md) — Pointer helpers: `String()`, `Bool()`, `Int64()`, `Timestamp()`

---

[README](../README.md)
