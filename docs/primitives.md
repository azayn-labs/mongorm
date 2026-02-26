# MongORM Primitives

Primitives are type-safe field wrappers used in your schema struct. Each field type maps to a Go type and exposes query-builder methods that return `bson.M` filters, ready to be passed to `Where()`.

## Importing

```go
import "github.com/CdTgr/mongorm/primitives"
```

## Available Field Types

| Field Type | Go Type | Use for |
| --- | --- | --- |
| `ObjectIDField` | `bson.ObjectID` | MongoDB `_id` and foreign key fields |
| `StringField` | `string` | Text fields |
| `Int64Field` | `int64` | Integer numeric fields (also handles int32, int8, int) |
| `Float64Field` | `float64` | Floating-point fields (also handles float32) |
| `BoolField` | `bool` | Boolean fields |
| `TimestampField` | `time.Time` | Date/time fields |
| `GeoField` | `mongorm.GeoPoint` / `mongorm.GeoLineString` / `mongorm.GeoPolygon` | Geospatial fields |
| `GenericField` | any | Fallback for unmapped types (only `BSONName()` available) |

---

## ObjectIDField

**Package:** `primitives.ObjectIDField`

```go
type ToDoSchema struct {
    ID *primitives.ObjectIDField
}
var ToDoFields = mongorm.FieldsOf[ToDo, ToDoSchema]()
```

### ObjectID Methods

| Method | MongoDB operator | Description |
| --- | --- | --- |
| `Eq(v bson.ObjectID)` | `$eq` | Equals |
| `Ne(v bson.ObjectID)` | `$ne` | Not equals |
| `In(v ...bson.ObjectID)` | `$in` | In a list |
| `Nin(v ...bson.ObjectID)` | `$nin` | Not in a list |
| `Gt(v bson.ObjectID)` | `$gt` | Greater than |
| `Gte(v bson.ObjectID)` | `$gte` | Greater than or equal |
| `Lt(v bson.ObjectID)` | `$lt` | Less than |
| `Lte(v bson.ObjectID)` | `$lte` | Less than or equal |
| `Exists()` | `$exists: true` | Field exists |
| `NotExists()` | `$exists: false` | Field does not exist |
| `IsNull()` | `$eq: null` | Field is null |
| `IsNotNull()` | `$ne: null` | Field is not null |

```go
orm.Where(ToDoFields.ID.Eq(someID))
orm.Where(ToDoFields.ID.In(id1, id2, id3))
```

---

## StringField

**Package:** `primitives.StringField`

### String Methods

| Method | MongoDB operator | Description |
| --- | --- | --- |
| `Eq(v string)` | `$eq` | Equals |
| `Ne(v string)` | `$ne` | Not equals |
| `Reg(pattern string)` | `$regex` | Regex match |
| `In(v ...string)` | `$in` | In a list |
| `Nin(v ...string)` | `$nin` | Not in a list |
| `Exists()` | `$exists: true` | Field exists |
| `NotExists()` | `$exists: false` | Field does not exist |
| `IsNull()` | `$eq: null` | Field is null |
| `IsNotNull()` | `$ne: null` | Field is not null |

```go
orm.Where(ToDoFields.Text.Eq("Buy groceries"))
orm.Where(ToDoFields.Text.Reg("groceries$")) // ends with "groceries"
```

---

## Int64Field

**Package:** `primitives.Int64Field`

Handles `int64`, `int32`, `int8`, and `int` model fields.

### Int64 Methods

| Method | MongoDB operator | Description |
| --- | --- | --- |
| `Eq(v int64)` | `$eq` | Equals |
| `Ne(v int64)` | `$ne` | Not equals |
| `In(v ...int64)` | `$in` | In a list |
| `Nin(v ...int64)` | `$nin` | Not in a list |
| `Gt(v int64)` | `$gt` | Greater than |
| `Gte(v int64)` | `$gte` | Greater than or equal |
| `Lt(v int64)` | `$lt` | Less than |
| `Lte(v int64)` | `$lte` | Less than or equal |
| `Exists()` | `$exists: true` | Field exists |
| `NotExists()` | `$exists: false` | Field does not exist |
| `IsNull()` | `$eq: null` | Field is null |
| `IsNotNull()` | `$ne: null` | Field is not null |

```go
orm.Where(TaskFields.Priority.Gte(3))
orm.Where(TaskFields.Count.In(1, 2, 5))
```

---

## Float64Field

**Package:** `primitives.Float64Field`

Handles `float64` and `float32` model fields.

### Float64 Methods

| Method | MongoDB operator | Description |
| --- | --- | --- |
| `Eq(v float64)` | `$eq` | Equals |
| `Ne(v float64)` | `$ne` | Not equals |
| `In(v ...float64)` | `$in` | In a list |
| `Nin(v ...float64)` | `$nin` | Not in a list |
| `Gt(v float64)` | `$gt` | Greater than |
| `Gte(v float64)` | `$gte` | Greater than or equal |
| `Lt(v float64)` | `$lt` | Less than |
| `Lte(v float64)` | `$lte` | Less than or equal |
| `Exists()` | `$exists: true` | Field exists |
| `NotExists()` | `$exists: false` | Field does not exist |
| `IsNull()` | `$eq: null` | Field is null |
| `IsNotNull()` | `$ne: null` | Field is not null |

```go
orm.Where(ProductFields.Price.Lt(99.99))
```

---

## BoolField

**Package:** `primitives.BoolField`

### Bool Methods

| Method | MongoDB operator | Description |
| --- | --- | --- |
| `Eq(v bool)` | `$eq` | Equals |
| `Ne(v bool)` | `$ne` | Not equals |
| `In(v ...bool)` | `$in` | In a list |
| `Nin(v ...bool)` | `$nin` | Not in a list |
| `Exists()` | `$exists: true` | Field exists |
| `NotExists()` | `$exists: false` | Field does not exist |
| `IsNull()` | `$eq: null` | Field is null |
| `IsNotNull()` | `$ne: null` | Field is not null |

```go
orm.Where(ToDoFields.Done.Eq(false))
```

---

## TimestampField

**Package:** `primitives.TimestampField`

### Timestamp Methods

| Method | MongoDB operator | Description |
| --- | --- | --- |
| `Eq(v time.Time)` | `$eq` | Equals |
| `Ne(v time.Time)` | `$ne` | Not equals |
| `In(v ...time.Time)` | `$in` | In a list |
| `Nin(v ...time.Time)` | `$nin` | Not in a list |
| `Gt(v time.Time)` | `$gt` | After |
| `Gte(v time.Time)` | `$gte` | After or at |
| `Lt(v time.Time)` | `$lt` | Before |
| `Lte(v time.Time)` | `$lte` | Before or at |
| `Exists()` | `$exists: true` | Field exists |
| `NotExists()` | `$exists: false` | Field does not exist |
| `IsNull()` | `$eq: null` | Field is null |
| `IsNotNull()` | `$ne: null` | Field is not null |

```go
cutoff := time.Now().Add(-24 * time.Hour)
orm.Where(ToDoFields.CreatedAt.Gte(cutoff))
```

---

## GeoField

**Package:** `primitives.GeoField`

Use `GeoField` for geospatial queries.

### Supported model types

- `*mongorm.GeoPoint`
- `*mongorm.GeoLineString`
- `*mongorm.GeoPolygon`

### Geo Methods

| Method | MongoDB operator | Description |
| --- | --- | --- |
| `Eq(v any)` | `$eq` | Equals geometry |
| `Ne(v any)` | `$ne` | Not equals geometry |
| `Near(geometry any)` | `$near` | Near geometry |
| `NearWithDistance(geometry, min, max)` | `$near` | Near with min/max distance |
| `NearSphere(geometry any)` | `$nearSphere` | Spherical near |
| `NearSphereWithDistance(geometry, min, max)` | `$nearSphere` | Spherical near with min/max distance |
| `Within(geometry any)` | `$geoWithin` | Within geometry |
| `WithinBox(bottomLeft, upperRight)` | `$geoWithin/$box` | Within box |
| `WithinCenter(center, radius)` | `$geoWithin/$center` | Within flat circle |
| `WithinCenterSphere(center, radius)` | `$geoWithin/$centerSphere` | Within spherical circle |
| `Intersects(geometry any)` | `$geoIntersects` | Geometry intersection |
| `Exists()` | `$exists: true` | Field exists |
| `NotExists()` | `$exists: false` | Field does not exist |
| `IsNull()` | `$eq: null` | Field is null |
| `IsNotNull()` | `$ne: null` | Field is not null |

```go
point := mongorm.NewGeoPoint(12.9716, 77.5946)

orm.Where(ToDoFields.Location.Near(point))

cityBounds := mongorm.NewGeoPolygon(
    [][]float64{{77.4, 12.8}, {77.8, 12.8}, {77.8, 13.1}, {77.4, 13.1}, {77.4, 12.8}},
)
orm.Where(ToDoFields.Location.Within(cityBounds))
```

---

## GenericField

**Package:** `primitives.GenericField`

Fallback for any model field type that does not map to one of the specific primitives above. Only provides the `BSONName()` method; use raw `bson.M` with `Where()` for custom queries.

```go
orm.Where(bson.M{"myField": bson.M{"$exists": true}})
```

---

[Back to Documentation Index](./index.md) | [README](../README.md)
