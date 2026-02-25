# Basic Usage

```go
package main

import (
 "context"
 "fmt"

 "github.com/CdTgr/mongorm/mongorm"
 "github.com/CdTgr/mongorm/mongorm/primitives"
 "go.mongodb.org/mongo-driver/v2/bson"
 "go.mongodb.org/mongo-driver/v2/mongo"
 "go.mongodb.org/mongo-driver/v2/mongo/options"
)

type ToDo struct {
 ID   *bson.ObjectID `json:"id" bson:"_id,omitempty" mongorm:"primary"`
 Text *string        `json:"text" bson:"text,omitempty"`
}

type ToDoSchema struct {
 ID   *primitives.ObjectIDField
 Text *primitives.StringField
}

var ToDoFields = mongorm.FieldsOf[ToDo, ToDoSchema]()

func main() {
 dbURL := "mongodb://localhost:27017"
 // Get a MongoDB client
 client, err := mongo.Connect(
  options.Client().
   ApplyURI(dbURL).
   SetServerAPIOptions(options.ServerAPI(options.ServerAPIVersion1)),
 )
 if err != nil {
  panic(err)
 }

 dbName := "my-orm-test-db"

 // Create a new ToDo instance
 todo := &ToDo{
  Text: mongorm.String("Buy groceries"),
 }
 options := &mongorm.ModelOptions{
  CollectionName: "todo",
  Timestamps:     true,
  DB:             client.Database(dbName),
 }

 // Create a new model instance
 orm := mongorm.NewModel(todo, options)

 if err := orm.Save(context.TODO()); err != nil {
  panic(err.Error())
 }

 fmt.Printf("ToDo added %+v\n", todo)
}
```
