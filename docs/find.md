# Finding a document

Finding a document is as simple as calling the `find` method on the model. You can find a single document by its ID or by any other field.
Here is where the ORM really shines, as it allows you to use the same syntax for almost all the operations, making it very intuitive to use.

## Basic Usage

```go
...
...

// Create a new ToDo instance
todo := &ToDo{}

// Create a new model instance
orm := mongorm.NewModel(todo, options)

// Use your ObjectID reference
myId := bson.NewObjectID() // NewObjectID() is just an example.

// Query the database using the ORM
orm.Where(ToDoFields.ID.Eq(myId))

if err := orm.First(context.TODO()); err != nil {
  panic(err)
}

fmt.Printf("Found document %+v\n", todo)
```

All the primitive operations can be found in the [Primitives](./primitives.md) section.
