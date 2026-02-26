# Configuration

MongORM supports three flexible configuration modes so teams can choose static, runtime, or hybrid connection management patterns. All modes can be combined, making it easy to evolve from local development setups to production deployment environments.

## Mode A — Struct Tags

Embed private fields with `mongorm` connection tags directly in your model struct. This is the most self-contained approach because the connection details travel with the type definition.

```go
type ToDo struct {
    // Data fields
    ID   *bson.ObjectID `bson:"_id,omitempty" mongorm:"primary"`
    Text *string        `bson:"text,omitempty"`

    // Connection tags (private fields, values are the defaults)
    connectionString *string `mongorm:"mongodb://localhost:27017,connection:url"`
    database         *string `mongorm:"mydb,connection:database"`
    collection       *string `mongorm:"todos,connection:collection"`
}

orm := mongorm.New(&ToDo{})
```

Connection tag format: `mongorm:"<value>,<tag>"`.

| Tag | Purpose |
| --- | --- |
| `connection:url` | MongoDB connection string |
| `connection:database` | Database name |
| `connection:collection` | Collection name |

MongORM caches the MongoDB client by connection string, so multiple instances reuse the same underlying driver connection.

## Mode B — MongORMOptions

Pass a `*mongorm.MongORMOptions` struct to `FromOptions()`. This is useful when connection details are loaded from environment variables or config files at runtime.

```go
import (
    "github.com/CdTgr/mongorm"
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

orm := mongorm.FromOptions(&ToDo{}, opts)
```

You can also use `mongorm.NewClient(...)` as a convenience helper:

```go
client, err := mongorm.NewClient("mongodb://localhost:27017")
if err != nil {
    panic(err)
}

orm := mongorm.FromOptions(&ToDo{}, &mongorm.MongORMOptions{
    MongoClient:    client,
    DatabaseName:   mongorm.String("mydb"),
    CollectionName: mongorm.String("todos"),
})
```

### MongORMOptions fields

| Field | Type | Description |
| --- | --- | --- |
| `Timestamps` | `bool` | Enable automatic `CreatedAt` / `UpdatedAt` management |
| `CollectionName` | `*string` | MongoDB collection name |
| `DatabaseName` | `*string` | MongoDB database name |
| `MongoClient` | `*mongo.Client` | Pre-configured MongoDB client |

## Mode C — Mixed

Struct tags and `MongORMOptions` can be combined. `MongORMOptions` values take precedence when both are present.

```go
// Model provides collection name and connection string via tags
type ToDo struct {
    ID   *bson.ObjectID `bson:"_id,omitempty" mongorm:"primary"`
    Text *string        `bson:"text,omitempty"`

    connectionString *string `mongorm:"mongodb://localhost:27017,connection:url"`
    collection       *string `mongorm:"todos,connection:collection"`
    // No database tag — will be provided via options
}

opts := &mongorm.MongORMOptions{
    Timestamps:   true,
    DatabaseName: mongorm.String("mydb"), // overrides struct tag (if any)
}

orm := mongorm.FromOptions(&ToDo{}, opts)
```

## Field Tags

Field-level `mongorm` tags control how individual struct fields are handled during operations.

| Tag | Purpose |
| --- | --- |
| `primary` | Marks the primary key field (required). Used to determine insert vs update. |
| `version` | Enables optimistic locking on the field (typically `_version`). |
| `readonly` | Field is never written during `Set()` or `Unset()` operations. |
| `timestamp:created_at` | Field receives the insert timestamp and is never updated after that. |
| `timestamp:updated_at` | Field is updated to `time.Now()` on every `Save()` call. |

```go
type ToDo struct {
    ID        *bson.ObjectID `bson:"_id,omitempty"      mongorm:"primary"`
    ReadOnlyField *string    `bson:"readOnly,omitempty" mongorm:"readonly"`
    CreatedAt *time.Time     `bson:"createdAt,omitempty" mongorm:"true,timestamp:created_at"`
    UpdatedAt *time.Time     `bson:"updatedAt,omitempty" mongorm:"true,timestamp:updated_at"`
}
```

> The `true` value before `timestamp:created_at` / `timestamp:updated_at` enables timestamps automatically. Alternatively, set `MongORMOptions.Timestamps = true`.

See [Timestamps](./timestamps.md) for more details.

---

[Back to Documentation Index](./index.md) | [README](../README.md)
