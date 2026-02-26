# Deleting Documents

Use `Delete()` to remove a single document that matches the current filters.

## Basic Delete

Filter to the target document using `Where()`, then call `Delete()`:

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

    targetID := bson.NewObjectID() // use a real existing ID

    todo := &ToDo{}
    orm  := mongorm.New(todo)
    orm.Where(ToDoFields.ID.Eq(targetID))

    if err := orm.Delete(ctx); err != nil {
        panic(err)
    }

    fmt.Println("Document deleted")
}
```

## Delete Using an Already-Loaded Document

If you already have a document with a primary key populated (e.g., after a `First()` call), you can delete it without calling `Where()`:

```go
// Load the document first
todo := &ToDo{}
orm  := mongorm.New(todo)
orm.Where(ToDoFields.Text.Eq("Buy groceries"))

if err := orm.First(ctx); err != nil {
    panic(err)
}

// Now delete it — MongORM uses todo.ID as the filter automatically
if err := orm.Delete(ctx); err != nil {
    panic(err)
}
```

## Hooks

Before and after a delete, the `BeforeDeleteHook` and `AfterDeleteHook` are invoked if implemented on your model. See [Hooks](./hooks.md).

```go
func (t *ToDo) BeforeDelete(m *mongorm.MongORM[ToDo], filter *bson.M) error {
    fmt.Printf("Deleting document with filter: %+v\n", *filter)
    return nil
}

func (t *ToDo) AfterDelete(m *mongorm.MongORM[ToDo]) error {
    fmt.Println("Document deleted")
    return nil
}
```

## Notes

- `Delete()` only removes **one** document (the first match).
- There is no `DeleteMulti()` at this time. To delete multiple documents, iterate with `FindAll()` and call `Delete()` on each.

---

[Back to Documentation Index](./index.md) | [README](../README.md)
