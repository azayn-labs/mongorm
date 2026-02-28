# Hooks

MongORM provides a lifecycle hook system for every CRUD operation. Hooks are implemented as methods on your model struct that satisfy Go interfaces. MongORM detects them automatically using interface assertions — no registration required.

## Overview

| Hook Interface | Method Signature | Called |
| --- | --- | --- |
| `BeforeCreateHook[T]` | `BeforeCreate(*MongORM[T]) error` | Before inserting a new document |
| `AfterCreateHook[T]` | `AfterCreate(*MongORM[T]) error` | After inserting a new document |
| `BeforeSaveHook[T]` | `BeforeSave(*MongORM[T], *bson.M) error` | Before save (second arg is filter, or nil for insert) |
| `AfterSaveHook[T]` | `AfterSave(*MongORM[T]) error` | After save (insert or update) |
| `BeforeUpdateHook[T]` | `BeforeUpdate(*MongORM[T], *bson.M, *bson.M) error` | Before updating (filter, update doc) |
| `AfterUpdateHook[T]` | `AfterUpdate(*MongORM[T]) error` | After updating |
| `BeforeFindHook[T]` | `BeforeFind(*MongORM[T], *bson.M) error` | Before finding a document |
| `AfterFindHook[T]` | `AfterFind(*MongORM[T]) error` | After finding a document |
| `BeforeDeleteHook[T]` | `BeforeDelete(*MongORM[T], *bson.M) error` | Before deleting (filter doc) |
| `AfterDeleteHook[T]` | `AfterDelete(*MongORM[T]) error` | After deleting |
| `BeforeFinalizeHook[T]` | `BeforeFinalize(*MongORM[T]) error` | Before applying a fetched document to the schema |
| `AfterFinalizeHook[T]` | `AfterFinalize(*MongORM[T]) error` | After applying a fetched document to the schema |

## Change tracking inside hooks

During `BeforeSave`, `BeforeUpdate`, and `BeforeCreate`, you can inspect changed fields:

- `m.IsModified(UserFields.Email)`
- `m.IsModified(UserFields.Auth.Provider)`
- `m.IsModified(mongorm.FieldPath(mongorm.RawField("devices"), "token"))`
- `m.ModifiedFields()` (`[]mongorm.Field`)
- `m.ModifiedValue(UserFields.Email)` (`oldValue, newValue, ok`)

Nested paths are supported.

```go
func (u *User) BeforeUpdate(m *mongorm.MongORM[User], filter *bson.M, update *bson.M) error {
    if m.IsModified(UserFields.Auth.Provider) {
        // do something when nested provider changed
    }
    return nil
}
```

### Better end-to-end example

```go
type AuthProfile struct {
    Provider *string  `bson:"provider,omitempty"`
    Scopes   []string `bson:"scopes,omitempty"`
}

type User struct {
    ID       *bson.ObjectID `bson:"_id,omitempty" mongorm:"primary"`
    Email    *string        `bson:"email,omitempty"`
    Password *string        `bson:"password,omitempty"`
    Auth     *AuthProfile   `bson:"auth,omitempty"`
}

type AuthProfileSchema struct {
    Provider *primitives.StringField
    Scopes   *primitives.StringArrayField
}

type UserSchema struct {
    ID       *primitives.ObjectIDField
    Email    *primitives.StringField
    Password *primitives.StringField
    Auth     *AuthProfileSchema
}

var UserFields = mongorm.FieldsOf[User, UserSchema]()

func (u *User) BeforeSave(m *mongorm.MongORM[User], filter *bson.M) error {
    // Works for both create (filter=nil) and update (filter!=nil)
    if m.IsModified(UserFields.Password) {
        // hash password, rotate sessions, etc.
    }

    if m.IsModified(UserFields.Auth.Provider) || m.IsModified(UserFields.Auth.Scopes) {
        // re-sync OAuth metadata
    }

    // For arrays/slices of struct, check parent or deep path
    devicesField := mongorm.RawField("devices")
    if m.IsModified(devicesField) || m.IsModified(mongorm.FieldPath(devicesField, "token")) {
        // revoke old device tokens
    }

    // Optional: inspect all changed paths
    changed := m.ModifiedFields()
    _ = changed

    oldEmail, newEmail, emailChanged := m.ModifiedValue(UserFields.Email)
    _ = oldEmail
    _ = newEmail
    _ = emailChanged

    return nil
}

func (u *User) BeforeUpdate(m *mongorm.MongORM[User], filter *bson.M, update *bson.M) error {
    if m.IsModified(UserFields.Email) {
        // verify new email / update search index
    }

    return nil
}
```

With this model, changes like `auth.provider` match both `m.IsModified(UserFields.Auth.Provider)` and `m.IsModified(mongorm.RawField("auth"))`.

## Implementation

Implement any combination of hooks by adding the corresponding methods to your model struct. You only need to implement the hooks you care about.

```go
package main

import (
    "fmt"

    "github.com/azayn-labs/mongorm"
    "go.mongodb.org/mongo-driver/v2/bson"
)

type ToDo struct {
    ID   *bson.ObjectID `bson:"_id,omitempty" mongorm:"primary"`
    Text *string        `bson:"text,omitempty"`

    connectionString *string `mongorm:"mongodb://localhost:27017,connection:url"`
    database         *string `mongorm:"mydb,connection:database"`
    collection       *string `mongorm:"todos,connection:collection"`
}

// --- Find hooks ---

func (t *ToDo) BeforeFind(m *mongorm.MongORM[ToDo], query *bson.M) error {
    fmt.Printf("[HOOK] BeforeFind — query: %+v\n", *query)
    return nil
}

func (t *ToDo) AfterFind(m *mongorm.MongORM[ToDo]) error {
    fmt.Println("[HOOK] AfterFind")
    return nil
}

// --- Create hooks ---

func (t *ToDo) BeforeCreate(m *mongorm.MongORM[ToDo]) error {
    fmt.Println("[HOOK] BeforeCreate")
    return nil
}

func (t *ToDo) AfterCreate(m *mongorm.MongORM[ToDo]) error {
    fmt.Println("[HOOK] AfterCreate")
    return nil
}

// --- Save hooks (covers both insert and update) ---

func (t *ToDo) BeforeSave(m *mongorm.MongORM[ToDo], filter *bson.M) error {
    if filter != nil {
        fmt.Printf("[HOOK] BeforeSave (update) — filter: %+v\n", *filter)
    } else {
        fmt.Println("[HOOK] BeforeSave (insert)")
    }
    return nil
}

func (t *ToDo) AfterSave(m *mongorm.MongORM[ToDo]) error {
    fmt.Println("[HOOK] AfterSave")
    return nil
}

// --- Update hooks ---

func (t *ToDo) BeforeUpdate(m *mongorm.MongORM[ToDo], filter *bson.M, update *bson.M) error {
    fmt.Printf("[HOOK] BeforeUpdate — filter: %+v, update: %+v\n", *filter, *update)
    return nil
}

func (t *ToDo) AfterUpdate(m *mongorm.MongORM[ToDo]) error {
    fmt.Println("[HOOK] AfterUpdate")
    return nil
}

// --- Delete hooks ---

func (t *ToDo) BeforeDelete(m *mongorm.MongORM[ToDo], filter *bson.M) error {
    fmt.Printf("[HOOK] BeforeDelete — filter: %+v\n", *filter)
    return nil
}

func (t *ToDo) AfterDelete(m *mongorm.MongORM[ToDo]) error {
    fmt.Println("[HOOK] AfterDelete")
    return nil
}

// --- Finalize hooks (called when a document is decoded from the DB) ---

func (t *ToDo) BeforeFinalize(m *mongorm.MongORM[ToDo]) error {
    fmt.Println("[HOOK] BeforeFinalize")
    return nil
}

func (t *ToDo) AfterFinalize(m *mongorm.MongORM[ToDo]) error {
    fmt.Println("[HOOK] AfterFinalize")
    return nil
}
```

## Hook Execution Order

### Save (insert)

1. `BeforeSave` (filter is `nil`)
2. `BeforeCreate`
3. Document inserted into MongoDB
4. `AfterCreate`
5. `AfterSave`

### Save (update)

1. `BeforeSave` (filter is the query document)
2. `BeforeUpdate` (filter + update documents)
3. Document updated in MongoDB
4. `AfterUpdate`
5. `AfterSave`

### First / Find

1. `BeforeFind` (query document)
2. Document fetched from MongoDB
3. `BeforeFinalize`
4. Document decoded into schema
5. `AfterFinalize`
6. `AfterFind`

### Delete

1. `BeforeDelete` (filter document)
2. Document deleted from MongoDB
3. `AfterDelete`

## Returning Errors from Hooks

If a hook returns a non-nil error, the operation is aborted and the error is propagated to the caller:

```go
func (t *ToDo) BeforeCreate(m *mongorm.MongORM[ToDo]) error {
    if t.Text == nil || *t.Text == "" {
        return fmt.Errorf("text field is required")
    }
    return nil
}
```

---

[Back to Documentation Index](./index.md) | [README](../README.md)
