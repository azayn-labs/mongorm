# Indexes

MongORM provides field-based index helpers so users can create indexes without hardcoding BSON field names.

## Quick Start

```go
orm := mongorm.New(&ToDo{})

// 2dsphere index for geo queries
_, err := orm.Ensure2DSphereIndex(ctx, ToDoFields.Location)
if err != nil {
    panic(err)
}

// one-call geo defaults: 2dsphere + optional support index
_, err = orm.EnsureGeoDefaults(
    ctx,
    ToDoFields.Location,
    []bson.E{mongorm.Asc(ToDoFields.CreatedAt)},
)
if err != nil {
    panic(err)
}

// named compound index
model := mongorm.NamedIndexModelFromKeys(
    "todo_text_count_idx",
    mongorm.Asc(ToDoFields.Text),
    mongorm.Desc(ToDoFields.Count),
)

_, err = orm.EnsureIndex(ctx, model)
if err != nil {
    panic(err)
}
```

## Key Builders

Use these helpers to avoid hardcoded field names:

- `Asc(field)`
- `Desc(field)`
- `Text(field)`
- `Geo2DSphere(field)`
- `Geo2D(field)`

Each helper uses `field.BSONName()` internally.

## Index Model Builders

- `IndexModelFromKeys(keys...)`
- `UniqueIndexModelFromKeys(keys...)`
- `NamedIndexModelFromKeys(name, keys...)`

## Execution Methods

- `EnsureIndex(ctx, model)`
- `EnsureIndexes(ctx, models)`
- `Ensure2DSphereIndex(ctx, field)`
- `EnsureGeoDefaults(ctx, geoField, supportingKeys)`

## Create Multiple Indexes

```go
models := []mongo.IndexModel{
    mongorm.IndexModelFromKeys(mongorm.Text(ToDoFields.Text)),
    mongorm.UniqueIndexModelFromKeys(mongorm.Asc(ToDoFields.ID)),
}

names, err := orm.EnsureIndexes(ctx, models)
if err != nil {
    panic(err)
}
fmt.Println(names)
```

---

[Back to Documentation Index](./index.md) | [README](../README.md)
