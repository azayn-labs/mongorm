# Timestamps

MongORM can automatically manage `CreatedAt` and `UpdatedAt` fields for you. When enabled, `CreatedAt` is set once on the first insert, and `UpdatedAt` is updated to `time.Now()` on every `Save()` / `Update()` call.

## Enable Timestamps

### Via struct tags (automatic detection)

Add `timestamp:created_at` and / or `timestamp:updated_at` tags to your model fields. MongORM automatically enables timestamp management as soon as it detects any timestamp tag.

```go
import "time"

type ToDo struct {
    ID        *bson.ObjectID `bson:"_id,omitempty"       mongorm:"primary"`
    Text      *string        `bson:"text,omitempty"`
    CreatedAt *time.Time     `bson:"createdAt,omitempty" mongorm:"true,timestamp:created_at"`
    UpdatedAt *time.Time     `bson:"updatedAt,omitempty" mongorm:"true,timestamp:updated_at"`

    connectionString *string `mongorm:"mongodb://localhost:27017,connection:url"`
    database         *string `mongorm:"mydb,connection:database"`
    collection       *string `mongorm:"todos,connection:collection"`
}
```

The leading `true,` in `mongorm:"true,timestamp:created_at"` explicitly enables the timestamp feature. The feature is also auto-detected if any timestamp tag is present.

### Via MongORMOptions

```go
opts := &mongorm.MongORMOptions{
    Timestamps:     true,
    CollectionName: mongorm.String("todos"),
    DatabaseName:   mongorm.String("mydb"),
    MongoClient:    client,
}
orm := mongorm.FromOptions(&ToDo{}, opts)
```

## Behaviour

| Event | CreatedAt | UpdatedAt |
| --- | --- | --- |
| First `Save()` (insert) | Set to `time.Now()` | Set to `time.Now()` |
| Subsequent `Save()` (update) | Unchanged | Updated to `time.Now()` |
| `SaveMulti()` | Not managed | Not managed automatically |

> `SaveMulti()` does not apply automatic timestamps. Manage them manually if needed.

If your model defines only one timestamp field, MongORM manages that field independently:

- only `timestamp:created_at` → `CreatedAt` is set on insert and not changed afterwards.
- only `timestamp:updated_at` → `UpdatedAt` is refreshed on each `Save()` / `Update()`.

## Field Requirements

- Both fields must be of type `*time.Time`.
- The BSON tag should match what your application expects (e.g., `bson:"createdAt,omitempty"`).
- Only one of the two fields is required — you can have `CreatedAt` without `UpdatedAt` and vice versa.

## Protection

- `CreatedAt` is **never** overwritten after the initial insert, even if you include it in a `Set()` call.
- Timestamp fields are **never** removed by `Unset()`.

## Example

```go
package main

import (
    "context"
    "fmt"
    "time"

    "github.com/azayn-labs/mongorm"
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

func main() {
    ctx := context.Background()

    todo := &ToDo{Text: mongorm.String("Buy groceries")}
    orm  := mongorm.New(todo)

    // Insert
    if err := orm.Save(ctx); err != nil {
        panic(err)
    }
    fmt.Printf("Created at: %s\n", todo.CreatedAt)
    fmt.Printf("Updated at: %s\n", todo.UpdatedAt)

    // Update
    update := &ToDo{Text: mongorm.String("Buy organic groceries")}
    orm.Set(update)
    if err := orm.Save(ctx); err != nil {
        panic(err)
    }
    fmt.Printf("CreatedAt unchanged: %s\n", todo.CreatedAt)
    fmt.Printf("UpdatedAt changed:   %s\n", todo.UpdatedAt)
}
```

---

[Back to Documentation Index](./index.md) | [README](../README.md)
