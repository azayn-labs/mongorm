
# Getting Started

## Define a Model

You can define a model by creating a struct that follows the syntax as that of the official MongoDB Go driver. You can use struct tags to specify the JSON and BSON field names.

```go
package main

import "go.mongodb.org/mongo-driver/v2/bson"

type ToDo struct {
 ID   *bson.ObjectID `json:"id" bson:"_id,omitempty" mongorm:"primary"`
 Text *string        `json:"text" bson:"text,omitempty"`
}

```

> Note: Use the `mongorm:"primary"` tag to specify the primary key field. This is a mandatory requirement for all models. Refer [tags](./tags.md) for more details on the available struct tags.

## Defining the model schema

A model schema is a struct that defines the fields of the model. It is used to perform queries on the model. You can define a model schema by creating a struct that contains fields of type `*primitives.FieldType`.
The field names should match the field names of the model. You can find more details on the available field types in the [primitives documentation](./primitives.md).

```go
package main

import "github.com/CdTgr/mongorm/mongorm/primitives"

type ToDoSchema struct {
 ID   *primitives.ObjectIDField
 Text *primitives.StringField
}
```

## Getting the model fields

MongoFields are the core features of this library for querying. You will need to get the model fields to perform any query. You can get the model fields by using the `FieldsOf` function.
This function takes the model type and the schema type as type parameters and returns a struct containing the fields

_Definition:_

```go
mongorm.FieldsOf[M, S]()
```

_Example:_

```go
package main

import (
  "github.com/CdTgr/mongorm/mongorm"
  "github.com/CdTgr/mongorm/mongorm/primitives"
  "go.mongodb.org/mongo-driver/v2/bson"
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
```

---
---
[README](../README.md)
